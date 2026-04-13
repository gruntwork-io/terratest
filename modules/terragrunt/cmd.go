package terragrunt

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// runTerragruntStackCommandE executes terragrunt stack commands.
// It handles argument construction, retry logic, and error handling for all stack commands.
func runTerragruntStackCommandE(
	t testing.TestingT,
	ctx context.Context,
	opts *Options,
	subCommand string,
	additionalArgs ...string,
) (string, error) {
	// Build the base command arguments starting with "stack"
	commandArgs := []string{"stack"}
	if subCommand != "" {
		commandArgs = append(commandArgs, subCommand)
	}

	return executeTerragruntCommand(t, ctx, opts, commandArgs, additionalArgs...)
}

// runTerragruntCommandE is the core function that executes regular tg commands.
// It handles argument construction, retry logic, and error handling for non-stack commands.
func runTerragruntCommandE(
	t testing.TestingT,
	ctx context.Context,
	opts *Options,
	command string,
	additionalArgs ...string,
) (string, error) {
	// Build the base command arguments starting with the command
	commandArgs := []string{command}

	return executeTerragruntCommand(t, ctx, opts, commandArgs, additionalArgs...)
}

// executeTerragruntCommand is the common execution function for all tg commands.
// It handles validation, argument construction, retry logic, and error handling.
func executeTerragruntCommand(t testing.TestingT, ctx context.Context, opts *Options, baseCommandArgs []string,
	additionalArgs ...string) (string, error) {
	// Validate and prepare options
	if err := PrepareOptions(opts); err != nil {
		return "", err
	}

	// Build args and generate command
	finalArgs := BuildTerragruntArgs(opts, append(baseCommandArgs, additionalArgs...)...)
	execCommand := generateCommand(opts, finalArgs...)
	commandDescription := fmt.Sprintf("%s %v", opts.TerragruntBinary, finalArgs)

	// Execute the command with retry logic and error handling
	return retry.DoWithRetryableErrorsContextE(
		t,
		ctx,
		commandDescription,
		opts.RetryableTerraformErrors,
		opts.MaxRetries,
		opts.TimeBetweenRetries,
		func() (string, error) {
			output, err := shell.RunCommandContextAndGetOutputE(t, ctx, &execCommand)
			if err != nil {
				return output, err
			}

			// Check for warnings that should be treated as errors
			if warningErr := HasWarning(opts, output); warningErr != nil {
				return output, warningErr
			}

			return output, nil
		},
	)
}

// HasWarning checks if the command output contains any warnings that should be treated as errors.
// It uses regex patterns defined in opts.WarningsAsErrors to match warning messages.
func HasWarning(opts *Options, commandOutput string) error {
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

// PrepareOptions validates options and sets defaults.
func PrepareOptions(opts *Options) error {
	if err := ValidateOptions(opts); err != nil {
		return err
	}

	if opts.TerragruntBinary == "" {
		opts.TerragruntBinary = DefaultTerragruntBinary
	}

	setTerragruntLogFormatting(opts)

	return nil
}

// BuildTerragruntArgs constructs the final argument list for a terragrunt command.
// Arguments are ordered as: TerragruntArgs → --non-interactive → commandArgs → TerraformArgs.
func BuildTerragruntArgs(opts *Options, commandArgs ...string) []string {
	var args []string

	args = append(args, opts.TerragruntArgs...)
	args = append(args, NonInteractiveFlag)
	args = append(args, commandArgs...)

	if len(opts.TerraformArgs) > 0 {
		args = append(args, opts.TerraformArgs...)
	}

	return args
}

// ValidateOptions validates that required options are provided.
func ValidateOptions(opts *Options) error {
	if opts == nil {
		return ErrNilOptions
	}

	if opts.TerragruntDir == "" {
		return ErrMissingTerragruntDir
	}

	return nil
}

// defaultSuccessExitCode is the exit code returned when the OpenTofu/Terraform command succeeds
const defaultSuccessExitCode = 0

// defaultErrorExitCode is the exit code returned when the OpenTofu/Terraform command fails
const defaultErrorExitCode = 1

// getExitCodeForTerragruntCommandE runs terragrunt with the given arguments and options and returns exit code.
func getExitCodeForTerragruntCommandE(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) (int, error) {
	// Validate and prepare options
	if err := PrepareOptions(additionalOptions); err != nil {
		return defaultErrorExitCode, err
	}

	// Build args and generate command
	args := BuildTerragruntArgs(additionalOptions, additionalArgs...)
	additionalOptions.Logger.Logf(t, "Running terragrunt with args %v", args)
	cmd := generateCommand(additionalOptions, args...)

	_, err := shell.RunCommandContextAndGetOutputE(t, ctx, &cmd)
	if err == nil {
		return defaultSuccessExitCode, nil
	}

	exitCode, getExitCodeErr := shell.GetExitCodeForRunCommandError(err)
	if getExitCodeErr == nil {
		return exitCode, nil
	}

	return defaultErrorExitCode, getExitCodeErr
}

// BuildRunArgs constructs the argument list for a terragrunt run command.
// The -- separator disambiguates Terragrunt flags from OpenTofu/Terraform flags:
//
//	run [tgArgs...] -- [tfArgs...]
func BuildRunArgs(tgArgs []string, tfArgs []string) []string {
	args := make([]string, 0, len(tgArgs)+1+len(tfArgs))
	args = append(args, tgArgs...)
	args = append(args, "--")
	args = append(args, tfArgs...)

	return args
}

// generateCommand creates a shell.Command with the specified tg options and arguments.
// This function encapsulates the command creation logic for consistency.
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
