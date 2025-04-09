package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoWithRetry(t *testing.T) {
	t.Parallel()

	expectedOutput := "expected"
	expectedError := fmt.Errorf("expected error")

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
		description   string
		maxRetries    int
		expectedError error
		action        func() (string, error)
	}{
		{"Return value on first try", 10, nil, actionAlwaysReturnsExpected},
		{"Return error on all retries", 10, MaxRetriesExceeded{Description: "Return error on all retries", MaxRetries: 10}, actionAlwaysReturnsError},
		{"Return value after 5 retries", 10, nil, createActionThatReturnsExpectedAfterFiveRetries()},
		{"Return value after 5 retries, but only do 4 retries", 4, MaxRetriesExceeded{Description: "Return value after 5 retries, but only do 4 retries", MaxRetries: 4}, createActionThatReturnsExpectedAfterFiveRetries()},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := DoWithRetryE(t, testCase.description, testCase.maxRetries, 1*time.Millisecond, testCase.action)
			assert.Equal(t, expectedOutput, actualOutput)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

func TestDoWithTimeout(t *testing.T) {
	t.Parallel()

	expectedOutput := "expected"
	expectedError := fmt.Errorf("expected error")

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
		description   string
		timeout       time.Duration
		expectedError error
		action        func() (string, error)
	}{
		{"Returns value immediately", 5 * time.Second, nil, actionReturnsValueImmediately},
		{"Returns error immediately", 5 * time.Second, expectedError, actionReturnsErrorImmediately},
		{"Returns value after 2 seconds", 5 * time.Second, nil, createActionThatReturnsValueAfterDelay(2 * time.Second)},
		{"Returns error after 2 seconds", 5 * time.Second, expectedError, createActionThatReturnsErrorAfterDelay(2 * time.Second)},
		{"Returns value after timeout exceeded", 5 * time.Second, TimeoutExceeded{Description: "Returns value after timeout exceeded", Timeout: 5 * time.Second}, createActionThatReturnsValueAfterDelay(10 * time.Second)},
		{"Returns error after timeout exceeded", 5 * time.Second, TimeoutExceeded{Description: "Returns error after timeout exceeded", Timeout: 5 * time.Second}, createActionThatReturnsErrorAfterDelay(10 * time.Second)},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := DoWithTimeoutE(t, testCase.description, testCase.timeout, testCase.action)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
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

	stop := DoInBackgroundUntilStopped(t, t.Name(), sleepBetweenRetries, func() {
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
	expectedError := fmt.Errorf("expected error")
	unexpectedError := fmt.Errorf("some other error")

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
		description     string
		retryableErrors map[string]string
		maxRetries      int
		expectedError   error
		action          func() (string, error)
	}{
		{"Return value on first try", noRetryableErrors, 10, nil, actionAlwaysReturnsExpected},
		{"Return expected error, but no retryable errors requested", noRetryableErrors, 10, FatalError{Underlying: expectedError}, actionAlwaysReturnsExpectedError},
		{"Return expected error, but retryable errors do not match", retryOnErrorsThatWontMatch, 10, FatalError{Underlying: expectedError}, actionAlwaysReturnsExpectedError},
		{"Return expected error on all retries, use match all regex", retryOnAllErrors, 10, MaxRetriesExceeded{Description: "Return expected error on all retries, use match all regex", MaxRetries: 10}, actionAlwaysReturnsExpectedError},
		{"Return expected error on all retries, use match exactly regex", retryOnExpectedErrorExactMatch, 3, MaxRetriesExceeded{Description: "Return expected error on all retries, use match exactly regex", MaxRetries: 3}, actionAlwaysReturnsExpectedError},
		{"Return expected error on all retries, use regex", retryOnExpectedErrorRegexpMatch, 1, MaxRetriesExceeded{Description: "Return expected error on all retries, use regex", MaxRetries: 1}, actionAlwaysReturnsExpectedError},
		{"Return expected error on all retries, use regex amidst others", retryOnExpectedErrorRegexpMatchWithOthers, 1, MaxRetriesExceeded{Description: "Return expected error on all retries, use regex amidst others", MaxRetries: 1}, actionAlwaysReturnsExpectedError},
		{"Return unexpected error on all retries, but match stdout exactly", retryOnExpectedStdoutExactMatch, 10, MaxRetriesExceeded{Description: "Return unexpected error on all retries, but match stdout exactly", MaxRetries: 10}, actionAlwaysReturnsUnexpectedError},
		{"Return unexpected error on all retries, but match stdout with regex", retryOnExpectedStdoutRegex, 3, MaxRetriesExceeded{Description: "Return unexpected error on all retries, but match stdout with regex", MaxRetries: 3}, actionAlwaysReturnsUnexpectedError},
		{"Return value after 5 retries with expected error, match all", retryOnAllErrors, 10, nil, createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{"Return value after 5 retries with expected error, match exactly", retryOnExpectedErrorExactMatch, 10, nil, createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{"Return value after 5 retries with expected error, match regex", retryOnExpectedErrorRegexpMatch, 10, nil, createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{"Return value after 5 retries with expected error, match multiple regex", retryOnExpectedErrorRegexpMatchWithOthers, 10, nil, createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{"Return value after 5 retries with expected error, match stdout exactly", retryOnExpectedStdoutExactMatch, 10, nil, createActionThatReturnsExpectedAfterFiveRetriesOfUnexpectedErrors()},
		{"Return value after 5 retries with expected error, match stdout with regex", retryOnExpectedStdoutRegex, 10, nil, createActionThatReturnsExpectedAfterFiveRetriesOfUnexpectedErrors()},
		{"Return value after 5 retries with expected error, match exactly, but only do 4 retries", retryOnExpectedErrorExactMatch, 4, MaxRetriesExceeded{Description: "Return value after 5 retries with expected error, match exactly, but only do 4 retries", MaxRetries: 4}, createActionThatReturnsExpectedAfterFiveRetriesOfExpectedErrors()},
		{"Return unexpected error after 5 retries with expected error, match exactly", retryOnExpectedErrorExactMatch, 10, FatalError{Underlying: ErrorCounter(6)}, createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
		{"Return unexpected error after 5 retries with expected error, match regex", retryOnExpectedErrorRegexpMatch, 10, FatalError{Underlying: ErrorCounter(6)}, createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
		{"Return unexpected error after 5 retries with expected error, match multiple regex", retryOnExpectedErrorRegexpMatchWithOthers, 10, FatalError{Underlying: ErrorCounter(6)}, createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
		{"Return unexpected error after 5 retries with expected error, match all", retryOnAllErrors, 10, MaxRetriesExceeded{Description: "Return unexpected error after 5 retries with expected error, match all", MaxRetries: 10}, createActionThatReturnsErrorCounterAfterFiveRetriesOfExpectedErrors()},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable for each test case

		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actualOutput, err := DoWithRetryableErrorsE(t, testCase.description, testCase.retryableErrors, testCase.maxRetries, 1*time.Millisecond, testCase.action)
			assert.Equal(t, expectedOutput, actualOutput)
			if testCase.expectedError != nil {
				assert.Equal(t, testCase.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedOutput, actualOutput)
			}
		})
	}
}

type ErrorCounter int

func (count ErrorCounter) Error() string {
	return fmt.Sprintf("%d", int(count))
}
