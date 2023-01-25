// +build kubeall kubernetes

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

func TestGetStatefulSetEReturnsErrorForNonExistentStatefulSet(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	_, err := GetStatefulSetE(t, options, "sample-statefulset")
	require.Error(t, err)
}

func TestGetStatefulSetEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_STATEFULSET_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	statefulSet := GetStatefulSet(t, options, "sample-statefulset")
	require.Equal(t, "sample-statefulset", statefulSet.Name)
	require.Equal(t, statefulSet.Namespace, uniqueID)
}

func TestListStatefulSetsReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_STATEFULSET_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	statefulSets := ListStatefulSets(t, options, metav1.ListOptions{})
	require.Equal(t, 1, len(statefulSets))

	statefulSet := statefulSets[0]
	require.Equal(t, "sample-statefulset", statefulSet.Name)
	require.Equal(t, uniqueID, statefulSet.Namespace)
}

const EXAMPLE_STATEFULSET_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Service
metadata:
  name: sample-statefulsetservice
  labels:
    app: sample-statefulset
spec:
  ports:
  - port: 80
    name: web
  clusterIP: None
  selector:
    app: sample-statefulset
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: sample-statefulset
  namespace: %s
  labels:
    k8s-app: sample-statefulset
spec:
  serviceName: sample-statefulsetservice
  replicas: 1
  selector:
    matchLabels:
      name: sample-statefulset
  template:
    metadata:
      labels:
        name: sample-statefulset
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: alpine
        image: alpine:3.8
        command: ['sh', '-c', 'echo Hello Terratest! && sleep 99999']
`
