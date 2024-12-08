package k8s

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// LoadKubeConfig try to load the kubeconfig file from the default path.
// if not found, it will try to load the in-cluster config.
func LoadKubeConfig() (*rest.Config, error) {
	kubeconfigPath := clientcmd.RecommendedHomeFile
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		inClusterConfig, inErr := rest.InClusterConfig()
		if inErr != nil {
			return nil, fmt.Errorf("could not load kubeconfig nor in-cluster config: kubeconfigError=%v, inClusterError=%v", err, inErr)
		}
		return inClusterConfig, nil
	}
	return config, nil
}
