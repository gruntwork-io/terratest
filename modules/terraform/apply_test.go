package terraform

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
)

func TestApplyNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-no-error", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &Options{
		TerraformDir: testFolder,
		NoColor:      true,
	}

	out := InitAndApply(t, options)

	assert.Contains(t, out, "Hello, World")

	// Check that NoColor correctly doesn't output the colour escape codes which look like [0m,[1m or [32m
	assert.NotRegexp(t, `\[\d*m`, out, "Output should not contain color escape codes")
}

func TestApplyWithErrorNoRetry(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-error", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &Options{
		TerraformDir: testFolder,
	}

	out, err := InitAndApplyE(t, options)

	assert.Error(t, err)
	assert.Contains(t, out, "This is the first run, exiting with an error")
}

func TestApplyWithErrorWithRetry(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-with-error", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &Options{
		TerraformDir: testFolder,
		MaxRetries:   1,
		RetryableTerraformErrors: map[string]string{
			"This is the first run, exiting with an error": "Intentional failure in test fixture",
		},
	}

	out := InitAndApply(t, options)

	assert.Contains(t, out, "This is the first run, exiting with an error")
}
