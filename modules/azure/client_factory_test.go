//go:build azure
// +build azure

// This file contains unit tests for the client factory implementation(s).

package azure //nolint:testpackage // tests access unexported functions

import (
	"os"
	"reflect"
	"testing"

	autorest "github.com/Azure/go-autorest/autorest/azure"
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
		{"GovCloud/PublicIPAddressesClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint, false},
		{"PublicCloud/PublicIPAddressesClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"ChinaCloud/PublicIPAddressesClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint, false},
		{"GermanCloud/PublicIPAddressesClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint, true},
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
		{"GovCloud/LoadBalancersClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint, false},
		{"PublicCloud/LoadBalancersClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"ChinaCloud/LoadBalancersClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint, false},
		{"GermanCloud/LoadBalancersClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint, true},
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

func TestFrontDoorClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
	}{
		{"GovCloud/FrontDoorClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{"PublicCloud/FrontDoorClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
		{"ChinaCloud/FrontDoorClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
		{"GermanCloud/FrontDoorClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a Front Door client
			client, err := CreateFrontDoorClientE("")
			require.NoError(t, err)

			// Check for correct ARM URI
			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
		})
	}
}

func TestFrontDoorFrontendEndpointClientBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
	}{
		{"GovCloud/FrontDoorClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{"PublicCloud/FrontDoorClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
		{"ChinaCloud/FrontDoorClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
		{"GermanCloud/FrontDoorClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
	}

	for _, tt := range cases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		tt := tt
		t.Run(tt.CaseName, func(t *testing.T) {
			// Override env setting
			t.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a AFD frontend endpoint client
			client, err := CreateFrontDoorFrontendEndpointClientE("")
			require.NoError(t, err)

			// Check for correct ARM URI
			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
		})
	}
}

func TestCreateManagedEnvironmentsClientEEndpointURISetCorrectly(t *testing.T) {
	var cases = []struct {
		CaseName        string
		EnvironmentName string
		ExpectedBaseURI string
		ExpectErr       bool
	}{
		{"Default/ManagedEnvironmentsClient", "", autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"PublicCloud/ManagedEnvironmentsClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"GovCloud/ManagedEnvironmentsClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint, false},
		{"ChinaCloud/ManagedEnvironmentsClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint, false},
		{"GermanCloud/ManagedEnvironmentsClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint, true}, // GermanCloud is deleted as of 2021-10-21 https://learn.microsoft.com/en-us/previous-versions/azure/germany/germany-welcome
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
		{"Default/ContainerAppsClient", "", autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"PublicCloud/ContainerAppsClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"GovCloud/ContainerAppsClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint, false},
		{"ChinaCloud/ContainerAppsClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint, false},
		{"GermanCloud/ContainerAppsClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint, true}, // GermanCloud is deleted as of 2021-10-21 https://learn.microsoft.com/en-us/previous-versions/azure/germany/germany-welcome
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
		{"Default/ContainerAppJobsClient", "", autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"PublicCloud/ContainerAppJobsClient", publicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint, false},
		{"GovCloud/ContainerAppJobsClient", govCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint, false},
		{"ChinaCloud/ContainerAppJobsClient", chinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint, false},
		{"GermanCloud/ContainerAppJobsClient", germanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint, true}, // GermanCloud is deleted as of 2021-10-21 https://learn.microsoft.com/en-us/previous-versions/azure/germany/germany-welcome
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
