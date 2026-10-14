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
	"google.golang.org/api/config/v1"
	"google.golang.org/api/option"
)

// newFakeInfraManagerService points a real Infrastructure Manager client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeInfraManagerService(t *testing.T, handler http.Handler) *config.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := config.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetInfraManagerDeploymentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a deployment a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/deployments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"state":"ACTIVE","latestRevision":"projects/gw-library-test-project/locations/us-central1/deployments/gw-library-test/revisions/r-0","serviceAccount":"projects/gw-library-test-project/serviceAccounts/gw-library-test@gw-library-test-project.iam.gserviceaccount.com","terraformBlueprint":{"gcsSource":"gs://gw-library-test/blueprint.tgz"},"labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetInfraManagerDeploymentAttrsWithClient(context.Background(), newFakeInfraManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", result.State)
	assert.Contains(t, result.LatestRevision, "/revisions/r-0")
	require.NotNil(t, result.TerraformBlueprint)
	assert.Equal(t, "gs://gw-library-test/blueprint.tgz", result.TerraformBlueprint.GcsSource)
}

func TestGetInfraManagerDeploymentAttrsWithClientReportsAMissingDeployment(t *testing.T) {
	t.Parallel()

	// A caller who names a deployment that is not there should read a sentence about that deployment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetInfraManagerDeploymentAttrsWithClient(context.Background(), newFakeInfraManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
