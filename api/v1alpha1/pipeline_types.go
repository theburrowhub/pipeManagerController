/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineSpec defines the desired state of Pipeline.
type PipelineSpec struct {
	// +kubebuilder:validation:MinLength=0
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:MinLength=0
	// +optional
	Description string    `json:"description,omitempty"`
	Namespace   Namespace `json:"namespace,omitempty"`
	// +optional
	CloneRepository CloneRepositoryConfig `json:"cloneRepository,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +optional
	CloneDepth int `json:"cloneDepth,omitempty"`
	// +kubebuilder:validation:MinLength=0
	// +optional
	SshSecretName string `json:"sshSecretName,omitempty"`
	// +optional
	Launch Launch `json:"launch,omitempty"`
	// +optional
	Params map[string]string `json:"params"`
	// +optional
	Workspace corev1.Volume   `json:"workspaceDir,omitempty"`
	Tasks     map[string]Task `json:"tasks"`
	// +optional
	FinishTasks FinishTasks `json:"finishTasks,omitempty"`
}

// CloneRepositoryOptions define the options for the clone repository step in the pipeline
type CloneRepositoryOptions struct {
	// +optional
	Cache bool `json:"cache,omitempty"`
	// +optional
	Artifacts bool `json:"artifacts,omitempty"`
}

// CloneRepositoryConfig define the configuration for the clone repository step in the pipeline
type CloneRepositoryConfig struct {
	// +optional
	Enable bool `json:"enable,omitempty"`
	// +optional
	Options CloneRepositoryOptions `json:"options,omitempty"`
}

// FinishTasks is a struct to store the finish tasks data
type FinishTasks struct {
	// +optional
	Fail map[string]Task `json:"fail,omitempty"`
	// +optional
	Success map[string]Task `json:"success,omitempty"`
}

// Task is a struct to store the task data
type Task struct {
	// +kubebuilder:validation:MinLength=0
	// +optional
	Description string `json:"description"`
	// +optional
	RunAfter []string `json:"runAfter,omitempty"`
	// +optional
	Batch map[string]Batch `json:"batch,omitempty"`
	// +optional
	Params map[string]string `json:"params,omitempty"`
	// +optional
	CloneRepository CloneRepositoryConfig `json:"cloneRepository,omitempty"`
	// +optional
	Paths Paths `json:"paths,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +optional
	CloneDepth int `json:"cloneDepth,omitempty"`
	// +kubebuilder:validation:MinItems=1
	Steps []Step `json:"steps"`
	// +optional
	Sidecars []corev1.Container `json:"sidecars,omitempty"`
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`
}

// Paths define the paths for the artifacts and cache in the task
type Paths struct {
	// +optional
	Artifacts []string `json:"artifacts,omitempty"`
	// +optional
	Cache []string `json:"cache,omitempty"`
}

// Step is a struct to store the step data
type Step struct {
	// +kubebuilder:validation:MinLength=0
	Name string `json:"name"`
	// +kubebuilder:validation:MinLength=0
	Image string `json:"image"`
	// +kubebuilder:validation:MinLength=0
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// +kubebuilder:validation:MinItems=1
	// +optional
	Command []string `json:"command,omitempty"`
	// +kubebuilder:validation:MinItems=1
	// +optional
	Args []string `json:"args,omitempty"`
	// +kubebuilder:validation:MinLength=1
	// +optional
	Script string `json:"script,omitempty"`
}

// Batch define the batch data like params for the batch task
type Batch map[string]string

// Launch is a struct to store the launch data in the pipeline if it is defined
// It contains the steps to launch the next pipeline in the chain when the current pipeline finishes
// successfully or fails
type Launch struct {
	// +optional
	WhenFail []string `json:"whenFail,omitempty"`
	// +optional
	WhenSuccess []string `json:"whenSuccess,omitempty"`
}

// Namespace is a struct to store the namespace data in the pipeline
// It contains the name of the namespace and the labels to apply to the namespace, if any. Also, a boolean to
// indicate if the namespace should be created or not
type Namespace struct {
	// +optional
	Create bool `json:"create,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +kubebuilder:validation:MinLength=0
	Name string `json:"name"`
}

// PipelineStatus defines the observed state of Pipeline.
type PipelineStatus struct {
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Pipeline is the Schema for the pipelines API.
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineList contains a list of Pipeline.
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
	err := AddToScheme(Scheme)
	if err != nil {
		return
	}
}
