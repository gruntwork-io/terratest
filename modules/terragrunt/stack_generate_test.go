package terragrunt

import (
	"path"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestStackGenerate(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	Init(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})

	out := StackGenerate(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})

	require.True(t, containsEitherString(out, "Processing unit", "Generating unit"))
	require.DirExists(t, path.Join(testFolder, "live", ".terragrunt-stack"))
}

func TestStackGenerateE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = InitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Then generate the stack
	out, err := StackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
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

func TestStackGenerateNonExistentDir(t *testing.T) {
	t.Parallel()

	// Test with non-existent directory
	_, err := StackGenerateE(t, &Options{
		TerragruntDir:    "/non/existent/path",
		TerragruntBinary: "terragrunt",
	})
	require.Error(t, err)
}

// containsEitherString checks if the output contains at least one of the provided strings
func containsEitherString(output, str1, str2 string) bool {
	return strings.Contains(output, str1) || strings.Contains(output, str2)
}

// TestStackGenerateWithArgs verifies stack commands respect TerragruntArgs
func TestStackGenerateWithArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// Initialize first
	_, err = InitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with TerragruntArgs
	out, err := StackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--log-level", "error"},
	})
	require.NoError(t, err)
	// Verify args were respected
	require.NotContains(t, out, "level=info")
}
