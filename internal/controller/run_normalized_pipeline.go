package controller

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
	"github.com/sergiotejon/pipeManagerController/internal/runpipeline/tekton"

	tektonpipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"

	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

// tektonPipelineRun run the pipeline based on the pipeline type and returns the deployed object
// (currently only supports Tekton pipelines).
func (r *PipelineReconciler) tektonPipelineRun(pipelineObject *pipemanagerv1alpha1.Pipeline,
	normalizedPipelineSpec *v1alpha1.PipelineSpec) (error, *tektonpipelinev1.PipelineRun) {

	var dep *tektonpipelinev1.PipelineRun
	var err error

	err, dep = tekton.CreateKubernetesObject(string(pipelineObject.GetUID()), normalizedPipelineSpec)
	if err != nil {
		return err, nil
	}

	// Set the pipeline as the owner of the object created to run the pipeline.
	// Depending on the pipeline type, the object created can be any kind of Kubernetes object
	err = controllerutil.SetControllerReference(pipelineObject, dep, r.Scheme)
	if err != nil {
		return fmt.Errorf("error setting controller reference: %v", err), nil
	}

	return nil, dep
}

// TODO: In the future new functions will be added to run pipelines of different types
