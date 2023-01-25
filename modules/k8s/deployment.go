// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.
package k8s

import (
	"context"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListDeployments will look for deployments in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
func ListDeployments(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.Deployment {
	deployment, err := ListDeploymentsE(t, options, filters)
	require.NoError(t, err)
	return deployment
}

// ListDeploymentsE will look for deployments in the given namespace that match the given filters and return them.
func ListDeploymentsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.Deployment, error) {
	clientSet, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	resp, err := clientSet.AppsV1().Deployments(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetDeployment returns a Kubernetes deployment resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetDeployment(t testing.TestingT, options *KubectlOptions, deploymentName string) *appsv1.Deployment {
	deployment, err := GetDeploymentE(t, options, deploymentName)
	require.NoError(t, err)
	return deployment
}

// GetDeploymentE returns a Kubernetes deployment resource in the provided namespace with the given name.
func GetDeploymentE(t testing.TestingT, options *KubectlOptions, deploymentName string) (*appsv1.Deployment, error) {
	clientSet, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientSet.AppsV1().Deployments(options.Namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
}
