// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests that:
// - kubectl is properly configured to talk to a kubernetes cluster
// - GetNodes will return a list of nodes registered with kubernetes
func TestGetNodes(t *testing.T) {
	t.Parallel()

	// Assumes local kubernetes (minikube or docker-for-desktop kube), where there is only one node
	options := NewKubectlOptions("", "", "default")
	nodes := GetNodes(t, options)
	require.Equal(t, len(nodes), 1)

	node := nodes[0]
	// Make sure node name is not blank, indicating an uninitialized Node object
	assert.NotEqual(t, node.Name, "")
}

// Tests that:
// - kubectl is properly configured to talk to a kubernetes cluster
// - GetReadyNodes will return a list of ready nodes registered with kubernetes
func TestGetReadyNodes(t *testing.T) {
	t.Parallel()

	// Assumes local kubernetes (minikube or docker-for-desktop kube), where there is only one node
	options := NewKubectlOptions("", "", "default")
	nodes := GetReadyNodes(t, options)
	require.Equal(t, len(nodes), 1)

	node := nodes[0]
	// Make sure node name is not blank, indicating an uninitialized Node object
	assert.NotEqual(t, node.Name, "")
}

// Tests that:
// - kubectl is properly configured to talk to a kubernetes cluster
// - WaitUntilAllNodesReady checks if all nodes in the cluster are ready
func TestWaitUntilAllNodesReady(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")

	WaitUntilAllNodesReady(t, options, 12, 5*time.Second)

	nodes := GetNodes(t, options)
	nodeNames := map[string]bool{}
	for _, node := range nodes {
		nodeNames[node.Name] = true
	}

	readyNodes := GetReadyNodes(t, options)
	readyNodeNames := map[string]bool{}
	for _, node := range readyNodes {
		readyNodeNames[node.Name] = true
	}

	assert.Equal(t, nodeNames, readyNodeNames)
}
