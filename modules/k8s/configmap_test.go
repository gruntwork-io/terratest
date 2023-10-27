//go:build kubeall || kubernetes
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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetConfigMapEReturnsErrorForNonExistantConfigMap(t *testing.T) {
	t.Parallel()

	options := NewKubectlOptions("", "", "default")
	_, err := GetConfigMapE(t, options, "test-config-map")
	require.Error(t, err)
}

func TestGetConfigMapEReturnsCorrectConfigMapInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_CONFIGMAP_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)
	KubectlApplyFromString(t, options, configData)

	configMap := GetConfigMap(t, options, "test-config-map")
	require.Equal(t, configMap.Name, "test-config-map")
	require.Equal(t, configMap.Namespace, uniqueID)
}

func TestWaitUntilConfigMapAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	options := NewKubectlOptions("", "", uniqueID)
	configData := fmt.Sprintf(EXAMPLE_CONFIGMAP_YAML_TEMPLATE, uniqueID, uniqueID)
	defer KubectlDeleteFromString(t, options, configData)

	KubectlApplyFromString(t, options, configData)
	WaitUntilConfigMapAvailable(t, options, "test-config-map", 10, 1*time.Second)
}

const EXAMPLE_CONFIGMAP_YAML_TEMPLATE = `---
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
