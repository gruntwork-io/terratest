//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
)

func TestGetClusterRoleEReturnsErrorForNonExistantClusterRole(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetClusterRoleE(t, options, "non-existing-role")
	require.Error(t, err)
}

func TestGetClusterRoleEReturnsCorrectClusterRoleInCorrectNamespace(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	defer k8s.KubectlDeleteFromString(t, options, exampleClusterRoleYAMLTemplate)

	k8s.KubectlApplyFromString(t, options, exampleClusterRoleYAMLTemplate)

	role := k8s.GetClusterRole(t, options, "terratest-cluster-role")
	require.Equal(t, "terratest-cluster-role", role.Name)
}

const exampleClusterRoleYAMLTemplate = `---
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
