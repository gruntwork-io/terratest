package ssh

import (
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/core/v2/teststate"
)

// sshKeyPairFilename is the name of the file, within a test folder's test data directory, used to store a KeyPair.
const sshKeyPairFilename = "SshKeyPair.json"

// SaveSSHKeyPair serializes and saves an SSH key pair into the given folder. This allows you to create an SSH key pair
// during setup and to reuse that key pair later during validation and teardown.
//
// The key pair is saved with teststate.SaveRedacted so that KeyPair.PrivateKey is not written to the test log.
func SaveSSHKeyPair(t testing.TestingT, testFolder string, keyPair *KeyPair) {
	teststate.SaveRedacted(t, formatSSHKeyPairPath(testFolder), true, keyPair)
}

// LoadSSHKeyPair loads and unserializes an SSH key pair from the given folder. This allows you to reuse an SSH key pair
// that was created during an earlier setup step in later validation and teardown steps.
func LoadSSHKeyPair(t testing.TestingT, testFolder string) *KeyPair {
	var keyPair KeyPair
	teststate.Load(t, formatSSHKeyPairPath(testFolder), &keyPair)

	return &keyPair
}

// formatSSHKeyPairPath formats a path to save an SSH key pair in the given folder.
func formatSSHKeyPairPath(testFolder string) string {
	return teststate.FormatPath(testFolder, sshKeyPairFilename)
}
