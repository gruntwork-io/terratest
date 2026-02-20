//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
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

	// Clean up any stale SSH keys from previous test runs to avoid
	// "Login profile size exceeds 32 KiB" errors.
	user := GetGoogleIdentityEmailEnvVar(t)
	purgeAllSSHKeys(t, user)

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

// purgeAllSSHKeys deletes all SSH keys from the user's OS Login profile.
// This prevents "Login profile size exceeds 32 KiB" errors caused by
// stale keys accumulating from previous test runs.
func purgeAllSSHKeys(t *testing.T, user string) {
	profile, err := GetLoginProfileE(t, user)
	if err != nil {
		t.Logf("Warning: could not get login profile to purge keys: %v", err)
		return
	}

	if len(profile.SshPublicKeys) == 0 {
		return
	}

	logger.Default.Logf(t, "Purging %d stale SSH keys from OS Login profile for user %s", len(profile.SshPublicKeys), user)

	service, err := NewOSLoginServiceE(t)
	if err != nil {
		t.Logf("Warning: could not create OS Login service to purge keys: %v", err)
		return
	}

	ctx := context.Background()
	for fingerprint := range profile.SshPublicKeys {
		path := fmt.Sprintf("users/%s/sshPublicKeys/%s", user, fingerprint)
		if _, err := service.Users.SshPublicKeys.Delete(path).Context(ctx).Do(); err != nil {
			t.Logf("Warning: could not delete SSH key %s: %v", fingerprint, err)
		}
	}
}
