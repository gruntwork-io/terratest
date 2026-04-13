package k8s //nolint:dupl // structural pattern for k8s resource operations

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetConfigMapContextE returns a Kubernetes configmap resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
func GetConfigMapContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configMapName string) (*corev1.ConfigMap, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().ConfigMaps(options.Namespace).Get(ctx, configMapName, metav1.GetOptions{})
}

// GetConfigMapContext returns a Kubernetes configmap resource in the provided namespace with the given name. The
// namespace used is the one provided in the KubectlOptions.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetConfigMapContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configMapName string) *corev1.ConfigMap {
	t.Helper()
	configMap, err := GetConfigMapContextE(t, ctx, options, configMapName)
	require.NoError(t, err)

	return configMap
}

// GetConfigMap returns a Kubernetes configmap resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
//
// Deprecated: Use [GetConfigMapContext] instead.
func GetConfigMap(t testing.TestingT, options *KubectlOptions, configMapName string) *corev1.ConfigMap {
	t.Helper()

	return GetConfigMapContext(t, context.Background(), options, configMapName)
}

// GetConfigMapE returns a Kubernetes configmap resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
//
// Deprecated: Use [GetConfigMapContextE] instead.
func GetConfigMapE(t testing.TestingT, options *KubectlOptions, configMapName string) (*corev1.ConfigMap, error) {
	return GetConfigMapContextE(t, context.Background(), options, configMapName)
}

// WaitUntilConfigMapAvailableContextE waits until the configmap is present on the cluster in cases where it is not
// immediately available (for example, when using ClusterIssuer to request a certificate).
// The ctx parameter supports cancellation and timeouts.
func WaitUntilConfigMapAvailableContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, configMapName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for configmap %s to be provisioned.", configMapName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := GetConfigMapContextE(t, ctx, options, configMapName)
			if err != nil {
				return "", err
			}

			return "configmap is now available", nil
		},
	)
	if err != nil {
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilConfigMapAvailableContext waits until the configmap is present on the cluster in cases where it is not
// immediately available (for example, when using ClusterIssuer to request a certificate).
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilConfigMapAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, configMapName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	err := WaitUntilConfigMapAvailableContextE(t, ctx, options, configMapName, retries, sleepBetweenRetries)
	require.NoError(t, err)
}

// WaitUntilConfigMapAvailable waits until the configmap is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
//
// Deprecated: Use [WaitUntilConfigMapAvailableContext] instead.
func WaitUntilConfigMapAvailable(t testing.TestingT, options *KubectlOptions, configMapName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	WaitUntilConfigMapAvailableContext(t, context.Background(), options, configMapName, retries, sleepBetweenRetries)
}
