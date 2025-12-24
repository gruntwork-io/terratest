package terragrunt

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files/v2"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file demonstrates two approaches for testing Terragrunt configurations:
//
// 1. SINGLE-MODULE TESTING: Use the terraform module with TerraformBinary set to "terragrunt".
//    This works because terragrunt is a thin wrapper around terraform for single modules.
//    See: TestTerragruntExample, TestTerragruntConsole
//
// 2. MULTI-MODULE TESTING: Use the dedicated terragrunt module with ApplyAll/DestroyAll.
//    This is for testing multiple Terragrunt modules with dependencies using --all commands.
//    See: TestTerragruntMultiModuleExample

// TestTerragruntExample demonstrates testing a single Terragrunt module using the terraform package.
// For single-module testing, use terraform.Options with TerraformBinary set to "terragrunt".
func TestTerragruntExample(t *testing.T) {
	t.Parallel()

	// Copy the example folder to a temp folder to avoid state conflicts between parallel tests.
	testFolder, err := files.CopyTerragruntFolderToTemp("../../examples/terragrunt-example", t.Name())
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terragrunt module that will be tested.
		TerraformDir: testFolder,
		// Set the terraform binary path to terragrunt so that terratest uses terragrunt
		// instead of terraform. You must ensure that you have terragrunt downloaded and
		// available in your PATH.
		TerraformBinary: "terragrunt",
	})

	// Clean up resources with "terragrunt destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terragrunt apply". Under the hood, terragrunt will run "terraform init" and
	// "terraform apply". Fail the test if there are any errors.
	terraform.Apply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have
	// the expected values.
	// Note: When using terragrunt, OutputAll is recommended because terragrunt returns
	// all outputs in the full JSON format even when a specific key is requested.
	outputs := terraform.OutputAll(t, terraformOptions)
	assert.Equal(t, "one input another input", outputs["output"])
}

// TestTerragruntConsole demonstrates running terragrunt console command.
func TestTerragruntConsole(t *testing.T) {
	t.Parallel()

	// Copy the example folder to a temp folder to avoid state conflicts between parallel tests.
	testFolder, err := files.CopyTerragruntFolderToTemp("../../examples/terragrunt-example", t.Name())
	require.NoError(t, err)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir:    testFolder,
		TerraformBinary: "terragrunt",
		Stdin:           strings.NewReader("local.mylocal"),
	})

	defer terraform.Destroy(t, terraformOptions)

	// Run "terragrunt run -- console".
	out := terraform.RunTerraformCommand(t, terraformOptions, "run", "--", "console")
	assert.Contains(t, out, `"local variable named mylocal"`)
}

// TestTerragruntMultiModuleExample demonstrates testing multiple Terragrunt modules
// using the dedicated terragrunt package. Use this approach when you have multiple
// modules with dependencies that need to be applied/destroyed together using --all.
func TestTerragruntMultiModuleExample(t *testing.T) {
	t.Parallel()

	// Copy the entire example folder (including modules) to a temp folder.
	// We copy the parent folder because terragrunt.hcl files reference ../modules.
	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../examples/terragrunt-multi-module-example", t.Name())
	require.NoError(t, err)

	options := &Options{
		// Run from the live subfolder where the terragrunt configs are
		TerragruntDir: filepath.Join(testFolder, "live"),
		// Optional: Set log level for cleaner output
		TerragruntArgs: []string{"--log-level", "error"},
	}

	// Clean up all modules with "terragrunt destroy --all" at the end of the test.
	// DestroyAll respects the reverse dependency order.
	defer DestroyAll(t, options)

	// Run "terragrunt apply --all". This applies all modules in dependency order.
	ApplyAll(t, options)

	// Verify the plan shows no changes (infrastructure is up-to-date)
	exitCode := PlanAllExitCode(t, options)
	assert.Equal(t, 0, exitCode, "Plan should show no changes after apply")
}
