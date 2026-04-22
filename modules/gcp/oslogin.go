package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/oslogin/v1"
)

// ImportSSHKey will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// This will fail the test if there is an error.
//
// Deprecated: Use [ImportSSHKeyContext] instead.
func ImportSSHKey(t testing.TestingT, user, key string) {
	ImportSSHKeyContext(t, context.Background(), user, key)
}

// ImportSSHKeyContext will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ImportSSHKeyContext(t testing.TestingT, ctx context.Context, user, key string) {
	require.NoErrorf(t, ImportSSHKeyContextE(t, ctx, user, key), "Could not add SSH Key to user %s", user)
}

// ImportSSHKeyE will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
//
// Deprecated: Use [ImportSSHKeyContextE] instead.
func ImportSSHKeyE(t testing.TestingT, user, key string) error {
	return ImportSSHKeyContextE(t, context.Background(), user, key)
}

// ImportSSHKeyContextE will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// The ctx parameter supports cancellation and timeouts.
func ImportSSHKeyContextE(t testing.TestingT, ctx context.Context, user, key string) error {
	logger.Default.Logf(t, "Importing SSH key for user %s", user)

	service, err := NewOSLoginServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return ImportSSHKeyWithClient(ctx, service, user, key)
}

// ImportProjectSSHKey will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// The projectID parameter should be the chosen project ID.
// This will fail the test if there is an error.
//
// Deprecated: Use [ImportProjectSSHKeyContext] instead.
func ImportProjectSSHKey(t testing.TestingT, user, key, projectID string) {
	ImportProjectSSHKeyContext(t, context.Background(), user, key, projectID)
}

// ImportProjectSSHKeyContext will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// The projectID parameter should be the chosen project ID.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ImportProjectSSHKeyContext(t testing.TestingT, ctx context.Context, user, key, projectID string) {
	require.NoErrorf(t, ImportProjectSSHKeyContextE(t, ctx, user, key, projectID), "Could not add SSH Key to user %s", user)
}

// ImportProjectSSHKeyE will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// The projectID parameter should be the chosen project ID.
//
// Deprecated: Use [ImportProjectSSHKeyContextE] instead.
func ImportProjectSSHKeyE(t testing.TestingT, user, key, projectID string) error {
	return ImportProjectSSHKeyContextE(t, context.Background(), user, key, projectID)
}

// ImportProjectSSHKeyContextE will import an SSH key to GCP under the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key being uploaded.
// The projectID parameter should be the chosen project ID.
// The ctx parameter supports cancellation and timeouts.
func ImportProjectSSHKeyContextE(t testing.TestingT, ctx context.Context, user, key, projectID string) error {
	logger.Default.Logf(t, "Importing SSH key for user %s (project %s)", user, projectID)

	service, err := NewOSLoginServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return ImportProjectSSHKeyWithClient(ctx, service, user, key, projectID)
}

// ImportSSHKeyWithClient imports an SSH key to GCP under the provided user identity using the supplied
// *oslogin.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see oslogin_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func ImportSSHKeyWithClient(ctx context.Context, service *oslogin.Service, user, key string) error {
	_, err := service.Users.ImportSshPublicKey("users/"+user, &oslogin.SshPublicKey{Key: key}).Context(ctx).Do()

	return err
}

// ImportProjectSSHKeyWithClient imports an SSH key to GCP under the provided user identity, scoped to
// projectID, using the supplied *oslogin.Service.
// The ctx parameter supports cancellation and timeouts.
func ImportProjectSSHKeyWithClient(ctx context.Context, service *oslogin.Service, user, key, projectID string) error {
	_, err := service.Users.ImportSshPublicKey("users/"+user, &oslogin.SshPublicKey{Key: key}).ProjectId(projectID).Context(ctx).Do()

	return err
}

// DeleteSSHKey will delete an SSH key attached to the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key that was uploaded.
// This will fail the test if there is an error.
//
// Deprecated: Use [DeleteSSHKeyContext] instead.
func DeleteSSHKey(t testing.TestingT, user, key string) {
	DeleteSSHKeyContext(t, context.Background(), user, key)
}

// DeleteSSHKeyContext will delete an SSH key attached to the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key that was uploaded.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteSSHKeyContext(t testing.TestingT, ctx context.Context, user, key string) {
	require.NoErrorf(t, DeleteSSHKeyContextE(t, ctx, user, key), "Could not delete SSH Key for user %s", user)
}

// DeleteSSHKeyE will delete an SSH key attached to the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key that was uploaded.
//
// Deprecated: Use [DeleteSSHKeyContextE] instead.
func DeleteSSHKeyE(t testing.TestingT, user, key string) error {
	return DeleteSSHKeyContextE(t, context.Background(), user, key)
}

