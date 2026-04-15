package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// NewAzureCredentialE creates a new Azure credential using DefaultAzureCredential.
func NewAzureCredentialE() (*azidentity.DefaultAzureCredential, error) {
	return azidentity.NewDefaultAzureCredential(nil)
}

// KeyVaultSecretExistsContext indicates whether a key vault secret exists; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultSecretExistsContext(t testing.TestingT, ctx context.Context, keyVaultName string, secretName string) bool {
	t.Helper()

	result, err := KeyVaultSecretExistsContextE(ctx, keyVaultName, secretName)
	require.NoError(t, err)

	return result
}

// KeyVaultSecretExists indicates whether a key vault secret exists; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [KeyVaultSecretExistsContext] instead.
func KeyVaultSecretExists(t testing.TestingT, keyVaultName string, secretName string) bool {
	t.Helper()

	return KeyVaultSecretExistsContext(t, context.Background(), keyVaultName, secretName) //nolint:staticcheck
}

// KeyVaultKeyExistsContext indicates whether a key vault key exists; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultKeyExistsContext(t testing.TestingT, ctx context.Context, keyVaultName string, keyName string) bool {
	t.Helper()

	result, err := KeyVaultKeyExistsContextE(ctx, keyVaultName, keyName)
	require.NoError(t, err)

	return result
}

// KeyVaultKeyExists indicates whether a key vault key exists; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [KeyVaultKeyExistsContext] instead.
func KeyVaultKeyExists(t testing.TestingT, keyVaultName string, keyName string) bool {
	t.Helper()

	return KeyVaultKeyExistsContext(t, context.Background(), keyVaultName, keyName) //nolint:staticcheck
}

// KeyVaultCertificateExistsContext indicates whether a key vault certificate exists; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultCertificateExistsContext(t testing.TestingT, ctx context.Context, keyVaultName string, certificateName string) bool {
	t.Helper()

	result, err := KeyVaultCertificateExistsContextE(ctx, keyVaultName, certificateName)
	require.NoError(t, err)

	return result
}

// KeyVaultCertificateExists indicates whether a key vault certificate exists; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [KeyVaultCertificateExistsContext] instead.
func KeyVaultCertificateExists(t testing.TestingT, keyVaultName string, certificateName string) bool {
	t.Helper()

	return KeyVaultCertificateExistsContext(t, context.Background(), keyVaultName, certificateName) //nolint:staticcheck
}

// KeyVaultCertificateExistsContextE indicates whether a certificate exists in key vault; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultCertificateExistsContextE(ctx context.Context, keyVaultName, certificateName string) (bool, error) {
	client, err := GetKeyVaultCertificatesClientContextE(ctx, keyVaultName)
	if err != nil {
		return false, err
	}

	pager := client.NewListCertificatePropertiesVersionsPager(certificateName, nil)

	if pager.More() {
		_, err := pager.NextPage(ctx)
		if err != nil {
			if ResourceNotFoundErrorExists(err) {
				return false, nil
			}

			return false, err
		}

		return true, nil
	}

	return false, nil
}

// KeyVaultCertificateExistsE indicates whether a certificate exists in key vault; otherwise false.
//
// Deprecated: Use [KeyVaultCertificateExistsContextE] instead.
func KeyVaultCertificateExistsE(keyVaultName, certificateName string) (bool, error) {
	return KeyVaultCertificateExistsContextE(context.Background(), keyVaultName, certificateName)
}

// KeyVaultKeyExistsContextE indicates whether a key exists in the key vault; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultKeyExistsContextE(ctx context.Context, keyVaultName, keyName string) (bool, error) {
	client, err := GetKeyVaultKeysClientContextE(ctx, keyVaultName)
	if err != nil {
		return false, err
	}

	pager := client.NewListKeyPropertiesVersionsPager(keyName, nil)

	if pager.More() {
		_, err := pager.NextPage(ctx)
		if err != nil {
			if ResourceNotFoundErrorExists(err) {
				return false, nil
			}

			return false, err
		}

		return true, nil
	}

	return false, nil
}

// KeyVaultKeyExistsE indicates whether a key exists in the key vault; otherwise false.
//
// Deprecated: Use [KeyVaultKeyExistsContextE] instead.
func KeyVaultKeyExistsE(keyVaultName, keyName string) (bool, error) {
	return KeyVaultKeyExistsContextE(context.Background(), keyVaultName, keyName)
}

// KeyVaultSecretExistsContextE indicates whether a secret exists in the key vault; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultSecretExistsContextE(ctx context.Context, keyVaultName, secretName string) (bool, error) {
	client, err := GetKeyVaultSecretsClientContextE(ctx, keyVaultName)
	if err != nil {
		return false, err
	}

	pager := client.NewListSecretPropertiesVersionsPager(secretName, nil)

	if pager.More() {
		_, err := pager.NextPage(ctx)
		if err != nil {
			if ResourceNotFoundErrorExists(err) {
				return false, nil
			}

			return false, err
		}

		return true, nil
	}

	return false, nil
}

// KeyVaultSecretExistsE indicates whether a secret exists in the key vault; otherwise false.
//
// Deprecated: Use [KeyVaultSecretExistsContextE] instead.
func KeyVaultSecretExistsE(keyVaultName, secretName string) (bool, error) {
	return KeyVaultSecretExistsContextE(context.Background(), keyVaultName, secretName)
}

