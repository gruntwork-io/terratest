package packer_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/packer/v2"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndLoadPackerOptions(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &packer.Options{
		Template:   "packer.json",
		Only:       "amazon-ebs",
		WorkingDir: "/tmp/packer",
		Vars: map[string]string{
			"aws_region": "us-east-1",
		},
	}
	packer.SavePackerOptions(t, tmpFolder, expectedData)

	actualData := packer.LoadPackerOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}
