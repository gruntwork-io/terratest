package test_structure

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/logger"
)

const SKIP_STAGE_ENV_VAR_PREFIX = "SKIP_"

// Execute the given test stage (e.g., setup, teardown, validation) if an environment variable of the name
// `SKIP_<stageName>` (e.g., SKIP_teardown) is not set.
func RunTestStage(t *testing.T, stageName string, stage func()) {
	envVarName := fmt.Sprintf("%s%s", SKIP_STAGE_ENV_VAR_PREFIX, stageName)
	if os.Getenv(envVarName) == "" {
		logger.Logf(t, "The '%s' environment variable is not set, so executing stage '%s'.", envVarName, stageName)
		stage()
	} else {
		logger.Logf(t, "The '%s' environment variable is set, so skipping stage '%s'.", envVarName, stageName)
	}
}

// Returns true if an environment variable is set instructing Terratest to skip a test stage. This can be an easy way
// to tell if the tests are running in a local dev environment vs a CI server.
func SkipStageEnvVarSet() bool {
	for _, environmentVariable := range os.Environ() {
		if strings.HasPrefix(environmentVariable, SKIP_STAGE_ENV_VAR_PREFIX) {
			return true
		}
	}

	return false
}

// Copy the given root folder to a randomly-named temp folder and return the path to the given examples folder within
// the new temp root folder. This is useful when running multiple tests in parallel against the same set of Terraform
// files to ensure the tests don't overwrite each other's .terraform working directory and terraform.tfstate files. To
// ensure relative paths work, we copy over the entire root folder to a temp folder, and then return the path within
// that temp folder to the given example dir, which is where the actual test will be running.
//
// Note that if any of the SKIP_<stage> environment variables is set, we assume this is a test in the local dev where
// there are no other concurrent tests running and we want to be able to cache test data between test stages, so in
// that case, we do NOT copy anything to a temp folder, and return the path to the original examples folder instead.
func CopyTerraformFolderToTemp(t *testing.T, rootFolder string, examplesFolder string) string {
	if SkipStageEnvVarSet() {
		logger.Logf(t, "A SKIP_XXX environment variable is set. Using original examples folder rather than a temp folder so we can cache data between stages for faster local testing.")
		return filepath.Join(rootFolder, examplesFolder)
	}

	tmpRootFolder, err := files.CopyTerraformFolderToTemp(rootFolder, cleanName(t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Join(tmpRootFolder, examplesFolder)
}

func cleanName(originalName string) string {
	parts := strings.Split(originalName, "/")
	return parts[len(parts)-1]
}
