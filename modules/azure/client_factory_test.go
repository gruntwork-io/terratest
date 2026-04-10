//go:build azure
// +build azure

// This file contains unit tests for the client factory implementation(s).

package azure_test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Local consts for this file only
const govCloudEnvName = "AzureUSGovernmentCloud"
const publicCloudEnvName = "AzurePublicCloud"
const chinaCloudEnvName = "AzureChinaCloud"
const germanyCloudEnvName = "AzureGermanCloud"

func TestDefaultEnvIsPublicWhenNotSet(t *testing.T) {
	// Set env var to missing value
	t.Setenv(azure.AzureEnvironmentEnvName, "")

	// get the default
	env := azure.GetDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, "AzurePublicCloud", env)
}

func TestDefaultEnvSetToGov(t *testing.T) {
	// Set env var to gov
	t.Setenv(azure.AzureEnvironmentEnvName, govCloudEnvName)

	// get the default
	env := azure.GetDefaultEnvironmentName()

	// Make sure it's gov cloud
	assert.Equal(t, govCloudEnvName, env)
}

func TestSubscriptionClientCreation(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"GovCloud/SubscriptionClient", govCloudEnvName, false},
		{"PublicCloud/SubscriptionClient", publicCloudEnvName, false},
		{"ChinaCloud/SubscriptionClient", chinaCloudEnvName, false},
		{"GermanCloud/SubscriptionClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a subscriptions client
			client, err := azure.CreateSubscriptionsClientE()
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestGetClientCloudConfig(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"PublicCloud", publicCloudEnvName, false},
		{"GovCloud", govCloudEnvName, false},
		{"ChinaCloud", chinaCloudEnvName, false},
		{"GermanCloud", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			config, err := azure.GetClientCloudConfig()
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, config.ActiveDirectoryAuthorityHost)
			}
		})
	}
}

func TestCreateVirtualMachinesClientE(t *testing.T) {
	t.Setenv(azure.AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := azure.CreateVirtualMachinesClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestManagedClustersClientCreation(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"GovCloud/ManagedClustersClient", govCloudEnvName, false},
		{"PublicCloud/ManagedClustersClient", publicCloudEnvName, false},
		{"ChinaCloud/ManagedClustersClient", chinaCloudEnvName, false},
		{"GermanCloud/ManagedClustersClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a managed clusters client
			client, err := azure.CreateManagedClustersClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestCosmosDBAccountClientCreation(t *testing.T) {
	t.Setenv(azure.AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := azure.CreateCosmosDBAccountClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestCosmosDBSQLClientCreation(t *testing.T) {
	t.Setenv(azure.AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := azure.CreateCosmosDBSQLClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestPublicIPAddressesClientCreation(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"GovCloud/PublicIPAddressesClient", govCloudEnvName, false},
		{"PublicCloud/PublicIPAddressesClient", publicCloudEnvName, false},
		{"ChinaCloud/PublicIPAddressesClient", chinaCloudEnvName, false},
		{"GermanCloud/PublicIPAddressesClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a PublicIPAddresses client
			client, err := azure.CreatePublicIPAddressesClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestLoadBalancerClientCreation(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"GovCloud/LoadBalancersClient", govCloudEnvName, false},
		{"PublicCloud/LoadBalancersClient", publicCloudEnvName, false},
		{"ChinaCloud/LoadBalancersClient", chinaCloudEnvName, false},
		{"GermanCloud/LoadBalancersClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a LoadBalancer client
			client, err := azure.CreateLoadBalancerClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestFrontDoorClientCreation(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"GovCloud/FrontDoorClient", govCloudEnvName, false},
		{"PublicCloud/FrontDoorClient", publicCloudEnvName, false},
		{"ChinaCloud/FrontDoorClient", chinaCloudEnvName, false},
		{"GermanCloud/FrontDoorClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a Front Door client
			client, err := azure.CreateFrontDoorClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestFrontDoorFrontendEndpointClientCreation(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"GovCloud/FrontDoorClient", govCloudEnvName, false},
		{"PublicCloud/FrontDoorClient", publicCloudEnvName, false},
		{"ChinaCloud/FrontDoorClient", chinaCloudEnvName, false},
		{"GermanCloud/FrontDoorClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a AFD frontend endpoint client
			client, err := azure.CreateFrontDoorFrontendEndpointClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestCreateManagedEnvironmentsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"Default/ManagedEnvironmentsClient", "", false},
		{"PublicCloud/ManagedEnvironmentsClient", publicCloudEnvName, false},
		{"GovCloud/ManagedEnvironmentsClient", govCloudEnvName, false},
		{"ChinaCloud/ManagedEnvironmentsClient", chinaCloudEnvName, false},
		{"GermanCloud/ManagedEnvironmentsClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			if tt.EnvironmentName != "" {
				t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)
			} else {
				t.Setenv(azure.AzureEnvironmentEnvName, "")
				os.Unsetenv(azure.AzureEnvironmentEnvName)
			}

			// Get a ManagedEnvironmentsClient client
			client, err := azure.CreateManagedEnvironmentsClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestCreateContainerAppsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"Default/ContainerAppsClient", "", false},
		{"PublicCloud/ContainerAppsClient", publicCloudEnvName, false},
		{"GovCloud/ContainerAppsClient", govCloudEnvName, false},
		{"ChinaCloud/ContainerAppsClient", chinaCloudEnvName, false},
		{"GermanCloud/ContainerAppsClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			if tt.EnvironmentName != "" {
				t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)
			} else {
				t.Setenv(azure.AzureEnvironmentEnvName, "")
				os.Unsetenv(azure.AzureEnvironmentEnvName)
			}

			// Get a ContainerApps client
			client, err := azure.CreateContainerAppsClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestCreateContainerAppJobsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectErr       bool
	}{
		{"Default/ContainerAppJobsClient", "", false},
		{"PublicCloud/ContainerAppJobsClient", publicCloudEnvName, false},
		{"GovCloud/ContainerAppJobsClient", govCloudEnvName, false},
		{"ChinaCloud/ContainerAppJobsClient", chinaCloudEnvName, false},
		{"GermanCloud/ContainerAppJobsClient", germanyCloudEnvName, true},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			if tt.EnvironmentName != "" {
				t.Setenv(azure.AzureEnvironmentEnvName, tt.EnvironmentName)
			} else {
				t.Setenv(azure.AzureEnvironmentEnvName, "")
				os.Unsetenv(azure.AzureEnvironmentEnvName)
			}

			// Get a ContainerAppJobs client
			client, err := azure.CreateContainerAppJobsClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}
