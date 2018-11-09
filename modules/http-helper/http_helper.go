// Package http_helper contains helpers to interact with deployed resources through HTTP.
package http_helper

import (
	"fmt"
	// TODO Guido: remove this import
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/http-helper"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
)

// HttpGet performs an HTTP GET on the given URL and return the HTTP status code and body. If there's any error, fail the test.
func HttpGet(t *testing.T, url string) (int, string) {
	statusCode, body, err := HttpGetE(t, url)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, body
}

// HttpGetE performs an HTTP GET on the given URL and return the HTTP status code, body, and any error.
func HttpGetE(t *testing.T, url string) (int, string, error) {
	logger.Logf(t, "Making an HTTP GET call to URL %s", url)

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return -1, "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(body)), nil
}

// HttpGetWithValidation performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, fail the test.
func HttpGetWithValidation(t *testing.T, url string, expectedStatusCode int, expectedBody string) {
	err := HttpGetWithValidationE(t, url, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithValidationE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
func HttpGetWithValidationE(t *testing.T, url string, expectedStatusCode int, expectedBody string) error {
	return HttpGetWithCustomValidationE(t, url, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HttpGetWithCustomValidation performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidation(t *testing.T, url string, validateResponse func(int, string) bool) {
	err := HttpGetWithCustomValidationE(t, url, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithCustomValidationE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidationE(t *testing.T, url string, validateResponse func(int, string) bool) error {
	statusCode, body, err := HttpGetE(t, url)

	if err != nil {
		return err
	}

	if !validateResponse(statusCode, body) {
		return ValidationFunctionFailed{Url: url, Status: statusCode, Body: body}
	}

	return nil
}

// HttpGetWithRetry repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetry(t *testing.T, url string, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpGetWithRetryE(t, url, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryE(t *testing.T, url string, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithValidationE(t, url, expectedStatus, expectedBody)
	})

	return err
}

// HttpGetWithRetryWithCustomValidation repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidation(t *testing.T, url string, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	err := HttpGetWithRetryWithCustomValidationE(t, url, retries, sleepBetweenRetries, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryWithCustomValidationE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationE(t *testing.T, url string, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithCustomValidationE(t, url, validateResponse)
	})

	return err
}

// HttpDo performs the given HTTP method on the given URL and return the HTTP status code and body.
// If there's any error, fail the test.
func HttpDo(t *testing.T, method string, url string, body io.Reader, headers map[string]string) (int, string) {
	statusCode, respBody, err := HttpDoE(t, method, url, body, headers)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, respBody
}

// HttpDoE performs the given HTTP method on the given URL and return the HTTP status code, body, and any error.
func HttpDoE(t *testing.T, method string, url string, body io.Reader, headers map[string]string) (int, string, error) {
	logger.Logf(t, "Making an HTTP %s call to URL %s", method, url)

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
	}

	req := newRequest(method, url, body, headers)
	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(respBody)), nil
}

// HttpDoWithRetry repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
func HttpDoWithRetry(t *testing.T, method string, url string, body io.Reader, headers map[string]string, expectedStatus int, retries int, sleepBetweenRetries time.Duration) string {
	out, err := HttpDoWithRetryE(t, method, url, body, headers, expectedStatus, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// HttpDoWithRetryE repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
func HttpDoWithRetryE(t *testing.T, method string, url string, body io.Reader, headers map[string]string, expectedStatus int, retries int, sleepBetweenRetries time.Duration) (string, error) {
	out, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP %s to URL %s", method, url), retries, sleepBetweenRetries, func() (string, error) {
		statusCode, out, err := HttpDoE(t, method, url, body, headers)
		if err != nil {
			return "", err
		}
		logger.Logf(t, "output: %v", out)
		if statusCode != expectedStatus {
			return "", http_helper.ValidationFunctionFailed{Url: url, Status: statusCode}
		}
		return out, nil
	})

	return out, err
}

// HttpDoWithValidationRetry repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
func HttpDoWithValidationRetry(t *testing.T, method string, url string, body io.Reader, headers map[string]string, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpDoWithValidationRetryE(t, method, url, body, headers, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpDoWithValidaitonRetryE repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
func HttpDoWithValidationRetryE(t *testing.T, method string, url string, body io.Reader, headers map[string]string, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpDoWithValidationE(t, method, url, body, headers, expectedStatus, expectedBody)
	})

	return err
}

// HttpDoWithValidation performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, fail the test.
func HttpDoWithValidation(t *testing.T, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string) {
	err := HttpDoWithValidationE(t, method, url, body, headers, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpDoWithValidationE performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, return an error.
func HttpDoWithValidationE(t *testing.T, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string) error {
	return HttpDoWithCustomValidationE(t, method, url, body, headers, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HttpDoWithCustomValidation performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HttpDoWithCustomValidation(t *testing.T, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool) {
	err := HttpDoWithCustomValidationE(t, method, url, body, headers, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpDoWithCustomValidationE performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HttpDoWithCustomValidationE(t *testing.T, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool) error {
	statusCode, respBody, err := HttpDoE(t, method, url, body, headers)

	if err != nil {
		return err
	}

	if !validateResponse(statusCode, respBody) {
		return http_helper.ValidationFunctionFailed{Url: url, Status: statusCode, Body: respBody}
	}

	return nil
}

func newRequest(method string, url string, body io.Reader, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req
}

// ValidationFunctionFailed is an error that occurs if a validation function fails.
type ValidationFunctionFailed struct {
	Url    string
	Status int
	Body   string
}

func (err ValidationFunctionFailed) Error() string {
	return fmt.Sprintf("Validation failed for URL %s. Response status: %d. Response body:\n%s", err.Url, err.Status, err.Body)
}
