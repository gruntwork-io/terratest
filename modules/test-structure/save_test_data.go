package test_structure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// SaveTerraformOptions serializes and saves TerraformOptions into the given folder. This allows you to create TerraformOptions during setup
// and to reuse that TerraformOptions later during validation and teardown.
func SaveTerraformOptions(t testing.TestingT, testFolder string, terraformOptions *terraform.Options) {
	SaveTestData(t, formatTerraformOptionsPath(testFolder), terraformOptions)
}

// LoadTerraformOptions loads and unserializes TerraformOptions from the given folder. This allows you to reuse a TerraformOptions that was
// created during an earlier setup step in later validation and teardown steps.
func LoadTerraformOptions(t testing.TestingT, testFolder string) *terraform.Options {
	var terraformOptions terraform.Options
	LoadTestData(t, formatTerraformOptionsPath(testFolder), &terraformOptions)
	return &terraformOptions
}

// formatTerraformOptionsPath formats a path to save TerraformOptions in the given folder.
func formatTerraformOptionsPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "TerraformOptions.json")
}

// SavePackerOptions serializes and saves PackerOptions into the given folder. This allows you to create PackerOptions during setup
// and to reuse that PackerOptions later during validation and teardown.
func SavePackerOptions(t testing.TestingT, testFolder string, packerOptions *packer.Options) {
	SaveTestData(t, formatPackerOptionsPath(testFolder), packerOptions)
}

// LoadPackerOptions loads and unserializes PackerOptions from the given folder. This allows you to reuse a PackerOptions that was
// created during an earlier setup step in later validation and teardown steps.
func LoadPackerOptions(t testing.TestingT, testFolder string) *packer.Options {
	var packerOptions packer.Options
	LoadTestData(t, formatPackerOptionsPath(testFolder), &packerOptions)
	return &packerOptions
}

// formatPackerOptionsPath formats a path to save PackerOptions in the given folder.
func formatPackerOptionsPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "PackerOptions.json")
}

// SaveEc2KeyPair serializes and saves an Ec2KeyPair into the given folder. This allows you to create an Ec2KeyPair during setup
// and to reuse that Ec2KeyPair later during validation and teardown.
func SaveEc2KeyPair(t testing.TestingT, testFolder string, keyPair *aws.Ec2Keypair) {
	SaveTestData(t, formatEc2KeyPairPath(testFolder), keyPair)
}

// LoadEc2KeyPair loads and unserializes an Ec2KeyPair from the given folder. This allows you to reuse an Ec2KeyPair that was
// created during an earlier setup step in later validation and teardown steps.
func LoadEc2KeyPair(t testing.TestingT, testFolder string) *aws.Ec2Keypair {
	var keyPair aws.Ec2Keypair
	LoadTestData(t, formatEc2KeyPairPath(testFolder), &keyPair)
	return &keyPair
}

// formatEc2KeyPairPath formats a path to save an Ec2KeyPair in the given folder.
func formatEc2KeyPairPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "Ec2KeyPair.json")
}

// SaveSshKeyPair serializes and saves an SshKeyPair into the given folder. This allows you to create an SshKeyPair during setup
// and to reuse that SshKeyPair later during validation and teardown.
func SaveSshKeyPair(t testing.TestingT, testFolder string, keyPair *ssh.KeyPair) {
	SaveTestData(t, formatSshKeyPairPath(testFolder), keyPair)
}

// LoadSshKeyPair loads and unserializes an SshKeyPair from the given folder. This allows you to reuse an SshKeyPair that was
// created during an earlier setup step in later validation and teardown steps.
func LoadSshKeyPair(t testing.TestingT, testFolder string) *ssh.KeyPair {
	var keyPair ssh.KeyPair
	LoadTestData(t, formatSshKeyPairPath(testFolder), &keyPair)
	return &keyPair
}

// formatSshKeyPairPath formats a path to save an SshKeyPair in the given folder.
func formatSshKeyPairPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "SshKeyPair.json")
}

// SaveKubectlOptions serializes and saves KubectlOptions into the given folder. This allows you to create a KubectlOptions during setup
// and reuse that KubectlOptions later during validation and teardown.
func SaveKubectlOptions(t testing.TestingT, testFolder string, kubectlOptions *k8s.KubectlOptions) {
	SaveTestData(t, formatKubectlOptionsPath(testFolder), kubectlOptions)
}

// LoadKubectlOptions loads and unserializes a KubectlOptions from the given folder. This allows you to reuse a KubectlOptions that was
// created during an earlier setup step in later validation and teardown steps.
func LoadKubectlOptions(t testing.TestingT, testFolder string) *k8s.KubectlOptions {
	var kubectlOptions k8s.KubectlOptions
	LoadTestData(t, formatKubectlOptionsPath(testFolder), &kubectlOptions)
	return &kubectlOptions
}

// formatKubectlOptionsPath formats a path to save a KubectlOptions in the given folder.
func formatKubectlOptionsPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "KubectlOptions.json")
}

// SaveString serializes and saves a uniquely named string value into the given folder. This allows you to create one or more string
// values during one stage -- each with a unique name -- and to reuse those values during later stages.
func SaveString(t testing.TestingT, testFolder string, name string, val string) {
	path := formatNamedTestDataPath(testFolder, name)
	SaveTestData(t, path, val)
}

// LoadString loads and unserializes a uniquely named string value from the given folder. This allows you to reuse one or more string
// values that were created during an earlier setup step in later steps.
func LoadString(t testing.TestingT, testFolder string, name string) string {
	var val string
	LoadTestData(t, formatNamedTestDataPath(testFolder, name), &val)
	return val
}

