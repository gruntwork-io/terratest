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

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetDeploymentEReturnsErrorForNonExistentDeployment(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	_, err := GetDeploymentE(t, options, "sample-deployment")
	require.Error(t, err)
}

func TestGetDeploymentEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	deployment := GetDeployment(t, options, "sample-deployment")
	require.Equal(t, "sample-deployment", deployment.Name)
	require.Equal(t, uniqueID, deployment.Namespace)
}

func TestListDeploymentsReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DEPLOYMENT_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	deployments := ListDeployments(t, options, metav1.ListOptions{})
	require.Equal(t, 1, len(deployments))

	deployment := deployments[0]
	require.Equal(t, "sample-deployment", deployment.Name)
	require.Equal(t, uniqueID, deployment.Namespace)
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
  name: sample-deployment
  namespace: %s
  labels:
    k8s-app: sample-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      name: sample-deployment
  template:
    metadata:
      labels:
        name: sample-deployment
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: alpine
        image: alpine:3.8
        command: ['sh', '-c', 'echo Hello Terratest! && sleep 99999']
`
