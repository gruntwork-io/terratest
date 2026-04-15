package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	resfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeResourceGroupsClient(t *testing.T, srv *resfake.ResourceGroupsServer) *armresources.ResourceGroupsClient {
	t.Helper()

	transport := resfake.NewResourceGroupsServerTransport(srv)
	client, err := armresources.NewResourceGroupsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetAResourceGroupWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  resfake.ResourceGroupsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: resfake.ResourceGroupsServer{
				Get: func(_ context.Context, resourceGroupName string, _ *armresources.ResourceGroupsClientGetOptions) (resp azfake.Responder[armresources.ResourceGroupsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armresources.ResourceGroupsClientGetResponse{
						ResourceGroup: armresources.ResourceGroup{
							Name:     to.Ptr(resourceGroupName),
							Location: to.Ptr("eastus"),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: resfake.ResourceGroupsServer{
				Get: func(_ context.Context, _ string, _ *armresources.ResourceGroupsClientGetOptions) (resp azfake.Responder[armresources.ResourceGroupsClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceGroupNotFound")
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
			client := newFakeResourceGroupsClient(t, &srv)

			rg, err := azure.GetAResourceGroupWithClient(t.Context(), client, "my-rg")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-rg", *rg.Name)
			assert.Equal(t, "eastus", *rg.Location)
		})
	}
}

func TestListResourceGroupsByTagWithClient(t *testing.T) {
	t.Parallel()

	srv := &resfake.ResourceGroupsServer{
		NewListPager: func(_ *armresources.ResourceGroupsClientListOptions) (resp azfake.PagerResponder[armresources.ResourceGroupsClientListResponse]) {
			resp.AddPage(http.StatusOK, armresources.ResourceGroupsClientListResponse{
				ResourceGroupListResult: armresources.ResourceGroupListResult{
					Value: []*armresources.ResourceGroup{
						{Name: to.Ptr("rg-alpha"), Location: to.Ptr("eastus")},
						{Name: to.Ptr("rg-beta"), Location: to.Ptr("westus")},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeResourceGroupsClient(t, srv)

	results, err := azure.ListResourceGroupsByTagWithClient(t.Context(), client, "env")
	require.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "rg-alpha", *results[0].Name)
	assert.Equal(t, "rg-beta", *results[1].Name)
}
