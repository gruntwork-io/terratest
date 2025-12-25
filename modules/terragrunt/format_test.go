package terragrunt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files/v2"
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

	// Run format command
	FormatAll(t, options)

	// Read the formatted file to verify it was actually formatted
	formattedContent, err := os.ReadFile(tgFile)
	require.NoError(t, err)

	// Verify the file was formatted (should have proper spacing now)
	require.Contains(t, string(formattedContent), `source = "git::git@github.com:foo/modules.git//app"`,
		"Expected file to be formatted with proper spacing")
	require.Contains(t, string(formattedContent), `inputs = {`,
		"Expected inputs block to be formatted with spaces around =")
	require.Contains(t, string(formattedContent), `foo = "bar"`,
		"Expected key-value pairs to be formatted with spaces around =")
}

func TestFormatAllE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	// Create an unformatted file to ensure the command actually does something
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

	// Run format command - should succeed
	_, err = FormatAllE(t, options)
	require.NoError(t, err)

	// Verify the file was actually formatted by reading it
	formattedContent, err := os.ReadFile(tgFile)
	require.NoError(t, err)
	require.Contains(t, string(formattedContent), `inputs = {`,
		"File should be formatted with proper spacing")
}
