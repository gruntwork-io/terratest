package k8s //nolint:dupl // structural pattern for k8s resource operations

import (
	"context"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListDaemonSetsContextE looks up daemonsets in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDaemonSetsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.DaemonSet, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.AppsV1().DaemonSets(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListDaemonSetsContext looks up daemonsets in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDaemonSetsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []appsv1.DaemonSet {
	t.Helper()
	daemonset, err := ListDaemonSetsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return daemonset
}

// ListDaemonSets will look for daemonsets in the given namespace that match the given filters and return them. This will
// fail the test if there is an error.
//
// Deprecated: Use [ListDaemonSetsContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDaemonSets(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []appsv1.DaemonSet {
	t.Helper()

	return ListDaemonSetsContext(t, context.Background(), options, filters)
}

// ListDaemonSetsE will look for daemonsets in the given namespace that match the given filters and return them.
//
// Deprecated: Use [ListDaemonSetsContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListDaemonSetsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]appsv1.DaemonSet, error) {
	return ListDaemonSetsContextE(t, context.Background(), options, filters)
}

// GetDaemonSetContextE returns a Kubernetes daemonset resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetDaemonSetContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, daemonSetName string) (*appsv1.DaemonSet, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.AppsV1().DaemonSets(options.Namespace).Get(ctx, daemonSetName, metav1.GetOptions{})
}

// GetDaemonSetContext returns a Kubernetes daemonset resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetDaemonSetContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, daemonSetName string) *appsv1.DaemonSet {
	t.Helper()
	daemonset, err := GetDaemonSetContextE(t, ctx, options, daemonSetName)
	require.NoError(t, err)

	return daemonset
}

// GetDaemonSet returns a Kubernetes daemonset resource in the provided namespace with the given name. This will
// fail the test if there is an error.
//
// Deprecated: Use [GetDaemonSetContext] instead.
func GetDaemonSet(t testing.TestingT, options *KubectlOptions, daemonSetName string) *appsv1.DaemonSet {
	t.Helper()

	return GetDaemonSetContext(t, context.Background(), options, daemonSetName)
}

// GetDaemonSetE returns a Kubernetes daemonset resource in the provided namespace with the given name.
//
// Deprecated: Use [GetDaemonSetContextE] instead.
func GetDaemonSetE(t testing.TestingT, options *KubectlOptions, daemonSetName string) (*appsv1.DaemonSet, error) {
	return GetDaemonSetContextE(t, context.Background(), options, daemonSetName)
}
