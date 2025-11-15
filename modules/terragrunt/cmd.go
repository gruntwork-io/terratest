package terragrunt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// runTerragruntStackCommandE executes tg stack commands
// It handles argument construction, retry logic, and error handling for all stack commands
func runTerragruntStackCommandE(
	t testing.TestingT, opts *Options, subCommand string, additionalArgs ...string) (string, error) {
	// Build the base command arguments starting with "stack"
	commandArgs := []string{"stack"}
	if subCommand != "" {
		commandArgs = append(commandArgs, subCommand)
	}

	return executeTerragruntCommand(t, opts, commandArgs, additionalArgs...)
}

// runTerragruntCommandE is the core function that executes regular tg commands
// It handles argument construction, retry logic, and error handling for non-stack commands
func runTerragruntCommandE(t testing.TestingT, opts *Options, command string,
	additionalArgs ...string) (string, error) {
	// Build the base command arguments starting with the command
	commandArgs := []string{command}

	return executeTerragruntCommand(t, opts, commandArgs, additionalArgs...)
}

// executeTerragruntCommand is the common execution function for all tg commands
// It handles validation, argument construction, retry logic, and error handling
func executeTerragruntCommand(t testing.TestingT, opts *Options, baseCommandArgs []string,
	additionalArgs ...string) (string, error) {
	// Validate and prepare options
	if err := prepareOptions(opts); err != nil {
		return "", err
	}

	// Build args and generate command
	finalArgs := buildTerragruntArgs(opts, append(baseCommandArgs, additionalArgs...)...)
	execCommand := generateCommand(opts, finalArgs...)
	commandDescription := fmt.Sprintf("%s %v", opts.TerragruntBinary, finalArgs)

	// Execute the command with retry logic and error handling
	return retry.DoWithRetryableErrorsE(
		t,
		commandDescription,
		opts.RetryableTerraformErrors,
		opts.MaxRetries,
		opts.TimeBetweenRetries,
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

// prepareOptions validates options and sets defaults
func prepareOptions(opts *Options) error {
	if err := validateOptions(opts); err != nil {
		return err
	}
	if opts.TerragruntBinary == "" {
		opts.TerragruntBinary = DefaultTerragruntBinary
	}
	setTerragruntLogFormatting(opts)
	return nil
}

// buildTerragruntArgs constructs the final argument list for a terragrunt command
// Arguments are ordered as: TerragruntArgs → --non-interactive → commandArgs → TerraformArgs
func buildTerragruntArgs(opts *Options, commandArgs ...string) []string {
	var args []string
	args = append(args, opts.TerragruntArgs...)
	args = append(args, NonInteractiveFlag)
	args = append(args, commandArgs...)

	if len(opts.TerraformArgs) > 0 {
		args = append(args, opts.TerraformArgs...)
	}

	return args
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

// getExitCodeForTerragruntCommandE runs terragrunt with the given arguments and options and returns exit code
func getExitCodeForTerragruntCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (int, error) {
	// Validate and prepare options
	if err := prepareOptions(additionalOptions); err != nil {
		return defaultErrorExitCode, err
	}

	// Build args and generate command
	args := buildTerragruntArgs(additionalOptions, additionalArgs...)
	additionalOptions.Logger.Logf(t, "Running terragrunt with args %v", args)
	cmd := generateCommand(additionalOptions, args...)
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
