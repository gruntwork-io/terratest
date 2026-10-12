//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/shell"
	"github.com/gruntwork-io/terratest/modules/k8s/v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type KubectlVersion struct {
	ServerVersion struct {
		GitVersion string `json:"gitVersion"`
	} `json:"serverVersion"`
}

// kubectlServerVersion returns the server version kubectl reports. It reads stdout alone: when the
// client and server are more than one minor version apart, kubectl prints a skew warning on stderr,
// and that text would break the JSON.
func kubectlServerVersion(t *testing.T) string {
	t.Helper()

	out, err := shell.RunCommandContextAndGetStdOutE(t, context.Background(), &shell.Command{
		Command: "kubectl",
		Args:    []string{"version", "-o", "json"},
	})
	require.NoError(t, err)

	var kctlClusterVersion KubectlVersion
	require.NoError(t, json.Unmarshal([]byte(out), &kctlClusterVersion))

	return kctlClusterVersion.ServerVersion.GitVersion
}

func TestGetKubernetesClusterVersionE(t *testing.T) {
	t.Parallel()

	kubernetesClusterVersion, err := k8s.GetKubernetesClusterVersionContextE(t, context.Background())
	require.NoError(t, err)

	assert.Equal(t, kubernetesClusterVersion, kubectlServerVersion(t))
}

func TestGetKubernetesClusterVersionWithOptionsE(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "")
	kubernetesClusterVersion, err := k8s.GetKubernetesClusterVersionWithOptionsContextE(t, context.Background(), options)
	require.NoError(t, err)

	assert.Equal(t, kubernetesClusterVersion, kubectlServerVersion(t))
}
