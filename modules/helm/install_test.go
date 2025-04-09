//go:build kubeall || helm
// +build kubeall helm

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly,
// helm can overload the minikube system and thus interfere with the other kubernetes tests. To avoid overloading the
// system, we run the kubernetes tests and helm tests separately from the others.

package helm

import (
	"crypto/tls"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

const (
	remoteChartSource  = "https://charts.bitnami.com/bitnami"
	remoteChartName    = "nginx"
	remoteChartVersion = "13.2.23"
)

// Test that we can install a remote chart (e.g bitnami/nginx)
func TestRemoteChartInstall(t *testing.T) {
	t.Parallel()

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
		Version: remoteChartVersion,
	}

	// Add the stable repo under a random name so as not to touch existing repo configs
	uniqueName := strings.ToLower(fmt.Sprintf("terratest-%s", random.UniqueId()))
	defer RemoveRepo(t, options, uniqueName)
	AddRepo(t, options, uniqueName, remoteChartSource)
	helmChart := fmt.Sprintf("%s/%s", uniqueName, remoteChartName)

	// Generate a unique release name so we can defer the delete before installing
	releaseName := fmt.Sprintf(
		"%s-%s",
		remoteChartName, strings.ToLower(random.UniqueId()),
	)
	defer Delete(t, options, releaseName, true)

	// Test if helm.install will return an error if the chart version is incorrect
	options.Version = "notValidVersion.0.0.0"
	require.Error(t, InstallE(t, options, helmChart, releaseName))

	// Fix chart version and retry install
	options.Version = remoteChartVersion
	// Test that passing extra arguments doesn't error, by changing default timeout
	options.ExtraArgs = map[string][]string{"install": []string{"--timeout", "5m1s"}}
	options.ExtraArgs["delete"] = []string{"--timeout", "5m1s"}
	require.NoError(t, InstallE(t, options, helmChart, releaseName))
	waitForRemoteChartPods(t, kubectlOptions, releaseName, 1)

	// Verify service is accessible. Wait for it to become available and then hit the endpoint.
	serviceName := releaseName
	k8s.WaitUntilServiceAvailable(t, kubectlOptions, serviceName, 10, 1*time.Second)
	service := k8s.GetService(t, kubectlOptions, serviceName)
	endpoint := k8s.GetServiceEndpoint(t, kubectlOptions, service, 80)

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
}

// Test deployment of helm chart with dependencies.
func TestHelmDependencyInstall(t *testing.T) {
	t.Parallel()

	// Path to the helm chart with dependencies which we will test
	helmChartPath, err := filepath.Abs("../../examples/helm-dependency-example")
	require.NoError(t, err)

	// Custom namespace name.
	namespaceName := fmt.Sprintf("helm-dependency-example-%s", strings.ToLower(random.UniqueId()))

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	// Helm chart deployment options.
	options := &Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"containerImageRepo":       "nginx",
			"containerImageTag":        "1.15.8",
			"basic.containerImageRepo": "nginx",
			"basic.containerImageTag":  "1.15.8",
		},
		BuildDependencies: true,
	}
	// We generate a unique release name so that we can refer to after deployment.
	// By doing so, we can schedule the delete call here so that at the end of the test, we run
	// `helm delete RELEASE_NAME` to clean up any resources that were created.
	releaseName := fmt.Sprintf(
		"helm-dependency-example-%s",
		strings.ToLower(random.UniqueId()),
	)
	defer Delete(t, options, releaseName, true)

	// Deploy the chart using `helm install`.
	err = InstallE(t, options, helmChartPath, releaseName)
	assert.NoError(t, err)

	// Verify that Kubernetes service is available after helm chart deployment.
	_, err = k8s.GetServiceE(t, kubectlOptions, releaseName)
	assert.NoError(t, err)
}

func waitForRemoteChartPods(t *testing.T, kubectlOptions *k8s.KubectlOptions, releaseName string, podCount int) {
	// Get pod and wait for it to be avaialable
	// To get the pod, we need to filter it using the labels that the helm chart creates
	filters := metav1.ListOptions{
		LabelSelector: fmt.Sprintf(
			"app.kubernetes.io/name=%s,app.kubernetes.io/instance=%s",
			remoteChartName, releaseName,
		),
	}
	k8s.WaitUntilNumPodsCreated(t, kubectlOptions, filters, podCount, 30, 10*time.Second)
	pods := k8s.ListPods(t, kubectlOptions, filters)
	for _, pod := range pods {
		k8s.WaitUntilPodAvailable(t, kubectlOptions, pod.Name, 30, 10*time.Second)
	}
}
