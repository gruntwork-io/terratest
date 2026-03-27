package terragrunt_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
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

	terragrunt.ApplyAll(t, options)
	destroyOut := terragrunt.DestroyAll(t, options)
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

	terragrunt.Apply(t, options)
	destroyOut := terragrunt.Destroy(t, options)
	require.NotEmpty(t, destroyOut)
}

// TestDestroyAllWithArgs verifies DestroyAll respects TerragruntArgs
func TestDestroyAllWithArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Apply first
	terragrunt.ApplyAll(t, &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	// Destroy with TerragruntArgs
	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
		TerragruntArgs:   []string{"--log-level", "error"},
	}

	destroyOut := terragrunt.DestroyAll(t, options)
	require.NotEmpty(t, destroyOut)
	// With --log-level error, should not see info logs
	require.NotContains(t, destroyOut, "level=info")
}
