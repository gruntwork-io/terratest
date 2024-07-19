// Package http_helper contains helpers to interact with deployed resources through HTTP.
package http_helper

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

type HttpGetOptions struct {
	Url       string
	TlsConfig *tls.Config
	Timeout   int
}

type HttpDoOptions struct {
	Method    string
	Url       string
	Body      io.Reader
	Headers   map[string]string
	TlsConfig *tls.Config
	Timeout   int
}

// HttpGet performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code and body. If there's any error, fail the test.
func HttpGet(t testing.TestingT, url string, tlsConfig *tls.Config) (int, string) {
	return HttpGetWithOptions(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10})
}

// HttpGetWithOptions performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code and body. If there's any error, fail the test.
func HttpGetWithOptions(t testing.TestingT, options HttpGetOptions) (int, string) {
	statusCode, body, err := HttpGetWithOptionsE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return statusCode, body
}

// HttpGetE performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code, body, and any error.
func HttpGetE(t testing.TestingT, url string, tlsConfig *tls.Config) (int, string, error) {
	return HttpGetWithOptionsE(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10})
}

// HttpGetWithOptionsE performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code, body, and any error.
func HttpGetWithOptionsE(t testing.TestingT, options HttpGetOptions) (int, string, error) {
	logger.Logf(t, "Making an HTTP GET call to URL %s", options.Url)

	// Set HTTP client transport config
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = options.TlsConfig

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: time.Duration(options.Timeout) * time.Second,
		// Include the previously created transport config
		Transport: tr,
	}

	resp, err := client.Get(options.Url)
	if err != nil {
		return -1, "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(body)), nil
}

// HttpGetWithValidation performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, fail the test.
func HttpGetWithValidation(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}
	HttpGetWithValidationWithOptions(t, options, expectedStatusCode, expectedBody)
}

// HttpGetWithValidationWithOptions performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, fail the test.
func HttpGetWithValidationWithOptions(t testing.TestingT, options HttpGetOptions, expectedStatusCode int, expectedBody string) {
	err := HttpGetWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithValidationE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
func HttpGetWithValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) error {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}
	return HttpGetWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
}

// HttpGetWithValidationWithOptionsE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
func HttpGetWithValidationWithOptionsE(t testing.TestingT, options HttpGetOptions, expectedStatusCode int, expectedBody string) error {
	return HttpGetWithCustomValidationWithOptionsE(t, options, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HttpGetWithCustomValidation performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidation(t testing.TestingT, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) {
	HttpGetWithCustomValidationWithOptions(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}, validateResponse)
}

// HttpGetWithCustomValidationWithOptions performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidationWithOptions(t testing.TestingT, options HttpGetOptions, validateResponse func(int, string) bool) {
	err := HttpGetWithCustomValidationWithOptionsE(t, options, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithCustomValidationE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) error {
	return HttpGetWithCustomValidationWithOptionsE(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}, validateResponse)
}

// HttpGetWithCustomValidationWithOptionsE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
func HttpGetWithCustomValidationWithOptionsE(t testing.TestingT, options HttpGetOptions, validateResponse func(int, string) bool) error {
	statusCode, body, err := HttpGetWithOptionsE(t, options)

	if err != nil {
		return err
	}

	if !validateResponse(statusCode, body) {
		return ValidationFunctionFailed{Url: options.Url, Status: statusCode, Body: body}
	}

	return nil
}

// HttpGetWithRetry repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetry(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}
	HttpGetWithRetryWithOptions(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HttpGetWithRetryWithOptions repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryWithOptions(t testing.TestingT, options HttpGetOptions, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpGetWithRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryE(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}
	return HttpGetWithRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HttpGetWithRetryWithOptionsE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
func HttpGetWithRetryWithOptionsE(t testing.TestingT, options HttpGetOptions, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", options.Url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithValidationWithOptionsE(t, options, expectedStatus, expectedBody)
	})

	return err
}

// HttpGetWithRetryWithCustomValidation repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidation(t testing.TestingT, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}
	HttpGetWithRetryWithCustomValidationWithOptions(t, options, retries, sleepBetweenRetries, validateResponse)
}

// HttpGetWithRetryWithCustomValidationWithOptions repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationWithOptions(t testing.TestingT, options HttpGetOptions, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	err := HttpGetWithRetryWithCustomValidationWithOptionsE(t, options, retries, sleepBetweenRetries, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryWithCustomValidationE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: 10}
	return HttpGetWithRetryWithCustomValidationWithOptionsE(t, options, retries, sleepBetweenRetries, validateResponse)
}

// HttpGetWithRetryWithCustomValidationWithOptionsE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
func HttpGetWithRetryWithCustomValidationWithOptionsE(t testing.TestingT, options HttpGetOptions, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP GET to URL %s", options.Url), retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithCustomValidationWithOptionsE(t, options, validateResponse)
	})

	return err
}

// HTTPDo performs the given HTTP method on the given URL and return the HTTP status code and body.
// If there's any error, fail the test.
func HTTPDo(
	t testing.TestingT, method string, url string, body io.Reader,
	headers map[string]string, tlsConfig *tls.Config,
) (int, string) {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}
	return HTTPDoWithOptions(t, options)
}

// HTTPDoWithOptions performs the given HTTP method on the given URL and return the HTTP status code and body.
// If there's any error, fail the test.
func HTTPDoWithOptions(
	t testing.TestingT, options HttpDoOptions,
) (int, string) {
	statusCode, respBody, err := HTTPDoWithOptionsE(t, options)
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
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		Timeout:   10,
		TlsConfig: tlsConfig}
	return HTTPDoWithOptionsE(t, options)
}

