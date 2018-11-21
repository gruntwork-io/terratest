package http_helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	doBodyPath         = "/doBody"
	doHeadersPath      = "/doHeaders"
	wrongStatusPath    = "/wrongStatus"
	requestTimeoutPath = "/requestTimeout"
	retryPath          = "/retry"
)

var baseUrl string

func TestHttpDo(t *testing.T) {
	t.Parallel()

	handlers := map[string]func(http.ResponseWriter, *http.Request){
		doBodyPath:         bodyCopyHandler,
		doHeadersPath:      headersCopyHandler,
		wrongStatusPath:    wrongStatusHandler,
		requestTimeoutPath: sleepingHandler,
		retryPath:          retryHandler,
	}
	listener, port := RunDummyHandlerServer(t, handlers)
	defer shutDownServer(t, listener)

	baseUrl = fmt.Sprintf("http://localhost:%d", port)

	t.Run("okBody", okBody)
	t.Run("okHeaders", okHeaders)
	t.Run("wrongStatus", wrongStatus)
	t.Run("requestTimeout", requestTimeout)
	t.Run("okWithRetry", okWithRetry)
	t.Run("errorWithRetry", errorWithRetry)
}

func okBody(t *testing.T) {
	url := fmt.Sprintf("%s%s", baseUrl, doBodyPath)
	expectedBody := "Hello, Terratest!"
	body := bytes.NewReader([]byte(expectedBody))
	statusCode, respBody := HttpDo(t, "POST", url, body, nil)

	expectedCode := 200
	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, expectedCode)
	}
	if respBody != expectedBody {
		t.Errorf("handler returned wrong body: got %v want %v", respBody, expectedBody)
	}
}

func okHeaders(t *testing.T) {
	url := fmt.Sprintf("%s%s", baseUrl, doHeadersPath)
	headers := map[string]string{"Authorization": "Bearer 1a2b3c99ff"}
	statusCode, respBody := HttpDo(t, "POST", url, nil, headers)

	expectedCode := 200
	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, expectedCode)
	}
	expectedLine := "Authorization: Bearer 1a2b3c99ff"
	if !strings.Contains(respBody, expectedLine) {
		t.Errorf("handler returned wrong body: got %v want %v", respBody, expectedLine)
	}
}

func wrongStatus(t *testing.T) {
	url := fmt.Sprintf("%s%s", baseUrl, wrongStatusPath)
	statusCode, _ := HttpDo(t, "POST", url, nil, nil)

	expectedCode := 500
	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", statusCode, expectedCode)
	}
}

func requestTimeout(t *testing.T) {
	url := fmt.Sprintf("%s%s", baseUrl, requestTimeoutPath)
	_, _, err := HttpDoE(t, "DELETE", url, nil, nil)

	if err == nil {
		t.Error("handler didn't return a timeout error")
	}
	if !strings.Contains(err.Error(), "request canceled") {
		t.Errorf("handler didn't return an expected error, got %q", err)
	}
}

var counter int

func okWithRetry(t *testing.T) {
	counter = 3
	url := fmt.Sprintf("%s%s", baseUrl, retryPath)
	HttpDoWithRetry(t, "POST", url, nil, nil, 200, 10, time.Second)
}

func errorWithRetry(t *testing.T) {
	counter = 3
	url := fmt.Sprintf("%s%s", baseUrl, retryPath)
	_, err := HttpDoWithRetryE(t, "POST", url, nil, nil, 200, 2, time.Second)

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

func retryHandler(w http.ResponseWriter, r *http.Request) {
	if counter > 0 {
		counter--
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
