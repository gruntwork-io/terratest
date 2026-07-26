package ssh_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoadSSHKeyPair(t *testing.T) {
	t.Parallel()

	expectedData, err := ssh.GenerateRSAKeyPairE(t, 2048) //nolint:mnd // RSA key size for testing
	require.NoError(t, err)

	tmpFolder := t.TempDir()
	ssh.SaveSSHKeyPair(t, tmpFolder, expectedData)

	actualData := ssh.LoadSSHKeyPair(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}
