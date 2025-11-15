package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTgStackRunPlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Then run plan on the stack
	out, err := TgStackRunE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"plan"},
	})
	require.NoError(t, err)

	// Validate that generate command produced output
	// Terragrunt v0.80.4+ outputs "Processing unit", older versions output "Generating unit"
	require.True(t, containsEitherString(out, "Processing unit", "Generating unit"), "Output should contain either 'Processing unit' or 'Generating unit'")

	// Verify that the .terragrunt-stack directory was created
	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")
	require.DirExists(t, stackDir)

	// Verify that the expected unit directories were created
	expectedUnits := []string{"mother", "father", "chicks/chick-1", "chicks/chick-2"}
	for _, unit := range expectedUnits {
		unitPath := path.Join(stackDir, unit)
		require.DirExists(t, unitPath)
	}
}

func TestTgStackRunPlanWithNoColor(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Run plan with no-color option
	out, err := TgStackRunE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--no-color"},
		TerraformArgs:    []string{"plan"},
	})
	require.NoError(t, err)

	// Validate that generate command produced output
	// Terragrunt v0.80.4+ outputs "Processing unit", older versions output "Generating unit"
	require.True(t, containsEitherString(out, "Processing unit", "Generating unit"), "Output should contain either 'Processing unit' or 'Generating unit'")

	// Verify that the .terragrunt-stack directory was created
	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")
	require.DirExists(t, stackDir)
}

func TestTgStackRunNonExistentDir(t *testing.T) {
	t.Parallel()

	// Test with non-existent directory
	_, err := TgStackRunE(t, &Options{
		TerragruntDir:    "/non/existent/path",
		TerragruntBinary: "terragrunt",
	})
	require.Error(t, err)
}

func TestTgStackRunEmptyOptions(t *testing.T) {
	t.Parallel()

	// Test with minimal options to verify default behavior
	_, err := TgStackRunE(t, &Options{})
	require.Error(t, err)
	// Should fail due to missing TerragruntDir
}
