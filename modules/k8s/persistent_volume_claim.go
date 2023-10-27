package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListPersistentVolumeClaims will look for PersistentVolumeClaims in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
func ListPersistentVolumeClaims(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.PersistentVolumeClaim {
	pvcs, err := ListPersistentVolumeClaimsE(t, options, filters)
	require.NoError(t, err)
	return pvcs
}

// ListPersistentVolumeClaimsE will look for PersistentVolumeClaims in the given namespace that match the given filters and return them.
func ListPersistentVolumeClaimsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.PersistentVolumeClaim, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().PersistentVolumeClaims(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetPersistentVolumeClaim returns a Kubernetes PersistentVolumeClaim resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetPersistentVolumeClaim(t testing.TestingT, options *KubectlOptions, pvcName string) *corev1.PersistentVolumeClaim {
	pvc, err := GetPersistentVolumeClaimE(t, options, pvcName)
	require.NoError(t, err)
	return pvc
}

// GetPersistentVolumeClaimE returns a Kubernetes PersistentVolumeClaim resource in the provided namespace with the given name.
func GetPersistentVolumeClaimE(t testing.TestingT, options *KubectlOptions, pvcName string) (*corev1.PersistentVolumeClaim, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().PersistentVolumeClaims(options.Namespace).Get(context.Background(), pvcName, metav1.GetOptions{})
}

// WaitUntilPersistentVolumeClaimInStatus waits until the given PersistentVolumeClaim is the given status phase,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
// This will fail the test if there is an error.
func WaitUntilPersistentVolumeClaimInStatus(t testing.TestingT, options *KubectlOptions, pvcName string, pvcStatusPhase *corev1.PersistentVolumeClaimPhase, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilPersistentVolumeClaimInStatusE(t, options, pvcName, pvcStatusPhase, retries, sleepBetweenRetries))
}

// WaitUntilPersistentVolumeClaimInStatusE waits until the given PersistentVolumeClaim is the given status phase,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
// This will fail the test if there is an error.
func WaitUntilPersistentVolumeClaimInStatusE(t testing.TestingT, options *KubectlOptions, pvcName string, pvcStatusPhase *corev1.PersistentVolumeClaimPhase, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for PersistentVolumeClaim %s to be '%s'.", pvcName, *pvcStatusPhase)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pvc, err := GetPersistentVolumeClaimE(t, options, pvcName)
			if err != nil {
				return "", err
			}
			if !IsPersistentVolumeClaimInStatus(pvc, pvcStatusPhase) {
				return "", NewPersistentVolumeClaimNotInStatusError(pvc, pvcStatusPhase)
			}
			return fmt.Sprintf("PersistentVolumeClaim is now '%s'", *pvcStatusPhase), nil
		},
	)
	if err != nil {
		logger.Default.Logf(t, "Timeout waiting for PersistentVolumeClaim to be '%s': %s", *pvcStatusPhase, err)
		return err
	}
	logger.Default.Logf(t, message)
	return nil
}

// IsPersistentVolumeClaimInStatus returns true if the given PersistentVolumeClaim is in the given status phase
func IsPersistentVolumeClaimInStatus(pvc *corev1.PersistentVolumeClaim, pvcStatusPhase *corev1.PersistentVolumeClaimPhase) bool {
	return pvc != nil && pvc.Status.Phase == *pvcStatusPhase
}
