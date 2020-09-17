package k8s

import (
	"context"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListDaemonSets will look for daemonsets in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
func ListDaemonSets(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.DaemonSet {
	daemonset, err := ListDaemonSetsE(t, options, filters)
	require.NoError(t, err)
	return daemonset
}

// ListDaemonSetsE will look for daemonsets in the given namespace that match the given filters and return them.
func ListDaemonSetsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.DaemonSet, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	resp, err := clientset.AppsV1().DaemonSets(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetDaemonSet returns a Kubernetes daemonset resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetDaemonSet(t testing.TestingT, options *KubectlOptions, daemonSetName string) *appsv1.DaemonSet {
	daemonset, err := GetDaemonSetE(t, options, daemonSetName)
	require.NoError(t, err)
	return daemonset
}

// GetDaemonSetE returns a Kubernetes daemonset resource in the provided namespace with the given name.
func GetDaemonSetE(t testing.TestingT, options *KubectlOptions, daemonSetName string) (*appsv1.DaemonSet, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.AppsV1().DaemonSets(options.Namespace).Get(context.Background(), daemonSetName, metav1.GetOptions{})
}
