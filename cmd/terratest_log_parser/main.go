// A CLI command to parse parallel terratest output to produce test summaries and break out interleaved test output.
//
// This command will take as input a terratest log output from either stdin (through a pipe) or from a file, and output
// to a directory the following files:
// outputDir
//   |-> TEST_NAME.log
//   |-> summary.log
//   |-> report.xml
// where:
// - `TEST_NAME.log` is a log for each test run that only includes the relevant logs for that test.
// - `summary.log` is a summary of all the tests in the suite, including PASS/FAIL information.
// - `report.xml` is the test summary in junit XML format to be consumed by a CI engine.
//
// Certain tradeoffs were made in the decision to implement this functionality as a separate parsing command, as opposed
// to being built into the logger module as part of `Logf`. Specifically, this implementation avoids the difficulties of
// hooking into go's testing framework to be able to extract the summary logs, at the expense of a more complicated
// implementation in handling various corner cases due to logging flexibility. Here are the list of pros and cons of
// this approach that were considered:
//
// Pros:
// - Robust to unexpected failures in testing code, like `ctrl+c`, panics, OOM kills and the like since the parser is
//   not tied to the testing process. This approach is less likely to miss these entries, and can be surfaced to the
//   summary view for easy viewing in CI engines (no need to scroll), like the panic example.
// - Can combine `go test` output (e.g `--- PASS` entries) with the log entries for the test in a single log file.
// - Can extract out the summary view (those are all `go test` log entries).
// - Can build `junit.xml` report that CI engines can use for test insights.
//
// Cons:
// - Complicated implementation that is potentially brittle. E.g if someone decides to change the logging format then
//   this will all break. If we hook during the test, then the implementation is easier because those logs are all emitted
//   at certain points in code, the information of which is lost in the final log and have to parse out.
// - Have to store all the logs twice (the full interleaved version, and the broken out version) because the parsing
//   depends on logs being available. (NOTE: this is avoidable with a pipe).

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/go-commons/entrypoint"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/go-commons/logging"
	"github.com/gruntwork-io/terratest/modules/logger/parser"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var logger = logging.GetLogger("terratest_log_parser")

const CUSTOM_USAGE_TEXT = `Usage: terratest_log_parser [--help] [--log-level=info] [--testlog=LOG_INPUT] [--outputdir=OUTPUT_DIR]

A tool for parsing parallel terratest output to produce a test summary and to break out the interleaved logs by test for better debuggability.

Options:
   --log-level LEVEL  Set the log level to LEVEL. Must be one of: [panic fatal error warning info debug]
                      (default: "info")
   --testlog value    Path to file containing test log. If unset will use stdin.
   --outputdir value  Path to directory to output test output to. If unset will use the current directory.
   --help, -h         show help
`

func run(cliContext *cli.Context) error {
	filename := cliContext.String("testlog")
	outputDir := cliContext.String("outputdir")
	logLevel := cliContext.String("log-level")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return errors.WithStackTrace(err)
	}
	logger.SetLevel(level)

	var file *os.File
	if filename != "" {
		logger.Infof("reading from file")
		file, err = os.Open(filename)
		if err != nil {
			logger.Fatalf("Error opening file: %s", err)
		}
	} else {
		logger.Infof("reading from stdin")
		file = os.Stdin
	}
	defer file.Close()

	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		logger.Fatalf("Error extracting absolute path of output directory: %s", err)
	}

	parser.SpawnParsers(logger, file, outputDir)
	return nil
}

func main() {
	app := entrypoint.NewApp()
	cli.AppHelpTemplate = CUSTOM_USAGE_TEXT
	entrypoint.HelpTextLineWidth = 120

	app.Name = "terratest_log_parser"
	app.Author = "Gruntwork <www.gruntwork.io>"
	app.Description = `A tool for parsing parallel terratest output to produce a test summary and to break out the interleaved logs by test for better debuggability.`
	app.Action = run

	currentDir, err := os.Getwd()
	if err != nil {
		logger.Fatalf("Error finding current directory: %s", err)
	}
	defaultOutputDir := filepath.Join(currentDir, "out")

	logInputFlag := cli.StringFlag{
		Name:  "testlog, l",
		Value: "",
		Usage: "Path to file containing test log. If unset will use stdin.",
	}
	outputDirFlag := cli.StringFlag{
		Name:  "outputdir, o",
		Value: defaultOutputDir,
		Usage: "Path to directory to output test output to. If unset will use the current directory.",
	}
	logLevelFlag := cli.StringFlag{
		Name:  "log-level",
		Value: logrus.InfoLevel.String(),
		Usage: fmt.Sprintf("Set the log level to `LEVEL`. Must be one of: %v", logrus.AllLevels),
	}
	app.Flags = []cli.Flag{
		logLevelFlag,
		logInputFlag,
		outputDirFlag,
	}

	entrypoint.RunApp(app)
}
