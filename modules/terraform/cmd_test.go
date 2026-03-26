package terraform_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformCommand(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-error", strings.ReplaceAll(t.Name(), "/", "-"))
		require.NoError(t, err)

		options := &terraform.Options{
			TerraformDir: testFolder,
		}
		terraform.Init(t, options)

		stdout, stderr, code, err := terraform.RunTerraformCommandAndGetStdOutErrCodeE(t, options, "apply", "-input=false", "-auto-approve")
		require.Error(t, err)
		assert.Contains(t, stdout, "Creating...", "should capture stdout")
		assert.Contains(t, stderr, "Error: ", "should capture stderr")
		assert.Positive(t, code)
	})

	t.Run("WithWarning", func(t *testing.T) {
		t.Parallel()

		testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-warning", strings.ReplaceAll(t.Name(), "/", "-"))
		require.NoError(t, err)

		options := &terraform.Options{
			TerraformDir: testFolder,
			WarningsAsErrors: map[string]string{
				".*lorem ipsum.*": "this warning message should shown.",
			},
		}
		terraform.Init(t, options)

		stdout, stderr, code, err := terraform.RunTerraformCommandAndGetStdOutErrCodeE(t, options, "apply", "-input=false", "-auto-approve")
		require.Error(t, err)
		assert.Contains(t, stdout, "Creating...", "should capture stdout")
		assert.Contains(t, stderr, "", "should capture stderr")
		assert.Positive(t, code)
	})

	t.Run("NoError", func(t *testing.T) {
		t.Parallel()

		testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", strings.ReplaceAll(t.Name(), "/", "-"))
		require.NoError(t, err)

		options := &terraform.Options{
			TerraformDir: testFolder,
		}

		{
			stdout, stderr, code := terraform.RunTerraformCommandAndGetStdOutErrCode(t, options, "apply", "-input=false", "-auto-approve")
			assert.Contains(t, stdout, `test = "Hello, World"`, "should capture stdout")
			assert.Equal(t, 0, code)
			assert.Empty(t, stderr)
		}

		{
			stdout := terraform.RunTerraformCommandAndGetStdout(t, options, "apply", "-input=false", "-auto-approve")
			assert.Contains(t, stdout, `test = "Hello, World"`, "should capture stdout")
		}
	})
}
