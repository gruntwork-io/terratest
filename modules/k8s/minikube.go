package k8s

import (
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	corev1 "k8s.io/api/core/v1"
)

// IsMinikubeE returns true if the underlying kubernetes cluster is Minikube. This is determined by getting the
// associated nodes and checking if all nodes has at least one label namespaced with "minikube.k8s.io".
func IsMinikubeE(t testing.TestingT, options *KubectlOptions) (bool, error) {
	nodes, err := GetNodesE(t, options)
	if err != nil {
		return false, err
	}

	// ASSUMPTION: All minikube setups will have nodes with labels that are namespaced with minikube.k8s.io
	for _, node := range nodes {
		if !nodeHasMinikubeLabel(node) {
			return false, nil
		}
	}

	// At this point we know that all the nodes in the cluster has the minikube label, so we return true.
	return true, nil
}

// nodeHasMinikubeLabel returns true if any of the labels on the node is namespaced with minikube.k8s.io
func nodeHasMinikubeLabel(node corev1.Node) bool {
	labels := node.GetLabels()
	for key, _ := range labels {
		if strings.HasPrefix(key, "minikube.k8s.io") {
			return true
		}
	}
	return false
}
