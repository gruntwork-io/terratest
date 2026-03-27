package terraform_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitAndPlanWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	_, err = terraform.InitAndPlanE(t, options)
	require.Error(t, err)
}

func TestInitAndPlanWithNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	// In Terraform 0.12 and below, if there were no resources to create, update, or destroy, 'plan' command would
	// report "No changes. Infrastructure is up-to-date." However, with 0.13 and above, if the Terraform configuration
	// has never been applied at all, 'plan' always shows changes. So we have to run 'apply' first, and can then
	// check that 'plan' returns the message we expect.
	terraform.InitAndApply(t, options)

	out, err := terraform.PlanE(t, options)
	require.NoError(t, err)
	require.Contains(t, out, "No changes.")
}

func TestInitAndPlanWithOutput(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
		Vars: map[string]any{
			"cnt": 1,
		},
	}

	out, err := terraform.InitAndPlanE(t, options)
	require.NoError(t, err)
	require.Contains(t, out, "1 to add, 0 to change, 0 to destroy.")
}

func TestInitAndPlanWithPlanFile(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	planFilePath := filepath.Join(testFolder, "plan.out")

	options := &terraform.Options{
		TerraformDir: testFolder,
		Vars: map[string]any{
			"cnt": 1,
		},
		PlanFilePath: planFilePath,
	}

	out, err := terraform.InitAndPlanE(t, options)
	require.NoError(t, err)

	// clean output to be consistent in checks
	out = strings.ReplaceAll(out, "\n", "")
	assert.Contains(t, out, "1 to add, 0 to change, 0 to destroy.")
	assert.Contains(t, out, "Saved the plan to:"+planFilePath)
	assert.FileExists(t, planFilePath, "Plan file was not saved to expected location:", planFilePath)
}

func TestInitAndPlanAndShowWithStructNoLogTempPlanFile(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
		Vars: map[string]any{
			"cnt": 1,
		},
	}

	planStruct := terraform.InitAndPlanAndShowWithStructNoLogTempPlanFile(t, options)
	assert.Len(t, planStruct.ResourceChangesMap, 1)
}

func TestPlanWithExitCodeWithNoChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	// In Terraform 0.12 and below, if there were no resources to create, update, or destroy, the -detailed-exitcode
	// would return a code of 0. However, with 0.13 and above, if the Terraform configuration has never been applied
	// at all, -detailed-exitcode always returns an exit code of 2. So we have to run 'apply' first, and can then
	// check that 'plan' returns the exit code we expect.
	terraform.InitAndApply(t, options)

	exitCode := terraform.PlanExitCode(t, options)
	require.Equal(t, terraform.DefaultSuccessExitCode, exitCode)
}

func TestPlanWithExitCodeWithChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
		Vars: map[string]any{
			"cnt": 1,
		},
	}

	exitCode := terraform.InitAndPlanWithExitCode(t, options)
	require.Equal(t, terraform.TerraformPlanChangesPresentExitCode, exitCode)
}

func TestPlanWithExitCodeWithFailure(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	exitCode, getExitCodeErr := terraform.InitAndPlanWithExitCodeE(t, options)
	require.NoError(t, getExitCodeErr)
	require.Equal(t, 1, exitCode)
}
