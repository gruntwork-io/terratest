package shell_test

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/gruntwork-io/terratest/modules/shell/v2"
)

func TestRunCommandAndGetOutput(t *testing.T) {
	t.Parallel()

	text := "Hello, World"
	cmd := &shell.Command{
		Command: "echo",
		Args:    []string{text},
	}

	out := shell.RunCommandContextAndGetOutput(t, t.Context(), cmd)
	assert.Equal(t, text, strings.TrimSpace(out))
}

func TestRunCommandAndGetOutputOrder(t *testing.T) {
	t.Parallel()

	stderrText := "Hello, Error"
	stdoutText := "Hello, World"
	expectedText := "Hello, Error\nHello, World\nHello, Error\nHello, World\nHello, Error\nHello, Error"
	bashCode := fmt.Sprintf(`
echo_stderr(){
	(>&2 echo "%s")
	# Add sleep to stabilize the test
	sleep .01s
}
echo_stdout(){
	echo "%s"
	# Add sleep to stabilize the test
	sleep .01s
}
echo_stderr
echo_stdout
echo_stderr
echo_stdout
echo_stderr
echo_stderr
`,
		stderrText,
		stdoutText,
	)
	cmd := &shell.Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
	}

	out := shell.RunCommandContextAndGetOutput(t, t.Context(), cmd)
	assert.Equal(t, expectedText, strings.TrimSpace(out))
}

func TestRunCommandGetExitCode(t *testing.T) {
	t.Parallel()

	cmd := &shell.Command{
		Command: "bash",
		Args:    []string{"-c", "exit 42"},
		Logger:  logger.Discard,
	}

	out, err := shell.RunCommandContextAndGetOutputE(t, t.Context(), cmd)
	assert.Empty(t, out)
	require.Error(t, err)

	code, err := shell.GetExitCodeForRunCommandError(err)
	require.NoError(t, err)
	assert.Equal(t, 42, code)
}

func TestRunCommandAndGetOutputConcurrency(t *testing.T) {
	t.Parallel()

	uniqueStderr := random.UniqueID()
	uniqueStdout := random.UniqueID()

	bashCode := fmt.Sprintf(`
echo_stderr(){
	sleep .0$[ ( $RANDOM %% 10 ) + 1 ]s
	(>&2 echo "%s")
}
echo_stdout(){
	sleep .0$[ ( $RANDOM %% 10 ) + 1 ]s
	echo "%s"
}
for i in {1..500}
do
	echo_stderr &
	echo_stdout &
done
wait
`,
		uniqueStderr,
		uniqueStdout,
	)
	cmd := &shell.Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
		Logger:  logger.Discard,
	}

	out := shell.RunCommandContextAndGetOutput(t, t.Context(), cmd)

	stdoutReg := regexp.MustCompile(uniqueStdout)
	stderrReg := regexp.MustCompile(uniqueStderr)

	assert.Len(t, stdoutReg.FindAllString(out, -1), 500)
	assert.Len(t, stderrReg.FindAllString(out, -1), 500)
}

func TestRunCommandWithHugeLineOutput(t *testing.T) {
	t.Parallel()

	// generate a ~100KB line
	bashCode := `
for i in {0..35000}
do
  echo -n foo
done
echo
`

	cmd := &shell.Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
		Logger:  logger.Discard, // don't print that line to stdout
	}

	out, err := shell.RunCommandContextAndGetOutputE(t, t.Context(), cmd)
	require.NoError(t, err)

	var buffer bytes.Buffer

	for i := 0; i <= 35000; i++ {
		buffer.WriteString("foo")
	}

	assert.Equal(t, out, buffer.String())
}

// TestRunCommandOutputError ensures that getting the output never panics, even if no command was ever run.
func TestRunCommandOutputError(t *testing.T) {
	t.Parallel()

	cmd := &shell.Command{
		Command: "thisbinarydoesnotexistbecausenobodyusesnamesthatlong",
		Args:    []string{"-no-flag"},
		Logger:  logger.Discard,
	}

	out, err := shell.RunCommandContextAndGetOutputE(t, t.Context(), cmd)
	assert.Empty(t, out)
	assert.Error(t, err)
}

func TestCommandOutputType(t *testing.T) {
	t.Parallel()

	stdout := "hello world"
	stderr := "this command has failed"

	_, err := shell.RunCommandContextAndGetOutputE(t, t.Context(), &shell.Command{
		Command: "sh",
		Args:    []string{"-c", `echo "` + stdout + `" && echo "` + stderr + `" >&2 && exit 1`},
		Logger:  logger.Discard,
	})
	if err != nil {
		var o *shell.ErrWithCmdOutput
		if !errors.As(err, &o) {
			t.Fatalf("did not get correct type. got=%T", err)
		}

		assert.Len(t, o.Output.Stdout(), len(stdout))
		assert.Len(t, o.Output.Stderr(), len(stderr))
		assert.Len(t, o.Output.Combined(), len(stdout)+len(stderr)+1) // +1 for newline
	}
}

func TestCommandWithStdoutAndStdErr(t *testing.T) {
	t.Parallel()

	stdout := "hello world"
	stderr := "this command has failed"
	command := &shell.Command{
		Command: "sh",
		Args:    []string{"-c", `echo "` + stdout + `" && echo "` + stderr + `" >&2`},
		Logger:  logger.Discard,
	}

	t.Run("MustNotError", func(t *testing.T) {
		t.Parallel()

		ostdout, ostderr := shell.RunCommandContextAndGetStdOutErr(t, t.Context(), command)
		assert.Equal(t, stdout, ostdout)
		assert.Equal(t, stderr, ostderr)
	})

	t.Run("ReturnError", func(t *testing.T) {
		t.Parallel()

		ostdout, ostderr, err := shell.RunCommandContextAndGetStdOutErrE(t, t.Context(), command)
		require.NoError(t, err)
		assert.Equal(t, stdout, ostdout)
		assert.Equal(t, stderr, ostderr)
	})
}

func TestRunCommandWithStdinAndGetOutput(t *testing.T) {
	t.Parallel()

	text := "Hello, World"
	cmd := &shell.Command{
		Command: "cat",
		Stdin:   strings.NewReader(text),
	}

	out := shell.RunCommandContextAndGetOutput(t, t.Context(), cmd)
	assert.Equal(t, text, strings.TrimSpace(out))
}
