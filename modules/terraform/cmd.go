package terraform

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
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
	"run-all",
	"apply-all",
	"destroy-all",
}

const (
	// TofuDefaultPath command to run tofu
	TofuDefaultPath = "tofu"

	// TerraformDefaultPath to run terraform
	TerraformDefaultPath = "terraform"

	// TerragruntDefaultPath to run terragrunt
	TerragruntDefaultPath = "terragrunt"
)

var DefaultExecutable = defaultTerraformExecutable()

// GetCommonOptions extracts commons terraform options
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	if options.TerraformBinary == "" {
		options.TerraformBinary = DefaultExecutable
	}

	if options.TerraformBinary == TerragruntDefaultPath {
		args = append(args, "--terragrunt-non-interactive")

		// for newer Terragrunt version, setting simplified log formatting
		setTerragruntLogFormatting(options)
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
		s, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			return s, err
		}
		if err := hasWarning(additionalOptions, s); err != nil {
			return s, err
		}
		return s, err
	})

}

// RunTerraformCommandAndGetStdout runs terraform with the given arguments and options and returns solely its stdout
// (but not stderr).
func RunTerraformCommandAndGetStdout(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) string {
	out, err := RunTerraformCommandAndGetStdoutE(t, additionalOptions, additionalArgs...)
	require.NoError(t, err)
	return out
}

// RunTerraformCommandAndGetStdoutE runs terraform with the given arguments and options and returns solely its stdout
// (but not stderr).
func RunTerraformCommandAndGetStdoutE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	out, _, _, err := RunTerraformCommandAndGetStdOutErrCodeE(t, additionalOptions, additionalArgs...)
	return out, err
}

// RunTerraformCommandAndGetStdOutErrCode runs terraform with the given arguments and options and returns its stdout, stderr, and exitcode
func RunTerraformCommandAndGetStdOutErrCode(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (stdout string, stderr string, exit int) {
	stdout, stderr, exit, err := RunTerraformCommandAndGetStdOutErrCodeE(t, additionalOptions, additionalArgs...)
	require.NoError(t, err)
	return stdout, stderr, exit
}

// RunTerraformCommandAndGetStdOutErrCodeE runs terraform with the given arguments and options and returns its stdout, stderr, and exitcode
func RunTerraformCommandAndGetStdOutErrCodeE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (stdout string, stderr string, exit int, err error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)

	exit = DefaultErrorExitCode
	_, err = retry.DoWithRetryableErrorsE(t, description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		stdout, stderr, err = shell.RunCommandAndGetStdOutErrE(t, cmd)
		if err != nil {
			exitCode, getExitCodeErr := shell.GetExitCodeForRunCommandError(err)
			if getExitCodeErr == nil {
				exit = exitCode
			}
			return "", err
		}

		if err = hasWarning(additionalOptions, stdout); err != nil {
			return "", err
		}

		exit = DefaultSuccessExitCode
		return "", nil
	})

	return
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

func defaultTerraformExecutable() string {
	cmd := exec.Command(TerraformDefaultPath, "-version")
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err == nil {
		return TerraformDefaultPath
	}

	// fallback to Tofu if terraform is not available
	return TofuDefaultPath
}

func hasWarning(opts *Options, out string) error {
	for k, v := range opts.WarningsAsErrors {
		str := fmt.Sprintf("\n.*(?i:Warning): %s[^\n]*\n", k)
		re, err := regexp.Compile(str)
		if err != nil {
			return fmt.Errorf("cannot compile regex for warning detection: %w", err)
		}
		m := re.FindAllString(out, -1)
		if len(m) == 0 {
			continue
		}
		return fmt.Errorf("warning(s) were found: %s:\n%s", v, strings.Join(m, ""))
	}
	return nil
}

// setTerragruntLogFormatting sets a default log formatting for terragrunt
// if it is not already set in options.EnvVars or OS environment vars
func setTerragruntLogFormatting(options *Options) {
	const (
		tgLogFormatKey       = "TERRAGRUNT_LOG_FORMAT"
		tgLogCustomFormatKey = "TERRAGRUNT_LOG_CUSTOM_FORMAT"
	)

	if options.EnvVars == nil {
		options.EnvVars = map[string]string{}
	}

	_, inOpts := options.EnvVars[tgLogFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(tgLogFormatKey)
		if !inEnv {
			// key-value format for terragrunt logs to avoid colors and have plain form
			// https://terragrunt.gruntwork.io/docs/reference/cli-options/#terragrunt-log-format
			options.EnvVars[tgLogFormatKey] = "key-value"
		}
	}

	_, inOpts = options.EnvVars[tgLogCustomFormatKey]
	if !inOpts {
		_, inEnv := os.LookupEnv(tgLogCustomFormatKey)
		if !inEnv {
			options.EnvVars[tgLogCustomFormatKey] = "%msg(color=disable)"
		}
	}
}
