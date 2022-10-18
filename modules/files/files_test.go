package files

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const copyFolderContentsFixtureRoot = "../../test/fixtures/copy-folder-contents"

func TestFileExists(t *testing.T) {
	t.Parallel()

	currentFile, err := filepath.Abs(os.Args[0])
	require.NoError(t, err)

	assert.True(t, FileExists(currentFile))
	assert.False(t, FileExists("/not/a/real/path"))
}

func TestIsExistingFile(t *testing.T) {
	t.Parallel()

	currentFile, err := filepath.Abs(os.Args[0])
	require.NoError(t, err)
	currentFileDir := filepath.Dir(currentFile)

	assert.True(t, IsExistingFile(currentFile))
	assert.False(t, IsExistingFile("/not/a/real/path"))
	assert.False(t, IsExistingFile(currentFileDir))
}

func TestIsExistingDir(t *testing.T) {
	t.Parallel()

	currentFile, err := filepath.Abs(os.Args[0])
	require.NoError(t, err)
	currentFileDir := filepath.Dir(currentFile)

	assert.False(t, IsExistingDir(currentFile))
	assert.False(t, IsExistingDir("/not/a/real/path"))
	assert.True(t, IsExistingDir(currentFileDir))
}

func TestCopyFolderToDest(t *testing.T) {
	t.Parallel()

	tempFolderPrefix := "someprefix"
	destFolder := os.TempDir()
	tmpDir := t.TempDir()

	filter := func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path) && !PathContainsTerraformState(path)
	}

	folder, err := CopyFolderToDest("/not/a/real/path", destFolder, tempFolderPrefix, filter)
	require.Error(t, err)
	assert.False(t, FileExists(folder))

	folder, err = CopyFolderToDest(tmpDir, destFolder, tempFolderPrefix, filter)
	assert.DirExists(t, folder)
	assert.NoError(t, err)
}

func TestCopyFolderContents(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "full-copy")
	tmpDir := t.TempDir()

	err := CopyFolderContents(originalDir, tmpDir)
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyFolderContentsWithHiddenFilesFilter(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-hidden-files")
	tmpDir := t.TempDir()

	err := CopyFolderContentsWithFilter(originalDir, tmpDir, func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path)
	})
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

// Test copying a folder that contains symlinks
func TestCopyFolderContentsWithSymLinks(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "symlinks")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "symlinks")
	tmpDir := t.TempDir()

	err := CopyFolderContentsWithFilter(originalDir, tmpDir, func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path)
	})
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

// Test copying a folder that contains symlinks that point to a non-existent file
func TestCopyFolderContentsWithBrokenSymLinks(t *testing.T) {
	t.Parallel()

	// Creating broken symlink
	pathToFile := filepath.Join(copyFolderContentsFixtureRoot, "symlinks-broken/nonexistent-folder/bar.txt")
	pathToSymlink := filepath.Join(copyFolderContentsFixtureRoot, "symlinks-broken/bar.txt")
	defer func() {
		if err := os.Remove(pathToSymlink); err != nil {
			t.Fatal(fmt.Errorf("Failed to remove link: %+v", err))
		}
	}()
	if err := os.Symlink(pathToFile, pathToSymlink); err != nil {
		t.Fatal(fmt.Errorf("Failed to create broken link for test: %+v", err))
	}

	// Test copying folder
	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "symlinks-broken")
	tmpDir := t.TempDir()

	err := CopyFolderContentsWithFilter(originalDir, tmpDir, func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path)
	})
	require.NoError(t, err)

	// This requireDirectoriesEqual command uses GNU diff under the hood, but unfortunately we cannot instruct diff to
	// compare symlinks in two directories without attempting to dereference any symlinks until diff version 3.3.0.
	// Because many environments are still using diff < 3.3.0, we disregard this test for now.
	// Per https://unix.stackexchange.com/a/119406/129208
	//requireDirectoriesEqual(t, expectedDir, tmpDir)
	fmt.Println("Test completed without error, however due to a limitation in GNU diff < 3.3.0, directories have not been compared for equivalency.")
}

func TestCopyTerraformFolderToTemp(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-hidden-files-no-terraform-files")

	tmpDir, err := CopyTerraformFolderToTemp(originalDir, "TestCopyTerraformFolderToTemp")
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyTerraformFolderToDest(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-hidden-files-no-terraform-files")
	destFolder := os.TempDir()

	tmpDir, err := CopyTerraformFolderToDest(originalDir, destFolder, "TestCopyTerraformFolderToTemp")
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyTerragruntFolderToTemp(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "terragrunt-files")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-state-files")

	tmpDir, err := CopyTerragruntFolderToTemp(originalDir, t.Name())
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyTerragruntFolderToDest(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "terragrunt-files")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-state-files")
	destFolder := os.TempDir()

	tmpDir, err := CopyTerragruntFolderToDest(originalDir, destFolder, t.Name())
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestPathContainsTerraformStateOrVars(t *testing.T) {
	var data = []struct {
		desc     string
		path     string
		contains bool
	}{
		{"contains tfvars", "./folder/terraform.tfvars", true},
		{"contains tfvars.json", "./folder/hello/terraform.tfvars.json", true},
		{"contains state", "./folder/hello/helloagain/terraform.tfstate", true},
		{"contains state backup", "./folder/hey/terraform.tfstate.backup", true},
		{"does not contain any", "./folder/salut/terraform.json", false},
	}

	for _, tt := range data {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			result := PathContainsTerraformStateOrVars(tt.path)
			if result != tt.contains {
				if tt.contains {
					t.Errorf("Expected %s to contain Terraform related file", tt.path)
				} else {
					t.Errorf("Expected %s to not contain Terraform related file", tt.path)
				}
			}
		})
	}
}

// Diffing two directories to ensure they have the exact same files, contents, etc and showing exactly what's different
// takes a lot of code. Why waste time on that when this functionality is already nicely implemented in the Unix/Linux
// "diff" command? We shell out to that command at test time.
func requireDirectoriesEqual(t *testing.T, folderWithExpectedContents string, folderWithActualContents string) {
	cmd := exec.Command("diff", "-r", "-u", folderWithExpectedContents, folderWithActualContents)

	bytes, err := cmd.Output()
	output := string(bytes)

	require.NoError(t, err, "diff command exited with an error. This likely means the contents of %s and %s are different. Here is the output of the diff command:\n%s", folderWithExpectedContents, folderWithActualContents, output)
}
