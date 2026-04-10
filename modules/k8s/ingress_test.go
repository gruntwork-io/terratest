//go:build kubernetes
// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

const ExampleIngressName = "nginx-service-ingress"

func TestGetIngressEReturnsErrorForNonExistantIngress(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetIngressE(t, options, "i-dont-exist")
	require.Error(t, err)
}

func TestGetIngressEReturnsCorrectIngressInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(exampleIngressDeploymentYamlTemplate, uniqueID, uniqueID, uniqueID, uniqueID, uniqueID)

	k8s.KubectlApplyFromString(t, options, configData)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	service := k8s.GetIngress(t, options, "nginx-service-ingress")
	require.Equal(t, "nginx-service-ingress", service.Name)
	require.Equal(t, service.Namespace, uniqueID)
}

func TestListIngressesReturnsCorrectIngressInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(exampleIngressDeploymentYamlTemplate, uniqueID, uniqueID, uniqueID, uniqueID, uniqueID)

	k8s.KubectlApplyFromString(t, options, configData)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	ingresses := k8s.ListIngresses(t, options, metav1.ListOptions{})
	require.Len(t, ingresses, 1)

	ingress := ingresses[0]
	require.Equal(t, ExampleIngressName, ingress.Name)
	require.Equal(t, ingress.Namespace, uniqueID)
}

func TestWaitUntilIngressAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(exampleIngressDeploymentYamlTemplate, uniqueID, uniqueID, uniqueID, uniqueID, uniqueID)

	k8s.KubectlApplyFromString(t, options, configData)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.WaitUntilIngressAvailable(t, options, ExampleIngressName, 60, 5*time.Second)
}

const exampleIngressDeploymentYamlTemplate = `---
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
apiVersion: networking.k8s.io/v1
metadata:
  name: nginx-service-ingress
  namespace: %s
spec:
  rules:
  - http:
      paths:
      - path: /app-%s
        pathType: Prefix
        backend:
          service:
            name: nginx-service
            port:
              number: 80
`
