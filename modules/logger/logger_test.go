package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	tftesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
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
