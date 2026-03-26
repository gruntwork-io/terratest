package terragrunt_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	"github.com/stretchr/testify/require"
)

func TestValidateAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	terragrunt.ValidateAll(t, &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
}

func TestValidate(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	terragrunt.Validate(t, &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
}

func TestInitAndValidate(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	out := terragrunt.InitAndValidate(t, &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
	require.NotEmpty(t, out)
}
