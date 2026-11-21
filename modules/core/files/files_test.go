package files_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const copyFolderContentsFixtureRoot = "../../../test/fixtures/copy-folder-contents"

func TestFileExists(t *testing.T) {
	t.Parallel()

	currentFile, err := filepath.Abs(os.Args[0])
	require.NoError(t, err)

	assert.True(t, files.FileExists(currentFile))
	assert.False(t, files.FileExists("/not/a/real/path"))
}

func TestIsExistingFile(t *testing.T) {
	t.Parallel()

	currentFile, err := filepath.Abs(os.Args[0])
	require.NoError(t, err)

	currentFileDir := filepath.Dir(currentFile)

	assert.True(t, files.IsExistingFile(currentFile))
	assert.False(t, files.IsExistingFile("/not/a/real/path"))
	assert.False(t, files.IsExistingFile(currentFileDir))
}

func TestIsExistingDir(t *testing.T) {
	t.Parallel()

	currentFile, err := filepath.Abs(os.Args[0])
	require.NoError(t, err)

	currentFileDir := filepath.Dir(currentFile)

	assert.False(t, files.IsExistingDir(currentFile))
	assert.False(t, files.IsExistingDir("/not/a/real/path"))
	assert.True(t, files.IsExistingDir(currentFileDir))
}

func TestCopyFolderToDest(t *testing.T) {
	t.Parallel()

	tempFolderPrefix := "someprefix"
	destFolder := os.TempDir()
	tmpDir := t.TempDir()

	filter := func(path string) bool {
		return !files.PathContainsHiddenFileOrFolder(path) && !files.PathContainsTerraformState(path)
	}

	folder, err := files.CopyFolderToDest("/not/a/real/path", destFolder, tempFolderPrefix, filter)
	require.Error(t, err)
	assert.False(t, files.FileExists(folder))

	folder, err = files.CopyFolderToDest(tmpDir, destFolder, tempFolderPrefix, filter)
	assert.DirExists(t, folder)
	assert.NoError(t, err)
}

func TestCopyFolderContents(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "full-copy")
	tmpDir := t.TempDir()

	err := files.CopyFolderContents(originalDir, tmpDir)
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyFolderContentsWithHiddenFilesFilter(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-hidden-files")
	tmpDir := t.TempDir()

	err := files.CopyFolderContentsWithFilter(originalDir, tmpDir, func(path string) bool {
		return !files.PathContainsHiddenFileOrFolder(path)
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

	err := files.CopyFolderContentsWithFilter(originalDir, tmpDir, func(path string) bool {
		return !files.PathContainsHiddenFileOrFolder(path)
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
			t.Fatal(fmt.Errorf("failed to remove link: %w", err))
		}
	}()

	if err := os.Symlink(pathToFile, pathToSymlink); err != nil {
		t.Fatal(fmt.Errorf("failed to create broken link for test: %w", err))
	}

	// Test copying folder
	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "symlinks-broken")
	tmpDir := t.TempDir()

	err := files.CopyFolderContentsWithFilter(originalDir, tmpDir, func(path string) bool {
		return !files.PathContainsHiddenFileOrFolder(path)
	})
	require.NoError(t, err)

	// This requireDirectoriesEqual command uses GNU diff under the hood, but unfortunately we cannot instruct diff to
	// compare symlinks in two directories without attempting to dereference any symlinks until diff version 3.3.0.
	// Because many environments are still using diff < 3.3.0, we disregard this test for now.
	// Per https://unix.stackexchange.com/a/119406/129208
	// requireDirectoriesEqual(t, expectedDir, tmpDir)
	fmt.Println("Test completed without error, however due to a limitation in GNU diff < 3.3.0, directories have not been compared for equivalency.")
}

func TestCopyTerraformFolderToTemp(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-hidden-files-no-terraform-files")

	tmpDir, err := files.CopyTerraformFolderToTemp(originalDir, "TestCopyTerraformFolderToTemp")
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyTerraformFolderToDest(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "original")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-hidden-files-no-terraform-files")
	destFolder := os.TempDir()

	tmpDir, err := files.CopyTerraformFolderToDest(originalDir, destFolder, "TestCopyTerraformFolderToTemp")
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyTerragruntFolderToTemp(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "terragrunt-files")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-state-files")

	tmpDir, err := files.CopyTerragruntFolderToTemp(originalDir, t.Name())
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

func TestCopyTerragruntFolderToDest(t *testing.T) {
	t.Parallel()

	originalDir := filepath.Join(copyFolderContentsFixtureRoot, "terragrunt-files")
	expectedDir := filepath.Join(copyFolderContentsFixtureRoot, "no-state-files")
	destFolder := os.TempDir()

	tmpDir, err := files.CopyTerragruntFolderToDest(originalDir, destFolder, t.Name())
	require.NoError(t, err)

	requireDirectoriesEqual(t, expectedDir, tmpDir)
}

// A source folder that itself sits under a hidden directory (for example ~/.cache/project) must still be copied. Only
// hidden files and folders below the source folder should be skipped.
func TestCopyTerraformFolderToDestUnderHiddenParent(t *testing.T) {
	t.Parallel()

	originalDir := createFolderUnderHiddenParent(t)

	tmpDir, err := files.CopyTerraformFolderToDest(originalDir, t.TempDir(), "TestCopyTerraformFolderToDestUnderHiddenParent")
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "main.tf"))
	assert.FileExists(t, filepath.Join(tmpDir, "sub", "nested.tf"))
	assert.NoDirExists(t, filepath.Join(tmpDir, ".terraform"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "terraform.tfstate"))
}

