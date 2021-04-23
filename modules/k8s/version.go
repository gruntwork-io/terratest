package k8s

import "github.com/gruntwork-io/terratest/modules/testing"

// GetKubernetesClusterVersion returns the Kubernetes cluster version.
func GetKubernetesClusterVersion(t testing.TestingT) (string, error) {
	clientset, err := GetKubernetesClientE(t)
	if err != nil {
		return "", err
	}

	versionInfo, err := clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		return "", err
	}

	return versionInfo.String(), nil
}
