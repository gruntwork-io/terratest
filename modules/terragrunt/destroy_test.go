package terragrunt_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
	"github.com/stretchr/testify/require"
)

func TestDestroyAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	terragrunt.ApplyAllContext(t, context.Background(), options)
	destroyOut := terragrunt.DestroyAllContext(t, context.Background(), options)
	require.NotEmpty(t, destroyOut)
}

func TestDestroy(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	terragrunt.ApplyContext(t, context.Background(), options)
	destroyOut := terragrunt.DestroyContext(t, context.Background(), options)
	require.NotEmpty(t, destroyOut)
}

// TestDestroyAllWithArgs verifies DestroyAll respects TerragruntArgs
func TestDestroyAllWithArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Apply first
	terragrunt.ApplyAllContext(t, context.Background(), &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	// Destroy with TerragruntArgs
	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--log-level", "error"},
	}

	destroyOut := terragrunt.DestroyAllContext(t, context.Background(), options)
	require.NotEmpty(t, destroyOut)
	// With --log-level error, should not see info logs
	require.NotContains(t, destroyOut, "level=info")
}
