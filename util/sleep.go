package util

import (
	"time"
	"github.com/sirupsen/logrus"
)

//Sleep for a specified time.Duration and write a message to logger with a reason why
func SleepWithMessage(logger *logrus.Logger, duration time.Duration, whySleepMessage string) {
	logger.Printf("Sleeping %v: %s\n", duration, whySleepMessage)
	time.Sleep(duration)
}
