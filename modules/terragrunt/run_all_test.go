package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestRunAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Test with validate command
	out := RunAll(t, options, "validate")
	require.NotEmpty(t, out)
}

func TestRunAllE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Test with validate command
	out, err := RunAllE(t, options, "validate")
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func TestRunAllWithPlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Test with plan command - verify output contains expected terraform plan text
	out, err := RunAllE(t, options, "plan")
	require.NoError(t, err)
	require.Contains(t, out, "Changes to Outputs")
}
