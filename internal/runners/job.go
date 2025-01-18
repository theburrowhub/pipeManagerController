//go:build job

package runners

import (
	v1 "k8s.io/api/batch/v1"

	"github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

type Runner struct {
	Object *v1.Job
	List   *v1.JobList
}

func GetRunner() *Runner {
	return &Runner{
		Object: &v1.Job{},
		List:   &v1.JobList{},
	}
}

func (r *Runner) BuildPipeline(pipelineUID string, spec *v1alpha1.PipelineSpec) error {

	// WIP: This is a work in progress. The goal is to create a pipeline runner using Kubernetes Job.

	return nil
}
