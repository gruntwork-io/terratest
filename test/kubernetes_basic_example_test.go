// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// An example of how to test the Kubernetes resource config in examples/kubernetes-basic-example using Terratest.
func TestKubernetesBasicExample(t *testing.T) {
	t.Parallel()

	// website::tag::1::Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/kubernetes-basic-example/nginx-deployment.yml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("kubernetes-basic-example-%s", strings.ToLower(random.UniqueId()))

	// website::tag::2::Setup the kubectl config and context.
	// Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	// - Random namespace
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)
	// website::tag::5::Make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespace(t, options, namespaceName)

	// website::tag::6::At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	// website::tag::3::Apply kubectl with 'kubectl apply -f RESOURCE_CONFIG' command.
	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, options, kubeResourcePath)

	// website::tag::4::Check if NGINX service was deployed successfully.
	// This will get the service resource and verify that it exists and was retrieved successfully. This function will
	// fail the test if the there is an error retrieving the service resource from Kubernetes.
	service := k8s.GetService(t, options, "nginx-service")
	require.Equal(t, service.Name, "nginx-service")
}
