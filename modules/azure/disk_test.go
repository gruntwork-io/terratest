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

	client, err := armcompute.NewDisksClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: computefake.NewDisksServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// DiskExists tests (via fake client replicating DiskExistsContextE logic)
// ---------------------------------------------------------------------------

func TestDiskExists(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name   string
		want   bool
		server computefake.DisksServer
	}{
		{
			name: "Exists",
			server: computefake.DisksServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcompute.DisksClientGetResponse{
						Disk: armcompute.Disk{Name: to.Ptr("my-disk")},
					}, nil)

					return
				},
			},
			want: true,
		},
		{
			name: "NotFound",
			server: computefake.DisksServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeDisksClient(t, &tc.server)

			_, err := client.Get(context.Background(), "rg", "my-disk", nil)
			if err != nil {
				assert.False(t, tc.want)
				assert.True(t, azure.ResourceNotFoundErrorExists(err))

				return
			}

			assert.True(t, tc.want)
		})
	}
}

// ---------------------------------------------------------------------------
// GetDisk tests (via fake client replicating GetDiskContextE logic)
// ---------------------------------------------------------------------------

func TestGetDisk(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    computefake.DisksServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: computefake.DisksServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcompute.DisksClientGetResponse{
						Disk: armcompute.Disk{
							Name: to.Ptr("os-disk"),
							Properties: &armcompute.DiskProperties{
								DiskSizeGB: to.Ptr[int32](128),
							},
						},
					}, nil)

					return
				},
			},
			wantName: "os-disk",
		},
		{
			name: "NotFound",
			server: computefake.DisksServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcompute.DisksClientGetOptions) (resp azfake.Responder[armcompute.DisksClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			wantErr:   true,
			errSubstr: "ResourceNotFound",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeDisksClient(t, &tc.server)

			resp, err := client.Get(context.Background(), "rg", "disk", nil)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *resp.Disk.Name)
		})
	}
}
