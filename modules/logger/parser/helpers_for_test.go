package parser

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func NewTestLogger(t *testing.T) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&LogTestFormatter{TestName: t.Name()})
	return logger
}

type LogTestFormatter struct {
	TestName string
}

func (formatter *LogTestFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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
