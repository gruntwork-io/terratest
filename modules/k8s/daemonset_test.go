// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.
package k8s

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestGetDaemonSetEReturnsErrorForNonExistantDaemonSet(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	_, err := GetDaemonSetE(t, options, "sample-ds")
	require.Error(t, err)
}

func TestGetDaemonSetEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DAEMONSET_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	daemonSet := GetDaemonSet(t, options, "sample-ds")
	require.Equal(t, daemonSet.Name, "sample-ds")
	require.Equal(t, daemonSet.Namespace, uniqueID)
}

func TestListDaemonSetsReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_DAEMONSET_YAML_TEMPLATE, uniqueID, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	daemonSets := ListDaemonSets(t, options, metav1.ListOptions{})
	require.Equal(t, len(daemonSets), 1)

	daemonSet := daemonSets[0]
	require.Equal(t, daemonSet.Name, "sample-ds")
	require.Equal(t, daemonSet.Namespace, uniqueID)
}

const EXAMPLE_DAEMONSET_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: sample-ds
  namespace: %s
  labels:
    k8s-app: sample-ds
spec:
  selector:
    matchLabels:
      name: sample-ds
  template:
    metadata:
      labels:
        name: sample-ds
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      containers:
      - name: alpine
        image: alpine:3.8
        command: ['sh', '-c', 'echo Hello Terratest! && sleep 99999']
`
