//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	authv1 "k8s.io/api/authorization/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetServiceAccountWithAuthTokenGetsTokenThatCanBeUsedForAuth(t *testing.T) {
	t.Parallel()

	// make a copy of kubeconfig to namespace it
	tmpConfigPath := k8s.CopyHomeKubeConfigToTemp(t)

	// Create a new namespace to work in
	namespaceName := strings.ToLower(random.UniqueID())

	options := k8s.NewKubectlOptions("", tmpConfigPath, namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)
	defer k8s.DeleteNamespace(t, options, namespaceName)

	// Create service account
	serviceAccountName := strings.ToLower(random.UniqueID())
	k8s.CreateServiceAccount(t, options, serviceAccountName)
	token := k8s.GetServiceAccountAuthToken(t, options, serviceAccountName)
	require.NoError(t, k8s.AddConfigContextForServiceAccountE(t, options, serviceAccountName, serviceAccountName, token))

	// Now validate auth as service account. This is a bit tricky because we don't have an API endpoint in k8s that
	// tells you who you are, so we will rely on the self subject access review and see if we have access to the
	// kube-system namespace.
	serviceAccountOptions := k8s.NewKubectlOptions(serviceAccountName, tmpConfigPath, namespaceName)
	action := authv1.ResourceAttributes{
		Namespace: "kube-system",
		Verb:      "list",
		Resource:  "pod",
	}
	require.False(t, k8s.CanIDo(t, serviceAccountOptions, action))
}

func TestGetServiceAccountEReturnsErrorForNonExistantServiceAccount(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetServiceAccountE(t, options, "terratest")
	require.Error(t, err)
}

func TestGetServiceAccountEReturnsCorrectServiceAccountInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleServiceAccountYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	serviceAccount := k8s.GetServiceAccount(t, options, "terratest")
	require.Equal(t, "terratest", serviceAccount.Name)
	require.Equal(t, serviceAccount.Namespace, uniqueID)
}

func TestCreateServiceAccountECreatesServiceAccountInNamespaceWithGivenName(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())

	options := k8s.NewKubectlOptions("", "", uniqueID)
	defer k8s.DeleteNamespace(t, options, options.Namespace)

	k8s.CreateNamespace(t, options, options.Namespace)

	// Note: We don't need to delete this at the end of test, because deleting the namespace automatically deletes
	// everything created in the namespace.
	k8s.CreateServiceAccount(t, options, "terratest")
	serviceAccount := k8s.GetServiceAccount(t, options, "terratest")
	require.Equal(t, "terratest", serviceAccount.Name)
	require.Equal(t, serviceAccount.Namespace, uniqueID)
}

const exampleServiceAccountYAMLTemplate = `---
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
