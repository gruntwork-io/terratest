package ssh_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh/v2"
	"github.com/stretchr/testify/assert"
)

// Basic test to ensure we can successfully generate key pairs (no explicit validation for now).
func TestGenerateRSAKeyPair(t *testing.T) {
	t.Parallel()

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	assert.Contains(t, keyPair.PublicKey, "ssh-rsa")
	assert.Contains(t, keyPair.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----")
}
