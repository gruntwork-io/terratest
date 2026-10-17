package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/api/oslogin/v1"
)

// newFakeOsLoginService points a *oslogin.Service at a local httptest server.
func newFakeOsLoginService(t *testing.T, handler http.Handler) *oslogin.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	svc, err := oslogin.NewService(context.Background(),
		option.WithEndpoint(server.URL+"/"), option.WithoutAuthentication())
	require.NoError(t, err)

	return svc
}

func TestImportSSHKeyWithClient(t *testing.T) {
	t.Parallel()

	const user = "u@example.com"

	svc := newFakeOsLoginService(t, respond(t, http.MethodPost, "/users/"+user+":importSshPublicKey", http.StatusOK, `{"loginProfile":{"name":"users/`+user+`"}}`))

	require.NoError(t, gcp.ImportSSHKeyWithClient(context.Background(), svc, user, "ssh-rsa A"))
}

func TestImportProjectSSHKeyWithClient(t *testing.T) {
	t.Parallel()

	const (
		user      = "u@example.com"
		projectID = "my-project"
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/users/"+user+":importSshPublicKey")
		assert.Equal(t, projectID, r.URL.Query().Get("projectId"), "expected projectId query param")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"loginProfile":{"name":"users/` + user + `"}}`))
	})

	svc := newFakeOsLoginService(t, handler)

	require.NoError(t, gcp.ImportProjectSSHKeyWithClient(context.Background(), svc, user, "ssh-rsa A", projectID))
}

func TestGetLoginProfileWithClient(t *testing.T) {
	t.Parallel()

	const user = "u@example.com"

	body := `{"name":"users/` + user + `","sshPublicKeys":{"abc":{"key":"ssh-rsa A","fingerprint":"abc"}}}`
	svc := newFakeOsLoginService(t, respond(t, http.MethodGet, "/users/"+user+"/loginProfile", http.StatusOK, body))

	profile, err := gcp.GetLoginProfileWithClient(context.Background(), svc, user)
	require.NoError(t, err)
	assert.Equal(t, "users/"+user, profile.Name)
	assert.Len(t, profile.SshPublicKeys, 1)
}

func TestDeleteSSHKeyWithClient(t *testing.T) {
	t.Parallel()

	const (
		user        = "u@example.com"
		key         = "ssh-rsa A matching"
		fingerprint = "abc123"
	)

	t.Run("deletes matching key", func(t *testing.T) {
		t.Parallel()

		var deleteCalled bool

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			switch r.Method {
			case http.MethodGet:
				_, _ = w.Write([]byte(`{"name":"users/` + user + `","sshPublicKeys":{"` + fingerprint + `":{"key":"` + key + `","fingerprint":"` + fingerprint + `"}}}`))
			case http.MethodDelete:
				assert.Contains(t, r.URL.Path, "/sshPublicKeys/"+fingerprint)

				deleteCalled = true

				_, _ = w.Write([]byte(`{}`))
			}
		})

		svc := newFakeOsLoginService(t, handler)

		require.NoError(t, gcp.DeleteSSHKeyWithClient(context.Background(), svc, user, key))
		assert.True(t, deleteCalled, "Delete should fire when the key is present")
	})

	t.Run("no-op when no matching key", func(t *testing.T) {
		t.Parallel()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				t.Fatalf("Delete must not fire when no key matches")
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"name":"users/` + user + `","sshPublicKeys":{"` + fingerprint + `":{"key":"ssh-rsa A other","fingerprint":"` + fingerprint + `"}}}`))
		})

		svc := newFakeOsLoginService(t, handler)

		require.NoError(t, gcp.DeleteSSHKeyWithClient(context.Background(), svc, user, key))
	})
}

func TestGetOSLoginSSHPublicKeyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a key the OS Login SSH public key module
	// imported, not a copy of any one fixture's values. The expiry comes back as a quoted number.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/users/terratest@example.com/sshPublicKeys/9c3d1e0a"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"users/terratest@example.com/sshPublicKeys/9c3d1e0a","key":"ssh-rsa AAAAB3NzaC1yc2E terratest","fingerprint":"9c3d1e0a","expirationTimeUsec":"1790000000000000"}`))
	})

	key, err := gcp.GetOSLoginSSHPublicKeyAttrsWithClient(context.Background(), newFakeOsLoginService(t, handler), "terratest@example.com", "9c3d1e0a")
	require.NoError(t, err)

	assert.Equal(t, "9c3d1e0a", key.Fingerprint)
	assert.Contains(t, key.Key, "ssh-rsa")
	assert.Equal(t, int64(1790000000000000), key.ExpirationTimeUsec)
}

func TestGetOSLoginSSHPublicKeyAttrsWithClientReportsAMissingKey(t *testing.T) {
	t.Parallel()

	// A caller who asks for a key that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetOSLoginSSHPublicKeyAttrsWithClient(context.Background(), newFakeOsLoginService(t, handler), "terratest@example.com", "deadbeef")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has no OS Login SSH public key")
}
