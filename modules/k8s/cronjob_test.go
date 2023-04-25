//go:build kubeall || kubernetes
// +build kubeall kubernetes

package k8s

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"k8s.io/api/batch/v1beta1"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestListCronJobsReturnsCronJobsInNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleCronjobYamlTemplate, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	jobs := ListCronJobs(t, options, metav1.ListOptions{})
	require.Equal(t, len(jobs), 1)
	job := jobs[0]
	require.Equal(t, job.Name, "cron-job")
	require.Equal(t, job.Namespace, uniqueID)
}

func TestGetCronJobEReturnErrorForNotExistingCronJob(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetJobE(t, options, random.UniqueId())
	require.Error(t, err)
}

func TestGetCronJobEReturnsCorrectJobInNamespace(t *testing.T) {
	t.Parallel()
	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleCronjobYamlTemplate, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	job := GetCronJob(t, options, "cron-job")
	require.Equal(t, job.Name, "cron-job")
	require.Equal(t, job.Namespace, uniqueID)
}

func TestWaitUntilCronJobScheduleSuccessfullyContainer(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleCronjobYamlTemplate, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	WaitUntilCronJobSucceed(t, options, "cron-job", 60, 5*time.Second)
}

func TestIsCronJobSucceeded(t *testing.T) {

	cases := []struct {
		title          string
		cronJob        *v1beta1.CronJob
		expectedResult bool
	}{
		{
			title: "CronJobScheduledContainer",
			cronJob: &v1beta1.CronJob{
				Status: v1beta1.CronJobStatus{
					LastScheduleTime: &metav1.Time{},
				},
			},
			expectedResult: true,
		},
		{
			title: "CronJobNotScheduledContainer",
			cronJob: &v1beta1.CronJob{
				Status: v1beta1.CronJobStatus{
					LastScheduleTime: nil,
				},
			},
			expectedResult: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			actualResult := IsCronJobSucceeded(tc.cronJob)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

const ExampleCronjobYamlTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: cron-job
  namespace: %s
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: ubuntu
            image: ubuntu:20.04
            command: ["sh", "-c", "ls"]
          restartPolicy: OnFailure
`
