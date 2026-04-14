//go:build azure
// +build azure

// This file contains unit tests for the client factory implementation(s).

package azure //nolint:testpackage // tests access unexported functions

import (
	"os"
	"reflect"
	"testing"

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
	t.Setenv(AzureEnvironmentEnvName, "")

	// get the default
	env := getDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, "AzurePublicCloud", env)
}

func TestDefaultEnvSetToGov(t *testing.T) {
	// Set env var to gov
	t.Setenv(AzureEnvironmentEnvName, govCloudEnvName)

	// get the default
	env := getDefaultEnvironmentName()

	// Make sure it's gov cloud
	assert.Equal(t, govCloudEnvName, env)
}

func TestSubscriptionClientCreation(t *testing.T) {
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateSubscriptionsClientContextE(t.Context())
	require.NoError(t, err)
	require.NotNil(t, client)
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
			t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			config, err := getClientCloudConfig()
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
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateVirtualMachinesClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestManagedClustersClientCreation(t *testing.T) {
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateManagedClustersClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestCosmosDBAccountClientCreation(t *testing.T) {
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateCosmosDBAccountClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestCosmosDBSQLClientCreation(t *testing.T) {
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateCosmosDBSQLClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestPublicIPAddressesClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
		ExpectErr       bool
	}{
		{"GovCloud/PublicIPAddressesClient", govCloudEnvName, "https://management.usgovcloudapi.net/", false},
		{"PublicCloud/PublicIPAddressesClient", publicCloudEnvName, "https://management.azure.com/", false},
		{"ChinaCloud/PublicIPAddressesClient", chinaCloudEnvName, "https://management.chinacloudapi.cn/", false},
		{"GermanCloud/PublicIPAddressesClient", germanyCloudEnvName, "https://management.microsoftazure.de/", true},
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a PublicIPAddresses client
			client, err := CreatePublicIPAddressesClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				// Not ideal, but to get the base URI we need to access the internal field
				internalField := reflect.ValueOf(client).Elem().FieldByName("internal")
				require.True(t, internalField.IsValid(), "internal field not found - Azure SDK may have changed")
				epField := internalField.Elem().FieldByName("ep")
				require.True(t, epField.IsValid(), "ep field not found - Azure SDK may have changed")
				assert.Equal(t, epField.String()+"/", tt.ExpectedBaseURI)
			}
		})
	}
}

func TestLoadBalancerClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
		ExpectErr       bool
	}{
		{"GovCloud/LoadBalancersClient", govCloudEnvName, "https://management.usgovcloudapi.net/", false},
		{"PublicCloud/LoadBalancersClient", publicCloudEnvName, "https://management.azure.com/", false},
		{"ChinaCloud/LoadBalancersClient", chinaCloudEnvName, "https://management.chinacloudapi.cn/", false},
		{"GermanCloud/LoadBalancersClient", germanyCloudEnvName, "https://management.microsoftazure.de/", true},
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a LoadBalancer client
			client, err := CreateLoadBalancerClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				// Not ideal, but to get the base URI we need to access the internal field
				internalField := reflect.ValueOf(client).Elem().FieldByName("internal")
				require.True(t, internalField.IsValid(), "internal field not found - Azure SDK may have changed")
				epField := internalField.Elem().FieldByName("ep")
				require.True(t, epField.IsValid(), "ep field not found - Azure SDK may have changed")
				assert.Equal(t, epField.String()+"/", tt.ExpectedBaseURI)
			}
		})
	}
}

