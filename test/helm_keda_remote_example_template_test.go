//go:build kubeall || helm
// +build kubeall helm

// **NOTE**: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly, helm
// can overload the minikube system and thus interfere with the other kubernetes tests. Specifically, many of the tests
// start to fail with `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes
// tests and helm tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.
// We recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

// This file contains an example of how to use terratest to test *remote* helm chart template logic by rendering the templates
// using `helm template`, and then reading in the rendered templates.
// - TestHelmKedaRemoteExampleTemplateRenderedDeployment: An example of how to read in the rendered object and check the
//   computed values.

// An example of how to verify the rendered template object of a Helm Chart given various inputs.
func TestHelmKedaRemoteExampleTemplateRenderedDeployment(t *testing.T) {
	t.Parallel()

	// chart name
	releaseName := "keda"

	// Set up the namespace; confirm that the template renders the expected value for the namespace.
	namespaceName := "medieval-" + strings.ToLower(random.UniqueId())
	logger.Logf(t, "Namespace: %s\n", namespaceName)

	// Setup the args. For this test, we will set the following input values:
	options := &helm.Options{
		SetValues: map[string]string{
			"metricsServer.replicaCount":           "999",
			"resources.metricServer.limits.memory": "1234Mi",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName),
		Logger:         logger.Discard,
	}

	// Run RenderTemplate to render the *remote* template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	// Additionally, we path a the templateFile for which we are setting test values to
	// demonstrate how to select individual templates to render.
	output := helm.RenderRemoteTemplate(t, options, "https://kedacore.github.io/charts", releaseName, []string{"templates/metrics-server/deployment.yaml"})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	// Verify the namespace matches the expected supplied namespace.
	require.Equal(t, namespaceName, deployment.Namespace)

	// Finally, we verify the deployment pod template spec is set to the expected container image value
	var expectedMetricsServerReplica int32
	expectedMetricsServerReplica = 999
	deploymentMetricsServerReplica := *deployment.Spec.Replicas
	require.Equal(t, expectedMetricsServerReplica, deploymentMetricsServerReplica)
	expectedContainerRLM := "1234Mi"
	deploymentContainers := deployment.Spec.Template.Spec.Containers
	require.Equal(t, len(deploymentContainers), 1)
	currentContainerRLM := deploymentContainers[0].Resources.Limits.Memory().String()
	require.Equal(t, currentContainerRLM, expectedContainerRLM)
}

// An example of how to verify the rendered template object of a Helm Chart given input from a `values.yaml` file.
func TestHelmKedaRemoteExampleTemplateRenderedValuesFileFixtureDeployment(t *testing.T) {
	t.Parallel()

	// chart name
	releaseName := "keda"

	// Set up the namespace; confirm that the template renders the expected value for the namespace.
	namespaceName := "medieval-" + strings.ToLower(random.UniqueId())
	logger.Logf(t, "Namespace: %s\n", namespaceName)
	options := &helm.Options{
		ValuesFiles:    []string{"./fixtures/helm/keda-values.yaml"},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespaceName),
		Logger:         logger.Discard,
	}

	// Run RenderTemplate to render the *remote* template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	// Additionally, we path a the templateFile for which we are setting test values to
	// demonstrate how to select individual templates to render.
	output := helm.RenderRemoteTemplate(t, options, "https://kedacore.github.io/charts", releaseName, []string{"templates/metrics-server/deployment.yaml"})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	// Verify the namespace matches the expected supplied namespace.
	require.Equal(t, namespaceName, deployment.Namespace)

	// Finally, we verify the deployment pod template spec is set to the expected value
	var expectedMetricsServerReplica int32
	expectedMetricsServerReplica = 3
	deploymentMetricsServerReplica := *deployment.Spec.Replicas
	require.Equal(t, expectedMetricsServerReplica, deploymentMetricsServerReplica)
	expectedContainerRLM := "1234Mi"
	deploymentContainers := deployment.Spec.Template.Spec.Containers
	require.Equal(t, len(deploymentContainers), 1)
	currentContainerRLM := deploymentContainers[0].Resources.Limits.Memory().String()
	require.Equal(t, currentContainerRLM, expectedContainerRLM)
}
