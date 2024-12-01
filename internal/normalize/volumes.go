package normalize

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

// addDefaultVolumes adds the default volumes to the task
func addDefaultVolumes(task v1alpha1.Task, workspace corev1.Volume, sshSecretName string) v1alpha1.Task {
	// Volumes for the workspaceDir and the ssh secret if it is defined
	task.Volumes = append(task.Volumes, workspaceVolume(workspace))
	if sshSecretName != "" {
		task.Volumes = append(task.Volumes, sshSecretVolume(sshSecretName))
	}

	return task
}

// workspaceVolume defines the volume for the workspaceDir
func workspaceVolume(w corev1.Volume) corev1.Volume {
	if (w != corev1.Volume{}) {
		return w
	} else {
		return corev1.Volume{
			Name: "workspace",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
	}
}

// sshSecretVolume defines the volume for the ssh secret
func sshSecretVolume(sshSecretName string) corev1.Volume {
	return corev1.Volume{
		Name: "ssh-credentials",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  sshSecretName,
				DefaultMode: func(i int) *int32 { v := int32(i); return &v }(256),
			},
		},
	}
}

// addDefaultVolumeMounts adds the default volume mounts to the step
func addDefaultVolumeMounts(step v1alpha1.Step, workspaceDir, sshSecretName string) v1alpha1.Step {
	// Volume mounts for the workspaceDir and the ssh secret if it is defined
	step.VolumeMounts = append(step.VolumeMounts, workspaceVolumeMount(workspaceDir))
	if sshSecretName != "" {
		step.VolumeMounts = append(step.VolumeMounts, sshSecretVolumeMount())
	}

	return step
}

// workspaceVolumeMount defines the volume mount for the workspaceDir
func workspaceVolumeMount(mountPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "workspace",
		MountPath: mountPath,
	}
}

// sshSecretVolumeMount defines the volume mount for the ssh secret
func sshSecretVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "ssh-credentials",
		MountPath: "/root/.ssh",
		ReadOnly:  true,
	}
}
