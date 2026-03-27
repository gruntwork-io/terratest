package terragrunt_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	output := terragrunt.Graph(t, &terragrunt.Options{
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

	output, err := terragrunt.GraphE(t, &terragrunt.Options{TerragruntDir: tmpDir})
	require.NoError(t, err)
	require.Contains(t, output, "digraph")
	// Invalid config produces a minimal graph with just the current unit
	require.NotContains(t, output, "->")
}
