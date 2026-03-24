package http_helper_test //nolint:staticcheck // package name determined by directory

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	httphelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestServerForFunction(handler func(w http.ResponseWriter,
	r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestOkBody(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(bodyCopyHandler)
	defer ts.Close()

	url := ts.URL
	expectedBody := "Hello, Terratest!"
	body := bytes.NewReader([]byte(expectedBody))
	statusCode, respBody := httphelper.HTTPDo(t, "POST", url, body, nil, nil)

	expectedCode := 200

	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, expectedCode)
	}

	if respBody != expectedBody {
		t.Errorf("handler returned wrong body: got %v want %v", respBody, expectedBody)
	}
}

func TestHTTPDoWithValidation(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(bodyCopyHandler)
	defer ts.Close()

	url := ts.URL
	expectedBody := "Hello, Terratest!"
	body := bytes.NewReader([]byte(expectedBody))
	httphelper.HTTPDoWithValidation(t, "POST", url, body, nil, 200, expectedBody, nil)
}

func TestHTTPDoWithCustomValidation(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(bodyCopyHandler)
	defer ts.Close()

	url := ts.URL
	expectedBody := "Hello, Terratest!"
	body := bytes.NewReader([]byte(expectedBody))

	customValidation := func(statusCode int, response string) bool {
		return statusCode == 200 && response == expectedBody
	}

	httphelper.HTTPDoWithCustomValidation(t, "POST", url, body, nil, customValidation, nil)
}

func TestOkHeaders(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(headersCopyHandler)
	defer ts.Close()

	url := ts.URL
	headers := map[string]string{"Authorization": "Bearer 1a2b3c99ff"}
	statusCode, respBody := httphelper.HTTPDo(t, "POST", url, nil, headers, nil)

	expectedCode := 200

	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, expectedCode)
	}

	expectedLine := "Authorization: Bearer 1a2b3c99ff"

	if !strings.Contains(respBody, expectedLine) {
		t.Errorf("handler returned wrong body: got %v want %v", respBody, expectedLine)
	}
}

func TestWrongStatus(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(wrongStatusHandler)
	defer ts.Close()

	url := ts.URL
	statusCode, _ := httphelper.HTTPDo(t, "POST", url, nil, nil, nil)

	expectedCode := 500

	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, expectedCode)
	}
}

func TestRequestTimeout(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(sleepingHandler)
	defer ts.Close()

	url := ts.URL

	_, _, err := httphelper.HTTPDoE(t, "DELETE", url, nil, nil, nil)
	if err == nil {
		t.Error("handler didn't return a timeout error")
	}

	if !strings.Contains(err.Error(), "Client.Timeout") {
		t.Errorf("handler didn't return an expected error, got %q", err)
	}
}

func TestOkWithRetry(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(retryHandler)
	defer ts.Close()

	body := "TEST_CONTENT"
	bodyBytes := []byte(body)

	url := ts.URL
	counter = 3
	response := httphelper.HTTPDoWithRetry(t, "POST", url, bodyBytes, nil, 200, 10, time.Second, nil)
	require.Equal(t, body, response)
}

func TestErrorWithRetry(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(failRetryHandler)
	defer ts.Close()

	failCounter = 3

	url := ts.URL

	_, err := httphelper.HTTPDoWithRetryE(t, "POST", url, nil, nil, 200, 2, time.Second, nil)
	if err == nil {
		t.Error("handler didn't return a retry error")
	}

	pattern := `unsuccessful after \d+ retries`

	match, _ := regexp.MatchString(pattern, err.Error())
	if !match {
		t.Errorf("handler didn't return an expected error, got %q", err)
	}
}

func TestEmptyRequestBodyWithRetryWithOptions(t *testing.T) {
	t.Parallel()

	ts := getTestServerForFunction(bodyCopyHandler)
	defer ts.Close()

	options := httphelper.HttpDoOptions{
		Method: "GET",
		Url:    ts.URL,
		Body:   nil,
	}

	response := httphelper.HTTPDoWithRetryWithOptions(t, options, 200, 0, time.Second)
	require.Empty(t, response)
}

func bodyCopyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	body, _ := io.ReadAll(r.Body)

	_, _ = w.Write(body)
}

func headersCopyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	var buffer bytes.Buffer

	for key, values := range r.Header {
		fmt.Fprintf(&buffer, "%s: %s\n", key, strings.Join(values, ","))
	}

	_, _ = w.Write(buffer.Bytes())
}

func wrongStatusHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func sleepingHandler(_ http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Second * 15)
}

var counter int

func retryHandler(w http.ResponseWriter, r *http.Request) {
	if counter > 0 {
		counter--

		w.WriteHeader(http.StatusServiceUnavailable)

		_, _ = io.ReadAll(r.Body)
	} else {
		w.WriteHeader(http.StatusOK)

		bytes, _ := io.ReadAll(r.Body)

		_, _ = w.Write(bytes)
	}
}

var failCounter int

func failRetryHandler(w http.ResponseWriter, r *http.Request) {
	if failCounter > 0 {
		failCounter--

		w.WriteHeader(http.StatusServiceUnavailable)

		_, _ = io.ReadAll(r.Body)
	} else {
		w.WriteHeader(http.StatusOK)

		bytes, _ := io.ReadAll(r.Body)

		_, _ = w.Write(bytes)
	}
}

func TestGlobalProxy(t *testing.T) {
	proxiedURL := ""

	httpProxy := getTestServerForFunction(func(w http.ResponseWriter, r *http.Request) {
		proxiedURL = r.RequestURI
		bodyCopyHandler(w, r)
	})
	t.Cleanup(httpProxy.Close)

	t.Setenv("HTTP_PROXY", httpProxy.URL)

	targetURL := "http://www.notexist.com/"
	body := "should be copied"

	st, b, err := httphelper.HTTPDoWithOptionsE(t, httphelper.HttpDoOptions{
		Url:    targetURL,
		Method: http.MethodPost,
		Body:   strings.NewReader(body),
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, st)
	assert.Equal(t, targetURL, proxiedURL)
	assert.Equal(t, body, b)
}
