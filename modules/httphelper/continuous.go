package httphelper

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// responseChannelBufferSize is the buffer size for the ContinuouslyCheckURL response channel.
const responseChannelBufferSize = 1000

// GetResponse represents the response from an HTTP GET request.
type GetResponse struct {
	Body       string
	StatusCode int
}

// ContinuouslyCheckURLContext continuously checks the given URL at the specified interval until the stopChecking
// channel receives a signal to stop. It returns a sync.WaitGroup that can be used to wait for the checking to stop,
// and a read-only channel to stream the responses for each check. The channel has a buffer of 1000 entries, after
// which it will start to drop send events. The provided context is used for each HTTP request made during checking.
func ContinuouslyCheckURLContext(
	t testing.TestingT,
	ctx context.Context,
	url string,
	stopChecking <-chan bool,
	sleepBetweenChecks time.Duration,
) (*sync.WaitGroup, <-chan GetResponse) {
	var wg sync.WaitGroup

	wg.Add(1)

	responses := make(chan GetResponse, responseChannelBufferSize)

	go func() {
		defer wg.Done()
		defer close(responses)

		for {
			select {
			case <-stopChecking:
				logger.Default.Logf(t, "Got signal to stop downtime checks for URL %s.\n", url)
				return
			case <-time.After(sleepBetweenChecks):
				statusCode, body, err := HTTPGetContextE(t, ctx, url, &tls.Config{})

				select {
				case responses <- GetResponse{StatusCode: statusCode, Body: body}:

				default:
					logger.Default.Logf(t, "WARNING: ContinuouslyCheckURLContext responses channel buffer is full")
				}

				logger.Default.Logf(t, "Got response %d and err %v from URL at %s", statusCode, err, url)

				if err != nil {

					t.Errorf("Failed to make HTTP request to the URL at %s: %s\n", url, err.Error())
				} else if statusCode != http.StatusOK {

					t.Errorf("Got a non-200 response (%d) from the URL at %s, which means there was downtime! Response body: %s", statusCode, url, body)
				}
			}
		}
	}()

	return &wg, responses
}
