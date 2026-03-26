// Package test_structure provides helpers for structuring Terraform tests into stages.
//
// Test stages allow you to break up a test into setup, validation, and teardown phases that
// can be selectively skipped via environment variables (e.g., SKIP_teardown). This enables
// faster local development by caching test data between stages while still running full
// end-to-end tests in CI.
//
// The package also provides utilities for copying Terraform folders to temporary directories
// to avoid conflicts when running tests in parallel, and for recursively discovering and
// validating all Terraform modules under a given root directory.
package test_structure //nolint:staticcheck // package name determined by directory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gruntwork-io/terratest/modules/git"

	go_test "testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/opa"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// SkipStageEnvVarPrefix is the prefix used for skipping stage environment variables.
const SkipStageEnvVarPrefix = "SKIP_"

// SKIP_STAGE_ENV_VAR_PREFIX is the prefix used for skipping stage environment variables.
//
// Deprecated: Use [SkipStageEnvVarPrefix] instead.
const SKIP_STAGE_ENV_VAR_PREFIX = SkipStageEnvVarPrefix //nolint:staticcheck,revive // preserving existing constant name

// RunTestStage executes the given test stage (e.g., setup, teardown, validation) if an environment variable of the name
// `SKIP_<stageName>` (e.g., SKIP_teardown) is not set.
func RunTestStage(t testing.TestingT, stageName string, stage func()) {
	envVarName := fmt.Sprintf("%s%s", SkipStageEnvVarPrefix, stageName)
	if os.Getenv(envVarName) == "" {
		logger.Default.Logf(t, "The '%s' environment variable is not set, so executing stage '%s'.", envVarName, stageName)
		stage()
	} else {
		logger.Default.Logf(t, "The '%s' environment variable is set, so skipping stage '%s'.", envVarName, stageName)
	}
}

// SkipStageEnvVarSet returns true if an environment variable is set instructing Terratest to skip a test stage. This can be an easy way
// to tell if the tests are running in a local dev environment vs a CI server.
func SkipStageEnvVarSet() bool {
	for _, environmentVariable := range os.Environ() {
		if strings.HasPrefix(environmentVariable, SkipStageEnvVarPrefix) {
			return true
		}
	}

	return false
}

// CopyTerraformFolderToTemp copies the given root folder to a randomly-named temp folder and return the path to the
// given terraform modules folder within the new temp root folder. This is useful when running multiple tests in
// parallel against the same set of Terraform files to ensure the tests don't overwrite each other's .terraform working
// directory and terraform.tfstate files. To ensure relative paths work, we copy over the entire root folder to a temp
// folder, and then return the path within that temp folder to the given terraform module dir, which is where the actual
// test will be running.
// For example, suppose you had the target terraform folder you want to test in "/examples/terraform-aws-example"
// relative to the repo root. If your tests reside in the "/test" relative to the root, then you will use this as
// follows:
//
//	// Root folder where terraform files should be (relative to the test folder)
//	rootFolder := ".."
//
//	// Relative path to terraform module being tested from the root folder
//	terraformFolderRelativeToRoot := "examples/terraform-aws-example"
//
//	// Copy the terraform folder to a temp folder
//	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)
//
//	// Make sure to use the temp test folder in the terraform options
//	terraformOptions := &terraform.Options{
//			TerraformDir: tempTestFolder,
//	}
//
// Note that if any of the SKIP_<stage> environment variables is set, we assume this is a test in the local dev where
// there are no other concurrent tests running and we want to be able to cache test data between test stages, so in that
// case, we do NOT copy anything to a temp folder, and return the path to the original terraform module folder instead.
func CopyTerraformFolderToTemp(t testing.TestingT, rootFolder string, terraformModuleFolder string) string {
	return CopyTerraformFolderToDest(t, rootFolder, terraformModuleFolder, os.TempDir())
}

