package parser

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/go-commons/files"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func createLogWriter(t *testing.T) LogWriter {
	dir := getTempDir(t)
	logWriter := LogWriter{
		lookup:    make(map[string]*os.File),
		outputDir: dir,
	}
	return logWriter
}

func TestEnsureDirectoryExistsCreatesDirectory(t *testing.T) {
	t.Parallel()

	dir := getTempDir(t)
	defer os.RemoveAll(dir)

	logger := NewTestLogger(t)
	tmpd := filepath.Join(dir, "tmpdir")
	assert.False(t, files.IsDir(tmpd))
	ensureDirectoryExists(logger, tmpd)
	assert.True(t, files.IsDir(tmpd))
}

func TestEnsureDirectoryExistsHandlesExistingDirectory(t *testing.T) {
	t.Parallel()

	dir := getTempDir(t)
	defer os.RemoveAll(dir)

	logger := NewTestLogger(t)
	assert.True(t, files.IsDir(dir))
	ensureDirectoryExists(logger, dir)
	assert.True(t, files.IsDir(dir))
}

func TestGetOrCreateFileCreatesNewFile(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	testFileName := filepath.Join(logWriter.outputDir, t.Name()+".log")
	assert.False(t, files.FileExists(testFileName))
	file, err := logWriter.getOrCreateFile(logger, t.Name())
	defer file.Close()
	assert.Nil(t, err)
	assert.NotNil(t, file)
	assert.True(t, files.FileExists(testFileName))
}

func TestGetOrCreateFileCreatesNewFileIfTestNameHasDir(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	dirName := filepath.Join(logWriter.outputDir, "TestMain")
	testFileName := filepath.Join(dirName, t.Name()+".log")
	assert.False(t, files.IsDir(dirName))
	assert.False(t, files.FileExists(testFileName))
	file, err := logWriter.getOrCreateFile(logger, filepath.Join("TestMain", t.Name()))
	defer file.Close()
	assert.Nil(t, err)
	assert.NotNil(t, file)
	assert.True(t, files.IsDir(dirName))
	assert.True(t, files.FileExists(testFileName))
}

func TestGetOrCreateChannelReturnsExistingFileHandle(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	testName := t.Name()
	logger := NewTestLogger(t)
	testFileName := filepath.Join(logWriter.outputDir, t.Name())
	file, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", testFileName)
	}
	defer file.Close()

	logWriter.lookup[testName] = file
	lookupFile, err := logWriter.getOrCreateFile(logger, testName)
	assert.Nil(t, err)
	assert.Equal(t, lookupFile, file)
}

func TestCloseFilesClosesAll(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	testName := t.Name()
	testFileName := filepath.Join(logWriter.outputDir, testName)
	testFile, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", testFileName)
	}
	alternativeTestName := t.Name() + "Alternative"
	alternativeTestFileName := filepath.Join(logWriter.outputDir, alternativeTestName)
	alternativeTestFile, err := os.Create(alternativeTestFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", alternativeTestFileName)
	}
	logWriter.lookup[testName] = testFile
	logWriter.lookup[alternativeTestName] = alternativeTestFile

	logWriter.closeFiles(logger)
	err = testFile.Close()
	assert.Contains(t, err.Error(), os.ErrClosed.Error())
	err = alternativeTestFile.Close()
	assert.Contains(t, err.Error(), os.ErrClosed.Error())
}

func TestWriteLogWritesToCorrectLogFile(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	testName := t.Name()
	testFileName := filepath.Join(logWriter.outputDir, testName)
	testFile, err := os.Create(testFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", testFileName)
	}
	defer testFile.Close()
	alternativeTestName := t.Name() + "Alternative"
	alternativeTestFileName := filepath.Join(logWriter.outputDir, alternativeTestName)
	alternativeTestFile, err := os.Create(alternativeTestFileName)
	if err != nil {
		t.Fatalf("error creating test file %s", alternativeTestFileName)
	}
	defer alternativeTestFile.Close()
	logWriter.lookup[testName] = testFile
	logWriter.lookup[alternativeTestName] = alternativeTestFile

	randomString := random.UniqueId()
	err = logWriter.writeLog(logger, testName, randomString)
	assert.Nil(t, err)
	alternativeRandomString := random.UniqueId()
	err = logWriter.writeLog(logger, alternativeTestName, alternativeRandomString)
	assert.Nil(t, err)

	buf, err := ioutil.ReadFile(testFileName)
	assert.Nil(t, err)
	assert.Equal(t, string(buf), randomString+"\n")
	buf, err = ioutil.ReadFile(alternativeTestFileName)
	assert.Nil(t, err)
	assert.Equal(t, string(buf), alternativeRandomString+"\n")
}

func TestWriteLogCreatesLogFileIfNotExists(t *testing.T) {
	t.Parallel()

	logWriter := createLogWriter(t)
	defer os.RemoveAll(logWriter.outputDir)

	logger := NewTestLogger(t)
	testName := t.Name()
	testFileName := filepath.Join(logWriter.outputDir, testName+".log")

	randomString := random.UniqueId()
	err := logWriter.writeLog(logger, testName, randomString)
	assert.Nil(t, err)

	assert.True(t, files.FileExists(testFileName))
	buf, err := ioutil.ReadFile(testFileName)
	assert.Nil(t, err)
	assert.Equal(t, string(buf), randomString+"\n")
}
