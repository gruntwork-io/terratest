// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

const ExampleIngressName = "nginx-service-ingress"

func TestGetIngressEReturnsErrorForNonExistantIngress(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetIngressV1Beta1E(t, options, "i-dont-exist")
	require.Error(t, err)
}

func TestGetIngressEReturnsCorrectIngressInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_INGRESS_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	service := GetIngressV1Beta1(t, options, "nginx-service-ingress")
	require.Equal(t, service.Name, "nginx-service-ingress")
	require.Equal(t, service.Namespace, uniqueID)
}

func TestListIngressesReturnsCorrectIngressInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_INGRESS_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	ingresses := ListIngressesV1Beta1(t, options, metav1.ListOptions{})
	require.Equal(t, len(ingresses), 1)

	ingress := ingresses[0]
	require.Equal(t, ingress.Name, ExampleIngressName)
	require.Equal(t, ingress.Namespace, uniqueID)
}

func TestWaitUntilIngressAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_INGRESS_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	WaitUntilIngressAvailableV1Beta1(t, options, ExampleIngressName, 60, 5*time.Second)
}

const EXAMPLE_INGRESS_DEPLOYMENT_YAML_TEMPLATE = `---
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
---
kind: Ingress
apiVersion: extensions/v1beta1
metadata:
  name: nginx-service-ingress
  namespace: %s
spec:
  rules:
  - http:
      paths:
      - path: /app
        backend:
          serviceName: nginx-service
          servicePort: 80
`
