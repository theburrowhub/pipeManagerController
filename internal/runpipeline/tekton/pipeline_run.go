package tekton

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/selection"
	"os"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
	"github.com/sergiotejon/pipeManagerController/internal/normalize"

	tektonpipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

const (
	workspaceName = "workspace"
)

var (
	logger   *logr.Logger
	pipeline *v1alpha1.PipelineSpec
)

func Deploy(l *logr.Logger, p *v1alpha1.PipelineSpec) error {
	var err error

	logger = l
	pipeline = p

	// DEBUG
	fmt.Println("-------> Pipeline Spec:")
	jsonData, err := json.MarshalIndent(pipeline, "", " ")
	if err != nil {
		logger.Error(err, "Error marshaling normalizedPipelineSpec to JSON")
	} else {
		err = os.WriteFile("/tmp/pipeline-spec.json", jsonData, 0666)
		if err != nil {
			logger.Error(err, "No funciona esto")
		}
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

	err = runDeploy(&tektonPipelineRun)
	if err != nil {
		logger.Error(err, "Error deploying PipelineRun")
		return err
	}

	return nil
}

func runDeploy(pipelineRun *tektonpipelinev1.PipelineRun) error {
	// TODO: Implement the logic to deploy the PipelineRun to the Kubernetes cluster
	// DEBUG
	fmt.Println("-------> PipelineRun:")
	jsonData, err := json.MarshalIndent(pipelineRun, "", " ")
	if err != nil {
		logger.Error(err, "Error marshaling pipeline run to JSON")
	} else {
		err = os.WriteFile("/tmp/pipeline-run.json", jsonData, 0666)
		if err != nil {
			logger.Error(err, "No funciona esto")
		}
	}
	// DEBUG

	return nil
}

func buildPipelineRunSpec() tektonpipelinev1.PipelineRunSpec {
	spec := tektonpipelinev1.PipelineRunSpec{
		Params:     buildParams(),
		Workspaces: buildWorkspacesBinding(),
		// TODO: Implement timeouts by global config and pipeline configuration as optional
		//Timeouts: ...
		PipelineSpec: buildPipelineSpec(),
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

func buildWorkspacesBinding() []tektonpipelinev1.WorkspaceBinding {
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

func buildPipelineSpec() *tektonpipelinev1.PipelineSpec {
	pipelineSpec := &tektonpipelinev1.PipelineSpec{
		Description: pipeline.Description,
		DisplayName: pipeline.Name,
		Tasks:       buildTasks(),
		Finally:     buildFinallyTasks(),
	}

	return pipelineSpec
}

func buildTasks() []tektonpipelinev1.PipelineTask {
	var tasks []tektonpipelinev1.PipelineTask

	for name, taskData := range pipeline.Tasks {
		tasks = append(tasks, tektonpipelinev1.PipelineTask{
			Name:        name,
			DisplayName: name,
			Description: taskData.Description,
			RunAfter:    taskData.RunAfter,
			TaskSpec:    buildTaskSpec(name, taskData),
			// TODO:
			// Retries:
		})

	}
	return tasks
}

func buildFinallyTasks() []tektonpipelinev1.PipelineTask {
	var tasks []tektonpipelinev1.PipelineTask

	for name, taskData := range pipeline.FinishTasks.Success {
		tasks = append(tasks, tektonpipelinev1.PipelineTask{
			Name:        name,
			DisplayName: name,
			Description: taskData.Description,
			RunAfter:    taskData.RunAfter,
			TaskSpec:    buildTaskSpec(name, taskData),
			When:        buildWhenSuccessExpresion(),
			// TODO:
			// Retries:
		})
	}

	for name, taskData := range pipeline.FinishTasks.Fail {
		tasks = append(tasks, tektonpipelinev1.PipelineTask{
			Name:        name,
			DisplayName: name,
			Description: taskData.Description,
			RunAfter:    taskData.RunAfter,
			TaskSpec:    buildTaskSpec(name, taskData),
			When:        buildWhenFailExpresion(),
			// TODO:
			// Retries:
		})
	}

	return tasks
}

func buildTaskSpec(name string, data v1alpha1.Task) *tektonpipelinev1.EmbeddedTask {
	return &tektonpipelinev1.EmbeddedTask{
		TaskSpec: tektonpipelinev1.TaskSpec{
			DisplayName: name,
			Description: data.Description,
			Params:      buildTaskParams(data),
			Steps:       buildTaskSteps(data),
			Volumes:     buildTaskVolumes(data),
			Sidecars:    buildTaskSidecars(data),
		},
	}
}

func buildTaskParams(data v1alpha1.Task) tektonpipelinev1.ParamSpecs {
	var paramSpecs tektonpipelinev1.ParamSpecs

	for paramName, _ := range data.Params {
		paramSpecs = append(paramSpecs, tektonpipelinev1.ParamSpec{
			Name: paramName,
			Type: tektonpipelinev1.ParamTypeString,
		})
	}

	return paramSpecs
}

func buildTaskSteps(data v1alpha1.Task) []tektonpipelinev1.Step {
	var steps []tektonpipelinev1.Step

	for _, stepData := range data.Steps {
		steps = append(steps, tektonpipelinev1.Step{
			Name:         stepData.Name,
			Image:        stepData.Image,
			Env:          stepData.Env,
			Command:      stepData.Command,
			Args:         stepData.Args,
			Script:       stepData.Script,
			VolumeMounts: stepData.VolumeMounts,
			WorkingDir:   normalize.GetWorkspaceDir(),

			// TODO:
			//ComputeResources:
			//Timeout: (task timeout)
			//When:
		})
	}

	return steps
}

func buildTaskVolumes(data v1alpha1.Task) []v1.Volume {
	var volumes []v1.Volume

	volumes = append(volumes, pipeline.Workspace)
	for _, volume := range data.Volumes {
		volumes = append(volumes, volume)
	}

	return volumes
}

func buildTaskSidecars(data v1alpha1.Task) []tektonpipelinev1.Sidecar {
	var sidecars []tektonpipelinev1.Sidecar

	for _, sidecar := range data.Sidecars {
		sidecars = append(sidecars, tektonpipelinev1.Sidecar{
			Name:                     sidecar.Name,
			Image:                    sidecar.Image,
			Command:                  sidecar.Command,
			Args:                     sidecar.Args,
			WorkingDir:               sidecar.WorkingDir,
			Ports:                    sidecar.Ports,
			EnvFrom:                  sidecar.EnvFrom,
			Env:                      sidecar.Env,
			VolumeMounts:             sidecar.VolumeMounts,
			VolumeDevices:            sidecar.VolumeDevices,
			LivenessProbe:            sidecar.LivenessProbe,
			ReadinessProbe:           sidecar.ReadinessProbe,
			StartupProbe:             sidecar.StartupProbe,
			Lifecycle:                sidecar.Lifecycle,
			TerminationMessagePath:   sidecar.TerminationMessagePath,
			TerminationMessagePolicy: sidecar.TerminationMessagePolicy,
			ImagePullPolicy:          sidecar.ImagePullPolicy,
			SecurityContext:          sidecar.SecurityContext,
			Stdin:                    sidecar.Stdin,
			StdinOnce:                sidecar.StdinOnce,
			TTY:                      sidecar.TTY,
			RestartPolicy:            sidecar.RestartPolicy,
		})
	}

	return sidecars
}

func buildWhenSuccessExpresion() tektonpipelinev1.WhenExpressions {
	var exp tektonpipelinev1.WhenExpressions

	exp = append(exp, tektonpipelinev1.WhenExpression{
		Input:    "$(tasks.status)",
		Operator: selection.NotIn,
		Values:   []string{"Failed"},
	})

	return exp
}

func buildWhenFailExpresion() tektonpipelinev1.WhenExpressions {
	var exp tektonpipelinev1.WhenExpressions

	exp = append(exp, tektonpipelinev1.WhenExpression{
		Input:    "$(tasks.status)",
		Operator: selection.In,
		Values:   []string{"Failed"},
	})

	return exp
}
