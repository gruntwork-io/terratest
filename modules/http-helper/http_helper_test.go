package http_helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestHttpDoBody(t *testing.T) {
	t.Parallel()

	path := "/doBody"
	listener, port := RunDummyHandlerServer(t, path, bodyCopyHandler)
	defer shutDownServer(t, listener)

	url := fmt.Sprintf("http://localhost:%d%s", port, path)
	expectedBody := "Hello, Terratest!"
	body := bytes.NewReader([]byte(expectedBody))
	statusCode, respBody := HttpDo(t, "POST", url, body, nil)

	expectedCode := 200
	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %q want %q", statusCode, expectedCode)
	}
	if respBody != expectedBody {
		t.Errorf("handler returned wrong body: got %q want %q", respBody, expectedBody)
	}
}

func TestHttpDoHeaders(t *testing.T) {
	t.Parallel()

	path := "/doHeaders"
	listener, port := RunDummyHandlerServer(t, path, headersCopyHandler)
	defer shutDownServer(t, listener)

	url := fmt.Sprintf("http://localhost:%d%s", port, path)
	headers := map[string]string{"Authorization": "Bearer 1a2b3c99ff"}
	statusCode, respBody := HttpDo(t, "POST", url, nil, headers)

	expectedCode := 200
	if statusCode != expectedCode {
		t.Errorf("handler returned wrong status code: got %q want %q", statusCode, expectedCode)
	}
	expectedLine := "Authorization: Bearer 1a2b3c99ff"
	if !strings.Contains(respBody, expectedLine) {
		t.Errorf("handler returned wrong body: got %q want %q", respBody, expectedLine)
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
	time.Sleep(time.Second)
}
