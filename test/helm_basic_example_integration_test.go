//go:build kubeall || helm
// +build kubeall helm

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly,
// helm can overload the minikube system and thus interfere with the other kubernetes tests. To avoid overloading the
// system, we run the kubernetes tests and helm tests separately from the others.

package test_test

import (
	"crypto/tls"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/helm"
	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// This file contains examples of how to use terratest to test helm charts by deploying the chart and verifying the
// deployment by hitting the service endpoint.
func TestHelmBasicExampleDeployment(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../examples/helm-basic-example")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := "helm-basic-example-" + strings.ToLower(random.UniqueID())

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespaceContext(t, t.Context(), kubectlOptions, namespaceName)
	// ... and make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespaceContext(t, t.Context(), kubectlOptions, namespaceName)

	// Setup the args. For this test, we will set the following input values:
	// - containerImageRepo=nginx
	// - containerImageTag=1.15.8
	options := &helm.Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"containerImageRepo": "nginx",
			"containerImageTag":  "1.15.8",
		},
		ExtraArgs: map[string][]string{
			"install": []string{"--wait", "--timeout", "1m30s"},
		},
	}

	// We generate a unique release name so that we can refer to after deployment.
	// By doing so, we can schedule the delete call here so that at the end of the test, we run
	// `helm delete RELEASE_NAME` to clean up any resources that were created.
	releaseName := "nginx-service-" + strings.ToLower(random.UniqueID())
	defer helm.DeleteContext(t, t.Context(), options, releaseName, true)

	// Deploy the chart using `helm install`. Note that we use the version without `E`, since we want to assert the
	// install succeeds without any errors.
	helm.InstallContext(t, t.Context(), options, helmChartPath, releaseName)

	// Now let's verify the deployment. We will get the service endpoint and try to access it.

	// First we need to get the service name. We will use domain knowledge of the chart here, where the name is
	// RELEASE_NAME-CHART_NAME
	serviceName := releaseName + "-helm-basic-example"

	// Next we wait until the service is available. This will wait up to 10 seconds for the service to become available,
	// to ensure that we can access it.
	k8s.WaitUntilServiceAvailableContext(t, t.Context(), kubectlOptions, serviceName, 10, 1*time.Second)

	// Now we open a tunnel to port forward service port to localhost
	tunnel := k8s.NewTunnel(
		kubectlOptions, k8s.ResourceTypeService, serviceName, 0, 80)
	defer tunnel.Close()

	tunnel.ForwardPort(t)
	// Get endpoint
	endpoint := tunnel.Endpoint()
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
			return statusCode == http.StatusOK
		},
	)
}
