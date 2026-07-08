//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.
package k8s_test

import (
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/require"
)

func TestGetDaemonSetEReturnsErrorForNonExistantDaemonSet(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "")
	_, err := k8s.GetDaemonSetE(t, options, "sample-ds")
	require.Error(t, err)
}

func TestGetDaemonSetEReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(exampleDaemonSetYAMLTemplate, uniqueID, uniqueID)

	k8s.KubectlApplyFromString(t, options, configData)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	daemonSet := k8s.GetDaemonSet(t, options, "sample-ds")
	require.Equal(t, "sample-ds", daemonSet.Name)
	require.Equal(t, daemonSet.Namespace, uniqueID)
}

func TestListDaemonSetsReturnsCorrectServiceInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(exampleDaemonSetYAMLTemplate, uniqueID, uniqueID)

	k8s.KubectlApplyFromString(t, options, configData)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	daemonSets := k8s.ListDaemonSets(t, options, metav1.ListOptions{})
	require.Len(t, daemonSets, 1)

	daemonSet := daemonSets[0]
	require.Equal(t, "sample-ds", daemonSet.Name)
	require.Equal(t, daemonSet.Namespace, uniqueID)
}

func TestWaitUntilDaemonSetAvailable(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(exampleDaemonSetYAMLTemplate, uniqueID, uniqueID)

	k8s.KubectlApplyFromString(t, options, configData)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.WaitUntilDaemonSetAvailable(t, options, "sample-ds", 60, 1*time.Second)
}

func TestIsDaemonSetAvailable(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ds             *appsv1.DaemonSet
		title          string
		expectedResult bool
	}{
		{
			title: "AvailableWhenAllPodsUpdatedAndAvailable",
			ds: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 3,
					UpdatedNumberScheduled: 3,
					NumberAvailable:        3,
				},
			},
			expectedResult: true,
		},
		{
			title: "AvailableWhenNoNodesMatchSelector",
			ds: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 0,
					UpdatedNumberScheduled: 0,
					NumberAvailable:        0,
				},
			},
			expectedResult: true,
		},
		{
			title: "NotAvailableWhenSomePodsNotYetAvailable",
			ds: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 1},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 3,
					UpdatedNumberScheduled: 3,
					NumberAvailable:        2,
				},
			},
			expectedResult: false,
		},
		{
			title: "NotAvailableMidRollingUpdate",
			ds: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 2},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     2,
					DesiredNumberScheduled: 3,
					UpdatedNumberScheduled: 1,
					NumberAvailable:        3,
				},
			},
			expectedResult: false,
		},
		{
			title: "NotAvailableWhenObservedGenerationStale",
			ds: &appsv1.DaemonSet{
				ObjectMeta: metav1.ObjectMeta{Generation: 2},
				Status: appsv1.DaemonSetStatus{
					ObservedGeneration:     1,
					DesiredNumberScheduled: 3,
					UpdatedNumberScheduled: 3,
					NumberAvailable:        3,
				},
			},
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			actualResult := k8s.IsDaemonSetAvailable(tc.ds)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

const exampleDaemonSetYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: sample-ds
  namespace: %s
  labels:
    k8s-app: sample-ds
spec:
  selector:
    matchLabels:
      name: sample-ds
  template:
    metadata:
      labels:
        name: sample-ds
    spec:
      tolerations:
      - key: node-role.kubernetes.io/master
        effect: NoSchedule
      - key: node-role.kubernetes.io/control-plane
        effect: NoSchedule
      containers:
      - name: alpine
        image: alpine:3.8
        command: ['sh', '-c', 'echo Hello Terratest! && sleep 99999']
`
