package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DiagnosticSettingsResourceExists indicates whether the diagnostic settings resource exists
// This function would fail the test if there is an error.
func DiagnosticSettingsResourceExists(t testing.TestingT, diagnosticSettingsResourceName string, resourceURI string, subscriptionID string) bool {
	exists, err := DiagnosticSettingsResourceExistsE(diagnosticSettingsResourceName, resourceURI, subscriptionID)
	require.NoError(t, err)

	return exists
}

// DiagnosticSettingsResourceExistsE indicates whether the diagnostic settings resource exists
func DiagnosticSettingsResourceExistsE(diagnosticSettingsResourceName string, resourceURI string, subscriptionID string) (bool, error) {
	_, err := GetDiagnosticsSettingsResourceE(diagnosticSettingsResourceName, resourceURI, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetDiagnosticsSettingsResource gets the diagnostics settings for a specified resource
// This function would fail the test if there is an error.
func GetDiagnosticsSettingsResource(t testing.TestingT, name string, resourceURI string, subscriptionID string) *insights.DiagnosticSettingsResource {
	resource, err := GetDiagnosticsSettingsResourceE(name, resourceURI, subscriptionID)
	require.NoError(t, err)
	return resource
}

// GetDiagnosticsSettingsResourceE gets the diagnostics settings for a specified resource
func GetDiagnosticsSettingsResourceE(name string, resourceURI string, subscriptionID string) (*insights.DiagnosticSettingsResource, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client, err := CreateDiagnosticsSettingsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	settings, err := client.Get(context.Background(), resourceURI, name)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// GetDiagnosticsSettingsClientE returns a diagnostics settings client
// TODO: delete in next version
func GetDiagnosticsSettingsClientE(subscriptionID string) (*insights.DiagnosticSettingsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := insights.NewDiagnosticSettingsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// GetVMInsightsOnboardingStatus get diagnostics VM onboarding status
// This function would fail the test if there is an error.
func GetVMInsightsOnboardingStatus(t testing.TestingT, resourceURI string, subscriptionID string) *insights.VMInsightsOnboardingStatus {
	status, err := GetVMInsightsOnboardingStatusE(t, resourceURI, subscriptionID)
	require.NoError(t, err)

	return status
}

// GetVMInsightsOnboardingStatusE get diagnostics VM onboarding status
func GetVMInsightsOnboardingStatusE(t testing.TestingT, resourceURI string, subscriptionID string) (*insights.VMInsightsOnboardingStatus, error) {
	client, err := CreateVMInsightsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	status, err := client.GetOnboardingStatus(context.Background(), resourceURI)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// GetVMInsightsClientE gets a VM Insights client
// TODO: delete in next version
func GetVMInsightsClientE(t testing.TestingT, subscriptionID string) (*insights.VMInsightsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := insights.NewVMInsightsClient(subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// GetActivityLogAlertResource gets a Action Group in the specified Azure Resource Group
// This function would fail the test if there is an error.
func GetActivityLogAlertResource(t testing.TestingT, activityLogAlertName string, resGroupName string, subscriptionID string) *insights.ActivityLogAlertResource {
	activityLogAlertResource, err := GetActivityLogAlertResourceE(activityLogAlertName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return activityLogAlertResource
}

// GetActivityLogAlertResourceE gets a Action Group in the specified Azure Resource Group
func GetActivityLogAlertResourceE(activityLogAlertName string, resGroupName string, subscriptionID string) (*insights.ActivityLogAlertResource, error) {
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

	// Get the Action Group
	activityLogAlertResource, err := client.Get(context.Background(), resGroupName, activityLogAlertName)
	if err != nil {
		return nil, err
	}

	return &activityLogAlertResource, nil
}

// GetActivityLogAlertsClientE gets an Action Groups client in the specified Azure Subscription
// TODO: delete in next version
func GetActivityLogAlertsClientE(subscriptionID string) (*insights.ActivityLogAlertsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Action Groups client
	client := insights.NewActivityLogAlertsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}
