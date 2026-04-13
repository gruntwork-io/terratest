// Package http_helper contains helpers to interact with deployed resources through HTTP.
package http_helper //nolint:staticcheck // package name determined by directory

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// defaultTimeoutSeconds is the default timeout in seconds for HTTP requests.
const defaultTimeoutSeconds = 10

// HttpGetOptions defines options for HTTP GET requests.
type HttpGetOptions struct { //nolint:staticcheck,revive // preserving existing type name
	// Context for the HTTP request. If nil, context.Background() is used.
	Context   context.Context
	TlsConfig *tls.Config //nolint:staticcheck,revive // preserving existing field name
	Url       string      //nolint:staticcheck,revive // preserving existing field name
	Timeout   int
}

// HttpDoOptions defines options for HTTP requests using an arbitrary method.
type HttpDoOptions struct { //nolint:staticcheck,revive // preserving existing type name
	// Context for the HTTP request. If nil, context.Background() is used.
	Context   context.Context
	Body      io.Reader
	Headers   map[string]string
	TlsConfig *tls.Config //nolint:staticcheck,revive // preserving existing field name
	Method    string
	Url       string //nolint:staticcheck,revive // preserving existing field name
	Timeout   int
}

// optionsContext returns the context from an HttpGetOptions, defaulting to context.Background() if nil.
func optionsContext(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}

	return context.Background()
}

// HTTPGetContext performs an HTTP GET on the given URL with an optional custom TLS configuration and returns the HTTP
// status code and body. The provided context is used for the HTTP request. If there's any error, fail the test.
func HTTPGetContext(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config) (int, string) {
	return HttpGetWithOptions(t, HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}) //nolint:contextcheck // context is passed via options struct
}

// HTTPGetContextE performs an HTTP GET on the given URL with an optional custom TLS configuration and returns the HTTP
// status code, body, and any error. The provided context is used for the HTTP request.
func HTTPGetContextE(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config) (int, string, error) {
	return HttpGetWithOptionsE(t, HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}) //nolint:contextcheck // context is passed via options struct
}

// HttpGet performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code and body. If there's any error, fail the test.
//
// Deprecated: Use [HTTPGetContext] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGet(t testing.TestingT, url string, tlsConfig *tls.Config) (int, string) {
	return HttpGetWithOptions(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds})
}

// HttpGetWithOptions performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code and body. If there's any error, fail the test.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithOptions(t testing.TestingT, options HttpGetOptions) (int, string) {
	statusCode, body, err := HttpGetWithOptionsE(t, options)
	if err != nil {
		t.Fatal(err)
	}

	return statusCode, body
}

// HttpGetE performs an HTTP GET, with an optional pointer to a custom TLS configuration, on the given URL and
// return the HTTP status code, body, and any error.
//
// Deprecated: Use [HTTPGetContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetE(t testing.TestingT, url string, tlsConfig *tls.Config) (int, string, error) {
	return HttpGetWithOptionsE(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds})
}

// HttpGetWithOptionsE performs an HTTP GET on the given URL with the given options and returns the HTTP status code,
// body, and any error.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithOptionsE(t testing.TestingT, options HttpGetOptions) (int, string, error) {
	logger.Default.Logf(t, "Making an HTTP GET call to URL %s", options.Url)

	ctx := optionsContext(options.Context)

	// Set HTTP client transport config
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = options.TlsConfig

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout: time.Duration(options.Timeout) * time.Second,
		// Include the previously created transport config
		Transport: tr,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, options.Url, nil)
	if err != nil {
		return -1, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(body)), nil
}

// HTTPGetWithValidationContext performs an HTTP GET on the given URL and verifies that the response has the expected
// status code and body. The provided context is used for the HTTP request. If either doesn't match, fail the test.
func HTTPGetWithValidationContext(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) {
	options := HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}
	HttpGetWithValidationWithOptions(t, options, expectedStatusCode, expectedBody) //nolint:contextcheck // context is passed via options struct
}

// HTTPGetWithValidationContextE performs an HTTP GET on the given URL and verifies that the response has the expected
// status code and body. The provided context is used for the HTTP request. If either doesn't match, return an error.
func HTTPGetWithValidationContextE(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) error {
	options := HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}

	return HttpGetWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody) //nolint:contextcheck // context is passed via options struct
}

// HttpGetWithValidation performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, fail the test.
//
// Deprecated: Use [HTTPGetWithValidationContext] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithValidation(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}
	HttpGetWithValidationWithOptions(t, options, expectedStatusCode, expectedBody)
}

