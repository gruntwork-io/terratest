package terragrunt

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// terragruntStackCommandE executes a terragrunt stack command without any subcommand
// This is used for commands like "terragrunt stack generate"
func terragruntStackCommandE(t testing.TestingT, opts *Options, additionalArgs ...string) (string, error) {
	return runTerragruntStackCommandE(t, opts, "", additionalArgs...)
}

// runTerragruntStackCommandE is the unified function that executes terragrunt stack commands
// It handles argument construction, retry logic, and error handling for all stack commands
func runTerragruntStackCommandE(t testing.TestingT, opts *Options, subCommand string, additionalArgs ...string) (string, error) {
	// Validate required options
	if err := validateOptions(opts); err != nil {
		return "", err
	}

	// Build the base command arguments starting with "stack"
	commandArgs := []string{"stack"}
	if subCommand != "" {
		commandArgs = append(commandArgs, subCommand)
	}

	// Apply common terragrunt options and get the final command arguments
	terragruntOptions, finalArgs := GetCommonOptions(opts, commandArgs...)

	// Append additional arguments with "--" separator for stack commands
	// Exception: "output" subcommand doesn't use the separator
	if len(additionalArgs) > 0 {
		if subCommand == "output" {
			// For output commands, append arguments directly without separator
			finalArgs = append(finalArgs, additionalArgs...)
		} else {
			// For other stack commands, use the separator
			finalArgs = append(finalArgs, slices.Insert(additionalArgs, 0, ArgSeparator)...)
		}
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

// runTerragruntCommandE is the core function that executes regular terragrunt commands
// It handles argument construction, retry logic, and error handling for non-stack commands
func runTerragruntCommandE(t testing.TestingT, opts *Options, command string, additionalArgs ...string) (string, error) {
	// Validate required options
	if err := validateOptions(opts); err != nil {
		return "", err
	}

	// Build the base command arguments starting with the command
	commandArgs := []string{command}

	// Apply common terragrunt options and get the final command arguments
	terragruntOptions, finalArgs := GetCommonOptions(opts, commandArgs...)

	// Append additional arguments
	finalArgs = append(finalArgs, additionalArgs...)

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

// generateCommand creates a shell.Command with the specified terragrunt options and arguments
// This function encapsulates the command creation logic for consistency
func generateCommand(terragruntOptions *Options, commandArgs ...string) shell.Command {
	return shell.Command{
		Command:    terragruntOptions.TerragruntBinary,
		Args:       commandArgs,
		WorkingDir: terragruntOptions.TerragruntDir,
		Env:        terragruntOptions.EnvVars,
		Logger:     terragruntOptions.Logger,
	}
}
