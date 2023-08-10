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

// ListPersistentVolumes will look for PersistentVolumes in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
func ListPersistentVolumes(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.PersistentVolume {
	pvs, err := ListPersistentVolumesE(t, options, filters)
	require.NoError(t, err)
	return pvs
}

// ListPersistentVolumesE will look for PersistentVolumes that match the given filters and return them.
func ListPersistentVolumesE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.PersistentVolume, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().PersistentVolumes().List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetPersistentVolume returns a Kubernetes PersistentVolume resource with the given name. This will fail the test if there is an error.
func GetPersistentVolume(t testing.TestingT, options *KubectlOptions, name string) *corev1.PersistentVolume {
	pv, err := GetPersistentVolumeE(t, options, name)
	require.NoError(t, err)
	return pv
}

// GetPersistentVolumeE returns a Kubernetes PersistentVolume resource with the given name.
func GetPersistentVolumeE(t testing.TestingT, options *KubectlOptions, name string) (*corev1.PersistentVolume, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().PersistentVolumes().Get(context.Background(), name, metav1.GetOptions{})
}

// WaitUntilPersistentVolumeInStatus waits until the given Persistent Volume is the given status phase,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
// This will fail the test if there is an error.
func WaitUntilPersistentVolumeInStatus(t testing.TestingT, options *KubectlOptions, pvName string, pvStatusPhase *corev1.PersistentVolumePhase, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilPersistentVolumeInStatusE(t, options, pvName, pvStatusPhase, retries, sleepBetweenRetries))
}

// WaitUntilPersistentVolumeInStatusE waits until the given PersistentVolume is in the given status phase,
// retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
func WaitUntilPersistentVolumeInStatusE(
	t testing.TestingT,
	options *KubectlOptions,
	pvName string,
	pvStatusPhase *corev1.PersistentVolumePhase,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for Persistent Volume %s to be '%s'", pvName, *pvStatusPhase)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pv, err := GetPersistentVolumeE(t, options, pvName)
			if err != nil {
				return "", err
			}
			if !IsPersistentVolumeInStatus(pv, pvStatusPhase) {
				return "", NewPersistentVolumeNotInStatusError(pv, pvStatusPhase)
			}
			return fmt.Sprintf("Persistent Volume is now '%s'", *pvStatusPhase), nil
		},
	)
	if err != nil {
		logger.Logf(t, "Timeout waiting for PersistentVolume to be '%s': %s", *pvStatusPhase, err)
		return err
	}
	logger.Logf(t, message)
	return nil
}

// IsPersistentVolumeInStatus returns true if the given PersistentVolume is in the given status phase
func IsPersistentVolumeInStatus(pv *corev1.PersistentVolume, pvStatusPhase *corev1.PersistentVolumePhase) bool {
	return pv != nil && pv.Status.Phase == *pvStatusPhase
}
