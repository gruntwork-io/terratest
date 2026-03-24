package http_helper_test //nolint:staticcheck // package name determined by directory

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	httphelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestRunDummyServer(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	text := "dummy-server-" + uniqueID

	listener, port := httphelper.RunDummyServer(t, text)
	defer shutDownServer(t, listener)

	url := fmt.Sprintf("http://localhost:%d", port)
	httphelper.HttpGetWithValidation(t, url, &tls.Config{}, 200, text)
}

func TestContinuouslyCheck(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	text := "dummy-server-" + uniqueID
	stopChecking := make(chan bool, 1)

	listener, port := httphelper.RunDummyServer(t, text)

	url := fmt.Sprintf("http://localhost:%d", port)
	wg, responses := httphelper.ContinuouslyCheckUrl(t, url, stopChecking, 1*time.Second)

	defer func() {
		stopChecking <- true

		counts := 0

		for response := range responses {
			counts++

			assert.Equal(t, 200, response.StatusCode)
			assert.Equal(t, text, response.Body)
		}

		wg.Wait()
		// Make sure we made at least one call
		assert.NotEqual(t, 0, counts)
		shutDownServer(t, listener)
	}()

	time.Sleep(5 * time.Second)
}

func TestRunDummyServersWithHandlers(t *testing.T) {
	// Given:
	//   several dummy servers, each with the same path
	// When:
	//   all of them are started at the same time
	// Then:
	//   every one of them can be started and serves their unique content
	t.Parallel()

	numServers := 2

	type testData struct {
		text string
		port int
	}

	data := make([]testData, numServers)

	for idx := 0; idx < numServers; idx++ {
		uniqueID := random.UniqueID()
		text := "dummy-server-" + uniqueID

		handler := func(w http.ResponseWriter, r *http.Request) {
			_, _ = io.WriteString(w, text)
		}

		handlerMap := map[string]func(http.ResponseWriter, *http.Request){
			// The same endpoint is provided for each dummy server.
			"/v1/endpoint": handler,
		}

		listener, port := httphelper.RunDummyServerWithHandlers(t, handlerMap)
		defer shutDownServer(t, listener)

		data[idx] = testData{text: text, port: port}
	}

	for _, testInstance := range data {
		url := fmt.Sprintf("http://localhost:%d/v1/endpoint", testInstance.port)
		httphelper.HttpGetWithValidation(t, url, &tls.Config{}, 200, testInstance.text)
	}
}

func shutDownServer(t *testing.T, listener io.Closer) {
	t.Helper()

	err := listener.Close()
	assert.NoError(t, err)
}
