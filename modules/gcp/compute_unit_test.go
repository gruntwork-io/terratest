package gcp_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// newFakeComputeService points a *compute.Service at a local httptest server. Gives unit tests a
// credential-free way to exercise *WithClient variants, analogous to the Azure azfake pattern.
func newFakeComputeService(t *testing.T, handler http.Handler) *compute.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	svc, err := compute.NewService(context.Background(),
		option.WithEndpoint(server.URL+"/"), option.WithoutAuthentication())
	require.NoError(t, err)

	return svc
}

// respond returns a handler that asserts the HTTP method (when non-empty) and path substring,
// then writes body with the given status. Cuts per-test boilerplate.
func respond(t *testing.T, method, pathContains string, status int, body string) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		if method != "" {
			assert.Equal(t, method, r.Method, "unexpected HTTP method")
		}

		assert.Contains(t, r.URL.Path, pathContains, "unexpected API path")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

// fetchInstanceForTest resolves a gcp.Instance via the aggregated-list endpoint so that its
// unexported projectID is populated — required for methods on *Instance.
func fetchInstanceForTest(t *testing.T, projectID, name, zoneURL string) *gcp.Instance {
	t.Helper()

	body := fmt.Sprintf(`{"items":{"zones/us-central1-a":{"instances":[{"name":%q,"zone":%q}]}}}`, name, zoneURL)
	svc := newFakeComputeService(t, respond(t, "", "aggregated/instances", http.StatusOK, body))

	inst, err := gcp.FetchInstanceWithClient(context.Background(), svc, projectID, name)
	require.NoError(t, err)

	return inst
}

// TestNewMetadataPreservesExisting — regression test for issue #1655.
func TestNewMetadataPreservesExisting(t *testing.T) {
	t.Parallel()

	v := "old"
	old := &compute.Metadata{Fingerprint: "fp", Items: []*compute.MetadataItems{{Key: "k1", Value: &v}}}

	result := gcp.NewMetadata(old, map[string]string{"k2": "new"})

	got := make(map[string]string)
	for _, item := range result.Items {
		got[item.Key] = *item.Value
	}

	assert.Equal(t, "fp", result.Fingerprint)
	assert.Equal(t, "old", got["k1"])
	assert.Equal(t, "new", got["k2"])
}

func TestGetPublicIPContextE(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		wantIP  string
		nics    []*compute.NetworkInterface
		wantErr bool
	}{
		"returns external IP": {
			wantIP: "1.2.3.4",
			nics:   []*compute.NetworkInterface{{AccessConfigs: []*compute.AccessConfig{{NatIP: "1.2.3.4"}}}},
		},
		"no network interfaces": {wantErr: true},
		"no access configs":     {wantErr: true, nics: []*compute.NetworkInterface{{}}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			inst := &gcp.Instance{Instance: &compute.Instance{Name: "x", NetworkInterfaces: tc.nics}}

			ip, err := inst.GetPublicIPContextE(t, context.Background())
			if tc.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantIP, ip)
		})
	}
}

func TestFetchInstanceWithClient(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		body := `{"items":{"zones/us-central1-a":{"instances":[{"name":"x"}]}}}`
		svc := newFakeComputeService(t, respond(t, http.MethodGet, "aggregated/instances", http.StatusOK, body))

		inst, err := gcp.FetchInstanceWithClient(context.Background(), svc, "p", "x")
		require.NoError(t, err)
		assert.Equal(t, "x", inst.Name)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		svc := newFakeComputeService(t, respond(t, "", "aggregated/instances", http.StatusOK, `{"items":{}}`))

		_, err := gcp.FetchInstanceWithClient(context.Background(), svc, "p", "x")
		require.ErrorContains(t, err, "could not be found")
	})
}

func TestFetchImageWithClient(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		svc := newFakeComputeService(t, respond(t, http.MethodGet, "/global/images/my-image", http.StatusOK, `{"name":"my-image"}`))

		img, err := gcp.FetchImageWithClient(context.Background(), svc, "p", "my-image")
		require.NoError(t, err)
		assert.Equal(t, "my-image", img.Name)
	})

	t.Run("404 error propagates", func(t *testing.T) {
		t.Parallel()

		svc := newFakeComputeService(t, respond(t, "", "/global/images/", http.StatusNotFound, `{"error":{"code":404,"message":"image not found"}}`))

		_, err := gcp.FetchImageWithClient(context.Background(), svc, "p", "missing")
		require.ErrorContains(t, err, "image not found")
	})
}

func TestFetchZonalInstanceGroupWithClient(t *testing.T) {
	t.Parallel()

	body := `{"name":"zig","zone":"https://www.googleapis.com/compute/v1/projects/p/zones/us-central1-a"}`
	svc := newFakeComputeService(t, respond(t, http.MethodGet, "/zones/us-central1-a/instanceGroups/zig", http.StatusOK, body))

	ig, err := gcp.FetchZonalInstanceGroupWithClient(context.Background(), svc, "p", "us-central1-a", "zig")
	require.NoError(t, err)
	assert.Equal(t, "zig", ig.Name)
}

