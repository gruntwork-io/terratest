package gcp

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetKeyRingAttrs returns the settings Google Cloud holds for the given KMS key ring, so a test can
// assert on what was actually created rather than only that it exists. A key ring cannot be
// deleted, so one outlives every test that creates it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKeyRingAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string) *cloudkms.KeyRing {
	keyRing, err := GetKeyRingAttrsE(t, ctx, projectID, location, keyRingID)
	require.NoError(t, err)

	return keyRing
}

// GetKeyRingAttrsE returns the settings Google Cloud holds for the given KMS key ring.
// The ctx parameter supports cancellation and timeouts.
func GetKeyRingAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string) (*cloudkms.KeyRing, error) {
	logger.Default.Logf(t, "Getting settings for KMS key ring %s in location %s in project %s", keyRingID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetKeyRingAttrsWithClient(ctx, service, projectID, location, keyRingID)
}

// GetKeyRingAttrsWithClient returns the settings Google Cloud holds for the given KMS key ring
// using the supplied *cloudkms.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetKeyRingAttrsWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, keyRingID string) (*cloudkms.KeyRing, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", projectID, location, keyRingID)

	keyRing, err := service.Projects.Locations.KeyRings.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("KMS key ring %s does not exist in location %s in project %s", keyRingID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for KMS key ring %s in location %s in project %s: %w", keyRingID, location, projectID, err)
	}

	return keyRing, nil
}

// GetCryptoKeyAttrs returns the settings Google Cloud holds for the given KMS key, so a test can
// assert on what was actually created rather than only that it exists. A key is named by its ring
// as well as its own id, and destroying one only schedules its versions for deletion.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCryptoKeyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string, keyID string) *cloudkms.CryptoKey {
	cryptoKey, err := GetCryptoKeyAttrsE(t, ctx, projectID, location, keyRingID, keyID)
	require.NoError(t, err)

	return cryptoKey
}

// GetCryptoKeyAttrsE returns the settings Google Cloud holds for the given KMS key.
// The ctx parameter supports cancellation and timeouts.
func GetCryptoKeyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string, keyID string) (*cloudkms.CryptoKey, error) {
	logger.Default.Logf(t, "Getting settings for KMS key %s in key ring %s in location %s in project %s", keyID, keyRingID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCryptoKeyAttrsWithClient(ctx, service, projectID, location, keyRingID, keyID)
}

// GetCryptoKeyAttrsWithClient returns the settings Google Cloud holds for the given KMS key using
// the supplied *cloudkms.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCryptoKeyAttrsWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, keyRingID string, keyID string) (*cloudkms.CryptoKey, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectID, location, keyRingID, keyID)

	cryptoKey, err := service.Projects.Locations.KeyRings.CryptoKeys.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("KMS key %s does not exist in key ring %s in location %s in project %s", keyID, keyRingID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for KMS key %s in key ring %s in location %s in project %s: %w", keyID, keyRingID, location, projectID, err)
	}

	return cryptoKey, nil
}

// NewCloudKMSServiceE creates a Cloud KMS service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudKMSServiceE(t testing.TestingT, ctx context.Context) (*cloudkms.Service, error) {
	return cloudkms.NewService(ctx, append(withOptions(), option.WithScopes(cloudkms.CloudPlatformScope))...)
}

// DecryptSecretCiphertext returns the plaintext Cloud KMS gets back from the given ciphertext, so a
// test can assert that what a module encrypted is what it was given. The ciphertext is base64 as the
// provider returns it, and the plaintext comes back decoded.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DecryptSecretCiphertext(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string, cryptoKeyID string, ciphertext string) string {
	plaintext, err := DecryptSecretCiphertextE(t, ctx, projectID, location, keyRingID, cryptoKeyID, ciphertext)
	require.NoError(t, err)

	return plaintext
}

// DecryptSecretCiphertextE returns the plaintext Cloud KMS gets back from the given ciphertext.
// The ctx parameter supports cancellation and timeouts.
func DecryptSecretCiphertextE(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string, cryptoKeyID string, ciphertext string) (string, error) {
	logger.Default.Logf(t, "Decrypting a ciphertext with key %s in key ring %s in %s in project %s", cryptoKeyID, keyRingID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return "", err
	}

	return DecryptSecretCiphertextWithClient(ctx, service, projectID, location, keyRingID, cryptoKeyID, ciphertext)
}

// DecryptSecretCiphertextWithClient returns the plaintext Cloud KMS gets back from the given
// ciphertext using the supplied *cloudkms.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DecryptSecretCiphertextWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, keyRingID string, cryptoKeyID string, ciphertext string) (string, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectID, location, keyRingID, cryptoKeyID)

	// The call takes a request body rather than a plain resource name, and returns the plaintext
	// base64 encoded whatever the plaintext was.
	response, err := service.Projects.Locations.KeyRings.CryptoKeys.Decrypt(name,
		&cloudkms.DecryptRequest{Ciphertext: ciphertext}).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to decrypt a ciphertext with key %s in key ring %s in %s in project %s: %w", cryptoKeyID, keyRingID, location, projectID, err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(response.Plaintext)
	if err != nil {
		return "", fmt.Errorf("Cloud KMS returned a plaintext that is not base64: %w", err)
	}

	return string(plaintext), nil
}
