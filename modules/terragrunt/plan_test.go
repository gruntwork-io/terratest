package terragrunt_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanAllExitCode(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.DestroyAllContext(t, context.Background(), options)

	terragrunt.ApplyAllContext(t, context.Background(), options)
	exitCode := terragrunt.PlanAllExitCodeContext(t, context.Background(), options)
	require.Equal(t, 0, exitCode)
}

func TestPlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out := terragrunt.PlanContext(t, context.Background(), options)
	require.NotEmpty(t, out)
}

func TestPlanExitCode(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Apply first so plan shows no changes (exit code 0)
	terragrunt.ApplyContext(t, context.Background(), options)
	defer terragrunt.DestroyContext(t, context.Background(), options)

	exitCode := terragrunt.PlanExitCodeContext(t, context.Background(), options)
	assert.Equal(t, 0, exitCode)
}

func TestInitAndPlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out := terragrunt.InitAndPlanContext(t, context.Background(), options)
	require.NotEmpty(t, out)
}

func TestPlanAllWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	getExitCode, errExitCode := terragrunt.PlanAllExitCodeContextE(t, context.Background(), options)
	// GetExitCodeForRunCommandError was unable to determine the exit code correctly
	require.NoError(t, errExitCode)

	require.Equal(t, 1, getExitCode)
}

func TestAssertPlanAllExitCodeNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.DestroyAllContext(t, context.Background(), options)

	getExitCode, errExitCode := terragrunt.PlanAllExitCodeContextE(t, context.Background(), options)
	if errExitCode != nil {
		t.Fatal(errExitCode)
	}

	// since there is no state file we expect `2` to be the success exit code
	assert.Equal(t, 2, getExitCode)
	assertPlanAllExitCode(t, getExitCode, true)

	terragrunt.ApplyAllContext(t, context.Background(), options)

	getExitCode, errExitCode = terragrunt.PlanAllExitCodeContextE(t, context.Background(), options)
	if errExitCode != nil {
		t.Fatal(errExitCode)
	}

	// since there is a state file we expect `0` to be the success exit code
	assert.Equal(t, 0, getExitCode)
	assertPlanAllExitCode(t, getExitCode, true)
}

func TestAssertPlanAllExitCodeWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	getExitCode, errExitCode := terragrunt.PlanAllExitCodeContextE(t, context.Background(), options)
	require.NoError(t, errExitCode)

	assertPlanAllExitCode(t, getExitCode, false)
}

func assertPlanAllExitCode(t *testing.T, exitCode int, assertTrue bool) {
	t.Helper()

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
