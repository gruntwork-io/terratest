package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	computefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeDisksClient(t *testing.T, srv *computefake.DisksServer) *armcompute.DisksClient {
	t.Helper()

	transport := computefake.NewDisksServerTransport(srv)
	client, err := armcompute.NewDisksClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetDiskWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  computefake.DisksServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: computefake.DisksServer{
				Get: func(_ context.Context, _ string, diskName string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcompute.DisksClientGetResponse{
						Disk: armcompute.Disk{
							Name: to.Ptr(diskName),
							Properties: &armcompute.DiskProperties{
								DiskSizeGB: to.Ptr[int32](128),
							},
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: computefake.DisksServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeDisksClient(t, &srv)

			disk, err := azure.GetDiskWithClient(t.Context(), client, "rg", "my-disk")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-disk", *disk.Name)
			assert.Equal(t, int32(128), *disk.Properties.DiskSizeGB)
		})
	}
}

func TestDiskExistsWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		get     func(context.Context, string, string, *armcompute.DisksClientGetOptions) (azfake.Responder[armcompute.DisksClientGetResponse], azfake.ErrorResponder)
		name    string
		want    bool
		wantErr bool
	}{
		{
			name: "disk exists",
			want: true,
			get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
				resp.SetResponse(http.StatusOK, armcompute.DisksClientGetResponse{
					Disk: armcompute.Disk{Name: to.Ptr("d")},
				}, nil)

				return
			},
		},
		{
			name: "not found",
			want: false,
			get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
				errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")
				return
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := &computefake.DisksServer{Get: tc.get}
			client := newFakeDisksClient(t, srv)

			exists, err := azure.DiskExistsWithClient(t.Context(), client, "rg", "d")
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.want, exists)
		})
	}
}
