//go:build azure
// +build azure

// This file contains unit tests for the client factory implementation(s).

package azure

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Local consts for this file only
const govCloudEnvName = "AzureUSGovernmentCloud"
const publicCloudEnvName = "AzurePublicCloud"
const chinaCloudEnvName = "AzureChinaCloud"
const germanyCloudEnvName = "AzureGermanCloud"

// expectedCloudConfig maps environment names to their expected cloud configurations.
var expectedCloudConfig = map[string]cloud.Configuration{
	publicCloudEnvName: cloud.AzurePublic,
	govCloudEnvName:    cloud.AzureGovernment,
	chinaCloudEnvName:  cloud.AzureChina,
}

func TestDefaultEnvIsPublicWhenNotSet(t *testing.T) {
	// save any current env value and restore on exit
	originalEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, originalEnv)

	// Set env var to missing value
	os.Setenv(AzureEnvironmentEnvName, "")

	// get the default
	env := getDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, "AzurePublicCloud", env)
}

func TestDefaultEnvSetToGov(t *testing.T) {
	// save any current env value and restore on exit
	originalEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, originalEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, govCloudEnvName)

	// get the default
	env := getDefaultEnvironmentName()

	// Make sure it's gov cloud
	assert.Equal(t, govCloudEnvName, env)
}

func TestGetClientCloudConfig(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedConfig  cloud.Configuration
		ExpectErr       bool
	}{
		{"PublicCloud", publicCloudEnvName, cloud.AzurePublic, false},
		{"GovCloud", govCloudEnvName, cloud.AzureGovernment, false},
		{"ChinaCloud", chinaCloudEnvName, cloud.AzureChina, false},
		{"GermanCloud", germanyCloudEnvName, cloud.Configuration{}, true},
	}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			config, err := getClientCloudConfig()
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.ExpectedConfig.ActiveDirectoryAuthorityHost, config.ActiveDirectoryAuthorityHost)
			}
		})
	}
}

func TestGetKeyVaultURISuffixE(t *testing.T) {
	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	os.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	suffix, err := GetKeyVaultURISuffixE()
	require.NoError(t, err)
	assert.Equal(t, "vault.azure.net", suffix)
}

func TestGetStorageURISuffixE(t *testing.T) {
	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	os.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	suffix, err := GetStorageURISuffixE()
	require.NoError(t, err)
	assert.Equal(t, "core.windows.net", suffix)
}

func TestClientCreation(t *testing.T) {
	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)
	os.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	t.Run("VirtualMachinesClient", func(t *testing.T) {
		client, err := CreateVirtualMachinesClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("ManagedClustersClient", func(t *testing.T) {
		client, err := CreateManagedClustersClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("SubscriptionsClient", func(t *testing.T) {
		client, err := CreateSubscriptionsClientE()
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("CosmosDBAccountClient", func(t *testing.T) {
		client, err := CreateCosmosDBAccountClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("PublicIPAddressesClient", func(t *testing.T) {
		client, err := CreatePublicIPAddressesClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("LoadBalancerClient", func(t *testing.T) {
		client, err := CreateLoadBalancerClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("FrontDoorClient", func(t *testing.T) {
		client, err := CreateFrontDoorClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("ManagedEnvironmentsClient", func(t *testing.T) {
		client, err := CreateManagedEnvironmentsClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})

	t.Run("ContainerAppsClient", func(t *testing.T) {
		client, err := CreateContainerAppsClientE("")
		require.NoError(t, err)
		require.NotNil(t, client)
	})
}

func TestClientCreationFailsForUnsupportedCloud(t *testing.T) {
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	os.Setenv(AzureEnvironmentEnvName, germanyCloudEnvName)

	_, err := CreateManagedEnvironmentsClientE("")
	require.Error(t, err)
}
