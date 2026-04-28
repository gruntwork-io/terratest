package shell

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Command is a simpler struct for defining commands than Go's built-in Cmd.
type Command struct {
	// Use the specified logger for the command's output. Use logger.Discard to not print the output while executing the command.
	Logger     *logger.Logger
	Stdin      io.Reader
	Env        map[string]string // Additional environment variables to set
	Command    string            // The command to run
	WorkingDir string            // The working directory
	Args       []string          // The args to pass to the command
}

// RunCommand runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself. If
// there are any errors, fail the test.
//
// Deprecated: Use RunCommandContext instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommand(t testing.TestingT, command Command) {
	RunCommandContext(t, context.Background(), &command)
}

// RunCommandContext is like RunCommand but includes a context.
func RunCommandContext(t testing.TestingT, ctx context.Context, command *Command) {
	err := RunCommandContextE(t, ctx, command)
	require.NoError(t, err)
}

// RunCommandE runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself. Any
// returned error will be of type ErrWithCmdOutput, containing the output streams and the underlying error.
//
// Deprecated: Use RunCommandContextE instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandE(t testing.TestingT, command Command) error {
	return RunCommandContextE(t, context.Background(), &command)
}

// RunCommandContextE is like RunCommandE but includes a context.
func RunCommandContextE(t testing.TestingT, ctx context.Context, command *Command) error {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return &ErrWithCmdOutput{err, output}
	}

	return nil
}

// RunCommandAndGetOutput runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of
// that command will also be logged with Command.Log to make debugging easier. If there are any errors, fail the test.
//
// Deprecated: Use RunCommandContextAndGetOutput instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandAndGetOutput(t testing.TestingT, command Command) string {
	return RunCommandContextAndGetOutput(t, context.Background(), &command)
}

// RunCommandContextAndGetOutput is like RunCommandAndGetOutput but includes a context.
func RunCommandContextAndGetOutput(t testing.TestingT, ctx context.Context, command *Command) string {
	out, err := RunCommandContextAndGetOutputE(t, ctx, command)
	require.NoError(t, err)

	return out
}

// RunCommandAndGetOutputE runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of
// that command will also be logged with Command.Log to make debugging easier. Any returned error will be of type
// ErrWithCmdOutput, containing the output streams and the underlying error.
//
// Deprecated: Use RunCommandContextAndGetOutputE instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandAndGetOutputE(t testing.TestingT, command Command) (string, error) {
	return RunCommandContextAndGetOutputE(t, context.Background(), &command)
}

// RunCommandContextAndGetOutputE is like RunCommandAndGetOutputE but includes a context.
func RunCommandContextAndGetOutputE(t testing.TestingT, ctx context.Context, command *Command) (string, error) {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return output.Combined(), &ErrWithCmdOutput{err, output}
	}

	return output.Combined(), nil
}

// RunCommandAndGetStdOut runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout and
// stderr of that command will also be logged with Command.Log to make debugging easier. If there are any errors, fail
// the test.
//
// Deprecated: Use RunCommandContextAndGetStdOut instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandAndGetStdOut(t testing.TestingT, command Command) string {
	return RunCommandContextAndGetStdOut(t, context.Background(), &command)
}

// RunCommandContextAndGetStdOut is like RunCommandAndGetStdOut but includes a context.
func RunCommandContextAndGetStdOut(t testing.TestingT, ctx context.Context, command *Command) string {
	output, err := RunCommandContextAndGetStdOutE(t, ctx, command)
	require.NoError(t, err)

	return output
}

// RunCommandAndGetStdOutE runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout
// and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging easier.
// Any returned error will be of type ErrWithCmdOutput, containing the output streams and the underlying error.
//
// Deprecated: Use RunCommandContextAndGetStdOutE instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandAndGetStdOutE(t testing.TestingT, command Command) (string, error) {
	return RunCommandContextAndGetStdOutE(t, context.Background(), &command)
}

// RunCommandContextAndGetStdOutE is like RunCommandAndGetStdOutE but includes a context.
func RunCommandContextAndGetStdOutE(t testing.TestingT, ctx context.Context, command *Command) (string, error) {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return output.Stdout(), &ErrWithCmdOutput{err, output}
	}

	return output.Stdout(), nil
}

// RunCommandAndGetStdOutErr runs a shell command and returns solely its stdout and stderr as a string. The stdout and
// stderr of that command will also be logged with Command.Log to make debugging easier. If there are any errors, fail
// the test.
//
// Deprecated: Use RunCommandContextAndGetStdOutErr instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandAndGetStdOutErr(t testing.TestingT, command Command) (stdout string, stderr string) {
	return RunCommandContextAndGetStdOutErr(t, context.Background(), &command)
}

