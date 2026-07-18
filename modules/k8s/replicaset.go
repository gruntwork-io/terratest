package k8s

import (
	"context"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListReplicaSetsContextE looks up replicasets in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListReplicaSetsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.ReplicaSet, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	replicasets, err := clientset.AppsV1().ReplicaSets(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return replicasets.Items, nil
}

// ListReplicaSetsContext looks up replicasets in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListReplicaSetsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []appsv1.ReplicaSet {
	t.Helper()
	replicaset, err := ListReplicaSetsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return replicaset
}

// GetReplicaSetContextE returns a Kubernetes replicaset resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetReplicaSetContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, replicaSetName string) (*appsv1.ReplicaSet, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.AppsV1().ReplicaSets(options.Namespace).Get(ctx, replicaSetName, metav1.GetOptions{})
}

// GetReplicaSetContext returns a Kubernetes replicaset resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetReplicaSetContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, replicaSetName string) *appsv1.ReplicaSet {
	t.Helper()
	replicaset, err := GetReplicaSetContextE(t, ctx, options, replicaSetName)
	require.NoError(t, err)

	return replicaset
}
