package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanAllExitCode(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	ApplyAll(t, options)
	exitCode := PlanAllExitCode(t, options)
	require.Equal(t, 0, exitCode)
}

func TestPlanAllExitCodeE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// In Terraform 0.12 and below, if there were no resources to create, update, or destroy, the -detailed-exitcode
	// would return a code of 0. However, with 0.13 and above, if the Terraform configuration has never been applied
	// at all, -detailed-exitcode always returns an exit code of 2. So we have to run 'apply' first, and can then
	// check that 'plan' returns the exit code we expect.
	ApplyAll(t, options)
	getExitCode, errExitCode := PlanAllExitCodeE(t, options)
	// GetExitCodeForRunCommandError was unable to determine the exit code correctly
	if errExitCode != nil {
		t.Fatal(errExitCode)
	}

	// Since PlanAllExitCodeE returns error codes, we want to compare against 1
	require.Equal(t, 0, getExitCode)
}

func TestPlanAllWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	getExitCode, errExitCode := PlanAllExitCodeE(t, options)
	// GetExitCodeForRunCommandError was unable to determine the exit code correctly
	require.NoError(t, errExitCode)

	require.Equal(t, 1, getExitCode)
}

func TestAssertPlanAllExitCodeNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	getExitCode, errExitCode := PlanAllExitCodeE(t, options)
	if errExitCode != nil {
		t.Fatal(errExitCode)
	}

	// since there is no state file we expect `2` to be the success exit code
	assert.Equal(t, 2, getExitCode)
	assertPlanAllExitCode(t, getExitCode, true)

	ApplyAll(t, options)

	getExitCode, errExitCode = PlanAllExitCodeE(t, options)
	if errExitCode != nil {
		t.Fatal(errExitCode)
	}

	// since there is a state file we expect `0` to be the success exit code
	assert.Equal(t, 0, getExitCode)
	assertPlanAllExitCode(t, getExitCode, true)
}

func TestAssertPlanAllExitCodeWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	getExitCode, errExitCode := PlanAllExitCodeE(t, options)
	require.NoError(t, errExitCode)

	assertPlanAllExitCode(t, getExitCode, false)
}

func assertPlanAllExitCode(t *testing.T, exitCode int, assertTrue bool) {

	validExitCodes := map[int]bool{
		0: true,
		2: true,
	}

	_, hasKey := validExitCodes[exitCode]
	if assertTrue {
		assert.True(t, hasKey)
	} else {
		assert.False(t, hasKey)
	}
}
