package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/automation/mgmt/automation"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// AutomationAccountExists indicates whether the specified Azure Automation Account exists.
// This function would fail the test if there is an error.
func AutomationAccountExists(t testing.TestingT, automationAccountName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := AutomationAccountExistsE(t, automationAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AutomationAccountExistsE indicates whether the specified Azure Automation Account exists.
func AutomationAccountExistsE(t testing.TestingT, automationAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAutomationAccountE(t, automationAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// AutomationAccountDscExists indicates whether the specified Azure Automation Account DSC exists.
// This function would fail the test if there is an error.
func AutomationAccountDscExists(t testing.TestingT, dscConfiguraitonName string, automationAccountName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := AutomationAccountDscExistsE(t, dscConfiguraitonName, automationAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AutomationAccountDscExistsE indicates whether the specified Azure Automation Account exists.
func AutomationAccountDscExistsE(t testing.TestingT, dscConfiguraitonName string, automationAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAutomationAccountDscConfigurationE(t, dscConfiguraitonName, automationAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// WaitUntilDscCompiled indicates whether the specified Azure Automation Account DSC compiled successfully.
// This function would fail the test if there is an error.
func WaitUntilDscCompiled(t testing.TestingT, dscConfiguraitonName string, automationAccountName string, resourceGroupName string, subscriptionID string) bool {
	compiled, err := WaitUntilDscCompiledE(t, dscConfiguraitonName, automationAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return compiled
}

// WaitUntilDscCompiledE indicates whether the specified Azure Automation Account DSC compiled successfully.
// DSC compilation is performed via a Terraform null_resource using PowerShell Core, executing async.
// Compilation can take a few minutes to spin up resources requiring a retry mechanism to allow for this
func WaitUntilDscCompiledE(t testing.TestingT, dscConfiguraitonName string, automationAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	seconds := 30 // 30 second initial delay
	sleep := time.Duration(seconds) * time.Second
	maxAttempts := 5 // try 5 times, doubling the delay each time.

	dscCompileJobStatus, err := retry.DoWithRetryE(t, "Wait for Azure Automation Account DSC to Compile", maxAttempts, sleep, func() (string, error) {
		return AutomationAccountDscCompileJobStatusE(t, dscConfiguraitonName, automationAccountName, resourceGroupName, subscriptionID)
	})

	if err != nil {
		return false, err
	}

	if dscCompileJobStatus == "Completed" {
		return true, nil
	}

	return false, nil // Return false if the DSC compilation job completed but failed to compile such as from an error in the DSC code
}

// AutomationAccountRunAsCertificateThumbprintMatches indicates whether the specified Azure Automation Account RunAs Certificate exists.
// This function would fail the test if there is an error.
func AutomationAccountRunAsCertificateThumbprintMatches(t testing.TestingT, runAsCertificateThumbprint string, runAsCertificateName string, automationAccountName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := AutomationAccountRunAsCertificateThumbprintMatchesE(t, runAsCertificateThumbprint, runAsCertificateName, automationAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AutomationAccountRunAsCertificateThumbprintMatchesE indicates whether the specified Azure Automation Account RunAs Certificate exists.
func AutomationAccountRunAsCertificateThumbprintMatchesE(t testing.TestingT, runAsCertificateThumbprint string, runAsCertificateName string, automationAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	certificate, err := GetAutomationAccountCertificateE(t, runAsCertificateName, automationAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return *certificate.CertificateProperties.Thumbprint == runAsCertificateThumbprint, nil
}

// AutomationAccountRunAsConnectionExists indicates whether the specified Azure Automation Account RunAs Account exists.
// This function would fail the test if there is an error.
func AutomationAccountRunAsConnectionExists(t testing.TestingT, automationAccountrunAsAccountName string, runAsConnectionType string, runAsCertificateThumbprint string, automationAccountName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := AutomationAccountRunAsConnectionExistsE(t, automationAccountrunAsAccountName, runAsConnectionType, runAsCertificateThumbprint, automationAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AutomationAccountRunAsConnectionExistsE indicates whether the specified Azure Automation Account RunAs Account exists.
func AutomationAccountRunAsConnectionExistsE(t testing.TestingT, automationAccountrunAsAccountName string, runAsConnectionType string, runAsCertificateThumbprint string, automationAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	runAsAccountConnection, err := GetAutomationAccountRunAsConnectionE(t, automationAccountrunAsAccountName, automationAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return *runAsAccountConnection.Name == automationAccountrunAsAccountName &&
		*runAsAccountConnection.ConnectionProperties.ConnectionType.Name == runAsConnectionType &&
		*runAsAccountConnection.ConnectionProperties.FieldDefinitionValues["CertificateThumbprint"] == runAsCertificateThumbprint, nil
}

// AutomationAccountDscAppliedSuccessfullyToVM indicates whether the specified Azure Automation Account compiled DSC has successfully been applied to the target VM.
// This function would fail the test if there is an error.
func AutomationAccountDscAppliedSuccessfullyToVM(t testing.TestingT, dscConfiguraitonName string, vmName string, automationAccountName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := AutomationAccountDscAppliedSuccessfullyToVME(t, dscConfiguraitonName, vmName, automationAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AutomationAccountDscAppliedSuccessfullyToVME indicates whether the specified Azure Automation Account compiled DSC has successfully been applied to the target VM.
func AutomationAccountDscAppliedSuccessfullyToVME(t testing.TestingT, dscConfiguraitonName string, vmName string, automationAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	dscNodeConfig, err := GetAutomationAccountDscNodeConfigurationE(t, dscConfiguraitonName, vmName, automationAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return *dscNodeConfig.Status == "Compliant", nil
}

// GetAutomationAccount returns the Azure Automation Account by name if it exists in the subscription
func GetAutomationAccount(t testing.TestingT, automationAccountName string, resourceGroupName string, subscriptionID string) *automation.Account {
	// Validate resource group name and subscription ID
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	require.NoError(t, err)

	client, err := CreateAutomationAccountClientE(subscriptionID)
	require.NoError(t, err)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	require.NoError(t, err)
	client.Authorizer = *authorizer

	automationAccount, err := client.Get(context.Background(), resourceGroupName, automationAccountName)
	require.NoError(t, err)

	return &automationAccount
}

// GetAutomationAccountE returns the Azure Automation Account by name if it exists in the subscription
func GetAutomationAccountE(t testing.TestingT, automationAccountName string, resourceGroupName string, subscriptionID string) (*automation.Account, error) {
	// Validate resource group name and subscription ID
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateAutomationAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	automationAccount, err := client.Get(context.Background(), resourceGroupName, automationAccountName)
	if err != nil {
		return nil, err
	}

	return &automationAccount, nil
}

// GetAutomationAccountDscConfiguration returns the Azure Automation Account DscConfiguration by Automation Account name if it exists in the subscription
func GetAutomationAccountDscConfiguration(t testing.TestingT, dscConfigurationName string, automationAccountName string, resourceGroupName string, subscriptionID string) *automation.DscConfiguration {
	client, err := CreateAutomationAccountDscConfigurationClientE(subscriptionID)
	require.NoError(t, err)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	require.NoError(t, err)
	client.Authorizer = *authorizer

	dscConfiguration, err := client.Get(context.Background(), resourceGroupName, automationAccountName, dscConfigurationName)
	require.NoError(t, err)

	return &dscConfiguration
}

// GetAutomationAccountDscConfigurationE returns the Azure Automation Account DscConfiguration by Automation Account name if it exists in the subscription
func GetAutomationAccountDscConfigurationE(t testing.TestingT, dscConfigurationName string, automationAccountName string, resourceGroupName string, subscriptionID string) (*automation.DscConfiguration, error) {
	client, err := CreateAutomationAccountDscConfigurationClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	dscConfiguration, err := client.Get(context.Background(), resourceGroupName, automationAccountName, dscConfigurationName)
	if err != nil {
		return nil, err
	}

	return &dscConfiguration, nil
}

// AutomationAccountDscCompileJobStatus returns the Azure Automation Account DscConfiguration by Automation Account name if it exists in the subscription
func AutomationAccountDscCompileJobStatus(t testing.TestingT, dscConfigurationName string, automationAccountName string, resourceGroupName string, subscriptionID string) string {
	client, err := CreateAutomationAccountDscCompilationJobClientE(subscriptionID)
	require.NoError(t, err)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	require.NoError(t, err)
	client.Authorizer = *authorizer

	filter := fmt.Sprintf("properties/configuration/name eq '%s'", dscConfigurationName)
	dscCompilationJobListResultPage, err := client.ListByAutomationAccount(context.Background(), resourceGroupName, automationAccountName, filter)
	require.NoError(t, err)

	var dscCompilationJobs []automation.DscCompilationJob
	var mostRecentCompileJobTick int64
	var mostRecentCompileJobStatus string
	// Loop through filtered pages of DSC compilation jobs to find latest compilation job
	for dscCompilationJobListResultPage.NotDone() {
		dscCompilationJobs = dscCompilationJobListResultPage.Values()
		// Loop through compilation jobs in the current page
		for _, element := range dscCompilationJobs {
			if element.CreationTime.Unix() > mostRecentCompileJobTick {
				mostRecentCompileJobTick = element.CreationTime.Unix()
				mostRecentCompileJobStatus = (string)(element.Status)
			}
		}
		err := dscCompilationJobListResultPage.Next()
		require.NoError(t, err)
	}
	// Check to ensure  DSC compilation jobs are present (i.e. mostRecentCompileJobTick is non zero)
	if mostRecentCompileJobTick == 0 {
		panic("No compilation jobs present for this DSC configuraiton, or compilation jobs are 'suspeneded' in the Automation Account.")
	} else {
		return mostRecentCompileJobStatus
	}
}

// AutomationAccountDscCompileJobStatusE returns the Azure Automation Account DscConfiguration by Automation Account name if it exists in the subscription
func AutomationAccountDscCompileJobStatusE(t testing.TestingT, dscConfigurationName string, automationAccountName string, resourceGroupName string, subscriptionID string) (string, error) {
	client, err := CreateAutomationAccountDscCompilationJobClientE(subscriptionID)
	if err != nil {
		return "", err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return "", err
	}
	client.Authorizer = *authorizer

	filter := fmt.Sprintf("properties/configuration/name eq '%s'", dscConfigurationName)
	dscCompilationJobListResultPage, err := client.ListByAutomationAccount(context.Background(), resourceGroupName, automationAccountName, filter)
	if err != nil {
		return "", err
	}

	var dscCompilationJobs []automation.DscCompilationJob
	var mostRecentCompileJobTick int64
	var mostRecentCompileJobStatus string
	// Loop through filtered pages of DSC compilation jobs to find latest compilation job
	for dscCompilationJobListResultPage.NotDone() {
		dscCompilationJobs = dscCompilationJobListResultPage.Values()
		// Loop through compilation jobs in the current page
		for _, element := range dscCompilationJobs {
			if element.CreationTime.Unix() > mostRecentCompileJobTick {
				mostRecentCompileJobTick = element.CreationTime.Unix()
				mostRecentCompileJobStatus = (string)(element.Status)
			}
		}
		err := dscCompilationJobListResultPage.Next()
		if err != nil {
			return "", err
		}
	}
	// Check to ensure  DSC compilation jobs are present (i.e. mostRecentCompileJobTick is non zero)
	if mostRecentCompileJobTick == 0 {
		panic("No compilation jobs present for this DSC configuraiton, or compilation jobs are 'suspeneded' in the Automation Account.")
	} else {
		return mostRecentCompileJobStatus, nil
	}
}

// GetAutomationAccountCertificate gets the RunAs Connection Certificate if it exists in the Azure Automation Account
func GetAutomationAccountCertificate(t testing.TestingT, automationAccountCertificateName string, automationAccountName string, resourceGroupName string, subscriptionID string) *automation.Certificate {
	client, err := CreateAutomationAccountCertficateClientE(subscriptionID)
	require.NoError(t, err)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	require.NoError(t, err)
	client.Authorizer = *authorizer

	// Get Automation Account Connection
	automationAccountCertificate, err := client.Get(context.Background(), resourceGroupName, automationAccountName, automationAccountCertificateName)
	require.NoError(t, err)

	return &automationAccountCertificate
}

// GetAutomationAccountCertificateE gets the RunAs Connection Certificate if it exists in the Azure Automation Account
func GetAutomationAccountCertificateE(t testing.TestingT, automationAccountCertificateName string, automationAccountName string, resourceGroupName string, subscriptionID string) (*automation.Certificate, error) {
	client, err := CreateAutomationAccountCertficateClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	// Get Automation Account Connection
	automationAccountCertificate, err := client.Get(context.Background(), resourceGroupName, automationAccountName, automationAccountCertificateName)
	if err != nil {
		return nil, err
	}

	return &automationAccountCertificate, nil
}

// GetAutomationAccountDscNodeConfiguration gets the Node Configuration if it exists in the Azure Automation Account
func GetAutomationAccountDscNodeConfiguration(t testing.TestingT, dscConfiguraitonName string, vmName string, automationAccountName string, resourceGroupName string, subscriptionID string) *automation.DscNode {
	client, err := CreateAutomationAccountDscNodeConfigClientE(subscriptionID)
	require.NoError(t, err)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	require.NoError(t, err)
	client.Authorizer = *authorizer

	filter := fmt.Sprintf("name eq '%s'", vmName)
	dscNodeListResultPage, err := client.ListByAutomationAccount(context.Background(), resourceGroupName, automationAccountName, filter)
	require.NoError(t, err)

	var dscNodeList []automation.DscNode
	var dscNodeID string
	for dscNodeListResultPage.NotDone() {
		dscNodeList = dscNodeListResultPage.Values()
		// Loop through compilation jobs in the current page
		for _, element := range dscNodeList {
			if *element.Name == vmName && *element.NodeConfiguration.Name == dscConfiguraitonName {
				dscNodeID = *element.NodeID
			}
		}

		err := dscNodeListResultPage.Next()
		require.NoError(t, err)
	}

	// Get Automation Account Connection
	dscNodeConfig, err := client.Get(context.Background(), resourceGroupName, automationAccountName, dscNodeID)
	require.NoError(t, err)

	return &dscNodeConfig
}

// GetAutomationAccountDscNodeConfigurationE gets the Node Configuration if it exists in the Azure Automation Account
func GetAutomationAccountDscNodeConfigurationE(t testing.TestingT, dscConfiguraitonName string, vmName string, automationAccountName string, resourceGroupName string, subscriptionID string) (*automation.DscNode, error) {
	client, err := CreateAutomationAccountDscNodeConfigClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	filter := fmt.Sprintf("name eq '%s'", vmName)
	dscNodeListResultPage, err := client.ListByAutomationAccount(context.Background(), resourceGroupName, automationAccountName, filter)
	if err != nil {
		return nil, err
	}

	var dscNodeList []automation.DscNode
	var dscNodeID string
	for dscNodeListResultPage.NotDone() {
		dscNodeList = dscNodeListResultPage.Values()
		// Loop through compilation jobs in the current page
		for _, element := range dscNodeList {
			if *element.Name == vmName && *element.NodeConfiguration.Name == dscConfiguraitonName {
				dscNodeID = *element.NodeID
			}
		}

		err := dscNodeListResultPage.Next()
		if err != nil {
			return nil, err
		}
	}

	// Get Automation Account Connection
	dscNodeConfig, err := client.Get(context.Background(), resourceGroupName, automationAccountName, dscNodeID)
	if err != nil {
		return nil, err
	}

	return &dscNodeConfig, nil
}

// GetAutomationAccountRunAsConnection returns the RunAs Connection if it exists in the Azure Automation Account
func GetAutomationAccountRunAsConnection(t testing.TestingT, automationAccountRunAsConnectionName string, automationAccountName string, resourceGroupName string, subscriptionID string) *automation.Connection {
	client, err := CreateAutomationAccountRunAsConnectionClientE(subscriptionID)
	require.NoError(t, err)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	require.NoError(t, err)
	client.Authorizer = *authorizer

	// Get Automation Account Connection
	automationAccountRunAsConnection, err := client.Get(context.Background(), resourceGroupName, automationAccountName, automationAccountRunAsConnectionName)
	require.NoError(t, err)

	return &automationAccountRunAsConnection
}

// GetAutomationAccountRunAsConnectionE gets the RunAs Connection if it exists in the Azure Automation Account
func GetAutomationAccountRunAsConnectionE(t testing.TestingT, automationAccountRunAsConnectionName string, automationAccountName string, resourceGroupName string, subscriptionID string) (*automation.Connection, error) {
	client, err := CreateAutomationAccountRunAsConnectionClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	// Get Automation Account Connection
	automationAccountRunAsConnection, err := client.Get(context.Background(), resourceGroupName, automationAccountName, automationAccountRunAsConnectionName)
	if err != nil {
		return nil, err
	}

	return &automationAccountRunAsConnection, nil
}
