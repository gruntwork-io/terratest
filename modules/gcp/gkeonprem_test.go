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
	"google.golang.org/api/gkeonprem/v1"
	"google.golang.org/api/option"
)

// newFakeGKEOnPremService points a real Anthos on premises client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeGKEOnPremService(t *testing.T, handler http.Handler) *gkeonprem.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := gkeonprem.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetGKEOnPremBareMetalAdminClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an admin cluster a Gruntwork module registered, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/bareMetalAdminClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/bareMetalAdminClusters/gw-library-test","description":"terratest admin cluster","bareMetalVersion":"1.29.0","state":"RUNNING","networkConfig":{"islandModeCidr":{"serviceAddressCidrBlocks":["10.96.0.0/20"],"podAddressCidrBlocks":["192.168.0.0/16"]}}}`))
	})

	cluster, err := gcp.GetGKEOnPremBareMetalAdminClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest admin cluster", cluster.Description)
	assert.Equal(t, "1.29.0", cluster.BareMetalVersion)
	assert.Equal(t, "RUNNING", cluster.State)
	require.NotNil(t, cluster.NetworkConfig)
	require.NotNil(t, cluster.NetworkConfig.IslandModeCidr)
	assert.Equal(t, []string{"10.96.0.0/20"}, cluster.NetworkConfig.IslandModeCidr.ServiceAddressCidrBlocks)
}

func TestGetGKEOnPremBareMetalAdminClusterAttrsWithClientReportsAMissingCluster(t *testing.T) {
	t.Parallel()

	// A caller who names a cluster that is not there should read a sentence about that cluster, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEOnPremBareMetalAdminClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEOnPremBareMetalClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A user cluster reports the admin cluster it was created through, which is the field that says the two were
	// wired together.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us-central1/bareMetalClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/bareMetalClusters/gw-library-test","adminClusterName":"projects/gw-library-test-project/locations/us-central1/bareMetalAdminClusters/gw-library-admin","bareMetalVersion":"1.29.0","state":"RUNNING","nodeConfig":{"maxPodsPerNode":"110"}}`))
	})

	cluster, err := gcp.GetGKEOnPremBareMetalClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Contains(t, cluster.AdminClusterName, "/bareMetalAdminClusters/gw-library-admin")
	assert.Equal(t, "RUNNING", cluster.State)
	require.NotNil(t, cluster.NodeConfig)
	assert.Equal(t, int64(110), cluster.NodeConfig.MaxPodsPerNode)
}

func TestGetGKEOnPremBareMetalClusterAttrsWithClientReportsAMissingCluster(t *testing.T) {
	t.Parallel()

	// A caller who names a cluster that is not there should read a sentence about that cluster, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEOnPremBareMetalClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEOnPremBareMetalNodePoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A node pool is named under its cluster, so the path carries both and the test checks both.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/bareMetalClusters/gw-library-cluster/bareMetalNodePools/gw-library-pool"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/bareMetalClusters/gw-library-cluster/bareMetalNodePools/gw-library-pool","displayName":"terratest pool","state":"RUNNING","nodePoolConfig":{"nodeConfigs":[{"nodeIp":"10.200.0.11"}]}}`))
	})

	pool, err := gcp.GetGKEOnPremBareMetalNodePoolAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-cluster", "gw-library-pool")
	require.NoError(t, err)

	assert.Equal(t, "terratest pool", pool.DisplayName)
	require.NotNil(t, pool.NodePoolConfig)
	require.Len(t, pool.NodePoolConfig.NodeConfigs, 1)
	assert.Equal(t, "10.200.0.11", pool.NodePoolConfig.NodeConfigs[0].NodeIp)
}

func TestGetGKEOnPremBareMetalNodePoolAttrsWithClientReportsAMissingPool(t *testing.T) {
	t.Parallel()

	// The sentence has to name the pool and the cluster it was looked for on, because a pool name is only unique
	// within its cluster.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEOnPremBareMetalNodePoolAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-cluster", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "gw-library-cluster")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEOnPremVMwareAdminClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A VMware admin cluster reports the vCenter it runs on, which nothing in a module's own state holds.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us-central1/vmwareAdminClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/vmwareAdminClusters/gw-library-test","description":"terratest VMware admin cluster","onPremVersion":"1.16.0","state":"RUNNING","vcenter":{"datacenter":"terratest-dc","cluster":"terratest-cluster"}}`))
	})

	cluster, err := gcp.GetGKEOnPremVMwareAdminClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest VMware admin cluster", cluster.Description)
	assert.Equal(t, "1.16.0", cluster.OnPremVersion)
	require.NotNil(t, cluster.Vcenter)
	assert.Equal(t, "terratest-dc", cluster.Vcenter.Datacenter)
}

func TestGetGKEOnPremVMwareAdminClusterAttrsWithClientReportsAMissingCluster(t *testing.T) {
	t.Parallel()

	// A caller who names a cluster that is not there should read a sentence about that cluster, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEOnPremVMwareAdminClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEOnPremVMwareClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The control plane node count a user cluster settled on is reported here, and a create only asks for it.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us-central1/vmwareClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/vmwareClusters/gw-library-test","adminClusterName":"projects/gw-library-test-project/locations/us-central1/vmwareAdminClusters/gw-library-admin","onPremVersion":"1.16.0","state":"RUNNING","controlPlaneNode":{"cpus":"4","memory":"8192","replicas":"1"}}`))
	})

	cluster, err := gcp.GetGKEOnPremVMwareClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Contains(t, cluster.AdminClusterName, "/vmwareAdminClusters/gw-library-admin")
	require.NotNil(t, cluster.ControlPlaneNode)
	assert.Equal(t, int64(1), cluster.ControlPlaneNode.Replicas)
	assert.Equal(t, int64(4), cluster.ControlPlaneNode.Cpus)
}

func TestGetGKEOnPremVMwareClusterAttrsWithClientReportsAMissingCluster(t *testing.T) {
	t.Parallel()

	// A caller who names a cluster that is not there should read a sentence about that cluster, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEOnPremVMwareClusterAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEOnPremVMwareNodePoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A VMware node pool is named under its cluster, and the machine shape it runs is reported on the pool.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/vmwareClusters/gw-library-cluster/vmwareNodePools/gw-library-pool"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/vmwareClusters/gw-library-cluster/vmwareNodePools/gw-library-pool","displayName":"terratest pool","state":"RUNNING","config":{"cpus":"4","memoryMb":"8192","replicas":"3","imageType":"ubuntu_containerd"}}`))
	})

	pool, err := gcp.GetGKEOnPremVMwareNodePoolAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-cluster", "gw-library-pool")
	require.NoError(t, err)

	assert.Equal(t, "terratest pool", pool.DisplayName)
	require.NotNil(t, pool.Config)
	assert.Equal(t, int64(3), pool.Config.Replicas)
	assert.Equal(t, "ubuntu_containerd", pool.Config.ImageType)
}

func TestGetGKEOnPremVMwareNodePoolAttrsWithClientReportsAMissingPool(t *testing.T) {
	t.Parallel()

	// The sentence has to name the pool and the cluster it was looked for on, because a pool name is only unique
	// within its cluster.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEOnPremVMwareNodePoolAttrsWithClient(context.Background(), newFakeGKEOnPremService(t, handler), "gw-library-test-project", "us-central1", "gw-library-cluster", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "gw-library-cluster")
	assert.Contains(t, err.Error(), "does not exist")
}
