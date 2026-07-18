package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListJobsContextE looks up Jobs in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListJobsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]batchv1.Job, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.BatchV1().Jobs(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListJobsContext looks up Jobs in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListJobsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []batchv1.Job {
	t.Helper()
	jobs, err := ListJobsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return jobs
}

// GetJobContextE returns a Kubernetes job resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetJobContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, jobName string) (*batchv1.Job, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.BatchV1().Jobs(options.Namespace).Get(ctx, jobName, metav1.GetOptions{})
}

// GetJobContext returns a Kubernetes job resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetJobContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, jobName string) *batchv1.Job {
	t.Helper()
	job, err := GetJobContextE(t, ctx, options, jobName)
	require.NoError(t, err)

	return job
}

// WaitUntilJobSucceedContextE waits until requested job is succeeded, retrying the check for the specified amount of
// times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilJobSucceedContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, jobName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for job %s to be provisioned.", jobName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			job, err := GetJobContextE(t, ctx, options, jobName)
			if err != nil {
				return "", err
			}

			if !IsJobSucceeded(job) {
				return "", NewJobNotSucceeded(job)
			}

			return "Job is now Succeeded", nil
		},
	)
	if err != nil {
		options.Logger.Logf(t, "Timed out waiting for Job to be provisioned: %s", err)
		return err
	}

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilJobSucceedContext waits until requested job is succeeded, retrying the check for the specified amount of
// times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilJobSucceedContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, jobName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	require.NoError(t, WaitUntilJobSucceedContextE(t, ctx, options, jobName, retries, sleepBetweenRetries))
}

// IsJobSucceeded returns true when the job status condition "Complete" is true. This behavior is documented in the kubernetes API reference:
// https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/job-v1/#JobStatus
func IsJobSucceeded(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}

// CreateJobFromCronJobContextE creates a Job from the specified CronJob in the given namespace and returns the created Job.
// The ctx parameter supports cancellation and timeouts.
// This function is similar to running `kubectl create job --from=cronjob/<cron-job-name> <new-job-name>`.
func CreateJobFromCronJobContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, cronJobName, newJobName string) (*batchv1.Job, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	cronJob, err := GetCronJobContextE(t, ctx, options, cronJobName)
	if err != nil {
		return nil, err
	}

	annotations := make(map[string]string)
	for k, v := range cronJob.Spec.JobTemplate.Annotations {
		annotations[k] = v
	}

	job := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        newJobName,
			Namespace:   options.Namespace,
			Labels:      cronJob.Spec.JobTemplate.Labels,
			Annotations: annotations,
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	createdJob, err := clientset.BatchV1().Jobs(options.Namespace).Create(ctx, job, metav1.CreateOptions{})

	return createdJob, err
}

// CreateJobFromCronJobContext creates a Job from the specified CronJob in the given namespace and returns the created Job.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func CreateJobFromCronJobContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, cronJobName, newJobName string) *batchv1.Job {
	t.Helper()
	job, err := CreateJobFromCronJobContextE(t, ctx, options, cronJobName, newJobName)
	require.NoError(t, err)

	return job
}
