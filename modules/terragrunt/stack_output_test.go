package terragrunt

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration test using actual terragrunt stack fixture
func TestStackOutputIntegration(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply tg using stack commands
	_, err = InitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		TerraformArgs:    []string{"apply"}, // stack run auto-approves by default
	}
	_, err = StackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			TerraformArgs:    []string{"destroy"}, // stack run auto-approves by default
		}
		_, _ = StackRunE(t, destroyOptions)
	}()

	// Test string stack output - get output from mother unit
	strOutput := StackOutput(t, options, "mother")
	assert.Contains(t, strOutput, "./test.txt")

	// Test getting stack output as JSON using the StackOutputJson function
	jsonOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	strOutputJson := StackOutputJson(t, jsonOptions, "mother")
	// The JSON output for a single value should still be cleaned to just show the value
	assert.Contains(t, strOutputJson, "./test.txt")

	// Test getting all stack outputs as JSON
	allOutputsJson := StackOutputJson(t, jsonOptions, "")
	require.NotEmpty(t, allOutputsJson)

	// For JSON output of all outputs, we should get valid JSON
	// But our function cleans it, so let's test it as-is
	// The JSON structure should be valid and contain our expected data
	if strings.Contains(allOutputsJson, "{") {
		// Parse and validate the JSON structure
		var allOutputs map[string]interface{}
		err = json.Unmarshal([]byte(allOutputsJson), &allOutputs)
		require.NoError(t, err)

		// Verify all expected stack outputs are present
		require.Contains(t, allOutputs, "mother")
		require.Contains(t, allOutputs, "father")
		require.Contains(t, allOutputs, "chick_1")
		require.Contains(t, allOutputs, "chick_2")

		// Verify the structure of outputs
		motherOutputMap := allOutputs["mother"].(map[string]interface{})
		assert.Equal(t, "./test.txt", motherOutputMap["output"])
	} else {
		// If not JSON format, at least verify it contains our expected values
		assert.Contains(t, allOutputsJson, "mother")
		assert.Contains(t, allOutputsJson, "father")
		assert.Contains(t, allOutputsJson, "chick_1")
		assert.Contains(t, allOutputsJson, "chick_2")
	}
}

// Test error handling with non-existent stack output
func TestStackOutputErrorHandling(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-error-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply tg using stack commands
	_, err = InitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		TerraformArgs:    []string{"apply"}, // stack run auto-approves by default
	}
	_, err = StackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			TerraformArgs:    []string{"destroy"}, // stack run auto-approves by default
		}
		_, _ = StackRunE(t, destroyOptions)
	}()

	// Test that non-existent stack output returns error or empty string
	output, err := StackOutputE(t, options, "non_existent_output")
	// Tg stack output might return empty string for non-existent outputs
	// rather than an error, so we need to handle both cases
	if err != nil {
		assert.Contains(t, strings.ToLower(err.Error()), "output")
	} else {
		assert.Empty(t, output, "Expected empty output for non-existent stack output")
	}
}

// Test StackOutputAll to get all stack outputs as a map
func TestStackOutputAll(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-all-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply tg using stack commands
	_, err = InitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		TerraformArgs:    []string{"apply"}, // stack run auto-approves by default
	}
	_, err = StackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			TerraformArgs:    []string{"destroy"}, // stack run auto-approves by default
		}
		_, _ = StackRunE(t, destroyOptions)
	}()

	// Test StackOutputAll - get all outputs as a map
	allOutputs := StackOutputAll(t, options)
	require.NotEmpty(t, allOutputs)

	// Verify expected outputs are present
	require.Contains(t, allOutputs, "mother")
	require.Contains(t, allOutputs, "father")
	require.Contains(t, allOutputs, "chick_1")
	require.Contains(t, allOutputs, "chick_2")

	// Verify we can access specific output values
	motherOutput := allOutputs["mother"].(map[string]interface{})
	assert.Equal(t, "./test.txt", motherOutput["output"])
}

// Test StackOutputListAll to get all stack output keys
func TestStackOutputListAll(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-list-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply using stack commands
	_, err = InitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		TerraformArgs:    []string{"apply"},
	}
	_, err = StackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			TerraformArgs:    []string{"destroy"},
		}
		_, _ = StackRunE(t, destroyOptions)
	}()

	// Test StackOutputListAll - get all output keys
	keys := StackOutputListAll(t, options)
	require.NotEmpty(t, keys)

	// Verify expected keys are present
	require.Contains(t, keys, "mother")
	require.Contains(t, keys, "father")
	require.Contains(t, keys, "chick_1")
	require.Contains(t, keys, "chick_2")

	// Verify we got all 4 keys
	require.Len(t, keys, 4)
}

// Test StackOutputListAllE
func TestStackOutputListAllE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-list-e-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	_, err = InitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		TerraformArgs:    []string{"apply"},
	}
	_, err = StackRunE(t, applyOptions)
	require.NoError(t, err)

	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			TerraformArgs:    []string{"destroy"},
		}
		_, _ = StackRunE(t, destroyOptions)
	}()

	keys, err := StackOutputListAllE(t, options)
	require.NoError(t, err)
	require.NotEmpty(t, keys)
	require.Contains(t, keys, "mother")
}
