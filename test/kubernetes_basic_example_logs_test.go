//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setupLogsTest(t *testing.T) (*k8s.KubectlOptions, v1.Pod) {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../examples/kubernetes-basic-example/podinfo-daemonset.yml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := strings.ToLower(random.UniqueId())

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	options := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, options, namespaceName)
	// ... and make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespace(t, options, namespaceName)

	// At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, options, kubeResourcePath)

	// Wait for at least 1 Pod to be ready from the DaemonSet
	retries := 10
	sleep := time.Second * 1
	for i := 1; i < retries; i++ {
		podsReady := k8s.GetDaemonSet(t, options, "podinfo-deamonset").Status.NumberReady
		if podsReady > 0 {
			break
		}
		time.Sleep(sleep)
	}

	// listOptions are used to select the pods with label app=podinfo
	listOptions := new(metav1.ListOptions)
	listOptions.LabelSelector = "app=podinfo"

	// Get a list of Pods. The pods are not guaranteed to be in running state.
	pods := k8s.ListPods(t, options, *listOptions)

	// Check that we did not timeout waiting for the Pod of the DaemonSet to be ready
	require.Greater(t, len(pods), 0)

	pod := pods[0]

	// Wait fot the pod to be started and ready
	k8s.WaitUntilPodAvailable(t, options, pod.Name, 5, 10*time.Second)

	return options, pod
}

func TestKubernetesBasicExampleLogsCheckWithContainerName(t *testing.T) {
	options, pod := setupLogsTest(t)
	logs := k8s.GetPodLogs(t, options, &pod, "podinfo")

	require.Contains(t, logs, "Starting podinfo")
}

func TestKubernetesBasicExampleLogsCheckWithNoContainerName(t *testing.T) {
	options, pod := setupLogsTest(t)
	logs := k8s.GetPodLogs(t, options, &pod, "")

	require.Contains(t, logs, "Starting podinfo")
}
