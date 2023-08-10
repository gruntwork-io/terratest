//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"crypto/tls"
	"fmt"
	"strings"
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestTunnelOpensAPortForwardTunnelToPod(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_POD_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)
	WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)

	// Open a tunnel to pod from any available port locally
	tunnel := NewTunnel(options, ResourceTypePod, "nginx-pod", 0, 80)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Try to access the nginx service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", tunnel.Endpoint()),
		&tlsConfig,
		60,
		5*time.Second,
		verifyNginxWelcomePage,
	)
}

func TestTunnelOpensAPortForwardTunnelToDeployment(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleDeploymentYAMLTemplate, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)
	WaitUntilDeploymentAvailable(t, options, "nginx-deployment", 60, 1*time.Second)

	// Open a tunnel to pod from any available port locally
	tunnel := NewTunnel(options, ResourceTypeDeployment, "nginx-deployment", 0, 80)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Try to access the nginx service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", tunnel.Endpoint()),
		&tlsConfig,
		60,
		5*time.Second,
		verifyNginxWelcomePage,
	)
}

func TestTunnelOpensAPortForwardTunnelToService(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExamplePodWithServiceYAMLTemplate, uniqueID, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)
	WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)
	WaitUntilServiceAvailable(t, options, "nginx-service", 60, 1*time.Second)

	// Open a tunnel from any available port locally
	tunnel := NewTunnel(options, ResourceTypeService, "nginx-service", 0, 80)
	defer tunnel.Close()
	tunnel.ForwardPort(t)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Try to access the nginx service on the local port, retrying until we get a good response for up to 5 minutes
	http_helper.HttpGetWithRetryWithCustomValidation(
		t,
		fmt.Sprintf("http://%s", tunnel.Endpoint()),
		&tlsConfig,
		60,
		5*time.Second,
		verifyNginxWelcomePage,
	)
}

func verifyNginxWelcomePage(statusCode int, body string) bool {
	if statusCode != 200 {
		return false
	}
	return strings.Contains(body, "Welcome to nginx")
}

const ExamplePodWithServiceYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  namespace: %s
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx:1.15.7
    ports:
    - containerPort: 80
    readinessProbe:
      httpGet:
        path: /
        port: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
  namespace: %s
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    targetPort: 80
    port: 80
`
