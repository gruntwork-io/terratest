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
// Options carrying a RestConfig cannot be saved and will fail the test. RestConfig is not serializable: beyond its
// func-typed fields it holds interfaces and exec credential plugin wiring that cannot be rebuilt from JSON, and the
// fields that would survive are the credentials themselves, which have no business being written to .test-data.
// Failing here is deliberate. Dropping the config silently would let LoadKubectlOptions return options that fall
// back to the ambient kubeconfig and authenticate against a different cluster than the one under test.
func SaveKubectlOptions(t testing.TestingT, testFolder string, kubectlOptions *KubectlOptions) {
	if kubectlOptions != nil && kubectlOptions.RestConfig != nil {
		t.Fatalf(
			"SaveKubectlOptions cannot save options built with a RestConfig, because a rest.Config cannot be "+
				"serialized. Save the values needed to rebuild it instead, or use options built from a kubeconfig "+
				"path and context name. Path that would have been written: %s",
			formatKubectlOptionsPath(testFolder),
		)

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