func TestFetchRegionalInstanceGroupWithClient(t *testing.T) {
	t.Parallel()

	body := `{"name":"rig","region":"https://www.googleapis.com/compute/v1/projects/p/regions/us-central1"}`
	svc := newFakeComputeService(t, respond(t, http.MethodGet, "/regions/us-central1/instanceGroups/rig", http.StatusOK, body))

	ig, err := gcp.FetchRegionalInstanceGroupWithClient(context.Background(), svc, "p", "us-central1", "rig")
	require.NoError(t, err)
	assert.Equal(t, "rig", ig.Name)
}

func TestSetLabelsWithClient(t *testing.T) {
	t.Parallel()

	zoneURL := "https://www.googleapis.com/compute/v1/projects/p/zones/us-central1-a"
	inst := fetchInstanceForTest(t, "p", "i", zoneURL)

	svc := newFakeComputeService(t, respond(t, http.MethodPost, "/instances/i/setLabels", http.StatusOK, `{"name":"op","status":"DONE"}`))

	require.NoError(t, inst.SetLabelsWithClient(context.Background(), svc, map[string]string{"env": "unit"}))
}

func TestSetMetadataWithClient(t *testing.T) {
	t.Parallel()

	zoneURL := "https://www.googleapis.com/compute/v1/projects/p/zones/us-central1-a"
	inst := fetchInstanceForTest(t, "p", "i", zoneURL)

	svc := newFakeComputeService(t, respond(t, http.MethodPost, "/instances/i/setMetadata", http.StatusOK, `{"name":"op","status":"DONE"}`))

	require.NoError(t, inst.SetMetadataWithClient(context.Background(), svc, map[string]string{"k": "v"}))
}

func TestAddSSHKeyWithClient(t *testing.T) {
	t.Parallel()

	zoneURL := "https://www.googleapis.com/compute/v1/projects/p/zones/us-central1-a"
	inst := fetchInstanceForTest(t, "p", "i", zoneURL)

	t.Run("happy", func(t *testing.T) {
		t.Parallel()

		svc := newFakeComputeService(t, respond(t, http.MethodPost, "/instances/i/setMetadata", http.StatusOK, `{"name":"op","status":"DONE"}`))

		require.NoError(t, inst.AddSSHKeyWithClient(context.Background(), svc, "alice", "ssh-rsa A alice@h"))
	})

	t.Run("SDK error is wrapped", func(t *testing.T) {
		t.Parallel()

		svc := newFakeComputeService(t, respond(t, "", "/instances/i/setMetadata", http.StatusBadRequest, `{"error":{"code":400,"message":"bad"}}`))

		err := inst.AddSSHKeyWithClient(context.Background(), svc, "alice", "ssh-rsa A alice@h")
		require.ErrorContains(t, err, "failed to add SSH key")
	})
}

func TestDeleteImageWithClient(t *testing.T) {
	t.Parallel()

	imgSvc := newFakeComputeService(t, respond(t, http.MethodGet, "/global/images/img", http.StatusOK, `{"name":"img"}`))
	img, err := gcp.FetchImageWithClient(context.Background(), imgSvc, "p", "img")
	require.NoError(t, err)

	svc := newFakeComputeService(t, respond(t, http.MethodDelete, "/global/images/img", http.StatusOK, `{"name":"op"}`))

	require.NoError(t, img.DeleteImageWithClient(context.Background(), svc))
}

func TestZonalInstanceGroupGetInstanceIDsWithClient(t *testing.T) {
	t.Parallel()

	igBody := `{"name":"zig","zone":"https://www.googleapis.com/compute/v1/projects/p/zones/us-central1-a"}`
	svc := newFakeComputeService(t, respond(t, "", "zig", http.StatusOK, igBody))
	ig, err := gcp.FetchZonalInstanceGroupWithClient(context.Background(), svc, "p", "us-central1-a", "zig")
	require.NoError(t, err)

	listBody := `{"items":[{"instance":"https://.../instances/a"},{"instance":"https://.../instances/b"}]}`
	svc2 := newFakeComputeService(t, respond(t, http.MethodPost, "/instanceGroups/zig/listInstances", http.StatusOK, listBody))

	ids, err := ig.GetInstanceIDsWithClient(context.Background(), svc2)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, ids)
}

func TestRegionalInstanceGroupGetInstanceIDsWithClient(t *testing.T) {
	t.Parallel()

	igBody := `{"name":"rig","region":"https://www.googleapis.com/compute/v1/projects/p/regions/us-central1"}`
	svc := newFakeComputeService(t, respond(t, "", "rig", http.StatusOK, igBody))
	ig, err := gcp.FetchRegionalInstanceGroupWithClient(context.Background(), svc, "p", "us-central1", "rig")
	require.NoError(t, err)

	listBody := `{"items":[{"instance":"https://.../instances/c"},{"instance":"https://.../instances/d"}]}`
	svc2 := newFakeComputeService(t, respond(t, http.MethodPost, "/instanceGroups/rig/listInstances", http.StatusOK, listBody))

	ids, err := ig.GetInstanceIDsWithClient(context.Background(), svc2)
	require.NoError(t, err)
	assert.Equal(t, []string{"c", "d"}, ids)
}
