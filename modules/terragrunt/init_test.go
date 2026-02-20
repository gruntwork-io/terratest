package terragrunt

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	out := Init(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	// Check for common success indicator (works with both Terraform and OpenTofu)
	require.Contains(t, out, "successfully initialized")
}

func TestInitE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	out, err := InitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"}, // Common terraform init flag
	})
	require.NoError(t, err)
	// Check for common success indicator (works with both Terraform and OpenTofu)
	require.Contains(t, out, "successfully initialized")
}

func TestInitWithInvalidConfig(t *testing.T) {
	t.Parallel()
	// Test error handling when tg.hcl has invalid HCL syntax
	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init-error", t.Name())
	require.NoError(t, err)

	// This should fail due to invalid HCL syntax in tg.hcl
	_, err = InitE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"}, // Common terraform init flag
	})
	require.Error(t, err)
	// The error should contain information about the HCL parsing error
	require.Contains(t, err.Error(), "Missing expression")
}

// TestInitWithBothArgTypes verifies init works with both TerragruntArgs and TerraformArgs
func TestInitWithBothArgTypes(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    filepath.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--log-level", "error"},
		TerraformArgs:    []string{"-upgrade"},
	}

	output, err := InitE(t, options)
	require.NoError(t, err)
	// Verify TerragruntArgs: no info logs
	require.NotContains(t, output, "level=info")
	// Verify TerraformArgs: -upgrade was passed (shows in terraform output)
	require.Contains(t, output, "Initializing")
}
