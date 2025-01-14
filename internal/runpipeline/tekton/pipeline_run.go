package tekton

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/selection"

	tektonpipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
	"github.com/sergiotejon/pipeManagerController/internal/k8s"
	"github.com/sergiotejon/pipeManagerController/internal/normalize"
)

const (
	workspaceName = "workspace"
)

var (
	pipelineSpec *v1alpha1.PipelineSpec
)

func CreateKubernetesObject(pipelineUID string, p *v1alpha1.PipelineSpec) (error, *tektonpipelinev1.PipelineRun) {
	pipelineSpec = p

	tektonPipelineRun := tektonpipelinev1.PipelineRun{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PipelineRun",
			APIVersion: "tekton.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: pipelineSpec.Name + "-",
			Namespace:    pipelineSpec.Namespace.Name,
			Labels: map[string]string{
				"pipelineRef": pipelineUID,
			},
		},
		Spec:   buildPipelineRunSpec(),
		Status: tektonpipelinev1.PipelineRunStatus{},
	}

	return nil, &tektonPipelineRun
}

func buildPipelineRunSpec() tektonpipelinev1.PipelineRunSpec {
	spec := tektonpipelinev1.PipelineRunSpec{
		Params:     buildParams(),
		Workspaces: buildWorkspacesBinding(),
		// TODO: Implement timeouts by global config and pipelineSpec configuration as optional
		//Timeouts: ...
		PipelineSpec: buildPipelineSpec(),
	}

	return spec
}

func buildParams() []tektonpipelinev1.Param {
	var params []tektonpipelinev1.Param

	for name, value := range pipelineSpec.Params {
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

	// TODO: test this code with a pipelineSpec that has a workspace defined in the configuration section

	// If the pipelineSpec has not defined a workspace, create an emptyDir workspace by default
	emptyDir := pipelineSpec.Workspace.EmptyDir
	if pipelineSpec.Workspace.EmptyDir == nil &&
		pipelineSpec.Workspace.PersistentVolumeClaim == nil &&
		pipelineSpec.Workspace.ConfigMap == nil &&
		pipelineSpec.Workspace.Secret == nil &&
		pipelineSpec.Workspace.Projected == nil &&
		pipelineSpec.Workspace.CSI == nil {
		emptyDir = &v1.EmptyDirVolumeSource{}
	}

	workspaces = append(workspaces, tektonpipelinev1.WorkspaceBinding{
		Name: workspaceName,

		// Set every kind of workspace to the same value. Only one of them will be used, that is the one that is not nil
		EmptyDir:              emptyDir,
		PersistentVolumeClaim: pipelineSpec.Workspace.PersistentVolumeClaim,
		ConfigMap:             pipelineSpec.Workspace.ConfigMap,
		Secret:                pipelineSpec.Workspace.Secret,
		Projected:             pipelineSpec.Workspace.Projected,
		CSI:                   pipelineSpec.Workspace.CSI,
	})

	return workspaces
}

func buildPipelineSpec() *tektonpipelinev1.PipelineSpec {
	pipelineSpec := &tektonpipelinev1.PipelineSpec{
		Description: pipelineSpec.Description,
		DisplayName: pipelineSpec.Name,
		Tasks:       buildTasks(),
		Finally:     buildFinallyTasks(),
	}

	return pipelineSpec
}

func buildTasks() []tektonpipelinev1.PipelineTask {
	var tasks []tektonpipelinev1.PipelineTask

	for name, taskData := range pipelineSpec.Tasks {
		tasks = append(tasks, tektonpipelinev1.PipelineTask{
			Name:        k8s.ToK8sName(name),
			DisplayName: name,
			Description: taskData.Description,
			RunAfter:    k8s.ToK8sNames(taskData.RunAfter),
			Params:      buildTaskParams(taskData),
			TaskSpec:    buildTaskSpec(name, taskData),
			// TODO:
			// Retries:
		})

	}
	return tasks
}

func buildFinallyTasks() []tektonpipelinev1.PipelineTask {
	var tasks []tektonpipelinev1.PipelineTask

	for name, taskData := range pipelineSpec.FinishTasks.Success {
		tasks = append(tasks, tektonpipelinev1.PipelineTask{
			Name:        k8s.ToK8sName("success_" + name),
			DisplayName: name,
			Description: taskData.Description,
			RunAfter:    k8s.ToK8sNames(taskData.RunAfter),
			Params:      buildTaskParams(taskData),
			TaskSpec:    buildTaskSpec(name, taskData),
			When:        buildWhenSuccessExpresion(),
			// TODO:
			// Retries:
		})
	}

	for name, taskData := range pipelineSpec.FinishTasks.Fail {
		tasks = append(tasks, tektonpipelinev1.PipelineTask{
			Name:        k8s.ToK8sName("fail_" + name),
			DisplayName: name,
			Description: taskData.Description,
			RunAfter:    k8s.ToK8sNames(taskData.RunAfter),
			Params:      buildTaskParams(taskData),
			TaskSpec:    buildTaskSpec(name, taskData),
			When:        buildWhenFailExpresion(),
			// TODO:
			// Retries:
		})
	}

	return tasks
}

func buildTaskParams(data v1alpha1.Task) []tektonpipelinev1.Param {
	var params []tektonpipelinev1.Param

	for paramName, paramValue := range data.Params {
		params = append(params, tektonpipelinev1.Param{
			Name: paramName,
			Value: tektonpipelinev1.ParamValue{
				StringVal: paramValue,
				Type:      tektonpipelinev1.ParamTypeString,
			},
		})
	}

	return params
}

func buildTaskSpec(name string, data v1alpha1.Task) *tektonpipelinev1.EmbeddedTask {
	return &tektonpipelinev1.EmbeddedTask{
		TaskSpec: tektonpipelinev1.TaskSpec{
			DisplayName: name,
			Description: data.Description,
			Params:      buildTaskParamSpecs(data),
			Steps:       buildTaskSteps(data),
			Volumes:     buildTaskVolumes(data),
			Sidecars:    buildTaskSidecars(data),
		},
	}
}

func buildTaskParamSpecs(data v1alpha1.Task) tektonpipelinev1.ParamSpecs {
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
			Name:         k8s.ToK8sName(stepData.Name),
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

	if pipelineSpec.Workspace.Name != "" {
		volumes = append(volumes, pipelineSpec.Workspace)
	}
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
