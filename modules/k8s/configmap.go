package k8s

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

// GetConfigMap returns a Kubernetes configmap resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions. This will fail the test if there is an error.
func GetConfigMap(t testing.TestingT, options *KubectlOptions, configMapName string) *corev1.ConfigMap {
	configMap, err := GetConfigMapE(t, options, configMapName)
	require.NoError(t, err)
	return configMap
}

// GetConfigMapE returns a Kubernetes configmap resource in the provided namespace with the given name. The namespace used
// is the one provided in the KubectlOptions.
func GetConfigMapE(t testing.TestingT, options *KubectlOptions, configMapName string) (*corev1.ConfigMap, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().ConfigMaps(options.Namespace).Get(context.Background(), configMapName, metav1.GetOptions{})
}

// WaitUntilConfigMapAvailable waits until the configmap is present on the cluster in cases where it is not immediately
// available (for example, when using ClusterIssuer to request a certificate).
func WaitUntilConfigMapAvailable(t testing.TestingT, options *KubectlOptions, configMapName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for configmap %s to be provisioned.", configMapName)
	message := retry.DoWithRetry(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			_, err := GetConfigMapE(t, options, configMapName)
			if err != nil {
				return "", err
			}

			return "configmap is now available", nil
		},
	)
	options.Logger.Logf(t, message)
}
