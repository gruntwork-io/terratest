package k8s //nolint:dupl // structural pattern for k8s resource operations

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetSecretContextE returns a Kubernetes secret resource in the provided namespace with the given name. The namespace
// used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
func GetSecretContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, secretName string) (*corev1.Secret, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Secrets(options.Namespace).Get(ctx, secretName, metav1.GetOptions{})
}

// GetSecretContext returns a Kubernetes secret resource in the provided namespace with the given name. The namespace
// used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetSecretContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, secretName string) *corev1.Secret {
	t.Helper()
	secret, err := GetSecretContextE(t, ctx, options, secretName)
	require.NoError(t, err)

	return secret
}

// GetSecret returns a Kubernetes secret resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
//
// Deprecated: Use [GetSecretContext] instead.
func GetSecret(t testing.TestingT, options *KubectlOptions, secretName string) *corev1.Secret {
	t.Helper()

	return GetSecretContext(t, context.Background(), options, secretName)
}

// GetSecretE returns a Kubernetes secret resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
//
// Deprecated: Use [GetSecretContextE] instead.
func GetSecretE(t testing.TestingT, options *KubectlOptions, secretName string) (*corev1.Secret, error) {
	return GetSecretContextE(t, context.Background(), options, secretName)
}

// WaitUntilSecretAvailableContextE waits until the secret is present on the cluster in cases where it is not
// immediately available (for example, when using ClusterIssuer to request a certificate).
// The ctx parameter supports cancellation and timeouts.
func WaitUntilSecretAvailableContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, secretName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for secret %s to be provisioned.", secretName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := GetSecretContextE(t, ctx, options, secretName)
			if err != nil {
				return "", err
			}

			return "Secret is now available", nil
		},
	)
	if err != nil {
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilSecretAvailableContext waits until the secret is present on the cluster in cases where it is not
// immediately available (for example, when using ClusterIssuer to request a certificate).
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilSecretAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, secretName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilSecretAvailableContextE(t, ctx, options, secretName, retries, sleepBetweenRetries)
	require.NoError(t, err)
}

// WaitUntilSecretAvailable waits until the secret is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
//
// Deprecated: Use [WaitUntilSecretAvailableContext] instead.
func WaitUntilSecretAvailable(t testing.TestingT, options *KubectlOptions, secretName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	WaitUntilSecretAvailableContext(t, context.Background(), options, secretName, retries, sleepBetweenRetries)
}
