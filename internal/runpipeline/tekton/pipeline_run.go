package tekton

import (
	"fmt"
	"github.com/go-logr/logr"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"

	tektonpipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const workspaceName = "workspace"

var (
	logger   *logr.Logger
	pipeline *v1alpha1.PipelineSpec
)

func Deploy(l *logr.Logger, p *v1alpha1.PipelineSpec) error {
	var err error

	logger = l
	pipeline = p

	//// Aquí puedes iterar sobre las tareas y gestionar cada una
	//for taskName, task := range pipeline.Tasks {
	//	logger.Info("Processing Task", "TaskName", taskName, "Description", task.Description)
	//	for _, step := range task.Steps {
	//		logger.Info("  Step", "Name", step.Name, "Image", step.Image)
	//		// Implementa la lógica para manejar cada step
	//		// Por ejemplo, crear Pods o Jobs según los steps
	//	}
	//}
	//

	// DEBUG
	fmt.Println("-------> Pipeline:")
	// Marshal the normalizedPipelineSpec to YAML
	yamlData, err := yaml.Marshal(pipeline)
	if err != nil {
		logger.Error(err, "Error marshaling normalizedPipelineSpec to YAML")
	} else {
		fmt.Println(string(yamlData))
	}
	// DEBUG

	tektonPipelineRun := tektonpipelinev1.PipelineRun{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PipelineRun",
			APIVersion: "tekton.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pipeline.Name + "-",
			Namespace:    pipeline.Namespace.Name,
		},
		Spec:   buildPipelineRunSpec(),
		Status: tektonpipelinev1.PipelineRunStatus{},
	}

	// DEBUG
	fmt.Println("-------> PipelineRun:")
	// Marshal the normalizedPipelineSpec to YAML
	yamlData, err = yaml.Marshal(tektonPipelineRun)
	if err != nil {
		logger.Error(err, "Error marshaling normalizedPipelineSpec to YAML")
	} else {
		fmt.Println(string(yamlData))
	}
	// DEBUG

	err = runDeploy(&tektonPipelineRun)
	if err != nil {
		logger.Error(err, "Error deploying PipelineRun")
		return err
	}

	return nil
}

func runDeploy(pipelineRun *tektonpipelinev1.PipelineRun) error {
	// TODO: Implement the logic to deploy the PipelineRun to the Kubernetes cluster
	return nil
}

func buildPipelineRunSpec() tektonpipelinev1.PipelineRunSpec {
	spec := tektonpipelinev1.PipelineRunSpec{
		Params:     buildParams(),
		Workspaces: buildWorkspaces(),
		// TODO: implement the rest...
	}

	return spec
}

func buildParams() []tektonpipelinev1.Param {
	var params []tektonpipelinev1.Param

	for name, value := range pipeline.Params {
		params = append(params, tektonpipelinev1.Param{
			Name: name,
			Value: tektonpipelinev1.ParamValue{
				StringVal: value,
				Type:      tektonpipelinev1.ParamTypeString,
			},
		})
	}

	return params
}

func buildWorkspaces() []tektonpipelinev1.WorkspaceBinding {
	var workspaces []tektonpipelinev1.WorkspaceBinding

	// TODO: test this code with a pipeline that has a workspace defined in the configuration section

	workspaces = append(workspaces, tektonpipelinev1.WorkspaceBinding{
		Name: workspaceName,

		// Set every kind of workspace to the same value. Only one of them will be used, that is the one that is not nil
		EmptyDir:              pipeline.Workspace.EmptyDir,
		PersistentVolumeClaim: pipeline.Workspace.PersistentVolumeClaim,
		ConfigMap:             pipeline.Workspace.ConfigMap,
		Secret:                pipeline.Workspace.Secret,
		Projected:             pipeline.Workspace.Projected,
		CSI:                   pipeline.Workspace.CSI,
	})

	return workspaces
}
