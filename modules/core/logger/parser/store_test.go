package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/go-commons/files"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger/parser"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createLogWriter(t *testing.T) parser.LogWriter {
	t.Helper()

	logWriter := parser.LogWriter{
		Lookup:    make(map[string]*os.File),
		OutputDir: t.TempDir(),
	}

	return logWriter
}

func TestEnsureDirectoryExistsCreatesDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	logger := NewTestLogger(t)
	tmpd := filepath.Join(dir, "tmpdir")
	assert.False(t, files.IsDir(tmpd))

	parser.EnsureDirectoryExists(logger, tmpd)
	assert.True(t, files.IsDir(tmpd))
}

func TestEnsureDirectoryExistsHandlesExistingDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	logger := NewTestLogger(t)
	assert.True(t, files.IsDir(dir))

	parser.EnsureDirectoryExists(logger, dir)
	assert.True(t, files.IsDir(dir))
}

func TestGetOrCreateFileCreatesNewFile(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)

	logger := NewTestLogger(t)
	testFileName := filepath.Join(logWriter.OutputDir, t.Name()+".log")
	assert.False(t, files.FileExists(testFileName))

	file, err := logWriter.GetOrCreateFile(logger, t.Name())
	require.NoError(t, err)

	defer file.Close()

	assert.NotNil(t, file)
	assert.True(t, files.FileExists(testFileName))
}

func TestGetOrCreateFileCreatesNewFileIfTestNameHasDir(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)

	logger := NewTestLogger(t)
	dirName := filepath.Join(logWriter.OutputDir, "TestMain")
	testFileName := filepath.Join(dirName, t.Name()+".log")
	assert.False(t, files.IsDir(dirName))
	assert.False(t, files.FileExists(testFileName))

	file, err := logWriter.GetOrCreateFile(logger, filepath.Join("TestMain", t.Name()))
	require.NoError(t, err)

	defer file.Close()

	assert.NotNil(t, file)
	assert.True(t, files.IsDir(dirName))
	assert.True(t, files.FileExists(testFileName))
}

func TestGetOrCreateChannelReturnsExistingFileHandle(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)

	testName := t.Name()
	logger := NewTestLogger(t)
	testFileName := filepath.Join(logWriter.OutputDir, t.Name())

	file, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", testFileName)
	}

	defer file.Close()

	logWriter.Lookup[testName] = file

	lookupFile, err := logWriter.GetOrCreateFile(logger, testName)
	require.NoError(t, err)
	assert.Equal(t, file, lookupFile)
}

func TestCloseFilesClosesAll(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)

	logger := NewTestLogger(t)
	testName := t.Name()
	testFileName := filepath.Join(logWriter.OutputDir, testName)

	testFile, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", testFileName)
	}

	alternativeTestName := t.Name() + "Alternative"
	alternativeTestFileName := filepath.Join(logWriter.OutputDir, alternativeTestName)

	alternativeTestFile, err := os.Create(alternativeTestFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", alternativeTestFileName)
	}

	logWriter.Lookup[testName] = testFile
	logWriter.Lookup[alternativeTestName] = alternativeTestFile

	logWriter.CloseFiles(logger)

	err = testFile.Close()
	assert.Contains(t, err.Error(), os.ErrClosed.Error())

	err = alternativeTestFile.Close()
	assert.Contains(t, err.Error(), os.ErrClosed.Error())
}

func TestWriteLogWritesToCorrectLogFile(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)

	logger := NewTestLogger(t)
	testName := t.Name()
	testFileName := filepath.Join(logWriter.OutputDir, testName)

	testFile, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", testFileName)
	}

	defer testFile.Close()

	alternativeTestName := t.Name() + "Alternative"
	alternativeTestFileName := filepath.Join(logWriter.OutputDir, alternativeTestName)

	alternativeTestFile, err := os.Create(alternativeTestFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", alternativeTestFileName)
	}

	defer alternativeTestFile.Close()

	logWriter.Lookup[testName] = testFile
	logWriter.Lookup[alternativeTestName] = alternativeTestFile

	randomString := random.UniqueID()

	err = logWriter.WriteLog(logger, testName, randomString)
	require.NoError(t, err)

	alternativeRandomString := random.UniqueID()

	err = logWriter.WriteLog(logger, alternativeTestName, alternativeRandomString)
	require.NoError(t, err)

	buf, err := os.ReadFile(testFileName)
	require.NoError(t, err)
	assert.Equal(t, randomString+"\n", string(buf))

	buf, err = os.ReadFile(alternativeTestFileName)
	require.NoError(t, err)
	assert.Equal(t, alternativeRandomString+"\n", string(buf))
}

func TestWriteLogCreatesLogFileIfNotExists(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)

	logger := NewTestLogger(t)
	testName := t.Name()
	testFileName := filepath.Join(logWriter.OutputDir, testName+".log")

	randomString := random.UniqueID()

	err := logWriter.WriteLog(logger, testName, randomString)
	require.NoError(t, err)

	assert.True(t, files.FileExists(testFileName))

	buf, err := os.ReadFile(testFileName)
	require.NoError(t, err)
	assert.Equal(t, randomString+"\n", string(buf))
}
