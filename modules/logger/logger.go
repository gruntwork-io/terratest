// Package logger contains different methods to log.
package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	gotesting "testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/testing"
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
	Terratest = New(terratestLogger{})
	// TestingT can be used to use Go's testing.T to log. If this is used, but no testing.T is provided, it will fallback
	// to Default.
	TestingT = New(testingT{})
)

type TestLogger interface {
	Logf(t testing.TestingT, format string, args ...interface{})
}

type Logger struct {
	l TestLogger
}

func New(l TestLogger) *Logger {
	return &Logger{
		l,
	}
}

func (l *Logger) Logf(t testing.TestingT, format string, args ...interface{}) {
	if tt, ok := t.(helper); ok {
		tt.Helper()
	}

	// methods can be called on (typed) nil pointers. In this case, use the Default function to log. This enables the
	// caller to do `var l *Logger` and then use the logger already.
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

func (_ discardLogger) Logf(_ testing.TestingT, format string, args ...interface{}) {}

type testingT struct{}

func (_ testingT) Logf(t testing.TestingT, format string, args ...interface{}) {
	// this should never fail
	tt, ok := t.(*gotesting.T)
	if !ok {
		// fallback
		DoLog(t, 2, os.Stdout, fmt.Sprintf(format, args...))
		return
	}

	tt.Helper()
	tt.Logf(format, args...)
	return
}

type terratestLogger struct{}

func (_ terratestLogger) Logf(t testing.TestingT, format string, args ...interface{}) {
	DoLog(t, 3, os.Stdout, fmt.Sprintf(format, args...))
}

// Deprecated: use Logger instead, as it provides more flexibility on logging.
// Logf logs the given format and arguments, formatted using fmt.Sprintf, to stdout, along with a timestamp and information
// about what test and file is doing the logging. Before Go 1.14, this is an alternative to t.Logf as it logs to stdout
// immediately, rather than buffering all log output and only displaying it at the very end of the test. This is useful
// because:
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
// Although t.Logf now supports streaming output since Go 1.14, this is kept for compatibility purposes.
func Logf(t testing.TestingT, format string, args ...interface{}) {
	if tt, ok := t.(helper); ok {
		tt.Helper()
	}

	DoLog(t, 2, os.Stdout, fmt.Sprintf(format, args...))
}

// Log logs the given arguments to stdout, along with a timestamp and information about what test and file is doing the
// logging. This is an alternative to t.Logf that logs to stdout immediately, rather than buffering all log output and
// only displaying it at the very end of the test. See the Logf method for more info.
func Log(t testing.TestingT, args ...interface{}) {
	if tt, ok := t.(helper); ok {
		tt.Helper()
	}

	DoLog(t, 2, os.Stdout, args...)
}

// DoLog logs the given arguments to the given writer, along with a timestamp and information about what test and file is
// doing the logging.
func DoLog(t testing.TestingT, callDepth int, writer io.Writer, args ...interface{}) {
	date := time.Now()
	prefix := fmt.Sprintf("%s %s %s:", t.Name(), date.Format(time.RFC3339), CallerPrefix(callDepth+1))
	allArgs := append([]interface{}{prefix}, args...)
	fmt.Fprintln(writer, allArgs...)
}

// CallerPrefix returns the file and line number information about the methods that called this method, based on the current
// goroutine's stack. The argument callDepth is the number of stack frames to ascend, with 0 identifying the method
// that called CallerPrefix, 1 identifying the method that called that method, and so on.
//
// This code is adapted from testing.go, where it is in a private method called decorate.
func CallerPrefix(callDepth int) string {
	_, file, line, ok := runtime.Caller(callDepth)
	if ok {
		// Truncate file name at last file name separator.
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
