package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestApplyAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer DestroyAll(t, options)
	out := ApplyAll(t, options)
	require.Contains(t, out, "Hello, World")
}

func TestApply(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer Destroy(t, options)
	out := Apply(t, options)
	require.Contains(t, out, "Hello, World")
}

func TestInitAndApply(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer Destroy(t, options)
	out := InitAndApply(t, options)
	require.Contains(t, out, "Hello, World")
}

// TestInitAndApplyE_InitFailure verifies that when init fails, apply is skipped
// and the init error is propagated.
func TestInitAndApplyE_InitFailure(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(
		"testdata/terragrunt-stack-init-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out, err := InitAndApplyE(t, options)
	require.Error(t, err, "InitAndApplyE should propagate init failure")
	require.Empty(t, out, "Output should be empty when init fails")
	require.Contains(t, err.Error(), "Missing expression",
		"Error should be from init, not apply")
}
