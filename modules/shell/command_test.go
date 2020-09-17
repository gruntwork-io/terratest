package shell

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestRunCommandAndGetOutput(t *testing.T) {
	t.Parallel()

	text := "Hello, World"
	cmd := Command{
		Command: "echo",
		Args:    []string{text},
	}

	out := RunCommandAndGetOutput(t, cmd)
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
	cmd := Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
	}

	out := RunCommandAndGetOutput(t, cmd)
	assert.Equal(t, expectedText, strings.TrimSpace(out))
}

func TestRunCommandGetExitCode(t *testing.T) {
	t.Parallel()

	cmd := Command{
		Command: "bash",
		Args:    []string{"-c", "exit 42"},
		Logger:  logger.Discard,
	}

	out, err := RunCommandAndGetOutputE(t, cmd)
	assert.Equal(t, "", out)
	assert.NotNil(t, err)
	code, err := GetExitCodeForRunCommandError(err)
	assert.Nil(t, err)
	assert.Equal(t, code, 42)
}

func TestRunCommandAndGetOutputConcurrency(t *testing.T) {
	t.Parallel()

	uniqueStderr := random.UniqueId()
	uniqueStdout := random.UniqueId()

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
	cmd := Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
		Logger:  logger.Discard,
	}

	out := RunCommandAndGetOutput(t, cmd)
	stdoutReg := regexp.MustCompile(uniqueStdout)
	stderrReg := regexp.MustCompile(uniqueStderr)
	assert.Equal(t, 500, len(stdoutReg.FindAllString(out, -1)))
	assert.Equal(t, 500, len(stderrReg.FindAllString(out, -1)))
}

func TestRunCommandWithHugeLineOutput(t *testing.T) {
	t.Parallel()

	// generate a ~100KB line
	bashCode := fmt.Sprintf(`
for i in {0..35000}
do
  echo -n foo
done
echo
`)

	cmd := Command{
		Command: "bash",
		Args:    []string{"-c", bashCode},
		Logger:  logger.Discard, // don't print that line to stdout
	}

	out, err := RunCommandAndGetOutputE(t, cmd)
	assert.NoError(t, err)

	var buffer bytes.Buffer
	for i := 0; i <= 35000; i++ {
		buffer.WriteString("foo")
	}

	assert.Equal(t, out, buffer.String())
}

// TestRunCommandOutputError ensures that getting the output never panics, even if no command was ever run.
func TestRunCommandOutputError(t *testing.T) {
	t.Parallel()

	cmd := Command{
		Command: "thisbinarydoesnotexistbecausenobodyusesnamesthatlong",
		Args:    []string{"-no-flag"},
		Logger:  logger.Discard,
	}

	out, err := RunCommandAndGetOutputE(t, cmd)
	assert.Equal(t, "", out)
	assert.NotNil(t, err)
}

func TestCommandOutputType(t *testing.T) {
	t.Parallel()

	stdout := "hello world"
	stderr := "this command has failed"

	_, err := RunCommandAndGetOutputE(t, Command{
		Command: "sh",
		Args:    []string{"-c", `echo "` + stdout + `" && echo "` + stderr + `" >&2 && exit 1`},
		Logger:  logger.Discard,
	})

	if err != nil {
		o, ok := err.(*ErrWithCmdOutput)
		if !ok {
			t.Fatalf("did not get correct type. got=%T", err)
		}
		assert.Len(t, o.Output.Stdout(), len(stdout))
		assert.Len(t, o.Output.Stderr(), len(stderr))
		assert.Len(t, o.Output.Combined(), len(stdout)+len(stderr)+1) // +1 for newline
	}
}
