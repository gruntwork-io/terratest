package terraform

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"

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
		Stdin:      options.Stdin,
	}

	return cmd
}

var commandsWithParallelism = []string{
	"plan",
	"apply",
	"destroy",
}

const (
	// TofuDefaultPath command to run tofu
	TofuDefaultPath = "tofu"

	// TerraformDefaultPath to run terraform
	TerraformDefaultPath = "terraform"
)

// DefaultExecutable is the default terraform executable to use. It is set to "terraform" if the terraform binary
// is available, otherwise it falls back to "tofu".
var DefaultExecutable = defaultTerraformExecutable()

// GetCommonOptions extracts commons terraform options
func GetCommonOptions(options *Options, args ...string) (*Options, []string) {
	if options.TerraformBinary == "" {
		options.TerraformBinary = DefaultExecutable
	}

	if options.Parallelism > 0 && len(args) > 0 && slices.Contains(commandsWithParallelism, args[0]) {
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

// RunTerraformCommandContext runs terraform with the given arguments and options and returns stdout/stderr.
func RunTerraformCommandContext(t testing.TestingT, ctx context.Context, additionalOptions *Options, args ...string) string {
	out, err := RunTerraformCommandContextE(t, ctx, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// RunTerraformCommandContextE runs terraform with the given arguments and options and returns stdout/stderr.
func RunTerraformCommandContextE(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) (string, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)

	return retry.DoWithRetryableErrorsContextE(t, ctx, description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		s, err := shell.RunCommandContextAndGetOutputE(t, ctx, &cmd)
		if err != nil {
			return s, err
		}

		if err := hasWarning(additionalOptions, s); err != nil {
			return s, err
		}

		return s, err
	})
}

// RunTerraformCommand runs terraform with the given arguments and options and return stdout/stderr.
//
// Deprecated: Use [RunTerraformCommandContext] instead.
func RunTerraformCommand(t testing.TestingT, additionalOptions *Options, args ...string) string {
	return RunTerraformCommandContext(t, context.Background(), additionalOptions, args...)
}

// RunTerraformCommandE runs terraform with the given arguments and options and return stdout/stderr.
//
// Deprecated: Use [RunTerraformCommandContextE] instead.
func RunTerraformCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	return RunTerraformCommandContextE(t, context.Background(), additionalOptions, additionalArgs...)
}

// RunTerraformCommandAndGetStdoutContext runs terraform with the given arguments and options and returns solely its
// stdout (but not stderr).
func RunTerraformCommandAndGetStdoutContext(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) string {
	out, err := RunTerraformCommandAndGetStdoutContextE(t, ctx, additionalOptions, additionalArgs...)
	require.NoError(t, err)

	return out
}

// RunTerraformCommandAndGetStdoutContextE runs terraform with the given arguments and options and returns solely its
// stdout (but not stderr).
func RunTerraformCommandAndGetStdoutContextE(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) (string, error) {
	out, _, _, err := RunTerraformCommandAndGetStdOutErrCodeContextE(t, ctx, additionalOptions, additionalArgs...)

	return out, err
}

// RunTerraformCommandAndGetStdout runs terraform with the given arguments and options and returns solely its stdout
// (but not stderr).
//
// Deprecated: Use [RunTerraformCommandAndGetStdoutContext] instead.
func RunTerraformCommandAndGetStdout(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) string {
	return RunTerraformCommandAndGetStdoutContext(t, context.Background(), additionalOptions, additionalArgs...)
}

// RunTerraformCommandAndGetStdoutE runs terraform with the given arguments and options and returns solely its stdout
// (but not stderr).
//
// Deprecated: Use [RunTerraformCommandAndGetStdoutContextE] instead.
func RunTerraformCommandAndGetStdoutE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (string, error) {
	return RunTerraformCommandAndGetStdoutContextE(t, context.Background(), additionalOptions, additionalArgs...)
}

// RunTerraformCommandAndGetStdOutErrCodeContext runs terraform with the given arguments and options and returns its
// stdout, stderr, and exitcode.
func RunTerraformCommandAndGetStdOutErrCodeContext(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) (stdout string, stderr string, exit int) {
	stdout, stderr, exit, err := RunTerraformCommandAndGetStdOutErrCodeContextE(t, ctx, additionalOptions, additionalArgs...)
	require.NoError(t, err)

	return stdout, stderr, exit
}

// RunTerraformCommandAndGetStdOutErrCodeContextE runs terraform with the given arguments and options and returns its
// stdout, stderr, and exitcode.
func RunTerraformCommandAndGetStdOutErrCodeContextE(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) (stdout string, stderr string, exit int, err error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	cmd := generateCommand(options, args...)
	description := fmt.Sprintf("%s %v", options.TerraformBinary, args)

	exit = DefaultErrorExitCode

	_, err = retry.DoWithRetryableErrorsContextE(t, ctx, description, options.RetryableTerraformErrors, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		stdout, stderr, err = shell.RunCommandContextAndGetStdOutErrE(t, ctx, &cmd)
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

// RunTerraformCommandAndGetStdOutErrCode runs terraform with the given arguments and options and returns its stdout,
// stderr, and exitcode.
//
// Deprecated: Use [RunTerraformCommandAndGetStdOutErrCodeContext] instead.
func RunTerraformCommandAndGetStdOutErrCode(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (stdout string, stderr string, exit int) {
	return RunTerraformCommandAndGetStdOutErrCodeContext(t, context.Background(), additionalOptions, additionalArgs...)
}

// RunTerraformCommandAndGetStdOutErrCodeE runs terraform with the given arguments and options and returns its stdout,
// stderr, and exitcode.
//
// Deprecated: Use [RunTerraformCommandAndGetStdOutErrCodeContextE] instead.
func RunTerraformCommandAndGetStdOutErrCodeE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (stdout string, stderr string, exit int, err error) {
	return RunTerraformCommandAndGetStdOutErrCodeContextE(t, context.Background(), additionalOptions, additionalArgs...)
}

// GetExitCodeForTerraformCommandContext runs terraform with the given arguments and options and returns exit code.
func GetExitCodeForTerraformCommandContext(t testing.TestingT, ctx context.Context, additionalOptions *Options, args ...string) int {
	exitCode, err := GetExitCodeForTerraformCommandContextE(t, ctx, additionalOptions, args...)
	if err != nil {
		t.Fatal(err)
	}

	return exitCode
}

// GetExitCodeForTerraformCommandContextE runs terraform with the given arguments and options and returns exit code.
func GetExitCodeForTerraformCommandContextE(t testing.TestingT, ctx context.Context, additionalOptions *Options, additionalArgs ...string) (int, error) {
	options, args := GetCommonOptions(additionalOptions, additionalArgs...)

	additionalOptions.Logger.Logf(t, "Running %s with args %v", options.TerraformBinary, args)

	cmd := generateCommand(options, args...)

	_, err := shell.RunCommandContextAndGetOutputE(t, ctx, &cmd)
	if err == nil {
		return DefaultSuccessExitCode, nil
	}

	exitCode, getExitCodeErr := shell.GetExitCodeForRunCommandError(err)
	if getExitCodeErr == nil {
		return exitCode, nil
	}

	return DefaultErrorExitCode, getExitCodeErr
}

// GetExitCodeForTerraformCommand runs terraform with the given arguments and options and returns exit code.
//
// Deprecated: Use [GetExitCodeForTerraformCommandContext] instead.
func GetExitCodeForTerraformCommand(t testing.TestingT, additionalOptions *Options, args ...string) int {
	return GetExitCodeForTerraformCommandContext(t, context.Background(), additionalOptions, args...)
}

// GetExitCodeForTerraformCommandE runs terraform with the given arguments and options and returns exit code.
//
// Deprecated: Use [GetExitCodeForTerraformCommandContextE] instead.
func GetExitCodeForTerraformCommandE(t testing.TestingT, additionalOptions *Options, additionalArgs ...string) (int, error) {
	return GetExitCodeForTerraformCommandContextE(t, context.Background(), additionalOptions, additionalArgs...)
}

func defaultTerraformExecutable() string {
	cmd := exec.CommandContext(context.Background(), TerraformDefaultPath, "-version")
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
