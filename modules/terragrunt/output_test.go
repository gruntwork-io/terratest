package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputJson(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-output", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	Apply(t, options)
	defer Destroy(t, options)

	json := OutputJson(t, options, "str")
	assert.Contains(t, json, "str")
}

func TestOutputJsonAllKeys(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-output", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	Apply(t, options)
	defer Destroy(t, options)

	json := OutputJson(t, options, "")
	assert.Contains(t, json, "str")
	assert.Contains(t, json, "list")
	assert.Contains(t, json, "map")
}

func TestOutputJsonE_Error(t *testing.T) {
	t.Parallel()

	options := &Options{
		TerragruntDir:    t.TempDir(),
		TerragruntBinary: "terragrunt",
	}

	_, err := OutputJsonE(t, options, "")
	require.Error(t, err)
}
