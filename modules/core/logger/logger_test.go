package logger_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	tftesting "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoLog(t *testing.T) {
	t.Parallel()

	text := "test-do-log"

	var buffer bytes.Buffer

	logger.DoLog(t, 1, &buffer, text)

	assert.Regexp(t, fmt.Sprintf("^%s .+? [[:word:]]+.go:[0-9]+: %s$", t.Name(), text), strings.TrimSpace(buffer.String()))
}

type customLogger struct {
	logs []string
}

func (c *customLogger) Logf(_ tftesting.TestingT, format string, args ...any) {
	c.logs = append(c.logs, fmt.Sprintf(format, args...))
}

//nolint:paralleltest // test verifies nil Logger behavior and uses subtests that interact with shared state
func TestCustomLogger(t *testing.T) {
	logger.Logf(t, "this should be logged with the default logger")

	var l *logger.Logger
	l.Logf(t, "this should be logged with the default logger too")

	l = logger.New(nil)
	l.Logf(t, "this should be logged with the default logger too!")

	c := &customLogger{}
	l = logger.New(c)
	l.Logf(t, "log output 1")
	l.Logf(t, "log output 2")

	t.Run("logger-subtest", func(t *testing.T) {
		l.Logf(t, "subtest log")
	})

	assert.Len(t, c.logs, 3)
	assert.Equal(t, "log output 1", c.logs[0])
	assert.Equal(t, "log output 2", c.logs[1])
	assert.Equal(t, "subtest log", c.logs[2])
}

// TestLockedLog makes sure that Log and Logf which use stdout are thread-safe.
//
//nolint:paralleltest // test modifies os.Stdout
func TestLockedLog(t *testing.T) {
	stdout := os.Stdout

	t.Cleanup(func() {
		os.Stdout = stdout
	})

	data := []struct {
		fn   func(string)
		name string
	}{
		{
			fn: func(s string) {
				logger.Log(t, s)
			},
			name: "Log",
		},
		{
			fn: func(s string) {
				logger.Logf(t, "%s", s)
			},
			name: "Logf",
		},
	}

	for _, d := range data {
		logger.MutexStdout.Lock()
		str := "Logging something" + t.Name()

		r, w, _ := os.Pipe()
		os.Stdout = w
		ch := make(chan struct{})

		go func() {
			d.fn(str)
			w.Close()
			close(ch)
		}()

		select {
		case <-ch:
			t.Error("Log should be locked")
		default:
		}

		logger.MutexStdout.Unlock()

		b, err := io.ReadAll(r)
		require.NoError(t, err, "log should be unlocked")
		assert.Contains(t, string(b), str, "should contains logged string")
	}
}
