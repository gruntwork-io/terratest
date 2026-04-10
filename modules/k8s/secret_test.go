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

func TestGetSecretEReturnsErrorForNonExistantSecret(t *testing.T) {
	t.Parallel()

	options := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetSecretE(t, options, "master-password")
	require.Error(t, err)
}

func TestGetSecretEReturnsCorrectSecretInCorrectNamespace(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleSecretYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)

	secret := k8s.GetSecret(t, options, "master-password")
	require.Equal(t, "master-password", secret.Name)
	require.Equal(t, secret.Namespace, uniqueID)
}

func TestWaitUntilSecretAvailableReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", uniqueID)

	configData := fmt.Sprintf(exampleSecretYAMLTemplate, uniqueID, uniqueID)
	defer k8s.KubectlDeleteFromString(t, options, configData)

	k8s.KubectlApplyFromString(t, options, configData)
	k8s.WaitUntilSecretAvailable(t, options, "master-password", 10, 1*time.Second)
}

const exampleSecretYAMLTemplate = `---
apiVersion: v1
kind: Namespace
metadata:
  name: %s
---
apiVersion: v1
kind: Secret
metadata:
  name: master-password
  namespace: %s
`
