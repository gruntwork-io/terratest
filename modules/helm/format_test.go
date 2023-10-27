package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatSetValuesAsArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setValues     map[string]string
		setStrValues  map[string]string
		setJsonValues map[string]string
		expected      []string
		expectedStr   []string
		expectedJson  []string
	}{
		{
			"EmptyValue",
			map[string]string{},
			map[string]string{},
			map[string]string{},
			[]string{},
			[]string{},
			[]string{},
		},
		{
			"SingleValue",
			map[string]string{"containerImage": "null"},
			map[string]string{"numericString": "123123123123"},
			map[string]string{"limits": `{"cpu": 1}`},
			[]string{"--set", "containerImage=null"},
			[]string{"--set-string", "numericString=123123123123"},
			[]string{"--set-json", fmt.Sprintf("limits=%s", `{"cpu": 1}`)},
		},
		{
			"MultipleValues",
			map[string]string{
				"containerImage.repository": "nginx",
				"containerImage.tag":        "v1.15.4",
			},
			map[string]string{
				"numericString": "123123123123",
				"otherString":   "null",
			},
			map[string]string{
				"containerImage": `{"repository": "nginx", "tag": "v1.15.4"}`,
				"otherString":    "{}",
			},
			[]string{
				"--set", "containerImage.repository=nginx",
				"--set", "containerImage.tag=v1.15.4",
			},
			[]string{
				"--set-string", "numericString=123123123123",
				"--set-string", "otherString=null",
			},
			[]string{
				"--set-json", fmt.Sprintf("containerImage=%s", `{"repository": "nginx", "tag": "v1.15.4"}`),
				"--set-json", "otherString={}",
			},
		},
	}

	for _, testCase := range testCases {
		// Capture the range value and force it into this scope. Otherwise, it is defined outside this block so it can
		// change when the subtests parallelize and switch contexts.
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, formatSetValuesAsArgs(testCase.setValues, "--set"), testCase.expected)
			assert.Equal(t, formatSetValuesAsArgs(testCase.setStrValues, "--set-string"), testCase.expectedStr)
			assert.Equal(t, formatSetValuesAsArgs(testCase.setJsonValues, "--set-json"), testCase.expectedJson)
		})
	}
}

func TestFormatSetFilesAsArgs(t *testing.T) {
	t.Parallel()

	paths, err := createTempFiles(2)
	defer deleteTempFiles(paths)
	require.NoError(t, err)
	absPathList := absPaths(t, paths)

	testCases := []struct {
		name     string
		setFiles map[string]string
		expected []string
	}{
		{
			"EmptyValue",
			map[string]string{},
			[]string{},
		},
		{
			"SingleValue",
			map[string]string{"containerImage": paths[0]},
			[]string{"--set-file", fmt.Sprintf("containerImage=%s", absPathList[0])},
		},
		{
			"MultipleValues",
			map[string]string{
				"containerImage.repository": paths[0],
				"containerImage.tag":        paths[1],
			},
			[]string{
				"--set-file", fmt.Sprintf("containerImage.repository=%s", absPathList[0]),
				"--set-file", fmt.Sprintf("containerImage.tag=%s", absPathList[1]),
			},
		},
	}

	// We create a subtest group that is NOT parallel, so the main test waits for all the tests to finish. This way, we
	// don't delete the files until the subtests finish.
	t.Run("group", func(t *testing.T) {
		for _, testCase := range testCases {
			// Capture the range value and force it into this scope. Otherwise, it is defined outside this block so it can
			// change when the subtests parallelize and switch contexts.
			testCase := testCase

			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, formatSetFilesAsArgs(t, testCase.setFiles), testCase.expected)
			})
		}
	})
}

func TestFormatValuesFilesAsArgs(t *testing.T) {
	t.Parallel()

	paths, err := createTempFiles(2)
	defer deleteTempFiles(paths)
	require.NoError(t, err)
	absPathList := absPaths(t, paths)

	testCases := []struct {
		name        string
		valuesFiles []string
		expected    []string
	}{
		{
			"EmptyValue",
			[]string{},
			[]string{},
		},
		{
			"SingleValue",
			[]string{paths[0]},
			[]string{"-f", absPathList[0]},
		},
		{
			"MultipleValues",
			paths,
			[]string{
				"-f", absPathList[0],
				"-f", absPathList[1],
			},
		},
	}

	// We create a subtest group that is NOT parallel, so the main test waits for all the tests to finish. This way, we
	// don't delete the files until the subtests finish.
	t.Run("group", func(t *testing.T) {
		for _, testCase := range testCases {
			// Capture the range value and force it into this scope. Otherwise, it is defined outside this block so it can
			// change when the subtests parallelize and switch contexts.
			testCase := testCase

			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()
				assert.Equal(t, formatValuesFilesAsArgs(t, testCase.valuesFiles), testCase.expected)
			})
		}
	})
}

// createTempFiles will create numFiles temporary files that can pass the abspath checks.
func createTempFiles(numFiles int) ([]string, error) {
	paths := []string{}
	for i := 0; i < numFiles; i++ {
		tmpFile, err := os.CreateTemp("", "")
		defer tmpFile.Close()
		// We don't use require or t.Fatal here so that we give a chance to delete any temp files that were created
		// before this error
		if err != nil {
			return paths, err
		}
		paths = append(paths, tmpFile.Name())
	}
	return paths, nil
}

// deleteTempFiles will delete all the given temp file paths
func deleteTempFiles(paths []string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

// absPaths will return the absolute paths of each path in the list
func absPaths(t *testing.T, paths []string) []string {
	out := []string{}
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		require.NoError(t, err)
		out = append(out, absPath)
	}
	return out
}
