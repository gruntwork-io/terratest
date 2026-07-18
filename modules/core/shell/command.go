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

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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

// RunCommandContext is like RunCommand but includes a context.
func RunCommandContext(t testing.TestingT, ctx context.Context, command *Command) {
	err := RunCommandContextE(t, ctx, command)
	require.NoError(t, err)
}

// RunCommandContextE is like RunCommandE but includes a context.
func RunCommandContextE(t testing.TestingT, ctx context.Context, command *Command) error {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return &ErrWithCmdOutput{err, output}
	}

	return nil
}

// RunCommandContextAndGetOutput is like RunCommandAndGetOutput but includes a context.
func RunCommandContextAndGetOutput(t testing.TestingT, ctx context.Context, command *Command) string {
	out, err := RunCommandContextAndGetOutputE(t, ctx, command)
	require.NoError(t, err)

	return out
}

// RunCommandContextAndGetOutputE is like RunCommandAndGetOutputE but includes a context.
func RunCommandContextAndGetOutputE(t testing.TestingT, ctx context.Context, command *Command) (string, error) {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return output.Combined(), &ErrWithCmdOutput{err, output}
	}

	return output.Combined(), nil
}

// RunCommandContextAndGetStdOut is like RunCommandAndGetStdOut but includes a context.
func RunCommandContextAndGetStdOut(t testing.TestingT, ctx context.Context, command *Command) string {
	output, err := RunCommandContextAndGetStdOutE(t, ctx, command)
	require.NoError(t, err)

	return output
}

// RunCommandContextAndGetStdOutE is like RunCommandAndGetStdOutE but includes a context.
func RunCommandContextAndGetStdOutE(t testing.TestingT, ctx context.Context, command *Command) (string, error) {
	output, err := runCommand(t, ctx, command)
	if err != nil {
		return output.Stdout(), &ErrWithCmdOutput{err, output}
	}

	return output.Stdout(), nil
}

// RunCommandContextAndGetStdOutErr is like RunCommandAndGetStdOutErr but includes a context.
func RunCommandContextAndGetStdOutErr(t testing.TestingT, ctx context.Context, command *Command) (stdout string, stderr string) {
	stdout, stderr, err := RunCommandContextAndGetStdOutErrE(t, ctx, command)
	require.NoError(t, err)

	return stdout, stderr
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

	wg.Add(2)

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

		line = strings.TrimSuffix(line, "\n")

		if len(line) == 0 && readErr == io.EOF {
			break
		}

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
