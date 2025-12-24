package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files/v2"
	"github.com/stretchr/testify/require"
)

func TestDestroyAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	ApplyAll(t, options)
	destroyOut := DestroyAll(t, options)
	require.NotEmpty(t, destroyOut)
}

func TestDestroyAllE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out, err := ApplyAllE(t, options)
	require.NoError(t, err)
	require.Contains(t, out, "Hello, World")

	// Test that destroy completes successfully
	destroyOut, err := DestroyAllE(t, options)
	require.NoError(t, err)
	require.NotEmpty(t, destroyOut, "Destroy output should not be empty")
}

// TestDestroyAllWithArgs verifies DestroyAll respects TerragruntArgs
func TestDestroyAllWithArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Apply first
	ApplyAll(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	// Destroy with TerragruntArgs
	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--log-level", "error"},
	}

	destroyOut := DestroyAll(t, options)
	require.NotEmpty(t, destroyOut)
	// With --log-level error, should not see info logs
	require.NotContains(t, destroyOut, "level=info")
}
