package aws

import (
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/gruntwork-io/terratest/modules/core/v2/teststate"
)

// ec2KeyPairFilename is the name of the file, within a test folder's test data directory, used to store an Ec2Keypair.
const ec2KeyPairFilename = "Ec2KeyPair.json"

// SaveEc2KeyPair serializes and saves an Ec2Keypair into the given folder. This allows you to create an Ec2Keypair
// during setup and to reuse that Ec2Keypair later during validation and teardown.
//
// The key pair contains a private key, so unlike most test data, it is never written to the test log.
func SaveEc2KeyPair(t testing.TestingT, testFolder string, keyPair *Ec2Keypair) {
	teststate.SaveRedacted(t, formatEc2KeyPairPath(testFolder), true, keyPair)
}

// LoadEc2KeyPair loads and unserializes an Ec2Keypair from the given folder. This allows you to reuse an Ec2Keypair
// that was created during an earlier setup step in later validation and teardown steps.
func LoadEc2KeyPair(t testing.TestingT, testFolder string) *Ec2Keypair {
	var keyPair Ec2Keypair
	teststate.Load(t, formatEc2KeyPairPath(testFolder), &keyPair)

	return &keyPair
}

// formatEc2KeyPairPath formats a path to save an Ec2Keypair in the given folder.
func formatEc2KeyPairPath(testFolder string) string {
	return teststate.FormatPath(testFolder, ec2KeyPairFilename)
}