func TestCopyTerragruntFolderToDestUnderHiddenParent(t *testing.T) {
	t.Parallel()

	originalDir := createFolderUnderHiddenParent(t)

	tmpDir, err := files.CopyTerragruntFolderToDest(originalDir, t.TempDir(), "TestCopyTerragruntFolderToDestUnderHiddenParent")
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmpDir, "terragrunt.hcl"))
	assert.FileExists(t, filepath.Join(tmpDir, "sub", "nested.tf"))
	assert.NoDirExists(t, filepath.Join(tmpDir, ".terraform"))
	assert.NoFileExists(t, filepath.Join(tmpDir, "terraform.tfstate"))
}

func TestFindTerraformSourceFilesInDirUnderHiddenParent(t *testing.T) {
	t.Parallel()

	originalDir := createFolderUnderHiddenParent(t)

	tfFiles, err := files.FindTerraformSourceFilesInDir(originalDir)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{
		filepath.Join(originalDir, "main.tf"),
		filepath.Join(originalDir, "sub", "nested.tf"),
	}, tfFiles)
}

func TestPathContainsTerraformStateOrVars(t *testing.T) {
	t.Parallel()

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
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			result := files.PathContainsTerraformStateOrVars(tt.path)
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

// createFolderUnderHiddenParent creates a small Terraform/Terragrunt folder whose parent directory is hidden, and returns
// the path of the folder. It also contains hidden and state files that the copy functions are expected to skip.
func createFolderUnderHiddenParent(t *testing.T) string {
	t.Helper()

	dir := filepath.Join(t.TempDir(), ".hidden-parent", "module")

	for _, name := range []string{
		"main.tf",
		"terragrunt.hcl",
		"terraform.tfstate",
		filepath.Join("sub", "nested.tf"),
		filepath.Join(".terraform", "hidden.tf"),
	} {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
		require.NoError(t, os.WriteFile(path, []byte("# test\n"), 0o600))
	}

	return dir
}

// Diffing two directories to ensure they have the exact same files, contents, etc and showing exactly what's different
// takes a lot of code. Why waste time on that when this functionality is already nicely implemented in the Unix/Linux
// "diff" command? We shell out to that command at test time.
func requireDirectoriesEqual(t *testing.T, folderWithExpectedContents string, folderWithActualContents string) {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), "diff", "-r", "-u", folderWithExpectedContents, folderWithActualContents)

	bytes, err := cmd.Output()
	output := string(bytes)

	require.NoError(t, err, "diff command exited with an error. This likely means the contents of %s and %s are different. Here is the output of the diff command:\n%s", folderWithExpectedContents, folderWithActualContents, output)
}
