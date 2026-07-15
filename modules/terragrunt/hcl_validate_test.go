package terragrunt_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
	"github.com/stretchr/testify/require"
)

func TestHclValidate(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	terragrunt.HclValidateContext(t, context.Background(), &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	})
}

func TestHclValidateE_InvalidConfig(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "terragrunt.hcl"), []byte("not_valid!!!"), 0644))

	_, err := terragrunt.HclValidateContextE(t, context.Background(), &terragrunt.Options{TerragruntDir: tmpDir})
	require.Error(t, err)
}
