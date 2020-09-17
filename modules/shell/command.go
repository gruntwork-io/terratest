package shell

import (
	"bufio"
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
	Command    string            // The command to run
	Args       []string          // The args to pass to the command
	WorkingDir string            // The working directory
	Env        map[string]string // Additional environment variables to set
	// Use the specified logger for the command's output. Use logger.Discard to not print the output while executing the command.
	Logger *logger.Logger
}

// RunCommand runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself. If
// there are any errors, fail the test.
func RunCommand(t testing.TestingT, command Command) {
	err := RunCommandE(t, command)
	require.NoError(t, err)
}

// RunCommandE runs a shell command and redirects its stdout and stderr to the stdout of the atomic script itself. Any
// returned error will be of type ErrWithCmdOutput, containing the output streams and the underlying error.
func RunCommandE(t testing.TestingT, command Command) error {
	output, err := runCommand(t, command)
	if err != nil {
		return &ErrWithCmdOutput{err, output}
	}
	return nil
}

// RunCommandAndGetOutput runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of
// that command will also be logged with Command.Log to make debugging easier. If there are any errors, fail the test.
func RunCommandAndGetOutput(t testing.TestingT, command Command) string {
	out, err := RunCommandAndGetOutputE(t, command)
	require.NoError(t, err)
	return out
}

// RunCommandAndGetOutputE runs a shell command and returns its stdout and stderr as a string. The stdout and stderr of
// that command will also be logged with Command.Log to make debugging easier. Any returned error will be of type
// ErrWithCmdOutput, containing the output streams and the underlying error.
func RunCommandAndGetOutputE(t testing.TestingT, command Command) (string, error) {
	output, err := runCommand(t, command)
	if err != nil {
		return output.Combined(), &ErrWithCmdOutput{err, output}
	}

	return output.Combined(), nil
}

// RunCommandAndGetStdOut runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout and
// stderr of that command will also be logged with Command.Log to make debugging easier. If there are any errors, fail
// the test.
func RunCommandAndGetStdOut(t testing.TestingT, command Command) string {
	output, err := RunCommandAndGetStdOutE(t, command)
	require.NoError(t, err)
	return output
}

// RunCommandAndGetStdOutE runs a shell command and returns solely its stdout (but not stderr) as a string. The stdout
// and stderr of that command will also be printed to the stdout and stderr of this Go program to make debugging easier.
// Any returned error will be of type ErrWithCmdOutput, containing the output streams and the underlying error.
func RunCommandAndGetStdOutE(t testing.TestingT, command Command) (string, error) {
	output, err := runCommand(t, command)
	if err != nil {
		return output.Stdout(), &ErrWithCmdOutput{err, output}
	}

	return output.Stdout(), nil
}

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
func runCommand(t testing.TestingT, command Command) (*output, error) {
	command.Logger.Logf(t, "Running command %s with args %s", command.Command, command.Args)

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.WorkingDir
	cmd.Stdin = os.Stdin
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
	var line string
	var readErr error
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

		log.Logf(t, line)
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
	if errWithOutput, ok := err.(*ErrWithCmdOutput); ok {
		err = errWithOutput.Underlying
	}

	// http://stackoverflow.com/a/10385867/483528
	if exitErr, ok := err.(*exec.ExitError); ok {
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

func formatEnvVars(command Command) []string {
	env := os.Environ()
	for key, value := range command.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	return env
}