func TestFrontDoorClientCreation(t *testing.T) {
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateFrontDoorClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestFrontDoorFrontendEndpointClientCreation(t *testing.T) {
	t.Setenv(AzureEnvironmentEnvName, publicCloudEnvName)

	client, err := CreateFrontDoorFrontendEndpointClientE("")
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestCreateManagedEnvironmentsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
		ExpectErr       bool
	}{
		{"Default/ManagedEnvironmentsClient", "", "https://management.azure.com/", false},
		{"PublicCloud/ManagedEnvironmentsClient", publicCloudEnvName, "https://management.azure.com/", false},
		{"GovCloud/ManagedEnvironmentsClient", govCloudEnvName, "https://management.usgovcloudapi.net/", false},
		{"ChinaCloud/ManagedEnvironmentsClient", chinaCloudEnvName, "https://management.chinacloudapi.cn/", false},
		{"GermanCloud/ManagedEnvironmentsClient", germanyCloudEnvName, "https://management.microsoftazure.de/", true}, // GermanCloud is deleted as of 2021-10-21 https://learn.microsoft.com/en-us/previous-versions/azure/germany/germany-welcome
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			if tt.EnvironmentName != "" {
				t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)
			} else {
				t.Setenv(AzureEnvironmentEnvName, "")
				os.Unsetenv(AzureEnvironmentEnvName)
			}

			// Get a ManagedEnvironmentsClient client
			client, err := CreateManagedEnvironmentsClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				// Not ideal, but to get the base URI we need to access the internal field
				internalField := reflect.ValueOf(client).Elem().FieldByName("internal")
				require.True(t, internalField.IsValid(), "internal field not found - Azure SDK may have changed")
				epField := internalField.Elem().FieldByName("ep")
				require.True(t, epField.IsValid(), "ep field not found - Azure SDK may have changed")
				assert.Equal(t, epField.String()+"/", tt.ExpectedBaseURI)
			}
		})
	}
}

func TestCreateContainerAppsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
		ExpectErr       bool
	}{
		{"Default/ContainerAppsClient", "", "https://management.azure.com/", false},
		{"PublicCloud/ContainerAppsClient", publicCloudEnvName, "https://management.azure.com/", false},
		{"GovCloud/ContainerAppsClient", govCloudEnvName, "https://management.usgovcloudapi.net/", false},
		{"ChinaCloud/ContainerAppsClient", chinaCloudEnvName, "https://management.chinacloudapi.cn/", false},
		{"GermanCloud/ContainerAppsClient", germanyCloudEnvName, "https://management.microsoftazure.de/", true}, // GermanCloud is deleted as of 2021-10-21 https://learn.microsoft.com/en-us/previous-versions/azure/germany/germany-welcome
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			if tt.EnvironmentName != "" {
				t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)
			} else {
				t.Setenv(AzureEnvironmentEnvName, "")
				os.Unsetenv(AzureEnvironmentEnvName)
			}

			// Get a ContainerApps client
			client, err := CreateContainerAppsClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				// Not ideal, but to get the base URI we need to access the internal field
				internalField := reflect.ValueOf(client).Elem().FieldByName("internal")
				require.True(t, internalField.IsValid(), "internal field not found - Azure SDK may have changed")
				epField := internalField.Elem().FieldByName("ep")
				require.True(t, epField.IsValid(), "ep field not found - Azure SDK may have changed")
				assert.Equal(t, epField.String()+"/", tt.ExpectedBaseURI)
			}
		})
	}
}

func TestCreateContainerAppJobsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
		ExpectErr       bool
	}{
		{"Default/ContainerAppJobsClient", "", "https://management.azure.com/", false},
		{"PublicCloud/ContainerAppJobsClient", publicCloudEnvName, "https://management.azure.com/", false},
		{"GovCloud/ContainerAppJobsClient", govCloudEnvName, "https://management.usgovcloudapi.net/", false},
		{"ChinaCloud/ContainerAppJobsClient", chinaCloudEnvName, "https://management.chinacloudapi.cn/", false},
		{"GermanCloud/ContainerAppJobsClient", germanyCloudEnvName, "https://management.microsoftazure.de/", true}, // GermanCloud is deleted as of 2021-10-21 https://learn.microsoft.com/en-us/previous-versions/azure/germany/germany-welcome
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			if tt.EnvironmentName != "" {
				t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)
			} else {
				t.Setenv(AzureEnvironmentEnvName, "")
				os.Unsetenv(AzureEnvironmentEnvName)
			}

			// Get a ContainerAppJobs client
			client, err := CreateContainerAppJobsClientE("")
			if tt.ExpectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				// Not ideal, but to get the base URI we need to access the internal field
				internalField := reflect.ValueOf(client).Elem().FieldByName("internal")
				require.True(t, internalField.IsValid(), "internal field not found - Azure SDK may have changed")
				epField := internalField.Elem().FieldByName("ep")
				require.True(t, epField.IsValid(), "ep field not found - Azure SDK may have changed")
				assert.Equal(t, epField.String()+"/", tt.ExpectedBaseURI)
			}
		})
	}
}
