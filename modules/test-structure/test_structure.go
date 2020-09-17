// Package test_structure allows to set up tests and their environment.
package test_structure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
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

	tmpRootFolder, err := files.CopyTerraformFolderToTemp(rootFolder, cleanName(t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	tmpTestFolder := filepath.Join(tmpRootFolder, terraformModuleFolder)

	// Log temp folder so we can see it
	logger.Logf(t, "Copied terraform folder %s to %s", filepath.Join(rootFolder, terraformModuleFolder), tmpTestFolder)

	return tmpTestFolder
}

func cleanName(originalName string) string {
	parts := strings.Split(originalName, "/")
	return parts[len(parts)-1]
}
