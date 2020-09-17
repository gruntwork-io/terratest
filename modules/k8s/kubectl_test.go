// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

// Test that RunKubectlAndGetOutputE will run kubectl and return the output by running a can-i command call.
func TestRunKubectlAndGetOutputReturnsOutput(t *testing.T) {
	namespaceName := fmt.Sprintf("kubectl-test-%s", strings.ToLower(random.UniqueId()))
	options := NewKubectlOptions("", "", namespaceName)
	output, err := RunKubectlAndGetOutputE(t, options, "auth", "can-i", "get", "pods")
	require.NoError(t, err)
	require.Equal(t, output, "yes")
}
