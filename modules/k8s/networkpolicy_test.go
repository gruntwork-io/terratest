//go:build kubeall || kubernetes
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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetNetworkPolicyEReturnsErrorForNonExistantNetworkPolicy(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetNetworkPolicyE(t, options, "test-network-policy")
	require.Error(t, err)
}

func TestGetNetworkPolicyEReturnsCorrectNetworkPolicyInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_NETWORK_POLICY_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	networkPolicy := GetNetworkPolicy(t, options, "test-network-policy")
	require.Equal(t, networkPolicy.Name, "test-network-policy")
	require.Equal(t, networkPolicy.Namespace, uniqueID)
}

func TestWaitUntilNetworkPolicyAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_NETWORK_POLICY_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)

	KubectlApplyFromString(t, options, configData)
	WaitUntilNetworkPolicyAvailable(t, options, "test-network-policy", 10, 1*time.Second)
}

const EXAMPLE_NETWORK_POLICY_YAML_TEMPLATE = `---
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