// CopyTerraformFolderToDest copies the given root folder to a randomly-named temp folder and return the path to the
// given terraform modules folder within the new temp root folder. This is useful when running multiple tests in
// parallel against the same set of Terraform files to ensure the tests don't overwrite each other's .terraform working
// directory and terraform.tfstate files. To ensure relative paths work, we copy over the entire root folder to a temp
// folder, and then return the path within that temp folder to the given terraform module dir, which is where the actual
// test will be running.
// For example, suppose you had the target terraform folder you want to test in "/examples/terraform-aws-example"
// relative to the repo root. If your tests reside in the "/test" relative to the root, then you will use this as
// follows:
//
//	// Destination for the copy of the files.  In this example we are using the Azure Dev Ops variable
//	// for the folder that is cleaned after each pipeline job.
//	destRootFolder := os.Getenv("AGENT_TEMPDIRECTORY")
//
//	// Root folder where terraform files should be (relative to the test folder)
//	rootFolder := ".."
//
//	// Relative path to terraform module being tested from the root folder
//	terraformFolderRelativeToRoot := "examples/terraform-aws-example"
//
//	// Copy the terraform folder to a temp folder
//	tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot, destRootFolder)
//
//	// Make sure to use the temp test folder in the terraform options
//	terraformOptions := &terraform.Options{
//			TerraformDir: tempTestFolder,
//	}
//
// Note that if any of the SKIP_<stage> environment variables is set, we assume this is a test in the local dev where
// there are no other concurrent tests running and we want to be able to cache test data between test stages, so in that
// case, we do NOT copy anything to a temp folder, and return the path to the original terraform module folder instead.
func CopyTerraformFolderToDest(t testing.TestingT, rootFolder string, terraformModuleFolder string, destRootFolder string) string {
	if SkipStageEnvVarSet() {
		logger.Default.Logf(t, "A SKIP_XXX environment variable is set. Using original examples folder rather than a temp folder so we can cache data between stages for faster local testing.")

		return filepath.Join(rootFolder, terraformModuleFolder)
	}

	fullTerraformModuleFolder := filepath.Join(rootFolder, terraformModuleFolder)

	exists, err := files.FileExistsE(fullTerraformModuleFolder)
	require.NoError(t, err)

	if !exists {
		t.Fatal(files.DirNotFoundError{Directory: fullTerraformModuleFolder})
	}

	tmpRootFolder, err := files.CopyTerraformFolderToDest(rootFolder, destRootFolder, cleanName(t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	tmpTestFolder := filepath.Join(tmpRootFolder, terraformModuleFolder)

	// Log temp folder so we can see it
	logger.Default.Logf(t, "Copied terraform folder %s to %s", fullTerraformModuleFolder, tmpTestFolder)

	return tmpTestFolder
}

func cleanName(originalName string) string {
	parts := strings.Split(originalName, "/")

	return parts[len(parts)-1]
}

// ValidateAllTerraformModulesContext automatically finds all folders specified in RootDir that contain .tf files and runs
// InitAndValidate in all of them. The provided context is passed through to the underlying command execution,
// allowing for timeout and cancellation control.
// Filters down to only those paths passed in ValidationOptions.IncludeDirs, if passed.
// Excludes any folders specified in the ValidationOptions.ExcludeDirs. IncludeDirs will take precedence over ExcludeDirs
// Use the NewValidationOptions method to pass relative paths for either of these options to have the full paths built
// Note that go_test is an alias to Golang's native testing package created to avoid naming conflicts with Terratest's
// own testing package. We are using the native testing.T here because Terratest's testing.T struct does not implement Run
// Note that we have opted to place the ValidateAllTerraformModulesContext function here instead of in the terraform package
// to avoid import cycling
func ValidateAllTerraformModulesContext(t *go_test.T, ctx context.Context, opts *ValidationOptions) {
	t.Helper()

	runValidateOnAllTerraformModulesContext(
		t,
		ctx,
		opts,
		func(t *go_test.T, ctx context.Context, _ ValidateFileType, tfOpts *terraform.Options) {
			t.Helper()
			terraform.InitAndValidateContext(t, ctx, tfOpts)
		},
	)
}

// ValidateAllTerraformModules automatically finds all folders specified in RootDir that contain .tf files and runs
// InitAndValidate in all of them.
// Filters down to only those paths passed in ValidationOptions.IncludeDirs, if passed.
// Excludes any folders specified in the ValidationOptions.ExcludeDirs. IncludeDirs will take precedence over ExcludeDirs
// Use the NewValidationOptions method to pass relative paths for either of these options to have the full paths built
// Note that go_test is an alias to Golang's native testing package created to avoid naming conflicts with Terratest's
// own testing package. We are using the native testing.T here because Terratest's testing.T struct does not implement Run
// Note that we have opted to place the ValidateAllTerraformModules function here instead of in the terraform package
// to avoid import cycling
//
// Deprecated: Use [ValidateAllTerraformModulesContext] instead.
func ValidateAllTerraformModules(t *go_test.T, opts *ValidationOptions) {
	t.Helper()

	ValidateAllTerraformModulesContext(t, context.Background(), opts)
}

// OPAEvalAllTerraformModulesContext automatically finds all folders specified in RootDir that contain .tf files and runs
// OPAEval in all of them. The provided context is passed through to the underlying command execution,
// allowing for timeout and cancellation control. The behavior of this function is similar to
// ValidateAllTerraformModulesContext. Refer to the docs of that function for more details.
func OPAEvalAllTerraformModulesContext(
	t *go_test.T,
	ctx context.Context,
	opts *ValidationOptions,
	opaEvalOpts *opa.EvalOptions,
	resultQuery string,
) {
	t.Helper()

	if opts.FileType != TF {
		t.Fatalf("OPAEvalAllTerraformModulesContext currently only works with Terraform modules")
	}

	runValidateOnAllTerraformModulesContext(
		t,
		ctx,
		opts,
		func(t *go_test.T, ctx context.Context, _ ValidateFileType, tfOpts *terraform.Options) {
			t.Helper()
			terraform.OPAEvalContext(t, ctx, tfOpts, opaEvalOpts, resultQuery)
		},
	)
}

// OPAEvalAllTerraformModules automatically finds all folders specified in RootDir that contain .tf files and runs
// OPAEval in all of them. The behavior of this function is similar to ValidateAllTerraformModules. Refer to the docs of
// that function for more details.
//
// Deprecated: Use [OPAEvalAllTerraformModulesContext] instead.
func OPAEvalAllTerraformModules(
	t *go_test.T,
	opts *ValidationOptions,
	opaEvalOpts *opa.EvalOptions,
	resultQuery string,
) {
	t.Helper()

	OPAEvalAllTerraformModulesContext(t, context.Background(), opts, opaEvalOpts, resultQuery)
}

// runValidateOnAllTerraformModulesContext is the main driver for ValidateAllTerraformModulesContext and
// OPAEvalAllTerraformModulesContext. Refer to the function docs of ValidateAllTerraformModulesContext for more details.
func runValidateOnAllTerraformModulesContext(
	t *go_test.T,
	ctx context.Context,
	opts *ValidationOptions,
	validationFunc func(t *go_test.T, ctx context.Context, fileType ValidateFileType, tfOps *terraform.Options),
) {
	t.Helper()

	// Find the Git root
	gitRoot, err := git.GetRepoRootForDirContextE(t, ctx, opts.RootDir)
	require.NoError(t, err)

	// Find the relative path between the root dir and the git root
	relPath, err := filepath.Rel(gitRoot, opts.RootDir)
	require.NoError(t, err)

	// Copy git root to tmp
	testFolder := CopyTerraformFolderToTemp(t, gitRoot, relPath)
	require.NotNil(t, testFolder)

	// Clone opts and override the root dir to the temp folder
	clonedOpts, err := CloneWithNewRootDir(opts, testFolder)
	require.NoError(t, err)

	// Find TF modules
	dirsToValidate, readErr := FindTerraformModulePathsInRootE(clonedOpts)
	require.NoError(t, readErr)

	for _, dir := range dirsToValidate {
		t.Run(strings.TrimLeft(dir, "/"), func(t *go_test.T) {
			// Run the validation function on the test folder that was copied to /tmp to avoid any potential conflicts
			// with tests that may not use the same copy to /tmp behavior
			tfOpts := &terraform.Options{TerraformDir: dir}
			validationFunc(t, ctx, clonedOpts.FileType, tfOpts)
		})
	}
}
