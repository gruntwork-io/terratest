package parser_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

// NewTestLogger creates a logrus.Logger configured for test output.
func NewTestLogger(t *testing.T) *logrus.Logger {
	t.Helper()

	logger := logrus.New()
	logger.SetFormatter(&logTestFormatter{TestName: t.Name()})

	return logger
}

type logTestFormatter struct {
	TestName string
}

func (formatter *logTestFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := bytes.Buffer{}
	outStr := fmt.Sprintf(
		"%s %s %s %s\n",
		formatter.TestName,
		strings.ToUpper(entry.Level.String()),
		entry.Time.Format(time.RFC3339),
		entry.Message,
	)
	b.WriteString(outStr)

	return b.Bytes(), nil
}
