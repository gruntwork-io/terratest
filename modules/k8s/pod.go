package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// ListPodsContextE looks up pods in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListPodsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) ([]corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	resp, err := clientset.CoreV1().Pods(options.Namespace).List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// ListPodsContext looks up pods in the given namespace that match the given filters and return them.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func ListPodsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, filters metav1.ListOptions) []corev1.Pod {
	t.Helper()
	pods, err := ListPodsContextE(t, ctx, options, filters)
	require.NoError(t, err)

	return pods
}

// GetPodContextE returns a Kubernetes pod resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
func GetPodContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, podName string) (*corev1.Pod, error) {
	clientset, err := GetKubernetesClientFromOptionsContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return clientset.CoreV1().Pods(options.Namespace).Get(ctx, podName, metav1.GetOptions{})
}

// GetPodContext returns a Kubernetes pod resource in the provided namespace with the given name.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func GetPodContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, podName string) *corev1.Pod {
	t.Helper()
	pod, err := GetPodContextE(t, ctx, options, podName)
	require.NoError(t, err)

	return pod
}

// WaitUntilNumPodsCreatedContextE waits until the desired number of pods are created that match the provided filter.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func WaitUntilNumPodsCreatedContextE(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	filters metav1.ListOptions,
	desiredCount int,
	retries int,
	sleepBetweenRetries time.Duration,
) error {
	statusMsg := fmt.Sprintf("Wait for num pods created to match desired count %d.", desiredCount)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pods, err := ListPodsContextE(t, ctx, options, filters)
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

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilNumPodsCreatedContext waits until the desired number of pods are created that match the provided filter.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
//
//nolint:gocritic // hugeParam: cannot change public function signature
func WaitUntilNumPodsCreatedContext(
	t testing.TestingT,
	ctx context.Context,
	options *KubectlOptions,
	filters metav1.ListOptions,
	desiredCount int,
	retries int,
	sleepBetweenRetries time.Duration,
) {
	t.Helper()
	require.NoError(t, WaitUntilNumPodsCreatedContextE(t, ctx, options, filters, desiredCount, retries, sleepBetweenRetries))
}

// WaitUntilPodAvailableContextE waits until all of the containers within the pod are ready and started,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
func WaitUntilPodAvailableContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) error {
	statusMsg := fmt.Sprintf("Wait for pod %s to be provisioned.", podName)

	message, err := retry.DoWithRetryContextE(
		t,
		ctx,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pod, err := GetPodContextE(t, ctx, options, podName)
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

	options.Logger.Logf(t, "%s", message)

	return nil
}

// WaitUntilPodAvailableContext waits until all of the containers within the pod are ready and started,
// retrying the check for the specified amount of times, sleeping for the provided duration between each try.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error.
func WaitUntilPodAvailableContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) {
	t.Helper()
	require.NoError(t, WaitUntilPodAvailableContextE(t, ctx, options, podName, retries, sleepBetweenRetries))
}

// IsPodAvailable returns true if the all of the containers within the pod are ready and started
func IsPodAvailable(pod *corev1.Pod) bool {

	if len(pod.Status.ContainerStatuses) != len(pod.Spec.Containers) {
		return false
	}

	for i := range pod.Status.ContainerStatuses {
		isContainerStarted := pod.Status.ContainerStatuses[i].Started
		isContainerReady := pod.Status.ContainerStatuses[i].Ready

		if !isContainerReady || (isContainerStarted != nil && !*isContainerStarted) {
			return false
		}
	}

	return pod.Status.Phase == corev1.PodRunning
}

// GetPodLogsContextE returns the logs of a Pod at the time when the function was called.
// The ctx parameter supports cancellation and timeouts.
// Pass container name if there are more containers in the Pod or set to "" if there is only one.
// If the Pod is not running an Error is returned.
// If the provided containerName is not the name of a container in the Pod an Error is returned.
func GetPodLogsContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, pod *corev1.Pod, containerName string) (string, error) {
	var (
		output string
		err    error
	)

	if containerName == "" {
		output, err = RunKubectlAndGetOutputContextE(t, ctx, options, "logs", pod.Name)
	} else {
		output, err = RunKubectlAndGetOutputContextE(t, ctx, options, "logs", pod.Name, "-c"+containerName)
	}

	if err != nil {
		return "", err
	}

	return output, nil
}

// GetPodLogsContext returns the logs of a Pod at the time when the function was called.
// The ctx parameter supports cancellation and timeouts.
// Pass container name if there are more containers in the Pod or set to "" if there is only one.
// This will fail the test if there is an error.
func GetPodLogsContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, pod *corev1.Pod, containerName string) string {
	t.Helper()
	logs, err := GetPodLogsContextE(t, ctx, options, pod, containerName)
	require.NoError(t, err)

	return logs
}

// ExecPodContextE executes a command in a container within a Kubernetes pod and returns the output.
// The ctx parameter supports cancellation and timeouts.
// Set containerName to "" if there is only one container in the pod.
func ExecPodContextE(t testing.TestingT, ctx context.Context, options *KubectlOptions, podName string, containerName string, command ...string) (string, error) {
	var args []string
	if containerName == "" {
		args = append([]string{"exec", podName, "--"}, command...)
	} else {
		args = append([]string{"exec", podName, "-c" + containerName, "--"}, command...)
	}

	return RunKubectlAndGetOutputContextE(t, ctx, options, args...)
}

// ExecPodContext executes a command in a container within a Kubernetes pod and returns the output.
// The ctx parameter supports cancellation and timeouts.
// This will fail the test if there is an error. Set containerName to "" if there is only one container in the pod.
func ExecPodContext(t testing.TestingT, ctx context.Context, options *KubectlOptions, podName string, containerName string, command ...string) string {
	t.Helper()
	o, err := ExecPodContextE(t, ctx, options, podName, containerName, command...)
	require.NoError(t, err)

	return o
}
