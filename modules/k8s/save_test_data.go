package k8s

import (
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/core/v2/teststate"
)

// kubectlOptionsFilename is the name of the file, within a test folder's test data directory, used to store a
// KubectlOptions.
const kubectlOptionsFilename = "KubectlOptions.json"

// SaveKubectlOptions serializes and saves KubectlOptions into the given folder. This allows you to create a
// KubectlOptions during setup and reuse that KubectlOptions later during validation and teardown.
//
// Options carrying a RestConfig cannot be saved and will fail the test, because a rest.Config cannot be rebuilt
// from JSON. Rebuild it in the later stage instead, or use options built from a kubeconfig path and context name.
func SaveKubectlOptions(t testing.TestingT, testFolder string, kubectlOptions *KubectlOptions) {
	if kubectlOptions != nil && kubectlOptions.RestConfig != nil {
		t.Fatalf("SaveKubectlOptions cannot save options built with a RestConfig, because a rest.Config cannot be " +
			"serialized. Rebuild it in the later stage instead, or use options built from a kubeconfig path and " +
			"context name.")

		return
	}

	teststate.Save(t, formatKubectlOptionsPath(testFolder), true, kubectlOptions)
}

// LoadKubectlOptions loads and unserializes a KubectlOptions from the given folder. This allows you to reuse a
// KubectlOptions that was created during an earlier setup step in later validation and teardown steps.
func LoadKubectlOptions(t testing.TestingT, testFolder string) *KubectlOptions {
	var kubectlOptions KubectlOptions
	teststate.Load(t, formatKubectlOptionsPath(testFolder), &kubectlOptions)

	return &kubectlOptions
}

// formatKubectlOptionsPath formats a path to save a KubectlOptions in the given folder.
func formatKubectlOptionsPath(testFolder string) string {
	return teststate.FormatPath(testFolder, kubectlOptionsFilename)
}
