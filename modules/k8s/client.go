package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

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
	var config *rest.Config

	if options.InClusterAuth {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		logger.Log(t, "Configuring Kubernetes client to use the in-cluster serviceaccount token")
	} else {
		kubeConfigPath, err := options.GetConfigPath(t)
		if err != nil {
			return nil, err
		}
		logger.Logf(t, "Configuring Kubernetes client using config file %s with context %s", kubeConfigPath, options.ContextName)
		// Load API config (instead of more low level ClientConfig)
		config, err = LoadApiClientConfigE(kubeConfigPath, options.ContextName)
		if err != nil {
			logger.Logf(t, "Error loading api client config, falling back to in-cluster authentication via serviceaccount token: %s", err)
			config, err = rest.InClusterConfig()
			if err != nil {
				return nil, err
			}
			logger.Log(t, "Configuring Kubernetes client to use the in-cluster serviceaccount token")
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
