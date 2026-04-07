package terragrunt_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputJSON(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-output", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	terragrunt.Apply(t, options)
	defer terragrunt.Destroy(t, options)

	json := terragrunt.OutputJSON(t, options, "str")
	assert.Contains(t, json, "str")
}

func TestOutputJSONAllKeys(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-output", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	terragrunt.Apply(t, options)
	defer terragrunt.Destroy(t, options)

	json := terragrunt.OutputJSON(t, options, "")
	assert.Contains(t, json, "str")
	assert.Contains(t, json, "list")
	assert.Contains(t, json, "map")
}

func TestOutputJSONE_Error(t *testing.T) {
	t.Parallel()

	options := &terragrunt.Options{
		TerragruntDir:    t.TempDir(),
		TerragruntBinary: "terragrunt",
	}

	_, err := terragrunt.OutputJSONE(t, options, "")
	require.Error(t, err)
}
