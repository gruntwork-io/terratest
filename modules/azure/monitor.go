package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
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
// TODO: this should be removed in the next version
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
// TODO: this should be removed in the next version
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


@@ -46,7 +46,7 @@ func GetDiagnosticsSettingsResourceE(name string, resourceURI string, subscripti
		return nil, err			return nil, err
	}		}


	client, err := GetDiagnosticsSettingsClientE(subscriptionID)		client, err := CreateDiagnosticsSettingsClientE(subscriptionID)
	if err != nil {		if err != nil {
		return nil, err			return nil, err
	}		}
@@ -59,25 +59,6 @@ func GetDiagnosticsSettingsResourceE(name string, resourceURI string, subscripti
	return &settings, nil		return &settings, nil
}	}


// GetDiagnosticsSettingsClientE returns a diagnostics settings client	
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

// GetVMInsightsOnboardingStatus get diagnostics VM onboarding status	// GetVMInsightsOnboardingStatus get diagnostics VM onboarding status
// This function would fail the test if there is an error.	// This function would fail the test if there is an error.
func GetVMInsightsOnboardingStatus(t testing.TestingT, resourceURI string, subscriptionID string) *insights.VMInsightsOnboardingStatus {	func GetVMInsightsOnboardingStatus(t testing.TestingT, resourceURI string, subscriptionID string) *insights.VMInsightsOnboardingStatus {
@@ -89,7 +70,7 @@ func GetVMInsightsOnboardingStatus(t testing.TestingT, resourceURI string, subsc


// GetVMInsightsOnboardingStatusE get diagnostics VM onboarding status	// GetVMInsightsOnboardingStatusE get diagnostics VM onboarding status
func GetVMInsightsOnboardingStatusE(t testing.TestingT, resourceURI string, subscriptionID string) (*insights.VMInsightsOnboardingStatus, error) {	func GetVMInsightsOnboardingStatusE(t testing.TestingT, resourceURI string, subscriptionID string) (*insights.VMInsightsOnboardingStatus, error) {
	client, err := GetVMInsightsClientE(t, subscriptionID)		client, err := CreateVMInsightsClientE(subscriptionID)
	if err != nil {		if err != nil {
		return nil, err			return nil, err
	}		}
@@ -102,26 +83,6 @@ func GetVMInsightsOnboardingStatusE(t testing.TestingT, resourceURI string, subs
	return &status, nil		return &status, nil
}	}


// GetVMInsightsClientE gets a VM Insights client	
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

// GetActivityLogAlertResource gets a Action Group in the specified Azure Resource Group	// GetActivityLogAlertResource gets a Action Group in the specified Azure Resource Group
// This function would fail the test if there is an error.	// This function would fail the test if there is an error.
func GetActivityLogAlertResource(t testing.TestingT, activityLogAlertName string, resGroupName string, subscriptionID string) *insights.ActivityLogAlertResource {	func GetActivityLogAlertResource(t testing.TestingT, activityLogAlertName string, resGroupName string, subscriptionID string) *insights.ActivityLogAlertResource {
@@ -140,7 +101,7 @@ func GetActivityLogAlertResourceE(activityLogAlertName string, resGroupName stri
	}		}


	// Get the client reference		// Get the client reference
	client, err := GetActivityLogAlertsClientE(subscriptionID)		client, err := CreateActivityLogAlertsClientE(subscriptionID)
	if err != nil {		if err != nil {
		return nil, err			return nil, err
	}		}
@@ -153,25 +114,3 @@ func GetActivityLogAlertResourceE(activityLogAlertName string, resGroupName stri


	return &activityLogAlertResource, nil		return &activityLogAlertResource, nil
}	}

// GetActivityLogAlertsClientE gets an Action Groups client in the specified Azure Subscription	
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

// GetActivityLogAlertsClientE gets an Action Groups client in the specified Azure Subscription
// TODO: this should be removed in the next version
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
