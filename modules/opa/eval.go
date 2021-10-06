package opa

import (
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

	// Path to rego file containing the OPA rules.
	RulePath string

	// Set a logger that should be used. See the logger package for more info.
	Logger *logger.Logger

	// The following options can be used to change the behavior of the related functions for debuggability.

	// When true, keep any temp files and folders that are created for the purpose of running opa eval.
	DebugKeepTempFiles bool
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
//     opa eval -i $JSONFile -d $RulePath $ResultQuery
// This will asynchronously run OPA on each file concurrently using goroutines.
// This will fail the test if any one of the files failed.
func Eval(t testing.TestingT, options *EvalOptions, jsonFilePaths []string, resultQuery string) {
	require.NoError(t, EvalE(t, options, jsonFilePaths, resultQuery))
}

// EvalE runs `opa eval` on the given JSON files using the configured policy file and result query. Translates to:
//     opa eval -i $JSONFile -d $RulePath $ResultQuery
// This will asynchronously run OPA on each file concurrently using goroutines.
func EvalE(t testing.TestingT, options *EvalOptions, jsonFilePaths []string, resultQuery string) error {
	wg := new(sync.WaitGroup)
	wg.Add(len(jsonFilePaths))
	errorsOccurred := new(multierror.Error)
	errChans := make([]chan error, len(jsonFilePaths))
	for i, jsonFilePath := range jsonFilePaths {
		errChan := make(chan error, 1)
		errChans[i] = errChan
		go asyncEval(t, wg, errChan, options, jsonFilePath, resultQuery)
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
	jsonFilePath string,
	resultQuery string,
) {
	defer wg.Done()
	cmd := shell.Command{
		Command: "opa",
		Args:    formatOPAEvalArgs(options, jsonFilePath, resultQuery),
		Logger:  options.Logger,
	}
	err := shell.RunCommandE(t, cmd)
	if err == nil {
		options.Logger.Logf(t, "opa eval passed on file %s", jsonFilePath)
	} else {
		options.Logger.Logf(t, "Failed opa eval on file %s", jsonFilePath)
	}
	errChan <- err
}

// formatOPAEvalArgs formats the arguments for the `opa eval` command.
func formatOPAEvalArgs(options *EvalOptions, jsonFilePath string, resultQuery string) []string {
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
			"-d", options.RulePath,
			resultQuery,
		}...,
	)
	return args
}
