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

func TestGetRoleEReturnsErrorForNonExistantRole(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetRoleE(t, options, "non-existing-role")
	require.Error(t, err)
}

func TestGetRoleEReturnsCorrectRoleInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_ROLE_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	role := GetRole(t, options, "terratest-role")
	require.Equal(t, role.Name, "terratest-role")
	require.Equal(t, role.Namespace, uniqueID)
}

const EXAMPLE_ROLE_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: '%s'
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
`
