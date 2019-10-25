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
)

func TestGetClusterRoleBindingEReturnsErrorForNonExistantClusterRoleBinding(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetClusterRoleBindingE(t, options, "non-existing-cluster-role-binding")
	require.Error(t, err)
}

func TestGetClusterRoleBindingEReturnsCorrectClusterRoleBinding(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_CLUSTER_ROLE_BINDING_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	crBinding := GetClusterRoleBinding(t, options, "terratest-cluster-role-binding")
	require.Equal(t, crBinding.Name, "terratest-cluster-role-binding")
	require.Equal(t, len(crBinding.Subjects), 1)
	require.Equal(t, crBinding.RoleRef.Name, "terratest-cluster-role")
	require.Equal(t, crBinding.RoleRef.Kind, "ClusterRole")
	require.Equal(t, crBinding.Subjects[0].Name, "terratest")
	require.Equal(t, crBinding.Subjects[0].Kind, "ServiceAccount")
}

const EXAMPLE_CLUSTER_ROLE_BINDING_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: '%s'
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: 'terratest'
  namespace: '%s'
---
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
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: 'terratest-cluster-role-binding'
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: 'terratest-cluster-role'
subjects:
- kind: ServiceAccount
  name: 'terratest'
  namespace: '%s'
`
