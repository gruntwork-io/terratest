//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test

import (
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setupFailingDeploymentTest(t *testing.T, file string) error {
	t.Parallel()

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs(path.Join("fixtures/kubernetes-failing-deployment", file))
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
	//// ... and make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespace(t, options, namespaceName)
	//
	//// At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, options, kubeResourcePath)

	// listOptions are used to select the pods with label app=podinfo
	listOptions := metav1.ListOptions{
		LabelSelector: "app=test-app",
	}

	// Wait for at least 1 Pod to be ready from the DaemonSet
	k8s.WaitUntilNumPodsCreated(t, options, listOptions, 1, 5, time.Second)

	// Get a list of Pods. The pods are not guaranteed to be in running state.
	pods := k8s.ListPods(t, options, listOptions)

	// Check that we did not timeout waiting for the Pod of the DaemonSet to be ready
	require.Greater(t, len(pods), 0)

	pod := pods[0]

	// Wait fot the pod to be started and ready
	err = k8s.WaitUntilPodConsistentlyAvailableE(t, options, pod.Name, 5, time.Second, 5)
	require.Error(t, err)

	return err
}

func TestKubernetesUnknownImage(t *testing.T) {
	err := setupFailingDeploymentTest(t, "unknown-image-deployment.yml")
	notAvailableError := k8s.PodNotAvailable{}
	require.ErrorAs(t, err, &notAvailableError)
	msg := notAvailableError.Error()
	require.Contains(t, msg, "ContainersNotReady")
	require.Contains(t, msg, "waiting(app):ErrImagePull")
	require.Contains(t, msg, "failed to pull and unpack image \"docker.io/not-an-place/not-an-app:latest\"")
}

func TestKubernetesTooMuchResource(t *testing.T) {
	err := setupFailingDeploymentTest(t, "too-much-resource-deployment.yml")
	notAvailableError := k8s.PodNotAvailable{}
	require.ErrorAs(t, err, &notAvailableError)
	msg := notAvailableError.Error()
	require.Contains(t, msg, "PodScheduled:Unschedulable")
}

func TestKubernetesFailsAfterASec(t *testing.T) {
	err := setupFailingDeploymentTest(t, "probe-slow-to-pass-deployment.yml")
	notAvailableError := k8s.PodNotAvailable{}
	require.ErrorAs(t, err, &notAvailableError)
	msg := notAvailableError.Error()
	require.Contains(t, msg, "Ready:ContainersNotReady")
}

func TestKubernetesHappyThenFails(t *testing.T) {
	err := setupFailingDeploymentTest(t, "happy-then-fail-deployment.yml")
	notAvailableError := k8s.PodNotAvailable{}
	require.ErrorAs(t, err, &notAvailableError)
	msg := notAvailableError.Error()
	require.Contains(t, msg, "Ready:ContainersNotReady")
	require.Contains(t, msg, "ContainersReady:ContainersNotReady->containers with unready status: [app]")
}

func TestKubernetesStopsAfter5s(t *testing.T) {
	err := setupFailingDeploymentTest(t, "container-stops-after-2s-deployment.yml")
	notAvailableError := k8s.PodNotAvailable{}
	require.ErrorAs(t, err, &notAvailableError)
	msg := notAvailableError.Error()
	require.Contains(t, msg, "Ready:ContainersNotReady")
	require.Contains(t, msg, "terminated(app):Completed")
}
