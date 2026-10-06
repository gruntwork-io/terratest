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
	"google.golang.org/api/apigee/v1"
	"google.golang.org/api/option"
)

// newFakeApigeeService points a real Apigee client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeApigeeService(t *testing.T, handler http.Handler) *apigee.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := apigee.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetApigeeOrganizationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee organization the terraform-google-api-management Apigee organization module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test-org","description":"created by terratest","runtimeType":"CLOUD","analyticsRegion":"us-central1","billingType":"PAYG"}`))
	})

	attrs, err := gcp.GetApigeeOrganizationAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "CLOUD", attrs.RuntimeType)
	assert.Equal(t, "us-central1", attrs.AnalyticsRegion)
}

func TestGetApigeeOrganizationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee organization that is not there should read a sentence about that Apigee organization, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeOrganizationAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeAPIProxyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee API proxy the terraform-google-api-management Apigee API proxy module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/apis/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","revision":["1"],"labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetApigeeAPIProxyAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
}

func TestGetApigeeAPIProxyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee API proxy that is not there should read a sentence about that Apigee API proxy, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeAPIProxyAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeAPIProductAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee API product the terraform-google-api-management Apigee API product module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/apiproducts/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","displayName":"terratest product","approvalType":"manual","environments":["gw-library-env"],"description":"created by terratest"}`))
	})

	attrs, err := gcp.GetApigeeAPIProductAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest product", attrs.DisplayName)
	assert.Equal(t, "manual", attrs.ApprovalType)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetApigeeAPIProductAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee API product that is not there should read a sentence about that Apigee API product, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeAPIProductAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeAppGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee app group the terraform-google-api-management Apigee app group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/appgroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","displayName":"terratest group","channelId":"terratest","status":"active"}`))
	})

	attrs, err := gcp.GetApigeeAppGroupAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest group", attrs.DisplayName)
	assert.Equal(t, "terratest", attrs.ChannelId)
	assert.Equal(t, "active", attrs.Status)
}

func TestGetApigeeAppGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee app group that is not there should read a sentence about that Apigee app group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeAppGroupAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeDeveloperAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee developer the terraform-google-api-management Apigee developer module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/developers/terratest@example.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"email":"terratest@example.com","firstName":"Terra","lastName":"Test","userName":"terratest","status":"active"}`))
	})

	attrs, err := gcp.GetApigeeDeveloperAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "terratest@example.com")
	require.NoError(t, err)

	assert.Equal(t, "terratest@example.com", attrs.Email)
	assert.Equal(t, "terratest", attrs.UserName)
	assert.Equal(t, "active", attrs.Status)
}

func TestGetApigeeDeveloperAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee developer that is not there should read a sentence about that Apigee developer, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeDeveloperAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeDeveloperAppAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee developer app the terraform-google-api-management Apigee developer app module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/developers/terratest@example.com/apps/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","appId":"6d4e2e6a","status":"approved","callbackUrl":"https://terratest.example.com/callback"}`))
	})

	attrs, err := gcp.GetApigeeDeveloperAppAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "terratest@example.com", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
	assert.Equal(t, "approved", attrs.Status)
	assert.Equal(t, "https://terratest.example.com/callback", attrs.CallbackUrl)
}

func TestGetApigeeDeveloperAppAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee developer app that is not there should read a sentence about that Apigee developer app, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeDeveloperAppAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "terratest@example.com", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment the terraform-google-api-management Apigee environment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","displayName":"terratest environment","description":"created by terratest","deploymentType":"PROXY","apiProxyType":"PROGRAMMABLE"}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest environment", attrs.DisplayName)
	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "PROXY", attrs.DeploymentType)
}

func TestGetApigeeEnvironmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment that is not there should read a sentence about that Apigee environment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment group the terraform-google-api-management Apigee environment group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/envgroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","hostnames":["terratest.example.com"],"state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentGroupAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetApigeeEnvironmentGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment group that is not there should read a sentence about that Apigee environment group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentGroupAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentGroupAttachmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment group attachment the terraform-google-api-management Apigee environment group attachment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/envgroups/gw-library-parent/attachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","environment":"gw-library-env"}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentGroupAttachmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-env", attrs.Environment)
}

func TestGetApigeeEnvironmentGroupAttachmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment group attachment that is not there should read a sentence about that Apigee environment group attachment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentGroupAttachmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee instance the terraform-google-api-management Apigee instance module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","location":"us-central1","description":"created by terratest","peeringCidrRange":"SLASH_22","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetApigeeInstanceAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "us-central1", attrs.Location)
	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "SLASH_22", attrs.PeeringCidrRange)
}

func TestGetApigeeInstanceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee instance that is not there should read a sentence about that Apigee instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeInstanceAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeInstanceAttachmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee instance attachment the terraform-google-api-management Apigee instance attachment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/instances/gw-library-parent/attachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","environment":"gw-library-env"}`))
	})

	attrs, err := gcp.GetApigeeInstanceAttachmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-env", attrs.Environment)
}

func TestGetApigeeInstanceAttachmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee instance attachment that is not there should read a sentence about that Apigee instance attachment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeInstanceAttachmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeNatAddressAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee NAT address the terraform-google-api-management Apigee NAT address module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/instances/gw-library-parent/natAddresses/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","ipAddress":"198.51.100.5","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetApigeeNatAddressAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "198.51.100.5", attrs.IpAddress)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetApigeeNatAddressAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee NAT address that is not there should read a sentence about that Apigee NAT address, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeNatAddressAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEndpointAttachmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee endpoint attachment the terraform-google-api-management Apigee endpoint attachment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/endpointAttachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","location":"us-central1","serviceAttachment":"projects/gw-library-test-project/regions/us-central1/serviceAttachments/gw-library-test","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetApigeeEndpointAttachmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "us-central1", attrs.Location)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetApigeeEndpointAttachmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee endpoint attachment that is not there should read a sentence about that Apigee endpoint attachment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEndpointAttachmentAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeDNSZoneAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee DNS zone the terraform-google-api-management Apigee DNS zone module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/dnsZones/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","domain":"terratest.example.com","description":"created by terratest","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetApigeeDNSZoneAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest.example.com", attrs.Domain)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetApigeeDNSZoneAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee DNS zone that is not there should read a sentence about that Apigee DNS zone, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeDNSZoneAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeDataCollectorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee data collector the terraform-google-api-management Apigee data collector module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/datacollectors/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","type":"STRING"}`))
	})

	attrs, err := gcp.GetApigeeDataCollectorAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "STRING", attrs.Type)
}

func TestGetApigeeDataCollectorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee data collector that is not there should read a sentence about that Apigee data collector, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeDataCollectorAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeDatastoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee datastore the terraform-google-api-management Apigee datastore module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/analytics/datastores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"self":"organizations/gw-library-test-org/analytics/datastores/gw-library-test","displayName":"terratest datastore","targetType":"gcs"}`))
	})

	attrs, err := gcp.GetApigeeDatastoreAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest datastore", attrs.DisplayName)
	assert.Equal(t, "gcs", attrs.TargetType)
}

func TestGetApigeeDatastoreAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee datastore that is not there should read a sentence about that Apigee datastore, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeDatastoreAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeSharedFlowAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee shared flow the terraform-google-api-management Apigee shared flow module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/sharedflows/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","revision":["1"],"latestRevisionId":"1"}`))
	})

	attrs, err := gcp.GetApigeeSharedFlowAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
	assert.Equal(t, "1", attrs.LatestRevisionId)
}

func TestGetApigeeSharedFlowAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee shared flow that is not there should read a sentence about that Apigee shared flow, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeSharedFlowAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeSpaceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee space the terraform-google-api-management Apigee space module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/spaces/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/spaces/gw-library-test","displayName":"terratest space"}`))
	})

	attrs, err := gcp.GetApigeeSpaceAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest space", attrs.DisplayName)
}

func TestGetApigeeSpaceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee space that is not there should read a sentence about that Apigee space, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeSpaceAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeTargetServerAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee target server the terraform-google-api-management Apigee target server module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/targetservers/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","host":"backend.terratest.example.com","port":8443,"description":"created by terratest","isEnabled":true}`))
	})

	attrs, err := gcp.GetApigeeTargetServerAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "backend.terratest.example.com", attrs.Host)
	assert.Equal(t, int64(8443), attrs.Port)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetApigeeTargetServerAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee target server that is not there should read a sentence about that Apigee target server, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeTargetServerAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeKeystoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee keystore the terraform-google-api-management Apigee keystore module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/keystores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","aliases":["gw-library-alias"]}`))
	})

	attrs, err := gcp.GetApigeeKeystoreAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
}

func TestGetApigeeKeystoreAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee keystore that is not there should read a sentence about that Apigee keystore, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeKeystoreAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeKeystoreAliasAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee keystore alias the terraform-google-api-management Apigee keystore alias module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/keystores/gw-library-parent/aliases/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"alias":"gw-library-test","type":"CERT"}`))
	})

	attrs, err := gcp.GetApigeeKeystoreAliasAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Alias)
	assert.Equal(t, "CERT", attrs.Type)
}

func TestGetApigeeKeystoreAliasAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee keystore alias that is not there should read a sentence about that Apigee keystore alias, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeKeystoreAliasAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeReferenceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee reference the terraform-google-api-management Apigee reference module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/references/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","refers":"gw-library-parent","resourceType":"KeyStore","description":"created by terratest"}`))
	})

	attrs, err := gcp.GetApigeeReferenceAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-parent", attrs.Refers)
	assert.Equal(t, "KeyStore", attrs.ResourceType)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetApigeeReferenceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee reference that is not there should read a sentence about that Apigee reference, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeReferenceAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeFlowHookAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee flow hook the terraform-google-api-management Apigee flow hook module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/flowhooks/PreProxyFlowHook"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"flowHookPoint":"PreProxyFlowHook","sharedFlow":"gw-library-flow","description":"created by terratest","continueOnError":true}`))
	})

	attrs, err := gcp.GetApigeeFlowHookAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "PreProxyFlowHook")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-flow", attrs.SharedFlow)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetApigeeFlowHookAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee flow hook that is not there should read a sentence about that Apigee flow hook, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeFlowHookAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentKeyValueMapAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment key value map the terraform-google-api-management Apigee environment key value map module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/keyvaluemaps/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","encrypted":true}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentKeyValueMapAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
	assert.True(t, attrs.Encrypted)
}

