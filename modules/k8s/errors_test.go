package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestErrorDeploymentNotAvailable(t *testing.T) {
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
			err := NewDeploymentNotAvailableError(tc.deploy)
			assert.EqualError(t, err, tc.expectedErr)
		})
	}
}
