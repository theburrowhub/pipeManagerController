// Package normalize provides functionality to normalize the pipeline data.
// It adds the necessary steps to the tasks to:
// - clone the repository
// - download and upload the artifacts and cache
// - expand each batch task in the pipeline
// It also adds the necessary finish tasks to:
// - launch the next pipeline in the chain
package normalize

import (
	"github.com/go-logr/logr"
	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

const (
	defaultShell     = "/bin/sh"       // Default shell for the automated steps
	defaultShellSets = "set -e"        // Default shell sets for the automated steps
	launcherBinary   = "/app/launcher" // Default launcher binary path

	workspaceDir = "/workspaceDir" // Default workspace directory for the all steps
)

// Normalize normalizes the pipeline data
// It adds the necessary steps to the tasks to:
// - clone the repository
// - download and upload the artifacts and cache
// - expand each batch task in the pipeline
// It also adds the necessary finish tasks to:
// - launch the next pipeline in the chain
func Normalize(logger logr.Logger, pipeline v1alpha1.PipelineSpec) (v1alpha1.PipelineSpec, error) {
	repository := pipeline.Params["REPOSITORY"]
	commit := pipeline.Params["COMMIT"]

	// Loop tasks for:
	// - cloning the repository
	// - download and upload the artifacts and cache
	// - expand each batch task in the pipeline
	for taskName, taskData := range pipeline.Tasks {
		// Process the task data to add the necessary steps
		taskData = processTask(logger, pipeline, taskName, taskData, repository, commit)

		// Explode batch tasks if they are defined
		if taskData.Batch != nil {
			// Explode the batch task
			tasks := processBatchTask(taskName, taskData)
			if tasks != nil {
				// Add all the new tasks to the pipeline
				for name, tData := range tasks {
					pipeline.Tasks[name] = tData
				}
				// Remove the original task from the pipeline
				delete(pipeline.Tasks, taskName)
			}
		} else {
			// Replace the task with the new data in the pipeline
			pipeline.Tasks[taskName] = taskData
		}

		// Normalize Fail tasks
		for FailTaskName, FailTaskData := range pipeline.FinishTasks.Fail {
			// Process the task data to add the necessary steps
			FailTaskData = processTask(logger, pipeline, FailTaskName, FailTaskData, repository, commit)

			// Explode batch tasks if they are defined
			if FailTaskData.Batch != nil {
				// Explode the batch task
				tasks := processBatchTask(FailTaskName, FailTaskData)
				if tasks != nil {
					// Add all the new tasks to the pipeline
					for name, tData := range tasks {
						pipeline.FinishTasks.Fail[name] = tData
					}
					// Remove the original task from the pipeline
					delete(pipeline.Tasks, FailTaskName)
				}
			} else {
				// Replace the task with the new data in the pipeline
				pipeline.FinishTasks.Fail[FailTaskName] = FailTaskData
			}

		}

		// Loop through the list of rawPipelines to launch when the current pipeline fails
		for _, launchPipelineName := range pipeline.Launch.WhenFail {
			failTaskName := k8sObjectName("launch", launchPipelineName)
			pipeline.FinishTasks.Fail[failTaskName] = defineLaunchPipelineTask(pipeline, repository, commit, launchPipelineName)
		}

		// Normalize Success tasks
		for successTaskName, SuccessTaskData := range pipeline.FinishTasks.Success {
			// Process the task data to add the necessary steps
			SuccessTaskData = processTask(logger, pipeline, successTaskName, SuccessTaskData, repository, commit)

			// Explode batch tasks if they are defined
			if SuccessTaskData.Batch != nil {
				// Explode the batch task
				tasks := processBatchTask(successTaskName, SuccessTaskData)
				if tasks != nil {
					// Add all the new tasks to the pipeline
					for name, tData := range tasks {
						pipeline.FinishTasks.Success[name] = tData
					}
					// Remove the original task from the pipeline
					delete(pipeline.Tasks, successTaskName)
				}
			} else {
				// Replace the task with the new data in the pipeline
				pipeline.FinishTasks.Success[successTaskName] = SuccessTaskData
			}

		}

		// Loop through the list of rawPipelines to launch when the current pipeline finishes successfully
		for _, launchPipelineName := range pipeline.Launch.WhenSuccess {
			successTaskName := k8sObjectName("launch", launchPipelineName)
			pipeline.FinishTasks.Success[successTaskName] = defineLaunchPipelineTask(pipeline, repository, commit, launchPipelineName)
		}

		// Clean
		// Remove unnecessary cloneRepository and launchPipeline from the pipeline
		pipeline.CloneRepository = v1alpha1.CloneRepositoryConfig{}
		pipeline.Launch = v1alpha1.Launch{}
	}

	return pipeline, nil
}

// processTask processes the task data to add the necessary steps to:
// - clone the repository
// - download and upload the artifacts and cache
// - expand each batch task in the pipeline
func processTask(logger logr.Logger, pipe v1alpha1.PipelineSpec, taskName string, taskData v1alpha1.Task, repository, commit string) v1alpha1.Task {
	logger.Info("Normalizing task", "taskName", taskName)

	// This is the list of steps that will be added to the task at the beginning
	var firstSteps []v1alpha1.Step

	// Add the clone repository step if it is defined as true in the pipe or in the task itself
	cloneRepository := pipe.CloneRepository.Enable || taskData.CloneRepository.Enable
	if cloneRepository {
		logger.Info("Adding clone repository step", "taskName", taskName)
		cloneRepositoryStep := defineCloneRepoStep(taskData, repository, commit)
		firstSteps = append(firstSteps, cloneRepositoryStep)

		// Download artifacts if it is defined as true in the pipe or in the task itself and the clone repository step is enabled
		artifacts := pipe.CloneRepository.Options.Artifacts || taskData.CloneRepository.Options.Artifacts
		if artifacts {
			logger.Info("Adding download artifacts step", "taskName", taskName)
			downloadArtifactsStep := defineDownloadArtifactsStep(logger, taskData, repository, commit)
			firstSteps = append(firstSteps, downloadArtifactsStep)
		}

		// idem for cache
		caches := pipe.CloneRepository.Options.Cache || taskData.CloneRepository.Options.Cache
		if caches {
			logger.Info("Adding download cache step", "taskName", taskName)
			downloadCacheStep := defineDownloadCacheStep(logger, taskData, repository, commit)
			firstSteps = append(firstSteps, downloadCacheStep)
		}
	}

	// Add all these automatic steps at the beginning of the task
	taskData.Steps = append(firstSteps, taskData.Steps...)

	// This is the list of steps that will be added at the end of the task
	var lastSteps []v1alpha1.Step

	// If the clone repository step is enabled, upload the artifacts and cache
	if cloneRepository {
		// Upload artifacts if it is defined as true in the pipe or in the task itself and the clone repository step is enabled
		artifacts := pipe.CloneRepository.Options.Artifacts || taskData.CloneRepository.Options.Artifacts
		if artifacts {
			logger.Info("Adding upload artifacts step", "taskName", taskName)
			uploadArtifactsStep := defineUploadArtifactsStep(logger, taskData, repository, commit)
			lastSteps = append(lastSteps, uploadArtifactsStep)
		}

		// idem for cache
		caches := pipe.CloneRepository.Options.Cache || taskData.CloneRepository.Options.Cache
		if caches {
			logger.Info("Adding upload cache step", "taskName", taskName)
			uploadCacheStep := defineUploadCacheStep(logger, taskData, repository, commit)
			lastSteps = append(lastSteps, uploadCacheStep)
		}
	}

	// Add all the upload steps at the end of the task
	taskData.Steps = append(taskData.Steps, lastSteps...)

	// Add the default volumes for the workspace and the ssh credentials secret if it is defined to the task
	taskData = addDefaultVolumes(taskData, pipe.Workspace, pipe.SshSecretName)

	// Add the volumeMounts for the workspaceDir and the ssh secret if it is defined to the steps
	for i := range taskData.Steps {
		taskData.Steps[i] = addDefaultVolumeMounts(taskData.Steps[i], workspaceDir, pipe.SshSecretName)
	}

	// Clean
	// Remove unnecessary cloneRepository and path from the task
	taskData.CloneRepository = v1alpha1.CloneRepositoryConfig{}
	taskData.Paths = v1alpha1.Paths{}

	return taskData
}

// GetWorkspaceDir returns the default workspace directory
func GetWorkspaceDir() string {
	return workspaceDir
}
