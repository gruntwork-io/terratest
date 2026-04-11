package azure_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// BuildNamespaceNamesList tests
// ---------------------------------------------------------------------------

func TestBuildNamespaceNamesList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		namespaces []*armservicebus.SBNamespace
		want       []string
	}{
		{
			name: "MultipleNamespaces",
			namespaces: []*armservicebus.SBNamespace{
				{Name: to.Ptr("ns-1")},
				{Name: to.Ptr("ns-2")},
				{Name: to.Ptr("ns-3")},
			},
			want: []string{"ns-1", "ns-2", "ns-3"},
		},
		{
			name:       "EmptyList",
			namespaces: []*armservicebus.SBNamespace{},
			want:       []string{},
		},
		{
			name: "SingleNamespace",
			namespaces: []*armservicebus.SBNamespace{
				{Name: to.Ptr("only-ns")},
			},
			want: []string{"only-ns"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := azure.BuildNamespaceNamesList(tc.namespaces)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// BuildNamespaceIdsList tests
// ---------------------------------------------------------------------------

func TestBuildNamespaceIdsList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		namespaces []*armservicebus.SBNamespace
		want       []string
	}{
		{
			name: "MultipleNamespaces",
			namespaces: []*armservicebus.SBNamespace{
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-1")},
				{ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-2")},
			},
			want: []string{
				"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-1",
				"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.ServiceBus/namespaces/ns-2",
			},
		},
		{
			name:       "EmptyList",
			namespaces: []*armservicebus.SBNamespace{},
			want:       []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := azure.BuildNamespaceIdsList(tc.namespaces)
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---------------------------------------------------------------------------
// ListServiceBusNamespaceContextE error handling (empty subscription)
// ---------------------------------------------------------------------------

func TestListServiceBusNamespaceContextE_EmptySubscription(t *testing.T) {
	t.Parallel()

	_, err := azure.ListServiceBusNamespaceContextE(t.Context(), "")
	assert.Error(t, err)
}

// ---------------------------------------------------------------------------
// ListNamespaceTopicsContextE error handling (empty subscription)
// ---------------------------------------------------------------------------

func TestListNamespaceTopicsContextE_EmptySubscription(t *testing.T) {
	t.Parallel()

	_, err := azure.ListNamespaceTopicsContextE(t.Context(), "", "", "")
	assert.Error(t, err)
}
