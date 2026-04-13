package helm

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/collections"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// FormatSetValuesAsArgs formats the given values as command line args for helm using the given flag (e.g flags of
// the format "--set"/"--set-string"/"--set-json" resulting in args like --set/set-string/set-json key=value...).
func FormatSetValuesAsArgs(setValues map[string]string, flag string) []string {
	args := make([]string, 0, len(setValues)*2) //nolint:mnd // Each key-value pair produces 2 args (flag + value).

	// To make it easier to test, go through the keys in sorted order
	keys := collections.Keys(setValues)
	for _, key := range keys {
		value := setValues[key]
		argValue := fmt.Sprintf("%s=%s", key, value)
		args = append(args, flag, argValue)
	}

	return args
}

// FormatValuesFilesAsArgs formats the given list of values file paths as command line args for helm (e.g of the format
// -f path). This will fail the test if one of the paths do not exist or the absolute path can not be determined.
//
// Deprecated: Use [FormatValuesFilesAsArgsContext] instead.
func FormatValuesFilesAsArgs(t testing.TestingT, valuesFiles []string) []string {
	return FormatValuesFilesAsArgsContext(t, context.Background(), valuesFiles)
}

// FormatValuesFilesAsArgsContext formats the given list of values file paths as command line args for helm (e.g of the
// format -f path). The ctx parameter is accepted for API consistency with other Context-aware helpers.
// This will fail the test if one of the paths do not exist or the absolute path can not be determined.
func FormatValuesFilesAsArgsContext(t testing.TestingT, ctx context.Context, valuesFiles []string) []string {
	args, err := FormatValuesFilesAsArgsContextE(ctx, valuesFiles)
	require.NoError(t, err)

	return args
}

// FormatValuesFilesAsArgsE formats the given list of values file paths as command line args for helm (e.g of the format
// -f path). This will error if the file does not exist.
//
// Deprecated: Use [FormatValuesFilesAsArgsContextE] instead.
func FormatValuesFilesAsArgsE(valuesFiles []string) ([]string, error) {
	return FormatValuesFilesAsArgsContextE(context.Background(), valuesFiles)
}

// FormatValuesFilesAsArgsContextE formats the given list of values file paths as command line args for helm (e.g of the
// format -f path). The ctx parameter is accepted for API consistency with other Context-aware helpers.
// This will error if the file does not exist.
func FormatValuesFilesAsArgsContextE(ctx context.Context, valuesFiles []string) ([]string, error) {
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

// FormatSetFilesAsArgs formats the given list of keys and file paths as command line args for helm to set from file
// (e.g of the format --set-file key=path). This will fail the test if one of the paths do not exist or the absolute
// path can not be determined.
//
// Deprecated: Use [FormatSetFilesAsArgsContext] instead.
func FormatSetFilesAsArgs(t testing.TestingT, setFiles map[string]string) []string {
	return FormatSetFilesAsArgsContext(t, context.Background(), setFiles)
}

// FormatSetFilesAsArgsContext formats the given list of keys and file paths as command line args for helm to set from
// file (e.g of the format --set-file key=path). The ctx parameter is accepted for API consistency with other
// Context-aware helpers. This will fail the test if one of the paths do not exist or the absolute path can not be
// determined.
func FormatSetFilesAsArgsContext(t testing.TestingT, ctx context.Context, setFiles map[string]string) []string {
	args, err := FormatSetFilesAsArgsContextE(ctx, setFiles)
	require.NoError(t, err)

	return args
}

// FormatSetFilesAsArgsE formats the given list of keys and file paths as command line args for helm to set from file
// (e.g of the format --set-file key=path).
//
// Deprecated: Use [FormatSetFilesAsArgsContextE] instead.
func FormatSetFilesAsArgsE(setFiles map[string]string) ([]string, error) {
	return FormatSetFilesAsArgsContextE(context.Background(), setFiles)
}

// FormatSetFilesAsArgsContextE formats the given list of keys and file paths as command line args for helm to set from
// file (e.g of the format --set-file key=path). The ctx parameter is accepted for API consistency with other
// Context-aware helpers.
func FormatSetFilesAsArgsContextE(ctx context.Context, setFiles map[string]string) ([]string, error) {
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
