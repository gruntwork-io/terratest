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

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetReplicaSetEReturnsError(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	_, err := GetReplicaSetE(t, options, "sample-rs")
	require.Error(t, err)
}

func TestGetReplicaSets(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_REPLICASET_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	replicaSet := GetReplicaSet(t, options, "sample-rs")
	require.Equal(t, replicaSet.Name, "sample-rs")
	require.Equal(t, replicaSet.Namespace, uniqueID)
}

func TestListReplicaSets(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_REPLICASET_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	replicaSets := ListReplicaSets(t, options, metav1.ListOptions{})
	require.Equal(t, len(replicaSets), 1)

	replicaSet := replicaSets[0]
	require.Equal(t, replicaSet.Name, "sample-rs")
	require.Equal(t, replicaSet.Namespace, uniqueID)
}

const EXAMPLE_REPLICASET_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: sample-rs
  namespace: %s
  labels:
    app: sample-rs
spec:
  selector:
    matchLabels:
      name: sample-rs
  template:
    metadata:
      labels:
        name: sample-rs
    spec:
      containers:
      - name: alpine
        image: alpine:3.8
        command: ['sh', '-c', 'echo Hello Terratest! && sleep 99999']
`
