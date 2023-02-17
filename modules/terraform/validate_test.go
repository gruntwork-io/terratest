package terraform

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestInitAndValidateWithNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	out := InitAndValidate(t, options)
	require.Contains(t, out, "The configuration is valid")
}

func TestInitAndValidateWithError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-plan-error", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
	}

	out, err := InitAndValidateE(t, options)
	require.Error(t, err)
	require.Contains(t, out, "Reference to undeclared input variable")
}
