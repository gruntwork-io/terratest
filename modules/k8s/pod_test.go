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
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestListPodsReturnsPodsInNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	pods := k8s.ListPods(t, options, metav1.ListOptions{})
	require.Len(t, pods, 1)
	pod := pods[0]
	require.Equal(t, "nginx-pod", pod.Name)
	require.Equal(t, pod.Namespace, uniqueID)
}

func TestGetPodEReturnsErrorForNonExistantPod(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetPodE(t, options, "nginx-pod")
	require.Error(t, err)
}

func TestGetPodEReturnsCorrectPodInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	pod := k8s.GetPod(t, options, "nginx-pod")
	require.Equal(t, "nginx-pod", pod.Name)
	require.Equal(t, pod.Namespace, uniqueID)
}

func TestWaitUntilNumPodsCreatedReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	k8s.WaitUntilNumPodsCreated(t, options, metav1.ListOptions{}, 1, 60, 1*time.Second)
}

func TestWaitUntilPodAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	k8s.WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)
}

func TestWaitUntilPodWithMultipleContainersAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodWithMultipleContainersYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	k8s.WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)
}

func TestWaitUntilPodAvailableWithReadinessProbe(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodWithReadinessProbe, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	k8s.WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)
}

func TestWaitUntilPodAvailableWithFailingReadinessProbe(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodWithFailingReadinessProbe, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	err := k8s.WaitUntilPodAvailableE(t, options, "nginx-pod", 60, 1*time.Second)
	require.Error(t, err)
}

const examplePodYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  namespace: %s
spec:
  containers:
  - name: nginx
    image: nginx:1.15.7
    env:
        - name: NAME
          value: "nginx"
    ports:
    - containerPort: 80
`

const examplePodWithMultipleContainersYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  namespace: %s
spec:
  containers:
  - name: nginx
    image: nginx:1.15.7
    env:
        - name: NAME
          value: "nginx"
    ports:
    - containerPort: 80
  - name: nginx-two
    image: nginx:1.15.7
    env:
        - name: NAME
          value: "nginx-two"
    ports:
    - containerPort: 8080
    command: ["sh", "-c", "sed -i 's/80/8080/' /etc/nginx/conf.d/default.conf && nginx -g 'daemon off;'"]
`

const examplePodWithReadinessProbe = examplePodYAMLTemplate + `
    readinessProbe:
      httpGet:
        path: /
        port: 80
`

const examplePodWithFailingReadinessProbe = examplePodYAMLTemplate + `
    readinessProbe:
      httpGet:
        path: /not-ready
        port: 80
      periodSeconds: 1
`

func TestIsPodAvailable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		pod            *corev1.Pod
		title          string
		expectedResult bool
	}{
		{
			title: "TestIsPodAvailableStartedButNotReady",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "container1"}},
				},
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:    "container1",
							Ready:   false,
							Started: &[]bool{true}[0],
						},
					},
					Phase: corev1.PodRunning,
				},
			},
			expectedResult: false,
		},
		{
			title: "TestIsPodAvailableStartedAndReady",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "container1"}},
				},
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:    "container1",
							Ready:   true,
							Started: &[]bool{true}[0],
						},
					},
					Phase: corev1.PodRunning,
				},
			},
			expectedResult: true,
		},
		{
			title: "TestIsPodAvailableMissingContainerStatus",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Name: "container1"}, {Name: "container2"}},
				},
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:    "container1",
							Ready:   true,
							Started: &[]bool{true}[0],
						},
					},
					Phase: corev1.PodRunning,
				},
			},
			expectedResult: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()

			actualResult := k8s.IsPodAvailable(tc.pod)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestExecPod(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(examplePodWithMultipleContainersYAMLTemplate, uniqueID, uniqueID)

	t.Cleanup(func() { k8s.KubectlDeleteFromString(t, options, configData) })

	k8s.KubectlApplyFromString(t, options, configData)

	k8s.WaitUntilPodAvailable(t, options, "nginx-pod", 60, 1*time.Second)

	t.Run("TestExecPodWithoutContainer", func(t *testing.T) {
		t.Parallel()

		stdout, err := k8s.ExecPodE(t, options, "nginx-pod", "", "env")
		require.NoError(t, err)
		require.Contains(t, stdout, "NAME=nginx\n")
	})

	t.Run("TestExecPodWithContainer", func(t *testing.T) {
		t.Parallel()

		stdout, err := k8s.ExecPodE(t, options, "nginx-pod", "nginx-two", "env")
		require.NoError(t, err)
		require.Contains(t, stdout, "NAME=nginx-two\n")
	})
}
