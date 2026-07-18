package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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

// GetDiagnosticsSettingsResourceContextE gets the diagnostics settings for a specified resource.
// The ctx parameter supports cancellation and timeouts.
func GetDiagnosticsSettingsResourceContextE(ctx context.Context, name string, resourceURI string, subscriptionID string) (*armmonitor.DiagnosticSettingsResource, error) {

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

// GetVMInsightsOnboardingStatusContextE gets diagnostics VM onboarding status.
// The ctx parameter supports cancellation and timeouts.
func GetVMInsightsOnboardingStatusContextE(t testing.TestingT, ctx context.Context, resourceURI string, subscriptionID string) (*armmonitor.VMInsightsOnboardingStatus, error) {
	client, err := CreateVMInsightsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetVMInsightsOnboardingStatusWithClient(ctx, client, resourceURI)
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

// GetActivityLogAlertResourceContextE gets an Activity Log Alert Resource in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetActivityLogAlertResourceContextE(ctx context.Context, activityLogAlertName string, resGroupName string, subscriptionID string) (*armmonitor.ActivityLogAlertResource, error) {

	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateActivityLogAlertsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetActivityLogAlertResourceWithClient(ctx, client, resGroupName, activityLogAlertName)
}

// GetActivityLogAlertResourceWithClient gets an Activity Log Alert Resource using the provided client.
func GetActivityLogAlertResourceWithClient(ctx context.Context, client *armmonitor.ActivityLogAlertsClient, resGroupName string, activityLogAlertName string) (*armmonitor.ActivityLogAlertResource, error) {
	resp, err := client.Get(ctx, resGroupName, activityLogAlertName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ActivityLogAlertResource, nil
}
