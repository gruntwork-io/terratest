//go:build azure
// +build azure

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

func TestDiagnosticsSettingsResourceExists(t *testing.T) {
	t.Parallel()

	diagnosticsSettingResourceName := "fakename"
	resGroupName := "fakeresgroup"
	subscriptionID := "fakesubid"

	_, err := azure.DiagnosticSettingsResourceExistsE(diagnosticsSettingResourceName, resGroupName, subscriptionID)
	require.Error(t, err)
}
