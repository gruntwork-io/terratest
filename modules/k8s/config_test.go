//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDeleteConfigContext(t *testing.T) {
	t.Parallel()

	path := StoreConfigToTempFile(t, BASIC_CONFIG_WITH_EXTRA_CONTEXT)
	defer os.Remove(path)

	err := DeleteConfigContextWithPathE(t, path, "extra_minikube")
	require.NoError(t, err)

	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	storedConfig := string(data)
	assert.Equal(t, BASIC_CONFIG, storedConfig)
}

func TestDeleteConfigContextWithAnotherContextRemaining(t *testing.T) {
	t.Parallel()

	path := StoreConfigToTempFile(t, BASIC_CONFIG_WITH_EXTRA_CONTEXT_NO_GARBAGE)
	defer os.Remove(path)

	err := DeleteConfigContextWithPathE(t, path, "extra_minikube")
	require.NoError(t, err)

	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	storedConfig := string(data)
	assert.Equal(t, EXPECTED_CONFIG_AFTER_EXTRA_MINIKUBE_DELETED_NO_GARBAGE, storedConfig)
}

func TestRemoveOrphanedClusterAndAuthInfoConfig(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			"TestExtraClusterRemoveOrphanedClusterAndAuthInfoed",
			BASIC_CONFIG_WITH_EXTRA_CLUSTER,
			BASIC_CONFIG,
		},
		{
			"TestExtraAuthInfoRemoveOrphanedClusterAndAuthInfoed",
			BASIC_CONFIG_WITH_EXTRA_AUTH_INFO,
			BASIC_CONFIG,
		},
	}
	for _, testCase := range testCases {
		// Capture range variable to scope within range
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			removeOrphanedClusterAndAuthInfoConfigTestFunc(t, testCase.in, testCase.out)
		})
	}
}

func removeOrphanedClusterAndAuthInfoConfigTestFunc(t *testing.T, inputConfig string, expectedOutputConfig string) {
	path := StoreConfigToTempFile(t, inputConfig)
	defer os.Remove(path)

	config := LoadConfigFromPath(path)
	rawConfig, err := config.RawConfig()
	require.NoError(t, err)
	RemoveOrphanedClusterAndAuthInfoConfig(&rawConfig)
	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false)
	require.NoError(t, err)
	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	storedConfig := string(data)
	assert.Equal(t, expectedOutputConfig, storedConfig)
}

// Various example configs used in testing the config manipulation functions

const BASIC_CONFIG = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
`

const BASIC_CONFIG_WITH_EXTRA_CLUSTER = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
`

const BASIC_CONFIG_WITH_EXTRA_AUTH_INFO = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
- name: extra_minikube
  user:
    client-certificate: /home/terratest/.minikube/extra_client.crt
    client-key: /home/terratest/.minikube/extra_client.key
`

const BASIC_CONFIG_WITH_EXTRA_CONTEXT = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
- context:
    cluster: extra_minikube
    user: extra_minikube
  name: extra_minikube
current-context: extra_minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
- name: extra_minikube
  user:
    client-certificate: /home/terratest/.minikube/extra_client.crt
    client-key: /home/terratest/.minikube/extra_client.key
`

const BASIC_CONFIG_WITH_EXTRA_CONTEXT_NO_GARBAGE = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
- context:
    cluster: extra_minikube
    user: extra_minikube
  name: extra_minikube
- context:
    cluster: extra_minikube
    user: minikube
  name: other_minikube

current-context: extra_minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
- name: extra_minikube
  user:
    client-certificate: /home/terratest/.minikube/extra_client.crt
    client-key: /home/terratest/.minikube/extra_client.key
`

const EXPECTED_CONFIG_AFTER_EXTRA_MINIKUBE_DELETED_NO_GARBAGE = `apiVersion: v1
clusters:
- cluster:
    certificate-authority: /home/terratest/.minikube/extra_ca.crt
    server: https://172.17.0.48:8443
  name: extra_minikube
- cluster:
    certificate-authority: /home/terratest/.minikube/ca.crt
    server: https://172.17.0.48:8443
  name: minikube
contexts:
- context:
    cluster: minikube
    user: minikube
  name: minikube
- context:
    cluster: extra_minikube
    user: minikube
  name: other_minikube
current-context: minikube
kind: Config
preferences: {}
users:
- name: minikube
  user:
    client-certificate: /home/terratest/.minikube/client.crt
    client-key: /home/terratest/.minikube/client.key
`
