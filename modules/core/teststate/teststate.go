// Package teststate provides primitives for serializing test data to disk so that it can be created in one test stage
// and reused in later validation and teardown stages.
//
// These primitives are type agnostic. Modules that own a type typically wrap them in a named helper (for example,
// aws.SaveEc2KeyPair or k8s.SaveKubectlOptions) so that callers do not have to spell out paths or filenames. Callers
// with types that have no owning module can use Save and Load directly.
//
// This package lives in core rather than teststructure so that modules such as aws, k8s, packer, and ssh can provide
// their own helpers without teststructure having to import every one of them.
package teststate

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// DirName is the name of the folder, relative to a test folder, in which test data is stored.
const DirName = ".test-data"

// redactedPlaceholder stands in for a value whose contents must not reach the test log.
const redactedPlaceholder = "[REDACTED]"

// FormatPath formats a path to save test data with the given filename in the given test folder.
func FormatPath(testFolder string, filename string) string {
	return filepath.Join(testFolder, DirName, filename)
}

// Save serializes and saves a value used at test time to the given path. This allows you to create some sort of test
// data (e.g., terraform.Options) during setup and to reuse this data later during validation and teardown. If
// `overwrite` is `true`, any contents that exist in the file found at `path` will be overwritten. This has the
// potential for causing duplicated resources and should be used with caution. If `overwrite` is `false`, the save will
// be skipped and a warning will be logged.
//
// The marshalled JSON is written to the test log. Use SaveRedacted for values that contain secrets.
func Save(t testing.TestingT, path string, overwrite bool, value any) {
	save(t, path, overwrite, value, true)
}

// SaveRedacted behaves exactly like Save, except that the marshalled JSON is not written to the test log. Use this for
// values that contain secrets, such as private keys.
func SaveRedacted(t testing.TestingT, path string, overwrite bool, value any) {
	save(t, path, overwrite, value, false)
}

// save serializes and saves a value used at test time to the given path. If `loggedVal` is `true`, the marshalled JSON
// is written to the test log.
func save(t testing.TestingT, path string, overwrite bool, value any, loggedVal bool) {
	logger.Default.Logf(t, "Storing test data in %s so it can be reused later", path)

	// The overwrite warnings render the value, so a redacted save must show a placeholder instead. Otherwise a
	// secret suppressed from the "Marshalled JSON" line below would still reach the log through this branch.
	loggedRepr := any(redactedPlaceholder)
	if loggedVal {
		loggedRepr = value
	}

	if IsPresent(t, path) {
		if overwrite {
			logger.Default.Logf(t, "[WARNING] The named test data at path %s is non-empty. Save operation will overwrite existing value with \"%v\".\n.", path, loggedRepr)
		} else {
			logger.Default.Logf(t, "[WARNING] The named test data at path %s is non-empty. Skipping save operation to prevent overwriting existing value with \"%v\".\n.", path, loggedRepr)

			return
		}
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to convert value %s to JSON: %v", path, err)
	}

	if loggedVal {
		logger.Default.Logf(t, "Marshalled JSON: %s", string(bytes))
	}

	parentDir := filepath.Dir(path)

	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		t.Fatalf("Failed to create folder %s: %v", parentDir, err)
	}

	// 0o600: this file can hold secrets, such as the private key in an aws.Ec2Keypair or an ssh.KeyPair.
	if err := os.WriteFile(path, bytes, 0o600); err != nil {
		t.Fatalf("Failed to save value %s: %v", path, err)
	}
}

// Load loads and unserializes a value stored at the given path. The value should be a pointer to a struct into which
// the value will be deserialized. This allows you to reuse some sort of test data (e.g., terraform.Options) from
// earlier setup steps in later validation and teardown steps.
func Load(t testing.TestingT, path string, value any) {
	logger.Default.Logf(t, "Loading test data from %s", path)

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load value from %s: %v", path, err)
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		t.Fatalf("Failed to parse JSON for value %s: %v", path, err)
	}
}

// IsPresent returns true if a file exists at $path and the test data there is non-empty.
func IsPresent(t testing.TestingT, path string) bool {
	exists, err := files.FileExistsE(path)
	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}

	if !exists {
		return false
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}

	if IsEmptyJSON(t, bytes) {
		return false
	}

	return true
}

// IsEmptyJSON returns true if the given bytes are empty, or in a valid JSON format that can reasonably be considered empty.
// The types used are based on the type possibilities listed at https://golang.org/src/encoding/json/decode.go?s=4062:4110#L51
func IsEmptyJSON(t testing.TestingT, bytes []byte) bool {
	var value any

	if len(bytes) == 0 {
		return true
	}

	if err := json.Unmarshal(bytes, &value); err != nil {
		t.Fatalf("Failed to parse JSON while testing whether it is empty: %v", err)
	}

	if value == nil {
		return true
	}

	valueBool, ok := value.(bool)
	if ok && !valueBool {
		return true
	}

	valueFloat64, ok := value.(float64)
	if ok && valueFloat64 == 0 {
		return true
	}

	valueString, ok := value.(string)
	if ok && valueString == "" {
		return true
	}

	valueSlice, ok := value.([]any)
	if ok && len(valueSlice) == 0 {
		return true
	}

	valueMap, ok := value.(map[string]any)
	if ok && len(valueMap) == 0 {
		return true
	}

	return false
}

// Cleanup cleans up the test data at the given path.
func Cleanup(t testing.TestingT, path string) {
	if files.FileExists(path) {
		logger.Default.Logf(t, "Cleaning up test data from %s", path)

		if err := os.Remove(path); err != nil {
			t.Fatalf("Failed to clean up file at %s: %v", path, err)
		}
	} else {
		logger.Default.Logf(t, "%s does not exist. Nothing to cleanup.", path)
	}
}

// CleanupFolder cleans up the .test-data folder inside the given folder.
// If there are any errors, fail the test.
func CleanupFolder(t testing.TestingT, path string) {
	err := CleanupFolderE(t, path)
	require.NoError(t, err)
}

// CleanupFolderE cleans up the .test-data folder inside the given folder.
func CleanupFolderE(t testing.TestingT, path string) error {
	path = filepath.Join(path, DirName)

	exists, err := files.FileExistsE(path)
	if err != nil {
		logger.Default.Logf(t, "Failed to clean up test data folder at %s: %v", path, err)

		return err
	}

	if !exists {
		logger.Default.Logf(t, "%s does not exist. Nothing to cleanup.", path)

		return nil
	}

	if err := os.RemoveAll(path); err != nil {
		logger.Default.Logf(t, "Failed to clean up test data folder at %s: %v", path, err)

		return err
	}

	return nil
}
