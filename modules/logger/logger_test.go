package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	tftesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoLog(t *testing.T) {
	t.Parallel()

	text := "test-do-log"
	var buffer bytes.Buffer

	DoLog(t, 1, &buffer, text)

	assert.Regexp(t, fmt.Sprintf("^%s .+? [[:word:]]+.go:[0-9]+: %s$", t.Name(), text), strings.TrimSpace(buffer.String()))
}

type customLogger struct {
	logs []string
}

func (c *customLogger) Logf(t tftesting.TestingT, format string, args ...interface{}) {
	c.logs = append(c.logs, fmt.Sprintf(format, args...))
}

func TestCustomLogger(t *testing.T) {
	Logf(t, "this should be logged with the default logger")

	var l *Logger
	l.Logf(t, "this should be logged with the default logger too")

	l = New(nil)
	l.Logf(t, "this should be logged with the default logger too!")

	c := &customLogger{}
	l = New(c)
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

// TestLockedLog make sure that Log and Logf which use stdout are thread-safe
func TestLockedLog(t *testing.T) {
	// should not call t.Parallel() since we are modifying os.Stdout
	stdout := os.Stdout
	t.Cleanup(func() {
		os.Stdout = stdout
	})

	data := []struct {
		name string
		fn   func(*testing.T, string)
	}{
		{
			name: "Log",
			fn: func(t *testing.T, s string) {
				Log(t, s)
			}},
		{
			name: "Logf",
			fn: func(t *testing.T, s string) {
				Logf(t, "%s", s)
			}},
	}

	for _, d := range data {
		mutexStdout.Lock()
		str := "Logging something" + t.Name()

		r, w, _ := os.Pipe()
		os.Stdout = w
		ch := make(chan struct{})
		go func() {
			d.fn(t, str)
			w.Close()
			close(ch)
		}()

		select {
		case <-ch:
			t.Error("Log should be locked")
		default:
		}

		mutexStdout.Unlock()
		b, err := io.ReadAll(r)
		require.NoError(t, err, "log should be unlocked")
		assert.Contains(t, string(b), str, "should contains logged string")
	}

}
