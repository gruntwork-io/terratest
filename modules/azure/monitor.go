package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DiagnosticSettingsResourceExistsContext indicates whether the diagnostic settings resource exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DiagnosticSettingsResourceExistsContext(t testing.TestingT, ctx context.Context, diagnosticSettingsResourceName string, resourceURI string, subscriptionID string) bool {
	t.Helper()

	exists, err := DiagnosticSettingsResourceExistsContextE(ctx, diagnosticSettingsResourceName, resourceURI, subscriptionID)
	require.NoError(t, err)

	return exists
}

// DiagnosticSettingsResourceExists indicates whether the diagnostic settings resource exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [DiagnosticSettingsResourceExistsContext] instead.
func DiagnosticSettingsResourceExists(t testing.TestingT, diagnosticSettingsResourceName string, resourceURI string, subscriptionID string) bool {
	t.Helper()

	return DiagnosticSettingsResourceExistsContext(t, context.Background(), diagnosticSettingsResourceName, resourceURI, subscriptionID)
}

// DiagnosticSettingsResourceExistsContextE indicates whether the diagnostic settings resource exists.
// The ctx parameter supports cancellation and timeouts.
func DiagnosticSettingsResourceExistsContextE(ctx context.Context, diagnosticSettingsResourceName string, resourceURI string, subscriptionID string) (bool, error) {
	_, err := GetDiagnosticsSettingsResourceContextE(ctx, diagnosticSettingsResourceName, resourceURI, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// DiagnosticSettingsResourceExistsE indicates whether the diagnostic settings resource exists.
//
// Deprecated: Use [DiagnosticSettingsResourceExistsContextE] instead.
func DiagnosticSettingsResourceExistsE(diagnosticSettingsResourceName string, resourceURI string, subscriptionID string) (bool, error) {
	return DiagnosticSettingsResourceExistsContextE(context.Background(), diagnosticSettingsResourceName, resourceURI, subscriptionID)
}

// DiagnosticSettingsResourceExistsWithClient checks if a diagnostic settings resource exists using the provided client.
func DiagnosticSettingsResourceExistsWithClient(ctx context.Context, client *armmonitor.DiagnosticSettingsClient, resourceURI string, name string) (bool, error) {
	_, err := GetDiagnosticsSettingsResourceWithClient(ctx, client, resourceURI, name)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetDiagnosticsSettingsResourceContext gets the diagnostics settings for a specified resource.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiagnosticsSettingsResourceContext(t testing.TestingT, ctx context.Context, name string, resourceURI string, subscriptionID string) *armmonitor.DiagnosticSettingsResource {
	t.Helper()

	resource, err := GetDiagnosticsSettingsResourceContextE(ctx, name, resourceURI, subscriptionID)
	require.NoError(t, err)

	return resource
}

// GetDiagnosticsSettingsResource gets the diagnostics settings for a specified resource.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetDiagnosticsSettingsResourceContext] instead.
func GetDiagnosticsSettingsResource(t testing.TestingT, name string, resourceURI string, subscriptionID string) *armmonitor.DiagnosticSettingsResource {
	t.Helper()

	return GetDiagnosticsSettingsResourceContext(t, context.Background(), name, resourceURI, subscriptionID)
}

// GetDiagnosticsSettingsResourceContextE gets the diagnostics settings for a specified resource.
// The ctx parameter supports cancellation and timeouts.
func GetDiagnosticsSettingsResourceContextE(ctx context.Context, name string, resourceURI string, subscriptionID string) (*armmonitor.DiagnosticSettingsResource, error) {
	// Validate Azure subscription ID
	_, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client, err := CreateDiagnosticsSettingsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetDiagnosticsSettingsResourceWithClient(ctx, client, resourceURI, name)
}

// GetDiagnosticsSettingsResourceE gets the diagnostics settings for a specified resource.
//
// Deprecated: Use [GetDiagnosticsSettingsResourceContextE] instead.
func GetDiagnosticsSettingsResourceE(name string, resourceURI string, subscriptionID string) (*armmonitor.DiagnosticSettingsResource, error) {
	return GetDiagnosticsSettingsResourceContextE(context.Background(), name, resourceURI, subscriptionID)
}

// GetDiagnosticsSettingsResourceWithClient gets the diagnostics settings for a specified resource using the provided client.
func GetDiagnosticsSettingsResourceWithClient(ctx context.Context, client *armmonitor.DiagnosticSettingsClient, resourceURI string, name string) (*armmonitor.DiagnosticSettingsResource, error) {
	resp, err := client.Get(ctx, resourceURI, name, nil)
	if err != nil {
		return nil, err
	}

	return &resp.DiagnosticSettingsResource, nil
}

// GetVMInsightsOnboardingStatusContext gets diagnostics VM onboarding status.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMInsightsOnboardingStatusContext(t testing.TestingT, ctx context.Context, resourceURI string, subscriptionID string) *armmonitor.VMInsightsOnboardingStatus {
	t.Helper()

	status, err := GetVMInsightsOnboardingStatusContextE(t, ctx, resourceURI, subscriptionID)
	require.NoError(t, err)

	return status
}

// GetVMInsightsOnboardingStatus gets diagnostics VM onboarding status.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVMInsightsOnboardingStatusContext] instead.
func GetVMInsightsOnboardingStatus(t testing.TestingT, resourceURI string, subscriptionID string) *armmonitor.VMInsightsOnboardingStatus {
	t.Helper()

	return GetVMInsightsOnboardingStatusContext(t, context.Background(), resourceURI, subscriptionID)
}

// GetVMInsightsOnboardingStatusContextE gets diagnostics VM onboarding status.
// The ctx parameter supports cancellation and timeouts.
func GetVMInsightsOnboardingStatusContextE(t testing.TestingT, ctx context.Context, resourceURI string, subscriptionID string) (*armmonitor.VMInsightsOnboardingStatus, error) {
	client, err := CreateVMInsightsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetVMInsightsOnboardingStatusWithClient(ctx, client, resourceURI)
}

// GetVMInsightsOnboardingStatusE gets diagnostics VM onboarding status.
//
// Deprecated: Use [GetVMInsightsOnboardingStatusContextE] instead.
func GetVMInsightsOnboardingStatusE(t testing.TestingT, resourceURI string, subscriptionID string) (*armmonitor.VMInsightsOnboardingStatus, error) {
	return GetVMInsightsOnboardingStatusContextE(t, context.Background(), resourceURI, subscriptionID)
}

// GetVMInsightsOnboardingStatusWithClient gets diagnostics VM onboarding status using the provided client.
func GetVMInsightsOnboardingStatusWithClient(ctx context.Context, client *armmonitor.VMInsightsClient, resourceURI string) (*armmonitor.VMInsightsOnboardingStatus, error) {
	resp, err := client.GetOnboardingStatus(ctx, resourceURI, nil)
	if err != nil {
		return nil, err
	}

	return &resp.VMInsightsOnboardingStatus, nil
}

// GetActivityLogAlertResourceContext gets an Activity Log Alert Resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetActivityLogAlertResourceContext(t testing.TestingT, ctx context.Context, activityLogAlertName string, resGroupName string, subscriptionID string) *armmonitor.ActivityLogAlertResource {
	t.Helper()

	activityLogAlertResource, err := GetActivityLogAlertResourceContextE(ctx, activityLogAlertName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return activityLogAlertResource
}

// GetActivityLogAlertResource gets an Activity Log Alert Resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetActivityLogAlertResourceContext] instead.
func GetActivityLogAlertResource(t testing.TestingT, activityLogAlertName string, resGroupName string, subscriptionID string) *armmonitor.ActivityLogAlertResource {
	t.Helper()

	return GetActivityLogAlertResourceContext(t, context.Background(), activityLogAlertName, resGroupName, subscriptionID)
}

// GetActivityLogAlertResourceContextE gets an Activity Log Alert Resource in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetActivityLogAlertResourceContextE(ctx context.Context, activityLogAlertName string, resGroupName string, subscriptionID string) (*armmonitor.ActivityLogAlertResource, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := CreateActivityLogAlertsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetActivityLogAlertResourceWithClient(ctx, client, resGroupName, activityLogAlertName)
}

// GetActivityLogAlertResourceE gets an Activity Log Alert Resource in the specified Azure Resource Group.
//
// Deprecated: Use [GetActivityLogAlertResourceContextE] instead.
func GetActivityLogAlertResourceE(activityLogAlertName string, resGroupName string, subscriptionID string) (*armmonitor.ActivityLogAlertResource, error) {
	return GetActivityLogAlertResourceContextE(context.Background(), activityLogAlertName, resGroupName, subscriptionID)
}

// GetActivityLogAlertResourceWithClient gets an Activity Log Alert Resource using the provided client.
func GetActivityLogAlertResourceWithClient(ctx context.Context, client *armmonitor.ActivityLogAlertsClient, resGroupName string, activityLogAlertName string) (*armmonitor.ActivityLogAlertResource, error) {
	resp, err := client.Get(ctx, resGroupName, activityLogAlertName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ActivityLogAlertResource, nil
}
