// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type KubectlVersion struct {
	ServerVersion struct {
		GitVersion string `json:"gitVersion"`
	} `json:"serverVersion"`
}

func TestGetKubernetesClusterVersionE(t *testing.T) {
	t.Parallel()

	kubernetesClusterVersion, err := GetKubernetesClusterVersionE(t)
	require.NoError(t, err)

	options := NewKubectlOptions("", "", "")
	kubernetesClusterVersionFromKubectl, err := RunKubectlAndGetOutputE(t, options, "version", "-o", "json")
	require.NoError(t, err)

	var kctlClusterVersion KubectlVersion
	require.NoError(
		t,
		json.Unmarshal([]byte(kubernetesClusterVersionFromKubectl), &kctlClusterVersion),
	)

	assert.EqualValues(t, kubernetesClusterVersion, kctlClusterVersion.ServerVersion.GitVersion)
}

func TestGetKubernetesClusterVersionWithOptionsE(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "")
	kubernetesClusterVersion, err := GetKubernetesClusterVersionWithOptionsE(t, options)
	require.NoError(t, err)

	kubernetesClusterVersionFromKubectl, err := RunKubectlAndGetOutputE(t, options, "version", "-o", "json")
	require.NoError(t, err)

	var kctlClusterVersion KubectlVersion
	require.NoError(
		t,
		json.Unmarshal([]byte(kubernetesClusterVersionFromKubectl), &kctlClusterVersion),
	)

	assert.EqualValues(t, kubernetesClusterVersion, kctlClusterVersion.ServerVersion.GitVersion)
}
