package k8s

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

// GetSecret returns a Kubernetes secret resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
func GetSecret(t testing.TestingT, options *KubectlOptions, secretName string) *corev1.Secret {
	secret, err := GetSecretE(t, options, secretName)
	require.NoError(t, err)
	return secret
}

// GetSecretE returns a Kubernetes secret resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
func GetSecretE(t testing.TestingT, options *KubectlOptions, secretName string) (*corev1.Secret, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Secrets(options.Namespace).Get(context.Background(), secretName, metav1.GetOptions{})
}

// WaitUntilSecretAvailable waits until the secret is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
func WaitUntilSecretAvailable(t testing.TestingT, options *KubectlOptions, secretName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for secret %s to be provisioned.", secretName)
	message := retry.DoWithRetry(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := GetSecretE(t, options, secretName)
			if err != nil {
				return "", err
			}

			return "Secret is now available", nil
		},
	)
	options.Logger.Logf(t, message)
}
