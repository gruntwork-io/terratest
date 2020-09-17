// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	authv1 "k8s.io/api/authorization/v1"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

// An example of how to test the Kubernetes resource config in examples/kubernetes-rbac-example using Terratest,
// including whether or not the permissions are set correctly.
func TestKubernetesRBACExample(t *testing.T) {
	t.Parallel()

	// These are pulled from the kubernetes resource config
	const serviceAccountName = "terratest-rbac-example-service-account"
	const namespaceName = "terratest-rbac-example-namespace"

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/kubernetes-rbac-example/namespace-service-account.yml")
	require.NoError(t, err)

	// Setup the kubectl config and context. Here we choose to create a new one because we will be manipulating the
	// entries to be able to add a new authentication option.
	tmpConfigPath := k8s.CopyHomeKubeConfigToTemp(t)
	defer os.Remove(tmpConfigPath)
	options := k8s.NewKubectlOptions("", tmpConfigPath, namespaceName)

	// At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, options, kubeResourcePath)

	// Retrieve authentication token for the newly created ServiceAccount
	token := k8s.GetServiceAccountAuthToken(t, options, serviceAccountName)

	// Now update the configuration to add a new context that can be used to make requests as that service account
	require.NoError(t, k8s.AddConfigContextForServiceAccountE(
		t,
		options,
		serviceAccountName, // for this test we will name the context after the ServiceAccount
		serviceAccountName,
		token,
	))
	serviceAccountKubectlOptions := k8s.NewKubectlOptions(serviceAccountName, tmpConfigPath, namespaceName)

	// At this point all requests made with serviceAccountKubectlOptions will be auth'd as that ServiceAccount. So let's
	// verify that! We will check:
	// - we can't access the kube-system namespace
	adminListPodAction := authv1.ResourceAttributes{
		Namespace: "kube-system",
		Verb:      "list",
		Resource:  "pod",
	}
	require.False(t, k8s.CanIDo(t, serviceAccountKubectlOptions, adminListPodAction))
	// - we can access the namespace the service account is in
	namespaceListPodAction := authv1.ResourceAttributes{
		Namespace: namespaceName,
		Verb:      "list",
		Resource:  "pod",
	}
	require.True(t, k8s.CanIDo(t, serviceAccountKubectlOptions, namespaceListPodAction))
}
