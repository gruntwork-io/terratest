package packer

import (
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/core/v2/teststate"
)

// packerOptionsFilename is the name of the file, within a test folder's test data directory, used to store an Options.
const packerOptionsFilename = "PackerOptions.json"

// SavePackerOptions serializes and saves PackerOptions into the given folder. This allows you to create PackerOptions
// during setup and to reuse that PackerOptions later during validation and teardown.
func SavePackerOptions(t testing.TestingT, testFolder string, packerOptions *Options) {
	teststate.Save(t, formatPackerOptionsPath(testFolder), true, packerOptions)
}

// LoadPackerOptions loads and unserializes PackerOptions from the given folder. This allows you to reuse a
// PackerOptions that was created during an earlier setup step in later validation and teardown steps.
func LoadPackerOptions(t testing.TestingT, testFolder string) *Options {
	var packerOptions Options
	teststate.Load(t, formatPackerOptionsPath(testFolder), &packerOptions)

	return &packerOptions
}

// formatPackerOptionsPath formats a path to save a PackerOptions in the given folder.
func formatPackerOptionsPath(testFolder string) string {
	return teststate.FormatPath(testFolder, packerOptionsFilename)
}
