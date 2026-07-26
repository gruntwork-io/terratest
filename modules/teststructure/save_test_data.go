package teststructure

import (
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/core/v2/teststate"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
)

// SaveTerraformOptions serializes and saves TerraformOptions into the given folder. This allows you to create TerraformOptions during setup
// and to reuse that TerraformOptions later during validation and teardown.
func SaveTerraformOptions(t testing.TestingT, testFolder string, terraformOptions *terraform.Options) {
	teststate.Save(t, formatTerraformOptionsPath(testFolder), true, terraformOptions)
}

// SaveTerraformOptionsIfNotPresent serializes and saves TerraformOptions into the given folder if the file does not exist or the json is
// empty. This allows you to create TerraformOptions during setup and to reuse that TerraformOptions later during validation and teardown,
// but will prevent overwriting the contents and potentially duplicating resources.
func SaveTerraformOptionsIfNotPresent(t testing.TestingT, testFolder string, terraformOptions *terraform.Options) {
	teststate.Save(t, formatTerraformOptionsPath(testFolder), false, terraformOptions)
}

// LoadTerraformOptions loads and unserializes TerraformOptions from the given folder. This allows you to reuse a TerraformOptions that was
// created during an earlier setup step in later validation and teardown steps.
func LoadTerraformOptions(t testing.TestingT, testFolder string) *terraform.Options {
	var terraformOptions terraform.Options
	teststate.Load(t, formatTerraformOptionsPath(testFolder), &terraformOptions)

	return &terraformOptions
}

// formatTerraformOptionsPath formats a path to save TerraformOptions in the given folder.
func formatTerraformOptionsPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "TerraformOptions.json")
}

// SaveString serializes and saves a uniquely named string value into the given folder. This allows you to create one or more string
// values during one stage -- each with a unique name -- and to reuse those values during later stages.
func SaveString(t testing.TestingT, testFolder string, name string, val string) {
	path := formatNamedTestDataPath(testFolder, name)
	teststate.Save(t, path, true, val)
}

// LoadString loads and unserializes a uniquely named string value from the given folder. This allows you to reuse one or more string
// values that were created during an earlier setup step in later steps.
func LoadString(t testing.TestingT, testFolder string, name string) string {
	var val string
	teststate.Load(t, formatNamedTestDataPath(testFolder, name), &val)

	return val
}

// SaveInt saves a uniquely named int value into the given folder. This allows you to create one or more int
// values during one stage -- each with a unique name -- and to reuse those values during later stages.
func SaveInt(t testing.TestingT, testFolder string, name string, val int) {
	path := formatNamedTestDataPath(testFolder, name)
	teststate.Save(t, path, true, val)
}

// LoadInt loads a uniquely named int value from the given folder. This allows you to reuse one or more int
// values that were created during an earlier setup step in later steps.
func LoadInt(t testing.TestingT, testFolder string, name string) int {
	var val int
	teststate.Load(t, formatNamedTestDataPath(testFolder, name), &val)

	return val
}

// SaveArtifactID serializes and saves an Artifact ID into the given folder. This allows you to build an Artifact during setup and to reuse that
// Artifact later during validation and teardown.
func SaveArtifactID(t testing.TestingT, testFolder string, artifactID string) {
	SaveString(t, testFolder, "Artifact", artifactID)
}

// LoadArtifactID loads and unserializes an Artifact ID from the given folder. This allows you to reuse an Artifact that was created during an
// earlier setup step in later validation and teardown steps.
func LoadArtifactID(t testing.TestingT, testFolder string) string {
	return LoadString(t, testFolder, "Artifact")
}

// formatNamedTestDataPath formats a path to save an arbitrary named value in the given folder.
func formatNamedTestDataPath(testFolder string, name string) string {
	filename := name + ".json"

	return FormatTestDataPath(testFolder, filename)
}

// FormatTestDataPath formats a path to save test data.
func FormatTestDataPath(testFolder string, filename string) string {
	return teststate.FormatPath(testFolder, filename)
}

// SaveTestData serializes and saves a value used at test time to the given path. This allows you to create some sort of test data
// (e.g., TerraformOptions) during setup and to reuse this data later during validation and teardown. If `overwrite` is `true`,
// any contents that exist in the file found at `path` will be overwritten. This has the potential for causing duplicated resources
// and should be used with caution. If `overwrite` is `false`, the save will be skipped and a warning will be logged.
func SaveTestData(t testing.TestingT, path string, overwrite bool, value any) {
	teststate.Save(t, path, overwrite, value)
}

// LoadTestData loads and unserializes a value stored at the given path. The value should be a pointer to a struct into which the
// value will be deserialized. This allows you to reuse some sort of test data (e.g., TerraformOptions) from earlier
// setup steps in later validation and teardown steps.
func LoadTestData(t testing.TestingT, path string, value any) {
	teststate.Load(t, path, value)
}

// IsTestDataPresent returns true if a file exists at $path and the test data there is non-empty.
func IsTestDataPresent(t testing.TestingT, path string) bool {
	return teststate.IsPresent(t, path)
}

// IsEmptyJSON returns true if the given bytes are empty, or in a valid JSON format that can reasonably be considered empty.
// The types used are based on the type possibilities listed at https://golang.org/src/encoding/json/decode.go?s=4062:4110#L51
func IsEmptyJSON(t testing.TestingT, bytes []byte) bool {
	return teststate.IsEmptyJSON(t, bytes)
}

// CleanupTestData cleans up the test data at the given path.
func CleanupTestData(t testing.TestingT, path string) {
	teststate.Cleanup(t, path)
}

// CleanupTestDataFolder cleans up the .test-data folder inside the given folder.
// If there are any errors, fail the test.
func CleanupTestDataFolder(t testing.TestingT, path string) {
	teststate.CleanupFolder(t, path)
}

// CleanupTestDataFolderE cleans up the .test-data folder inside the given folder.
func CleanupTestDataFolderE(t testing.TestingT, path string) error {
	return teststate.CleanupFolderE(t, path)
}
