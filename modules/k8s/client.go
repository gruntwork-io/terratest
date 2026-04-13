package k8s

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// The following line loads the gcp plugin which is required to authenticate against GKE clusters.
	// See: https://github.com/kubernetes/client-go/issues/242
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetKubernetesClientContextE returns a Kubernetes API client that can be used to make requests.
// The ctx parameter is accepted for API consistency.
func GetKubernetesClientContextE(t testing.TestingT, ctx context.Context) (*kubernetes.Clientset, error) {
	kubeConfigPath, err := GetKubeConfigPathContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	options := NewKubectlOptions("", kubeConfigPath, "default")

	return GetKubernetesClientFromOptionsContextE(t, ctx, options)
}

// GetKubernetesClientContext returns a Kubernetes API client that can be used to make requests.
// The ctx parameter is accepted for API consistency.
// This will fail the test if there is an error.
func GetKubernetesClientContext(t testing.TestingT, ctx context.Context) *kubernetes.Clientset {
	t.Helper()
	clientset, err := GetKubernetesClientContextE(t, ctx)
	require.NoError(t, err)

	return clientset
}

// GetKubernetesClientE returns a Kubernetes API client that can be used to make requests.
//
// Deprecated: Use [GetKubernetesClientContextE] instead.
func GetKubernetesClientE(t testing.TestingT) (*kubernetes.Clientset, error) {
	return GetKubernetesClientContextE(t, context.Background())
}

// GetKubernetesClientFromOptionsContextE returns a Kubernetes API client given a configured KubectlOptions object.
// The ctx parameter is accepted for API consistency.
func GetKubernetesClientFromOptionsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions) (*kubernetes.Clientset, error) { //nolint:contextcheck // GetConfigPath is a method that doesn't accept ctx
	var (
		err    error
		config *rest.Config
	)

	switch {
	case options.InClusterAuth:
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		options.Logger.Logf(t, "Configuring Kubernetes client to use the in-cluster serviceaccount token")
	case options.RestConfig != nil:
		config = options.RestConfig
		options.Logger.Logf(t, "Configuring Kubernetes client to use provided rest config object set with API server address: %s", config.Host)
	default:
		kubeConfigPath, err := options.GetConfigPath(t) //nolint:contextcheck // method doesn't accept ctx
		if err != nil {
			return nil, err
		}

		options.Logger.Logf(t, "Configuring Kubernetes client using config file %s with context %s", kubeConfigPath, options.ContextName)
		// Load API config (instead of more low level ClientConfig)
		config, err = LoadAPIClientConfigE(kubeConfigPath, options.ContextName)
		if err != nil {
			options.Logger.Logf(t, "Error loading api client config, falling back to in-cluster authentication via serviceaccount token: %s", err)

			config, err = rest.InClusterConfig()
			if err != nil {
				return nil, err
			}

			options.Logger.Logf(t, "Configuring Kubernetes client to use the in-cluster serviceaccount token")
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// GetKubernetesClientFromOptionsContext returns a Kubernetes API client given a configured KubectlOptions object.
// The ctx parameter is accepted for API consistency.
// This will fail the test if there is an error.
func GetKubernetesClientFromOptionsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions) *kubernetes.Clientset {
	t.Helper()
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	require.NoError(t, err)

	return clientset
}

// GetKubernetesClientFromOptionsE returns a Kubernetes API client given a configured KubectlOptions object.
//
// Deprecated: Use [GetKubernetesClientFromOptionsContextE] instead.
func GetKubernetesClientFromOptionsE(t testing.TestingT, options *KubectlOptions) (*kubernetes.Clientset, error) {
	return GetKubernetesClientFromOptionsContextE(t, context.Background(), options)
}
