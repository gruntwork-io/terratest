package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestDestroyAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	ApplyAll(t, options)
	destroyOut := DestroyAll(t, options)
	require.NotEmpty(t, destroyOut)
}

func TestDestroy(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	Apply(t, options)
	destroyOut := Destroy(t, options)
	require.NotEmpty(t, destroyOut)
}

// TestDestroyAllWithArgs verifies DestroyAll respects TerragruntArgs
func TestDestroyAllWithArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
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
