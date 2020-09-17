package k8s

import (
	"k8s.io/client-go/kubernetes"

	// The following line loads the gcp plugin which is required to authenticate against GKE clusters.
	// See: https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetKubernetesClientE returns a Kubernetes API client that can be used to make requests.
func GetKubernetesClientE(t testing.TestingT) (*kubernetes.Clientset, error) {
	kubeConfigPath, err := GetKubeConfigPathE(t)
	if err != nil {
		return nil, err
	}

	options := NewKubectlOptions("", kubeConfigPath, "default")
	return GetKubernetesClientFromOptionsE(t, options)
}

// GetKubernetesClientFromOptionsE returns a Kubernetes API client given a configured KubectlOptions object.
func GetKubernetesClientFromOptionsE(t testing.TestingT, options *KubectlOptions) (*kubernetes.Clientset, error) {
	var err error

	kubeConfigPath, err := options.GetConfigPath(t)
	if err != nil {
		return nil, err
	}
	logger.Logf(t, "Configuring kubectl using config file %s with context %s", kubeConfigPath, options.ContextName)
	// Load API config (instead of more low level ClientConfig)
	config, err := LoadApiClientConfigE(kubeConfigPath, options.ContextName)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
