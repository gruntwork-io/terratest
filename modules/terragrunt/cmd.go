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

// runTerragruntStackCommandE executes a terragrunt stack command with "run" subcommand
// This is used for commands that need to run terragrunt stack run <command>
func runTerragruntStackCommandE(t testing.TestingT, opts *Options, additionalArgs ...string) (string, error) {
	return runTerragruntStackSubCommandE(t, opts, "run", additionalArgs...)
}

// terragruntStackCommandE executes a terragrunt stack command without any subcommand
// This is used for commands like "terragrunt stack generate"
func terragruntStackCommandE(t testing.TestingT, opts *Options, additionalArgs ...string) (string, error) {
	return runTerragruntStackSubCommandE(t, opts, "", additionalArgs...)
}

// runTerragruntStackSubCommandE is the core function that executes terragrunt stack commands
// It handles experimental flag detection, argument construction, and retry logic
func runTerragruntStackSubCommandE(t testing.TestingT, opts *Options, subCommand string, additionalArgs ...string) (string, error) {
	// Build the base command arguments starting with "stack"
	commandArgs := []string{"stack"}
	if subCommand != "" {
		commandArgs = append(commandArgs, subCommand)
	}

	// Check if the terragrunt binary supports experimental stack features
	// This is done by testing if "-experiment stack" is a valid flag combination
	experimentalTestCmd := shell.Command{
		Command: opts.TerragruntBinary,
		Args:    []string{"-experiment", "stack"},
	}
	if err := shell.RunCommandE(t, experimentalTestCmd); err == nil {
		// If experimental flag is supported, prepend it to the command arguments
		commandArgs = slices.Insert(commandArgs, 0, "-experiment", "stack")
	}

	// Apply common terragrunt options and get the final command arguments
	terragruntOptions, finalArgs := GetCommonOptions(opts, commandArgs...)

	// Append additional arguments with "--" separator
	finalArgs = append(finalArgs, slices.Insert(additionalArgs, 0, "--")...)

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
