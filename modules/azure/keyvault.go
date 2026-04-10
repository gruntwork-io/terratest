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

// KeyVaultKeyExistsContext indicates whether a key vault key exists; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultKeyExistsContext(t testing.TestingT, ctx context.Context, keyVaultName string, keyName string) bool {
	t.Helper()

	result, err := KeyVaultKeyExistsContextE(ctx, keyVaultName, keyName)
	require.NoError(t, err)

	return result
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

// KeyVaultCertificateExistsContextE indicates whether a certificate exists in key vault; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultCertificateExistsContextE(ctx context.Context, keyVaultName, certificateName string) (bool, error) {
	client, err := GetKeyVaultCertificatesClientE(keyVaultName)
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

// KeyVaultKeyExistsContextE indicates whether a key exists in the key vault; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultKeyExistsContextE(ctx context.Context, keyVaultName, keyName string) (bool, error) {
	client, err := GetKeyVaultKeysClientE(keyVaultName)
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

// KeyVaultSecretExistsContextE indicates whether a secret exists in the key vault; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func KeyVaultSecretExistsContextE(ctx context.Context, keyVaultName, secretName string) (bool, error) {
	client, err := GetKeyVaultSecretsClientE(keyVaultName)
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

// GetKeyVaultSecretsClientE creates a KeyVault secrets client.
func GetKeyVaultSecretsClientE(keyVaultName string) (*azsecrets.Client, error) {
	keyVaultSuffix, err := GetKeyVaultURISuffixE()
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

// GetKeyVaultKeysClientE creates a KeyVault keys client.
func GetKeyVaultKeysClientE(keyVaultName string) (*azkeys.Client, error) {
	keyVaultSuffix, err := GetKeyVaultURISuffixE()
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

// GetKeyVaultCertificatesClientE creates a KeyVault certificates client.
func GetKeyVaultCertificatesClientE(keyVaultName string) (*azcertificates.Client, error) {
	keyVaultSuffix, err := GetKeyVaultURISuffixE()
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

// GetKeyVaultContext is a helper function that gets the keyvault management object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultContext(t testing.TestingT, ctx context.Context, resGroupName string, keyVaultName string, subscriptionID string) *armkeyvault.Vault {
	t.Helper()

	keyVault, err := GetKeyVaultContextE(t, ctx, resGroupName, keyVaultName, subscriptionID)
	require.NoError(t, err)

	return keyVault
}

// GetKeyVaultContextE is a helper function that gets the keyvault management object.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultContextE(t testing.TestingT, ctx context.Context, resGroupName string, keyVaultName string, subscriptionID string) (*armkeyvault.Vault, error) {
	t.Helper()

	// Create a key vault management client
	vaultClient, err := GetKeyVaultManagementClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding vault
	resp, err := vaultClient.Get(ctx, resGroupName, keyVaultName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Vault, nil
}

// GetKeyVaultManagementClientE is a helper function that will setup a key vault management client.
func GetKeyVaultManagementClientE(subscriptionID string) (*armkeyvault.VaultsClient, error) {
	clientFactory, err := getArmKeyVaultClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVaultsClient(), nil
}
