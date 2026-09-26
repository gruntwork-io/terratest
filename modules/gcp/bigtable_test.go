package gcp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigtableadmin/v2"
	"google.golang.org/api/option"
)

// newFakeBigtableAdminService points a real Bigtable client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigtableAdminService(t *testing.T, handler http.Handler) *bigtableadmin.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigtableadmin.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigtableInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test","displayName":"terratest instance","type":"DEVELOPMENT","state":"READY","labels":{"purpose":"terratest"}}`))
	})

	instance, err := gcp.GetBigtableInstanceAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest instance", instance.DisplayName)
	assert.Equal(t, "DEVELOPMENT", instance.Type)
	assert.Equal(t, "READY", instance.State)
	assert.Equal(t, "terratest", instance.Labels["purpose"])
}

func TestGetBigtableInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetBigtableInstanceAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetBigtableClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a cluster the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test/clusters/gw-library-test","location":"projects/gw-library-test-project/locations/us-central1-b","serveNodes":1,"defaultStorageType":"SSD","state":"READY"}`))
	})

	cluster, err := gcp.GetBigtableClusterAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler),
		"gw-library-test-project", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1-b", cluster.Location)
	assert.Equal(t, int64(1), cluster.ServeNodes)
	assert.Equal(t, "SSD", cluster.DefaultStorageType)
}

func TestGetBigtableClusterAttrsWithClientMissingCluster(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"not found"}}`))
	})

	_, err := gcp.GetBigtableClusterAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler),
		"gw-library-test-project", "gw-library-test", "gone")
	require.Error(t, err)
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetBigtableTableAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a table the terraform-google-data-storage table module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test/tables/gw-library-test"), "unexpected path %s", r.URL.Path)

		// The default answer carries no column families and no subset rows, so the read has to ask for
		// the full view and this is what fails if it stops doing so.
		assert.Equal(t, "FULL", r.URL.Query().Get("view"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test/tables/gw-library-test","columnFamilies":{"terratest":{"gcRule":{"maxNumVersions":3}}},"granularity":"MILLIS"}`))
	})

	table, err := gcp.GetBigtableTableAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	require.Contains(t, table.ColumnFamilies, "terratest", "the table should carry the column family the fixture gave it")
	require.NotNil(t, table.ColumnFamilies["terratest"].GcRule)
	assert.Equal(t, int64(3), table.ColumnFamilies["terratest"].GcRule.MaxNumVersions)
	assert.Equal(t, "MILLIS", table.Granularity)
}

func TestGetBigtableAppProfileAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a app profile the terraform-google-data-storage app profile module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test/appProfiles/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test/appProfiles/gw-library-test","description":"created by terratest","multiClusterRoutingUseAny":{}}`))
	})

	profile, err := gcp.GetBigtableAppProfileAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", profile.Description)
	require.NotNil(t, profile.MultiClusterRoutingUseAny, "the profile should route to any cluster")
}

func TestGetBigtableAuthorizedViewAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a authorized view the terraform-google-data-storage authorized view module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test/tables/gw-library-test/authorizedViews/gw-library-test"), "unexpected path %s", r.URL.Path)

		// The default answer carries no column families and no subset rows, so the read has to ask for
		// the full view and this is what fails if it stops doing so.
		assert.Equal(t, "FULL", r.URL.Query().Get("view"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test/tables/gw-library-test/authorizedViews/gw-library-test","deletionProtection":false,"subsetView":{"rowPrefixes":["dGVycmF0ZXN0"]}}`))
	})

	view, err := gcp.GetBigtableAuthorizedViewAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gw-library-test", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, view.SubsetView, "the view should carry the subset the fixture gave it")
	require.Len(t, view.SubsetView.RowPrefixes, 1)
	assert.False(t, view.DeletionProtection)
}

func TestGetBigtableTableIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a table IAM policy the terraform-google-data-storage table IAM module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test/tables/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)

		// A conditional binding only comes back at version 3, so the read has to ask for it, which for
		// this call means in the request body rather than in the query.
		var request struct {
			Options struct {
				RequestedPolicyVersion int64 `json:"requestedPolicyVersion"`
			} `json:"options"`
		}

		// assert rather than require: a failed require inside a handler stops the wrong goroutine.
		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&request)) {
			assert.Equal(t, int64(3), request.Options.RequestedPolicyVersion)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/bigtable.reader","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetBigtableTableIamPolicyAttrsWithClient(context.Background(), newFakeBigtableAdminService(t, handler), "gw-library-test-project", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/bigtable.reader", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}
