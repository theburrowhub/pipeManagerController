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
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"
	"github.com/sergiotejon/pipeManagerController/internal/validate"
)

// nolint:unused
// log is for logging in this package.
var pipelinelog = logf.Log.WithName("pipeline-resource")

// SetupPipelineWebhookWithManager registers the webhook for Pipeline in the manager.
func SetupPipelineWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&pipemanagerv1alpha1.Pipeline{}).
		WithValidator(&PipelineCustomValidator{}).
		WithDefaulter(&PipelineCustomDefaulter{
			CloneRepository: true,
			Artifacts:       false,
			Cache:           false,
		}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-pipemanager-sergiotejon-github-io-v1alpha1-pipeline,mutating=true,failurePolicy=fail,sideEffects=None,groups=pipemanager.sergiotejon.github.io,resources=pipelines,verbs=create;update,versions=v1alpha1,name=mpipeline-v1alpha1.kb.io,admissionReviewVersions=v1

// PipelineCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Pipeline when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PipelineCustomDefaulter struct {
	CloneRepository bool
	Artifacts       bool
	Cache           bool
}

var _ webhook.CustomDefaulter = &PipelineCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Pipeline.
func (d *PipelineCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	pipeline, ok := obj.(*pipemanagerv1alpha1.Pipeline)

	if !ok {
		return fmt.Errorf("expected an Pipeline object but got %T", obj)
	}
	pipelinelog.Info("Defaulting for Pipeline", "name", pipeline.GetName())

	// Set default values
	d.applyDefaults(pipeline)

	return nil
}

func (d *PipelineCustomDefaulter) applyDefaults(pipeline *pipemanagerv1alpha1.Pipeline) {
	if (pipeline.Spec.CloneRepository == pipemanagerv1alpha1.CloneRepositoryConfig{}) {
		pipeline.Spec.CloneRepository = pipemanagerv1alpha1.CloneRepositoryConfig{
			Enable: d.CloneRepository,
			Options: pipemanagerv1alpha1.CloneRepositoryOptions{
				Cache:     d.Cache,
				Artifacts: d.Artifacts,
			},
		}
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-pipemanager-sergiotejon-github-io-v1alpha1-pipeline,mutating=false,failurePolicy=fail,sideEffects=None,groups=pipemanager.sergiotejon.github.io,resources=pipelines,verbs=create;update,versions=v1alpha1,name=vpipeline-v1alpha1.kb.io,admissionReviewVersions=v1

// PipelineCustomValidator struct is responsible for validating the Pipeline resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type PipelineCustomValidator struct {
	//TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &PipelineCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Pipeline.
func (v *PipelineCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	pipeline, ok := obj.(*pipemanagerv1alpha1.Pipeline)
	if !ok {
		return nil, fmt.Errorf("expected a Pipeline object but got %T", obj)
	}
	pipelinelog.Info("Validation for Pipeline upon creation", "name", pipeline.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, validatePipeline(pipeline)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Pipeline.
func (v *PipelineCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	pipeline, ok := newObj.(*pipemanagerv1alpha1.Pipeline)
	if !ok {
		return nil, fmt.Errorf("expected a Pipeline object for the newObj but got %T", newObj)
	}
	pipelinelog.Info("Validation for Pipeline upon update", "name", pipeline.GetName())

	// TODO: Allow updates status only

	// Impede updates by returning an error
	return nil, fmt.Errorf("updates to Pipeline objects are not allowed")
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Pipeline.
func (v *PipelineCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	pipeline, ok := obj.(*pipemanagerv1alpha1.Pipeline)
	if !ok {
		return nil, fmt.Errorf("expected a Pipeline object but got %T", obj)
	}
	pipelinelog.Info("Validation for Pipeline upon deletion", "name", pipeline.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

// validatePipeline validates the Pipeline object
func validatePipeline(pipeline *pipemanagerv1alpha1.Pipeline) error {
	// Validate the pipeline object

	// Validate runAfter field if they are corrected referenced to other tasks in the pipeline
	err := validate.RunAfterValidation(pipeline)
	if err != nil {
		return err
	}

	// TODO(user): Add more validation logic here

	return nil
}
