package test_structure

import (
	"fmt"
	"path"
	"path/filepath"

	go_commons_collections "github.com/gruntwork-io/go-commons/collections"
	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/mattn/go-zglob"
)

// ValidateFileType is the underlying module type to search for when performing validation. Either Terraform or Terragrunt
// files are targeted during a given validation sweep
type ValidateFileType string

const (
	// TF represents repositories that contain Terraform code
	TF = "*.tf"
	// TG represents repositories that contain Terragrunt code
	TG = "terragrunt.hcl"
)

// ValidationOptions represent the configuration for a given validation sweep of a target repo
type ValidationOptions struct {
	// The target directory to recursively search for all Terraform directories (those that contain .tf files)
	// If you provide RootDir and do not pass entries in either IncludeDirs or ExcludeDirs, then all Terraform directories
	// From the RootDir, recursively, will be validated
	RootDir  string
	FileType ValidateFileType
	// If you only want to include certain sub directories of RootDir, add the absolute paths here. For example, if the
	// RootDir is /home/project and you want to only include /home/project/examples, add /home/project/examples here
	// Note that while the struct requires full paths, you can pass relative paths to the NewValidationOptions function
	// which will build the full paths based on the supplied RootDir
	IncludeDirs []string
	// If you want to explicitly exclude certain sub directories of RootDir, add the absolute paths here. For example, if the
	// RootDir is /home/project and you want to include everything EXCEPT /home/project/modules, add
	// /home/project/modules to this slice. Note that ExcludeDirs is only considered when IncludeDirs is not passed
	// Note that while the struct requires full paths, you can pass relative paths to the NewValidationOptions function
	// which will build the full paths based on the supplied RootDir
	ExcludeDirs []string
}

// CloneWithNewRootDir clones the given opts with a new root dir. Updates all include and exclude dirs to be relative
// to the new root dir.
func CloneWithNewRootDir(opts *ValidationOptions, newRootDir string) (*ValidationOptions, error) {
	includeDirs, err := buildRelPathsFromFull(opts.RootDir, opts.IncludeDirs)
	if err != nil {
		return nil, err
	}
	excludeDirs, err := buildRelPathsFromFull(opts.RootDir, opts.ExcludeDirs)
	if err != nil {
		return nil, err
	}

	out, err := NewValidationOptions(newRootDir, includeDirs, excludeDirs)
	if err != nil {
		return nil, err
	}
	out.FileType = opts.FileType
	return out, nil
}

// configureBaseValidationOptions returns a pointer to a ValidationOptions struct configured with sane, override-able defaults
// Note that the ValidationOptions's fields IncludeDirs and ExcludeDirs must be absolute paths, but this method will accept relative paths
// and build the absolute paths when instantiating the ValidationOptions struct,  making it the preferred means of configuring
// ValidationOptions.
//
// For example, if your RootDir is /home/project/ and you want to exclude "modules" and "test" you need
// only pass the relative paths in your excludeDirs slice like so:
// opts, err := NewValidationOptions("/home/project", []string{}, []string{"modules", "test"})
func configureBaseValidationOptions(rootDir string, includeDirs, excludeDirs []string) (*ValidationOptions, error) {
	vo := &ValidationOptions{
		RootDir:     "",
		IncludeDirs: []string{},
		ExcludeDirs: []string{},
	}

	if rootDir == "" {
		return nil, ValidationUndefinedRootDirErr{}
	}

	if !filepath.IsAbs(rootDir) {
		rootDirAbs, err := filepath.Abs(rootDir)
		if err != nil {
			return nil, ValidationAbsolutePathErr{rootDir: rootDir}
		}
		rootDir = rootDirAbs
	}

	vo.RootDir = filepath.Clean(rootDir)

	if len(includeDirs) > 0 {
		vo.IncludeDirs = buildFullPathsFromRelative(vo.RootDir, includeDirs)
	}

	if len(excludeDirs) > 0 {
		vo.ExcludeDirs = buildFullPathsFromRelative(vo.RootDir, excludeDirs)
	}

	return vo, nil
}

