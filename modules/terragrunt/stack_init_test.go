package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTerragruntStackInit(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	out, err := TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")
}

func TestTerragruntStackInitError(t *testing.T) {
	t.Parallel()

	// Test with invalid terragrunt binary to ensure errors are caught
	_, err := TgStackInitE(t, &Options{
		TerragruntDir:    "/nonexistent/path",
		TerragruntBinary: "nonexistent-binary",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid binary")
}

func TestTerragruntStackInitWithInvalidConfig(t *testing.T) {
	t.Parallel()
	// Test error handling when terragrunt.hcl has invalid HCL syntax
	// The .terragrunt-stack directory is missing because terragrunt stack init
	// fails to parse the malformed terragrunt.hcl file, preventing the stack
	// infrastructure from being created
	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-init-error", t.Name())
	require.NoError(t, err)

	// This should fail due to missing .terragrunt-stack directory during stack init
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
	require.Error(t, err)
	// The error should contain information about the missing .terragrunt-stack directory
	// This indicates that terragrunt is trying to run stack commands but can't find
	// the required stack infrastructure
	require.Contains(t, err.Error(), ".terragrunt-stack")
}