func TestGetApigeeEnvironmentKeyValueMapAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment key value map that is not there should read a sentence about that Apigee environment key value map, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentKeyValueMapAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentKeyValueEntryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment key value entry the terraform-google-api-management Apigee environment key value entry module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/keyvaluemaps/gw-library-parent/entries/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","value":"created by terratest"}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentKeyValueEntryAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", attrs.Name)
	assert.Equal(t, "created by terratest", attrs.Value)
}

func TestGetApigeeEnvironmentKeyValueEntryAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment key value entry that is not there should read a sentence about that Apigee environment key value entry, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentKeyValueEntryAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentAddonsConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment addons config the terraform-google-api-management Apigee environment addons config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-test/addonsConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"analyticsConfig":{"enabled":true,"state":"ENABLED"}}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentAddonsConfigAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ENABLED", attrs.AnalyticsConfig.State)
}

func TestGetApigeeEnvironmentAddonsConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment addons config that is not there should read a sentence about that Apigee environment addons config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentAddonsConfigAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeEnvironmentDebugMaskAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee environment debug mask the terraform-google-api-management Apigee environment debug mask module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-test/debugmask"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/environments/gw-library-test/debugmask","requestXpaths":["/terratest"],"namespaces":{"terratest":"https://terratest.example.com"}}`))
	})

	attrs, err := gcp.GetApigeeEnvironmentDebugMaskAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "/terratest", attrs.RequestXPaths[0])
}

func TestGetApigeeEnvironmentDebugMaskAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee environment debug mask that is not there should read a sentence about that Apigee environment debug mask, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeEnvironmentDebugMaskAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeSecurityActionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee security action the terraform-google-api-management Apigee security action module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/environments/gw-library-env/securityActions/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/environments/gw-library-env/securityActions/gw-library-test","description":"created by terratest","state":"ENABLED"}`))
	})

	attrs, err := gcp.GetApigeeSecurityActionAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "ENABLED", attrs.State)
}

func TestGetApigeeSecurityActionAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee security action that is not there should read a sentence about that Apigee security action, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeSecurityActionAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-env", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeSecurityFeedbackAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee security feedback the terraform-google-api-management Apigee security feedback module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/securityFeedback/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/securityFeedback/gw-library-test"}`))
	})

	attrs, err := gcp.GetApigeeSecurityFeedbackAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "organizations/gw-library-test-org/securityFeedback/gw-library-test", attrs.Name)
}

func TestGetApigeeSecurityFeedbackAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee security feedback that is not there should read a sentence about that Apigee security feedback, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeSecurityFeedbackAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeSecurityMonitoringConditionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee security monitoring condition the terraform-google-api-management Apigee security monitoring condition module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/securityMonitoringConditions/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/securityMonitoringConditions/gw-library-test","profile":"gw-library-profile","include":{"environments":["gw-library-env"]}}`))
	})

	attrs, err := gcp.GetApigeeSecurityMonitoringConditionAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-profile", attrs.Profile)
}

func TestGetApigeeSecurityMonitoringConditionAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee security monitoring condition that is not there should read a sentence about that Apigee security monitoring condition, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeSecurityMonitoringConditionAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeSecurityProfileV2AttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee security profile the terraform-google-api-management Apigee security profile module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/securityProfilesV2/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/securityProfilesV2/gw-library-test","description":"created by terratest"}`))
	})

	attrs, err := gcp.GetApigeeSecurityProfileV2AttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetApigeeSecurityProfileV2AttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee security profile that is not there should read a sentence about that Apigee security profile, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeSecurityProfileV2AttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetApigeeControlPlaneAccessAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Apigee control plane access the terraform-google-api-management Apigee control plane access module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/organizations/gw-library-test-org/controlPlaneAccess"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"organizations/gw-library-test-org/controlPlaneAccess","synchronizerIdentities":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}`))
	})

	attrs, err := gcp.GetApigeeControlPlaneAccessAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-test-org")
	require.NoError(t, err)

	assert.Equal(t, "serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com", attrs.SynchronizerIdentities[0])
}

func TestGetApigeeControlPlaneAccessAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Apigee control plane access that is not there should read a sentence about that Apigee control plane access, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetApigeeControlPlaneAccessAttrsWithClient(context.Background(), newFakeApigeeService(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
