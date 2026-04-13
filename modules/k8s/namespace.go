package k8s

import (
	"context"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
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

// CreateNamespace will create a new Kubernetes namespace on the cluster targeted by the provided options. This will
// fail the test if there is an error in creating the namespace.
//
// Deprecated: Use [CreateNamespaceContext] instead.
func CreateNamespace(t testing.TestingT, options *KubectlOptions, namespaceName string) {
	t.Helper()
	CreateNamespaceContext(t, context.Background(), options, namespaceName)
}

// CreateNamespaceE will create a new Kubernetes namespace on the cluster targeted by the provided options.
//
// Deprecated: Use [CreateNamespaceContextE] instead.
func CreateNamespaceE(t testing.TestingT, options *KubectlOptions, namespaceName string) error {
	return CreateNamespaceContextE(t, context.Background(), options, namespaceName)
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

// CreateNamespaceWithMetadataE will create a new Kubernetes namespace on the cluster targeted by the provided options and
// with the provided metadata. This method expects the entire namespace ObjectMeta to be passed in, so you'll need to set the name within the ObjectMeta struct yourself.
//
// Deprecated: Use [CreateNamespaceWithMetadataContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CreateNamespaceWithMetadataE(t testing.TestingT, options *KubectlOptions, namespaceObjectMeta metav1.ObjectMeta) error {
	return CreateNamespaceWithMetadataContextE(t, context.Background(), options, namespaceObjectMeta)
}

// CreateNamespaceWithMetadata will create a new Kubernetes namespace on the cluster targeted by the provided options and
// with the provided metadata. This method expects the entire namespace ObjectMeta to be passed in, so you'll need to set the name within the ObjectMeta struct yourself.
// This will fail the test if there is an error while creating the namespace.
//
// Deprecated: Use [CreateNamespaceWithMetadataContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func CreateNamespaceWithMetadata(t testing.TestingT, options *KubectlOptions, namespaceObjectMeta metav1.ObjectMeta) {
	t.Helper()
	CreateNamespaceWithMetadataContext(t, context.Background(), options, namespaceObjectMeta)
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

// GetNamespace will query the Kubernetes cluster targeted by the provided options for the requested namespace. This will
// fail the test if there is an error in getting the namespace or if the namespace doesn't exist.
//
// Deprecated: Use [GetNamespaceContext] instead.
func GetNamespace(t testing.TestingT, options *KubectlOptions, namespaceName string) *corev1.Namespace {
	t.Helper()

	return GetNamespaceContext(t, context.Background(), options, namespaceName)
}

// GetNamespaceE will query the Kubernetes cluster targeted by the provided options for the requested namespace.
//
// Deprecated: Use [GetNamespaceContextE] instead.
func GetNamespaceE(t testing.TestingT, options *KubectlOptions, namespaceName string) (*corev1.Namespace, error) {
	return GetNamespaceContextE(t, context.Background(), options, namespaceName)
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

// DeleteNamespace will delete the requested namespace from the Kubernetes cluster targeted by the provided options. This will
// fail the test if there is an error in creating the namespace.
//
// Deprecated: Use [DeleteNamespaceContext] instead.
func DeleteNamespace(t testing.TestingT, options *KubectlOptions, namespaceName string) {
	t.Helper()
	DeleteNamespaceContext(t, context.Background(), options, namespaceName)
}

// DeleteNamespaceE will delete the requested namespace from the Kubernetes cluster targeted by the provided options.
//
// Deprecated: Use [DeleteNamespaceContextE] instead.
func DeleteNamespaceE(t testing.TestingT, options *KubectlOptions, namespaceName string) error {
	return DeleteNamespaceContextE(t, context.Background(), options, namespaceName)
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

// ListNamespaces will list all namespaces in the Kubernetes cluster targeted by the provided options.
// This will fail the test if there is an error in listing the namespaces.
//
// Deprecated: Use [ListNamespacesContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListNamespaces(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Namespace {
	t.Helper()

	return ListNamespacesContext(t, context.Background(), options, filters)
}

// ListNamespacesE lists all namespaces in the Kubernetes cluster and returns them or an error.
//
// Deprecated: Use [ListNamespacesContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListNamespacesE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Namespace, error) {
	return ListNamespacesContextE(t, context.Background(), options, filters)
}
