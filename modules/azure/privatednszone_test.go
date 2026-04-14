package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	privatednsfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakePrivateZonesClient(t *testing.T, srv *privatednsfake.PrivateZonesServer) *armprivatedns.PrivateZonesClient {
	t.Helper()

	client, err := armprivatedns.NewPrivateZonesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: privatednsfake.NewPrivateZonesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func TestGetPrivateDNSZoneWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    privatednsfake.PrivateZonesServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: privatednsfake.PrivateZonesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armprivatedns.PrivateZonesClientGetOptions) (resp azfake.Responder[armprivatedns.PrivateZonesClientGetResponse], errResp azfake.ErrorResponder) {
					result := armprivatedns.PrivateZonesClientGetResponse{
						PrivateZone: armprivatedns.PrivateZone{
							Name: to.Ptr("test.private.zone"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test.private.zone",
		},
		{
			name: "NotFound",
			server: privatednsfake.PrivateZonesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armprivatedns.PrivateZonesClientGetOptions) (resp azfake.Responder[armprivatedns.PrivateZonesClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakePrivateZonesClient(t, &tc.server)

			zone, err := azure.GetPrivateDNSZoneWithClient(context.Background(), client, "rg", "test.private.zone")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *zone.Name)
		})
	}
}
