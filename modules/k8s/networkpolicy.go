package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNetworkPolicy returns a Kubernetes networkpolicy resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
func GetNetworkPolicy(t testing.TestingT, options *KubectlOptions, networkPolicyName string) *networkingv1.NetworkPolicy {
	networkPolicy, err := GetNetworkPolicyE(t, options, networkPolicyName)
	require.NoError(t, err)
	return networkPolicy
}

// GetNetworkPolicyE returns a Kubernetes networkpolicy resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
func GetNetworkPolicyE(t testing.TestingT, options *KubectlOptions, networkPolicyName string) (*networkingv1.NetworkPolicy, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.NetworkingV1().NetworkPolicies(options.Namespace).Get(context.Background(), networkPolicyName, metav1.GetOptions{})
}

// WaitUntilNetworkPolicyAvailable waits until the networkpolicy is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
func WaitUntilNetworkPolicyAvailable(t testing.TestingT, options *KubectlOptions, networkPolicyName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for networkpolicy %s to be provisioned.", networkPolicyName)
	message := retry.DoWithRetry(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := GetNetworkPolicyE(t, options, networkPolicyName)
			if err != nil {
				return "", err
			}

			return "networkpolicy is now available", nil
		},
	)
	options.Logger.Logf(t, message)
}
