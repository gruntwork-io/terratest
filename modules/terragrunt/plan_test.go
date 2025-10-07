package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestPlanAllExitCode(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Run plan with detailed exit code
	exitCode, err := PlanAllExitCodeE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Exit code should be 2 (changes present) or 0 (no changes)
	// depending on the state
	require.True(t, exitCode == 0 || exitCode == 2,
		"Expected exit code 0 or 2, got %d", exitCode)
}

func TestPlanAllExitCodeWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-with-error", t.Name())
	require.NoError(t, err)

	// Run plan - should get exit code 1 for terraform errors
	// Note: getExitCodeForTerraformCommandE returns the exit code without error
	// unless it cannot determine the exit code
	exitCode, err := PlanAllExitCodeE(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)       // No error determining exit code
	require.Equal(t, 1, exitCode) // Exit code 1 indicates terraform error
}

func TestPlanAllExitCodeDetectsChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-stack-init", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		TerraformArgs:    []string{"-upgrade=true"},
	})
	require.NoError(t, err)

	// Run plan - should detect changes (exit code 2)
	exitCode, err := PlanAllExitCodeE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
	require.Equal(t, 2, exitCode, "Expected exit code 2 (changes detected)")

	// Apply the changes
	_, err = ApplyAllE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Run plan again - should be no changes (exit code 0)
	exitCode, err = PlanAllExitCodeE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)
	require.Equal(t, 0, exitCode, "Expected exit code 0 (no changes)")
}
