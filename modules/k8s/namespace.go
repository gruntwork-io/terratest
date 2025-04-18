package k8s

import (
	"context"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNamespace will create a new Kubernetes namespace on the cluster targeted by the provided options. This will
// fail the test if there is an error in creating the namespace.
func CreateNamespace(t testing.TestingT, options *KubectlOptions, namespaceName string) {
	require.NoError(t, CreateNamespaceE(t, options, namespaceName))
}

// CreateNamespaceE will create a new Kubernetes namespace on the cluster targeted by the provided options.
func CreateNamespaceE(t testing.TestingT, options *KubectlOptions, namespaceName string) error {
	namespaceObject := metav1.ObjectMeta{
		Name: namespaceName,
	}
	return CreateNamespaceWithMetadataE(t, options, namespaceObject)
}

// CreateNamespaceWithMetadataE will create a new Kubernetes namespace on the cluster targeted by the provided options and
// with the provided metadata. This method expects the entire namespace ObjectMeta to be passed in, so you'll need to set the name within the ObjectMeta struct yourself.
func CreateNamespaceWithMetadataE(t testing.TestingT, options *KubectlOptions, namespaceObjectMeta metav1.ObjectMeta) error {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return err
	}

	namespace := corev1.Namespace{
		ObjectMeta: namespaceObjectMeta,
	}
	_, err = clientset.CoreV1().Namespaces().Create(context.Background(), &namespace, metav1.CreateOptions{})
	return err
}

// CreateNamespaceWithMetadata will create a new Kubernetes namespace on the cluster targeted by the provided options and
// with the provided metadata. This method expects the entire namespace ObjectMeta to be passed in, so you'll need to set the name within the ObjectMeta struct yourself.
// This will fail the test if there is an error while creating the namespace.
func CreateNamespaceWithMetadata(t testing.TestingT, options *KubectlOptions, namespaceObjectMeta metav1.ObjectMeta) {
	require.NoError(t, CreateNamespaceWithMetadataE(t, options, namespaceObjectMeta))
}

// GetNamespace will query the Kubernetes cluster targeted by the provided options for the requested namespace. This will
// fail the test if there is an error in getting the namespace or if the namespace doesn't exist.
func GetNamespace(t testing.TestingT, options *KubectlOptions, namespaceName string) *corev1.Namespace {
	namespace, err := GetNamespaceE(t, options, namespaceName)
	require.NoError(t, err)
	require.NotNil(t, namespace)
	return namespace
}

// GetNamespaceE will query the Kubernetes cluster targeted by the provided options for the requested namespace.
func GetNamespaceE(t testing.TestingT, options *KubectlOptions, namespaceName string) (*corev1.Namespace, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Namespaces().Get(context.Background(), namespaceName, metav1.GetOptions{})
}

// DeleteNamespace will delete the requested namespace from the Kubernetes cluster targeted by the provided options. This will
// fail the test if there is an error in creating the namespace.
func DeleteNamespace(t testing.TestingT, options *KubectlOptions, namespaceName string) {
	require.NoError(t, DeleteNamespaceE(t, options, namespaceName))
}

// DeleteNamespaceE will delete the requested namespace from the Kubernetes cluster targeted by the provided options.
func DeleteNamespaceE(t testing.TestingT, options *KubectlOptions, namespaceName string) error {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return err
	}

	return clientset.CoreV1().Namespaces().Delete(context.Background(), namespaceName, metav1.DeleteOptions{})
}

// ListNamespaces will list all namespaces in the Kubernetes cluster targeted by the provided options.
// This will fail the test if there is an error in listing the namespaces.
func ListNamespaces(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Namespace {
	namespaces, err := ListNamespacesE(t, options, filters)
	require.NoError(t, err)

	if len(namespaces) > 0 {
		var namespaceNames []string
		for _, ns := range namespaces {
			namespaceNames = append(namespaceNames, ns.Name)
		}
		options.Logger.Logf(t, "Found namespaces: %s", strings.Join(namespaceNames, ", "))
	} else {
		options.Logger.Logf(t, "No namespaces found matching the provided filters.")
	}

	return namespaces
}

// ListNamespacesE lists all namespaces in the Kubernetes cluster and returns them or an error.
func ListNamespacesE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Namespace, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	namespaceList, err := clientset.CoreV1().Namespaces().List(context.Background(), filters)
	if err != nil {
		return nil, err
	}

	return namespaceList.Items, nil
}
