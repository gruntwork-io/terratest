package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files/v2"
	"github.com/stretchr/testify/require"
)

func TestApplyAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	out := ApplyAll(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	require.Contains(t, out, "Hello, World")
}

func TestApplyAllE(t *testing.T) {
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
}
