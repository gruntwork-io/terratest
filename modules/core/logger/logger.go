// Package logger contains different methods to log.
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	gotesting "testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

const (
	// callDepthDirect is the call depth for functions called directly by the user.
	callDepthDirect = 2

	// callDepthWrapped is the call depth for functions called through one additional wrapper.
	callDepthWrapped = 3
)

var (
	// Default is the default logger that is used for the Logf function, if no one is provided. It uses the
	// TerratestLogger to log messages. This can be overwritten to change the logging globally.
	Default = New(terratestLogger{})
	// Discard discards all logging.
	Discard = New(discardLogger{})
	// Terratest logs the given format and arguments, formatted using fmt.Sprintf, to stdout, along with a timestamp and
	// information about what test and file is doing the logging. Before Go 1.14, this is an alternative to t.Logf as it
	// logs to stdout immediately, rather than buffering all log output and only displaying it at the very end of the test.
	// This is useful because:
	//
	// 1. It allows you to iterate faster locally, as you get feedback on whether your code changes are working as expected
	//    right away, rather than at the very end of the test run.
	//
	// 2. If you have a bug in your code that causes a test to never complete or if the test code crashes, t.Logf would
	//    show you no log output whatsoever, making debugging very hard, where as this method will show you all the log
	//    output available.
	//
	// 3. If you have a test that takes a long time to complete, some CI systems will kill the test suite prematurely
	//    because there is no log output with t.Logf (e.g., CircleCI kills tests after 10 minutes of no log output). With
	//    this log method, you get log output continuously.
	//
	// When a *testing.T is available, the log is emitted through t.Log rather than written to stdout directly, so that
	// `go test -json` attributes each line to the correct test even when tests run in parallel (see DoLog and issue
	// #1871). t.Log still streams immediately under `-v` on Go 1.14+, so the benefits above are preserved.
	//
	Terratest = New(terratestLogger{})
	// TestingT can be used to use Go's testing.T to log. If this is used, but no testing.T is provided, it will fallback
	// to Default.
	TestingT = New(testingT{})
)

// TestLogger is the interface for custom logger implementations that can be used with the Logger wrapper.
type TestLogger interface {
	Logf(t testing.TestingT, format string, args ...any)
}

// Logger wraps a TestLogger implementation and provides nil-safe logging.
type Logger struct {
	l TestLogger
}

// New creates a new Logger instance wrapping the given TestLogger implementation.
func New(l TestLogger) *Logger {
	return &Logger{
		l: l,
	}
}

// Logf logs the given format and arguments using the encapsulated TestLogger implementation.
func (l *Logger) Logf(t testing.TestingT, format string, args ...any) {
	if tt, ok := t.(helper); ok {
		tt.Helper()
	}

	if l == nil || l.l == nil {
		Default.Logf(t, format, args...)
		return
	}

	l.l.Logf(t, format, args...)
}

// helper is used to mark this library as a "helper", and thus not appearing in the line numbers. testing.T implements
// this interface, for example.
type helper interface {
	Helper()
}

type discardLogger struct{}

func (discardLogger) Logf(testing.TestingT, string, ...any) {}

type testingT struct{}

func (testingT) Logf(t testing.TestingT, format string, args ...any) {

	tt, ok := t.(*gotesting.T)
	if !ok {

		DoLog(t, callDepthDirect, os.Stdout, fmt.Sprintf(format, args...))
		return
	}

	tt.Helper()
	tt.Logf(format, args...)
}

type terratestLogger struct{}

func (terratestLogger) Logf(t testing.TestingT, format string, args ...any) {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	DoLog(t, callDepthWrapped, os.Stdout, fmt.Sprintf(format, args...))
}

// Log logs the given arguments to stdout, along with a timestamp and information about what test and file is doing the
// logging. This is an alternative to t.Logf that logs to stdout immediately, rather than buffering all log output and
// only displaying it at the very end of the test. See the Logf method for more info.
func Log(t testing.TestingT, args ...any) {
	if tt, ok := t.(helper); ok {
		tt.Helper()
	}

	MutexStdout.Lock()
	defer MutexStdout.Unlock()

	DoLog(t, callDepthDirect, os.Stdout, args...)
}

// MutexStdout is used to synchronize Log and Logf calls that write to stdout.
var MutexStdout sync.Mutex

// DoLog logs the given arguments to the given writer, along with a timestamp and information about what test and file is
// doing the logging.
func DoLog(t testing.TestingT, callDepth int, writer io.Writer, args ...any) {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	date := time.Now()
	prefix := fmt.Sprintf("%s %s %s:", t.Name(), date.Format(time.RFC3339), CallerPrefix(callDepth+1))
	allArgs := append([]any{prefix}, args...)

	// When we would otherwise write to stdout and a *testing.T is available, route the line through t.Log instead.
	// This lets `go test -json` attribute each line to the correct test, which it cannot do for raw stdout writes made
	// by tests running in parallel: such writes bypass the framework's per-test output coordination, so the JSON runner
	// tags them with whichever test happened to be active, mixing up the output of parallel tests (issue #1871). An
	// explicit non-stdout writer (e.g. a bytes.Buffer) is always honored as-is.
	if writer == os.Stdout {
		if sink, ok := t.(logSink); ok && logViaTestingT(sink, allArgs...) {
			return
		}
	}

	fmt.Fprintln(writer, allArgs...)
}

// logSink is satisfied by *testing.T (and *testing.B / *testing.F). See DoLog for why terratest routes stdout logging
// through it: doing so is what allows `go test -json` to attribute output to the right test under t.Parallel().
type logSink interface {
	Log(args ...any)
	Helper()
}

// logViaTestingT routes the given args through the testing.T's Log method, which formats them like fmt.Sprintln (the
// same as fmt.Fprintln) and streams them under the testing framework so they are attributed to the right test. It
// recovers if the test has already completed, since t.Log panics in that case, and reports false so the caller can fall
// back to writing directly to stdout rather than crash or drop the line.
func logViaTestingT(sink logSink, args ...any) (logged bool) {
	defer func() {
		if recover() != nil {
			logged = false
		}
	}()

	sink.Helper()
	sink.Log(args...)

	return true
}

// CallerPrefix returns the file and line number information about the methods that called this method, based on the current
// goroutine's stack. The argument callDepth is the number of stack frames to ascend, with 0 identifying the method
// that called CallerPrefix, 1 identifying the method that called that method, and so on.
//
// This code is adapted from testing.go, where it is in a private method called decorate.
func CallerPrefix(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if ok {

		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}

	return fmt.Sprintf("%s:%d", file, line)
}
