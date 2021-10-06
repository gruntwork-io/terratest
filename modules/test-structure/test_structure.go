package test_structure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	go_test "testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/opa"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// SKIP_STAGE_ENV_VAR_PREFIX is the prefix used for skipping stage environment variables.
const SKIP_STAGE_ENV_VAR_PREFIX = "SKIP_"

// RunTestStage executes the given test stage (e.g., setup, teardown, validation) if an environment variable of the name
// `SKIP_<stageName>` (e.g., SKIP_teardown) is not set.
func RunTestStage(t testing.TestingT, stageName string, stage func()) {
	envVarName := fmt.Sprintf("%s%s", SKIP_STAGE_ENV_VAR_PREFIX, stageName)
	if os.Getenv(envVarName) == "" {
		logger.Logf(t, "The '%s' environment variable is not set, so executing stage '%s'.", envVarName, stageName)
		stage()
	} else {
		logger.Logf(t, "The '%s' environment variable is set, so skipping stage '%s'.", envVarName, stageName)
	}
}

// SkipStageEnvVarSet returns true if an environment variable is set instructing Terratest to skip a test stage. This can be an easy way
// to tell if the tests are running in a local dev environment vs a CI server.
func SkipStageEnvVarSet() bool {
	for _, environmentVariable := range os.Environ() {
		if strings.HasPrefix(environmentVariable, SKIP_STAGE_ENV_VAR_PREFIX) {
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
//       // Root folder where terraform files should be (relative to the test folder)
//       rootFolder := ".."
//
//       // Relative path to terraform module being tested from the root folder
//       terraformFolderRelativeToRoot := "examples/terraform-aws-example"
//
//       // Copy the terraform folder to a temp folder
//       tempTestFolder := test_structure.CopyTerraformFolderToTemp(t, rootFolder, terraformFolderRelativeToRoot)
//
//       // Make sure to use the temp test folder in the terraform options
//       terraformOptions := &terraform.Options{
//       		TerraformDir: tempTestFolder,
//       }
//
// Note that if any of the SKIP_<stage> environment variables is set, we assume this is a test in the local dev where
// there are no other concurrent tests running and we want to be able to cache test data between test stages, so in that
// case, we do NOT copy anything to a temp folder, and return the path to the original terraform module folder instead.
func CopyTerraformFolderToTemp(t testing.TestingT, rootFolder string, terraformModuleFolder string) string {
	if SkipStageEnvVarSet() {
		logger.Logf(t, "A SKIP_XXX environment variable is set. Using original examples folder rather than a temp folder so we can cache data between stages for faster local testing.")
		return filepath.Join(rootFolder, terraformModuleFolder)
	}

	fullTerraformModuleFolder := filepath.Join(rootFolder, terraformModuleFolder)

	exists, err := files.FileExistsE(fullTerraformModuleFolder)
	require.NoError(t, err)
	if !exists {
		t.Fatal(files.DirNotFoundError{Directory: fullTerraformModuleFolder})
	}

	tmpRootFolder, err := files.CopyTerraformFolderToTemp(rootFolder, cleanName(t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	tmpTestFolder := filepath.Join(tmpRootFolder, terraformModuleFolder)

	// Log temp folder so we can see it
	logger.Logf(t, "Copied terraform folder %s to %s", fullTerraformModuleFolder, tmpTestFolder)

	return tmpTestFolder
}

func cleanName(originalName string) string {
	parts := strings.Split(originalName, "/")
	return parts[len(parts)-1]
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
func ValidateAllTerraformModules(t *go_test.T, opts *ValidationOptions) {
	runValidateOnAllTerraformModules(
		t,
		opts,
		func(t *go_test.T, fileType ValidateFileType, tfOpts *terraform.Options) {
			if fileType == TG {
				tfOpts.TerraformBinary = "terragrunt"
				// First call init and terraform validate
				terraform.InitAndValidate(t, tfOpts)
				// Next, call terragrunt validate-inputs which will catch mis-aligned inputs provided via Terragrunt
				terraform.ValidateInputs(t, tfOpts)
			} else if fileType == TF {
				terraform.InitAndValidate(t, tfOpts)
			}
		},
	)
}

// OPAEvalAllTerraformModules automatically finds all folders specified in RootDir that contain .tf files and runs
// OPAEval in all of them. The behavior of this function is similar to ValidateAllTerraformModules. Refer to the docs of
// that function for more details.
func OPAEvalAllTerraformModules(
	t *go_test.T,
	opts *ValidationOptions,
	opaEvalOpts *opa.EvalOptions,
	resultQuery string,
) {
	if opts.FileType != TF {
		t.Fatalf("OPAEvalAllTerraformModules currently only works with Terraform modules")
	}
	runValidateOnAllTerraformModules(
		t,
		opts,
		func(t *go_test.T, _ ValidateFileType, tfOpts *terraform.Options) {
			terraform.OPAEval(t, tfOpts, opaEvalOpts, resultQuery)
		},
	)
}

// runValidateOnAllTerraformModules main driver for ValidateAllTerraformModules and OPAEvalAllTerraformModules. Refer to
// the function docs of ValidateAllTerraformModules for more details.
func runValidateOnAllTerraformModules(
	t *go_test.T,
	opts *ValidationOptions,
	validationFunc func(t *go_test.T, fileType ValidateFileType, tfOps *terraform.Options),
) {
	dirsToValidate, readErr := FindTerraformModulePathsInRootE(opts)
	require.NoError(t, readErr)

	for _, dir := range dirsToValidate {
		dir := dir
		t.Run(strings.TrimLeft(dir, "/"), func(t *go_test.T) {
			// Determine the absolute path to the git repository root
			cwd, cwdErr := os.Getwd()
			require.NoError(t, cwdErr)
			gitRoot, gitRootErr := filepath.Abs(filepath.Join(cwd, "../../"))
			require.NoError(t, gitRootErr)

			// Determine the relative path to the example, module, etc that is currently being considered
			relativePath, pathErr := filepath.Rel(gitRoot, dir)
			require.NoError(t, pathErr)
			// Copy git root to tmp and supply the path to the current module to run init and validate on
			testFolder := CopyTerraformFolderToTemp(t, gitRoot, relativePath)
			require.NotNil(t, testFolder)

			// Run Terraform init and terraform validate on the test folder that was copied to /tmp
			// to avoid any potential conflicts with tests that may not use the same copy to /tmp behavior
			tfOpts := &terraform.Options{TerraformDir: testFolder}
			validationFunc(t, opts.FileType, tfOpts)
		})
	}
}
