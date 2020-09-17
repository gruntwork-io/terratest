// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetClusterRoleEReturnsErrorForNonExistantClusterRole(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetClusterRoleE(t, options, "non-existing-role")
	require.Error(t, err)
}

func TestGetClusterRoleEReturnsCorrectClusterRoleInCorrectNamespace(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	defer KubectlDeleteFromString(t, options, EXAMPLE_CLUSTER_ROLE_YAML_TEMPLATE)
	KubectlApplyFromString(t, options, EXAMPLE_CLUSTER_ROLE_YAML_TEMPLATE)

	role := GetClusterRole(t, options, "terratest-cluster-role")
	require.Equal(t, role.Name, "terratest-cluster-role")
}

const EXAMPLE_CLUSTER_ROLE_YAML_TEMPLATE = `---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: 'terratest-cluster-role'
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
`
