package k8s

import (
	"regexp"
	"strings"
)

// ToK8sName normalizes a string to be a valid Kubernetes resource name
func ToK8sName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)
	// Replace invalid characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	name = re.ReplaceAllString(name, "-")
	// Trim leading and trailing hyphens
	name = strings.Trim(name, "-")
	// Truncate to a maximum of 63 characters
	if len(name) > 63 {
		name = name[:63]
	}
	return name
}