// HTTPDoWithOptionsE performs the given HTTP method on the given URL and return the HTTP status code, body, and any error.
func HTTPDoWithOptionsE(
	t testing.TestingT, options HttpDoOptions,
) (int, string, error) {
	logger.Logf(t, "Making an HTTP %s call to URL %s", options.Method, options.Url)

	// Set HTTP client transport config
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = tlsConfig

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout:   time.Duration(options.Timeout) * time.Second,
		Transport: tr,
	}

	req := newRequest(options.Method, options.Url, options.Body, options.Headers)
	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

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
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}
	return HTTPDoWithRetryWithOptions(t, options, expectedStatus, retries, sleepBetweenRetries)
}

// HTTPDoWithRetryWithOptions repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
func HTTPDoWithRetryWithOptions(
	t testing.TestingT, options HttpDoOptions, expectedStatus int,
	retries int, sleepBetweenRetries time.Duration,
) string {
	out, err := HTTPDoWithRetryWithOptionsE(t, options, expectedStatus, retries, sleepBetweenRetries)
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
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	return HTTPDoWithRetryWithOptionsE(t, options, expectedStatus, retries, sleepBetweenRetries)
}

// HTTPDoWithRetryWithOptionsE repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
func HTTPDoWithRetryWithOptionsE(
	t testing.TestingT, options HttpDoOptions, expectedStatus int,
	retries int, sleepBetweenRetries time.Duration,
) (string, error) {
	var data []byte
	if options.Body != nil {
		// The request body is closed after a request is complete.
		// Read the underlying data and cache it, so we can reuse for retried requests.
		b, err := io.ReadAll(options.Body)
		if err != nil {
			return "", err
		}
		data = b
	}

	options.Body = nil

	out, err := retry.DoWithRetryE(
		t, fmt.Sprintf("HTTP %s to URL %s", options.Method, options.Url), retries,
		sleepBetweenRetries, func() (string, error) {
			options.Body = bytes.NewReader(data)
			statusCode, out, err := HTTPDoWithOptionsE(t, options)
			if err != nil {
				return "", err
			}
			logger.Logf(t, "output: %v", out)
			if statusCode != expectedStatus {
				return "", ValidationFunctionFailed{Url: options.Url, Status: statusCode}
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
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	HTTPDoWithValidationRetryWithOptions(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HTTPDoWithValidationRetryWithOptions repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
func HTTPDoWithValidationRetryWithOptions(
	t testing.TestingT, options HttpDoOptions, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration,
) {
	err := HTTPDoWithValidationRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
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
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	return HTTPDoWithValidationRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HTTPDoWithValidationRetryWithOptionsE repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
func HTTPDoWithValidationRetryWithOptionsE(
	t testing.TestingT, options HttpDoOptions, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration,
) error {
	_, err := retry.DoWithRetryE(t, fmt.Sprintf("HTTP %s to URL %s", options.Method, options.Url), retries,
		sleepBetweenRetries, func() (string, error) {
			return "", HTTPDoWithValidationWithOptionsE(t, options, expectedStatus, expectedBody)
		})

	return err
}

// HTTPDoWithValidation performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, fail the test.
func HTTPDoWithValidation(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	HTTPDoWithValidationWithOptions(t, options, expectedStatusCode, expectedBody)
}

// HTTPDoWithValidationWithOptions performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, fail the test.
func HTTPDoWithValidationWithOptions(t testing.TestingT, options HttpDoOptions, expectedStatusCode int, expectedBody string) {
	err := HTTPDoWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithValidationE performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, return an error.
func HTTPDoWithValidationE(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) error {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	return HTTPDoWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
}

// HTTPDoWithValidationWithOptionsE performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, return an error.
func HTTPDoWithValidationWithOptionsE(t testing.TestingT, options HttpDoOptions, expectedStatusCode int, expectedBody string) error {
	return HTTPDoWithCustomValidationWithOptionsE(t, options, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HTTPDoWithCustomValidation performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HTTPDoWithCustomValidation(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	HTTPDoWithCustomValidationWithOptions(t, options, validateResponse)
}

// HTTPDoWithCustomValidationWithOptions performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HTTPDoWithCustomValidationWithOptions(t testing.TestingT, options HttpDoOptions, validateResponse func(int, string) bool) {
	err := HTTPDoWithCustomValidationWithOptionsE(t, options, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithCustomValidationE performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HTTPDoWithCustomValidationE(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) error {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   10}

	return HTTPDoWithCustomValidationWithOptionsE(t, options, validateResponse)
}

// HTTPDoWithCustomValidationWithOptionsE performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
func HTTPDoWithCustomValidationWithOptionsE(t testing.TestingT, options HttpDoOptions, validateResponse func(int, string) bool) error {
	statusCode, respBody, err := HTTPDoWithOptionsE(t, options)

	if err != nil {
		return err
	}

	if !validateResponse(statusCode, respBody) {
		return ValidationFunctionFailed{Url: options.Url, Status: statusCode, Body: respBody}
	}

	return nil
}

func newRequest(method string, url string, body io.Reader, headers map[string]string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	for k, v := range headers {
		switch k {
		case "Host":
			req.Host = v
		default:
			req.Header.Add(k, v)
		}
	}
	return req
}
