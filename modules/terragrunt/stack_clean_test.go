package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackClean(t *testing.T) {
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

	// Generate stack to create .terragrunt-stack directory
	_, err = TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify that the .terragrunt-stack directory was created
	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")
	require.DirExists(t, stackDir)

	// Now run stack clean
	out, err := TgStackCleanE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Validate that clean command produced output
	require.Contains(t, out, "Cleaning stack")
}

func TestTerragruntStackCleanWithNoColor(t *testing.T) {
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

	// Generate stack to create .terragrunt-stack directory
	_, err = TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Run clean with no-color option
	out, err := TgStackCleanE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		ExtraArgs:        []string{"-no-color"},
	})
	require.NoError(t, err)

	// Validate that clean command produced output
	require.Contains(t, out, "Cleaning stack")
}

func TestTerragruntStackCleanNonExistentDir(t *testing.T) {
	t.Parallel()

	// Test with non-existent directory
	_, err := TgStackCleanE(t, &Options{
		TerragruntDir:    "/non/existent/path",
		TerragruntBinary: "terragrunt",
	})
	require.Error(t, err)
}

func TestTerragruntStackCleanEmptyOptions(t *testing.T) {
	t.Parallel()

	// Test with minimal options to verify default behavior
	_, err := TgStackCleanE(t, &Options{})
	require.Error(t, err)
	// Should fail due to missing TerragruntDir
}
