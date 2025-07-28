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
func TestTgOutputIntegration(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply terragrunt using stack commands
	_, err = TgStackInitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		ExtraArgs:        []string{"apply", "-auto-approve"},
	}
	_, err = TgStackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			ExtraArgs:        []string{"destroy", "-auto-approve"},
		}
		_, _ = TgStackRunE(t, destroyOptions)
	}()

	// Test string stack output - get output from mother unit
	strOutput := TgOutput(t, options, "mother")
	assert.Contains(t, strOutput, "./test.txt")

	// Test getting stack output as JSON - note that our cleaning function will still extract just the value
	jsonOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		ExtraArgs:        []string{"-json"},
	}

	strOutputJson := TgOutput(t, jsonOptions, "mother")
	// The JSON output for a single value should still be cleaned to just show the value
	assert.Contains(t, strOutputJson, "./test.txt")

	// Test getting all stack outputs as JSON
	allOutputsJson := TgOutput(t, jsonOptions, "")
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
func TestTgOutputErrorHandling(t *testing.T) {
	t.Parallel()

	// Create a temporary copy of the stack fixture
	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", "tg-stack-output-error-test")
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
	}

	// Initialize and apply terragrunt using stack commands
	_, err = TgStackInitE(t, options)
	require.NoError(t, err)

	applyOptions := &Options{
		TerragruntDir:    testFolder + "/live",
		TerragruntBinary: "terragrunt",
		Logger:           logger.Discard,
		ExtraArgs:        []string{"apply", "-auto-approve"},
	}
	_, err = TgStackRunE(t, applyOptions)
	require.NoError(t, err)

	// Clean up after test
	defer func() {
		destroyOptions := &Options{
			TerragruntDir:    testFolder + "/live",
			TerragruntBinary: "terragrunt",
			Logger:           logger.Discard,
			ExtraArgs:        []string{"destroy", "-auto-approve"},
		}
		_, _ = TgStackRunE(t, destroyOptions)
	}()

	// Test that non-existent stack output returns error or empty string
	output, err := TgOutputE(t, options, "non_existent_output")
	// Terragrunt stack output might return empty string for non-existent outputs
	// rather than an error, so we need to handle both cases
	if err != nil {
		assert.Contains(t, strings.ToLower(err.Error()), "output")
	} else {
		assert.Empty(t, output, "Expected empty output for non-existent stack output")
	}
}
