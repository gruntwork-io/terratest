package retry_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
)

func TestDoWithRetry(t *testing.T) {
	t.Parallel()

	expectedOutput := "expected"
	expectedError := errors.New("expected error")

	actionAlwaysReturnsExpected := func() (string, error) { return expectedOutput, nil }
	actionAlwaysReturnsError := func() (string, error) { return expectedOutput, expectedError }

	createActionThatReturnsExpectedAfterFiveRetries := func() func() (string, error) {
		count := 0

		return func() (string, error) {
			count++

			if count > 5 {
				return expectedOutput, nil
			}

			return expectedOutput, expectedError
		}
	}

	testCases := []struct {
		expectedError error
		action        func() (string, error)
		description   string
		maxRetries    int
	}{
		{description: "Return value on first try", maxRetries: 10, action: actionAlwaysReturnsExpected},
		{description: "Return error on all retries", maxRetries: 10, expectedError: retry.MaxRetriesExceeded{Description: "Return error on all retries", MaxRetries: 10}, action: actionAlwaysReturnsError},
		{description: "Return value after 5 retries", maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetries()},
		{description: "Return value after 5 retries, but only do 4 retries", maxRetries: 4, expectedError: retry.MaxRetriesExceeded{Description: "Return value after 5 retries, but only do 4 retries", MaxRetries: 4}, action: createActionThatReturnsExpectedAfterFiveRetries()},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := retry.DoWithRetryE(t, testCase.description, testCase.maxRetries, 1*time.Millisecond, testCase.action)
			assert.Equal(t, expectedOutput, actualOutput)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

func TestDoWithTimeout(t *testing.T) {
	t.Parallel()

	expectedOutput := "expected"
	expectedError := errors.New("expected error")

	actionReturnsValueImmediately := func() (string, error) { return expectedOutput, nil }
	actionReturnsErrorImmediately := func() (string, error) { return "", expectedError }

	createActionThatReturnsValueAfterDelay := func(delay time.Duration) func() (string, error) {
		return func() (string, error) {
			time.Sleep(delay)

			return expectedOutput, nil
		}
	}

	createActionThatReturnsErrorAfterDelay := func(delay time.Duration) func() (string, error) {
		return func() (string, error) {
			time.Sleep(delay)

			return "", expectedError
		}
	}

	testCases := []struct {
		expectedError error
		action        func() (string, error)
		description   string
		timeout       time.Duration
	}{
		{description: "Returns value immediately", timeout: 5 * time.Second, action: actionReturnsValueImmediately},
		{description: "Returns error immediately", timeout: 5 * time.Second, expectedError: expectedError, action: actionReturnsErrorImmediately},
		{description: "Returns value after 2 seconds", timeout: 5 * time.Second, action: createActionThatReturnsValueAfterDelay(2 * time.Second)},
		{description: "Returns error after 2 seconds", timeout: 5 * time.Second, expectedError: expectedError, action: createActionThatReturnsErrorAfterDelay(2 * time.Second)},
		{description: "Returns value after timeout exceeded", timeout: 5 * time.Second, expectedError: retry.TimeoutExceeded{Description: "Returns value after timeout exceeded", Timeout: 5 * time.Second}, action: createActionThatReturnsValueAfterDelay(10 * time.Second)},
		{description: "Returns error after timeout exceeded", timeout: 5 * time.Second, expectedError: retry.TimeoutExceeded{Description: "Returns error after timeout exceeded", Timeout: 5 * time.Second}, action: createActionThatReturnsErrorAfterDelay(10 * time.Second)},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := retry.DoWithTimeoutE(t, testCase.description, testCase.timeout, testCase.action)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

func TestDoInBackgroundUntilStopped(t *testing.T) {
	t.Parallel()

	sleepBetweenRetries := 2 * time.Second
	waitStop := sleepBetweenRetries*2 + sleepBetweenRetries/2
	counter := 0

	stop := retry.DoInBackgroundUntilStopped(t, t.Name(), sleepBetweenRetries, func() {
		counter++
		t.Log(time.Now(), counter)
	})

	time.Sleep(waitStop)
	stop.Done()
	assert.Equal(t, 3, counter)

	time.Sleep(waitStop)
	assert.Equal(t, 3, counter)
}

func TestDoWithRetryableErrors(t *testing.T) {
	t.Parallel()

	expectedOutput := "this is the expected output"
	expectedError := errors.New("expected error")
	unexpectedError := errors.New("some other error")

	actionAlwaysReturnsExpected := func() (string, error) { return expectedOutput, nil }
	actionAlwaysReturnsExpectedError := func() (string, error) { return expectedOutput, expectedError }
	actionAlwaysReturnsUnexpectedError := func() (string, error) { return expectedOutput, unexpectedError }

	createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors := func() func() (string, error) {
		count := 0

		return func() (string, error) {
			count++

			if count > 5 {
				return expectedOutput, nil
			}

			return expectedOutput, expectedError
		}
	}

	createActionThatReturnsExpectedAfterFiveRetriesOfUnexpectedErrors := func() func() (string, error) {
		count := 0

		return func() (string, error) {
			count++

			if count > 5 {
				return expectedOutput, nil
			}

			return expectedOutput, unexpectedError
		}
	}

	createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors := func() func() (string, error) {
		count := 0

		return func() (string, error) {
			count++

			if count > 5 {
				return expectedOutput, ErrorCounter(count)
			}

			return expectedOutput, expectedError
		}
	}

	matchAllRegexp := ".*"
	matchExpectedErrorExactRegexp := expectedError.Error()
	matchExpectedErrorRegexp := "^expected.*$"
	matchNothingRegexp1 := "this won't match any of our errors"
	matchNothingRegexp2 := "this also won't match any of our errors"
	matchStdoutExactlyRegexp := expectedOutput
	matchStdoutRegexp := "this.*output"

	noRetryableErrors := map[string]string{}
	retryOnAllErrors := map[string]string{
		matchAllRegexp: "match all errors",
	}
	retryOnExpectedErrorExactMatch := map[string]string{
		matchExpectedErrorExactRegexp: "match expected error exactly",
	}
	retryOnExpectedErrorRegexpMatch := map[string]string{
		matchExpectedErrorRegexp: "match expected error using a regex",
	}
	retryOnExpectedErrorRegexpMatchWithOthers := map[string]string{
		matchNothingRegexp1:      "unrelated regex that shouldn't match anything",
		matchExpectedErrorRegexp: "match expected error using a regex",
		matchNothingRegexp2:      "another unrelated regex that shouldn't match anything",
	}
	retryOnErrorsThatWontMatch := map[string]string{
		matchNothingRegexp1: "unrelated regex that shouldn't match anything",
		matchNothingRegexp2: "another unrelated regex that shouldn't match anything",
	}
	retryOnExpectedStdoutExactMatch := map[string]string{
		matchStdoutExactlyRegexp: "match expected stdout exactly",
	}
	retryOnExpectedStdoutRegex := map[string]string{
		matchStdoutRegexp: "match expected stdout using a regex",
	}

	testCases := []struct {
		expectedError   error
		retryableErrors map[string]string
		action          func() (string, error)
		description     string
		maxRetries      int
	}{
		{description: "Return value on first try", retryableErrors: noRetryableErrors, maxRetries: 10, action: actionAlwaysReturnsExpected},
		{description: "Return expected error, but no retryable errors requested", retryableErrors: noRetryableErrors, maxRetries: 10, expectedError: retry.FatalError{Underlying: expectedError}, action: actionAlwaysReturnsExpectedError},
		{description: "Return expected error, but retryable errors do not match", retryableErrors: retryOnErrorsThatWontMatch, maxRetries: 10, expectedError: retry.FatalError{Underlying: expectedError}, action: actionAlwaysReturnsExpectedError},
		{description: "Return expected error on all retries, use match all regex", retryableErrors: retryOnAllErrors, maxRetries: 10, expectedError: retry.MaxRetriesExceeded{Description: "Return expected error on all retries, use match all regex", MaxRetries: 10}, action: actionAlwaysReturnsExpectedError},
		{description: "Return expected error on all retries, use match exactly regex", retryableErrors: retryOnExpectedErrorExactMatch, maxRetries: 3, expectedError: retry.MaxRetriesExceeded{Description: "Return expected error on all retries, use match exactly regex", MaxRetries: 3}, action: actionAlwaysReturnsExpectedError},
		{description: "Return expected error on all retries, use regex", retryableErrors: retryOnExpectedErrorRegexpMatch, maxRetries: 1, expectedError: retry.MaxRetriesExceeded{Description: "Return expected error on all retries, use regex", MaxRetries: 1}, action: actionAlwaysReturnsExpectedError},
		{description: "Return expected error on all retries, use regex amidst others", retryableErrors: retryOnExpectedErrorRegexpMatchWithOthers, maxRetries: 1, expectedError: retry.MaxRetriesExceeded{Description: "Return expected error on all retries, use regex amidst others", MaxRetries: 1}, action: actionAlwaysReturnsExpectedError},
		{description: "Return unexpected error on all retries, but match stdout exactly", retryableErrors: retryOnExpectedStdoutExactMatch, maxRetries: 10, expectedError: retry.MaxRetriesExceeded{Description: "Return unexpected error on all retries, but match stdout exactly", MaxRetries: 10}, action: actionAlwaysReturnsUnexpectedError},
		{description: "Return unexpected error on all retries, but match stdout with regex", retryableErrors: retryOnExpectedStdoutRegex, maxRetries: 3, expectedError: retry.MaxRetriesExceeded{Description: "Return unexpected error on all retries, but match stdout with regex", MaxRetries: 3}, action: actionAlwaysReturnsUnexpectedError},
		{description: "Return value after 5 retries with expected error, match all", retryableErrors: retryOnAllErrors, maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{description: "Return value after 5 retries with expected error, match exactly", retryableErrors: retryOnExpectedErrorExactMatch, maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{description: "Return value after 5 retries with expected error, match regex", retryableErrors: retryOnExpectedErrorRegexpMatch, maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{description: "Return value after 5 retries with expected error, match multiple regex", retryableErrors: retryOnExpectedErrorRegexpMatchWithOthers, maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{description: "Return value after 5 retries with expected error, match stdout exactly", retryableErrors: retryOnExpectedStdoutExactMatch, maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetriesOfUnexpectedErrors()},
		{description: "Return value after 5 retries with expected error, match stdout with regex", retryableErrors: retryOnExpectedStdoutRegex, maxRetries: 10, action: createActionThatReturnsExpectedAfterFiveRetriesOfUnexpectedErrors()},
		{description: "Return value after 5 retries with expected error, match exactly, but only do 4 retries", retryableErrors: retryOnExpectedErrorExactMatch, maxRetries: 4, expectedError: retry.MaxRetriesExceeded{Description: "Return value after 5 retries with expected error, match exactly, but only do 4 retries", MaxRetries: 4}, action: createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{description: "Return unexpected error after 5 retries with expected error, match exactly", retryableErrors: retryOnExpectedErrorExactMatch, maxRetries: 10, expectedError: retry.FatalError{Underlying: ErrorCounter(6)}, action: createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
		{description: "Return unexpected error after 5 retries with expected error, match regex", retryableErrors: retryOnExpectedErrorRegexpMatch, maxRetries: 10, expectedError: retry.FatalError{Underlying: ErrorCounter(6)}, action: createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
		{description: "Return unexpected error after 5 retries with expected error, match multiple regex", retryableErrors: retryOnExpectedErrorRegexpMatchWithOthers, maxRetries: 10, expectedError: retry.FatalError{Underlying: ErrorCounter(6)}, action: createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
		{description: "Return unexpected error after 5 retries with expected error, match all", retryableErrors: retryOnAllErrors, maxRetries: 10, expectedError: retry.MaxRetriesExceeded{Description: "Return unexpected error after 5 retries with expected error, match all", MaxRetries: 10}, action: createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := retry.DoWithRetryableErrorsE(t, testCase.description, testCase.retryableErrors, testCase.maxRetries, 1*time.Millisecond, testCase.action)
			assert.Equal(t, expectedOutput, actualOutput)

			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

type ErrorCounter int

func (count ErrorCounter) Error() string {
	return strconv.Itoa(int(count))
}
