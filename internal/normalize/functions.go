package normalize

import "github.com/sergiotejon/pipeManagerController/api/v1alpha1"

// k8sObjectName returns a valid Kubernetes object name by combining two strings
func k8sObjectName(str1, str2 string) string {
	combined := str1 + "-" + str2
	if len(combined) > 60 {
		return combined[:60]
	}
	return combined
}

// getTaskNames returns a slice of task names from a map of tasks
func getTaskNames(tasks map[string]v1alpha1.Task) []string {
	var taskNames []string
	for name := range tasks {
		taskNames = append(taskNames, name)
	}
	return taskNames
}
