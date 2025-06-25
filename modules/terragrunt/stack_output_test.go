package terragrunt

import (
	"path"
	"runtime"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackOutput(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with no-color option
	out, err := TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		NoColor:          true,
	})
	require.NoError(t, err)

	// Get the output of stack output command
	out, err = TgStackOutputE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")
}

func TestTerragruntStackOutputError(t *testing.T) {
	t.Parallel()

	// Test with invalid terragrunt binary to ensure errors are caught
	_, err := TgStackOutputE(t, &Options{
		TerragruntDir:    "/nonexistent/path",
		TerragruntBinary: "nonexistent-binary",
	})
	require.Error(t, err)
	if runtime.GOOS == "linux" {
		require.Contains(t, err.Error(), "executable file not found in $PATH")
	}
}

func TestTerragruntStackOutputWithExtraArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with extra args
	out, err := TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		ExtraArgs: ExtraArgs{
			Apply: []string{"--terragrunt-log-level", "info"},
		},
	})
	require.NoError(t, err)

	// Validate that generate command produced output
	require.Contains(t, out, "Generating stack from")
	require.Contains(t, out, "Processing unit")

	// Verify that the .terragrunt-stack directory was created
	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")
	require.DirExists(t, stackDir)
}
