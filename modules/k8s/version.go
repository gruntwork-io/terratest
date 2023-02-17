package k8s

import "github.com/gruntwork-io/terratest/modules/testing"

// GetKubernetesClusterVersion returns the Kubernetes cluster version.
func GetKubernetesClusterVersionE(t testing.TestingT) (string, error) {
	kubeConfigPath, err := GetKubeConfigPathE(t)
	if err != nil {
		return "", err
	}

	options := NewKubectlOptions("", kubeConfigPath, "default")

	return GetKubernetesClusterVersionWithOptionsE(t, options)
}

// GetKubernetesClusterVersion returns the Kubernetes cluster version given a configured KubectlOptions object.
func GetKubernetesClusterVersionWithOptionsE(t testing.TestingT, kubectlOptions *KubectlOptions) (string, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, kubectlOptions)
	if err != nil {
		return "", err
	}

	versionInfo, err := clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		return "", err
	}

	return versionInfo.String(), nil
}
