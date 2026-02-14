package terragrunt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	output := Graph(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})

	require.Contains(t, output, "digraph")
	require.Contains(t, output, `"foo"`)
	require.Contains(t, output, `"bar"`)
}

func TestGraphE_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "terragrunt.hcl"), []byte("not_valid!!!"), 0644))

	_, err := GraphE(t, &Options{TerragruntDir: tmpDir})
	require.Error(t, err)
}
