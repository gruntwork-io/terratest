// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete automation account resources are added, these tests can be extended.
*/

func TestAutomationAccountExistsE(t *testing.T) {
	t.Parallel()

	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := AutomationAccountExistsE(t, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestAutomationAccountDscExistsE(t *testing.T) {
	t.Parallel()

	dscConfigurationName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := AutomationAccountDscExistsE(t, dscConfigurationName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestWaitUntilDscCompiled(t *testing.T) {
	t.Parallel()

	dscConfigurationName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := WaitUntilDscCompiledE(t, dscConfigurationName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestAutomationAccountRunAsCertificateThumbprintMatchesE(t *testing.T) {
	t.Parallel()

	runAsCertificateThumbprint := ""
	runAsCertificateName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := AutomationAccountRunAsCertificateThumbprintMatchesE(t, runAsCertificateThumbprint, runAsCertificateName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestAutomationAccountRunAsConnectionExistsE(t *testing.T) {
	t.Parallel()

	automationAccountrunAsAccountName := ""
	runAsConnectionType := ""
	runAsCertificateThumbprint := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := AutomationAccountRunAsConnectionExistsE(t, automationAccountrunAsAccountName, runAsConnectionType, runAsCertificateThumbprint, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestAutomationAccountDscAppliedSuccessfullyToVME(t *testing.T) {
	t.Parallel()

	dscConfigurationName := ""
	vmName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := AutomationAccountDscAppliedSuccessfullyToVME(t, dscConfigurationName, vmName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
}

func TestGetAutomationAccountE(t *testing.T) {
	t.Parallel()

	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	client, err := GetAutomationAccountE(t, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
	assert.Nil(t, client)
}

func TestGetAutomationAccountDscConfigurationE(t *testing.T) {
	t.Parallel()

	dscConfigurationName := ""
	resourceGroupName := ""
	automationAccountName := ""
	subscriptionID := ""

	dscConfiguration, err := GetAutomationAccountDscConfigurationE(t, dscConfigurationName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
	assert.Nil(t, dscConfiguration)
}

func TestAutomationAccountDscCompileJobStatusE(t *testing.T) {
	t.Parallel()

	dscConfigurationName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	status, err := AutomationAccountDscCompileJobStatusE(t, dscConfigurationName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
	assert.Empty(t, status)
}

func TestGetAutomationAccountCertificateE(t *testing.T) {
	t.Parallel()

	automationAccountCertificateName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	certificate, err := GetAutomationAccountCertificateE(t, automationAccountCertificateName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
	assert.Nil(t, certificate)
}

func TestGetAutomationAccountDscNodeConfigurationE(t *testing.T) {
	t.Parallel()

	dscConfigurationName := ""
	vmName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	dscNodeConfig, err := GetAutomationAccountDscNodeConfigurationE(t, dscConfigurationName, vmName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
	assert.Nil(t, dscNodeConfig)
}

func TestGetAutomationAccountRunAsConnectionE(t *testing.T) {
	t.Parallel()

	automationAccountRunAsConnectionName := ""
	automationAccountName := ""
	resourceGroupName := ""
	subscriptionID := ""

	connection, err := GetAutomationAccountRunAsConnectionE(t, automationAccountRunAsConnectionName, automationAccountName, resourceGroupName, subscriptionID)

	require.Error(t, err)
	assert.Nil(t, connection)
}
