// Package logger/parser contains methods to parse and restructure log output from go testing and terratest
package parser

import (
	"os"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/go-commons/files"
	junitformatter "github.com/jstemmer/go-junit-report/formatter"
	junitparser "github.com/jstemmer/go-junit-report/parser"
	"github.com/sirupsen/logrus"
)

type LogWriter struct {
	// Represents an open file to a log corresponding to a test (key = test name)
	lookup    map[string]*os.File
	outputDir string
}

// LogWriter.getOrCreateFile will get the corresponding file to a log for the provided test name, or create a new file.
func (logWriter LogWriter) getOrCreateFile(logger *logrus.Logger, testName string) (*os.File, error) {
	file, hasKey := logWriter.lookup[testName]
	if hasKey {
		return file, nil
	}

	filename := filepath.Join(logWriter.outputDir, testName+".log")
	file, err := createLogFile(logger, filename)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}
	logWriter.lookup[testName] = file
	return file, nil
}

// LogWriter.closeChannels closes all the channels in the lookup dictionary
func (logWriter LogWriter) closeFiles(logger *logrus.Logger) {
	logger.Infof("Closing all the files in log writer")
	for testName, file := range logWriter.lookup {
		err := file.Close()
		if err != nil {
			logger.Errorf("Error closing log file for test %s: %s", testName, err)
		}
	}
}

// writeLog will write the provided text to the corresponding log file for the provided test.
func (logWriter LogWriter) writeLog(logger *logrus.Logger, testName string, text string) error {
	file, err := logWriter.getOrCreateFile(logger, testName)
	if err != nil {
		logger.Errorf("Error retrieving log for test: %s", testName)
		return errors.WithStackTrace(err)
	}
	_, err = file.WriteString(text + "\n")
	if err != nil {
		logger.Errorf("Error (%s) writing log entry: %s", err, text)
		return errors.WithStackTrace(err)
	}
	file.Sync()
	return nil
}

// createLogFile will create and return the open file handle for the file at provided filename, creating all directories
// in the process.
func createLogFile(logger *logrus.Logger, filename string) (*os.File, error) {
	// We extract and create the directory for interpolated filename, to handle nested tests where testname contains '/'
	dirName := filepath.Dir(filename)
	err := ensureDirectoryExists(logger, dirName)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}
	file, err := os.Create(filename)
	if err != nil {
		return nil, errors.WithStackTrace(err)
	}
	return file, nil
}

// ensureDirectoryExists will only attempt to create the directory if it does not exist
func ensureDirectoryExists(logger *logrus.Logger, dirName string) error {
	if files.IsDir(dirName) {
		logger.Infof("Directory %s already exists", dirName)
		return nil
	}
	logger.Infof("Creating directory %s", dirName)
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		logger.Errorf("Error making directory %s: %s", dirName, err)
		return errors.WithStackTrace(err)
	}
	return nil
}

// storeJunitReport takes a parsed Junit report and stores it as report.xml in the output directory
func storeJunitReport(logger *logrus.Logger, outputDir string, report *junitparser.Report) {
	ensureDirectoryExists(logger, outputDir)
	filename := filepath.Join(outputDir, "report.xml")
	f, err := os.Create(filename)
	if err != nil {
		logger.Errorf("Error making file %s for junit report", filename)
		return
	}
	defer f.Close()

	err = junitformatter.JUnitReportXML(report, false, "", f)
	if err != nil {
		logger.Errorf("Error formatting junit xml report: %s", err)
		return
	}
}
