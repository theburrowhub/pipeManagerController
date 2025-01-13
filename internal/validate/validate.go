package validate

import (
	"fmt"

	pipemanagerv1alpha1 "github.com/sergiotejon/pipeManagerController/api/v1alpha1"
)

// RunAfterValidation runs after the validation of the Pipeline resource.
// Validate the runAfter field if they are correctly referenced to other tasks in the pipeline.
func RunAfterValidation(pipeline *pipemanagerv1alpha1.Pipeline) error {
	for taskName, task := range pipeline.Spec.Tasks {
		for _, runAfter := range task.RunAfter {
			if _, exists := pipeline.Spec.Tasks[runAfter]; !exists {
				return fmt.Errorf("_ %s in runAfter of _ %s does not exist", runAfter, taskName)
			}
		}
	}

	return nil
}