// DeleteSSHKeyContextE will delete an SSH key attached to the provided user identity.
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key that was uploaded.
// The ctx parameter supports cancellation and timeouts.
func DeleteSSHKeyContextE(t testing.TestingT, ctx context.Context, user, key string) error {
	logger.Default.Logf(t, "Deleting SSH key for user %s", user)

	service, err := NewOSLoginServiceContextE(t, ctx)
	if err != nil {
		return err
	}

	return DeleteSSHKeyWithClient(ctx, service, user, key)
}

// DeleteSSHKeyWithClient deletes an SSH key attached to the provided user identity using the supplied
// *oslogin.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see oslogin_unit_test.go for the pattern).
// The user parameter should be the email address of the user.
// The key parameter should be the public key of the SSH key that was uploaded.
// The ctx parameter supports cancellation and timeouts.
func DeleteSSHKeyWithClient(ctx context.Context, service *oslogin.Service, user, key string) error {
	loginProfile, err := GetLoginProfileWithClient(ctx, service, user)
	if err != nil {
		return err
	}

	for _, v := range loginProfile.SshPublicKeys {
		if key == v.Key {
			path := fmt.Sprintf("users/%s/sshPublicKeys/%s", user, v.Fingerprint)

			if _, err := service.Users.SshPublicKeys.Delete(path).Context(ctx).Do(); err != nil {
				return err
			}

			break
		}
	}

	return nil
}

// GetLoginProfile will retrieve the login profile for a user's Google identity. The login profile is a combination of
// OS Login + gcloud SSH keys and POSIX accounts the user will appear as. Generally, this will only be the OS Login
// key + account, but gcloud compute ssh could create temporary keys and profiles.
// The user parameter should be the email address of the user.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetLoginProfileContext] instead.
func GetLoginProfile(t testing.TestingT, user string) *oslogin.LoginProfile {
	return GetLoginProfileContext(t, context.Background(), user)
}

// GetLoginProfileContext will retrieve the login profile for a user's Google identity. The login profile is a
// combination of OS Login + gcloud SSH keys and POSIX accounts the user will appear as. Generally, this will only be
// the OS Login key + account, but gcloud compute ssh could create temporary keys and profiles.
// The user parameter should be the email address of the user.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLoginProfileContext(t testing.TestingT, ctx context.Context, user string) *oslogin.LoginProfile {
	profile, err := GetLoginProfileContextE(t, ctx, user)
	require.NoErrorf(t, err, "Could not get login profile for user %s", user)

	return profile
}

// GetLoginProfileE will retrieve the login profile for a user's Google identity. The login profile is a combination of
// OS Login + gcloud SSH keys and POSIX accounts the user will appear as. Generally, this will only be the OS Login
// key + account, but gcloud compute ssh could create temporary keys and profiles.
// The user parameter should be the email address of the user.
//
// Deprecated: Use [GetLoginProfileContextE] instead.
func GetLoginProfileE(t testing.TestingT, user string) (*oslogin.LoginProfile, error) {
	return GetLoginProfileContextE(t, context.Background(), user)
}

// GetLoginProfileContextE will retrieve the login profile for a user's Google identity. The login profile is a
// combination of OS Login + gcloud SSH keys and POSIX accounts the user will appear as. Generally, this will only be
// the OS Login key + account, but gcloud compute ssh could create temporary keys and profiles.
// The user parameter should be the email address of the user.
// The ctx parameter supports cancellation and timeouts.
func GetLoginProfileContextE(t testing.TestingT, ctx context.Context, user string) (*oslogin.LoginProfile, error) {
	logger.Default.Logf(t, "Getting login profile for user %s", user)

	service, err := NewOSLoginServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetLoginProfileWithClient(ctx, service, user)
}

// GetLoginProfileWithClient retrieves the login profile for a user's Google identity using the supplied
// *oslogin.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see oslogin_unit_test.go for the pattern).
// The user parameter should be the email address of the user.
// The ctx parameter supports cancellation and timeouts.
func GetLoginProfileWithClient(ctx context.Context, service *oslogin.Service, user string) (*oslogin.LoginProfile, error) {
	name := "users/" + user

	profile, err := service.Users.GetLoginProfile(name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// NewOSLoginServiceE creates a new OS Login service, which is used to make OS Login API calls.
//
// Deprecated: Use [NewOSLoginServiceContextE] instead.
func NewOSLoginServiceE(t testing.TestingT) (*oslogin.Service, error) {
	return NewOSLoginServiceContextE(t, context.Background())
}

// NewOSLoginServiceContextE creates a new OS Login service, which is used to make OS Login API calls.
// The ctx parameter supports cancellation and timeouts.
func NewOSLoginServiceContextE(t testing.TestingT, ctx context.Context) (*oslogin.Service, error) {
	if ts, ok := getStaticTokenSource(); ok {
		return oslogin.NewService(ctx, option.WithTokenSource(ts))
	}

	client, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("failed to get default client: %w", err)
	}

	service, err := oslogin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return service, nil
}
