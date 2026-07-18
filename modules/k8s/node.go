package k8s

import (
	"context"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// GetNodesContextE queries Kubernetes for information about the worker nodes registered to the cluster.
// The ctx parameter supports cancellation and timeouts.
func GetNodesContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions) ([]corev1.Node, error) {
	return GetNodesByFilterContextE(t, ctx, options, metav1.ListOptions{})
}

// GetNodesContext queries Kubernetes for information about the worker nodes registered to the cluster.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetNodesContext(t testing.TestingT, ctx context.Context, options *KubectlOptions) []corev1.Node {
	t.Helper()
	nodes, err := GetNodesContextE(t, ctx, options)
	require.NoError(t, err)

	return nodes
}

// GetNodesByFilterContextE queries Kubernetes for information about the worker nodes registered to the cluster,
// filtering the list of nodes using the provided ListOptions.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func GetNodesByFilterContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filter metav1.ListOptions) ([]corev1.Node, error) {
	options.Logger.Logf(t, "Getting list of nodes from Kubernetes")

	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return nodes.Items, err
}

// GetNodesByFilterContext queries Kubernetes for information about the worker nodes registered to the cluster,
// filtering the list of nodes using the provided ListOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func GetNodesByFilterContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filter metav1.ListOptions) []corev1.Node {
	t.Helper()
	nodes, err := GetNodesByFilterContextE(t, ctx, options, filter)
	require.NoError(t, err)

	return nodes
}

// GetReadyNodesContextE queries Kubernetes for information about the worker nodes registered to the cluster and only
// returns those that are in the ready state.
// The ctx parameter supports cancellation and timeouts.
func GetReadyNodesContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions) ([]corev1.Node, error) {
	nodes, err := GetNodesContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	options.Logger.Logf(t, "Filtering list of nodes from Kubernetes for Ready nodes")

	nodesFiltered := []corev1.Node{}

	for i := range nodes {
		if IsNodeReady(nodes[i]) {
			nodesFiltered = append(nodesFiltered, nodes[i])
		}
	}

	return nodesFiltered, nil
}

// GetReadyNodesContext queries Kubernetes for information about the worker nodes registered to the cluster and only
// returns those that are in the ready state.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetReadyNodesContext(t testing.TestingT, ctx context.Context, options *KubectlOptions) []corev1.Node {
	t.Helper()
	nodes, err := GetReadyNodesContextE(t, ctx, options)
	require.NoError(t, err)

	return nodes
}

// IsNodeReady takes a Kubernetes Node information object and checks if the Node is in the ready state.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func IsNodeReady(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}

	return false
}

// WaitUntilAllNodesReadyContextE continuously polls the Kubernetes cluster until all nodes in the cluster reach the
// ready state, or runs out of retries.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilAllNodesReadyContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, retries int, sleepBetweenRetries time.Duration) error {
	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		"Wait for all Kube Nodes to be ready",
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := AreAllNodesReadyContextE(t, ctx, options)
			if err != nil {
				return "", err
			}

			return "All nodes ready", nil
		},
	)
	options.Logger.Logf(t, "%s", message)

	return err
}

// WaitUntilAllNodesReadyContext continuously polls the Kubernetes cluster until all nodes in the cluster reach the
// ready state, or runs out of retries.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilAllNodesReadyContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilAllNodesReadyContextE(t, ctx, options, retries, sleepBetweenRetries)
	require.NoError(t, err)
}

// AreAllNodesReadyContextE checks if all nodes are ready in the Kubernetes cluster targeted by the current config
// context. The ctx parameter supports cancellation and timeouts.
// If false, returns an error indicating the reason.
func AreAllNodesReadyContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions) (bool, error) {
	nodes, err := GetNodesContextE(t, ctx, options)
	if err != nil {
		return false, err
	}

	if len(nodes) == 0 {
		return false, ErrNoNodesAvailable
	}

	for i := range nodes {
		if !IsNodeReady(nodes[i]) {
			return false, ErrNotAllNodesReady
		}
	}

	return true, nil
}
