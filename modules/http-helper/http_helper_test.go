package http_helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

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
	statusCode, respBody := HTTPDo(t, "POST", url, body, nil, nil)

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
	HTTPDoWithValidation(t, "POST", url, body, nil, 200, expectedBody, nil)
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

	HTTPDoWithCustomValidation(t, "POST", url, body, nil, customValidation, nil)
}

func TestOkHeaders(t *testing.T) {
	t.Parallel()
	ts := getTestServerForFunction(headersCopyHandler)
	defer ts.Close()
	url := ts.URL
	headers := map[string]string{"Authorization": "Bearer 1a2b3c99ff"}
	statusCode, respBody := HTTPDo(t, "POST", url, nil, headers, nil)

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
	statusCode, _ := HTTPDo(t, "POST", url, nil, nil, nil)

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
	_, _, err := HTTPDoE(t, "DELETE", url, nil, nil, nil)

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
	response := HTTPDoWithRetry(t, "POST", url, bodyBytes, nil, 200, 10, time.Second, nil)
	require.Equal(t, body, response)
}

func TestErrorWithRetry(t *testing.T) {
	t.Parallel()
	ts := getTestServerForFunction(failRetryHandler)
	defer ts.Close()
	failCounter = 3
	url := ts.URL
	_, err := HTTPDoWithRetryE(t, "POST", url, nil, nil, 200, 2, time.Second, nil)

	if err == nil {
		t.Error("handler didn't return a retry error")
	}

	pattern := `unsuccessful after \d+ retries`
	match, _ := regexp.MatchString(pattern, err.Error())
	if !match {
		t.Errorf("handler didn't return an expected error, got %q", err)
	}
}

func bodyCopyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(r.Body)
	w.Write(body)
}

func headersCopyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	var buffer bytes.Buffer
	for key, values := range r.Header {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(values, ",")))
	}
	w.Write(buffer.Bytes())
}

func wrongStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func sleepingHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 15)
}

var counter int

func retryHandler(w http.ResponseWriter, r *http.Request) {
	if counter > 0 {
		counter--
		w.WriteHeader(http.StatusServiceUnavailable)
		ioutil.ReadAll(r.Body)
	} else {
		w.WriteHeader(http.StatusOK)
		bytes, _ := ioutil.ReadAll(r.Body)
		w.Write(bytes)
	}
}

var failCounter int

func failRetryHandler(w http.ResponseWriter, r *http.Request) {
	if failCounter > 0 {
		failCounter--
		w.WriteHeader(http.StatusServiceUnavailable)
		ioutil.ReadAll(r.Body)
	} else {
		w.WriteHeader(http.StatusOK)
		bytes, _ := ioutil.ReadAll(r.Body)
		w.Write(bytes)
	}
}
