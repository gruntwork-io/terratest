//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetConfigMapEReturnsErrorForNonExistantConfigMap(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetConfigMapE(t, options, "test-config-map")
	require.Error(t, err)
}

func TestGetConfigMapEReturnsCorrectConfigMapInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleConfigMapYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	configMap := k8s.GetConfigMap(t, options, "test-config-map")
	require.Equal(t, "test-config-map", configMap.Name)
	require.Equal(t, configMap.Namespace, uniqueID)
}

func TestWaitUntilConfigMapAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleConfigMapYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)
	k8s.WaitUntilConfigMapAvailable(t, options, "test-config-map", 10, 1*time.Second)
}

const exampleConfigMapYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: test-config-map
  namespace: %s
`
