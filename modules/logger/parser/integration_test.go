package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/assert"
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
	logger := NewTestLogger(t)
	dir := getTempDir(t)
	defer os.RemoveAll(dir)
	logFileName := fmt.Sprintf("./fixtures/%s_example.log", example)
	expectedOutputDirName := fmt.Sprintf("./fixtures/%s_example_expected", example)
	file := openFile(t, logFileName)
	SpawnParsers(logger, file, dir)
	assert.True(t, DirectoryEqual(t, dir, expectedOutputDirName))
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
