package test_structure

import "testing"

func TestCopyToTempFolder(t *testing.T) {
	tempFolder := CopyTerraformFolderToTemp(t, "../../", "examples")
	t.Log(tempFolder)
}

func TestCopySubtestToTempFolder(t *testing.T) {
	t.Run("Subtest", func(t *testing.T) {
		tempFolder := CopyTerraformFolderToTemp(t, "../../", "examples")
		t.Log(tempFolder)
	})
}
