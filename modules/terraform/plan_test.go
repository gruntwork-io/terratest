package terraform

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitAndPlanWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	_, err = InitAndPlanE(t, options)
	require.Error(t, err)
}

func TestInitAndPlanWithNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	// In Terraform 0.12 and below, if there were no resources to create, update, or destroy, 'plan' command would
	// report "No changes. Infrastructure is up-to-date." However, with 0.13 and above, if the Terraform configuration
	// has never been applied at all, 'plan' always shows changes. So we have to run 'apply' first, and can then
	// check that 'plan' returns the message we expect.
	InitAndApply(t, options)
	out, err := PlanE(t, options)
	require.NoError(t, err)
	require.Contains(t, out, "No changes.")
}

func TestInitAndPlanWithOutput(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
	}

	out, err := InitAndPlanE(t, options)
	require.NoError(t, err)
	require.Contains(t, out, "1 to add, 0 to change, 0 to destroy.")
}

func TestInitAndPlanWithPlanFile(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)
	planFilePath := filepath.Join(testFolder, "plan.out")

	options := &Options{
		TerraformDir: testFolder,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
		PlanFilePath: planFilePath,
	}

	out, err := InitAndPlanE(t, options)
	require.NoError(t, err)
	assert.Contains(t, out, "1 to add, 0 to change, 0 to destroy.")
	assert.Contains(t, out, fmt.Sprintf("Saved the plan to: %s", planFilePath))
	assert.FileExists(t, planFilePath, "Plan file was not saved to expected location:", planFilePath)
}

func TestPlanWithExitCodeWithNoChanges(t *testing.T) {
	t.Parallel()
	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	// In Terraform 0.12 and below, if there were no resources to create, update, or destroy, the -detailed-exitcode
	// would return a code of 0. However, with 0.13 and above, if the Terraform configuration has never been applied
	// at all, -detailed-exitcode always returns an exit code of 2. So we have to run 'apply' first, and can then
	// check that 'plan' returns the exit code we expect.
	InitAndApply(t, options)
	exitCode := PlanExitCode(t, options)
	require.Equal(t, DefaultSuccessExitCode, exitCode)
}

func TestPlanWithExitCodeWithChanges(t *testing.T) {
	t.Parallel()
	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
	}
	exitCode := InitAndPlanWithExitCode(t, options)
	require.Equal(t, TerraformPlanChangesPresentExitCode, exitCode)
}

func TestPlanWithExitCodeWithFailure(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	exitCode, getExitCodeErr := InitAndPlanWithExitCodeE(t, options)
	require.NoError(t, getExitCodeErr)
	require.Equal(t, exitCode, 1)
}

func TestTgPlanAllNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir:    testFolder,
		TerraformBinary: "terragrunt",
	}

	// In Terraform 0.12 and below, if there were no resources to create, update, or destroy, the -detailed-exitcode
	// would return a code of 0. However, with 0.13 and above, if the Terraform configuration has never been applied
	// at all, -detailed-exitcode always returns an exit code of 2. So we have to run 'apply' first, and can then
	// check that 'plan' returns the exit code we expect.
	TgApplyAll(t, options)
	getExitCode, errExitCode := TgPlanAllExitCodeE(t, options)
	// GetExitCodeForRunCommandError was unable to determine the exit code correctly
	if errExitCode != nil {
		t.Fatal(errExitCode)
	}

	// Since PlanAllExitCodeTgE returns error codes, we want to compare against 1
	require.Equal(t, DefaultSuccessExitCode, getExitCode)
}

func TestTgPlanAllWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir:    testFolder,
		TerraformBinary: "terragrunt",
	}

	getExitCode, errExitCode := TgPlanAllExitCodeE(t, options)
	// GetExitCodeForRunCommandError was unable to determine the exit code correctly
	require.NoError(t, errExitCode)

	require.Equal(t, DefaultErrorExitCode, getExitCode)
}
