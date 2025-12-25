package parser

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files/v2"
	"github.com/gruntwork-io/terratest/modules/shell/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func DirectoryEqual(t *testing.T, dirA string, dirB string) bool {
	dirAAbs, err := filepath.Abs(dirA)
	if err != nil {
		t.Fatal(err)
	}
	dirBAbs, err := filepath.Abs(dirB)
	if err != nil {
		t.Fatal(err)
	}
	// We use diff here instead of using something in go for simplicity of comparing directories and file contents
	// recursively
	cmd := shell.Command{
		Command: "diff",
		Args:    []string{"-ar", dirAAbs, dirBAbs},
	}
	err = shell.RunCommandE(t, cmd)
	exitCode, err := shell.GetExitCodeForRunCommandError(err)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode == 0
}

func openFile(t *testing.T, filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Error opening file: %s", err)
	}
	return file
}

func testExample(t *testing.T, example string) {
	expected, output := path.Join(t.TempDir(), "expected"), path.Join(t.TempDir(), "output")
	require.NoError(t, os.Mkdir(expected, 0755))
	require.NoError(t, os.Mkdir(output, 0755))

	// prepare expected directory to diff against
	expectedOutputDirName := fmt.Sprintf("./fixtures/%s_example_expected", example)
	require.NoError(t, files.CopyFolderContents(expectedOutputDirName, expected))
	b, err := os.ReadFile(path.Join(expected, "report.xml"))
	require.NoError(t, err)
	b = bytes.ReplaceAll(b, []byte("go1.21.1"), []byte(runtime.Version())) // replace the harcoded go version of the fixture
	require.NoError(t, os.WriteFile(path.Join(expected, "report.xml"), b, 644))

	// run the parser
	logger := NewTestLogger(t)
	logFileName := fmt.Sprintf("./fixtures/%s_example.log", example)
	file := openFile(t, logFileName)
	SpawnParsers(logger, file, output)

	// assert
	assert.True(t, DirectoryEqual(t, expected, output))
}

func TestIntegrationBasicExample(t *testing.T) {
	t.Parallel()
	testExample(t, "basic")
}

func TestIntegrationFailingExample(t *testing.T) {
	t.Parallel()
	testExample(t, "failing")
}

func TestIntegrationPanicExample(t *testing.T) {
	t.Parallel()
	testExample(t, "panic")
}

func TestIntegrationNewGoExample(t *testing.T) {
	t.Parallel()
	testExample(t, "new_go_failing")
}
