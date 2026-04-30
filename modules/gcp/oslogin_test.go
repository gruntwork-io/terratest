package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
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
