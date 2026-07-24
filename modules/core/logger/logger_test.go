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
	logger.Default.Logf(t, "this should be logged with the default logger")

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

// fakeTestingT implements testing.TestingT but is deliberately NOT a *testing.T and has no Log method, so logging
// through it exercises DoLog's stdout fallback path. A real *testing.T is instead routed through t.Log (see DoLog).
type fakeTestingT struct{ name string }

func (fakeTestingT) Fail()                 {}
func (fakeTestingT) FailNow()              {}
func (fakeTestingT) Fatal(...any)          {}
func (fakeTestingT) Fatalf(string, ...any) {}
func (fakeTestingT) Error(...any)          {}
func (fakeTestingT) Errorf(string, ...any) {}
func (f fakeTestingT) Name() string        { return f.name }
func (fakeTestingT) Helper()               {}

// spyTestingT satisfies both testing.TestingT and the logSink that DoLog routes stdout logging through, recording each
// Log call so tests can assert on what was routed.
type spyTestingT struct {
	fakeTestingT
	logged []string
}

func (s *spyTestingT) Log(args ...any) { s.logged = append(s.logged, fmt.Sprintln(args...)) }

// TestDoLogRoutesThroughTestingT verifies that DoLog routes to t.Log when writing to stdout for a *testing.T (so that
// `go test -json` attributes output to the correct test), while always honoring an explicit non-stdout writer.
//
//nolint:paralleltest // asserts on os.Stdout routing
func TestDoLogRoutesThroughTestingT(t *testing.T) {
	// writer == os.Stdout with a testing.T-like sink: routed through Log, nothing written to the real stdout.
	spy := &spyTestingT{fakeTestingT: fakeTestingT{name: "TestApply1"}}
	logger.DoLog(spy, 1, os.Stdout, "routed-message")
	require.Len(t, spy.logged, 1)
	assert.Contains(t, spy.logged[0], "TestApply1")
	assert.Contains(t, spy.logged[0], "routed-message")

	// An explicit non-stdout writer is always honored and never routed to the sink.
	var buf bytes.Buffer

	spy2 := &spyTestingT{fakeTestingT: fakeTestingT{name: "TestApply2"}}

	logger.DoLog(spy2, 1, &buf, "buffered-message")
	assert.Empty(t, spy2.logged)
	assert.Contains(t, buf.String(), "buffered-message")
}

// TestLockedLog makes sure that Log which uses the stdout fallback path is thread-safe. It uses fakeTestingT (not a
// *testing.T) so DoLog writes to stdout under MutexStdout rather than routing through t.Log.
//
//nolint:paralleltest // test modifies os.Stdout
func TestLockedLog(t *testing.T) {
	stdout := os.Stdout

	t.Cleanup(func() {
		os.Stdout = stdout
	})

	ft := fakeTestingT{name: t.Name()}

	data := []struct {
		fn   func(string)
		name string
	}{
		{
			fn: func(s string) {
				logger.Log(ft, s)
			},
			name: "Log",
		},
		{
			fn: func(s string) {
				logger.Default.Logf(ft, "%s", s)
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
