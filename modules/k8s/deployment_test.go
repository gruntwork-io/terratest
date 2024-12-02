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
	"time"

	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetDeploymentEReturnsError(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	_, err := GetDeploymentE(t, options, "nginx-deployment")
	require.Error(t, err)
}

func TestGetDeployments(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleDeploymentYAMLTemplate, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	deployment := GetDeployment(t, options, "nginx-deployment")
	require.Equal(t, deployment.Name, "nginx-deployment")
	require.Equal(t, deployment.Namespace, uniqueID)
}

func TestListDeployments(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleDeploymentYAMLTemplate, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	deployments := ListDeployments(t, options, metav1.ListOptions{})
	require.Equal(t, len(deployments), 1)

	deployment := deployments[0]
	require.Equal(t, deployment.Name, "nginx-deployment")
	require.Equal(t, deployment.Namespace, uniqueID)
}

func TestWaitUntilDeploymentAvailable(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(ExampleDeploymentYAMLTemplate, uniqueID)
	KubectlApplyFromString(t, options, configData)
	defer KubectlDeleteFromString(t, options, configData)

	WaitUntilDeploymentAvailable(t, options, "nginx-deployment", 60, 1*time.Second)
}

func TestTestIsDeploymentAvailable(t *testing.T) {
	testCases := []struct {
		title          string
		deploy         *appsv1.Deployment
		expectedResult bool
	}{
		{
			title: "TestIsDeploymentAvailableWithProgressingNewReplicaSetAvailable",
			deploy: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:   appsv1.DeploymentProgressing,
							Status: v1.ConditionTrue,
							Reason: "NewReplicaSetAvailable",
						},
					},
				},
			},
			expectedResult: true,
		},
		{
			title: "TestIsDeploymentAvailableWithoutProgressingNewReplicaSetAvailable",
			deploy: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{
						{
							Type:   appsv1.DeploymentProgressing,
							Status: v1.ConditionTrue,
							Reason: "ReplicaSetUpdated",
						},
					},
				},
			},
			expectedResult: false,
		},
		{
			title: "TestIsDeploymentAvailableWithoutProgressingCondition",
			deploy: &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Conditions: []appsv1.DeploymentCondition{},
				},
			},
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			actualResult := IsDeploymentAvailable(tc.deploy)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

const ExampleDeploymentYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  strategy:
    rollingUpdate:
      maxSurge: 10%%
      maxUnavailable: 0
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.15.7
        ports:
        - containerPort: 80
        readinessProbe:
          httpGet:
            path: /
            port: 80
`
