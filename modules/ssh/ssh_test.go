package ssh

import (
	"errors"
	"fmt"
	"testing"

	grunttest "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
)

func TestHostWithDefaultPort(t *testing.T) {
	t.Parallel()

	host := Host{}

	assert.Equal(t, 22, host.getPort(), "host.getPort() did not return the default ssh port of 22")
}

func TestHostWithCustomPort(t *testing.T) {
	t.Parallel()

	customPort := 2222
	host := Host{CustomPort: customPort}

	assert.Equal(t, customPort, host.getPort(), "host.getPort() did not return the custom port number")
}

// global var for use in mock callback
var timesCalled int

func TestCheckSshConnectionWithRetryE(t *testing.T) {
	// Reset the global call count
	timesCalled = 0

	host := Host{Hostname: "Host"}
	retries := 10

	assert.Nil(t, CheckSshConnectionWithRetryE(t, host, retries, 3, mockSshConnectionE))
}

func TestCheckSshConnectionWithRetryEExceedsMaxRetries(t *testing.T) {
	// Reset the global call count
	timesCalled = 0

	host := Host{Hostname: "Host"}

	// Not enough retries
	retries := 3

	assert.Error(t, CheckSshConnectionWithRetryE(t, host, retries, 3, mockSshConnectionE))
}

func TestCheckSshConnectionWithRetry(t *testing.T) {
	// Reset the global call count
	timesCalled = 0

	host := Host{Hostname: "Host"}
	retries := 10

	CheckSshConnectionWithRetry(t, host, retries, 3, mockSshConnectionE)
}

func TestCheckSshCommandWithRetryE(t *testing.T) {
	// Reset the global call count
	timesCalled = 0

	host := Host{Hostname: "Host"}
	command := "echo -n hello world"
	retries := 10

	_, err := CheckSshCommandWithRetryE(t, host, command, retries, 3, mockSshCommandE)
	assert.Nil(t, err)
}

func TestCheckSshCommandWithRetryEExceedsRetries(t *testing.T) {
	// Reset the global call count
	timesCalled = 0

	host := Host{Hostname: "Host"}
	command := "echo -n hello world"

	// Not enough retries
	retries := 3

	_, err := CheckSshCommandWithRetryE(t, host, command, retries, 3, mockSshCommandE)
	assert.Error(t, err)
}

func TestCheckSshCommandWithRetry(t *testing.T) {
	// Reset the global call count
	timesCalled = 0

	host := Host{Hostname: "Host"}
	command := "echo -n hello world"
	retries := 10

	CheckSshCommandWithRetry(t, host, command, retries, 3, mockSshCommandE)
}

func mockSshConnectionE(t grunttest.TestingT, host Host) error {
	timesCalled += 1
	if timesCalled >= 5 {
		return nil
	} else {
		return errors.New(fmt.Sprintf("Called %v times", timesCalled))
	}
}

func mockSshCommandE(t grunttest.TestingT, host Host, command string) (string, error) {
	return "", mockSshConnectionE(t, host)
}
