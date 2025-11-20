package terragrunt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestFormatAll(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Create an unformatted terragrunt.hcl file in foo directory
	unformattedContent := `terraform {
source = "git::git@github.com:foo/modules.git//app"
}
inputs={
foo="bar"
}`
	tgFile := filepath.Join(testFolder, "foo", "terragrunt.hcl")
	err = os.WriteFile(tgFile, []byte(unformattedContent), 0644)
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out := FormatAll(t, options)
	require.NotEmpty(t, out)
}

func TestFormatAllE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Create an unformatted terragrunt.hcl file in bar directory
	unformattedContent := `terraform {
source = "git::git@github.com:foo/modules.git//app"
}
inputs={
foo="bar"
}`
	tgFile := filepath.Join(testFolder, "bar", "terragrunt.hcl")
	err = os.WriteFile(tgFile, []byte(unformattedContent), 0644)
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	out, err := FormatAllE(t, options)
	require.NoError(t, err)
	require.NotEmpty(t, out)
}
