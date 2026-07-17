package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNetworkPolicyContextE returns a Kubernetes networkpolicy resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkPolicyContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, networkPolicyName string) (*networkingv1.NetworkPolicy, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.NetworkingV1().NetworkPolicies(options.Namespace).Get(ctx, networkPolicyName, metav1.GetOptions{})
}

// GetNetworkPolicyContext returns a Kubernetes networkpolicy resource in the provided namespace with the given name.
// The namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetNetworkPolicyContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, networkPolicyName string) *networkingv1.NetworkPolicy {
	t.Helper()
	networkPolicy, err := GetNetworkPolicyContextE(t, ctx, options, networkPolicyName)
	require.NoError(t, err)

	return networkPolicy
}

// WaitUntilNetworkPolicyAvailableContextE waits until the networkpolicy is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
// The ctx parameter supports cancellation and timeouts.
func WaitUntilNetworkPolicyAvailableContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, networkPolicyName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for networkpolicy %s to be provisioned.", networkPolicyName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := GetNetworkPolicyContextE(t, ctx, options, networkPolicyName)
			if err != nil {
				return "", err
			}

			return "networkpolicy is now available", nil
		},
	)
	if err != nil {
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilNetworkPolicyAvailableContext waits until the networkpolicy is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilNetworkPolicyAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, networkPolicyName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilNetworkPolicyAvailableContextE(t, ctx, options, networkPolicyName, retries, sleepBetweenRetries)
	require.NoError(t, err)
}
