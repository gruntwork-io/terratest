package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// ListPods will look for pods in the given namespace that match the given filters and return them. This will fail the
// test if there is an error.
func ListPods(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) []corev1.Pod {
	pods, err := ListPodsE(t, options, filters)
	require.NoError(t, err)
	return pods
}

// ListPodsE will look for pods in the given namespace that match the given filters and return them.
func ListPodsE(t testing.TestingT, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().Pods(options.Namespace).List(context.Background(), filters)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetPod returns a Kubernetes pod resource in the provided namespace with the given name. This will
// fail the test if there is an error.
func GetPod(t testing.TestingT, options *KubectlOptions, podName string) *corev1.Pod {
	pod, err := GetPodE(t, options, podName)
	require.NoError(t, err)
	return pod
}

// GetPodE returns a Kubernetes pod resource in the provided namespace with the given name.
func GetPodE(t testing.TestingT, options *KubectlOptions, podName string) (*corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsE(t, options)
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Pods(options.Namespace).Get(context.Background(), podName, metav1.GetOptions{})
}

// WaitUntilNumPodsCreated waits until the desired number of pods are created that match the provided filter. This will
// retry the check for the specified amount of times, sleeping for the provided duration between each try. This will
// fail the test if the retry times out.
func WaitUntilNumPodsCreated(
	t testing.TestingT,
	options *KubectlOptions,
	filters metav1.ListOptions,
	desiredCount int,
	retries int,
	sleepBetweenRetries time.Duration,
) {
	require.NoError(t, WaitUntilNumPodsCreatedE(t, options, filters, desiredCount, retries, sleepBetweenRetries))
}

// WaitUntilNumPodsCreatedE waits until the desired number of pods are created that match the provided filter. This will
// retry the check for the specified amount of times, sleeping for the provided duration between each try.
func WaitUntilNumPodsCreatedE(
	t testing.TestingT,
	options *KubectlOptions,
	filters metav1.ListOptions,
	desiredCount int,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for num pods created to match desired count %d.", desiredCount)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pods, err := ListPodsE(t, options, filters)
			if err != nil {
				return "", err
			}
			if len(pods) != desiredCount {
				return "", DesiredNumberOfPodsNotCreated{Filter: filters, DesiredCount: desiredCount}
			}
			return "Desired number of Pods created", nil
		},
	)
	if err != nil {
		options.Logger.Logf(t, "Timedout waiting for the desired number of Pods to be created: %s", err)
		return err
	}
	options.Logger.Logf(t, message)
	return nil
}

// WaitUntilPodAvailable waits until all of the containers within the pod are ready and started, retrying the check for the specified amount of times, sleeping
// for the provided duration between each try. This will fail the test if there is an error or if the check times out.
func WaitUntilPodAvailable(t testing.TestingT, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) {
	require.NoError(t, WaitUntilPodAvailableE(t, options, podName, retries, sleepBetweenRetries))
}

// WaitUntilPodAvailableE waits until all of the containers within the pod are ready and started, retrying the check for the specified amount of times, sleeping
// for the provided duration between each try.
func WaitUntilPodAvailableE(t testing.TestingT, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for pod %s to be provisioned.", podName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pod, err := GetPodE(t, options, podName)
			if err != nil {
				return "", err
			}
			if !IsPodAvailable(pod) {
				return "", NewPodNotAvailableError(pod)
			}
			return "Pod is now available", nil
		},
	)
	if err != nil {
		options.Logger.Logf(t, "Timedout waiting for Pod to be provisioned: %s", err)
		return err
	}
	options.Logger.Logf(t, message)
	return nil
}

// IsPodAvailable returns true if the all of the containers within the pod are ready and started
func IsPodAvailable(pod *corev1.Pod) bool {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		isContainerStarted := containerStatus.Started
		isContainerReady := containerStatus.Ready

		if !isContainerReady || (isContainerStarted != nil && *isContainerStarted == false) {
			return false
		}
	}
	return pod.Status.Phase == corev1.PodRunning
}

// GetPodLogsE returns the logs of a Pod at the time when the function was called. Pass container name if there are more containers in the Pod or set to "" if there is only one.
// If the Pod is not running an Error is returned.
// If the provided containerName is not the name of a container in the Pod an Error is returned.
func GetPodLogsE(t testing.TestingT, options *KubectlOptions, pod *corev1.Pod, containerName string) (string, error) {
	var output string
	var err error
	if containerName == "" {
		output, err = RunKubectlAndGetOutputE(t, options, "logs", pod.Name)
	} else {
		output, err = RunKubectlAndGetOutputE(t, options, "logs", pod.Name, fmt.Sprintf("-c%s", containerName))
	}

	if err != nil {
		return "", err
	}
	return output, nil
}

// GetPodLogsE returns the logs of a Pod at the time when the function was called.  Pass container name if there are more containers in the Pod or set to "" if there is only one.
func GetPodLogs(t testing.TestingT, options *KubectlOptions, pod *corev1.Pod, containerName string) string {
	logs, err := GetPodLogsE(t, options, pod, containerName)
	require.NoError(t, err)
	return logs
}
