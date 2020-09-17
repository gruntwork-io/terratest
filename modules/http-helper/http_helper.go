// Package http_helper contains helpers to interact with deployed resources through HTTP.
package http_helper

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// HttpGet performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code and body. If there's any error, fail the test.
func HttpGet(t testing.TestingT, url string, tlsConfig *tls.Config) (int, string) {
	statusCode, body, err := HttpGetE(t, url, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, body
}

// HttpGetE performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code, body, and any error.
func HttpGetE(t testing.TestingT, url string, tlsConfig *tls.Config) (int, string, error) {
	logger.Logf(t, "Making an HTTP GET call to URL %s", url)

	// Set HTTP client transport config
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = tlsConfig

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: 10 * time.Second,
		// Include the previously created transport config
		Transport: tr,
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
func HttpGetWithValidation(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) {
	err := HttpGetWithValidationE(t, url, tlsConfig, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithValidationE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
func HttpGetWithValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) error {
	return HttpGetWithCustomValidationE(t, url, tlsConfig, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HttpGetWithCustomValidation performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidation(t testing.TestingT, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) {
	err := HttpGetWithCustomValidationE(t, url, tlsConfig, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithCustomValidationE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) error {
	statusCode, body, err := HttpGetE(t, url, tlsConfig)

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
func HttpGetWithRetry(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpGetWithRetryE(t, url, tlsConfig, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryE(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithValidationE(t, url, tlsConfig, expectedStatus, expectedBody)
	})

	return err
}

// HttpGetWithRetryWithCustomValidation repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidation(t testing.TestingT, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	err := HttpGetWithRetryWithCustomValidationE(t, url, tlsConfig, retries, sleepBetweenRetries, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryWithCustomValidationE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithCustomValidationE(t, url, tlsConfig, validateResponse)
	})

	return err
}

// HTTPDo performs the given HTTP method on the given URL and return the HTTP status code and body.
// If there's any error, fail the test.
func HTTPDo(
	t testing.TestingT, method string, url string, body io.Reader,
	headers map[string]string, tlsConfig *tls.Config,
) (int, string) {
	statusCode, respBody, err := HTTPDoE(t, method, url, body, headers, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, respBody
}

// HTTPDoE performs the given HTTP method on the given URL and return the HTTP status code, body, and any error.
func HTTPDoE(
	t testing.TestingT, method string, url string, body io.Reader,
	headers map[string]string, tlsConfig *tls.Config,
) (int, string, error) {
	logger.Logf(t, "Making an HTTP %s call to URL %s", method, url)

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout:   10 * time.Second,
		Transport: tr,
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

// HTTPDoWithRetry repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
func HTTPDoWithRetry(
	t testing.TestingT, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) string {
	out, err := HTTPDoWithRetryE(t, method, url, body,
		headers, expectedStatus, retries, sleepBetweenRetries, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// HTTPDoWithRetryE repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
func HTTPDoWithRetryE(
	t testing.TestingT, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) (string, error) {
	out, err := retry.DoWithRetryE(
		t, fmt.Sprintf("HTTP %s to URL %s", method, url), retries,
		sleepBetweenRetries, func() (string, error) {
			bodyReader := bytes.NewReader(body)
			statusCode, out, err := HTTPDoE(t, method, url, bodyReader, headers, tlsConfig)
			if err != nil {
				return "", err
			}
			logger.Logf(t, "output: %v", out)
			if statusCode != expectedStatus {
				return "", ValidationFunctionFailed{Url: url, Status: statusCode}
			}
			return out, nil
		})

	return out, err
}

// HTTPDoWithValidationRetry repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
func HTTPDoWithValidationRetry(
	t testing.TestingT, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) {
	err := HTTPDoWithValidationRetryE(t, method, url, body, headers, expectedStatus, expectedBody, retries, sleepBetweenRetries, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithValidationRetryE repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
func HTTPDoWithValidationRetryE(
	t testing.TestingT, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP %s to URL %s", method, url), retries,
		sleepBetweenRetries, func() (string, error) {
			bodyReader := bytes.NewReader(body)
			return "", HTTPDoWithValidationE(t, method, url, bodyReader, headers, expectedStatus, expectedBody, tlsConfig)
		})

	return err
}

// HTTPDoWithValidation performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, fail the test.
func HTTPDoWithValidation(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) {
	err := HTTPDoWithValidationE(t, method, url, body, headers, expectedStatusCode, expectedBody, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithValidationE performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, return an error.
func HTTPDoWithValidationE(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) error {
	return HTTPDoWithCustomValidationE(t, method, url, body, headers, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	}, tlsConfig)
}

// HTTPDoWithCustomValidation performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HTTPDoWithCustomValidation(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) {
	err := HTTPDoWithCustomValidationE(t, method, url, body, headers, validateResponse, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithCustomValidationE performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HTTPDoWithCustomValidationE(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) error {
	statusCode, respBody, err := HTTPDoE(t, method, url, body, headers, tlsConfig)

	if err != nil {
		return err
	}

	if !validateResponse(statusCode, respBody) {
		return ValidationFunctionFailed{Url: url, Status: statusCode, Body: respBody}
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