// HttpGetWithValidationWithOptions performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, fail the test.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithValidationWithOptions(t testing.TestingT, options HttpGetOptions, expectedStatusCode int, expectedBody string) {
	err := HttpGetWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithValidationE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
//
// Deprecated: Use [HTTPGetWithValidationContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatusCode int, expectedBody string) error {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}

	return HttpGetWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
}

// HttpGetWithValidationWithOptionsE performs an HTTP GET on the given URL and verify that you get back the expected status code and body. If either
// doesn't match, return an error.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithValidationWithOptionsE(t testing.TestingT, options HttpGetOptions, expectedStatusCode int, expectedBody string) error {
	return HttpGetWithCustomValidationWithOptionsE(t, options, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HTTPGetWithCustomValidationContext performs an HTTP GET on the given URL and validates the returned status code and
// body using the given function. The provided context is used for the HTTP request.
func HTTPGetWithCustomValidationContext(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) {
	HttpGetWithCustomValidationWithOptions(t, HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}, validateResponse) //nolint:contextcheck // context is passed via options struct
}

// HTTPGetWithCustomValidationContextE performs an HTTP GET on the given URL and validates the returned status code and
// body using the given function. The provided context is used for the HTTP request.
func HTTPGetWithCustomValidationContextE(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) error {
	return HttpGetWithCustomValidationWithOptionsE(t, HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}, validateResponse) //nolint:contextcheck // context is passed via options struct
}

// HttpGetWithCustomValidation performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
//
// Deprecated: Use [HTTPGetWithCustomValidationContext] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithCustomValidation(t testing.TestingT, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) {
	HttpGetWithCustomValidationWithOptions(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}, validateResponse)
}

// HttpGetWithCustomValidationWithOptions performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithCustomValidationWithOptions(t testing.TestingT, options HttpGetOptions, validateResponse func(int, string) bool) {
	err := HttpGetWithCustomValidationWithOptionsE(t, options, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithCustomValidationE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
//
// Deprecated: Use [HTTPGetWithCustomValidationContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithCustomValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, validateResponse func(int, string) bool) error {
	return HttpGetWithCustomValidationWithOptionsE(t, HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}, validateResponse)
}

// HttpGetWithCustomValidationWithOptionsE performs an HTTP GET on the given URL and validate the returned status code and body using the given function.
//
//nolint:staticcheck,revive // preserving existing function name
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

// HTTPGetWithRetryContext repeatedly performs an HTTP GET on the given URL until the given status code and body are
// returned or until max retries has been exceeded. The provided context is used for each HTTP request.
func HTTPGetWithRetryContext(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	options := HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}
	HttpGetWithRetryWithOptions(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries) //nolint:contextcheck // context is passed via options struct
}

// HTTPGetWithRetryContextE repeatedly performs an HTTP GET on the given URL until the given status code and body are
// returned or until max retries has been exceeded. The provided context is used for each HTTP request.
func HTTPGetWithRetryContextE(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	options := HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}

	return HttpGetWithRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries) //nolint:contextcheck // context is passed via options struct
}

// HttpGetWithRetry repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
//
// Deprecated: Use [HTTPGetWithRetryContext] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetry(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}
	HttpGetWithRetryWithOptions(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HttpGetWithRetryWithOptions repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryWithOptions(t testing.TestingT, options HttpGetOptions, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) {
	err := HttpGetWithRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
//
// Deprecated: Use [HTTPGetWithRetryContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryE(t testing.TestingT, url string, tlsConfig *tls.Config, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}

	return HttpGetWithRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HttpGetWithRetryWithOptionsE repeatedly performs an HTTP GET on the given URL until the given status code and body are returned or until max
// retries has been exceeded.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryWithOptionsE(t testing.TestingT, options HttpGetOptions, expectedStatus int, expectedBody string, retries int, sleepBetweenRetries time.Duration) error {
	_, err := retry.DoWithRetryE(t, "HTTP GET to URL "+options.Url, retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithValidationWithOptionsE(t, options, expectedStatus, expectedBody)
	})

	return err
}

// HTTPGetWithRetryWithCustomValidationContext repeatedly performs an HTTP GET on the given URL until the given
// validation function returns true or max retries has been exceeded. The provided context is used for each HTTP request.
func HTTPGetWithRetryWithCustomValidationContext(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	options := HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}
	HttpGetWithRetryWithCustomValidationWithOptions(t, options, retries, sleepBetweenRetries, validateResponse) //nolint:contextcheck // context is passed via options struct
}

