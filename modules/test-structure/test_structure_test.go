package test_structure_test //nolint:staticcheck // package name determined by directory

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	teststructure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyToTempFolder(t *testing.T) {
	t.Parallel()

	tempFolder := teststructure.CopyTerraformFolderToTemp(t, "../../", "examples")
	t.Log(tempFolder)
}

func TestCopySubtestToTempFolder(t *testing.T) {
	t.Parallel()

	t.Run("Subtest", func(t *testing.T) {
		t.Parallel()

		tempFolder := teststructure.CopyTerraformFolderToTemp(t, "../../", "examples")
		t.Log(tempFolder)
	})
}

// TestValidateAllTerraformModulesSucceedsOnValidTerraform points at a simple text fixture Terraform module that is
// known to be valid
func TestValidateAllTerraformModulesSucceedsOnValidTerraform(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Use the test fixtures directory as the RootDir for ValidationOptions
	projectRootDir := filepath.Join(cwd, "../../test/fixtures")

	opts, optsErr := teststructure.NewValidationOptions(projectRootDir, []string{"terraform-validation-valid"}, []string{})
	require.NoError(t, optsErr)

	teststructure.ValidateAllTerraformModulesContext(t, t.Context(), opts)
}

func TestNewValidationOptionsRejectsEmptyRootDir(t *testing.T) {
	t.Parallel()

	_, err := teststructure.NewValidationOptions("", []string{}, []string{})
	require.Error(t, err)
}

func TestFindTerraformModulePathsInRootEExamples(t *testing.T) {
	t.Parallel()

	cwd, cwdErr := os.Getwd()
	require.NoError(t, cwdErr)

	opts, optsErr := teststructure.NewValidationOptions(filepath.Join(cwd, "../../"), []string{}, []string{})
	require.NoError(t, optsErr)

	subDirs, err := teststructure.FindTerraformModulePathsInRootE(opts)
	require.NoError(t, err)
	// There are many valid Terraform modules in the root/examples directory of the Terratest project, so we should get back many results
	require.NotEmpty(t, subDirs)
}

// This test calls ValidateAllTerraformModules on the Terratest root directory
func TestValidateAllTerraformModulesOnTerratest(t *testing.T) {
	t.Parallel()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	projectRootDir := filepath.Join(cwd, "../..")

	opts, optsErr := teststructure.NewValidationOptions(projectRootDir, []string{}, []string{
		"test/fixtures/terraform-with-plan-error",
		"modules/terragrunt/testdata/terragrunt-with-plan-error",
		"examples/terraform-backend-example",
	})
	require.NoError(t, optsErr)

	teststructure.ValidateAllTerraformModulesContext(t, t.Context(), opts)
}

// Verify ExcludeDirs is working properly, by explicitly passing a list of two test fixture modules to exclude
// and ensuring at the end that they do not appear in the returned slice of sub directories to validate
// Then, re-run the function with no exclusions and ensure the excluded paths ARE returned in the result set when no
// exclusions are passed
func TestFindTerraformModulePathsInRootEWithResultsExclusion(t *testing.T) {
	t.Parallel()

	cwd, cwdErr := os.Getwd()
	require.NoError(t, cwdErr)

	projectRootDir := filepath.Join(cwd, "../..")

	// First, call the FindTerraformModulePathsInRootE method with several exclusions
	exclusions := []string{
		filepath.Join("test", "fixtures", "terraform-output"),
		filepath.Join("test", "fixtures", "terraform-output-map"),
	}

	opts, optsErr := teststructure.NewValidationOptions(projectRootDir, []string{}, exclusions)
	require.NoError(t, optsErr)

	subDirs, err := teststructure.FindTerraformModulePathsInRootE(opts)
	require.NoError(t, err)
	require.NotEmpty(t, subDirs)

	// Ensure none of the excluded paths were returned by FindTerraformModulePathsInRootE
	for _, exclusion := range exclusions {
		assert.False(t, slices.Contains(subDirs, filepath.Join(projectRootDir, exclusion)))
	}

	// Next, call the same function but this time without exclusions and ensure that the excluded paths
	// exist in the non-excluded result set
	optsWithoutExclusions, optswoErr := teststructure.NewValidationOptions(projectRootDir, []string{}, []string{})
	require.NoError(t, optswoErr)

	subDirsWithoutExclusions, woExErr := teststructure.FindTerraformModulePathsInRootE(optsWithoutExclusions)
	require.NoError(t, woExErr)
	require.NotEmpty(t, subDirsWithoutExclusions)

	for _, exclusion := range exclusions {
		assert.True(t, slices.Contains(subDirsWithoutExclusions, filepath.Join(projectRootDir, exclusion)))
	}
}
