package k8s

import (
	"context"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListStatefulSets will look for statefulSets in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
func ListStatefulSets(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.StatefulSet {
	statefulSet, err := ListStatefulSetsE(t, options, filters)
	require.NoError(t, err)
	return statefulSet
}

// ListStatefulSetsE will look for statefulSets in the given namespace that match the given filters and return them.
func ListStatefulSetsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.StatefulSet, error) {
	clientSet, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	resp, err := clientSet.AppsV1().StatefulSets(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetStatefulSet returns a Kubernetes statefulSet resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetStatefulSet(t testing.TestingT, options *KubectlOptions, statefulSetName string) *appsv1.StatefulSet {
	statefulSet, err := GetStatefulSetE(t, options, statefulSetName)
	require.NoError(t, err)
	return statefulSet
}

// GetStatefulSetsE returns a Kubernetes statefulSet resource in the provided namespace with the given name.
func GetStatefulSetE(t testing.TestingT, options *KubectlOptions, statefulSetName string) (*appsv1.StatefulSet, error) {
	clientSet, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientSet.AppsV1().StatefulSets(options.Namespace).Get(context.Background(), statefulSetName, metav1.GetOptions{})
}
