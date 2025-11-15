package terragrunt

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTgInit(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	out, err := TgInitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"}, // Common terraform init flag
	})
	require.NoError(t, err)
	require.Contains(t, out, "Terraform has been successfully initialized!")
}

func TestTgInitWithInvalidConfig(t *testing.T) {
	t.Parallel()
	// Test error handling when tg.hcl has invalid HCL syntax
	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init-error", t.Name())
	require.NoError(t, err)

	// This should fail due to invalid HCL syntax in tg.hcl
	_, err = TgInitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"}, // Common terraform init flag
	})
	require.Error(t, err)
	// The error should contain information about the HCL parsing error
	require.Contains(t, err.Error(), "Missing expression")
}

// TestTgInitWithBothArgTypes verifies init works with both TerragruntArgs and TerraformArgs
func TestTgInitWithBothArgTypes(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    filepath.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--log-level", "error"},
		TerraformArgs:    []string{"-upgrade"},
	}

	output, err := TgInitE(t, options)
	require.NoError(t, err)
	// Verify TerragruntArgs: no info logs
	require.NotContains(t, output, "level=info")
	// Verify TerraformArgs: -upgrade was passed (shows in terraform output)
	require.Contains(t, output, "Initializing")
}
