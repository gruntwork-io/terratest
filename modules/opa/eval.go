package opa

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

// EvalOptions defines options that can be passed to the 'opa eval' command for checking policies on arbitrary JSON data
// via OPA.
type EvalOptions struct {
	// Whether OPA should run checks with failure.
	FailMode FailMode

	// Path to rego file containing the OPA rules. Can also be a remote path defined in go-getter syntax. Refer to
	// https://github.com/hashicorp/go-getter#url-format for supported options.
	RulePath string

	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger

	// The following options can be used to change the behavior of the related functions for debuggability.

	// When true, keep any temp files and folders that are created for the purpose of running opa eval.
	DebugKeepTempFiles bool

	// When true, disable the functionality where terratest reruns the opa check on the same file and query all elements
	// on error. By default, terratest will rerun the opa eval call with `data` query so you can see all the contents
	// evaluated.
	DebugDisableQueryDataOnError bool
}

// FailMode signals whether `opa eval` should fail when the query returns an undefined value (FailUndefined), a
// defined value (FailDefined), or not at all (NoFail).
type FailMode int

const (
	FailUndefined FailMode = iota
	FailDefined
	NoFail
)

// EvalE runs `opa eval` on the given JSON files using the configured policy file and result query. Translates to:
//
//	opa eval -i $JSONFile -d $RulePath $ResultQuery
//
// This will asynchronously run OPA on each file concurrently using goroutines.
// This will fail the test if any one of the files failed.
func Eval(t testing.TestingT, options *EvalOptions, jsonFilePaths []string, resultQuery string) {
	require.NoError(t, EvalE(t, options, jsonFilePaths, resultQuery))
}

// EvalE runs `opa eval` on the given JSON files using the configured policy file and result query. Translates to:
//
//	opa eval -i $JSONFile -d $RulePath $ResultQuery
//
// This will asynchronously run OPA on each file concurrently using goroutines.
func EvalE(t testing.TestingT, options *EvalOptions, jsonFilePaths []string, resultQuery string) error {
	downloadedPolicyPath, err := DownloadPolicyE(t, options.RulePath)
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(jsonFilePaths))
	errorsOccurred := new(multierror.Error)
	errChans := make([]chan error, len(jsonFilePaths))
	for i, jsonFilePath := range jsonFilePaths {
		errChan := make(chan error, 1)
		errChans[i] = errChan
		go asyncEval(t, wg, errChan, options, downloadedPolicyPath, jsonFilePath, resultQuery)
	}
	wg.Wait()
	for _, errChan := range errChans {
		err := <-errChan
		if err != nil {
			errorsOccurred = multierror.Append(errorsOccurred, err)
		}
	}
	return errorsOccurred.ErrorOrNil()
}

// asyncEval is a function designed to be run in a goroutine to asynchronously call `opa eval` on a single input file.
func asyncEval(
	t testing.TestingT,
	wg *sync.WaitGroup,
	errChan chan error,
	options *EvalOptions,
	downloadedPolicyPath string,
	jsonFilePath string,
	resultQuery string,
) {
	defer wg.Done()
	cmd := shell.Command{
		Command: "opa",
		Args:    formatOPAEvalArgs(options, downloadedPolicyPath, jsonFilePath, resultQuery),

		// Do not log output from shell package so we can log the full json without breaking it up. This is ok, because
		// opa eval is typically very quick.
		Logger: logger.Discard,
	}
	err := runCommandWithFullLoggingE(t, options.Logger, cmd)
	ruleBasePath := filepath.Base(downloadedPolicyPath)
	if err == nil {
		options.Logger.Logf(t, "opa eval passed on file %s (policy %s; query %s)", jsonFilePath, ruleBasePath, resultQuery)
	} else {
		options.Logger.Logf(t, "Failed opa eval on file %s (policy %s; query %s)", jsonFilePath, ruleBasePath, resultQuery)
		if options.DebugDisableQueryDataOnError == false {
			options.Logger.Logf(t, "DEBUG: rerunning opa eval to query for full data.")
			cmd.Args = formatOPAEvalArgs(options, downloadedPolicyPath, jsonFilePath, "data")
			// We deliberately ignore the error here as we want to only return the original error.
			runCommandWithFullLoggingE(t, options.Logger, cmd)
		}
	}
	errChan <- err
}

// formatOPAEvalArgs formats the arguments for the `opa eval` command.
func formatOPAEvalArgs(options *EvalOptions, rulePath, jsonFilePath, resultQuery string) []string {
	args := []string{"eval"}

	switch options.FailMode {
	case FailUndefined:
		args = append(args, "--fail")
	case FailDefined:
		args = append(args, "--fail-defined")
	}

	args = append(
		args,
		[]string{
			"-i", jsonFilePath,
			"-d", rulePath,
			resultQuery,
		}...,
	)
	return args
}

// runCommandWithFullLogging will log the command output in its entirety with buffering. This avoids breaking up the
// logs when commands are run concurrently. This is a private function used in the context of opa only because opa runs
// very quickly, and the output of opa is hard to parse if it is broken up by interleaved logs.
func runCommandWithFullLoggingE(t testing.TestingT, logger *logger.Logger, cmd shell.Command) error {
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	logger.Logf(t, "Output of command `%s %s`:\n%s", cmd.Command, strings.Join(cmd.Args, " "), output)
	return err
}
