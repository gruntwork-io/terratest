package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListPersistentVolumeClaimsContextE will look for PersistentVolumeClaims in the given namespace that match the given
// filters and return them. The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListPersistentVolumeClaimsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.PersistentVolumeClaim, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().PersistentVolumeClaims(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListPersistentVolumeClaimsContext will look for PersistentVolumeClaims in the given namespace that match the given
// filters and return them. The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListPersistentVolumeClaimsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []corev1.PersistentVolumeClaim {
	t.Helper()
	pvcs, err := ListPersistentVolumeClaimsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return pvcs
}

// GetPersistentVolumeClaimContextE returns a Kubernetes PersistentVolumeClaim resource in the provided namespace with
// the given name. The ctx parameter supports cancellation and timeouts.
func GetPersistentVolumeClaimContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, pvcName string) (*corev1.PersistentVolumeClaim, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().PersistentVolumeClaims(options.Namespace).Get(ctx, pvcName, metav1.GetOptions{})
}

// GetPersistentVolumeClaimContext returns a Kubernetes PersistentVolumeClaim resource in the provided namespace with
// the given name. The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetPersistentVolumeClaimContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, pvcName string) *corev1.PersistentVolumeClaim {
	t.Helper()
	pvc, err := GetPersistentVolumeClaimContextE(t, ctx, options, pvcName)
	require.NoError(t, err)

	return pvc
}

// WaitUntilPersistentVolumeClaimInStatusContextE waits until the given PersistentVolumeClaim is the given status phase,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:dupl // structural pattern for k8s resource operations
func WaitUntilPersistentVolumeClaimInStatusContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, pvcName string, pvcStatusPhase *corev1.PersistentVolumeClaimPhase, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for PersistentVolumeClaim %s to be '%s'.", pvcName, *pvcStatusPhase)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pvc, err := GetPersistentVolumeClaimContextE(t, ctx, options, pvcName)
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

	logger.Default.Logf(t, "%s", message)

	return nil
}

// WaitUntilPersistentVolumeClaimInStatusContext waits until the given PersistentVolumeClaim is the given status phase,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilPersistentVolumeClaimInStatusContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, pvcName string, pvcStatusPhase *corev1.PersistentVolumeClaimPhase, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	require.NoError(t, WaitUntilPersistentVolumeClaimInStatusContextE(t, ctx, options, pvcName, pvcStatusPhase, retries, sleepBetweenRetries))
}

// IsPersistentVolumeClaimInStatus returns true if the given PersistentVolumeClaim is in the given status phase
func IsPersistentVolumeClaimInStatus(pvc *corev1.PersistentVolumeClaim, pvcStatusPhase *corev1.PersistentVolumeClaimPhase) bool {
	return pvc != nil && pvc.Status.Phase == *pvcStatusPhase
}
