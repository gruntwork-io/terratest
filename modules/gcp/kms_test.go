package gcp_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/option"
)

// newFakeCloudKMSService points a real Cloud KMS client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudKMSService(t *testing.T, handler http.Handler) *cloudkms.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudkms.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetKeyRingAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a key ring the terraform-google-
	// security module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test","createTime":"2026-09-22T10:00:00Z"}`))
	})

	keyRing, err := gcp.GetKeyRingAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test", keyRing.Name)
}

func TestGetKeyRingAttrsWithClientMissingKeyRing(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the key ring and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetKeyRingAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetCryptoKeyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a key the terraform-google-security
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/keyRings/gw-library-test/cryptoKeys/gw-key", "unexpected path")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-key","purpose":"ENCRYPT_DECRYPT","rotationPeriod":"7776000s","labels":{"purpose":"terratest"},"versionTemplate":{"algorithm":"GOOGLE_SYMMETRIC_ENCRYPTION","protectionLevel":"SOFTWARE"}}`))
	})

	key, err := gcp.GetCryptoKeyAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-key")
	require.NoError(t, err)

	assert.Equal(t, "ENCRYPT_DECRYPT", key.Purpose)
	assert.Equal(t, "7776000s", key.RotationPeriod)
	assert.Equal(t, "terratest", key.Labels["purpose"])
	require.NotNil(t, key.VersionTemplate)
	assert.Equal(t, "GOOGLE_SYMMETRIC_ENCRYPTION", key.VersionTemplate.Algorithm)
}

func TestGetCryptoKeyAttrsWithClientMissingKey(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the key and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetCryptoKeyAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-ring", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-ring")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestDecryptSecretCiphertextWithClient(t *testing.T) {
	t.Parallel()

	// Cloud KMS answers a decrypt with the plaintext base64 encoded, whatever the plaintext was.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/keyRings/gw-library-test/cryptoKeys/gw-key:decrypt"), "unexpected path %s", r.URL.Path)

		var request struct {
			Ciphertext                  string `json:"ciphertext"`
			AdditionalAuthenticatedData string `json:"additionalAuthenticatedData"`
		}

		// assert rather than require: a failed require inside a handler stops the wrong goroutine.
		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&request)) {
			assert.Equal(t, "Q2lwaGVyVGV4dA==", request.Ciphertext)
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("terratest")), request.AdditionalAuthenticatedData,
				"a ciphertext encrypted with additional authenticated data can only be decrypted with it")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"plaintext":"d3JpdHRlbiBieSB0ZXJyYXRlc3Q="}`))
	})

	plaintext, err := gcp.DecryptSecretCiphertextWithClient(context.Background(), newFakeCloudKMSService(t, handler),
		"gw-library-test-project", "us-central1", "gw-library-test", "gw-key", "Q2lwaGVyVGV4dA==", "terratest")
	require.NoError(t, err)

	assert.Equal(t, "written by terratest", plaintext)
}

func TestDecryptSecretCiphertextWithClientWithoutAuthenticatedData(t *testing.T) {
	t.Parallel()

	// A module that named no additional authenticated data must not have one sent on its behalf, or
	// Cloud KMS refuses the decrypt.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			AdditionalAuthenticatedData string `json:"additionalAuthenticatedData"`
		}

		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&request)) {
			assert.Empty(t, request.AdditionalAuthenticatedData)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"plaintext":"d3JpdHRlbiBieSB0ZXJyYXRlc3Q="}`))
	})

	plaintext, err := gcp.DecryptSecretCiphertextWithClient(context.Background(), newFakeCloudKMSService(t, handler),
		"gw-library-test-project", "us-central1", "gw-library-test", "gw-key", "Q2lwaGVyVGV4dA==", "")
	require.NoError(t, err)

	assert.Equal(t, "written by terratest", plaintext)
}
