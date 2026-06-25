package k8s_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/k8s"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestErrorDeploymentNotAvailable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		title       string
		deploy      *appsv1.Deployment
		expectedErr string
	}{
		{
			title: "NoProgressingCondition",
			deploy: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{},
				},
			},
			expectedErr: "Deployment foo is not available, missing 'Progressing' condition",
		},
		{
			title: "DeploymentNotComplete",
			deploy: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:    appsv1.DeploymentProgressing,
							Status:  v1.ConditionTrue,
							Reason:  "ReplicaSetUpdated",
							Message: "bar",
						},
					},
				},
			},
			expectedErr: "Deployment foo is not available as 'Progressing' condition indicates that the Deployment is not complete, status: True, reason: ReplicaSetUpdated, message: bar",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			err := k8s.NewDeploymentNotAvailableError(tc.deploy)
			assert.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestErrorPodNotAvailable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		pod         *v1.Pod
		title       string
		expectedErr string
	}{
		{
			title: "PodLevelReasonOnly",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Status:     v1.PodStatus{Reason: "Evicted", Message: "low memory"},
			},
			expectedErr: "Pod foo is not available, reason: Evicted, message: low memory",
		},
		{
			title: "ContainerWaiting",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Status: v1.PodStatus{
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name:  "web",
							Ready: false,
							State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{
								Reason:  "CrashLoopBackOff",
								Message: "back-off 5m restarting failed container",
							}},
						},
					},
				},
			},
			expectedErr: "Pod foo is not available, reason: , message: . container web waiting: CrashLoopBackOff back-off 5m restarting failed container",
		},
		{
			title: "InitContainerWaitingAndReadyContainerIgnored",
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Status: v1.PodStatus{
					InitContainerStatuses: []v1.ContainerStatus{
						{Name: "setup", Ready: false, State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: "ImagePullBackOff"}}},
					},
					ContainerStatuses: []v1.ContainerStatus{
						{Name: "web", Ready: true},
					},
				},
			},
			expectedErr: "Pod foo is not available, reason: , message: . init container setup waiting: ImagePullBackOff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			err := k8s.NewPodNotAvailableError(tc.pod)
			assert.EqualError(t, err, tc.expectedErr)
		})
	}
}
