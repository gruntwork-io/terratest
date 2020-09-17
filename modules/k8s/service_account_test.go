// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	authv1 "k8s.io/api/authorization/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetServiceAccountWithAuthTokenGetsTokenThatCanBeUsedForAuth(t *testing.T) {
	t.Parallel()

	// make a copy of kubeconfig to namespace it
	tmpConfigPath := CopyHomeKubeConfigToTemp(t)

	// Create a new namespace to work in
	namespaceName := strings.ToLower(random.UniqueId())

	options := NewKubectlOptions("", tmpConfigPath, namespaceName)

	CreateNamespace(t, options, namespaceName)
	defer DeleteNamespace(t, options, namespaceName)

	// Create service account
	serviceAccountName := strings.ToLower(random.UniqueId())
	CreateServiceAccount(t, options, serviceAccountName)
	token := GetServiceAccountAuthToken(t, options, serviceAccountName)
	require.NoError(t, AddConfigContextForServiceAccountE(t, options, serviceAccountName, serviceAccountName, token))

	// Now validate auth as service account. This is a bit tricky because we don't have an API endpoint in k8s that
	// tells you who you are, so we will rely on the self subject access review and see if we have access to the
	// kube-system namespace.
	serviceAccountOptions := NewKubectlOptions(serviceAccountName, tmpConfigPath, namespaceName)
	action := authv1.ResourceAttributes{
		Namespace: "kube-system",
		Verb:      "list",
		Resource:  "pod",
	}
	require.False(t, CanIDo(t, serviceAccountOptions, action))
}

func TestGetServiceAccountEReturnsErrorForNonExistantServiceAccount(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetServiceAccountE(t, options, "terratest")
	require.Error(t, err)
}

func TestGetServiceAccountEReturnsCorrectServiceAccountInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_SERVICEACCOUNT_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	serviceAccount := GetServiceAccount(t, options, "terratest")
	require.Equal(t, serviceAccount.Name, "terratest")
	require.Equal(t, serviceAccount.Namespace, uniqueID)
}

func TestCreateServiceAccountECreatesServiceAccountInNamespaceWithGivenName(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	defer DeleteNamespace(t, options, options.Namespace)
	CreateNamespace(t, options, options.Namespace)

	// Note: We don't need to delete this at the end of test, because deleting the namespace automatically deletes
	// everything created in the namespace.
	CreateServiceAccount(t, options, "terratest")
	serviceAccount := GetServiceAccount(t, options, "terratest")
	require.Equal(t, serviceAccount.Name, "terratest")
	require.Equal(t, serviceAccount.Namespace, uniqueID)
}

const EXAMPLE_SERVICEACCOUNT_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: terratest
  namespace: %s
`
