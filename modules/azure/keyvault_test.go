package azure_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	azure "github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete key vault resources are added, these tests can be extended.
*/

func TestKeyVaultSecretExists(t *testing.T) {
	t.Parallel()

	testKeyVaultName := "fakeKeyVault"
	testKeyVaultSecretName := "fakeSecretName"
	_, err := azure.KeyVaultSecretExistsContextE(t.Context(), testKeyVaultName, testKeyVaultSecretName)
	require.Error(t, err)
}

func TestKeyVaultKeyExists(t *testing.T) {
	t.Parallel()

	testKeyVaultName := "fakeKeyVault"
	testKeyVaultKeyName := "fakeKeyName"
	_, err := azure.KeyVaultKeyExistsContextE(t.Context(), testKeyVaultName, testKeyVaultKeyName)
	require.Error(t, err)
}

func TestKeyVaultCertificateExists(t *testing.T) {
	t.Parallel()

	testKeyVaultName := "fakeKeyVault"
	testKeyVaultCertName := "fakeCertName"
	_, err := azure.KeyVaultCertificateExistsContextE(t.Context(), testKeyVaultName, testKeyVaultCertName)
	require.Error(t, err)
}

func TestGetKeyVault(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	keyVaultName := ""
	subscriptionID := ""

	_, err := azure.GetKeyVaultContextE(t, t.Context(), resGroupName, keyVaultName, subscriptionID)
	require.Error(t, err)
}
