// +build kubeall helm

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly,
// helm can overload the minikube system and thus interfere with the other kubernetes tests. To avoid overloading the
// system, we run the kubernetes tests and helm tests separately from the others.

package helm

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Test that we can install, upgrade, and rollback a remote chart (e.g stable/chartmuseum)
func TestRemoteChartInstallUpgradeRollback(t *testing.T) {
	t.Parallel()

	helmChart := "stable/chartmuseum"

	namespaceName := fmt.Sprintf(
		"%s-%s",
		strings.ToLower(t.Name()),
		strings.ToLower(random.UniqueId()),
	)

	// Use default kubectl options to create a new namespace for this test, and then update the namespace for kubectl
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)

	// Override service type to node port
	options := &Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"service.type": "NodePort",
		},
	}

	// Generate a unique release name so we can defer the delete before installing
	releaseName := fmt.Sprintf(
		"chartmuseum-%s",
		strings.ToLower(random.UniqueId()),
	)
	defer Delete(t, options, releaseName, true)
	Install(t, options, helmChart, releaseName)
	waitForRemoteChartPods(t, kubectlOptions, releaseName, 1)

	// Setting replica count to 2 to check the upgrade functionality.
	// After successful upgrade , the count of pods should be equal to 2.
	options.SetValues = map[string]string{
		"replicaCount": "2",
		"service.type": "NodePort",
	}
	Upgrade(t, options, helmChart, releaseName)
	waitForRemoteChartPods(t, kubectlOptions, releaseName, 2)

	// Verify service is accessible. Wait for it to become available and then hit the endpoint.
	// Service name is RELEASE_NAME-CHART_NAME
	serviceName := fmt.Sprintf("%s-chartmuseum", releaseName)
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, serviceName, 10, 1*time.Second)
	service := k8s.GetService(t, kubectlOptions, serviceName)
	endpoint := k8s.GetServiceEndpoint(t, kubectlOptions, service, 8080)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", endpoint),
		&tlsConfig,
		30,
		10*time.Second,
		func(statusCode int, body string) bool {
			return statusCode == 200
		},
	)

	// Finally, test rollback functionality. When rolling back, we should see the pods go back down to 1.
	Rollback(t, options, releaseName, "")
	waitForRemoteChartPods(t, kubectlOptions, releaseName, 1)
}
