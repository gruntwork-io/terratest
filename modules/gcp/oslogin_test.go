// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
)

func TestImportSSHKeyOSLogin(t *testing.T) {
	t.Parallel()

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	key := keyPair.PublicKey

	user := GetGoogleIdentityEmailEnvVar(t)

	defer DeleteSSHKey(t, user, key)
	ImportSSHKey(t, user, key)
}

func TestGetLoginProfile(t *testing.T) {
	t.Parallel()

	user := GetGoogleIdentityEmailEnvVar(t)
	GetLoginProfile(t, user)
}

func TestSetOSLoginKey(t *testing.T) {
	t.Parallel()

	keyPair := ssh.GenerateRSAKeyPair(t, 2048)
	key := keyPair.PublicKey

	user := GetGoogleIdentityEmailEnvVar(t)

	defer DeleteSSHKey(t, user, key)
	ImportSSHKey(t, user, key)
	loginProfile := GetLoginProfile(t, user)

	found := false
	for _, v := range loginProfile.SshPublicKeys {
		if key == v.Key {
			found = true
		}
	}

	if found != true {
		t.Fatalf("Did not find key in login profile for user %s", user)
	}
}
