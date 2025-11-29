package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/require"
)

// NewAzureCredentialE creates a new Azure credential using DefaultAzureCredential.
func NewAzureCredentialE() (*azidentity.DefaultAzureCredential, error) {
	return azidentity.NewDefaultAzureCredential(nil)
}

// KeyVaultSecretExists indicates whether a key vault secret exists; otherwise false
// This function would fail the test if there is an error.
func KeyVaultSecretExists(t *testing.T, keyVaultName string, secretName string) bool {
	result, err := KeyVaultSecretExistsE(keyVaultName, secretName)
	require.NoError(t, err)
	return result
}

// KeyVaultKeyExists indicates whether a key vault key exists; otherwise false.
// This function would fail the test if there is an error.
func KeyVaultKeyExists(t *testing.T, keyVaultName string, keyName string) bool {
	result, err := KeyVaultKeyExistsE(keyVaultName, keyName)
	require.NoError(t, err)
	return result
}

// KeyVaultCertificateExists indicates whether a key vault certificate exists; otherwise false.
// This function would fail the test if there is an error.
func KeyVaultCertificateExists(t *testing.T, keyVaultName string, certificateName string) bool {
	result, err := KeyVaultCertificateExistsE(keyVaultName, certificateName)
	require.NoError(t, err)
	return result
}

// KeyVaultCertificateExistsE indicates whether a certificate exists in key vault; otherwise false.
func KeyVaultCertificateExistsE(keyVaultName, certificateName string) (bool, error) {
	client, err := GetKeyVaultCertificatesClientE(keyVaultName)
	if err != nil {
		return false, err
	}

	pager := client.NewListCertificatePropertiesVersionsPager(certificateName, nil)
	if pager.More() {
		_, err := pager.NextPage(context.Background())
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
func KeyVaultKeyExistsE(keyVaultName, keyName string) (bool, error) {
	client, err := GetKeyVaultKeysClientE(keyVaultName)
	if err != nil {
		return false, err
	}

	pager := client.NewListKeyPropertiesVersionsPager(keyName, nil)
	if pager.More() {
		_, err := pager.NextPage(context.Background())
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
func KeyVaultSecretExistsE(keyVaultName, secretName string) (bool, error) {
	client, err := GetKeyVaultSecretsClientE(keyVaultName)
	if err != nil {
		return false, err
	}

	pager := client.NewListSecretPropertiesVersionsPager(secretName, nil)
	if pager.More() {
		_, err := pager.NextPage(context.Background())
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

// GetKeyVault is a helper function that gets the keyvault management object.
// This function would fail the test if there is an error.
func GetKeyVault(t *testing.T, resGroupName string, keyVaultName string, subscriptionID string) *armkeyvault.Vault {
	keyVault, err := GetKeyVaultE(t, resGroupName, keyVaultName, subscriptionID)
	require.NoError(t, err)

	return keyVault
}

// GetKeyVaultE is a helper function that gets the keyvault management object.
func GetKeyVaultE(t *testing.T, resGroupName string, keyVaultName string, subscriptionID string) (*armkeyvault.Vault, error) {
	// Create a key vault management client
	vaultClient, err := GetKeyVaultManagementClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding vault
	resp, err := vaultClient.Get(context.Background(), resGroupName, keyVaultName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Vault, nil
}

// GetKeyVaultManagementClientE is a helper function that will setup a key vault management client
func GetKeyVaultManagementClientE(subscriptionID string) (*armkeyvault.VaultsClient, error) {
	clientFactory, err := getArmKeyVaultClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewVaultsClient(), nil
}
