//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetReplicaSetEReturnsError(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "")
	_, err := k8s.GetReplicaSetE(t, options, "sample-rs")
	require.Error(t, err)
}

func TestGetReplicaSets(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleReplicaSetYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	replicaSet := k8s.GetReplicaSet(t, options, "sample-rs")
	require.Equal(t, "sample-rs", replicaSet.Name)
	require.Equal(t, replicaSet.Namespace, uniqueID)
}

func TestListReplicaSets(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleReplicaSetYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	replicaSets := k8s.ListReplicaSets(t, options, metav1.ListOptions{})
	require.Len(t, replicaSets, 1)

	replicaSet := replicaSets[0]
	require.Equal(t, "sample-rs", replicaSet.Name)
	require.Equal(t, replicaSet.Namespace, uniqueID)
}

const exampleReplicaSetYAMLTemplate = `---
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
