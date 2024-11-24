package runpipeline

import (
	"github.com/go-logr/logr"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManagerController/internal/runpipeline/tekton"
	"github.com/sergiotejon/pipeManagerLauncher/pkg/config"
)

func NormalizedPipeline(l *logr.Logger, pipeline *v1alpha1.PipelineSpec) error {

	err := createNamespace(pipeline.Namespace.Name)
	if err != nil {
		return err
	}

	// Generate the pipeline object based on the pipeline type
	// Currently only supports Tekton pipelines
	switch config.Launcher.Data.PipelineType {
	case "tekton":
		err := tekton.Deploy(l, pipeline)
		if err != nil {
			return err
		}
	default:
		err := tekton.Deploy(l, pipeline)
		if err != nil {
			return err
		}
	}

	return nil
}

func createNamespace(namespace string) error {
	// TODO: Create Namespace if it doesn't exist
	return nil
}
