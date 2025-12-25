package http_helper

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random/v2"
	"github.com/stretchr/testify/assert"
)

func TestRunDummyServer(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueId()
	text := fmt.Sprintf("dummy-server-%s", uniqueID)

	listener, port := RunDummyServer(t, text)
	defer shutDownServer(t, listener)

	url := fmt.Sprintf("http://localhost:%d", port)
	HttpGetWithValidation(t, url, &tls.Config{}, 200, text)
}

func TestContinuouslyCheck(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueId()
	text := fmt.Sprintf("dummy-server-%s", uniqueID)
	stopChecking := make(chan bool, 1)

	listener, port := RunDummyServer(t, text)

	url := fmt.Sprintf("http://localhost:%d", port)
	wg, responses := ContinuouslyCheckUrl(t, url, stopChecking, 1*time.Second)
	defer func() {
		stopChecking <- true
		counts := 0
		for response := range responses {
			counts++
			assert.Equal(t, response.StatusCode, 200)
			assert.Equal(t, response.Body, text)
		}
		wg.Wait()
		// Make sure we made at least one call
		assert.NotEqual(t, counts, 0)
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
		uniqueID := random.UniqueId()
		text := fmt.Sprintf("dummy-server-%s", uniqueID)

		handler := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "%s", text)
		}

		handlerMap := map[string]func(http.ResponseWriter, *http.Request){
			// The same endpoint is provided for each dummy server.
			"/v1/endpoint": handler,
		}

		listener, port := RunDummyServerWithHandlers(t, handlerMap)
		defer shutDownServer(t, listener)

		data[idx] = testData{text: text, port: port}
	}

	for _, testInstance := range data {
		url := fmt.Sprintf("http://localhost:%d/v1/endpoint", testInstance.port)
		HttpGetWithValidation(t, url, &tls.Config{}, 200, testInstance.text)
	}
}

func shutDownServer(t *testing.T, listener io.Closer) {
	err := listener.Close()
	assert.NoError(t, err)
}
