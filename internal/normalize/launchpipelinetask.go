package normalize

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
)

const envVarPrefix = "PIPELINE_"

// defineLaunchPipelineTask defines the task to launch the next pipeline in the chain
func defineLaunchPipelineTask(currentPipeline v1alpha1.PipelineSpec, repository, commit, pipelineToLaunch string) v1alpha1.Task {
	var env []corev1.EnvVar

	// Set the parameters for the new pipeline through environment variables
	for name, value := range envvars.Variables {
		env = append(env, corev1.EnvVar{
			Name:  envVarPrefix + name,
			Value: value,
		})
	}

	// Set environment variable for configuration instead of using the configuration file
	env = append(env, corev1.EnvVar{
		Name:  envVarPrefix + "NAME",
		Value: pipelineToLaunch,
	})
	env = append(env, corev1.EnvVar{
		Name:  "COMMON_DATA_LOG_LEVEL",
		Value: config.Common.Data.Log.Level,
	})
	env = append(env, corev1.EnvVar{
		Name:  "COMMON_DATA_LOG_FORMAT",
		Value: config.Common.Data.Log.Format,
	})
	env = append(env, corev1.EnvVar{
		Name:  "COMMON_DATA_LOG_FILE",
		Value: config.Common.Data.Log.File,
	})
	env = append(env, corev1.EnvVar{
		Name:  "LAUNCHER_DATA_CLONEDEPTH",
		Value: fmt.Sprintf("%d", config.Launcher.Data.CloneDepth),
	})
	env = append(env, corev1.EnvVar{
		Name:  "LAUNCHER_DATA_NAMESPACE",
		Value: config.Launcher.Data.Namespace,
	})

	// Get the image from the configuration to launch the pipeline
	launcherImage := config.Launcher.Data.GetLauncherImage()

	// Define the task to launch the pipeline
	launchPipelineTask := v1alpha1.Task{
		Description: fmt.Sprintf("Launch the next pipeline '%s' in the chain", pipelineToLaunch),
		Steps: []v1alpha1.Step{
			{
				Name:        "launch-pipeline",
				Description: fmt.Sprintf("Launch the next pipeline '%s' in the chain", pipelineToLaunch),
				Image:       launcherImage,
				Env:         env,
				Script: strings.Join([]string{
					fmt.Sprintf("#!%s", defaultShell),
					defaultShellSets,
					fmt.Sprintf("%s %s", launcherBinary, "run"),
				}, "\n"),
			},
		},
	}

	// Append the clone repository step to the task steps to ensure the repository is cloned before launching the pipeline.
	// It's the only way to look for the pipeline definition of pipelineToLaunch in the repository.
	cloneRepositoryStep := defineCloneRepoStep(launchPipelineTask, repository, commit)
	launchPipelineTask.Steps = append([]v1alpha1.Step{cloneRepositoryStep}, launchPipelineTask.Steps...)

	return launchPipelineTask
}
