package k8s

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetKubernetesClusterVersionContextE returns the Kubernetes cluster version.
// The ctx parameter is accepted for API consistency.
func GetKubernetesClusterVersionContextE(t testing.TestingT, ctx context.Context) (string, error) {
	kubeConfigPath, err := GetKubeConfigPathContextE(t, ctx)
	if err != nil {
		return "", err
	}

	options := NewKubectlOptions("", kubeConfigPath, "default")

	return GetKubernetesClusterVersionWithOptionsContextE(t, ctx, options)
}

// GetKubernetesClusterVersionContext returns the Kubernetes cluster version.
// The ctx parameter is accepted for API consistency.
// This will fail the test if there is an error.
func GetKubernetesClusterVersionContext(t testing.TestingT, ctx context.Context) string {
	t.Helper()
	version, err := GetKubernetesClusterVersionContextE(t, ctx)
	require.NoError(t, err)

	return version
}

// GetKubernetesClusterVersionE returns the Kubernetes cluster version.
//
// Deprecated: Use [GetKubernetesClusterVersionContextE] instead.
func GetKubernetesClusterVersionE(t testing.TestingT) (string, error) {
	return GetKubernetesClusterVersionContextE(t, context.Background())
}

// GetKubernetesClusterVersionWithOptionsContextE returns the Kubernetes cluster version given a configured KubectlOptions object.
// The ctx parameter is accepted for API consistency.
func GetKubernetesClusterVersionWithOptionsContextE(t testing.TestingT, ctx context.Context, kubectlOptions *KubectlOptions) (string, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, kubectlOptions)
	if err != nil {
		return "", err
	}

	versionInfo, err := clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		return "", err
	}

	return versionInfo.String(), nil
}

// GetKubernetesClusterVersionWithOptionsContext returns the Kubernetes cluster version given a configured KubectlOptions object.
// The ctx parameter is accepted for API consistency.
// This will fail the test if there is an error.
func GetKubernetesClusterVersionWithOptionsContext(t testing.TestingT, ctx context.Context, kubectlOptions *KubectlOptions) string {
	t.Helper()
	version, err := GetKubernetesClusterVersionWithOptionsContextE(t, ctx, kubectlOptions)
	require.NoError(t, err)

	return version
}

// GetKubernetesClusterVersionWithOptionsE returns the Kubernetes cluster version given a configured KubectlOptions object.
//
// Deprecated: Use [GetKubernetesClusterVersionWithOptionsContextE] instead.
func GetKubernetesClusterVersionWithOptionsE(t testing.TestingT, kubectlOptions *KubectlOptions) (string, error) {
	return GetKubernetesClusterVersionWithOptionsContextE(t, context.Background(), kubectlOptions)
}