// RunCommandContextAndGetStdOutErr is like RunCommandAndGetStdOutErr but includes a context.
func RunCommandContextAndGetStdOutErr(t testing.TestingT, ctx context.Context, command *Command) (stdout string, stderr string) {
	stdout, stderr, err := RunCommandContextAndGetStdOutErrE(t, ctx, command)
	require.NoError(t, err)

	return stdout, stderr
}

// RunCommandAndGetStdOutErrE runs a shell command and returns solely its stdout and stderr as a string. The stdout
// and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging easier.
// Any returned error will be of type ErrWithCmdOutput, containing the output streams and the underlying error.
//
// Deprecated: Use RunCommandContextAndGetStdOutErrE instead.
//
//nolint:gocritic // hugeParam - changing to pointer would break public API
func RunCommandAndGetStdOutErrE(t testing.TestingT, command Command) (stdout string, stderr string, err error) {
	return RunCommandContextAndGetStdOutErrE(t, context.Background(), &command)
}

// RunCommandContextAndGetStdOutErrE is like RunCommandAndGetStdOutErrE but includes a context.
func RunCommandContextAndGetStdOutErrE(t testing.TestingT, ctx context.Context, command *Command) (stdout string, stderr string, err error) {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return output.Stdout(), output.Stderr(), &ErrWithCmdOutput{err, output}
	}

	return output.Stdout(), output.Stderr(), nil
}

// ErrWithCmdOutput wraps an underlying error with the captured stdout and stderr from the command that produced it.
type ErrWithCmdOutput struct {
	Underlying error
	Output     *output
}

func (e *ErrWithCmdOutput) Error() string {
	return fmt.Sprintf("error while running command: %v; %s", e.Underlying, e.Output.Stderr())
}

// runCommand runs a shell command and stores each line from stdout and stderr in Output. Depending on the logger, the
// stdout and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging
// easier.
func runCommand(t testing.TestingT, ctx context.Context, command *Command) (*output, error) {
	command.Logger.Logf(t, "Running command %s with args %s", command.Command, command.Args)

	cmd := exec.CommandContext(ctx, command.Command, command.Args...)

	cmd.Dir = command.WorkingDir
	if command.Stdin != nil {
		cmd.Stdin = command.Stdin
	} else {
		cmd.Stdin = os.Stdin
	}

	cmd.Env = formatEnvVars(command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	output, err := readStdoutAndStderr(t, command.Logger, stdout, stderr)
	if err != nil {
		return output, err
	}

	return output, cmd.Wait()
}

// This function captures stdout and stderr into the given variables while still printing it to the stdout and stderr
// of this Go program
func readStdoutAndStderr(t testing.TestingT, log *logger.Logger, stdout, stderr io.ReadCloser) (*output, error) {
	out := newOutput()
	stdoutReader := bufio.NewReader(stdout)
	stderrReader := bufio.NewReader(stderr)

	wg := &sync.WaitGroup{}

	wg.Add(2) //nolint:mnd // 2 goroutines: one for stdout, one for stderr

	var stdoutErr, stderrErr error

	go func() {
		defer wg.Done()

		stdoutErr = readData(t, log, stdoutReader, out.stdout)
	}()

	go func() {
		defer wg.Done()

		stderrErr = readData(t, log, stderrReader, out.stderr)
	}()

	wg.Wait()

	if stdoutErr != nil {
		return out, stdoutErr
	}

	if stderrErr != nil {
		return out, stderrErr
	}

	return out, nil
}

func readData(t testing.TestingT, log *logger.Logger, reader *bufio.Reader, writer io.StringWriter) error {
	var (
		line    string
		readErr error
	)

	for {
		line, readErr = reader.ReadString('\n')

		// remove newline, our output is in a slice,
		// one element per line.
		line = strings.TrimSuffix(line, "\n")

		// only return early if the line does not have
		// any contents. We could have a line that does
		// not not have a newline before io.EOF, we still
		// need to add it to the output.
		if len(line) == 0 && readErr == io.EOF {
			break
		}

		// logger.Logger has a Logf method, but not a Log method.
		// We have to use the format string indirection to avoid
		// interpreting any possible formatting characters in
		// the line.
		//
		// See https://github.com/gruntwork-io/terratest/issues/982.
		log.Logf(t, "%s", line)

		if _, err := writer.WriteString(line); err != nil {
			return err
		}

		if readErr != nil {
			break
		}
	}

	if readErr != io.EOF {
		return readErr
	}

	return nil
}

// GetExitCodeForRunCommandError tries to read the exit code for the error object returned from running a shell command. This is a bit tricky to do
// in a way that works across platforms.
func GetExitCodeForRunCommandError(err error) (int, error) {
	var errWithOutput *ErrWithCmdOutput
	if errors.As(err, &errWithOutput) {
		err = errWithOutput.Underlying
	}

	// http://stackoverflow.com/a/10385867/483528
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
		}

		return 1, errors.New("could not determine exit code")
	}

	return 0, nil
}

func formatEnvVars(command *Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	return env
}
