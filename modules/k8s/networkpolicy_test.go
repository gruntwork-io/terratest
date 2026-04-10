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
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetNetworkPolicyEReturnsErrorForNonExistantNetworkPolicy(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetNetworkPolicyE(t, options, "test-network-policy")
	require.Error(t, err)
}

func TestGetNetworkPolicyEReturnsCorrectNetworkPolicyInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleNetworkPolicyYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	networkPolicy := k8s.GetNetworkPolicy(t, options, "test-network-policy")
	require.Equal(t, "test-network-policy", networkPolicy.Name)
	require.Equal(t, networkPolicy.Namespace, uniqueID)
}

func TestWaitUntilNetworkPolicyAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleNetworkPolicyYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)
	k8s.WaitUntilNetworkPolicyAvailable(t, options, "test-network-policy", 10, 1*time.Second)
}

const exampleNetworkPolicyYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: test-network-policy
  namespace: %s
spec:
  podSelector: {}
  policyTypes:
    - Ingress
`