// HTTPGetWithRetryWithCustomValidationContextE repeatedly performs an HTTP GET on the given URL until the given
// validation function returns true or max retries has been exceeded. The provided context is used for each HTTP request.
func HTTPGetWithRetryWithCustomValidationContextE(t testing.TestingT, ctx context.Context, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	options := HttpGetOptions{Context: ctx, Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}

	return HttpGetWithRetryWithCustomValidationWithOptionsE(t, options, retries, sleepBetweenRetries, validateResponse) //nolint:contextcheck // context is passed via options struct
}

// HttpGetWithRetryWithCustomValidation repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
//
// Deprecated: Use [HTTPGetWithRetryWithCustomValidationContext] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryWithCustomValidation(t testing.TestingT, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}
	HttpGetWithRetryWithCustomValidationWithOptions(t, options, retries, sleepBetweenRetries, validateResponse)
}

// HttpGetWithRetryWithCustomValidationWithOptions repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryWithCustomValidationWithOptions(t testing.TestingT, options HttpGetOptions, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) {
	err := HttpGetWithRetryWithCustomValidationWithOptionsE(t, options, retries, sleepBetweenRetries, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HttpGetWithRetryWithCustomValidationE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
//
// Deprecated: Use [HTTPGetWithRetryWithCustomValidationContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryWithCustomValidationE(t testing.TestingT, url string, tlsConfig *tls.Config, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	options := HttpGetOptions{Url: url, TlsConfig: tlsConfig, Timeout: defaultTimeoutSeconds}

	return HttpGetWithRetryWithCustomValidationWithOptionsE(t, options, retries, sleepBetweenRetries, validateResponse)
}

// HttpGetWithRetryWithCustomValidationWithOptionsE repeatedly performs an HTTP GET on the given URL until the given validation function returns true or max retries
// has been exceeded.
//
//nolint:staticcheck,revive // preserving existing function name
func HttpGetWithRetryWithCustomValidationWithOptionsE(t testing.TestingT, options HttpGetOptions, retries int, sleepBetweenRetries time.Duration, validateResponse func(int, string) bool) error {
	_, err := retry.DoWithRetryE(t, "HTTP GET to URL "+options.Url, retries, sleepBetweenRetries, func() (string, error) {
		return "", HttpGetWithCustomValidationWithOptionsE(t, options, validateResponse)
	})

	return err
}

// HTTPDoContext performs the given HTTP method on the given URL and returns the HTTP status code and body.
// The provided context is used for the HTTP request. If there's any error, fail the test.
func HTTPDoContext(
	t testing.TestingT, ctx context.Context, method string, url string, body io.Reader,
	headers map[string]string, tlsConfig *tls.Config,
) (int, string) {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithOptions(t, options) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoContextE performs the given HTTP method on the given URL and returns the HTTP status code, body, and any error.
// The provided context is used for the HTTP request.
func HTTPDoContextE(
	t testing.TestingT, ctx context.Context, method string, url string, body io.Reader,
	headers map[string]string, tlsConfig *tls.Config,
) (int, string, error) {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithOptionsE(t, options) //nolint:contextcheck // context is passed via options struct
}

// HTTPDo performs the given HTTP method on the given URL and return the HTTP status code and body.
// If there's any error, fail the test.
//
// Deprecated: Use [HTTPDoContext] instead.
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
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithOptions(t, options)
}

// HTTPDoWithOptions performs the given HTTP method on the given URL and return the HTTP status code and body.
// If there's any error, fail the test.
//
//nolint:gocritic // cannot change public function signature
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
//
// Deprecated: Use [HTTPDoContextE] instead.
func HTTPDoE(
	t testing.TestingT, method string, url string, body io.Reader,
	headers map[string]string, tlsConfig *tls.Config,
) (int, string, error) {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithOptionsE(t, options)
}

// HTTPDoWithOptionsE performs the given HTTP method on the given URL and return the HTTP status code, body, and any error.
//
//nolint:gocritic // cannot change public function signature
func HTTPDoWithOptionsE(
	t testing.TestingT, options HttpDoOptions,
) (int, string, error) {
	logger.Default.Logf(t, "Making an HTTP %s call to URL %s", options.Method, options.Url)

	ctx := optionsContext(options.Context)

	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = options.TlsConfig

	client := http.Client{
		// By default, Go does not impose a timeout, so an HTTP connection attempt can hang for a LONG time.
		Timeout:   time.Duration(options.Timeout) * time.Second,
		Transport: tr,
	}

	req := newRequestWithContext(ctx, options.Method, options.Url, options.Body, options.Headers)

	resp, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}

	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, "", err
	}

	return resp.StatusCode, strings.TrimSpace(string(respBody)), nil
}

