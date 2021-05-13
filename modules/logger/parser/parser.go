// Package logger/parser contains methods to parse and restructure log output from go testing and terratest
package parser

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	junitparser "github.com/jstemmer/go-junit-report/parser"
	"github.com/sirupsen/logrus"
)

// SpawnParsers will spawn the log parser and junit report parsers off of a single reader.
func SpawnParsers(logger *logrus.Logger, reader io.Reader, outputDir string) {
	forkedReader, forkedWriter := io.Pipe()
	teedReader := io.TeeReader(reader, forkedWriter)
	var waitForParsers sync.WaitGroup
	waitForParsers.Add(2)
	go func() {
		// close pipe writer, because this section drains the tee reader indicating reader is done draining
		defer forkedWriter.Close()
		defer waitForParsers.Done()
		parseAndStoreTestOutput(logger, teedReader, outputDir)
	}()
	go func() {
		defer waitForParsers.Done()
		report, err := junitparser.Parse(forkedReader, "")
		if err == nil {
			storeJunitReport(logger, outputDir, report)
		} else {
			logger.Errorf("Error parsing test output into junit report: %s", err)
		}
	}()
	waitForParsers.Wait()
}

// RegEx for parsing test status lines. Pulled from jstemmer/go-junit-report
var (
	regexResult  = regexp.MustCompile(`--- (PASS|FAIL|SKIP): (.+) \((\d+\.\d+)(?: ?seconds|s)\)`)
	regexStatus  = regexp.MustCompile(`=== (RUN|PAUSE|CONT)\s+(.+)`)
	regexSummary = regexp.MustCompile(`(^FAIL$)|(^(ok|FAIL)\s+([^ ]+)\s+(?:(\d+\.\d+)s|\(cached\)|(\[\w+ failed]))(?:\s+coverage:\s+(\d+\.\d+)%\sof\sstatements(?:\sin\s.+)?)?$)`)
	regexPanic   = regexp.MustCompile(`^panic:`)
)

// getIndent takes a line and returns the indent string
// Example:
//   in:  "    --- FAIL: TestSnafu"
//   out: "    "
func getIndent(data string) string {
	re := regexp.MustCompile(`^\s+`)
	indent := re.FindString(data)
	return indent
}

// getTestNameFromResultLine takes a go testing result line and extracts out the test name
// Example:
//   in:  --- FAIL: TestSnafu
//   out: TestSnafu
func getTestNameFromResultLine(text string) string {
	m := regexResult.FindStringSubmatch(text)
	return m[2]
}

// isResultLine checks if a line of text matches a test result (begins with "--- FAIL" or "--- PASS")
func isResultLine(text string) bool {
	return regexResult.MatchString(text)
}

// getTestNameFromStatusLine takes a go testing status line and extracts out the test name
// Example:
//   in:  === RUN  TestSnafu
//   out: TestSnafu
func getTestNameFromStatusLine(text string) string {
	m := regexStatus.FindStringSubmatch(text)
	return m[2]
}

// isStatusLine checks if a line of text matches a test status
func isStatusLine(text string) bool {
	return regexStatus.MatchString(text)
}

// isSummaryLine checks if a line of text matches the test summary
func isSummaryLine(text string) bool {
	return regexSummary.MatchString(text)
}

// isPanicLine checks if a line of text matches a panic
func isPanicLine(text string) bool {
	return regexPanic.MatchString(text)
}

// parseAndStoreTestOutput will take test log entries from terratest and aggregate the output by test. Takes advantage
// of the fact that terratest logs are prefixed by the test name. This will store the broken out logs into files under
// the outputDir, named by test name.
// Additionally will take test result lines and collect them under a summary log file named `summary.log`.
// See the `fixtures` directory for some examples.
func parseAndStoreTestOutput(
	logger *logrus.Logger,
	read io.Reader,
	outputDir string,
) {
	logWriter := LogWriter{
		lookup:    make(map[string]*os.File),
		outputDir: outputDir,
	}
	defer logWriter.closeFiles(logger)

	// Track some state that persists across lines
	testResultMarkers := TestResultMarkerStack{}
	previousTestName := ""

	var err error
	reader := bufio.NewReader(read)
	for {
		var data string
		data, err = reader.ReadString('\n')
		if len(data) == 0 && err == io.EOF {
			break
		}

		data = strings.TrimSuffix(data, "\n")

		// separate block so that we do not overwrite the err variable that we need afterwards to check if we're done
		{
			indentLevel := len(getIndent(data))
			isIndented := indentLevel > 0

			// Garbage collection of test result markers. Primary purpose is to detect when we dedent out, which can only be
			// detected when we reach a dedented line.
			testResultMarkers = testResultMarkers.removeDedentedTestResultMarkers(indentLevel)

			// Handle each possible category of test lines
			switch {
			case isSummaryLine(data):
				logWriter.writeLog(logger, "summary", data)

			case isStatusLine(data):
				testName := getTestNameFromStatusLine(data)
				previousTestName = testName
				logWriter.writeLog(logger, testName, data)

			case strings.HasPrefix(data, "Test"):
				// Heuristic: `go test` will only execute test functions named `Test.*`, so we assume any line prefixed
				// with `Test` is a test output for a named test. Also assume that test output will be space delimeted and
				// test names can't contain spaces (because they are function names).
				// This must be modified when `logger.DoLog` changes.
				vals := strings.Split(data, " ")
				testName := vals[0]
				previousTestName = testName
				logWriter.writeLog(logger, testName, data)

			case isIndented && isResultLine(data):
				// In a nested test result block, so collect the line into all the test results we have seen so far.
				for _, marker := range testResultMarkers {
					logWriter.writeLog(logger, marker.TestName, data)
				}

			case isPanicLine(data):
				// When panic, we want all subsequent nonstandard test lines to roll up to the summary
				previousTestName = "summary"
				logWriter.writeLog(logger, "summary", data)

			case isResultLine(data):
				// We ignore result lines, because that is handled specially below.

			case previousTestName != "":
				// Base case: roll up to the previous test line, if it exists.
				// Handles case where terratest log has entries with newlines in them.
				logWriter.writeLog(logger, previousTestName, data)

			default:
				logger.Warnf("Found test line that does not match known cases: %s", data)
			}

			// This has to happen separately from main if block to handle the special case of nested tests (e.g table driven
			// tests). For those result lines, we want it to roll up to the parent test, so we need to run the handler in
			// the `isIndented` section. But for both root and indented result lines, we want to execute the following code,
			// hence this special block.
			if isResultLine(data) {
				testName := getTestNameFromResultLine(data)
				logWriter.writeLog(logger, testName, data)
				logWriter.writeLog(logger, "summary", data)

				marker := TestResultMarker{
					TestName:    testName,
					IndentLevel: indentLevel,
				}
				testResultMarkers = testResultMarkers.push(marker)
			}
		}

		if err != nil {
			break
		}
	}
	if err != io.EOF {
		logger.Fatalf("Error reading from Reader: %s", err)
	}
}
