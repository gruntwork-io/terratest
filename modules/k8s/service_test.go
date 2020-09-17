// +build kubernetes

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

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetServiceEReturnsErrorForNonExistantService(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetServiceE(t, options, "nginx-service")
	require.Error(t, err)
}

func TestGetServiceEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	service := GetService(t, options, "nginx-service")
	require.Equal(t, service.Name, "nginx-service")
	require.Equal(t, service.Namespace, uniqueID)
}

func TestListServicesReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	services := ListServices(t, options, metav1.ListOptions{})
	require.Equal(t, len(services), 1)

	service := services[0]
	require.Equal(t, service.Name, "nginx-service")
	require.Equal(t, service.Namespace, uniqueID)
}

func TestWaitUntilServiceAvailableReturnsSuccessfullyOnNodePortType(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	WaitUntilServiceAvailable(t, options, "nginx-service", 10, 1*time.Second)
}

func TestGetServiceEndpointEReturnsAccessibleEndpointForNodePort(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	service := GetService(t, options, "nginx-service")
	endpoint := GetServiceEndpoint(t, options, service, 80)

	// Setup a TLS configuration to submit with the helper, a blank struct is acceptable
	tlsConfig := tls.Config{}

	// Test up to 5 minutes
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

const EXAMPLE_DEPLOYMENT_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: %s
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.15.7
        ports:
        - containerPort: 80
---
kind: Service
apiVersion: v1
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
  type: NodePort
`
