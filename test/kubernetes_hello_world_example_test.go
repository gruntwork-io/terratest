//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: See the notes in the other Kubernetes example tests for why this build tag is included.

package test_test

import (
	"testing"
	"time"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestKubernetesHelloWorldExample(t *testing.T) {
	t.Parallel()

	// website::tag::1:: Path to the Kubernetes resource config we will test.
	kubeResourcePath := "../examples/kubernetes-hello-world-example/hello-world-deployment.yml"

	// website::tag::2:: Setup the kubectl config and context.
	options := k8s.NewKubectlOptions("", "", "default")

	// website::tag::6:: At the end of the test, run "kubectl delete" to clean up any resources that were created.
	defer k8s.KubectlDeleteContext(t, t.Context(), options, kubeResourcePath)

	// website::tag::3:: Run `kubectl apply` to deploy. Fail the test if there are any errors.
	k8s.KubectlApplyContext(t, t.Context(), options, kubeResourcePath)

	// website::tag::4:: Verify the service is available and get the URL for it.
	k8s.WaitUntilServiceAvailableContext(t, t.Context(), options, "hello-world-service", 10, 1*time.Second)
	service := k8s.GetServiceContext(t, t.Context(), options, "hello-world-service")
	url := "http://" + k8s.GetServiceEndpointContext(t, t.Context(), options, service, 5000)

	// website::tag::5:: Make an HTTP request to the URL and make sure it returns a 200 OK with the body "Hello, World!".
	http_helper.HTTPGetWithRetryContext(t, t.Context(), url, nil, 200, "Hello, World!", 30, 3*time.Second)
}
