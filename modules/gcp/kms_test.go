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

func TestGetKMSEkmConnectionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a connection the EKM connection module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/ekmConnections/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/ekmConnections/gw-library-test","keyManagementMode":"MANUAL","serviceResolvers":[{"serviceDirectoryService":"projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test/services/gw-library-test","hostname":"ekm.example.com","endpointFilter":"","serverCertificates":[{"rawDer":"MIIB"}]}]}`))
	})

	connection, err := gcp.GetKMSEkmConnectionAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "MANUAL", connection.KeyManagementMode)
	require.Len(t, connection.ServiceResolvers, 1)
	assert.Equal(t, "ekm.example.com", connection.ServiceResolvers[0].Hostname)
}

func TestGetKMSEkmConnectionAttrsWithClientReportsAMissingConnection(t *testing.T) {
	t.Parallel()

	// A caller who asks for a connection that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetKMSEkmConnectionAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetKMSEkmConnectionIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The policy is shaped like the one Google returns after the EKM connection IAM policy module has
	// applied a binding, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/ekmConnections/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("options.requestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwYabc","bindings":[{"role":"roles/cloudkms.viewer","members":["serviceAccount:terratest@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetKMSEkmConnectionIamPolicyAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/cloudkms.viewer", policy.Bindings[0].Role)
}

func TestGetKMSEkmConnectionIamPolicyAttrsWithClientReportsAMissingConnection(t *testing.T) {
	t.Parallel()

	// Asking for the policy of a connection that is not there should name the connection rather than
	// hand back the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetKMSEkmConnectionIamPolicyAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetKMSKeyHandleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a handle the key handle module created,
	// not a copy of any one fixture's values. The key is the one Autokey provisioned.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/keyHandles/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/keyHandles/gw-library-test","kmsKey":"projects/gw-library-keys/locations/us-central1/keyRings/autokey/cryptoKeys/storage-gw-library","resourceTypeSelector":"storage.googleapis.com/Bucket"}`))
	})

	handle, err := gcp.GetKMSKeyHandleAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "storage.googleapis.com/Bucket", handle.ResourceTypeSelector)
	assert.Contains(t, handle.KmsKey, "/cryptoKeys/storage-gw-library")
}

func TestGetKMSKeyHandleAttrsWithClientReportsAMissingHandle(t *testing.T) {
	t.Parallel()

	// A caller who asks for a handle that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetKMSKeyHandleAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetKMSProjectAutokeyConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns once the project Autokey config module has
	// pointed a project at a key project, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/autokeyConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/autokeyConfig","keyProject":"projects/gw-library-keys","state":"ACTIVE"}`))
	})

	config, err := gcp.GetKMSProjectAutokeyConfigAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-keys", config.KeyProject)
	assert.Equal(t, "ACTIVE", config.State)
}

func TestGetKMSProjectAutokeyConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A project Autokey has never been set up on answers 404, and a caller should be told that in
	// those words rather than be handed the transport's own.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetKMSProjectAutokeyConfigAttrsWithClient(context.Background(), newFakeCloudKMSService(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has no Cloud KMS Autokey config")
}
