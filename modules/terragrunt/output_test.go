package terragrunt

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestOutputAllJson(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp(
		"../../test/fixtures/terragrunt/terragrunt-multi-plan", t.Name())
	require.NoError(t, err)

	options := &Options{TerragruntDir: testFolder}

	ApplyAll(t, options)
	defer DestroyAll(t, options)

	output := OutputAllJson(t, options)

	// Current terragrunt returns separate JSON objects per module
	require.Contains(t, output, `"value": "foo"`)
	require.Contains(t, output, `"value": "bar"`)
}
