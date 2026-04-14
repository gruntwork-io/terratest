package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	cifake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	crfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeRegistriesClient(t *testing.T, srv *crfake.RegistriesServer) *armcontainerregistry.RegistriesClient {
	t.Helper()

	transport := crfake.NewRegistriesServerTransport(srv)
	client, err := armcontainerregistry.NewRegistriesClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakeContainerGroupsClient(t *testing.T, srv *cifake.ContainerGroupsServer) *armcontainerinstance.ContainerGroupsClient {
	t.Helper()

	transport := cifake.NewContainerGroupsServerTransport(srv)
	client, err := armcontainerinstance.NewContainerGroupsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetContainerRegistryWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  crfake.RegistriesServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: crfake.RegistriesServer{
				Get: func(_ context.Context, _ string, registryName string, _ *armcontainerregistry.RegistriesClientGetOptions) (resp azfake.Responder[armcontainerregistry.RegistriesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcontainerregistry.RegistriesClientGetResponse{
						Registry: armcontainerregistry.Registry{
							Name: to.Ptr(registryName),
							Properties: &armcontainerregistry.RegistryProperties{
								LoginServer: to.Ptr(registryName + ".azurecr.io"),
							},
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: crfake.RegistriesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcontainerregistry.RegistriesClientGetOptions) (resp azfake.Responder[armcontainerregistry.RegistriesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeRegistriesClient(t, &srv)

			registry, err := azure.GetContainerRegistryWithClient(t.Context(), client, "rg", "my-registry")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-registry", *registry.Name)
			assert.Equal(t, "my-registry.azurecr.io", *registry.Properties.LoginServer)
		})
	}
}

func TestGetContainerInstanceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  cifake.ContainerGroupsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: cifake.ContainerGroupsServer{
				Get: func(_ context.Context, _ string, containerGroupName string, _ *armcontainerinstance.ContainerGroupsClientGetOptions) (resp azfake.Responder[armcontainerinstance.ContainerGroupsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcontainerinstance.ContainerGroupsClientGetResponse{
						ContainerGroup: armcontainerinstance.ContainerGroup{
							Name: to.Ptr(containerGroupName),
							Properties: &armcontainerinstance.ContainerGroupPropertiesProperties{
								OSType: to.Ptr(armcontainerinstance.OperatingSystemTypesLinux),
							},
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: cifake.ContainerGroupsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcontainerinstance.ContainerGroupsClientGetOptions) (resp azfake.Responder[armcontainerinstance.ContainerGroupsClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeContainerGroupsClient(t, &srv)

			instance, err := azure.GetContainerInstanceWithClient(t.Context(), client, "rg", "my-container")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-container", *instance.Name)
			assert.Equal(t, armcontainerinstance.OperatingSystemTypesLinux, *instance.Properties.OSType)
		})
	}
}
