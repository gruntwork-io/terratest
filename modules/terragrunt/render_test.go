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

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
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

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	output := RenderJson(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(output), &parsed), "output should be valid JSON")
	require.Contains(t, parsed, "terraform")
}

func TestRenderE_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "terragrunt.hcl"), []byte("not_valid!!!"), 0644))

	_, err := RenderE(t, &Options{TerragruntDir: tmpDir})
	require.Error(t, err)
}
