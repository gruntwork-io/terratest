package azure_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListServiceBusNamespaceWithClient(t *testing.T) {
	t.Parallel()

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := armservicebus.SBNamespaceListResult{
			Value: []*armservicebus.SBNamespace{
				{Name: to.Ptr("ns-alpha"), ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-alpha")},
				{Name: to.Ptr("ns-beta"), ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-beta")},
			},
		}

		// httptest used because the servicebus beta SDK (v2.0.0-beta.3) doesn't include azfake support
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	client, err := armservicebus.NewNamespacesClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: srv.Client(),
			Cloud: cloud.Configuration{
				Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
					cloud.ResourceManager: {Endpoint: srv.URL, Audience: srv.URL},
				},
			},
		},
	})
	require.NoError(t, err)

	results, err := azure.ListServiceBusNamespaceWithClient(t.Context(), client)
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "ns-alpha", *results[0].Name)
	assert.Equal(t, "ns-beta", *results[1].Name)
}

func TestBuildNamespaceNamesList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		namespaces []*armservicebus.SBNamespace
		expected   []string
	}{
		{
			name:       "empty list",
			namespaces: []*armservicebus.SBNamespace{},
			expected:   []string{},
		},
		{
			name: "multiple namespaces",
			namespaces: []*armservicebus.SBNamespace{
				{Name: to.Ptr("ns-alpha")},
				{Name: to.Ptr("ns-beta")},
				{Name: to.Ptr("ns-gamma")},
			},
			expected: []string{"ns-alpha", "ns-beta", "ns-gamma"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := azure.BuildNamespaceNamesList(tc.namespaces)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBuildNamespaceIdsList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		namespaces []*armservicebus.SBNamespace
		expected   []string
	}{
		{
			name:       "empty list",
			namespaces: []*armservicebus.SBNamespace{},
			expected:   []string{},
		},
		{
			name: "multiple namespaces",
			namespaces: []*armservicebus.SBNamespace{
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-1")},
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-2")},
			},
			expected: []string{
				"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-1",
				"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-2",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := azure.BuildNamespaceIdsList(tc.namespaces)
			assert.Equal(t, tc.expected, result)
		})
	}
}
