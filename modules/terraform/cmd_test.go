package terraform

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformCommand(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-error", strings.ReplaceAll(t.Name(), "/", "-"))
		require.NoError(t, err)
		options := &Options{
			TerraformDir: testFolder,
		}
		Init(t, options)

		stdout, stderr, code, err := RunTerraformCommandAndGetStdOutErrCodeE(t, options, "apply", "-input=false", "-auto-approve")
		assert.Error(t, err)
		assert.Contains(t, stdout, "Creating...", "should capture stdout")
		assert.Contains(t, stderr, "Error: ", "should capture stderr")
		assert.Greater(t, code, 0)
	})

	t.Run("WithWarning", func(t *testing.T) {
		testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-warning", strings.ReplaceAll(t.Name(), "/", "-"))
		require.NoError(t, err)
		options := &Options{
			TerraformDir: testFolder,
			WarningsAsErrors: map[string]string{
				".*lorem ipsum.*": "this warning message should shown.",
			},
		}
		Init(t, options)

		stdout, stderr, code, err := RunTerraformCommandAndGetStdOutErrCodeE(t, options, "apply", "-input=false", "-auto-approve")
		assert.Error(t, err)
		assert.Contains(t, stdout, "Creating...", "should capture stdout")
		assert.Contains(t, stderr, "", "should capture stderr")
		assert.Greater(t, code, 0)
	})

	t.Run("NoError", func(t *testing.T) {
		testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", strings.ReplaceAll(t.Name(), "/", "-"))
		require.NoError(t, err)
		options := &Options{
			TerraformDir: testFolder,
		}

		{
			stdout, stderr, code := RunTerraformCommandAndGetStdOutErrCode(t, options, "apply", "-input=false", "-auto-approve")
			assert.Contains(t, stdout, `test = "Hello, World"`, "should capture stdout")
			assert.Equal(t, code, 0)
			assert.Empty(t, stderr)
		}

		{
			stdout := RunTerraformCommandAndGetStdout(t, options, "apply", "-input=false", "-auto-approve")
			assert.Contains(t, stdout, `test = "Hello, World"`, "should capture stdout")
		}
	})
}
