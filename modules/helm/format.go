package helm

import (
	"fmt"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/collections"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// formatSetValuesAsArgs formats the given values as command line args for helm using the given flag (e.g flags of
// the format "--set"/"--set-string" resulting in args like --set/set-string key=value...)
func formatSetValuesAsArgs(setValues map[string]string, flag string) []string {
	args := []string{}

	// To make it easier to test, go through the keys in sorted order
	keys := collections.Keys(setValues)
	for _, key := range keys {
		value := setValues[key]
		argValue := fmt.Sprintf("%s=%s", key, value)
		args = append(args, flag, argValue)
	}

	return args
}

// formatValuesFilesAsArgs formats the given list of values file paths as command line args for helm (e.g of the format
// -f path). This will fail the test if one of the paths do not exist or the absolute path can not be determined.
func formatValuesFilesAsArgs(t testing.TestingT, valuesFiles []string) []string {
	args, err := formatValuesFilesAsArgsE(t, valuesFiles)
	require.NoError(t, err)
	return args
}

// formatValuesFilesAsArgsE formats the given list of values file paths as command line args for helm (e.g of the format
// -f path). This will error if the file does not exist.
func formatValuesFilesAsArgsE(t testing.TestingT, valuesFiles []string) ([]string, error) {
	args := []string{}

	for _, valuesFilePath := range valuesFiles {
		// Pass through filepath.Abs to clean the path, and then make sure this file exists
		absValuesFilePath, err := filepath.Abs(valuesFilePath)
		if err != nil {
			return args, errors.WithStackTrace(err)
		}
		if !files.FileExists(absValuesFilePath) {
			return args, errors.WithStackTrace(ValuesFileNotFoundError{valuesFilePath})
		}
		args = append(args, "-f", absValuesFilePath)
	}

	return args, nil
}

// formatSetFilesAsArgs formats the given list of keys and file paths as command line args for helm to set from file
// (e.g of the format --set-file key=path). This will fail the test if one of the paths do not exist or the absolute
// path can not be determined.
func formatSetFilesAsArgs(t testing.TestingT, setFiles map[string]string) []string {
	args, err := formatSetFilesAsArgsE(t, setFiles)
	require.NoError(t, err)
	return args
}

// formatSetFilesAsArgsE formats the given list of keys and file paths as command line args for helm to set from file
// (e.g of the format --set-file key=path)
func formatSetFilesAsArgsE(t testing.TestingT, setFiles map[string]string) ([]string, error) {
	args := []string{}

	// To make it easier to test, go through the keys in sorted order
	keys := collections.Keys(setFiles)
	for _, key := range keys {
		setFilePath := setFiles[key]
		// Pass through filepath.Abs to clean the path, and then make sure this file exists
		absSetFilePath, err := filepath.Abs(setFilePath)
		if err != nil {
			return args, errors.WithStackTrace(err)
		}
		if !files.FileExists(absSetFilePath) {
			return args, errors.WithStackTrace(SetFileNotFoundError{setFilePath})
		}
		argValue := fmt.Sprintf("%s=%s", key, absSetFilePath)
		args = append(args, "--set-file", argValue)
	}

	return args, nil
}
