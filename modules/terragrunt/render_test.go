package terragrunt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	output := Render(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	require.Contains(t, output, `source = "`)
	require.Contains(t, output, `extra_arguments`)
	// Verify log lines are stripped and indentation is preserved
	require.NotContains(t, output, "level=")
	require.Contains(t, output, "  source = ")
}

func TestRenderJson(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	output := RenderJson(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &parsed), "output should be valid JSON")
	require.Contains(t, parsed, "terraform")
}

func TestFilterLogLines(t *testing.T) {
	t.Parallel()

	input := "20:41:53.564 INFO   some log message\n  source = \"./modules/vpc\"\n\ntime=2023-07-11 level=info msg=hello\n  inputs = {\nGroup 1\n    name = \"test\"\n  }"
	result := filterLogLines(input)

	// Log lines and metadata lines should be stripped
	require.NotContains(t, result, "INFO")
	require.NotContains(t, result, "level=info")
	require.NotContains(t, result, "Group 1")

	// Indentation should be preserved (unlike removeLogLines which trims)
	require.Contains(t, result, "  source = ")
	require.Contains(t, result, "  inputs = {")
	require.Contains(t, result, "    name = ")
}

func TestRenderE_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "terragrunt.hcl"), []byte("not_valid!!!"), 0644))

	_, err := RenderE(t, &Options{TerragruntDir: tmpDir})
	require.Error(t, err)
}
