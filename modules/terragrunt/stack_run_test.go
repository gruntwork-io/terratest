package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackRunPlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Then run plan on the stack
	out, err := TgStackRunE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"plan"},
	})
	require.NoError(t, err)

	// Validate that generate command produced output
	require.Contains(t, out, "Generating stack from")
	require.Contains(t, out, "Processing unit")

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

func TestTerragruntStackRunPlanWithNoColor(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Run plan with no-color option
	out, err := TgStackRunE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"plan", "-no-color"},
	})
	require.NoError(t, err)

	// Validate that generate command produced output
	require.Contains(t, out, "Generating stack from")
	require.Contains(t, out, "Processing unit")

	// Verify that the .terragrunt-stack directory was created
	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")
	require.DirExists(t, stackDir)
}

func TestTerragruntStackRunNonExistentDir(t *testing.T) {
	t.Parallel()

	// Test with non-existent directory
	_, err := TgStackRunE(t, &Options{
		TerragruntDir:    "/non/existent/path",
		TerragruntBinary: "terragrunt",
	})
	require.Error(t, err)
}

func TestTerragruntStackRunEmptyOptions(t *testing.T) {
	t.Parallel()

	// Test with minimal options to verify default behavior
	_, err := TgStackRunE(t, &Options{})
	require.Error(t, err)
	// Should fail due to missing TerragruntDir
}
