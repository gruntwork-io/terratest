package test_structure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyToTempFolder(t *testing.T) {
	tempFolder := CopyTerraformFolderToTemp(t, "../../", "examples")
	t.Log(tempFolder)
}

func TestCopySubtestToTempFolder(t *testing.T) {
	t.Run("Subtest", func(t *testing.T) {
		tempFolder := CopyTerraformFolderToTemp(t, "../../", "examples")
		t.Log(tempFolder)
	})
}

// TestValidateAllTerraformModulesSucceedsOnValidTerraform points at a simple text fixture Terraform module that is
// known to be valid
func TestValidateAllTerraformModulesSucceedsOnValidTerraform(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Use the test fixtures directory as the RootDir for ValidationOptions
	projectRootDir := filepath.Join(cwd, "../../test/fixtures")

	opts, optsErr := NewValidationOptions(projectRootDir, []string{"terraform-validation-valid"}, []string{})
	require.NoError(t, optsErr)

	ValidateAllTerraformModules(t, opts)
}

func TestNewValidationOptionsRejectsEmptyRootDir(t *testing.T) {
	_, err := NewValidationOptions("", []string{}, []string{})
	require.Error(t, err)
}

func TestFindTerraformModulePathsInRootEExamples(t *testing.T) {
	cwd, cwdErr := os.Getwd()
	require.NoError(t, cwdErr)

	opts, optsErr := NewValidationOptions(filepath.Join(cwd, "../../"), []string{}, []string{})
	require.NoError(t, optsErr)

	subDirs, err := FindTerraformModulePathsInRootE(opts)
	require.NoError(t, err)
	// There are many valid Terraform modules in the root/examples directory of the Terratest project, so we should get back many results
	require.Greater(t, len(subDirs), 0)
}

// This test calls ValidateAllTerraformModules on the Terratest root directory
func TestValidateAllTerraformModulesOnTerratest(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	projectRootDir := filepath.Join(cwd, "../..")

	opts, optsErr := NewValidationOptions(projectRootDir, []string{}, []string{
		"test/fixtures/terraform-with-plan-error",
		"test/fixtures/terragrunt/terragrunt-with-plan-error",
		"examples/terraform-backend-example",
	})
	require.NoError(t, optsErr)

	ValidateAllTerraformModules(t, opts)
}

// Verify ExcludeDirs is working properly, by explicitly passing a list of two test fixture modules to exclude
// and ensuring at the end that they do not appear in the returned slice of sub directories to validate
// Then, re-run the function with no exclusions and ensure the excluded paths ARE returned in the result set when no
// exclusions are passed
func TestFindTerraformModulePathsInRootEWithResultsExclusion(t *testing.T) {

	cwd, cwdErr := os.Getwd()
	require.NoError(t, cwdErr)

	projectRootDir := filepath.Join(cwd, "../..")

	// First, call the FindTerraformModulePathsInRootE method with several exclusions
	exclusions := []string{
		filepath.Join("test", "fixtures", "terraform-output"),
		filepath.Join("test", "fixtures", "terraform-output-map"),
	}

	opts, optsErr := NewValidationOptions(projectRootDir, []string{}, exclusions)
	require.NoError(t, optsErr)

	subDirs, err := FindTerraformModulePathsInRootE(opts)
	require.NoError(t, err)
	require.Greater(t, len(subDirs), 0)
	// Ensure none of the excluded paths were returned by FindTerraformModulePathsInRootE
	for _, exclusion := range exclusions {
		assert.False(t, collections.ListContains(subDirs, filepath.Join(projectRootDir, exclusion)))
	}

	// Next, call the same function but this time without exclusions and ensure that the excluded paths
	// exist in the non-excluded result set
	optsWithoutExclusions, optswoErr := NewValidationOptions(projectRootDir, []string{}, []string{})
	require.NoError(t, optswoErr)

	subDirsWithoutExclusions, woExErr := FindTerraformModulePathsInRootE(optsWithoutExclusions)
	require.NoError(t, woExErr)
	require.Greater(t, len(subDirsWithoutExclusions), 0)
	for _, exclusion := range exclusions {
		assert.True(t, collections.ListContains(subDirsWithoutExclusions, filepath.Join(projectRootDir, exclusion)))
	}
}
