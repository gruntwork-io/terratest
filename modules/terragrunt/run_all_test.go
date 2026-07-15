package terragrunt_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
	"github.com/stretchr/testify/require"
)

func TestRunAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Test with validate command
	out := terragrunt.RunContext(t, context.Background(), options, []string{"--all"}, []string{"validate"})
	require.NotEmpty(t, out)
}

func TestRunAllE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Test with validate command
	out, err := terragrunt.RunContextE(t, context.Background(), options, []string{"--all"}, []string{"validate"})
	require.NoError(t, err)
	require.NotEmpty(t, out)
}

func TestRunAllWithPlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Test with plan command - verify output contains expected terraform plan text
	out, err := terragrunt.RunContextE(t, context.Background(), options, []string{"--all"}, []string{"plan"})
	require.NoError(t, err)
	require.Contains(t, out, "Changes to Outputs")
}
