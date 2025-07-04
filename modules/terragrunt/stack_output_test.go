package terragrunt

import (
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestTgStackOutput(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-simple", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with no-color option
	_, err = TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		NoColor:          true,
	})
	require.NoError(t, err)

	// Get the output of stack output command (plain text)
	out, err := TgStackOutputE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	}, "")
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")
}

func TestTgStackOutputE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-simple", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with no-color option
	_, err = TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		NoColor:          true,
	})
	require.NoError(t, err)

	// Get the output of stack output command
	out, err := TgStackOutputE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	}, "")
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")
}

func TestTgStackOutputJson(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terragrunt/terragrunt-stack-simple", t.Name())
	require.NoError(t, err)

	// First initialize the stack
	_, err = TgStackInitE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	})
	require.NoError(t, err)

	// Generate with no-color option
	_, err = TgStackGenerateE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
		NoColor:          true,
	})
	require.NoError(t, err)

	// Get the output of stack output command with -json
	out, err := TgStackOutputE(t, &Options{
		TerragruntDir:    path.Join(testFolder, "live"),
		TerragruntBinary: "terragrunt",
	}, "", "-json")
	require.NoError(t, err)
	require.Contains(t, out, ".terragrunt-stack")
	require.Contains(t, out, "has been successfully initialized!")
}
