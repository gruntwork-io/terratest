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
	"google.golang.org/api/beyondcorp/v1"
	"google.golang.org/api/option"
)

// newFakeBeyondCorpService points a real BeyondCorp client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeBeyondCorpService(t *testing.T, handler http.Handler) *beyondcorp.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := beyondcorp.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetBeyondCorpAppConnectionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp app connection the terraform-google-security BeyondCorp app connection module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/appConnections/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/appConnections/gw-library-test","displayName":"terratest connection","type":"TCP_PROXY","state":"CREATED","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetBeyondCorpAppConnectionAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest connection", attrs.DisplayName)
	assert.Equal(t, "TCP_PROXY", attrs.Type)
	assert.Equal(t, "CREATED", attrs.State)
}

func TestGetBeyondCorpAppConnectionAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp app connection that is not there should read a sentence about that BeyondCorp app connection, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpAppConnectionAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBeyondCorpAppConnectorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp app connector the terraform-google-security BeyondCorp app connector module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/appConnectors/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/appConnectors/gw-library-test","displayName":"terratest connector","state":"CREATED","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetBeyondCorpAppConnectorAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest connector", attrs.DisplayName)
	assert.Equal(t, "CREATED", attrs.State)
}

func TestGetBeyondCorpAppConnectorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp app connector that is not there should read a sentence about that BeyondCorp app connector, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpAppConnectorAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBeyondCorpAppGatewayAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp app gateway the terraform-google-security BeyondCorp app gateway module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/appGateways/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/appGateways/gw-library-test","displayName":"terratest gateway","type":"TCP_PROXY","hostType":"GCP_REGIONAL_MIG","state":"CREATED","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetBeyondCorpAppGatewayAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest gateway", attrs.DisplayName)
	assert.Equal(t, "TCP_PROXY", attrs.Type)
	assert.Equal(t, "GCP_REGIONAL_MIG", attrs.HostType)
}

func TestGetBeyondCorpAppGatewayAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp app gateway that is not there should read a sentence about that BeyondCorp app gateway, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpAppGatewayAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBeyondCorpSecurityGatewayAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp security gateway the terraform-google-security BeyondCorp security gateway module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/securityGateways/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/securityGateways/gw-library-test","displayName":"terratest gateway","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetBeyondCorpSecurityGatewayAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest gateway", attrs.DisplayName)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetBeyondCorpSecurityGatewayAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp security gateway that is not there should read a sentence about that BeyondCorp security gateway, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpSecurityGatewayAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBeyondCorpSecurityGatewayApplicationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp security gateway application the terraform-google-security BeyondCorp security gateway application module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/securityGateways/gw-library-parent/applications/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/securityGateways/gw-library-parent/applications/gw-library-test","displayName":"terratest application","endpointMatchers":[{"hostname":"terratest.example.com"}]}`))
	})

	attrs, err := gcp.GetBeyondCorpSecurityGatewayApplicationAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest application", attrs.DisplayName)
	assert.Equal(t, "terratest.example.com", attrs.EndpointMatchers[0].Hostname)
}

func TestGetBeyondCorpSecurityGatewayApplicationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp security gateway application that is not there should read a sentence about that BeyondCorp security gateway application, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpSecurityGatewayApplicationAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBeyondCorpSecurityGatewayIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp security gateway the terraform-google-security BeyondCorp security gateway module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/securityGateways/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/beyondcorp.securityGatewayUser","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetBeyondCorpSecurityGatewayIamPolicyAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/beyondcorp.securityGatewayUser", policy.Bindings[0].Role)
}

func TestGetBeyondCorpSecurityGatewayIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp security gateway that is not there should read a sentence about that BeyondCorp security gateway, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpSecurityGatewayIamPolicyAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BeyondCorp security gateway application the terraform-google-security BeyondCorp security gateway application module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/securityGateways/gw-library-parent/applications/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/beyondcorp.securityGatewayUser","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/beyondcorp.securityGatewayUser", policy.Bindings[0].Role)
}

func TestGetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a BeyondCorp security gateway application that is not there should read a sentence about that BeyondCorp security gateway application, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBeyondCorpSecurityGatewayApplicationIamPolicyAttrsWithClient(context.Background(), newFakeBeyondCorpService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
