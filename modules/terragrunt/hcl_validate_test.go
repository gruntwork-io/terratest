package terragrunt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestHclValidate(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	HclValidate(t, &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
}

func TestHclValidateE_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "terragrunt.hcl"), []byte("not_valid!!!"), 0644))

	_, err := HclValidateE(t, &Options{TerragruntDir: tmpDir})
	require.Error(t, err)
}
