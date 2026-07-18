package k8s

import (
	"context"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNamespaceContextE will create a new Kubernetes namespace on the cluster targeted by the provided options.
// The ctx parameter supports cancellation and timeouts.
func CreateNamespaceContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceName string) error {
	namespaceObject := metav1.ObjectMeta{
		Name: namespaceName,
	}

	return CreateNamespaceWithMetadataContextE(t, ctx, options, namespaceObject)
}

// CreateNamespaceContext will create a new Kubernetes namespace on the cluster targeted by the provided options.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error in creating the namespace.
func CreateNamespaceContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceName string) {
	t.Helper()
	require.NoError(t, CreateNamespaceContextE(t, ctx, options, namespaceName))
}

// CreateNamespaceWithMetadataContextE will create a new Kubernetes namespace on the cluster targeted by the provided
// options and with the provided metadata.
// The ctx parameter supports cancellation and timeouts.
// This method expects the entire namespace ObjectMeta to be passed in, so you'll need to set the name within the ObjectMeta struct yourself.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CreateNamespaceWithMetadataContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceObjectMeta metav1.ObjectMeta) error {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return err
	}

	namespace := corev1.Namespace{
		ObjectMeta: namespaceObjectMeta,
	}
	_, err = clientset.CoreV1().Namespaces().Create(ctx, &namespace, metav1.CreateOptions{})

	return err
}

// CreateNamespaceWithMetadataContext will create a new Kubernetes namespace on the cluster targeted by the provided
// options and with the provided metadata.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error while creating the namespace.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CreateNamespaceWithMetadataContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceObjectMeta metav1.ObjectMeta) {
	t.Helper()
	require.NoError(t, CreateNamespaceWithMetadataContextE(t, ctx, options, namespaceObjectMeta))
}

// GetNamespaceContextE will query the Kubernetes cluster targeted by the provided options for the requested namespace.
// The ctx parameter supports cancellation and timeouts.
func GetNamespaceContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceName string) (*corev1.Namespace, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
}

// GetNamespaceContext will query the Kubernetes cluster targeted by the provided options for the requested namespace.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error or if the namespace doesn't exist.
func GetNamespaceContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceName string) *corev1.Namespace {
	t.Helper()
	namespace, err := GetNamespaceContextE(t, ctx, options, namespaceName)
	require.NoError(t, err)
	require.NotNil(t, namespace)

	return namespace
}

// DeleteNamespaceContextE will delete the requested namespace from the Kubernetes cluster targeted by the provided options.
// The ctx parameter supports cancellation and timeouts.
func DeleteNamespaceContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceName string) error {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return err
	}

	return clientset.CoreV1().Namespaces().Delete(ctx, namespaceName, metav1.DeleteOptions{})
}

// DeleteNamespaceContext will delete the requested namespace from the Kubernetes cluster targeted by the provided options.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func DeleteNamespaceContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, namespaceName string) {
	t.Helper()
	require.NoError(t, DeleteNamespaceContextE(t, ctx, options, namespaceName))
}

// ListNamespacesContextE lists all namespaces in the Kubernetes cluster that match the given filters and returns them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListNamespacesContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Namespace, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	namespaceList, err := clientset.CoreV1().Namespaces().List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return namespaceList.Items, nil
}

// ListNamespacesContext lists all namespaces in the Kubernetes cluster that match the given filters and returns them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListNamespacesContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []corev1.Namespace {
	t.Helper()
	namespaces, err := ListNamespacesContextE(t, ctx, options, filters)
	require.NoError(t, err)

	if len(namespaces) > 0 {
		namespaceNames := make([]string, 0, len(namespaces))
		for _, ns := range namespaces {
			namespaceNames = append(namespaceNames, ns.Name)
		}

		options.Logger.Logf(t, "Found namespaces: %s", strings.Join(namespaceNames, ", "))
	} else {
		options.Logger.Logf(t, "No namespaces found matching the provided filters.")
	}

	return namespaces
}