// HTTPDoWithRetryContext repeatedly performs the given HTTP method on the given URL until the given status code is
// returned or until max retries has been exceeded. The provided context is used for each HTTP request. The function
// compares the expected status code against the received one and fails if they don't match.
func HTTPDoWithRetryContext(
	t testing.TestingT, ctx context.Context, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) string {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithRetryWithOptions(t, options, expectedStatus, retries, sleepBetweenRetries) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithRetryContextE repeatedly performs the given HTTP method on the given URL until the given status code is
// returned or until max retries has been exceeded. The provided context is used for each HTTP request. The function
// compares the expected status code against the received one and fails if they don't match.
func HTTPDoWithRetryContextE(
	t testing.TestingT, ctx context.Context, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) (string, error) {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithRetryWithOptionsE(t, options, expectedStatus, retries, sleepBetweenRetries) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithRetry repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
//
// Deprecated: Use [HTTPDoWithRetryContext] instead.
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
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithRetryWithOptions(t, options, expectedStatus, retries, sleepBetweenRetries)
}

// HTTPDoWithRetryWithOptions repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
//
//nolint:gocritic // cannot change public function signature
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
//
// Deprecated: Use [HTTPDoWithRetryContextE] instead.
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
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithRetryWithOptionsE(t, options, expectedStatus, retries, sleepBetweenRetries)
}

// HTTPDoWithRetryWithOptionsE repeatedly performs the given HTTP method on the given URL until the given status code and body are
// returned or until max retries has been exceeded.
// The function compares the expected status code against the received one and fails if they don't match.
//
//nolint:gocritic // cannot change public function signature
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

	out, err := retry.DoWithRetryE( //nolint:staticcheck // deprecated wrapper; use HTTPDoWithRetryContext for context support
		t, "HTTP "+options.Method+" to URL "+options.Url, retries,
		sleepBetweenRetries, func() (string, error) {
			options.Body = bytes.NewReader(data)

			statusCode, out, err := HTTPDoWithOptionsE(t, options)
			if err != nil {
				return "", err
			}

			logger.Default.Logf(t, "output: %v", out)

			if statusCode != expectedStatus {
				return "", ValidationFunctionFailed{Url: options.Url, Status: statusCode}
			}

			return out, nil
		})

	return out, err
}

