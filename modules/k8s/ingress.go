package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListIngresses will look for Ingress resources in the given namespace that match the given filters and return them.
// This will fail the test if there is an error.
func ListIngresses(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []networkingv1.Ingress {
	ingresses, err := ListIngressesE(t, options, filters)
	require.NoError(t, err)
	return ingresses
}

// ListIngressesE will look for Ingress resources in the given namespace that match the given filters and return them.
func ListIngressesE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]networkingv1.Ingress, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	resp, err := clientset.NetworkingV1().Ingresses(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil

}

// GetIngress returns a Kubernetes Ingress resource in the provided namespace with the given name. This will fail the
// test if there is an error.
func GetIngress(t testing.TestingT, options *KubectlOptions, ingressName string) *networkingv1.Ingress {
	ingress, err := GetIngressE(t, options, ingressName)
	require.NoError(t, err)
	return ingress
}

// GetIngressE returns a Kubernetes Ingress resource in the provided namespace with the given name.
func GetIngressE(t testing.TestingT, options *KubectlOptions, ingressName string) (*networkingv1.Ingress, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.NetworkingV1().Ingresses(options.Namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
}

// IsIngressAvailable returns true if the Ingress endpoint is provisioned and available.
func IsIngressAvailable(ingress *networkingv1.Ingress) bool {
	// Ingress is ready if it has at least one endpoint
	endpoints := ingress.Status.LoadBalancer.Ingress
	return len(endpoints) > 0
}

// WaitUntilIngressAvailable waits until the Ingress resource has an endpoint provisioned for it.
func WaitUntilIngressAvailable(t testing.TestingT, options *KubectlOptions, ingressName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for ingress %s to be provisioned.", ingressName)
	message := retry.DoWithRetry(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			ingress, err := GetIngressE(t, options, ingressName)
			if err != nil {
				return "", err
			}
			if !IsIngressAvailable(ingress) {
				return "", IngressNotAvailable{ingress: ingress}
			}
			return "Ingress is now available", nil
		},
	)
	logger.Logf(t, message)
}
