package terragrunt_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
	"github.com/stretchr/testify/require"
)

func TestApplyAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.DestroyAllContext(t, context.Background(), options)

	out := terragrunt.ApplyAllContext(t, context.Background(), options)
	require.Contains(t, out, "Hello, World")
}

func TestApply(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.DestroyContext(t, context.Background(), options)

	out := terragrunt.ApplyContext(t, context.Background(), options)
	require.Contains(t, out, "Hello, World")
}

func TestInitAndApply(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.DestroyContext(t, context.Background(), options)

	out := terragrunt.InitAndApplyContext(t, context.Background(), options)
	require.Contains(t, out, "Hello, World")
}

// TestInitAndApplyE_InitFailure verifies that when init fails, apply is skipped
// and the init error is propagated.
func TestInitAndApplyE_InitFailure(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out, err := terragrunt.InitAndApplyContextE(t, context.Background(), options)
	require.Error(t, err, "InitAndApplyE should propagate init failure")
	require.Empty(t, out, "Output should be empty when init fails")
	require.Contains(t, err.Error(), "Missing expression",
		"Error should be from init, not apply")
}
