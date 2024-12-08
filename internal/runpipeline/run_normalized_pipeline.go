package runpipeline

import (
	"fmt"
	"github.com/go-logr/logr"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	"github.com/sergiotejon/pipeManagerController/internal/runpipeline/tekton"
	"github.com/sergiotejon/pipeManagerLauncher/pkg/config"
)

type PipelineDeployFunc func(*logr.Logger, *v1alpha1.PipelineSpec) error

var pipelineDeployFuncs = map[string]PipelineDeployFunc{
	"tekton":  tekton.Deploy,
	"default": tekton.Deploy, // Assuming default is also tekton.Deploy
}

func NormalizedPipeline(l *logr.Logger, pipeline *v1alpha1.PipelineSpec) error {
	// Generate the pipeline object based on the pipeline type
	// Currently only supports Tekton pipelines
	pipelineType := config.Launcher.Data.PipelineType
	pipelineDeploy, exists := pipelineDeployFuncs[pipelineType]
	if !exists {
		return fmt.Errorf("unsupported pipeline type: %s", pipelineType)
	}

	err := pipelineDeploy(l, pipeline)
	if err != nil {
		return err
	}

	return nil
}
