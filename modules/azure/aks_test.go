package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6"
	csfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeManagedClustersClient(t *testing.T, srv *csfake.ManagedClustersServer) *armcontainerservice.ManagedClustersClient {
	t.Helper()

	transport := csfake.NewManagedClustersServerTransport(srv)
	client, err := armcontainerservice.NewManagedClustersClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetManagedClusterWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  csfake.ManagedClustersServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: csfake.ManagedClustersServer{
				Get: func(_ context.Context, _ string, clusterName string, _ *armcontainerservice.ManagedClustersClientGetOptions) (resp azfake.Responder[armcontainerservice.ManagedClustersClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcontainerservice.ManagedClustersClientGetResponse{
						ManagedCluster: armcontainerservice.ManagedCluster{
							Name: to.Ptr(clusterName),
							Properties: &armcontainerservice.ManagedClusterProperties{
								KubernetesVersion: to.Ptr("1.28.0"),
							},
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: csfake.ManagedClustersServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcontainerservice.ManagedClustersClientGetOptions) (resp azfake.Responder[armcontainerservice.ManagedClustersClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")
					return
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := tc.server
			client := newFakeManagedClustersClient(t, &srv)

			cluster, err := azure.GetManagedClusterWithClient(t.Context(), client, "rg", "my-cluster")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-cluster", *cluster.Name)
			assert.Equal(t, "1.28.0", *cluster.Properties.KubernetesVersion)
		})
	}
}
