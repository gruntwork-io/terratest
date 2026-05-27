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

// parserCount is the number of concurrent parsers spawned by SpawnParsers.
const parserCount = 2

// SpawnParsers will spawn the log parser and junit report parsers off of a single reader.
func SpawnParsers(logger *logrus.Logger, reader io.Reader, outputDir string) {
	forkedReader, forkedWriter := io.Pipe()
	teedReader := io.TeeReader(reader, forkedWriter)

	var waitForParsers sync.WaitGroup

	waitForParsers.Add(parserCount)

	go func() {
		// close pipe writer, because this section drains the tee reader indicating reader is done draining
		defer func() {
			if err := forkedWriter.Close(); err != nil {
				logger.Errorf("Error closing forked writer: %s", err)
			}
		}()
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

// GetIndent takes a line and returns the indent string
// Example:
//
//	in:  "    --- FAIL: TestSnafu"
//	out: "    "
func GetIndent(data string) string {
	re := regexp.MustCompile(`^\s+`)

	return re.FindString(data)
}

// GetTestNameFromResultLine takes a go testing result line and extracts out the test name
// Example:
//
//	in:  --- FAIL: TestSnafu
//	out: TestSnafu
func GetTestNameFromResultLine(text string) string {
	m := regexResult.FindStringSubmatch(text)
	return m[2]
}

// IsResultLine checks if a line of text matches a test result (begins with "--- FAIL" or "--- PASS")
func IsResultLine(text string) bool {
	return regexResult.MatchString(text)
}

// GetTestNameFromStatusLine takes a go testing status line and extracts out the test name
// Example:
//
//	in:  === RUN  TestSnafu
//	out: TestSnafu
func GetTestNameFromStatusLine(text string) string {
	m := regexStatus.FindStringSubmatch(text)
	return m[2]
}

// IsStatusLine checks if a line of text matches a test status
func IsStatusLine(text string) bool {
	return regexStatus.MatchString(text)
}

// IsSummaryLine checks if a line of text matches the test summary
func IsSummaryLine(text string) bool {
	return regexSummary.MatchString(text)
}

// IsPanicLine checks if a line of text matches a panic
func IsPanicLine(text string) bool {
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
		Lookup:    make(map[string]*os.File),
		OutputDir: outputDir,
	}
	defer logWriter.CloseFiles(logger)

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
			indentLevel := len(GetIndent(data))
			isIndented := indentLevel > 0

			// Garbage collection of test result markers. Primary purpose is to detect when we dedent out, which can only be
			// detected when we reach a dedented line.
			testResultMarkers = testResultMarkers.RemoveDedentedTestResultMarkers(indentLevel)

			// Handle each possible category of test lines
			switch {
			case IsSummaryLine(data):
				if writeErr := logWriter.WriteLog(logger, "summary", data); writeErr != nil {
					logger.Errorf("Error writing summary log: %s", writeErr)
				}

			case IsStatusLine(data):
				testName := GetTestNameFromStatusLine(data)
				previousTestName = testName

				if writeErr := logWriter.WriteLog(logger, testName, data); writeErr != nil {
					logger.Errorf("Error writing log for test %s: %s", testName, writeErr)
				}

			case strings.HasPrefix(data, "Test"):
				// Heuristic: `go test` will only execute test functions named `Test.*`, so we assume any line prefixed
				// with `Test` is a test output for a named test. Also assume that test output will be space delimeted and
				// test names can't contain spaces (because they are function names).
				// This must be modified when `logger.DoLog` changes.
				vals := strings.Split(data, " ")
				testName := vals[0]
				previousTestName = testName

				if writeErr := logWriter.WriteLog(logger, testName, data); writeErr != nil {
					logger.Errorf("Error writing log for test %s: %s", testName, writeErr)
				}

			case isIndented && IsResultLine(data):
				// In a nested test result block, so collect the line into all the test results we have seen so far.
				for _, marker := range testResultMarkers {
					if writeErr := logWriter.WriteLog(logger, marker.TestName, data); writeErr != nil {
						logger.Errorf("Error writing log for test %s: %s", marker.TestName, writeErr)
					}
				}

			case IsPanicLine(data):
				// When panic, we want all subsequent nonstandard test lines to roll up to the summary
				previousTestName = "summary"

				if writeErr := logWriter.WriteLog(logger, "summary", data); writeErr != nil {
					logger.Errorf("Error writing summary log: %s", writeErr)
				}

			case IsResultLine(data):
				// We ignore result lines, because that is handled specially below.

			case previousTestName != "":
				// Base case: roll up to the previous test line, if it exists.
				// Handles case where terratest log has entries with newlines in them.
				if writeErr := logWriter.WriteLog(logger, previousTestName, data); writeErr != nil {
					logger.Errorf("Error writing log for test %s: %s", previousTestName, writeErr)
				}

			default:
				logger.Warnf("Found test line that does not match known cases: %s", data)
			}

			// This has to happen separately from main if block to handle the special case of nested tests (e.g table driven
			// tests). For those result lines, we want it to roll up to the parent test, so we need to run the handler in
			// the `isIndented` section. But for both root and indented result lines, we want to execute the following code,
			// hence this special block.
			if IsResultLine(data) {
				testName := GetTestNameFromResultLine(data)

				if writeErr := logWriter.WriteLog(logger, testName, data); writeErr != nil {
					logger.Errorf("Error writing log for test %s: %s", testName, writeErr)
				}

				if writeErr := logWriter.WriteLog(logger, "summary", data); writeErr != nil {
					logger.Errorf("Error writing summary log: %s", writeErr)
				}

				marker := TestResultMarker{
					TestName:    testName,
					IndentLevel: indentLevel,
				}
				testResultMarkers = testResultMarkers.Push(marker)
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
