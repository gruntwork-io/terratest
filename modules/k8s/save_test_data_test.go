package k8s_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndLoadKubectlOptions(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &k8s.KubectlOptions{
		ContextName: "terratest-context",
		ConfigPath:  "~/.kube/config",
		Namespace:   "default",
		Env: map[string]string{
			"TERRATEST_ENV_VAR": "terratest",
		},
	}
	k8s.SaveKubectlOptions(t, tmpFolder, expectedData)

	actualData := k8s.LoadKubectlOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}
