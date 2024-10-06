// Package files allows to interact with files on a file system.
package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

// DefaultPathFilter is the default filter used by CopyTerraformFolderToDest to determine which files to copy.
func DefaultPathFilter(path string) bool {
	if PathIsTerraformVersionFile(path) || PathIsTerraformLockFile(path) {
		return true
	}
	if PathContainsHiddenFileOrFolder(path) || PathContainsTerraformStateOrVars(path) {
		return false
	}
	return true
}

// FileExists returns true if the given file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileExistsE returns true if the given file exists
// It will return an error if os.Stat error is not an ErrNotExist
func FileExistsE(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	return err == nil, nil
}

// IsExistingFile returns true if the path exists and is a file.
func IsExistingFile(path string) bool {
	fileInfo, err := os.Stat(path)
	return err == nil && !fileInfo.IsDir()
}

// IsExistingDir returns true if the path exists and is a directory
func IsExistingDir(path string) bool {
	fileInfo, err := os.Stat(path)
	return err == nil && fileInfo.IsDir()
}

// CopyTerraformFolderToDestWithFilter creates a copy of the given folder and all its contents in a specified folder with a unique name and the given prefix.
// This is useful when running multiple tests in parallel against the same set of Terraform files to ensure the
// tests don't overwrite each other's .terraform working directory and terraform.tfstate files. This method returns
// the path to the dest folder with the copied contents. You can configure which files you want to copy by passing in a filter function.
// This method is useful when running through a build tool so the files are copied to a destination that is cleaned on each run of the pipeline.
func CopyTerraformFolderToDestWithFilter(folderPath string, destRootFolder string, tempFolderPrefix string, filter func(string) bool) (string, error) {
	destFolder, err := CopyFolderToDest(folderPath, destRootFolder, tempFolderPrefix, filter)
	if err != nil {
		return "", err
	}

	return destFolder, nil
}

// CopyTerraformFolderToDest creates a copy of the given folder and all its contents in a specified folder with a unique name and the given prefix.
// This is useful when running multiple tests in parallel against the same set of Terraform files to ensure the
// tests don't overwrite each other's .terraform working directory and terraform.tfstate files. This method returns
// the path to the dest folder with the copied contents. Hidden files and folders (with the exception of the `.terraform-version` files used
// by the [mise tool](https://github.com/jdx/mise) and `.terraform.lock.hcl` used by Terraform to lock providers versions), Terraform state
// files, and terraform.tfvars files are not copied to this temp folder, as you typically don't want them interfering with your tests.
// This method is useful when running through a build tool so the files are copied to a destination that is cleaned on each run of the pipeline.
func CopyTerraformFolderToDest(folderPath string, destRootFolder string, tempFolderPrefix string) (string, error) {
	return CopyTerraformFolderToDestWithFilter(folderPath, destRootFolder, tempFolderPrefix, DefaultPathFilter)
}

// CopyTerraformFolderToTempWithFilter calls CopyTerraformFolderToDestWithFilter, passing os.TempDir() as the root destination folder.
func CopyTerraformFolderToTempWithFilter(folderPath string, tempFolderPrefix string, filter func(string) bool) (string, error) {
	return CopyTerraformFolderToDestWithFilter(folderPath, os.TempDir(), tempFolderPrefix, filter)
}

// CopyTerraformFolderToTemp calls CopyTerraformFolderToDest, passing os.TempDir() as the root destination folder.
func CopyTerraformFolderToTemp(folderPath string, tempFolderPrefix string) (string, error) {
	return CopyTerraformFolderToDest(folderPath, os.TempDir(), tempFolderPrefix)
}

// CopyTerragruntFolderToDest creates a copy of the given folder and all its contents in a specified folder with a unique name and the given prefix.
// Since terragrunt uses tfvars files to specify modules, they are copied to the directory as well.
// Terraform state files are excluded as well as .terragrunt-cache to avoid overwriting contents.
func CopyTerragruntFolderToDest(folderPath string, destRootFolder string, tempFolderPrefix string) (string, error) {
	filter := func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path) && !PathContainsTerraformState(path)
	}

	destFolder, err := CopyFolderToDest(folderPath, destRootFolder, tempFolderPrefix, filter)
	if err != nil {
		return "", err
	}

	return destFolder, nil
}

// CopyTerragruntFolderToTemp calls CopyTerragruntFolderToDest, passing os.TempDir() as the root destination folder.
func CopyTerragruntFolderToTemp(folderPath string, tempFolderPrefix string) (string, error) {
	return CopyTerragruntFolderToDest(folderPath, os.TempDir(), tempFolderPrefix)
}

// CopyFolderToDest creates a copy of the given folder and all its filtered contents in a temp folder
// with a unique name and the given prefix.
func CopyFolderToDest(folderPath string, destRootFolder string, tempFolderPrefix string, filter func(path string) bool) (string, error) {
	destRootExists, err := FileExistsE(destRootFolder)
	if err != nil {
		return "", err
	}
	if !destRootExists {
		return "", DirNotFoundError{Directory: destRootFolder}
	}

	exists, err := FileExistsE(folderPath)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", DirNotFoundError{Directory: folderPath}
	}

	tmpDir, err := os.MkdirTemp(destRootFolder, tempFolderPrefix)
	if err != nil {
		return "", err
	}

	// Inside of the temp folder, we create a subfolder that preserves the name of the folder we're copying from.
	absFolderPath, err := filepath.Abs(folderPath)
	if err != nil {
		return "", err
	}
	folderName := filepath.Base(absFolderPath)
	destFolder := filepath.Join(tmpDir, folderName)

	if err := os.MkdirAll(destFolder, 0777); err != nil {
		return "", err
	}

	if err := CopyFolderContentsWithFilter(folderPath, destFolder, filter); err != nil {
		return "", err
	}

	return destFolder, nil
}

