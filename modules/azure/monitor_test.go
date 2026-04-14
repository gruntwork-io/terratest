package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	monitorfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeDiagnosticSettingsClient(t *testing.T, srv *monitorfake.DiagnosticSettingsServer) *armmonitor.DiagnosticSettingsClient {
	t.Helper()

	transport := monitorfake.NewDiagnosticSettingsServerTransport(srv)
	client, err := armmonitor.NewDiagnosticSettingsClient(&azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakeVMInsightsClient(t *testing.T, srv *monitorfake.VMInsightsServer) *armmonitor.VMInsightsClient {
	t.Helper()

	transport := monitorfake.NewVMInsightsServerTransport(srv)
	client, err := armmonitor.NewVMInsightsClient(&azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakeActivityLogAlertsClient(t *testing.T, srv *monitorfake.ActivityLogAlertsServer) *armmonitor.ActivityLogAlertsClient {
	t.Helper()

	transport := monitorfake.NewActivityLogAlertsServerTransport(srv)
	client, err := armmonitor.NewActivityLogAlertsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetDiagnosticsSettingsResourceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  monitorfake.DiagnosticSettingsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: monitorfake.DiagnosticSettingsServer{
				Get: func(_ context.Context, resourceURI string, name string, _ *armmonitor.DiagnosticSettingsClientGetOptions) (resp azfake.Responder[armmonitor.DiagnosticSettingsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armmonitor.DiagnosticSettingsClientGetResponse{
						DiagnosticSettingsResource: armmonitor.DiagnosticSettingsResource{
							Name: to.Ptr(name),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: monitorfake.DiagnosticSettingsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armmonitor.DiagnosticSettingsClientGetOptions) (resp azfake.Responder[armmonitor.DiagnosticSettingsClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeDiagnosticSettingsClient(t, &srv)

			resource, err := azure.GetDiagnosticsSettingsResourceWithClient(t.Context(), client, "/subscriptions/fake/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm", "my-setting")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-setting", *resource.Name)
		})
	}
}

func TestDiagnosticSettingsResourceExistsWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		get  func(context.Context, string, string, *armmonitor.DiagnosticSettingsClientGetOptions) (azfake.Responder[armmonitor.DiagnosticSettingsClientGetResponse], azfake.ErrorResponder)
		name string
		want bool
	}{
		{
			name: "exists",
			want: true,
			get: func(_ context.Context, _ string, name string, _ *armmonitor.DiagnosticSettingsClientGetOptions) (resp azfake.Responder[armmonitor.DiagnosticSettingsClientGetResponse], errResp azfake.ErrorResponder) {
				resp.SetResponse(http.StatusOK, armmonitor.DiagnosticSettingsClientGetResponse{
					DiagnosticSettingsResource: armmonitor.DiagnosticSettingsResource{
						Name: to.Ptr(name),
					},
				}, nil)

				return
			},
		},
		{
			name: "not found",
			want: false,
			get: func(_ context.Context, _ string, _ string, _ *armmonitor.DiagnosticSettingsClientGetOptions) (resp azfake.Responder[armmonitor.DiagnosticSettingsClientGetResponse], errResp azfake.ErrorResponder) {
				errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")
				return
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := &monitorfake.DiagnosticSettingsServer{Get: tc.get}
			client := newFakeDiagnosticSettingsClient(t, srv)

			exists, err := azure.DiagnosticSettingsResourceExistsWithClient(t.Context(), client, "/subscriptions/fake/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm", "my-setting")
			require.NoError(t, err)
			assert.Equal(t, tc.want, exists)
		})
	}
}

func TestGetVMInsightsOnboardingStatusWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  monitorfake.VMInsightsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: monitorfake.VMInsightsServer{
				GetOnboardingStatus: func(_ context.Context, resourceURI string, _ *armmonitor.VMInsightsClientGetOnboardingStatusOptions) (resp azfake.Responder[armmonitor.VMInsightsClientGetOnboardingStatusResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armmonitor.VMInsightsClientGetOnboardingStatusResponse{
						VMInsightsOnboardingStatus: armmonitor.VMInsightsOnboardingStatus{
							Name: to.Ptr("default"),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: monitorfake.VMInsightsServer{
				GetOnboardingStatus: func(_ context.Context, _ string, _ *armmonitor.VMInsightsClientGetOnboardingStatusOptions) (resp azfake.Responder[armmonitor.VMInsightsClientGetOnboardingStatusResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeVMInsightsClient(t, &srv)

			status, err := azure.GetVMInsightsOnboardingStatusWithClient(t.Context(), client, "/subscriptions/fake/resourceGroups/rg/providers/Microsoft.Compute/virtualMachines/vm")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "default", *status.Name)
		})
	}
}

func TestGetActivityLogAlertResourceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  monitorfake.ActivityLogAlertsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: monitorfake.ActivityLogAlertsServer{
				Get: func(_ context.Context, _ string, alertName string, _ *armmonitor.ActivityLogAlertsClientGetOptions) (resp azfake.Responder[armmonitor.ActivityLogAlertsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armmonitor.ActivityLogAlertsClientGetResponse{
						ActivityLogAlertResource: armmonitor.ActivityLogAlertResource{
							Name: to.Ptr(alertName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: monitorfake.ActivityLogAlertsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armmonitor.ActivityLogAlertsClientGetOptions) (resp azfake.Responder[armmonitor.ActivityLogAlertsClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeActivityLogAlertsClient(t, &srv)

			resource, err := azure.GetActivityLogAlertResourceWithClient(t.Context(), client, "rg", "my-alert")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-alert", *resource.Name)
		})
	}
}
