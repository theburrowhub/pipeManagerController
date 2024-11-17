package normalize

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManager/internal/pkg/config"
	"github.com/sergiotejon/pipeManager/internal/pkg/envvars"
)

// defineDownloadArtifactsStep defines the download artifacts step in the task
func defineDownloadArtifactsStep(logger logr.Logger, taskData v1alpha1.Task) v1alpha1.Step {
	return defineArtifactBucketStep(logger, "artifacts", "download", taskData.Paths.Artifacts)
}

// defineDownloadCacheStep defines the download cache step in the task
func defineDownloadCacheStep(logger logr.Logger, taskData v1alpha1.Task) v1alpha1.Step {
	return defineArtifactBucketStep(logger, "cache", "download", taskData.Paths.Cache)
}

// defineUploadArtifactsStep defines the upload artifacts step in the task
func defineUploadArtifactsStep(logger logr.Logger, taskData v1alpha1.Task) v1alpha1.Step {
	return defineArtifactBucketStep(logger, "artifacts", "upload", taskData.Paths.Artifacts)
}

// defineUploadCacheStep defines the upload cache step in the task
func defineUploadCacheStep(logger logr.Logger, taskData v1alpha1.Task) v1alpha1.Step {
	return defineArtifactBucketStep(logger, "cache", "upload", taskData.Paths.Cache)
}

// defineArtifactBucketStep defines the download step in the task for artifacts or cache
func defineArtifactBucketStep(logger logr.Logger, commandType, commandAction string, paths []string) v1alpha1.Step {
	// Get the download image from the configuration
	image := config.Launcher.Data.GetLauncherImage()

	// Get configuration for the bucket to store the artifacts or cache
	bucketURL := config.Launcher.Data.ArtifactsBucket.URL
	basePath := config.Launcher.Data.ArtifactsBucket.BasePath
	bucketParameters, err := json.Marshal(config.Launcher.Data.ArtifactsBucket.Parameters)
	if err != nil {
		logger.Error(err, "Error marshalling bucketParameters")
		bucketParameters = []byte("{}")
	}
	credentials := config.Launcher.Data.ArtifactsBucket.Credentials

	// Define the download step
	env, volumeMounts, command := defineArtifactBucketStepConfig(
		commandType, commandAction, bucketURL, basePath, paths, bucketParameters, credentials, workspaceDir,
	)
	step := v1alpha1.Step{
		Name:         fmt.Sprintf("%s-%s", commandAction, commandType),
		Description:  fmt.Sprintf("Automatically %s the %s", commandAction, commandType),
		Image:        image,
		Env:          env,
		VolumeMounts: volumeMounts,
		Script: strings.Join([]string{
			fmt.Sprintf("#!%s", defaultShell),
			defaultShellSets,
			command,
		}, "\n"),
	}

	return step
}

// defineArtifactBucketStepConfig defines the configuration for the download step in the task for downloading artifacts or cache
func defineArtifactBucketStepConfig(commandType, commandAction, bucketURL, basePath string, paths []string,
	bucketParameters []byte, credentials config.BucketCredentials, workspaceDir string) ([]corev1.EnvVar, []corev1.VolumeMount, string) {

	// Define the environment variables for the download artifacts step
	// Environment variables example:
	//
	// # Configuration for the bucket to store the artifacts
	//     LAUNCHER_DATA_BUCKET_URL="s3://pipemanager"
	//     LAUNCHER_DATA_BUCKET_BASEPATH="pipe-manager"
	//     LAUNCHER_DATA_BUCKET_PARAMETERS="{\"endpoint\":\"localhost:9000\",\"disableSSL\":true,\"s3ForcePathStyle\":true,\"awssdk\":\"v1\"}"
	//
	// # Environment variables for the credentials
	//     AWS_ACCESS_KEY_ID="EXAMPLE_ACCESS_KEY_ID"
	//     AWS_SECRET_ACCESS_KEY="************"
	env := []corev1.EnvVar{
		{
			Name:  "LAUNCHER_DATA_BUCKET_URL",
			Value: bucketURL,
		},
		{
			Name:  "LAUNCHER_DATA_BUCKET_BASEPATH",
			Value: basePath,
		},
		{
			Name:  "LAUNCHER_DATA_BUCKET_PARAMETERS",
			Value: string(bucketParameters),
		},
	}
	for _, envVar := range credentials.Env {
		env = append(env, envVar)
	}

	// Define the volume mounts for the download artifacts step from the credentials if they are defined in the configuration
	// Example:
	//    volumeMounts:
	//	     - mountPath: /etc/s3-credentials
	//    name: s3-credentials
	//    readOnly: true
	var volumeMounts []corev1.VolumeMount
	for _, volumeMount := range credentials.VolumeMounts {
		volumeMounts = append(volumeMounts, volumeMount)
	}

	// Define the command to download the artifacts
	// Example: `/app/launcher artifacts download --commit XXXXXXXXX --destination /workdir \
	//				--project "git@github.com:my-user/my-project.git" --path exampleDir1 --path exampleDir2`
	pathArgs := ""
	for _, path := range paths {
		pathArgs += fmt.Sprintf("--path '%s' ", path)
	}
	command := ""
	if commandType == "artifacts" {
		command = fmt.Sprintf("%s %s %s --commit '%s' --destination '%s' --project '%s' %s",
			launcherBinary, commandType, commandAction, envvars.Variables["COMMIT"], workspaceDir, envvars.Variables["REPOSITORY"], pathArgs)
	} else {
		command = fmt.Sprintf("%s %s %s --destination '%s' --project '%s' %s",
			launcherBinary, commandType, commandAction, workspaceDir, envvars.Variables["REPOSITORY"], pathArgs)
	}

	return env, volumeMounts, command
}