// HTTPDoWithValidationRetryContext repeatedly performs the given HTTP method on the given URL until the given status
// code and body are returned or until max retries has been exceeded. The provided context is used for each HTTP request.
func HTTPDoWithValidationRetryContext(
	t testing.TestingT, ctx context.Context, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	HTTPDoWithValidationRetryWithOptions(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithValidationRetryContextE repeatedly performs the given HTTP method on the given URL until the given status
// code and body are returned or until max retries has been exceeded. The provided context is used for each HTTP request.
func HTTPDoWithValidationRetryContextE(
	t testing.TestingT, ctx context.Context, method string, url string,
	body []byte, headers map[string]string, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration, tlsConfig *tls.Config,
) error {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      bytes.NewReader(body),
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithValidationRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithValidationRetry repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
//
// Deprecated: Use [HTTPDoWithValidationRetryContext] instead.
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
		Timeout:   defaultTimeoutSeconds,
	}

	HTTPDoWithValidationRetryWithOptions(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HTTPDoWithValidationRetryWithOptions repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
//
//nolint:gocritic // cannot change public function signature
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
//
// Deprecated: Use [HTTPDoWithValidationRetryContextE] instead.
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
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithValidationRetryWithOptionsE(t, options, expectedStatus, expectedBody, retries, sleepBetweenRetries)
}

// HTTPDoWithValidationRetryWithOptionsE repeatedly performs the given HTTP method on the given URL until the given status code and
// body are returned or until max retries has been exceeded.
//
//nolint:gocritic // cannot change public function signature
func HTTPDoWithValidationRetryWithOptionsE(
	t testing.TestingT, options HttpDoOptions, expectedStatus int,
	expectedBody string, retries int, sleepBetweenRetries time.Duration,
) error {
	_, err := retry.DoWithRetryE(t, "HTTP "+options.Method+" to URL "+options.Url, retries, //nolint:staticcheck // deprecated wrapper; use HTTPDoWithValidationRetryContext for context support
		sleepBetweenRetries, func() (string, error) {
			return "", HTTPDoWithValidationWithOptionsE(t, options, expectedStatus, expectedBody)
		})

	return err
}

// HTTPDoWithValidationContext performs the given HTTP method on the given URL and verifies that the response has the
// expected status code and body. The provided context is used for the HTTP request. If either doesn't match, fail
// the test.
func HTTPDoWithValidationContext(t testing.TestingT, ctx context.Context, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	HTTPDoWithValidationWithOptions(t, options, expectedStatusCode, expectedBody) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithValidationContextE performs the given HTTP method on the given URL and verifies that the response has the
// expected status code and body. The provided context is used for the HTTP request. If either doesn't match, return
// an error.
func HTTPDoWithValidationContextE(t testing.TestingT, ctx context.Context, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) error {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithValidation performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, fail the test.
//
// Deprecated: Use [HTTPDoWithValidationContext] instead.
func HTTPDoWithValidation(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	HTTPDoWithValidationWithOptions(t, options, expectedStatusCode, expectedBody)
}

// HTTPDoWithValidationWithOptions performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, fail the test.
//
//nolint:gocritic // cannot change public function signature
func HTTPDoWithValidationWithOptions(t testing.TestingT, options HttpDoOptions, expectedStatusCode int, expectedBody string) {
	err := HTTPDoWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithValidationE performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, return an error.
//
// Deprecated: Use [HTTPDoWithValidationContextE] instead.
func HTTPDoWithValidationE(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, expectedStatusCode int, expectedBody string, tlsConfig *tls.Config) error {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithValidationWithOptionsE(t, options, expectedStatusCode, expectedBody)
}

// HTTPDoWithValidationWithOptionsE performs the given HTTP method on the given URL and verify that you get back the expected status
// code and body. If either doesn't match, return an error.
//
//nolint:gocritic // cannot change public function signature
func HTTPDoWithValidationWithOptionsE(t testing.TestingT, options HttpDoOptions, expectedStatusCode int, expectedBody string) error {
	return HTTPDoWithCustomValidationWithOptionsE(t, options, func(statusCode int, body string) bool {
		return statusCode == expectedStatusCode && body == expectedBody
	})
}

// HTTPDoWithCustomValidationContext performs the given HTTP method on the given URL and validates the returned status
// code and body using the given function. The provided context is used for the HTTP request.
func HTTPDoWithCustomValidationContext(t testing.TestingT, ctx context.Context, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	HTTPDoWithCustomValidationWithOptions(t, options, validateResponse) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithCustomValidationContextE performs the given HTTP method on the given URL and validates the returned status
// code and body using the given function. The provided context is used for the HTTP request.
func HTTPDoWithCustomValidationContextE(t testing.TestingT, ctx context.Context, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) error {
	options := HttpDoOptions{
		Context:   ctx,
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithCustomValidationWithOptionsE(t, options, validateResponse) //nolint:contextcheck // context is passed via options struct
}

// HTTPDoWithCustomValidation performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
//
// Deprecated: Use [HTTPDoWithCustomValidationContext] instead.
func HTTPDoWithCustomValidation(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	HTTPDoWithCustomValidationWithOptions(t, options, validateResponse)
}

// HTTPDoWithCustomValidationWithOptions performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
//
//nolint:gocritic // cannot change public function signature
func HTTPDoWithCustomValidationWithOptions(t testing.TestingT, options HttpDoOptions, validateResponse func(int, string) bool) {
	err := HTTPDoWithCustomValidationWithOptionsE(t, options, validateResponse)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPDoWithCustomValidationE performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
//
// Deprecated: Use [HTTPDoWithCustomValidationContextE] instead.
func HTTPDoWithCustomValidationE(t testing.TestingT, method string, url string, body io.Reader, headers map[string]string, validateResponse func(int, string) bool, tlsConfig *tls.Config) error {
	options := HttpDoOptions{
		Method:    method,
		Url:       url,
		Body:      body,
		Headers:   headers,
		TlsConfig: tlsConfig,
		Timeout:   defaultTimeoutSeconds,
	}

	return HTTPDoWithCustomValidationWithOptionsE(t, options, validateResponse)
}

// HTTPDoWithCustomValidationWithOptionsE performs the given HTTP method on the given URL and validate the returned status code and
// body using the given function.
//
//nolint:gocritic // cannot change public function signature
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

func newRequestWithContext(ctx context.Context, method string, url string, body io.Reader, headers map[string]string) *http.Request {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
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
