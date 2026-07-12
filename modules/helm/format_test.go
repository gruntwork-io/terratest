package helm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatSetValuesAsArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setValues     map[string]string
		setStrValues  map[string]string
		setJSONValues map[string]string
		expected      []string
		expectedStr   []string
		expectedJSON  []string
	}{
		{
			name:          "EmptyValue",
			setValues:     map[string]string{},
			setStrValues:  map[string]string{},
			setJSONValues: map[string]string{},
			expected:      []string{},
			expectedStr:   []string{},
			expectedJSON:  []string{},
		},
		{
			name:          "SingleValue",
			setValues:     map[string]string{"containerImage": "null"},
			setStrValues:  map[string]string{"numericString": "123123123123"},
			setJSONValues: map[string]string{"limits": `{"cpu": 1}`},
			expected:      []string{"--set", "containerImage=null"},
			expectedStr:   []string{"--set-string", "numericString=123123123123"},
			expectedJSON:  []string{"--set-json", "limits=" + `{"cpu": 1}`},
		},
		{
			name: "MultipleValues",
			setValues: map[string]string{
				"containerImage.repository": "nginx",
				"containerImage.tag":        "v1.15.4",
			},
			setStrValues: map[string]string{
				"numericString": "123123123123",
				"otherString":   "null",
			},
			setJSONValues: map[string]string{
				"containerImage": `{"repository": "nginx", "tag": "v1.15.4"}`,
				"otherString":    "{}",
			},
			expected: []string{
				"--set", "containerImage.repository=nginx",
				"--set", "containerImage.tag=v1.15.4",
			},
			expectedStr: []string{
				"--set-string", "numericString=123123123123",
				"--set-string", "otherString=null",
			},
			expectedJSON: []string{
				"--set-json", "containerImage=" + `{"repository": "nginx", "tag": "v1.15.4"}`,
				"--set-json", "otherString={}",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, helm.FormatSetValuesAsArgs(testCase.setValues, "--set"))
			assert.Equal(t, testCase.expectedStr, helm.FormatSetValuesAsArgs(testCase.setStrValues, "--set-string"))
			assert.Equal(t, testCase.expectedJSON, helm.FormatSetValuesAsArgs(testCase.setJSONValues, "--set-json"))
		})
	}
}

func TestFormatSetFilesAsArgs(t *testing.T) {
	t.Parallel()

	paths, err := createTempFiles(2)
	require.NoError(t, err)

	t.Cleanup(func() { deleteTempFiles(paths) })

	absPathList := absPaths(t, paths)

	testCases := []struct {
		name     string
		setFiles map[string]string
		expected []string
	}{
		{
			name:     "EmptyValue",
			setFiles: map[string]string{},
			expected: []string{},
		},
		{
			name:     "SingleValue",
			setFiles: map[string]string{"containerImage": paths[0]},
			expected: []string{"--set-file", "containerImage=" + absPathList[0]},
		},
		{
			name: "MultipleValues",
			setFiles: map[string]string{
				"containerImage.repository": paths[0],
				"containerImage.tag":        paths[1],
			},
			expected: []string{
				"--set-file", "containerImage.repository=" + absPathList[0],
				"--set-file", "containerImage.tag=" + absPathList[1],
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, helm.FormatSetFilesAsArgs(t, testCase.setFiles))
		})
	}
}

func TestFormatValuesFilesAsArgs(t *testing.T) {
	t.Parallel()

	paths, err := createTempFiles(2)
	require.NoError(t, err)

	t.Cleanup(func() { deleteTempFiles(paths) })

	absPathList := absPaths(t, paths)

	testCases := []struct {
		name        string
		valuesFiles []string
		expected    []string
	}{
		{
			name:        "EmptyValue",
			valuesFiles: []string{},
			expected:    []string{},
		},
		{
			name:        "SingleValue",
			valuesFiles: []string{paths[0]},
			expected:    []string{"-f", absPathList[0]},
		},
		{
			name:        "MultipleValues",
			valuesFiles: paths,
			expected: []string{
				"-f", absPathList[0],
				"-f", absPathList[1],
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, testCase.expected, helm.FormatValuesFilesAsArgs(t, testCase.valuesFiles))
		})
	}
}

// createTempFiles will create numFiles temporary files that can pass the abspath checks.
func createTempFiles(numFiles int) ([]string, error) {
	paths := []string{}

	for i := 0; i < numFiles; i++ {
		tmpFile, err := os.CreateTemp("", "")
		// We don't use require or t.Fatal here so that we give a chance to delete any temp files that were created
		// before this error
		if err != nil {
			return paths, err
		}

		defer tmpFile.Close()

		paths = append(paths, tmpFile.Name())
	}

	return paths, nil
}

// deleteTempFiles will delete all the given temp file paths.
func deleteTempFiles(paths []string) {
	for _, path := range paths {
		os.Remove(path)
	}
}

// absPaths will return the absolute paths of each path in the list.
func absPaths(t *testing.T, paths []string) []string {
	t.Helper()

	out := make([]string, 0, len(paths))

	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		require.NoError(t, err)

		out = append(out, absPath)
	}

	return out
}
