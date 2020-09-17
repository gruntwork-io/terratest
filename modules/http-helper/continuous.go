package http_helper

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

type GetResponse struct {
	StatusCode int
	Body       string
}

// Continuously check the given URL every 1 second until the stopChecking channel receives a signal to stop.
// This function will return a sync.WaitGroup that can be used to wait for the checking to stop, and a read only channel
// to stream the responses for each check.
// Note that the channel has a buffer of 1000, after which it will start to drop the send events
func ContinuouslyCheckUrl(
	t testing.TestingT,
	url string,
	stopChecking <-chan bool,
	sleepBetweenChecks time.Duration,
) (*sync.WaitGroup, <-chan GetResponse) {
	var wg sync.WaitGroup
	wg.Add(1)
	responses := make(chan GetResponse, 1000)
	go func() {
		defer wg.Done()
		defer close(responses)
		for {
			select {
			case <-stopChecking:
				logger.Logf(t, "Got signal to stop downtime checks for URL %s.\n", url)
				return
			case <-time.After(sleepBetweenChecks):
				statusCode, body, err := HttpGetE(t, url, &tls.Config{})
				// Non-blocking send, defaulting to logging a warning if there is no channel reader
				select {
				case responses <- GetResponse{StatusCode: statusCode, Body: body}:
					// do nothing since all we want to do is send the response
				default:
					logger.Logf(t, "WARNING: ContinuouslyCheckUrl responses channel buffer is full")
				}
				logger.Logf(t, "Got response %d and err %v from URL at %s", statusCode, err, url)
				if err != nil {
					// We use Errorf instead of Fatalf here because Fatalf is not goroutine safe, while Errorf is. Refer
					// to the docs on `T`: https://godoc.org/testing#T
					t.Errorf("Failed to make HTTP request to the URL at %s: %s\n", url, err.Error())
				} else if statusCode != 200 {
					// We use Errorf instead of Fatalf here because Fatalf is not goroutine safe, while Errorf is. Refer
					// to the docs on `T`: https://godoc.org/testing#T
					t.Errorf("Got a non-200 response (%d) from the URL at %s, which means there was downtime! Response body: %s", statusCode, url, body)
				}
			}
		}
	}()
	return &wg, responses
}
