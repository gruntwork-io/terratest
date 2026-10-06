package gcp_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
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

func TestSetLabelsWithClientMergesExisting(t *testing.T) {
	t.Parallel()

	zoneURL := "https://www.googleapis.com/compute/v1/projects/p/zones/us-central1-a"
	inst := fetchInstanceForTest(t, "p", "i", zoneURL)
	inst.Labels = map[string]string{"team": "platform"}

	var sentLabels map[string]string

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Contains(t, r.URL.Path, "/instances/i/setLabels")

		var req compute.InstancesSetLabelsRequest
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		sentLabels = req.Labels

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"name":"op","status":"DONE"}`))
	}

	svc := newFakeComputeService(t, http.HandlerFunc(handler))

	require.NoError(t, inst.SetLabelsWithClient(context.Background(), svc, map[string]string{"env": "unit"}))

	assert.Equal(t, "platform", sentLabels["team"], "existing label should be preserved")
	assert.Equal(t, "unit", sentLabels["env"], "new label should be set")
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

func TestFetchNetworkWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking network module sets, because the
	// point of reading settings back is asserting a module configured the network it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/networks/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"compute#network",
			"name":"gw-library-test",
			"description":"created by terratest",
			"autoCreateSubnetworks":false,
			"routingConfig":{"routingMode":"REGIONAL"},
			"mtu":1500
		}`))
	})

	network, err := gcp.FetchNetworkWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", network.Name)
	assert.Equal(t, "created by terratest", network.Description)
	assert.False(t, network.AutoCreateSubnetworks)
	require.NotNil(t, network.RoutingConfig)
	assert.Equal(t, "REGIONAL", network.RoutingConfig.RoutingMode)
	assert.Equal(t, int64(1500), network.Mtu)
}

func TestFetchFirewallWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking firewall module sets. A rule whose
	// ports or direction came out wrong is the case most worth catching here.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/firewalls/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"kind":"compute#firewall",
			"name":"gw-library-test",
			"description":"created by terratest",
			"direction":"INGRESS",
			"priority":1000,
			"disabled":true,
			"sourceRanges":["10.0.0.0/8"],
			"allowed":[{"IPProtocol":"tcp","ports":["443"]}]
		}`))
	})

	firewall, err := gcp.FetchFirewallWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "INGRESS", firewall.Direction)
	assert.Equal(t, int64(1000), firewall.Priority)
	assert.True(t, firewall.Disabled)
	assert.Equal(t, []string{"10.0.0.0/8"}, firewall.SourceRanges)
	require.Len(t, firewall.Allowed, 1)
	assert.Equal(t, "tcp", firewall.Allowed[0].IPProtocol)
	assert.Equal(t, []string{"443"}, firewall.Allowed[0].Ports)
}

func TestFetchProjectMetadataWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-compute project metadata module sets, because the point of reading settings back
	// is asserting a module configured the metadata it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test-project","commonInstanceMetadata":{"fingerprint":"abc123","items":[{"key":"gw-library-test","value":"created by terratest"}]}}`))
	})

	metadata, err := gcp.FetchProjectMetadataWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	require.Len(t, metadata.Items, 1, "the project should hold exactly the item the fixture set")
	assert.Equal(t, "gw-library-test", metadata.Items[0].Key)
	require.NotNil(t, metadata.Items[0].Value)
	assert.Equal(t, "created by terratest", *metadata.Items[0].Value)
}

func TestFetchNetworkAttachmentWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an attachment the terraform-google-networking network attachment module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/networkAttachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","connectionPreference":"ACCEPT_MANUAL","subnetworks":["https://www.googleapis.com/compute/v1/projects/gw-library-test-project/regions/us-central1/subnetworks/gw-library-test"],"producerAcceptLists":["gw-library-test-project"]}`))
	})

	attachment, err := gcp.FetchNetworkAttachmentWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACCEPT_MANUAL", attachment.ConnectionPreference)
	assert.Equal(t, "created by terratest", attachment.Description)
	require.Len(t, attachment.Subnetworks, 1)
	assert.Equal(t, []string{"gw-library-test-project"}, attachment.ProducerAcceptLists)
}

func TestFetchNodeGroupWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a sole-tenant node group the terraform-google-compute node group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/nodeGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","size":1,"maintenancePolicy":"RESTART_IN_PLACE","nodeTemplate":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/regions/us-central1/nodeTemplates/gw-library-test"}`))
	})

	group, err := gcp.FetchNodeGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", group.Description)
	assert.Equal(t, "RESTART_IN_PLACE", group.MaintenancePolicy)
	assert.Equal(t, int64(1), group.Size)
}

func TestFetchNodeGroupWithClientMissingGroup(t *testing.T) {
	t.Parallel()

	// A caller who names a sole-tenant node group that is not there should read a sentence about that sole-tenant node group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchNodeGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchReservationWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a reservation the terraform-google-compute reservation module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/reservations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","specificReservationRequired":true,"status":"READY"}`))
	})

	reservation, err := gcp.FetchReservationWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", reservation.Description)
	assert.True(t, reservation.SpecificReservationRequired)
	assert.Equal(t, "READY", reservation.Status)
}

