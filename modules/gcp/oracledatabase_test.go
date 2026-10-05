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
	"google.golang.org/api/oracledatabase/v1"
)

// newFakeOracleDatabaseService points a real Oracle Database client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeOracleDatabaseService(t *testing.T, handler http.Handler) *oracledatabase.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := oracledatabase.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetOracleAutonomousDatabaseAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle autonomous database the terraform-google-data-storage Oracle autonomous database module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/autonomousDatabases/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/autonomousDatabases/gw-library-test","database":"terratest","network":"projects/gw-library-test-project/global/networks/gw-library-test","cidr":"10.98.0.0/24","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleAutonomousDatabaseAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest", attrs.Database)
	assert.Equal(t, "10.98.0.0/24", attrs.Cidr)
}

func TestGetOracleAutonomousDatabaseAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle autonomous database that is not there should read a sentence about that Oracle autonomous database, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleAutonomousDatabaseAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleCloudExadataInfrastructureAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle Exadata infrastructure the terraform-google-data-storage Oracle Exadata infrastructure module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/cloudExadataInfrastructures/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/cloudExadataInfrastructures/gw-library-test","displayName":"terratest infrastructure","gcpOracleZone":"us-central1-a-r1","properties":{"shape":"Exadata.X9M"},"labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleCloudExadataInfrastructureAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest infrastructure", attrs.DisplayName)
	assert.Equal(t, "us-central1-a-r1", attrs.GcpOracleZone)
}

func TestGetOracleCloudExadataInfrastructureAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle Exadata infrastructure that is not there should read a sentence about that Oracle Exadata infrastructure, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleCloudExadataInfrastructureAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleCloudVMClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle Exadata VM cluster the terraform-google-data-storage Oracle Exadata VM cluster module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/cloudVmClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/cloudVmClusters/gw-library-test","displayName":"terratest cluster","exadataInfrastructure":"projects/gw-library-test-project/locations/us-central1/cloudExadataInfrastructures/gw-library-parent","cidr":"10.98.1.0/24","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleCloudVMClusterAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest cluster", attrs.DisplayName)
	assert.Equal(t, "10.98.1.0/24", attrs.Cidr)
}

func TestGetOracleCloudVMClusterAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle Exadata VM cluster that is not there should read a sentence about that Oracle Exadata VM cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleCloudVMClusterAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleDBSystemAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle database system the terraform-google-data-storage Oracle database system module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dbSystems/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/dbSystems/gw-library-test","displayName":"terratest system","gcpOracleZone":"us-central1-a-r1","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleDBSystemAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest system", attrs.DisplayName)
	assert.Equal(t, "us-central1-a-r1", attrs.GcpOracleZone)
}

func TestGetOracleDBSystemAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle database system that is not there should read a sentence about that Oracle database system, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleDBSystemAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleExadbVMClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle Exascale VM cluster the terraform-google-data-storage Oracle Exascale VM cluster module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/exadbVmClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/exadbVmClusters/gw-library-test","displayName":"terratest cluster","gcpOracleZone":"us-central1-a-r1","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleExadbVMClusterAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest cluster", attrs.DisplayName)
	assert.Equal(t, "us-central1-a-r1", attrs.GcpOracleZone)
}

func TestGetOracleExadbVMClusterAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle Exascale VM cluster that is not there should read a sentence about that Oracle Exascale VM cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleExadbVMClusterAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleExascaleDBStorageVaultAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle Exascale storage vault the terraform-google-data-storage Oracle Exascale storage vault module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/exascaleDbStorageVaults/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/exascaleDbStorageVaults/gw-library-test","displayName":"terratest vault","gcpOracleZone":"us-central1-a-r1","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleExascaleDBStorageVaultAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest vault", attrs.DisplayName)
	assert.Equal(t, "us-central1-a-r1", attrs.GcpOracleZone)
}

func TestGetOracleExascaleDBStorageVaultAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle Exascale storage vault that is not there should read a sentence about that Oracle Exascale storage vault, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleExascaleDBStorageVaultAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleODBNetworkAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle ODB network the terraform-google-data-storage Oracle ODB network module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/odbNetworks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/odbNetworks/gw-library-test","network":"projects/gw-library-test-project/global/networks/gw-library-test","gcpOracleZone":"us-central1-a-r1","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetOracleODBNetworkAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/global/networks/gw-library-test", attrs.Network)
}

func TestGetOracleODBNetworkAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle ODB network that is not there should read a sentence about that Oracle ODB network, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleODBNetworkAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetOracleODBSubnetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Oracle ODB subnet the terraform-google-data-storage Oracle ODB subnet module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/odbNetworks/gw-library-parent/odbSubnets/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/odbNetworks/gw-library-parent/odbSubnets/gw-library-test","cidrRange":"10.98.2.0/24","purpose":"CLIENT_SUBNET"}`))
	})

	attrs, err := gcp.GetOracleODBSubnetAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "10.98.2.0/24", attrs.CidrRange)
	assert.Equal(t, "CLIENT_SUBNET", attrs.Purpose)
}

func TestGetOracleODBSubnetAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Oracle ODB subnet that is not there should read a sentence about that Oracle ODB subnet, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetOracleODBSubnetAttrsWithClient(context.Background(), newFakeOracleDatabaseService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