// GetKeyVaultSecretsClientContextE creates a KeyVault secrets client.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultSecretsClientContextE(_ context.Context, keyVaultName string) (*azsecrets.Client, error) {
	keyVaultSuffix, err := GetKeyVaultURISuffixE() //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	vaultURL := fmt.Sprintf("https://%s.%s", keyVaultName, keyVaultSuffix)

	cred, err := NewAzureCredentialE()
	if err != nil {
		return nil, err
	}

	return azsecrets.NewClient(vaultURL, cred, nil)
}

// GetKeyVaultSecretsClientE creates a KeyVault secrets client.
//
// Deprecated: Use [GetKeyVaultSecretsClientContextE] instead.
func GetKeyVaultSecretsClientE(keyVaultName string) (*azsecrets.Client, error) {
	return GetKeyVaultSecretsClientContextE(context.Background(), keyVaultName)
}

// GetKeyVaultKeysClientContextE creates a KeyVault keys client.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultKeysClientContextE(_ context.Context, keyVaultName string) (*azkeys.Client, error) {
	keyVaultSuffix, err := GetKeyVaultURISuffixE() //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	vaultURL := fmt.Sprintf("https://%s.%s", keyVaultName, keyVaultSuffix)

	cred, err := NewAzureCredentialE()
	if err != nil {
		return nil, err
	}

	return azkeys.NewClient(vaultURL, cred, nil)
}

// GetKeyVaultKeysClientE creates a KeyVault keys client.
//
// Deprecated: Use [GetKeyVaultKeysClientContextE] instead.
func GetKeyVaultKeysClientE(keyVaultName string) (*azkeys.Client, error) {
	return GetKeyVaultKeysClientContextE(context.Background(), keyVaultName)
}

// GetKeyVaultCertificatesClientContextE creates a KeyVault certificates client.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultCertificatesClientContextE(_ context.Context, keyVaultName string) (*azcertificates.Client, error) {
	keyVaultSuffix, err := GetKeyVaultURISuffixE() //nolint:contextcheck
	if err != nil {
		return nil, err
	}

	vaultURL := fmt.Sprintf("https://%s.%s", keyVaultName, keyVaultSuffix)

	cred, err := NewAzureCredentialE()
	if err != nil {
		return nil, err
	}

	return azcertificates.NewClient(vaultURL, cred, nil)
}

// GetKeyVaultCertificatesClientE creates a KeyVault certificates client.
//
// Deprecated: Use [GetKeyVaultCertificatesClientContextE] instead.
func GetKeyVaultCertificatesClientE(keyVaultName string) (*azcertificates.Client, error) {
	return GetKeyVaultCertificatesClientContextE(context.Background(), keyVaultName)
}

// GetKeyVaultContext is a helper function that gets the keyvault management object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultContext(t testing.TestingT, ctx context.Context, resGroupName string, keyVaultName string, subscriptionID string) *armkeyvault.Vault {
	t.Helper()

	keyVault, err := GetKeyVaultContextE(t, ctx, resGroupName, keyVaultName, subscriptionID)
	require.NoError(t, err)

	return keyVault
}

// GetKeyVault is a helper function that gets the keyvault management object.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetKeyVaultContext] instead.
func GetKeyVault(t testing.TestingT, resGroupName string, keyVaultName string, subscriptionID string) *armkeyvault.Vault {
	t.Helper()

	return GetKeyVaultContext(t, context.Background(), resGroupName, keyVaultName, subscriptionID) //nolint:staticcheck
}

// GetKeyVaultContextE is a helper function that gets the keyvault management object.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultContextE(t testing.TestingT, ctx context.Context, resGroupName string, keyVaultName string, subscriptionID string) (*armkeyvault.Vault, error) {
	t.Helper()

	// Create a key vault management client
	vaultClient, err := GetKeyVaultManagementClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetKeyVaultWithClient(ctx, vaultClient, resGroupName, keyVaultName)
}

// GetKeyVaultWithClient gets the specified Key Vault using the provided VaultsClient.
// This variant is useful for testing with fake clients.
func GetKeyVaultWithClient(ctx context.Context, client *armkeyvault.VaultsClient, resGroupName string, keyVaultName string) (*armkeyvault.Vault, error) {
	resp, err := client.Get(ctx, resGroupName, keyVaultName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Vault, nil
}

// GetKeyVaultE is a helper function that gets the keyvault management object.
//
// Deprecated: Use [GetKeyVaultContextE] instead.
func GetKeyVaultE(t testing.TestingT, resGroupName string, keyVaultName string, subscriptionID string) (*armkeyvault.Vault, error) {
	t.Helper()

	return GetKeyVaultContextE(t, context.Background(), resGroupName, keyVaultName, subscriptionID)
}

// GetKeyVaultManagementClientContextE is a helper function that will setup a key vault management client.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultManagementClientContextE(_ context.Context, subscriptionID string) (*armkeyvault.VaultsClient, error) {
	clientFactory, err := getArmKeyVaultClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVaultsClient(), nil
}

// GetKeyVaultManagementClientE is a helper function that will setup a key vault management client.
//
// Deprecated: Use [GetKeyVaultManagementClientContextE] instead.
func GetKeyVaultManagementClientE(subscriptionID string) (*armkeyvault.VaultsClient, error) {
	return GetKeyVaultManagementClientContextE(context.Background(), subscriptionID)
}