func TestFetchReservationWithClientMissingReservation(t *testing.T) {
	t.Parallel()

	// A caller who names a reservation that is not there should read a sentence about that reservation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchReservationWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchRegionCommitmentWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a commitment the terraform-google-compute region commitment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/commitments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","plan":"TWELVE_MONTH","type":"GENERAL_PURPOSE_N2","status":"ACTIVE"}`))
	})

	commitment, err := gcp.FetchRegionCommitmentWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", commitment.Description)
	assert.Equal(t, "TWELVE_MONTH", commitment.Plan)
	assert.Equal(t, "GENERAL_PURPOSE_N2", commitment.Type)
}

func TestFetchRegionCommitmentWithClientMissingCommitment(t *testing.T) {
	t.Parallel()

	// A caller who names a commitment that is not there should read a sentence about that commitment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchRegionCommitmentWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchStoragePoolWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a storage pool the terraform-google-compute storage pool module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/storagePools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","poolProvisionedCapacityGb":"10240","poolProvisionedIops":"10000","capacityProvisioningType":"ADVANCED"}`))
	})

	pool, err := gcp.FetchStoragePoolWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", pool.Description)
	assert.Equal(t, int64(10240), pool.PoolProvisionedCapacityGb)
	assert.Equal(t, "ADVANCED", pool.CapacityProvisioningType)
}

func TestFetchStoragePoolWithClientMissingPool(t *testing.T) {
	t.Parallel()

	// A caller who names a storage pool that is not there should read a sentence about that storage pool, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchStoragePoolWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchPublicAdvertisedPrefixWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a public advertised prefix the terraform-google-networking public advertised prefix module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/publicAdvertisedPrefixes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","ipCidrRange":"198.51.100.0/24","dnsVerificationIp":"198.51.100.1","status":"VALIDATED"}`))
	})

	prefix, err := gcp.FetchPublicAdvertisedPrefixWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", prefix.Description)
	assert.Equal(t, "198.51.100.0/24", prefix.IpCidrRange)
	assert.Equal(t, "VALIDATED", prefix.Status)
}

func TestFetchPublicAdvertisedPrefixWithClientMissingPrefix(t *testing.T) {
	t.Parallel()

	// A caller who names a public advertised prefix that is not there should read a sentence about that public advertised prefix, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchPublicAdvertisedPrefixWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchPublicDelegatedPrefixWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a public delegated prefix the terraform-google-networking public delegated prefix module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/publicDelegatedPrefixes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","ipCidrRange":"198.51.100.0/26","parentPrefix":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/publicAdvertisedPrefixes/gw-library-parent","isLiveMigration":true}`))
	})

	prefix, err := gcp.FetchPublicDelegatedPrefixWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", prefix.Description)
	assert.Equal(t, "198.51.100.0/26", prefix.IpCidrRange)
	assert.True(t, prefix.IsLiveMigration)
}

func TestFetchPublicDelegatedPrefixWithClientMissingPrefix(t *testing.T) {
	t.Parallel()

	// A caller who names a public delegated prefix that is not there should read a sentence about that public delegated prefix, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchPublicDelegatedPrefixWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchInterconnectWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a interconnect the terraform-google-networking interconnect module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/interconnects/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","linkType":"LINK_TYPE_ETHERNET_10G_LR","requestedLinkCount":2,"location":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/interconnectLocations/iad-zone1-1"}`))
	})

	interconnect, err := gcp.FetchInterconnectWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", interconnect.Description)
	assert.Equal(t, "LINK_TYPE_ETHERNET_10G_LR", interconnect.LinkType)
	assert.Equal(t, int64(2), interconnect.RequestedLinkCount)
}

func TestFetchInterconnectWithClientMissingInterconnect(t *testing.T) {
	t.Parallel()

	// A caller who names a interconnect that is not there should read a sentence about that interconnect, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchInterconnectWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchInterconnectGroupWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a interconnect group the terraform-google-networking interconnect group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/interconnectGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","intent":{"topologyCapability":"PRODUCTION_NON_CRITICAL"}}`))
	})

	group, err := gcp.FetchInterconnectGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", group.Description)
}

func TestFetchInterconnectGroupWithClientMissingGroup(t *testing.T) {
	t.Parallel()

	// A caller who names a interconnect group that is not there should read a sentence about that interconnect group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchInterconnectGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchInterconnectAttachmentWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a interconnect attachment the terraform-google-networking interconnect attachment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/interconnectAttachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","bandwidth":"BPS_1G","type":"PARTNER","router":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/regions/us-central1/routers/gw-library-test"}`))
	})

	attachment, err := gcp.FetchInterconnectAttachmentWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attachment.Description)
	assert.Equal(t, "BPS_1G", attachment.Bandwidth)
	assert.Equal(t, "PARTNER", attachment.Type)
}

func TestFetchInterconnectAttachmentWithClientMissingAttachment(t *testing.T) {
	t.Parallel()

	// A caller who names a interconnect attachment that is not there should read a sentence about that interconnect attachment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchInterconnectAttachmentWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchCrossSiteNetworkWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a cross-site network the terraform-google-networking cross-site network module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/crossSiteNetworks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest"}`))
	})

	network, err := gcp.FetchCrossSiteNetworkWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", network.Description)
}

