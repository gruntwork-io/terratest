package terragrunt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// runTerragruntStackCommandE is the unified function that executes tg stack commands
// It handles argument construction, retry logic, and error handling for all stack commands
func runTerragruntStackCommandE(
	t testing.TestingT, opts *Options, subCommand string, additionalArgs ...string) (string, error) {
	// Default behavior: use arg separator (for backward compatibility)
	return runTerragruntStackCommandWithSeparatorE(t, opts, subCommand, true, additionalArgs...)
}

// runTerragruntStackCommandWithSeparatorE executes tg stack commands with control over the -- separator
// useArgSeparator controls whether the "--" separator is added before additional arguments
func runTerragruntStackCommandWithSeparatorE(t testing.TestingT, opts *Options,
	subCommand string, useArgSeparator bool, additionalArgs ...string) (string, error) {
	// Build the base command arguments starting with "stack"
	commandArgs := []string{"stack"}
	if subCommand != "" {
		commandArgs = append(commandArgs, subCommand)
	}

	return executeTerragruntCommand(t, opts, commandArgs, useArgSeparator, additionalArgs...)
}

// runTerragruntCommandE is the core function that executes regular tg commands
// It handles argument construction, retry logic, and error handling for non-stack commands
func runTerragruntCommandE(t testing.TestingT, opts *Options, command string,
	additionalArgs ...string) (string, error) {
	// Build the base command arguments starting with the command
	commandArgs := []string{command}

	// For non-stack commands, we typically don't use the separator
	return executeTerragruntCommand(t, opts, commandArgs, false, additionalArgs...)
}

// executeTerragruntCommand is the common execution function for all tg commands
// It handles validation, argument construction, retry logic, and error handling
func executeTerragruntCommand(t testing.TestingT, opts *Options, baseCommandArgs []string,
	useArgSeparator bool, additionalArgs ...string) (string, error) {
	// Validate required options
	if err := validateOptions(opts); err != nil {
		return "", err
	}

	// Apply common tg options and get the final command arguments
	terragruntOptions, finalArgs := GetCommonOptions(opts, baseCommandArgs...)

	// Append arguments from options using the new separation logic
	argsFromOptions := GetArgsForCommand(terragruntOptions, useArgSeparator)
	finalArgs = append(finalArgs, argsFromOptions...)

	// Append any additional arguments passed directly to this function
	if len(additionalArgs) > 0 {
		finalArgs = append(finalArgs, additionalArgs...)
	}

	// Generate the final shell command
	execCommand := generateCommand(terragruntOptions, finalArgs...)
	commandDescription := fmt.Sprintf("%s %v", terragruntOptions.TerragruntBinary, finalArgs)

	// Execute the command with retry logic and error handling
	return retry.DoWithRetryableErrorsE(
		t,
		commandDescription,
		terragruntOptions.RetryableTerraformErrors,
		terragruntOptions.MaxRetries,
		terragruntOptions.TimeBetweenRetries,
		func() (string, error) {
			output, err := shell.RunCommandAndGetOutputE(t, execCommand)
			if err != nil {
				return output, err
			}

			// Check for warnings that should be treated as errors
			if warningErr := hasWarning(opts, output); warningErr != nil {
				return output, warningErr
			}

			return output, nil
		},
	)
}

// hasWarning checks if the command output contains any warnings that should be treated as errors
// It uses regex patterns defined in opts.WarningsAsErrors to match warning messages
func hasWarning(opts *Options, commandOutput string) error {
	for warningPattern, errorMessage := range opts.WarningsAsErrors {
		// Create a regex pattern to match warnings with the specified pattern
		regexPattern := fmt.Sprintf("\nWarning: %s[^\n]*\n", warningPattern)
		compiledRegex, err := regexp.Compile(regexPattern)
		if err != nil {
			return fmt.Errorf("cannot compile regex for warning detection: %w", err)
		}

		// Find all matches of the warning pattern in the output
		matches := compiledRegex.FindAllString(commandOutput, -1)
		if len(matches) == 0 {
			continue
		}

		// If warnings are found, return an error with the specified message
		return fmt.Errorf("warning(s) were found: %s:\n%s", errorMessage, strings.Join(matches, ""))
	}
	return nil
}

// validateOptions validates that required options are provided
func validateOptions(opts *Options) error {
	if opts == nil {
		return fmt.Errorf("options cannot be nil")
	}
	if opts.TerragruntDir == "" {
		return fmt.Errorf("TerragruntDir is required")
	}
	return nil
}

// defaultSuccessExitCode is the exit code returned when terraform command succeeds
const defaultSuccessExitCode = 0

// defaultErrorExitCode is the exit code returned when terraform command fails
const defaultErrorExitCode = 1

// getExitCodeForTerraformCommandE runs terraform with the given arguments and options and returns exit code
func getExitCodeForTerraformCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (int, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	additionalOptions.Logger.Logf(t, "Running terragrunt with args %v", args)
	cmd := generateCommand(options, args...)
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err == nil {
		return defaultSuccessExitCode, nil
	}
	exitCode, getExitCodeErr := shell.GetExitCodeForRunCommandError(err)
	if getExitCodeErr == nil {
		return exitCode, nil
	}
	return defaultErrorExitCode, getExitCodeErr
}

// generateCommand creates a shell.Command with the specified tg options and arguments
// This function encapsulates the command creation logic for consistency
func generateCommand(terragruntOptions *Options, commandArgs ...string) shell.Command {
	return shell.Command{
		Command:    terragruntOptions.TerragruntBinary,
		Args:       commandArgs,
		WorkingDir: terragruntOptions.TerragruntDir,
		Env:        terragruntOptions.EnvVars,
		Logger:     terragruntOptions.Logger,
		Stdin:      terragruntOptions.Stdin,
	}
}
