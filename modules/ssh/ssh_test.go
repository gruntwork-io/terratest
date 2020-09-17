package ssh

import (
	"testing"

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
