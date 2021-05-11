// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKubernetesClusterVersionE(t *testing.T) {
	t.Parallel()

	kubernetesClusterVersion, err := GetKubernetesClusterVersionE(t)
	assert.NoError(t, err)
	assert.NotEmpty(t, kubernetesClusterVersion)
}

func TestGetKubernetesClusterVersionWithOptionsE(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	kubernetesClusterVersion, err := GetKubernetesClusterVersionWithOptionsE(t, options)
	assert.NoError(t, err)
	assert.NotEmpty(t, kubernetesClusterVersion)
}
