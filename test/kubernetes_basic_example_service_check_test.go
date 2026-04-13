//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test_test

import (
	"crypto/tls"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// An example of how to do more expanded verification of the Kubernetes resource config in examples/kubernetes-basic-example using Terratest.
func TestKubernetesBasicExampleServiceCheck(t *testing.T) {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/kubernetes-basic-example/nginx-deployment.yml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := strings.ToLower(random.UniqueID())

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespaceContext(t, t.Context(), options, namespaceName)
	// ... and make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespaceContext(t, t.Context(), options, namespaceName)

	// At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDeleteContext(t, t.Context(), options, kubeResourcePath)

	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApplyContext(t, t.Context(), options, kubeResourcePath)

	// This will wait up to 10 seconds for the service to become available, to ensure that we can access it.
	k8s.WaitUntilServiceAvailableContext(t, t.Context(), options, "nginx-service", 10, 1*time.Second)

	// Now we verify that the service will successfully boot and start serving requests
	service := k8s.GetServiceContext(t, t.Context(), options, "nginx-service")
	endpoint := k8s.GetServiceEndpointContext(t, t.Context(), options, service, 80)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Test the endpoint for up to 5 minutes. This will only fail if we timeout waiting for the service to return a 200
	// response.
	http_helper.HTTPGetWithRetryWithCustomValidationContext(
		t,
		t.Context(),
		"http://"+endpoint,
		&tlsConfig,
		30,
		10*time.Second,
		func(statusCode int, body string) bool {
			return statusCode == 200
		},
	)
}
