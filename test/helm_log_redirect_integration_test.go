//go:build kubeall || helm
// +build kubeall helm

// **NOTE**: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly, helm
// can overload the minikube system and thus interfere with the other kubernetes tests. Specifically, many of the tests
// start to fail with `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes
// tests and helm tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.
// We recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	tftesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Example how to redirect helm logs to custom logger
func TestHelmLogsRedirect(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../examples/helm-basic-example")
	require.NoError(t, err)

	// Namespace to deploy helm chart
	namespaceName := fmt.Sprintf("helm-logs-%s", strings.ToLower(random.UniqueId()))

	// Setup the kubectl config and context. Here we choose to use the defaults, which is:
	// - HOME/.kube/config for the kubectl config file
	// - Current context of the kubectl config file
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	defer k8s.DeleteNamespace(t, kubectlOptions, namespaceName)

	customLogger := helmLogger{}

	options := &helm.Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"containerImageRepo": "nginx",
			"containerImageTag":  "1.15.8",
		},
		Logger: logger.New(&customLogger),
	}

	// Generate a unique release to avoid conflicts with other tests
	releaseName := fmt.Sprintf(
		"nginx-service-%s",
		strings.ToLower(random.UniqueId()),
	)
	defer helm.Delete(t, options, releaseName, true)

	helm.Install(t, options, helmChartPath, releaseName)

	// Validate that logs were redirected to custom logger
	require.Contains(t, customLogger.logs, releaseName)
	require.Contains(t, customLogger.logs, "STATUS: deployed")
}

type helmLogger struct {
	logs string
}

func (c *helmLogger) Logf(t tftesting.TestingT, format string, args ...interface{}) {
	c.logs = fmt.Sprintf("%s\n%s", c.logs, fmt.Sprintf(format, args...))
}
