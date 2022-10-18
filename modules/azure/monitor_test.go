//go:build azure
// +build azure

package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnosticsSettingsResourceExists(t *testing.T) {
	t.Parallel()

	diagnosticsSettingResourceName := "fakename"
	resGroupName := "fakeresgroup"
	subscriptionID := "fakesubid"

	_, err := DiagnosticSettingsResourceExistsE(diagnosticsSettingResourceName, resGroupName, subscriptionID)
	require.Error(t, err)
}
