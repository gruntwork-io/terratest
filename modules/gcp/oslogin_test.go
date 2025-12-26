//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
)

// TestOSLogin groups all OS Login tests that mutate SSH keys for the same user.
// These tests cannot run in parallel with each other because Google's OS Login API
// returns "409: Multiple concurrent mutations" errors when multiple operations
// modify the same user's SSH keys simultaneously.
//
// By grouping them in a single test function with subtests (without t.Parallel()),
// we ensure they run sequentially while still allowing other GCP tests to run in parallel.
func TestOSLogin(t *testing.T) {
	t.Parallel() // This test can run in parallel with OTHER GCP tests

	// Subtests run sequentially (no t.Parallel() on subtests) to avoid 409 conflicts
	t.Run("ImportSSHKey", func(t *testing.T) {
		keyPair := ssh.GenerateRSAKeyPair(t, 2048)
		key := keyPair.PublicKey

		user := GetGoogleIdentityEmailEnvVar(t)

		defer DeleteSSHKey(t, user, key)
		ImportSSHKey(t, user, key)
	})

	t.Run("ImportProjectSSHKey", func(t *testing.T) {
		keyPair := ssh.GenerateRSAKeyPair(t, 2048)
		key := keyPair.PublicKey

		user := GetGoogleIdentityEmailEnvVar(t)
		projectID := GetGoogleProjectIDFromEnvVar(t)

		defer DeleteSSHKey(t, user, key)
		ImportProjectSSHKey(t, user, key, projectID)
	})

	t.Run("GetLoginProfile", func(t *testing.T) {
		user := GetGoogleIdentityEmailEnvVar(t)
		GetLoginProfile(t, user)
	})

	t.Run("SetOSLoginKey", func(t *testing.T) {
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
	})
}
