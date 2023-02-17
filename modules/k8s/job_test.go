//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestListJobsReturnsJobsInNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_JOB_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	jobs := ListJobs(t, options, metav1.ListOptions{})
	require.Equal(t, len(jobs), 1)
	job := jobs[0]
	require.Equal(t, job.Name, "pi-job")
	require.Equal(t, job.Namespace, uniqueID)
}

func TestGetJobEReturnsErrorForNonExistantJob(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetJobE(t, options, "pi-job")
	require.Error(t, err)
}

func TestGetJobEReturnsCorrectJobInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_JOB_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	job := GetJob(t, options, "pi-job")
	require.Equal(t, job.Name, "pi-job")
	require.Equal(t, job.Namespace, uniqueID)
}

func TestWaitUntilJobSucceedReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_JOB_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	WaitUntilJobSucceed(t, options, "pi-job", 60, 1*time.Second)
}

func TestIsJobSucceeded(t *testing.T) {
	t.Parallel()

	cases := []struct {
		title          string
		job            *batchv1.Job
		expectedResult bool
	}{
		{
			title: "TestIsJobSucceeded",
			job: &batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{
						batchv1.JobCondition{
							Type:   batchv1.JobComplete,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			expectedResult: true,
		},
		{
			title: "TestIsJobFailed",
			job: &batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{
						batchv1.JobCondition{
							Type:   batchv1.JobFailed,
							Status: corev1.ConditionTrue,
						},
					},
				},
			},
			expectedResult: false,
		},
		{
			title: "TestIsJobStarting",
			job: &batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{},
				},
			},
			expectedResult: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			actualResult := IsJobSucceeded(tc.job)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

const EXAMPLE_JOB_YAML_TEMPLATE = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: batch/v1
kind: Job
metadata:
  name: pi-job
  namespace: %s
spec:
  template:
    spec:
      containers:
      - name: pi
        image: "perl:5.34.1"
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
`
