//go:build kubernetes
// +build kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestListEventsEReturnsNilErrorWhenListingEvents(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "kube-system")
	events, err := k8s.ListEventsE(t, options, v1.ListOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, events)
}

func TestListEventsInNamespace(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "kube-system")
	events := k8s.ListEvents(t, options, v1.ListOptions{})
	require.NotEmpty(t, events)
}

func TestListEventsReturnsZeroEventsIfNoneCreated(t *testing.T) {
	t.Parallel()

	ns := "test-ns"

	options := k8s.NewKubectlOptions("", "", "")

	defer k8s.DeleteNamespace(t, options, ns)

	k8s.CreateNamespace(t, options, ns)

	options.Namespace = ns
	events := k8s.ListEvents(t, options, v1.ListOptions{})
	require.Empty(t, events)
}