// SaveInt saves a uniquely named int value into the given folder. This allows you to create one or more int
// values during one stage -- each with a unique name -- and to reuse those values during later stages.
func SaveInt(t testing.TestingT, testFolder string, name string, val int) {
	path := formatNamedTestDataPath(testFolder, name)
	SaveTestData(t, path, val)
}

// LoadInt loads a uniquely named int value from the given folder. This allows you to reuse one or more int
// values that were created during an earlier setup step in later steps.
func LoadInt(t testing.TestingT, testFolder string, name string) int {
	var val int
	LoadTestData(t, formatNamedTestDataPath(testFolder, name), &val)
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

// SaveAmiId serializes and saves an AMI ID into the given folder. This allows you to build an AMI during setup and to reuse that
// AMI later during validation and teardown.
//
// Deprecated: Use SaveArtifactID instead.
func SaveAmiId(t testing.TestingT, testFolder string, amiId string) {
	SaveString(t, testFolder, "AMI", amiId)
}

// LoadAmiId loads and unserializes an AMI ID from the given folder. This allows you to reuse an AMI  that was created during an
// earlier setup step in later validation and teardown steps.
//
// Deprecated: Use LoadArtifactID instead.
func LoadAmiId(t testing.TestingT, testFolder string) string {
	return LoadString(t, testFolder, "AMI")
}

// formatNamedTestDataPath formats a path to save an arbitrary named value in the given folder.
func formatNamedTestDataPath(testFolder string, name string) string {
	filename := fmt.Sprintf("%s.json", name)
	return FormatTestDataPath(testFolder, filename)
}

// FormatTestDataPath formats a path to save test data.
func FormatTestDataPath(testFolder string, filename string) string {
	return filepath.Join(testFolder, ".test-data", filename)
}

// SaveTestData serializes and saves a value used at test time to the given path. This allows you to create some sort of test data
// (e.g., TerraformOptions) during setup and to reuse this data later during validation and teardown.
func SaveTestData(t testing.TestingT, path string, value interface{}) {
	logger.Logf(t, "Storing test data in %s so it can be reused later", path)

	if IsTestDataPresent(t, path) {
		logger.Logf(t, "[WARNING] The named test data at path %s is non-empty. Save operation will overwrite existing value with \"%v\".\n.", path, value)
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to convert value %s to JSON: %v", path, err)
	}

	logger.Logf(t, "Marshalled JSON: %s", string(bytes))

	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0777); err != nil {
		t.Fatalf("Failed to create folder %s: %v", parentDir, err)
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		t.Fatalf("Failed to save value %s: %v", path, err)
	}
}

// LoadTestData loads and unserializes a value stored at the given path. The value should be a pointer to a struct into which the
// value will be deserialized. This allows you to reuse some sort of test data (e.g., TerraformOptions) from earlier
// setup steps in later validation and teardown steps.
func LoadTestData(t testing.TestingT, path string, value interface{}) {
	logger.Logf(t, "Loading test data from %s", path)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load value from %s: %v", path, err)
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		t.Fatalf("Failed to parse JSON for value %s: %v", path, err)
	}
}

// IsTestDataPresent returns true if a file exists at $path and the test data there is non-empty.
func IsTestDataPresent(t testing.TestingT, path string) bool {
	exists, err := files.FileExistsE(path)
	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}
	if !exists {
		return false
	}

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}

	if isEmptyJSON(t, bytes) {
		return false
	}

	return true
}

// isEmptyJSON returns true if the given bytes are empty, or in a valid JSON format that can reasonably be considered empty.
// The types used are based on the type possibilities listed at https://golang.org/src/encoding/json/decode.go?s=4062:4110#L51
func isEmptyJSON(t testing.TestingT, bytes []byte) bool {
	var value interface{}

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

	valueSlice, ok := value.([]interface{})
	if ok && len(valueSlice) == 0 {
		return true
	}

	valueMap, ok := value.(map[string]interface{})
	if ok && len(valueMap) == 0 {
		return true
	}

	return false
}

// CleanupTestData cleans up the test data at the given path.
func CleanupTestData(t testing.TestingT, path string) {
	if files.FileExists(path) {
		logger.Logf(t, "Cleaning up test data from %s", path)
		if err := os.Remove(path); err != nil {
			t.Fatalf("Failed to clean up file at %s: %v", path, err)
		}
	} else {
		logger.Logf(t, "%s does not exist. Nothing to cleanup.", path)
	}
}

// CleanupTestDataFolder cleans up the .test-data folder inside the given folder.
// If there are any errors, fail the test.
func CleanupTestDataFolder(t testing.TestingT, path string) {
	err := CleanupTestDataFolderE(t, path)
	require.NoError(t, err)
}

// CleanupTestDataFolderE cleans up the .test-data folder inside the given folder.
func CleanupTestDataFolderE(t testing.TestingT, path string) error {
	path = filepath.Join(path, ".test-data")
	exists, err := files.FileExistsE(path)
	if err != nil {
		logger.Logf(t, "Failed to clean up test data folder at %s: %v", path, err)
		return err
	}

	if !exists {
		logger.Logf(t, "%s does not exist. Nothing to cleanup.", path)
		return nil
	}
	if err := os.RemoveAll(path); err != nil {
		logger.Logf(t, "Failed to clean up test data folder at %s: %v", path, err)
		return err
	}

	return nil
}
