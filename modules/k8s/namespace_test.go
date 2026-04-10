//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestNamespaces(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	namespaceName := strings.ToLower(uniqueID)
	options := k8s.NewKubectlOptions("", "", namespaceName)
	k8s.CreateNamespace(t, options, namespaceName)

	defer func() {
		k8s.DeleteNamespace(t, options, namespaceName)
		namespace := k8s.GetNamespace(t, options, namespaceName)
		require.Equal(t, corev1.NamespaceTerminating, namespace.Status.Phase)
	}()

	namespace := k8s.GetNamespace(t, options, namespaceName)
	require.Equal(t, namespace.Name, namespaceName)
}

func TestNamespaceWithMetadata(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	namespaceName := strings.ToLower(uniqueID)
	options := k8s.NewKubectlOptions("", "", namespaceName)
	namespaceLabels := map[string]string{"foo": "bar"}
	namespaceObjectMetaWithLabels := metav1.ObjectMeta{
		Name:   namespaceName,
		Labels: namespaceLabels,
	}
	k8s.CreateNamespaceWithMetadata(t, options, namespaceObjectMetaWithLabels)

	defer func() {
		k8s.DeleteNamespace(t, options, namespaceName)
		namespace := k8s.GetNamespace(t, options, namespaceName)
		require.Equal(t, corev1.NamespaceTerminating, namespace.Status.Phase)
	}()

	namespace := k8s.GetNamespace(t, options, namespaceName)
	require.Equal(t, namespace.Name, namespaceName)

	for k, v := range namespaceLabels {
		require.Equal(t, v, namespace.Labels[k], "Expected label %s=%s", k, v)
	}
}

func TestListNamespaces(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	namespaceName := strings.ToLower(uniqueID)
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)
	t.Cleanup(func() { k8s.DeleteNamespace(t, options, namespaceName) })

	t.Run("List all namespaces and find the created one", func(t *testing.T) {
		t.Parallel()
		namespaces := k8s.ListNamespaces(t, options, metav1.ListOptions{})
		require.NotEmpty(t, namespaces, "Should find at least some namespaces")

		found := false

		for _, ns := range namespaces {
			if ns.Name == namespaceName {
				found = true
				break
			}
		}

		require.True(t, found, "Should find the created namespace in the list")
	})
}
