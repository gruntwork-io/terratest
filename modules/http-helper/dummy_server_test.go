package http_helper

import (
	"crypto/tls"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/random"
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

func shutDownServer(t *testing.T, listener io.Closer) {
	err := listener.Close()
	assert.NoError(t, err)
}
