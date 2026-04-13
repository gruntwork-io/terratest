package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListCronJobsContextE lists cron jobs in namespace that match provided filters and returns them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListCronJobsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]batchv1.CronJob, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.BatchV1().CronJobs(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListCronJobsContext lists cron jobs in namespace that match provided filters and returns them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListCronJobsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []batchv1.CronJob {
	t.Helper()
	cronJobs, err := ListCronJobsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return cronJobs
}

// ListCronJobs list cron jobs in namespace that match provided filters. This will fail the test if there is an error.
//
// Deprecated: Use [ListCronJobsContext] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListCronJobs(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []batchv1.CronJob {
	t.Helper()

	return ListCronJobsContext(t, context.Background(), options, filters)
}

// ListCronJobsE list cron jobs in namespace that match provided filters. This will return list or error.
//
// Deprecated: Use [ListCronJobsContextE] instead.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListCronJobsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]batchv1.CronJob, error) {
	return ListCronJobsContextE(t, context.Background(), options, filters)
}

// GetCronJobContextE returns a cron job resource from namespace by name.
// The ctx parameter supports cancellation and timeouts.
func GetCronJobContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, cronJobName string) (*batchv1.CronJob, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.BatchV1().CronJobs(options.Namespace).Get(ctx, cronJobName, metav1.GetOptions{})
}

// GetCronJobContext returns a cron job resource from namespace by name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetCronJobContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, cronJobName string) *batchv1.CronJob {
	t.Helper()
	job, err := GetCronJobContextE(t, ctx, options, cronJobName)
	require.NoError(t, err)

	return job
}

// GetCronJob return cron job resource from namespace by name. This will fail the test if there is an error.
//
// Deprecated: Use [GetCronJobContext] instead.
func GetCronJob(t testing.TestingT, options *KubectlOptions, cronJobName string) *batchv1.CronJob {
	t.Helper()

	return GetCronJobContext(t, context.Background(), options, cronJobName)
}

// GetCronJobE return cron job resource from namespace by name. This will return cron job or error.
//
// Deprecated: Use [GetCronJobContextE] instead.
func GetCronJobE(t testing.TestingT, options *KubectlOptions, cronJobName string) (*batchv1.CronJob, error) {
	return GetCronJobContextE(t, context.Background(), options, cronJobName)
}

// WaitUntilCronJobSucceedContextE waits until cron job will successfully complete a job, retrying the check for the
// specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilCronJobSucceedContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, cronJobName string, retries int, sleepBetweenRetries time.Duration) error { //nolint:dupl // similar retry pattern across resource types is intentional
	statusMsg := fmt.Sprintf("Wait for CronJob %s to successfully schedule container", cronJobName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			job, err := GetCronJobContextE(t, ctx, options, cronJobName)
			if err != nil {
				return "", err
			}

			if !IsCronJobSucceeded(job) {
				return "", NewCronJobNotSucceeded(job)
			}

			return "CronJob scheduled container", nil
		},
	)
	if err != nil {
		options.Logger.Logf(t, "Timed out waiting for CronJob to schedule job: %s", err)
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilCronJobSucceedContext waits until cron job will successfully complete a job, retrying the check for the
// specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilCronJobSucceedContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, cronJobName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	require.NoError(t, WaitUntilCronJobSucceedContextE(t, ctx, options, cronJobName, retries, sleepBetweenRetries))
}

// WaitUntilCronJobSucceed waits until cron job will successfully complete a job. This will fail the test if there is an
// error or if the check times out.
//
// Deprecated: Use [WaitUntilCronJobSucceedContext] instead.
func WaitUntilCronJobSucceed(t testing.TestingT, options *KubectlOptions, cronJobName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	WaitUntilCronJobSucceedContext(t, context.Background(), options, cronJobName, retries, sleepBetweenRetries)
}

// WaitUntilCronJobSucceedE waits until cron job will successfully complete a job, retrying the check for the specified
// amount of times, sleeping for the provided duration between each try.
//
// Deprecated: Use [WaitUntilCronJobSucceedContextE] instead.
func WaitUntilCronJobSucceedE(t testing.TestingT, options *KubectlOptions, cronJobName string, retries int, sleepBetweenRetries time.Duration) error {
	return WaitUntilCronJobSucceedContextE(t, context.Background(), options, cronJobName, retries, sleepBetweenRetries)
}

// IsCronJobSucceeded returns true if cron job successfully scheduled and completed job.
// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/cron-job-v1/#CronJobStatus
func IsCronJobSucceeded(cronJob *batchv1.CronJob) bool {
	return cronJob.Status.LastScheduleTime != nil
}
