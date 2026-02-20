package terragrunt

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

// TestTerragruntEndToEndIntegration is a comprehensive integration test that validates
// the complete terragrunt workflow with TerragruntArgs and TerraformArgs.
// This test exercises the fix for issue #1609 where args were being ignored.
func TestTerragruntEndToEndIntegration(t *testing.T) {
	t.Parallel()

	// Setup: Copy test fixture to temp directory
	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Configure options with TerragruntArgs
	options := &Options{
		TerragruntDir: testFolder,
		// TerragruntArgs: Global terragrunt flags that should be respected
		TerragruntArgs: []string{"--log-level", "error"},
	}

	// Step 1: Plan with exit code (original bug scenario from issue #1609)
	// This is the exact scenario from the bug report
	t.Log("Step 1: Testing PlanAllExitCode with TerragruntArgs (original bug scenario)")
	exitCode, err := PlanAllExitCodeE(t, options)
	require.NoError(t, err)
	// Should show changes (exit code 2) since nothing has been applied yet
	require.Equal(t, 2, exitCode, "Plan should detect changes")

	// Step 2: Apply all modules
	t.Log("Step 2: Testing ApplyAll with TerragruntArgs")
	applyOutput := ApplyAll(t, options)
	require.NotEmpty(t, applyOutput)
	// Verify TerragruntArgs: should not see info-level logs
	require.NotContains(t, applyOutput, "level=info", "TerragruntArgs should suppress info logs")

	// Step 3: Plan again - should show no changes (exit code 0)
	t.Log("Step 3: Verifying infrastructure is up-to-date")
	exitCode, err = PlanAllExitCodeE(t, options)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Plan should show no changes after apply")

	// Step 4: Clean up - Destroy all
	t.Log("Step 4: Testing DestroyAll with TerragruntArgs")
	destroyOutput := DestroyAll(t, options)
	require.NotEmpty(t, destroyOutput)
	// Verify TerragruntArgs: should not see info-level logs
	require.NotContains(t, destroyOutput, "level=info", "TerragruntArgs should suppress info logs")

	t.Log("Integration test completed successfully - all args were properly passed")
}

// TestStackEndToEndIntegration tests the complete stack workflow with args
func TestStackEndToEndIntegration(t *testing.T) {
	t.Parallel()

	// Setup: Copy stack test fixture
	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:  filepath.Join(testFolder, "live"),
		TerragruntArgs: []string{"--log-level", "error"},
	}

	// Step 1: Initialize stack
	t.Log("Step 1: Initializing stack with TerragruntArgs")
	output, err := InitE(t, options)
	require.NoError(t, err)
	require.NotContains(t, output, "level=info", "TerragruntArgs should suppress info logs")

	// Step 2: Generate stack
	t.Log("Step 2: Generating stack with TerragruntArgs")
	genOutput, err := StackGenerateE(t, options)
	require.NoError(t, err)
	require.NotContains(t, genOutput, "level=info", "TerragruntArgs should suppress info logs")

	// Step 3: Run stack plan
	t.Log("Step 3: Running stack plan with TerraformArgs")
	runOptions := *options
	runOptions.TerraformArgs = []string{"plan"}
	planOutput, err := StackRunE(t, &runOptions)
	require.NoError(t, err)
	// Check for common plan indicator (works with both Terraform and OpenTofu)
	require.Contains(t, planOutput, "will perform")

	// Step 4: Clean stack
	t.Log("Step 4: Cleaning stack")
	_, err = StackCleanE(t, options)
	require.NoError(t, err)

	t.Log("Stack integration test completed successfully")
}

// TestOutputAllJsonEndToEnd tests OutputAllJson extracts clean JSON from terragrunt output
func TestOutputAllJsonEndToEnd(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp(
		"testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{TerragruntDir: testFolder}

	ApplyAll(t, options)
	defer DestroyAll(t, options)

	output := OutputAllJson(t, options)

	// Contains module outputs, no log noise
	require.Contains(t, output, `"value": "foo"`)
	require.Contains(t, output, `"value": "bar"`)
	// Check for both old and new log format markers
	require.NotContains(t, output, "time=")
	require.NotContains(t, output, " INFO ")
	require.NotContains(t, output, " STDOUT ")
	require.NotContains(t, output, "Group 1")
	require.NotContains(t, output, "- Unit ")

	// Validate output contains at least 2 valid JSON objects (foo and bar modules)
	dec := json.NewDecoder(strings.NewReader(output))
	var jsonCount int
	for dec.More() {
		var obj json.RawMessage
		require.NoError(t, dec.Decode(&obj))
		jsonCount++
	}
	require.GreaterOrEqual(t, jsonCount, 2)
}