// NewValidationOptions returns a ValidationOptions struct, with override-able sane defaults, configured to find
// and process all directories containing .tf files
func NewValidationOptions(rootDir string, includeDirs, excludeDirs []string) (*ValidationOptions, error) {
	opts, err := configureBaseValidationOptions(rootDir, includeDirs, excludeDirs)
	if err != nil {
		return opts, err
	}
	opts.FileType = TF
	return opts, nil
}

// NewTerragruntValidationOptions returns a ValidationOptions struct, with override-able sane defaults, configured to find
// and process all directories containing .hcl files.
func NewTerragruntValidationOptions(rootDir string, includeDirs, excludeDirs []string) (*ValidationOptions, error) {
	opts, err := configureBaseValidationOptions(rootDir, includeDirs, excludeDirs)
	if err != nil {
		return opts, err
	}
	opts.FileType = TG
	return opts, nil
}

func buildRelPathsFromFull(rootDir string, fullPaths []string) ([]string, error) {
	var relPaths []string
	for _, maybeFullPath := range fullPaths {
		if filepath.IsAbs(maybeFullPath) {
			relPath, err := filepath.Rel(rootDir, maybeFullPath)
			if err != nil {
				return nil, err
			}
			relPaths = append(relPaths, relPath)
		} else {
			relPaths = append(relPaths, maybeFullPath)
		}
	}
	return relPaths, nil
}

func buildFullPathsFromRelative(rootDir string, relativePaths []string) []string {
	var fullPaths []string
	for _, maybeRelativePath := range relativePaths {
		// If the relativePath is already an absolute path, don't modify it
		if filepath.IsAbs(maybeRelativePath) {
			fullPaths = append(fullPaths, filepath.Clean(maybeRelativePath))
		} else {
			fullPaths = append(fullPaths, filepath.Clean(filepath.Join(rootDir, maybeRelativePath)))
		}
	}
	return fullPaths
}

// FindTerraformModulePathsInRootE returns a slice strings representing the filepaths for all valid Terraform / Terragrunt
// modules in the given RootDir, subject to the include / exclude filters.
func FindTerraformModulePathsInRootE(opts *ValidationOptions) ([]string, error) {
	// Find all Terraform / Terragrunt files (as specified by opts.FileType) from the configured RootDir
	pattern := fmt.Sprintf("%s/**/%s", opts.RootDir, opts.FileType)
	matches, err := zglob.Glob(pattern)
	if err != nil {
		return matches, err
	}
	// Keep a unique set of the base dirs that contain Terraform / Terragrunt files
	terraformDirSet := make(map[string]string)
	for _, match := range matches {
		// The glob match returns all full paths to every target file, whereas we're only interested in their root
		// directories for the purposes of running Terraform validate / terragrunt validate-inputs
		rootDir := path.Dir(match)
		// Don't include hidden .terraform directories when finding paths to validate
		if !files.PathContainsHiddenFileOrFolder(rootDir) {
			terraformDirSet[rootDir] = "exists"
		}
	}

	// Retrieve just the unique paths to each Terraform module directory from the map we're using as a set
	terraformDirs := go_commons_collections.Keys(terraformDirSet)

	if len(opts.IncludeDirs) > 0 {
		terraformDirs = collections.ListIntersection(terraformDirs, opts.IncludeDirs)
	}

	if len(opts.ExcludeDirs) > 0 {
		terraformDirs = collections.ListSubtract(terraformDirs, opts.ExcludeDirs)
	}

	// Filter out any filepaths that were explicitly included in opts.ExcludeDirs
	return terraformDirs, nil
}

// Custom error types

// ValidationAbsolutePathErr is returned when NewValidationOptions was unable to convert a non-absolute RootDir to
// an absolute path
type ValidationAbsolutePathErr struct {
	rootDir string
}

func (e ValidationAbsolutePathErr) Error() string {
	return fmt.Sprintf("Could not convert RootDir: %s to absolute path", e.rootDir)
}

// ValidationUndefinedRootDirErr is returned when NewValidationOptions is called without a RootDir argument
type ValidationUndefinedRootDirErr struct{}

func (e ValidationUndefinedRootDirErr) Error() string {
	return "RootDir must be defined in ValidationOptions passed to ValidateAllTerraformModules"
}
