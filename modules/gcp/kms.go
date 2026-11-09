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

// GetKMSEkmConnectionAttrs returns the settings Google Cloud holds for the given Cloud KMS EKM connection, so a test can assert on
// what was actually created rather than only that it exists.
// The service resolvers a connection reaches the external key manager through are reported in full,
// hostname and server certificates included.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKMSEkmConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *cloudkms.EkmConnection {
	result, err := GetKMSEkmConnectionAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return result
}

// GetKMSEkmConnectionAttrsE returns the settings Google Cloud holds for the given Cloud KMS EKM connection.
// The ctx parameter supports cancellation and timeouts.
func GetKMSEkmConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*cloudkms.EkmConnection, error) {
	logger.Default.Logf(t, "Getting settings for Cloud KMS EKM connection %s in %s of project %s", connectionID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetKMSEkmConnectionAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetKMSEkmConnectionAttrsWithClient returns the settings Google Cloud holds for the given Cloud KMS EKM connection using the supplied
// *cloudkms.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetKMSEkmConnectionAttrsWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, connectionID string) (*cloudkms.EkmConnection, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/ekmConnections/%s", projectID, location, connectionID)

	result, err := service.Projects.Locations.EkmConnections.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud KMS EKM connection %s does not exist in %s of project %s", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud KMS EKM connection %s in %s of project %s: %w", connectionID, location, projectID, err)
	}

	return result, nil
}

// GetKMSKeyHandleAttrs returns the settings Google Cloud holds for the given Cloud KMS key handle, so a test can assert on
// what was actually created rather than only that it exists.
// The key Autokey provisioned for the handle is named only on the handle, so a read is the only way to
// tell which key a resource ended up protected by.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKMSKeyHandleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, handleID string) *cloudkms.KeyHandle {
	result, err := GetKMSKeyHandleAttrsE(t, ctx, projectID, location, handleID)
	require.NoError(t, err)

	return result
}

// GetKMSKeyHandleAttrsE returns the settings Google Cloud holds for the given Cloud KMS key handle.
// The ctx parameter supports cancellation and timeouts.
func GetKMSKeyHandleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, handleID string) (*cloudkms.KeyHandle, error) {
	logger.Default.Logf(t, "Getting settings for Cloud KMS key handle %s in %s of project %s", handleID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetKMSKeyHandleAttrsWithClient(ctx, service, projectID, location, handleID)
}

// GetKMSKeyHandleAttrsWithClient returns the settings Google Cloud holds for the given Cloud KMS key handle using the supplied
// *cloudkms.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetKMSKeyHandleAttrsWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, handleID string) (*cloudkms.KeyHandle, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/keyHandles/%s", projectID, location, handleID)

	result, err := service.Projects.Locations.KeyHandles.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud KMS key handle %s does not exist in %s of project %s", handleID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud KMS key handle %s in %s of project %s: %w", handleID, location, projectID, err)
	}

	return result, nil
}

// GetKMSProjectAutokeyConfigAttrs returns the settings Google Cloud holds for the Cloud KMS Autokey config of the given project, so a test can assert on
// what was actually created rather than only that it exists.
// The config is a singleton per project and names the folder whose key project Autokey will draw
// from, so a read is how a test tells whether the module pointed it at the right one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKMSProjectAutokeyConfigAttrs(t testing.TestingT, ctx context.Context, projectID string) *cloudkms.AutokeyConfig {
	result, err := GetKMSProjectAutokeyConfigAttrsE(t, ctx, projectID)
	require.NoError(t, err)

	return result
}

// GetKMSProjectAutokeyConfigAttrsE returns the settings Google Cloud holds for the Cloud KMS Autokey config of the given project.
// The ctx parameter supports cancellation and timeouts.
func GetKMSProjectAutokeyConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string) (*cloudkms.AutokeyConfig, error) {
	logger.Default.Logf(t, "Getting the Cloud KMS Autokey config of project %s", projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetKMSProjectAutokeyConfigAttrsWithClient(ctx, service, projectID)
}

// GetKMSProjectAutokeyConfigAttrsWithClient returns the settings Google Cloud holds for the Cloud KMS Autokey config of the given project using the supplied
// *cloudkms.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetKMSProjectAutokeyConfigAttrsWithClient(ctx context.Context, service *cloudkms.Service, projectID string) (*cloudkms.AutokeyConfig, error) {
	name := fmt.Sprintf("projects/%s/autokeyConfig", projectID)

	result, err := service.Projects.GetAutokeyConfig(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("project %s has no Cloud KMS Autokey config", projectID)
		}

		return nil, fmt.Errorf("failed to get the Cloud KMS Autokey config of project %s: %w", projectID, err)
	}

	return result, nil
}

// GetKMSEkmConnectionIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Cloud
// KMS EKM connection, so a test can assert on the bindings that were applied rather than only that
// the apply finished.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKMSEkmConnectionIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *cloudkms.Policy {
	policy, err := GetKMSEkmConnectionIamPolicyAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return policy
}

// GetKMSEkmConnectionIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Cloud
// KMS EKM connection.
// The ctx parameter supports cancellation and timeouts.
func GetKMSEkmConnectionIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*cloudkms.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Cloud KMS EKM connection %s in %s of project %s", connectionID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetKMSEkmConnectionIamPolicyAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetKMSEkmConnectionIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the
// given Cloud KMS EKM connection using the supplied *cloudkms.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetKMSEkmConnectionIamPolicyAttrsWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, connectionID string) (*cloudkms.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/ekmConnections/%s", projectID, location, connectionID)

	policy, err := service.Projects.Locations.EkmConnections.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud KMS EKM connection %s does not exist in %s of project %s", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Cloud KMS EKM connection %s in %s of project %s: %w", connectionID, location, projectID, err)
	}

	return policy, nil
}

// NewCloudKMSServiceE creates a Cloud KMS service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudKMSServiceE(t testing.TestingT, ctx context.Context) (*cloudkms.Service, error) {
	return cloudkms.NewService(ctx, append(withOptions(), option.WithScopes(cloudkms.CloudPlatformScope))...)
}

// DecryptSecretCiphertext returns the plaintext Cloud KMS gets back from the given ciphertext, so a
// test can assert that what a module encrypted is what it was given. The ciphertext is base64 as the
// provider returns it, and the plaintext comes back decoded. A ciphertext encrypted with additional
// authenticated data can only be decrypted with the same data, so that is passed too, as the plain
// string the module was given; pass an empty string when the module named none.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DecryptSecretCiphertext(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string, cryptoKeyID string, ciphertext string, additionalAuthenticatedData string) string {
	plaintext, err := DecryptSecretCiphertextE(t, ctx, projectID, location, keyRingID, cryptoKeyID, ciphertext, additionalAuthenticatedData)
	require.NoError(t, err)

	return plaintext
}

// DecryptSecretCiphertextE returns the plaintext Cloud KMS gets back from the given ciphertext.
// The ctx parameter supports cancellation and timeouts.
func DecryptSecretCiphertextE(t testing.TestingT, ctx context.Context, projectID string, location string, keyRingID string, cryptoKeyID string, ciphertext string, additionalAuthenticatedData string) (string, error) {
	logger.Default.Logf(t, "Decrypting a ciphertext with key %s in key ring %s in %s in project %s", cryptoKeyID, keyRingID, location, projectID)

	service, err := NewCloudKMSServiceE(t, ctx)
	if err != nil {
		return "", err
	}

	return DecryptSecretCiphertextWithClient(ctx, service, projectID, location, keyRingID, cryptoKeyID, ciphertext, additionalAuthenticatedData)
}

// DecryptSecretCiphertextWithClient returns the plaintext Cloud KMS gets back from the given
// ciphertext using the supplied *cloudkms.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see kms_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DecryptSecretCiphertextWithClient(ctx context.Context, service *cloudkms.Service, projectID string, location string, keyRingID string, cryptoKeyID string, ciphertext string, additionalAuthenticatedData string) (string, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectID, location, keyRingID, cryptoKeyID)

	// The call takes a request body rather than a plain resource name, and returns the plaintext
	// base64 encoded whatever the plaintext was.
	request := &cloudkms.DecryptRequest{Ciphertext: ciphertext}
	if additionalAuthenticatedData != "" {
		request.AdditionalAuthenticatedData = base64.StdEncoding.EncodeToString([]byte(additionalAuthenticatedData))
	}

	response, err := service.Projects.Locations.KeyRings.CryptoKeys.Decrypt(name, request).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to decrypt a ciphertext with key %s in key ring %s in %s in project %s: %w", cryptoKeyID, keyRingID, location, projectID, err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(response.Plaintext)
	if err != nil {
		return "", fmt.Errorf("the plaintext Cloud KMS returned is not base64: %w", err)
	}

	return string(plaintext), nil
}
