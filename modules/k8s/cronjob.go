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

// ListCronJobs list cron jobs in namespace that match provided filters. This will fail the test if there is an error.
func ListCronJobs(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []batchv1.CronJob {
	cronJobs, err := ListCronJobsE(t, options, filters)
	require.NoError(t, err)
	return cronJobs
}

// ListCronJobsE list cron jobs in namespace that match provided filters. This will return list or error.
func ListCronJobsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]batchv1.CronJob, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	resp, err := clientset.BatchV1().CronJobs(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetCronJob return cron job resource from namespace by name. This will fail the test if there is an error.
func GetCronJob(t testing.TestingT, options *KubectlOptions, cronJobName string) *batchv1.CronJob {
	job, err := GetCronJobE(t, options, cronJobName)
	require.NoError(t, err)
	return job
}

// GetCronJobE return cron job resource from namespace by name. This will return cron job or error.
func GetCronJobE(t testing.TestingT, options *KubectlOptions, cronJobName string) (*batchv1.CronJob, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.BatchV1().CronJobs(options.Namespace).Get(context.Background(), cronJobName, metav1.GetOptions{})
}

// WaitUntilCronJobSucceed waits until cron job will successfully complete a job. This will fail the test if there is an
// error or if the check times out.
func WaitUntilCronJobSucceed(t testing.TestingT, options *KubectlOptions, cronJobName string, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilCronJobSucceedE(t, options, cronJobName, retries, sleepBetweenRetries))
}

// WaitUntilCronJobSucceedE waits until cron job will successfully complete a job, retrying the check for the specified
// amount of times, sleeping for the provided duration between each try.
func WaitUntilCronJobSucceedE(t testing.TestingT, options *KubectlOptions, cronJobName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for CronJob %s to successfully schedule container", cronJobName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			job, err := GetCronJobE(t, options, cronJobName)
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
	options.Logger.Logf(t, message)
	return nil
}

// IsCronJobSucceeded returns true if cron job successfully scheduled and completed job.
// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/cron-job-v1/#CronJobStatus
func IsCronJobSucceeded(cronJob *batchv1.CronJob) bool {
	return cronJob.Status.LastScheduleTime != nil
}
