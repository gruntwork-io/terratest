package http_helper

import (
	"fmt"
	"io"
	"testing"

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
	HttpGetWithValidation(t, url, 200, text)
}

func shutDownServer(t *testing.T, listener io.Closer) {
	err := listener.Close()
	assert.NoError(t, err)
}
