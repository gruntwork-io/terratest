//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestNamespaces(t *testing.T) {
	t.Parallel()

	uniqueId := random.UniqueId()
	namespaceName := strings.ToLower(uniqueId)
	options := NewKubectlOptions("", "", namespaceName)
	CreateNamespace(t, options, namespaceName)
	defer func() {
		DeleteNamespace(t, options, namespaceName)
		namespace := GetNamespace(t, options, namespaceName)
		require.Equal(t, namespace.Status.Phase, corev1.NamespaceTerminating)
	}()

	namespace := GetNamespace(t, options, namespaceName)
	require.Equal(t, namespace.Name, namespaceName)
}

func TestNamespaceWithMetadata(t *testing.T) {
	t.Parallel()

	uniqueId := random.UniqueId()
	namespaceName := strings.ToLower(uniqueId)
	options := NewKubectlOptions("", "", namespaceName)
	namespaceLabels := map[string]string{"foo": "bar"}
	namespaceObjectMetaWithLabels := metav1.ObjectMeta{
		Name:   namespaceName,
		Labels: namespaceLabels,
	}
	CreateNamespaceWithMetadata(t, options, namespaceObjectMetaWithLabels)
	defer func() {
		DeleteNamespace(t, options, namespaceName)
		namespace := GetNamespace(t, options, namespaceName)
		require.Equal(t, namespace.Status.Phase, corev1.NamespaceTerminating)
	}()

	namespace := GetNamespace(t, options, namespaceName)
	require.Equal(t, namespace.Name, namespaceName)
	require.Equal(t, namespace.Labels, namespaceLabels)
}
