package terraform

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
)

// GetCommonOptions extracts commons terraform options
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	if options.NoColor && !collections.ListContains(args, "-no-color") {
		args = append(args, "-no-color")
	}

	// If Lock is set, override the lock argument that has been passed in (for apply, plan and destroy commands)
	if collections.ListContains(args, "plan") || collections.ListContains(args, "apply") || collections.ListContains(args, "destroy") {
		lockKey := "-lock"
		for i, argValue := range args {
			parts := strings.Split(argValue, "=")
			if parts[0] == lockKey {
				args[i] = fmt.Sprintf("%s=%s", lockKey, strconv.FormatBool(options.Lock))
			}
		}

		if options.ExtraArguments != nil {
			args = append(args, options.ExtraArguments...)
		}
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
func RunTerraformCommand(t *testing.T, additionalOptions *Options, args ...string) string {
	out, err := RunTerraformCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunTerraformCommandE runs terraform with the given arguments and options and return stdout/stderr.
func RunTerraformCommandE(t *testing.T, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	description := fmt.Sprintf("Running terraform %v", args)
	return retry.DoWithRetryE(t, description, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		cmd := shell.Command{
			Command:    "terraform",
			Args:       args,
			WorkingDir: options.TerraformDir,
			Env:        options.EnvVars,
		}

		out, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil {
			return out, nil
		}

		for errorText, errorMessage := range options.RetryableTerraformErrors {
			if strings.Contains(out, errorText) {
				logger.Logf(t, "terraform failed with the error '%s' but this error was expected and warrants a retry. Further details: %s\n", errorText, errorMessage)
				return out, err
			}
		}

		return out, retry.FatalError{Underlying: err}
	})
}

// GetExitCodeForTerraformCommand runs terraform with the given arguments and options and returns exit code
func GetExitCodeForTerraformCommand(t *testing.T, additionalOptions *Options, args ...string) int {
	exitCode, err := GetExitCodeForTerraformCommandE(t, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// GetExitCodeForTerraformCommandE runs terraform with the given arguments and options and returns exit code
func GetExitCodeForTerraformCommandE(t *testing.T, additionalOptions *Options, additionalArgs ...string) (int, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	logger.Log(t, "Running terraform %v", args)
	cmd := shell.Command{
		Command:    "terraform",
		Args:       args,
		WorkingDir: options.TerraformDir,
		Env:        options.EnvVars,
	}

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
