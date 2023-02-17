package k8s

import (
	"context"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListReplicaSets will look for replicasets in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
func ListReplicaSets(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.ReplicaSet {
	replicaset, err := ListReplicaSetsE(t, options, filters)
	require.NoError(t, err)
	return replicaset
}

// ListReplicaSetsE will look for replicasets in the given namespace that match the given filters and return them.
func ListReplicaSetsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.ReplicaSet, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	replicasets, err := clientset.AppsV1().ReplicaSets(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return replicasets.Items, nil
}

// GetReplicaSet returns a Kubernetes replicaset resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetReplicaSet(t testing.TestingT, options *KubectlOptions, replicaSetName string) *appsv1.ReplicaSet {
	replicaset, err := GetReplicaSetE(t, options, replicaSetName)
	require.NoError(t, err)
	return replicaset
}

// GetReplicaSetE returns a Kubernetes replicaset resource in the provided namespace with the given name.
func GetReplicaSetE(t testing.TestingT, options *KubectlOptions, replicaSetName string) (*appsv1.ReplicaSet, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.AppsV1().ReplicaSets(options.Namespace).Get(context.Background(), replicaSetName, metav1.GetOptions{})
}
