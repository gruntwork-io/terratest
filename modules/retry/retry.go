// Package retry contains logic to retry actions with certain conditions.
package retry

import (
	"fmt"
	"regexp"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"golang.org/x/net/context"
)

// Either contains a result and potentially an error.
type Either struct {
	Result string
	Error  error
}

// DoWithTimeout runs the specified action and waits up to the specified timeout for it to complete. Return the output of the action if
// it completes on time or fail the test otherwise.
func DoWithTimeout(t testing.TestingT, actionDescription string, timeout time.Duration, action func() (string, error)) string {
	out, err := DoWithTimeoutE(t, actionDescription, timeout, action)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// DoWithTimeoutE runs the specified action and waits up to the specified timeout for it to complete. Return the output of the action if
// it completes on time or an error otherwise.
func DoWithTimeoutE(t testing.TestingT, actionDescription string, timeout time.Duration, action func() (string, error)) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChannel := make(chan Either, 1)

	go func() {
		out, err := action()
		resultChannel <- Either{Result: out, Error: err}
	}()

	select {
	case either := <-resultChannel:
		return either.Result, either.Error
	case <-ctx.Done():
		return "", TimeoutExceeded{Description: actionDescription, Timeout: timeout}
	}
}

// DoWithRetry runs the specified action. If it returns a value, return that value. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, fail the test.
func DoWithRetry(t testing.TestingT, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() (string, error)) string {
	out, err := DoWithRetryE(t, actionDescription, maxRetries, sleepBetweenRetries, action)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// DoWithRetryE runs the specified action. If it returns a value, return that value. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryE(t testing.TestingT, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() (string, error)) (string, error) {
	var output string
	var err error

	for i := 0; i <= maxRetries; i++ {
		logger.Log(t, actionDescription)

		output, err = action()
		if err == nil {
			return output, nil
		}

		if _, isFatalErr := err.(FatalError); isFatalErr {
			logger.Logf(t, "Returning due to fatal error: %v", err)
			return output, err
		}

		logger.Logf(t, "%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return output, MaxRetriesExceeded{Description: actionDescription, MaxRetries: maxRetries}
}

// DoWithRetryableErrors runs the specified action. If it returns a value, return that value. If it returns an error,
// check if error message or the string output from the action (which is often stdout/stderr from running some command)
// matches any of the regular expressions in the specified retryableErrors map. If there is a match, sleep for
// sleepBetweenRetries, and retry the specified action, up to a maximum of maxRetries retries. If there is no match,
// return that error immediately, wrapped in a FatalError. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryableErrors(t testing.TestingT, actionDescription string, retryableErrors map[string]string, maxRetries int, sleepBetweenRetries time.Duration, action func() (string, error)) string {
	out, err := DoWithRetryableErrorsE(t, actionDescription, retryableErrors, maxRetries, sleepBetweenRetries, action)
	require.NoError(t, err)
	return out
}

// DoWithRetryableErrorsE runs the specified action. If it returns a value, return that value. If it returns an error,
// check if error message or the string output from the action (which is often stdout/stderr from running some command)
// matches any of the regular expressions in the specified retryableErrors map. If there is a match, sleep for
// sleepBetweenRetries, and retry the specified action, up to a maximum of maxRetries retries. If there is no match,
// return that error immediately, wrapped in a FatalError. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func DoWithRetryableErrorsE(t testing.TestingT, actionDescription string, retryableErrors map[string]string, maxRetries int, sleepBetweenRetries time.Duration, action func() (string, error)) (string, error) {
	retryableErrorsRegexp := map[*regexp.Regexp]string{}
	for errorStr, errorMessage := range retryableErrors {
		errorRegex, err := regexp.Compile(errorStr)
		if err != nil {
			return "", FatalError{Underlying: err}
		}
		retryableErrorsRegexp[errorRegex] = errorMessage
	}

	return DoWithRetryE(t, actionDescription, maxRetries, sleepBetweenRetries, func() (string, error) {
		output, err := action()
		if err == nil {
			return output, nil
		}

		for errorRegexp, errorMessage := range retryableErrorsRegexp {
			if errorRegexp.MatchString(output) || errorRegexp.MatchString(err.Error()) {
				logger.Logf(t, "'%s' failed with the error '%s' but this error was expected and warrants a retry. Further details: %s\n", actionDescription, err.Error(), errorMessage)
				return output, err
			}
		}

		return output, FatalError{Underlying: err}
	})
}

// Done can be stopped.
type Done struct {
	stop chan bool
}

// Done stops the execution.
func (done Done) Done() {
	done.stop <- true
}

// DoInBackgroundUntilStopped runs the specified action in the background (in a goroutine) repeatedly, waiting the specified amount of time between
// repetitions. To stop this action, call the Done() function on the returned value.
func DoInBackgroundUntilStopped(t testing.TestingT, actionDescription string, sleepBetweenRepeats time.Duration, action func()) Done {
	stop := make(chan bool)

	go func() {
		for {
			logger.Logf(t, "Executing action '%s'", actionDescription)

			action()

			logger.Logf(t, "Sleeping for %s before repeating action '%s'", sleepBetweenRepeats, actionDescription)

			select {
			case <-time.After(sleepBetweenRepeats):
				// Nothing to do, just allow the loop to continue
			case <-stop:
				logger.Logf(t, "Received stop signal for action '%s'.", actionDescription)
				return
			}
		}
	}()

	return Done{stop: stop}
}

// Custom error types

// TimeoutExceeded is an error that occurs when a timeout is exceeded.
type TimeoutExceeded struct {
	Description string
	Timeout     time.Duration
}

func (err TimeoutExceeded) Error() string {
	return fmt.Sprintf("'%s' did not complete before timeout of %s", err.Description, err.Timeout)
}

// MaxRetriesExceeded is an error that occurs when the maximum amount of retries is exceeded.
type MaxRetriesExceeded struct {
	Description string
	MaxRetries  int
}

func (err MaxRetriesExceeded) Error() string {
	return fmt.Sprintf("'%s' unsuccessful after %d retries", err.Description, err.MaxRetries)
}

// FatalError is a marker interface for errors that should not be retried.
type FatalError struct {
	Underlying error
}

func (err FatalError) Error() string {
	return fmt.Sprintf("FatalError{Underlying: %v}", err.Underlying)
}