func TestFetchCrossSiteNetworkWithClientMissingNetwork(t *testing.T) {
	t.Parallel()

	// A caller who names a cross-site network that is not there should read a sentence about that cross-site network, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchCrossSiteNetworkWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchWireGroupWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a wire group the terraform-google-networking wire group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/crossSiteNetworks/gw-library-parent/wireGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","adminEnabled":true}`))
	})

	group, err := gcp.FetchWireGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", group.Description)
	assert.True(t, group.AdminEnabled)
}

func TestFetchWireGroupWithClientMissingGroup(t *testing.T) {
	t.Parallel()

	// A caller who names a wire group that is not there should read a sentence about that wire group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchWireGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchPreviewFeatureWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a preview feature the terraform-google-compute preview feature module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/previewFeatures/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","activationStatus":"ENABLED"}`))
	})

	feature, err := gcp.FetchPreviewFeatureWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", feature.Description)
	assert.Equal(t, "ENABLED", feature.ActivationStatus)
}

func TestFetchPreviewFeatureWithClientMissingFeature(t *testing.T) {
	t.Parallel()

	// A caller who names a preview feature that is not there should read a sentence about that preview feature, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.FetchPreviewFeatureWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestFetchDiskIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a disk a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/disks/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.storageAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchDiskIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.storageAdmin", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchRegionDiskIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a regional disk a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/disks/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.storageAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchRegionDiskIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.storageAdmin", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchSnapshotIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a snapshot a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/snapshots/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.storageAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchSnapshotIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.storageAdmin", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchInstantSnapshotIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instant snapshot a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/instantSnapshots/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.storageAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchInstantSnapshotIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.storageAdmin", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchRegionInstantSnapshotIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a regional instant snapshot a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/instantSnapshots/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.storageAdmin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchRegionInstantSnapshotIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.storageAdmin", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchInstanceTemplateIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance template a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/instanceTemplates/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.viewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchInstanceTemplateIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.viewer", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchInstanceIamPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance a Gruntwork module set a policy
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/instances/gw-library-test/getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"), "a conditional binding only comes back at version 3")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/compute.osLogin","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.FetchInstanceIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/compute.osLogin", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}
func TestFetchInstantSnapshotWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instant snapshot a Gruntwork module created, not
	// a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/instantSnapshots/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","diskSizeGb":"10","labels":{"purpose":"terratest"},"status":"READY"}`))
	})

	result, err := gcp.FetchInstantSnapshotWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Equal(t, int64(10), result.DiskSizeGb)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
}
func TestFetchRegionInstantSnapshotWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a regional instant snapshot a Gruntwork module created, not
	// a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/instantSnapshots/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","diskSizeGb":"10","status":"READY"}`))
	})

	result, err := gcp.FetchRegionInstantSnapshotWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", result.Description)
	assert.Equal(t, int64(10), result.DiskSizeGb)
}
func TestFetchNodeTemplateWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a sole tenant node template a Gruntwork module created, not
	// a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/nodeTemplates/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","nodeType":"n1-node-96-624","cpuOvercommitType":"NONE","nodeAffinityLabels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.FetchNodeTemplateWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "n1-node-96-624", result.NodeType)
	assert.Equal(t, "NONE", result.CpuOvercommitType)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.NodeAffinityLabels)
}
func TestFetchRegionAutoscalerWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a regional autoscaler a Gruntwork module created, not
	// a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/autoscalers/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","autoscalingPolicy":{"minNumReplicas":1,"maxNumReplicas":3,"coolDownPeriodSec":90,"mode":"OFF"}}`))
	})

	result, err := gcp.FetchRegionAutoscalerWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, result.AutoscalingPolicy)
	assert.Equal(t, int64(3), result.AutoscalingPolicy.MaxNumReplicas)
	assert.Equal(t, int64(90), result.AutoscalingPolicy.CoolDownPeriodSec)
	assert.Equal(t, "OFF", result.AutoscalingPolicy.Mode)
}
func TestFetchDiskIamPolicyWithClientReportsAMissingDisk(t *testing.T) {
	t.Parallel()

	// A caller who asks for a disk that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"The resource was not found."}}`))
	})

	_, err := gcp.FetchDiskIamPolicyWithClient(context.Background(), newFakeComputeService(t, handler),
		"gw-library-test-project", "us-central1-a", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
func TestFetchInstanceResourcePoliciesWithClient(t *testing.T) {
	t.Parallel()

	// An attachment is not a resource of its own, so Google answers for it as a list on the instance.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","resourcePolicies":["https://www.googleapis.com/compute/v1/projects/gw-library-test-project/regions/us-central1/resourcePolicies/gw-library-test"]}`))
	})

	policies, err := gcp.FetchInstanceResourcePoliciesWithClient(context.Background(), newFakeComputeService(t, handler),
		"gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policies, 1)
	assert.Contains(t, policies[0], "resourcePolicies/gw-library-test")
}
