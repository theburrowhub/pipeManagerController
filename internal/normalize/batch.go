package normalize

import (
	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

// processBatchTask processes a batch task and returns a map of tasks to run in parallel
func processBatchTask(taskName string, taskData v1alpha1.Task) map[string]v1alpha1.Task {
	tasks := make(map[string]v1alpha1.Task)

	// If no batch is defined, return nil
	if taskData.Batch == nil {
		return nil
	}

	for batchName, batchParams := range taskData.Batch {
		// Add the new task to the tasks map
		name := k8sObjectName(taskName, batchName)

		// Copy taskData to newTask removing the batch field
		newTask := taskData.DeepCopy()
		newTask.Batch = nil
		tasks[name] = *newTask

		// Copy the params from the batch task to the new task
		for key, value := range batchParams {
			tasks[name].Params[key] = value
		}
	}

	return tasks
}

// fixRunAfterReferencesForBatchTasks fixes the runAfter references in the pipeline
// TODO: Mix FinishTasks and Tasks in the same type of data
// TODO: Refactor this function to avoid code duplication
func fixRunAfterReferencesForBatchTasks(pipeline *v1alpha1.PipelineSpec, batchTaskName string, newTaskNames []string) {
	for t := range pipeline.Tasks {
		for i := range pipeline.Tasks[t].RunAfter {
			if pipeline.Tasks[t].RunAfter[i] == batchTaskName {
				newRunAfter := append(pipeline.Tasks[t].RunAfter[:i], append(newTaskNames, pipeline.Tasks[t].RunAfter[i+1:]...)...)
				task := pipeline.Tasks[t]
				task.RunAfter = newRunAfter
				pipeline.Tasks[t] = task
				break
			}
		}
	}
	for t := range pipeline.FinishTasks.Fail {
		for i := range pipeline.FinishTasks.Fail[t].RunAfter {
			if pipeline.FinishTasks.Fail[t].RunAfter[i] == batchTaskName {
				newRunAfter := append(pipeline.FinishTasks.Fail[t].RunAfter[:i], append(newTaskNames, pipeline.FinishTasks.Fail[t].RunAfter[i+1:]...)...)
				task := pipeline.FinishTasks.Fail[t]
				task.RunAfter = newRunAfter
				pipeline.FinishTasks.Fail[t] = task
				break
			}
		}
	}
	for t := range pipeline.FinishTasks.Success {
		for i := range pipeline.FinishTasks.Success[t].RunAfter {
			if pipeline.FinishTasks.Success[t].RunAfter[i] == batchTaskName {
				newRunAfter := append(pipeline.FinishTasks.Success[t].RunAfter[:i], append(newTaskNames, pipeline.FinishTasks.Success[t].RunAfter[i+1:]...)...)
				task := pipeline.FinishTasks.Success[t]
				task.RunAfter = newRunAfter
				pipeline.FinishTasks.Success[t] = task
				break
			}
		}
	}
}
