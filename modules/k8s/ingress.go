package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListIngressesContextE will look for Ingress resources in the given namespace that match the given filters and return
// them. The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListIngressesContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]networkingv1.Ingress, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.NetworkingV1().Ingresses(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListIngressesContext will look for Ingress resources in the given namespace that match the given filters and return
// them. The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListIngressesContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []networkingv1.Ingress {
	t.Helper()
	ingresses, err := ListIngressesContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return ingresses
}

// GetIngressContextE returns a Kubernetes Ingress resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetIngressContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string) (*networkingv1.Ingress, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.NetworkingV1().Ingresses(options.Namespace).Get(ctx, ingressName, metav1.GetOptions{})
}

// GetIngressContext returns a Kubernetes Ingress resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetIngressContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string) *networkingv1.Ingress {
	t.Helper()
	ingress, err := GetIngressContextE(t, ctx, options, ingressName)
	require.NoError(t, err)

	return ingress
}

// IsIngressAvailable returns true if the Ingress endpoint is provisioned and available.
func IsIngressAvailable(ingress *networkingv1.Ingress) bool {

	endpoints := ingress.Status.LoadBalancer.Ingress
	return len(endpoints) > 0
}

// WaitUntilIngressAvailableContextE waits until the Ingress resource has an endpoint provisioned for it.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilIngressAvailableContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for ingress %s to be provisioned.", ingressName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			ingress, err := GetIngressContextE(t, ctx, options, ingressName)
			if err != nil {
				return "", err
			}

			if !IsIngressAvailable(ingress) {
				return "", IngressNotAvailable{ingress: ingress}
			}

			return "Ingress is now available", nil
		},
	)
	if err != nil {
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilIngressAvailableContext waits until the Ingress resource has an endpoint provisioned for it.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilIngressAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilIngressAvailableContextE(t, ctx, options, ingressName, retries, sleepBetweenRetries)
	require.NoError(t, err)
}

// ListIngressesV1Beta1ContextE will look for Ingress resources in the given namespace that match the given filters and
// return them, using networking.k8s.io/v1beta1 API.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListIngressesV1Beta1ContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]networkingv1beta1.Ingress, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.NetworkingV1beta1().Ingresses(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListIngressesV1Beta1Context will look for Ingress resources in the given namespace that match the given filters and
// return them, using networking.k8s.io/v1beta1 API.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListIngressesV1Beta1Context(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []networkingv1beta1.Ingress {
	t.Helper()
	ingresses, err := ListIngressesV1Beta1ContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return ingresses
}

// GetIngressV1Beta1ContextE returns a Kubernetes Ingress resource in the provided namespace with the given name, using
// networking.k8s.io/v1beta1 API.
// The ctx parameter supports cancellation and timeouts.
func GetIngressV1Beta1ContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string) (*networkingv1beta1.Ingress, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.NetworkingV1beta1().Ingresses(options.Namespace).Get(ctx, ingressName, metav1.GetOptions{})
}

// GetIngressV1Beta1Context returns a Kubernetes Ingress resource in the provided namespace with the given name, using
// networking.k8s.io/v1beta1 API.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetIngressV1Beta1Context(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string) *networkingv1beta1.Ingress {
	t.Helper()
	ingress, err := GetIngressV1Beta1ContextE(t, ctx, options, ingressName)
	require.NoError(t, err)

	return ingress
}

// IsIngressAvailableV1Beta1 returns true if the Ingress endpoint is provisioned and available, using
// networking.k8s.io/v1beta1 API.
func IsIngressAvailableV1Beta1(ingress *networkingv1beta1.Ingress) bool {

	endpoints := ingress.Status.LoadBalancer.Ingress
	return len(endpoints) > 0
}

// WaitUntilIngressAvailableV1Beta1ContextE waits until the Ingress resource has an endpoint provisioned for it, using
// networking.k8s.io/v1beta1 API.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilIngressAvailableV1Beta1ContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for ingress %s to be provisioned.", ingressName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			ingress, err := GetIngressV1Beta1ContextE(t, ctx, options, ingressName)
			if err != nil {
				return "", err
			}

			if !IsIngressAvailableV1Beta1(ingress) {
				return "", IngressNotAvailableV1Beta1{ingress: ingress}
			}

			return "Ingress is now available", nil
		},
	)
	if err != nil {
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilIngressAvailableV1Beta1Context waits until the Ingress resource has an endpoint provisioned for it, using
// networking.k8s.io/v1beta1 API.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilIngressAvailableV1Beta1Context(t testing.TestingT, ctx context.Context, options *KubectlOptions, ingressName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilIngressAvailableV1Beta1ContextE(t, ctx, options, ingressName, retries, sleepBetweenRetries)
	require.NoError(t, err)
}
