package k8s

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// AssertAnyNodeHasLabelE will return an error if no nodes are found with the given label.
func AssertAnyNodeHasLabelE(t testing.TestingT, options *KubectlOptions, label string) error {
	nodes, err := GetNodesByFilterE(t, options, metav1.ListOptions{
		LabelSelector: label,
	})

	require.NoError(t, err)
	if len(nodes) == 0 {
		return fmt.Errorf("No nodes found with label %s", label)
	}
	return nil
}

// AssertAnyNodeHasLabel will fail the test if no nodes are found with the given label.
func AssertAnyNodeHasLabel(t testing.TestingT, options *KubectlOptions, label string) {
	err := AssertAnyNodeHasLabelE(t, options, label)
	require.NoError(t, err)
}

func AssertAnyNodeHasTaintE(t testing.TestingT, options *KubectlOptions, taintKey string, taintValue string, taintEffect corev1.TaintEffect) error {
	nodes, err := GetNodesE(t, options)
	require.NoError(t, err)

	for _, node := range nodes {
		for _, taint := range node.Spec.Taints {
			if taint.Key == taintKey && taint.Value == taintValue && taint.Effect == taintEffect {
				return nil
			}
		}
	}

	return fmt.Errorf("No nodes found with taint %s/%s/%s", taintKey, taintValue, taintEffect)
}

// AssertAnyNodeHasTaint will fail the test if no nodes are found with the given taint.
func AssertAnyNodeHasTaint(t testing.TestingT, options *KubectlOptions, taintKey string, taintValue string, taintEffect corev1.TaintEffect) {
	err := AssertAnyNodeHasTaintE(t, options, taintKey, taintValue, taintEffect)
	require.NoError(t, err)
}

// GetNodes queries Kubernetes for information about the worker nodes registered to the cluster. If anything goes wrong,
// the function will automatically fail the test.
func GetNodes(t testing.TestingT, options *KubectlOptions) []corev1.Node {
	nodes, err := GetNodesE(t, options)
	require.NoError(t, err)
	return nodes
}

// GetNodesE queries Kubernetes for information about the worker nodes registered to the cluster.
func GetNodesE(t testing.TestingT, options *KubectlOptions) ([]corev1.Node, error) {
	return GetNodesByFilterE(t, options, metav1.ListOptions{})
}

// GetNodesByFilterE queries Kubernetes for information about the worker nodes registered to the cluster, filtering the
// list of nodes using the provided ListOptions.
func GetNodesByFilterE(t testing.TestingT, options *KubectlOptions, filter metav1.ListOptions) ([]corev1.Node, error) {
	options.Logger.Logf(t, "Getting list of nodes from Kubernetes")

	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return nodes.Items, err
}

// GetReadyNodes queries Kubernetes for information about the worker nodes registered to the cluster and only returns
// those that are in the ready state. If anything goes wrong, the function will automatically fail the test.
func GetReadyNodes(t testing.TestingT, options *KubectlOptions) []corev1.Node {
	nodes, err := GetReadyNodesE(t, options)
	require.NoError(t, err)
	return nodes
}

// GetReadyNodesE queries Kubernetes for information about the worker nodes registered to the cluster and only returns
// those that are in the ready state.
func GetReadyNodesE(t testing.TestingT, options *KubectlOptions) ([]corev1.Node, error) {
	nodes, err := GetNodesE(t, options)
	if err != nil {
		return nil, err
	}
	options.Logger.Logf(t, "Filtering list of nodes from Kubernetes for Ready nodes")
	nodesFiltered := []corev1.Node{}
	for _, node := range nodes {
		if IsNodeReady(node) {
			nodesFiltered = append(nodesFiltered, node)
		}
	}
	return nodesFiltered, nil
}

// IsNodeReady takes a Kubernetes Node information object and checks if the Node is in the ready state.
func IsNodeReady(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

// WaitUntilAllNodesReady continuously polls the Kubernetes cluster until all nodes in the cluster reach the ready
// state, or runs out of retries. Will fail the test immediately if it times out.
func WaitUntilAllNodesReady(t testing.TestingT, options *KubectlOptions, retries int, sleepBetweenRetries time.Duration) {
	err := WaitUntilAllNodesReadyE(t, options, retries, sleepBetweenRetries)
	require.NoError(t, err)
}

// WaitUntilAllNodesReadyE continuously polls the Kubernetes cluster until all nodes in the cluster reach the ready
// state, or runs out of retries.
func WaitUntilAllNodesReadyE(t testing.TestingT, options *KubectlOptions, retries int, sleepBetweenRetries time.Duration) error {
	message, err := retry.DoWithRetryE(
		t,
		"Wait for all Kube Nodes to be ready",
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := AreAllNodesReadyE(t, options)
			if err != nil {
				return "", err
			}
			return "All nodes ready", nil
		},
	)
	options.Logger.Logf(t, message)
	return err
}

// AreAllNodesReady checks if all nodes are ready in the Kubernetes cluster targeted by the current config context
func AreAllNodesReady(t testing.TestingT, options *KubectlOptions) bool {
	nodesReady, _ := AreAllNodesReadyE(t, options)
	return nodesReady
}

// AreAllNodesReadyE checks if all nodes are ready in the Kubernetes cluster targeted by the current config context. If
// false, returns an error indicating the reason.
func AreAllNodesReadyE(t testing.TestingT, options *KubectlOptions) (bool, error) {
	nodes, err := GetNodesE(t, options)
	if err != nil {
		return false, err
	}
	if len(nodes) == 0 {
		return false, errors.New("No nodes available")
	}
	for _, node := range nodes {
		if !IsNodeReady(node) {
			return false, errors.New("Not all nodes ready")
		}
	}
	return true, nil
}
