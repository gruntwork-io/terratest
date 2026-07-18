package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListPersistentVolumesContextE will look for PersistentVolumes that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListPersistentVolumesContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.PersistentVolume, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().PersistentVolumes().List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListPersistentVolumesContext will look for PersistentVolumes that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListPersistentVolumesContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []corev1.PersistentVolume {
	t.Helper()
	pvs, err := ListPersistentVolumesContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return pvs
}

// GetPersistentVolumeContextE returns a Kubernetes PersistentVolume resource with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetPersistentVolumeContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, name string) (*corev1.PersistentVolume, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().PersistentVolumes().Get(ctx, name, metav1.GetOptions{})
}

// GetPersistentVolumeContext returns a Kubernetes PersistentVolume resource with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetPersistentVolumeContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, name string) *corev1.PersistentVolume {
	t.Helper()
	pv, err := GetPersistentVolumeContextE(t, ctx, options, name)
	require.NoError(t, err)

	return pv
}

// WaitUntilPersistentVolumeInStatusContextE waits until the given PersistentVolume is in the given status phase,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:dupl // structural pattern for k8s resource operations
func WaitUntilPersistentVolumeInStatusContextE(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	pvName string,
	pvStatusPhase *corev1.PersistentVolumePhase,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for Persistent Volume %s to be '%s'", pvName, *pvStatusPhase)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pv, err := GetPersistentVolumeContextE(t, ctx, options, pvName)
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
		options.Logger.Logf(t, "Timeout waiting for PersistentVolume to be '%s': %s", *pvStatusPhase, err)
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilPersistentVolumeInStatusContext waits until the given PersistentVolume is in the given status phase,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilPersistentVolumeInStatusContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, pvName string, pvStatusPhase *corev1.PersistentVolumePhase, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	require.NoError(t, WaitUntilPersistentVolumeInStatusContextE(t, ctx, options, pvName, pvStatusPhase, retries, sleepBetweenRetries))
}

// IsPersistentVolumeInStatus returns true if the given PersistentVolume is in the given status phase
func IsPersistentVolumeInStatus(pv *corev1.PersistentVolume, pvStatusPhase *corev1.PersistentVolumePhase) bool {
	return pv != nil && pv.Status.Phase == *pvStatusPhase
}
