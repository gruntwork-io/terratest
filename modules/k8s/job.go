package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListJobs will look for Jobs in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
func ListJobs(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []batchv1.Job {
	jobs, err := ListJobsE(t, options, filters)
	require.NoError(t, err)
	return jobs
}

// ListJobsE will look for jobs in the given namespace that match the given filters and return them.
func ListJobsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]batchv1.Job, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.BatchV1().Jobs(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetJob returns a Kubernetes job resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetJob(t testing.TestingT, options *KubectlOptions, jobName string) *batchv1.Job {
	job, err := GetJobE(t, options, jobName)
	require.NoError(t, err)
	return job
}

// GetJobE returns a Kubernetes job resource in the provided namespace with the given name.
func GetJobE(t testing.TestingT, options *KubectlOptions, jobName string) (*batchv1.Job, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.BatchV1().Jobs(options.Namespace).Get(context.Background(), jobName, metav1.GetOptions{})
}

// WaitUntilJobSucceed waits until requested job is suceeded, retrying the check for the specified amount of times, sleeping
// for the provided duration between each try. This will fail the test if there is an error or if the check times out.
func WaitUntilJobSucceed(t testing.TestingT, options *KubectlOptions, jobName string, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilJobSucceedE(t, options, jobName, retries, sleepBetweenRetries))
}

// WaitUntilJobSucceedE waits until requested job is succeeded, retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
func WaitUntilJobSucceedE(t testing.TestingT, options *KubectlOptions, jobName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for job %s to be provisioned.", jobName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			job, err := GetJobE(t, options, jobName)
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
	options.Logger.Logf(t, message)
	return nil
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
