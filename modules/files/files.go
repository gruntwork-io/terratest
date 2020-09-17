// Package files allows to interact with files on a file system.
package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

// CopyTerraformFolderToTemp creates a copy of the given folder and all its contents in a temp folder with a unique name and the given prefix.
// This is useful when running multiple tests in parallel against the same set of Terraform files to ensure the
// tests don't overwrite each other's .terraform working directory and terraform.tfstate files. This method returns
// the path to the temp folder with the copied contents. Hidden files and folders, Terraform state files, and
// terraform.tfvars files are not copied to this temp folder, as you typically don't want them interfering with your
// tests.
func CopyTerraformFolderToTemp(folderPath string, tempFolderPrefix string) (string, error) {
	filter := func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path) && !PathContainsTerraformStateOrVars(path)
	}

	destFolder, err := CopyFolderToTemp(folderPath, tempFolderPrefix, filter)
	if err != nil {
		return "", err
	}

	return destFolder, nil
}

// CopyTerragruntFolderToTemp creates a copy of the given folder and all its contents in a temp folder with a unique name and the given prefix.
// Since terragrunt uses tfvars files to specify modules, they are copied to the temporary directory as well.
// Terraform state files are excluded as well as .terragrunt-cache to avoid overwriting contents.
func CopyTerragruntFolderToTemp(folderPath string, tempFolderPrefix string) (string, error) {
	filter := func(path string) bool {
		return !PathContainsHiddenFileOrFolder(path) && !PathContainsTerraformState(path)
	}

	destFolder, err := CopyFolderToTemp(folderPath, tempFolderPrefix, filter)
	if err != nil {
		return "", err
	}

	return destFolder, nil
}

// CopyFolderToTemp creates a copy of the given folder and all its filtered contents in a temp folder
// with a unique name and the given prefix.
func CopyFolderToTemp(folderPath string, tempFolderPrefix string, filter func(path string) bool) (string, error) {
	tmpDir, err := ioutil.TempDir("", tempFolderPrefix)
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

// CopyFolderContents copies all the files and folders within the given source folder to the destination folder.
func CopyFolderContents(source string, destination string) error {
	return CopyFolderContentsWithFilter(source, destination, func(path string) bool {
		return true
	})
}

// CopyFolderContentsWithFilter copies the files and folders within the given source folder that pass the given filter (return true) to the
// destination folder.
func CopyFolderContentsWithFilter(source string, destination string, filter func(path string) bool) error {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	for _, file := range files {
		src := filepath.Join(source, file.Name())
		dest := filepath.Join(destination, file.Name())

		if !filter(src) {
			continue
		} else if file.IsDir() {
			if err := os.MkdirAll(dest, file.Mode()); err != nil {
				return err
			}

			if err := CopyFolderContentsWithFilter(src, dest, filter); err != nil {
				return err
			}

		} else if isSymLink(file) {
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

// PathContainsTerraformStateOrVars returns true if the path corresponds to a Terraform state file or .tfvars file.
func PathContainsTerraformStateOrVars(path string) bool {
	filename := filepath.Base(path)
	return filename == "terraform.tfstate" || filename == "terraform.tfstate.backup" || filename == "terraform.tfvars"
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

// CopyFile copies a file from source to destination.
func CopyFile(source string, destination string) error {
	contents, err := ioutil.ReadFile(source)
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

	return ioutil.WriteFile(destination, contents, fileInfo.Mode())
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
