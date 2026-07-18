package terragrunt_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt/v2"
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

	terragrunt.ApplyContext(t, context.Background(), options)
	defer terragrunt.DestroyContext(t, context.Background(), options)

	json := terragrunt.OutputJSONContext(t, context.Background(), options, "str")
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

	terragrunt.ApplyContext(t, context.Background(), options)
	defer terragrunt.DestroyContext(t, context.Background(), options)

	json := terragrunt.OutputJSONContext(t, context.Background(), options, "")
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

	_, err := terragrunt.OutputJSONContextE(t, context.Background(), options, "")
	require.Error(t, err)
}