// CopyFolderToTemp calls CopyFolderToDest, passing os.TempDir() as the root destination folder.
func CopyFolderToTemp(folderPath string, tempFolderPrefix string, filter func(path string) bool) (string, error) {
	return CopyFolderToDest(folderPath, os.TempDir(), tempFolderPrefix, filter)
}

// CopyFolderContents copies all the files and folders within the given source folder to the destination folder.
func CopyFolderContents(source string, destination string) error {
	return CopyFolderContentsWithFilter(source, destination, func(path string) bool {
		return true
	})
}

// CopyFolderContentsWithFilter copies the files and folders within the given source folder that pass the given filter (return true) to the
// destination folder.
func CopyFolderContentsWithFilter(source string, destination string, filter func(path string) bool) error {
	files, err := os.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		src := filepath.Join(source, file.Name())
		dest := filepath.Join(destination, file.Name())
		f, _ := file.Info()
		if !filter(src) {
			continue
		} else if file.IsDir() {
			if err := os.MkdirAll(dest, f.Mode()); err != nil {
				return err
			}

			if err := CopyFolderContentsWithFilter(src, dest, filter); err != nil {
				return err
			}

		} else if isSymLink(f) {
			if err := copySymLink(src, dest); err != nil {
				return err
			}
		} else {
			if err := CopyFile(src, dest); err != nil {
				return err
			}
		}
	}

	return nil
}

// PathContainsTerraformStateOrVars returns true if the path corresponds to a Terraform state file or .tfvars/.tfvars.json file.
func PathContainsTerraformStateOrVars(path string) bool {
	filename := filepath.Base(path)
	return filename == "terraform.tfstate" || filename == "terraform.tfstate.backup" || filename == "terraform.tfvars" || filename == "terraform.tfvars.json"
}

// PathContainsTerraformState returns true if the path corresponds to a Terraform state file.
func PathContainsTerraformState(path string) bool {
	filename := filepath.Base(path)
	return filename == "terraform.tfstate" || filename == "terraform.tfstate.backup"
}

// PathContainsHiddenFileOrFolder returns true if the given path contains a hidden file or folder.
func PathContainsHiddenFileOrFolder(path string) bool {
	pathParts := strings.Split(path, string(filepath.Separator))
	for _, pathPart := range pathParts {
		if strings.HasPrefix(pathPart, ".") && pathPart != "." && pathPart != ".." {
			return true
		}
	}
	return false
}

// PathIsTerraformVersionFile returns true if the given path is the special '.terraform-version' file used by the [mise](https://github.com/jdx/mise) tool.
func PathIsTerraformVersionFile(path string) bool {
	return filepath.Base(path) == ".terraform-version"
}

// PathIsTerraformLockFile return true if the given path is the special '.terraform.lock.hcl' file used by Terraform to lock providers versions
func PathIsTerraformLockFile(path string) bool {
	return filepath.Base(path) == ".terraform.lock.hcl"
}

// CopyFile copies a file from source to destination.
func CopyFile(source string, destination string) error {
	contents, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	return WriteFileWithSamePermissions(source, destination, contents)
}

// WriteFileWithSamePermissions writes a file to the given destination with the given contents using the same permissions as the file at source.
func WriteFileWithSamePermissions(source string, destination string, contents []byte) error {
	fileInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	return os.WriteFile(destination, contents, fileInfo.Mode())
}

// isSymLink returns true if the given file is a symbolic link
// Per https://stackoverflow.com/a/18062079/2308858
func isSymLink(file os.FileInfo) bool {
	return file.Mode()&os.ModeSymlink != 0
}

// copySymLink copies the source symbolic link to the given destination.
func copySymLink(source string, destination string) error {
	symlinkPath, err := os.Readlink(source)
	if err != nil {
		return err
	}

	err = os.Symlink(symlinkPath, destination)
	if err != nil {
		return err
	}

	return nil
}

// FindTerraformSourceFilesInDir given a directory path, finds all the terraform source files contained in it. This will
// recursively search subdirectories, but will ignore any hidden files (which in turn ignores terraform data dirs like
// .terraform folder).
func FindTerraformSourceFilesInDir(dirPath string) ([]string, error) {
	pattern := fmt.Sprintf("%s/**/*.tf", dirPath)
	matches, err := zglob.Glob(pattern)
	if err != nil {
		return nil, err
	}

	tfFiles := []string{}
	for _, match := range matches {
		// Don't include hidden .terraform directories when finding paths to validate
		if !PathContainsHiddenFileOrFolder(match) {
			tfFiles = append(tfFiles, match)
		}
	}
	return tfFiles, nil
}
