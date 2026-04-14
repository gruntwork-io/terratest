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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	appfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeManagedEnvironmentsClient(t *testing.T, srv *appfake.ManagedEnvironmentsServer) *armappcontainers.ManagedEnvironmentsClient {
	t.Helper()

	client, err := armappcontainers.NewManagedEnvironmentsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: appfake.NewManagedEnvironmentsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeContainerAppsClient(t *testing.T, srv *appfake.ContainerAppsServer) *armappcontainers.ContainerAppsClient {
	t.Helper()

	client, err := armappcontainers.NewContainerAppsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: appfake.NewContainerAppsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeJobsClient(t *testing.T, srv *appfake.JobsServer) *armappcontainers.JobsClient {
	t.Helper()

	client, err := armappcontainers.NewJobsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: appfake.NewJobsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetManagedEnvironmentWithClient tests
// ---------------------------------------------------------------------------

func TestGetManagedEnvironmentWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    appfake.ManagedEnvironmentsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: appfake.ManagedEnvironmentsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappcontainers.ManagedEnvironmentsClientGetOptions) (resp azfake.Responder[armappcontainers.ManagedEnvironmentsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armappcontainers.ManagedEnvironmentsClientGetResponse{
						ManagedEnvironment: armappcontainers.ManagedEnvironment{
							Name: to.Ptr("test-env"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-env",
		},
		{
			name: "NotFound",
			server: appfake.ManagedEnvironmentsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappcontainers.ManagedEnvironmentsClientGetOptions) (resp azfake.Responder[armappcontainers.ManagedEnvironmentsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeManagedEnvironmentsClient(t, &tc.server)

			env, err := azure.GetManagedEnvironmentWithClient(context.Background(), client, "rg", "test-env")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *env.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// GetContainerAppWithClient tests
// ---------------------------------------------------------------------------

func TestGetContainerAppWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    appfake.ContainerAppsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: appfake.ContainerAppsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappcontainers.ContainerAppsClientGetOptions) (resp azfake.Responder[armappcontainers.ContainerAppsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armappcontainers.ContainerAppsClientGetResponse{
						ContainerApp: armappcontainers.ContainerApp{
							Name: to.Ptr("test-app"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-app",
		},
		{
			name: "NotFound",
			server: appfake.ContainerAppsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappcontainers.ContainerAppsClientGetOptions) (resp azfake.Responder[armappcontainers.ContainerAppsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeContainerAppsClient(t, &tc.server)

			app, err := azure.GetContainerAppWithClient(context.Background(), client, "rg", "test-app")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *app.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// GetContainerAppJobWithClient tests
// ---------------------------------------------------------------------------

func TestGetContainerAppJobWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    appfake.JobsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: appfake.JobsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappcontainers.JobsClientGetOptions) (resp azfake.Responder[armappcontainers.JobsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armappcontainers.JobsClientGetResponse{
						Job: armappcontainers.Job{
							Name: to.Ptr("test-job"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-job",
		},
		{
			name: "NotFound",
			server: appfake.JobsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappcontainers.JobsClientGetOptions) (resp azfake.Responder[armappcontainers.JobsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeJobsClient(t, &tc.server)

			job, err := azure.GetContainerAppJobWithClient(context.Background(), client, "rg", "test-job")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *job.Name)
		})
	}
}
