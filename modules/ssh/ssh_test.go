package ssh_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	grunttest "github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/ssh/v2"
	"github.com/stretchr/testify/assert"
)

func TestHostWithDefaultPort(t *testing.T) {
	t.Parallel()

	host := ssh.Host{}

	assert.Equal(t, 22, host.GetPort(), "host.GetPort() did not return the default ssh port of 22")
}

func TestHostWithCustomPort(t *testing.T) {
	t.Parallel()

	customPort := 2222
	host := ssh.Host{CustomPort: customPort}

	assert.Equal(t, customPort, host.GetPort(), "host.GetPort() did not return the custom port number")
}

func TestCheckSshConnectionWithRetryE(t *testing.T) {
	t.Parallel()

	timesCalled := 0

	handler := func(_ grunttest.TestingT, _ context.Context, _ *ssh.Host) error {
		timesCalled++

		if timesCalled >= 5 {
			return nil
		}

		return fmt.Errorf("called %v times", timesCalled)
	}

	host := ssh.Host{Hostname: "Host"}
	retries := 10

	assert.NoError(t, ssh.CheckSSHConnectionWithRetryContextE(t, context.Background(), &host, retries, 3*time.Millisecond, handler))
}

func TestCheckSshConnectionWithRetryEExceedsMaxRetries(t *testing.T) {
	t.Parallel()

	timesCalled := 0

	handler := func(_ grunttest.TestingT, _ context.Context, _ *ssh.Host) error {
		timesCalled++

		if timesCalled >= 5 {
			return nil
		}

		return fmt.Errorf("called %v times", timesCalled)
	}

	host := ssh.Host{Hostname: "Host"}

	// Not enough retries.
	retries := 3

	assert.Error(t, ssh.CheckSSHConnectionWithRetryContextE(t, context.Background(), &host, retries, 3*time.Millisecond, handler))
}

func TestCheckSshConnectionWithRetry(t *testing.T) {
	t.Parallel()

	timesCalled := 0

	handler := func(_ grunttest.TestingT, _ context.Context, _ *ssh.Host) error {
		timesCalled++

		if timesCalled >= 5 {
			return nil
		}

		return fmt.Errorf("called %v times", timesCalled)
	}

	host := ssh.Host{Hostname: "Host"}
	retries := 10

	ssh.CheckSSHConnectionWithRetryContext(t, context.Background(), &host, retries, 3*time.Millisecond, handler)
}

func TestCheckSshCommandWithRetryE(t *testing.T) {
	t.Parallel()

	timesCalled := 0

	handler := func(_ grunttest.TestingT, _ context.Context, _ *ssh.Host, _ string) (string, error) {
		timesCalled++

		if timesCalled >= 5 {
			return "", nil
		}

		return "", fmt.Errorf("called %v times", timesCalled)
	}

	host := ssh.Host{Hostname: "Host"}
	command := "echo -n hello world"
	retries := 10

	_, err := ssh.CheckSSHCommandWithRetryContextE(t, context.Background(), &host, command, retries, 3*time.Millisecond, handler)
	assert.NoError(t, err)
}

func TestCheckSshCommandWithRetryEExceedsRetries(t *testing.T) {
	t.Parallel()

	timesCalled := 0

	handler := func(_ grunttest.TestingT, _ context.Context, _ *ssh.Host, _ string) (string, error) {
		timesCalled++

		if timesCalled >= 5 {
			return "", nil
		}

		return "", fmt.Errorf("called %v times", timesCalled)
	}

	host := ssh.Host{Hostname: "Host"}
	command := "echo -n hello world"

	// Not enough retries.
	retries := 3

	_, err := ssh.CheckSSHCommandWithRetryContextE(t, context.Background(), &host, command, retries, 3*time.Millisecond, handler)
	assert.Error(t, err)
}

func TestCheckSshCommandWithRetry(t *testing.T) {
	t.Parallel()

	timesCalled := 0

	handler := func(_ grunttest.TestingT, _ context.Context, _ *ssh.Host, _ string) (string, error) {
		timesCalled++

		if timesCalled >= 5 {
			return "", nil
		}

		return "", fmt.Errorf("called %v times", timesCalled)
	}

	host := ssh.Host{Hostname: "Host"}
	command := "echo -n hello world"
	retries := 10

	ssh.CheckSSHCommandWithRetryContext(t, context.Background(), &host, command, retries, 3*time.Millisecond, handler)
}
