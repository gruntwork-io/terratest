package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
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

	return DiagnosticSettingsResourceExistsContext(t, context.Background(), diagnosticSettingsResourceName, resourceURI, subscriptionID) //nolint:staticcheck
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

// GetDiagnosticsSettingsResourceContext gets the diagnostics settings for a specified resource.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDiagnosticsSettingsResourceContext(t testing.TestingT, ctx context.Context, name string, resourceURI string, subscriptionID string) *insights.DiagnosticSettingsResource {
	t.Helper()

	resource, err := GetDiagnosticsSettingsResourceContextE(ctx, name, resourceURI, subscriptionID)
	require.NoError(t, err)

	return resource
}

// GetDiagnosticsSettingsResource gets the diagnostics settings for a specified resource.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetDiagnosticsSettingsResourceContext] instead.
func GetDiagnosticsSettingsResource(t testing.TestingT, name string, resourceURI string, subscriptionID string) *insights.DiagnosticSettingsResource {
	t.Helper()

	return GetDiagnosticsSettingsResourceContext(t, context.Background(), name, resourceURI, subscriptionID) //nolint:staticcheck
}

// GetDiagnosticsSettingsResourceContextE gets the diagnostics settings for a specified resource.
// The ctx parameter supports cancellation and timeouts.
func GetDiagnosticsSettingsResourceContextE(ctx context.Context, name string, resourceURI string, subscriptionID string) (*insights.DiagnosticSettingsResource, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client, err := CreateDiagnosticsSettingsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	settings, err := client.Get(ctx, resourceURI, name)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// GetDiagnosticsSettingsResourceE gets the diagnostics settings for a specified resource.
//
// Deprecated: Use [GetDiagnosticsSettingsResourceContextE] instead.
func GetDiagnosticsSettingsResourceE(name string, resourceURI string, subscriptionID string) (*insights.DiagnosticSettingsResource, error) {
	return GetDiagnosticsSettingsResourceContextE(context.Background(), name, resourceURI, subscriptionID)
}

// GetVMInsightsOnboardingStatusContext gets diagnostics VM onboarding status.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMInsightsOnboardingStatusContext(t testing.TestingT, ctx context.Context, resourceURI string, subscriptionID string) *insights.VMInsightsOnboardingStatus {
	t.Helper()

	status, err := GetVMInsightsOnboardingStatusContextE(t, ctx, resourceURI, subscriptionID)
	require.NoError(t, err)

	return status
}

// GetVMInsightsOnboardingStatus gets diagnostics VM onboarding status.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVMInsightsOnboardingStatusContext] instead.
func GetVMInsightsOnboardingStatus(t testing.TestingT, resourceURI string, subscriptionID string) *insights.VMInsightsOnboardingStatus {
	t.Helper()

	return GetVMInsightsOnboardingStatusContext(t, context.Background(), resourceURI, subscriptionID) //nolint:staticcheck
}

// GetVMInsightsOnboardingStatusContextE gets diagnostics VM onboarding status.
// The ctx parameter supports cancellation and timeouts.
func GetVMInsightsOnboardingStatusContextE(t testing.TestingT, ctx context.Context, resourceURI string, subscriptionID string) (*insights.VMInsightsOnboardingStatus, error) {
	client, err := CreateVMInsightsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	status, err := client.GetOnboardingStatus(ctx, resourceURI)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// GetVMInsightsOnboardingStatusE gets diagnostics VM onboarding status.
//
// Deprecated: Use [GetVMInsightsOnboardingStatusContextE] instead.
func GetVMInsightsOnboardingStatusE(t testing.TestingT, resourceURI string, subscriptionID string) (*insights.VMInsightsOnboardingStatus, error) {
	return GetVMInsightsOnboardingStatusContextE(t, context.Background(), resourceURI, subscriptionID)
}

// GetActivityLogAlertResourceContext gets an Activity Log Alert Resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetActivityLogAlertResourceContext(t testing.TestingT, ctx context.Context, activityLogAlertName string, resGroupName string, subscriptionID string) *insights.ActivityLogAlertResource {
	t.Helper()

	activityLogAlertResource, err := GetActivityLogAlertResourceContextE(ctx, activityLogAlertName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return activityLogAlertResource
}

// GetActivityLogAlertResource gets an Activity Log Alert Resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetActivityLogAlertResourceContext] instead.
func GetActivityLogAlertResource(t testing.TestingT, activityLogAlertName string, resGroupName string, subscriptionID string) *insights.ActivityLogAlertResource {
	t.Helper()

	return GetActivityLogAlertResourceContext(t, context.Background(), activityLogAlertName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetActivityLogAlertResourceContextE gets an Activity Log Alert Resource in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetActivityLogAlertResourceContextE(ctx context.Context, activityLogAlertName string, resGroupName string, subscriptionID string) (*insights.ActivityLogAlertResource, error) {
	// Validate resource group name and subscription ID
	_, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := CreateActivityLogAlertsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Activity Log Alert Resource
	activityLogAlertResource, err := client.Get(ctx, resGroupName, activityLogAlertName)
	if err != nil {
		return nil, err
	}

	return &activityLogAlertResource, nil
}

// GetActivityLogAlertResourceE gets an Activity Log Alert Resource in the specified Azure Resource Group.
//
// Deprecated: Use [GetActivityLogAlertResourceContextE] instead.
func GetActivityLogAlertResourceE(activityLogAlertName string, resGroupName string, subscriptionID string) (*insights.ActivityLogAlertResource, error) {
	return GetActivityLogAlertResourceContextE(context.Background(), activityLogAlertName, resGroupName, subscriptionID)
}

