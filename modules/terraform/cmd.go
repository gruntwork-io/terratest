package terraform

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
)

func generateCommand(options *Options, args ...string) shell.Command {
	cmd := shell.Command{
		Command:    options.TerraformBinary,
		Args:       args,
		WorkingDir: options.TerraformDir,
		Env:        options.EnvVars,
		Logger:     options.Logger,
	}
	return cmd
}

var commandsWithParallelism = []string{
	"plan",
	"apply",
	"destroy",
	"plan-all",
	"apply-all",
	"destroy-all",
}

// GetCommonOptions extracts commons terraform options
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	if options.NoColor && !collections.ListContains(args, "-no-color") {
		args = append(args, "-no-color")
	}

	if options.TerraformBinary == "" {
		options.TerraformBinary = "terraform"
	}

	if options.TerraformBinary == "terragrunt" {
		args = append(args, "--terragrunt-non-interactive")
	}

	if options.Parallelism > 0 && len(args) > 0 && collections.ListContains(commandsWithParallelism, args[0]) {
		args = append(args, fmt.Sprintf("--parallelism=%d", options.Parallelism))
	}

	// if SshAgent is provided, override the local SSH agent with the socket of our in-process agent
	if options.SshAgent != nil {
		// Initialize EnvVars, if it hasn't been set yet
		if options.EnvVars == nil {
			options.EnvVars = map[string]string{}
		}
		options.EnvVars["SSH_AUTH_SOCK"] = options.SshAgent.SocketFile()
	}
	return options, args
}

// RunTerraformCommand runs terraform with the given arguments and options and return stdout/stderr.
func RunTerraformCommand(t testing.TestingT, additionalOptions *Options, args ...string) string {
	out, err := RunTerraformCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunTerraformCommandE runs terraform with the given arguments and options and return stdout/stderr.
func RunTerraformCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)
	return retry.DoWithRetryableErrorsE(t, description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandAndGetOutputE(t, cmd)
	})
}

// RunTerraformCommandAndGetStdoutE runs terraform with the given arguments and options and returns solely its stdout
// (but not stderr).
func RunTerraformCommandAndGetStdoutE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)
	return retry.DoWithRetryableErrorsE(t, description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		return shell.RunCommandAndGetStdOutE(t, cmd)
	})
}

// GetExitCodeForTerraformCommand runs terraform with the given arguments and options and returns exit code
func GetExitCodeForTerraformCommand(t testing.TestingT, additionalOptions *Options, args ...string) int {
	exitCode, err := GetExitCodeForTerraformCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// GetExitCodeForTerraformCommandE runs terraform with the given arguments and options and returns exit code
func GetExitCodeForTerraformCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (int, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	additionalOptions.Logger.Logf(t, "Running %s with args %v", options.TerraformBinary, args)
	cmd := generateCommand(options, args...)
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err == nil {
		return DefaultSuccessExitCode, nil
	}
	exitCode, getExitCodeErr := shell.GetExitCodeForRunCommandError(err)
	if getExitCodeErr == nil {
		return exitCode, nil
	}
	return DefaultErrorExitCode, getExitCodeErr
}
