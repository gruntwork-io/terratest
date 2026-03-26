package terragrunt_test

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	"github.com/stretchr/testify/require"
)

func TestStackClean(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	terragrunt.StackGenerate(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})

	require.DirExists(t, stackDir)

	out := terragrunt.StackClean(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})

	require.Contains(t, out, "Deleting stack directory")
	require.NoDirExists(t, stackDir)
}

func TestStackCleanE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	// First generate the stack to create .terragrunt-stack directory
	_, err = terragrunt.StackGenerateE(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify that the .terragrunt-stack directory was created
	require.DirExists(t, stackDir)

	// Clean the stack
	out, err := terragrunt.StackCleanE(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify clean command produced expected output
	require.Contains(t, out, "Deleting stack directory")

	// Verify that the .terragrunt-stack directory was removed
	require.NoDirExists(t, stackDir)
}

func TestStackCleanNonExistentStack(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	// Verify that the .terragrunt-stack directory doesn't exist
	require.NoDirExists(t, stackDir)

	// Clean should succeed even if .terragrunt-stack doesn't exist
	_, err = terragrunt.StackCleanE(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
}

func TestStackCleanAfterRun(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	stackDir := path.Join(testFolder, "live", ".terragrunt-stack")

	// Initialize the stack
	_, err = terragrunt.InitE(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Run plan to generate the stack
	_, err = terragrunt.StackRunE(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"plan"},
	})
	require.NoError(t, err)

	// Verify that the .terragrunt-stack directory was created
	require.DirExists(t, stackDir)

	// Clean the stack
	out, err := terragrunt.StackCleanE(t, &terragrunt.Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Verify clean command produced expected output
	require.Contains(t, out, "Deleting stack directory")

	// Verify that the .terragrunt-stack directory was removed
	require.NoDirExists(t, stackDir)
}
