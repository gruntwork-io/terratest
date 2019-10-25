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

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetRoleBindingEReturnsErrorForNonExistantRoleBinding(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetRoleBindingE(t, options, "non-existing-role-binding")
	require.Error(t, err)
}

func TestGetRoleBindingEReturnsCorrectRoleBindingInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_ROLE_BINDING_YAML_TEMPLATE, uniqueID, uniqueID, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	roleBinding := GetRoleBinding(t, options, "terratest-role-binding")
	require.Equal(t, roleBinding.Name, "terratest-role-binding")
	require.Equal(t, roleBinding.Namespace, uniqueID)
	require.Equal(t, roleBinding.RoleRef.Name, "terratest-role")
	require.Equal(t, roleBinding.RoleRef.Kind, "Role")
	require.Equal(t, len(roleBinding.Subjects), 1)
	require.Equal(t, roleBinding.Subjects[0].Name, "terratest")
	require.Equal(t, roleBinding.Subjects[0].Kind, "ServiceAccount")
}

const EXAMPLE_ROLE_BINDING_YAML_TEMPLATE = `---
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
kind: Role
metadata:
  name: 'terratest-role'
  namespace: '%s'
rules:
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - '*'
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: 'terratest-role-binding'
  namespace: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: 'terratest-role'
subjects:
- kind: ServiceAccount
  name: terratest
`
