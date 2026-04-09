/*

This file implements an Azure client factory that automatically handles setting up Base URI
values for sovereign cloud support. Note the list of clients below is not initially exhaustive;
rather, additional clients will be added as-needed.

*/

package azure

// snippet-tag-start::client_factory_example.imports

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/frontdoor/mgmt/frontdoor"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/privatedns/mgmt/privatedns"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
)

// snippet-tag-end::client_factory_example.imports

const (
	// AzureEnvironmentEnvName is the name of the Azure environment to use. Set to one of the following:
	//
	// "AzureChinaCloud":        ChinaCloud
	// "AzureGermanCloud":       GermanCloud
	// "AzurePublicCloud":       PublicCloud
	// "AzureUSGovernmentCloud": USGovernmentCloud
	// "AzureStackCloud":		 Azure stack
	AzureEnvironmentEnvName = "AZURE_ENVIRONMENT"

	// ResourceManagerEndpointName is the name of the ResourceManagerEndpoint field in the Environment struct.
	ResourceManagerEndpointName = "ResourceManagerEndpoint"

	// Azure environment name constants (upper-cased for case-insensitive switch matching).
	azurePublicCloud = "AZUREPUBLICCLOUD"
	azureUSGovCloud  = "AZUREUSGOVERNMENTCLOUD"
	azureChinaCloud  = "AZURECHINACLOUD"
	azureStackCloud  = "AZURESTACKCLOUD"
)

// ---- Credential & cloud config helpers ----

// newArmCredential creates a DefaultAzureCredential configured for the current cloud environment.
func newArmCredential() (*azidentity.DefaultAzureCredential, error) {
	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	return azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// newArmClientOptions returns arm.ClientOptions configured for the current cloud environment.
func newArmClientOptions() (*arm.ClientOptions, error) {
	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	}, nil
}

// ---- Private ARM factory functions ----

func getArmComputeClientFactory(subscriptionID string) (*armcompute.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armcompute.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getArmNetworkClientFactory(subscriptionID string) (*armnetwork.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armnetwork.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getArmStorageClientFactory(subscriptionID string) (*armstorage.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armstorage.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getArmCosmosClientFactory(subscriptionID string) (*armcosmos.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armcosmos.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getArmServiceBusClientFactory(subscriptionID string) (*armservicebus.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armservicebus.NewClientFactory(targetSubscriptionID, cred, opts)
}

// ---- Public client creator functions ----

// CreateSubscriptionsClientE returns a virtual machines client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateSubscriptionsClientE() (subscriptions.Client, error) {
	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return subscriptions.Client{}, err
	}

	// Create correct client based on type passed
	return subscriptions.NewClientWithBaseURI(baseURI), nil
}

// CreateVirtualMachinesClientE returns a virtual machines client.
func CreateVirtualMachinesClientE(subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualMachinesClient(), nil
}

// CreateManagedClustersClientE returns a virtual machines client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateManagedClustersClientE(subscriptionID string) (containerservice.ManagedClustersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	// Create correct client based on type passed
	return containerservice.NewManagedClustersClientWithBaseURI(baseURI, subscriptionID), nil
}

// CreateCosmosDBAccountClientE returns a Cosmos DB database accounts client.
func CreateCosmosDBAccountClientE(subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	clientFactory, err := getArmCosmosClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabaseAccountsClient(), nil
}

// CreateCosmosDBSQLClientE returns a Cosmos DB SQL resources client.
func CreateCosmosDBSQLClientE(subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
	clientFactory, err := getArmCosmosClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSQLResourcesClient(), nil
}

// getArmKeyVaultClientFactory gets an arm keyvault client factory
func getArmKeyVaultClientFactory(subscriptionID string) (*armkeyvault.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armkeyvault.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// getArmPostgreSQLClientFactory gets an arm postgresql client factory
func getArmPostgreSQLClientFactory(subscriptionID string) (*armpostgresql.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armpostgresql.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// CreateStorageAccountClientE creates a storage account client.
func CreateStorageAccountClientE(subscriptionID string) (*armstorage.AccountsClient, error) {
	clientFactory, err := getArmStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAccountsClient(), nil
}

// CreateStorageBlobContainerClientE creates a storage blob container client.
func CreateStorageBlobContainerClientE(subscriptionID string) (*armstorage.BlobContainersClient, error) {
	clientFactory, err := getArmStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBlobContainersClient(), nil
}

// CreateStorageFileSharesClientE creates a storage file shares client.
func CreateStorageFileSharesClientE(subscriptionID string) (*armstorage.FileSharesClient, error) {
	clientFactory, err := getArmStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFileSharesClient(), nil
}

// CreateServiceBusNamespacesClientE returns a service bus namespaces client.
func CreateServiceBusNamespacesClientE(subscriptionID string) (*armservicebus.NamespacesClient, error) {
	clientFactory, err := getArmServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewNamespacesClient(), nil
}

// CreateServiceBusTopicsClientE returns a service bus topics client.
func CreateServiceBusTopicsClientE(subscriptionID string) (*armservicebus.TopicsClient, error) {
	clientFactory, err := getArmServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewTopicsClient(), nil
}

// CreateServiceBusSubscriptionsClientE returns a service bus subscriptions client.
func CreateServiceBusSubscriptionsClientE(subscriptionID string) (*armservicebus.SubscriptionsClient, error) {
	clientFactory, err := getArmServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSubscriptionsClient(), nil
}

// CreateAvailabilitySetClientE creates a new Availability Set client.
func CreateAvailabilitySetClientE(subscriptionID string) (*armcompute.AvailabilitySetsClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAvailabilitySetsClient(), nil
}

// CreateResourceGroupClientE gets a resource group client in a subscription
func CreateResourceGroupClientE(subscriptionID string) (*resources.GroupsClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	resourceGroupClient := resources.NewGroupsClientWithBaseURI(baseURI, subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	resourceGroupClient.Authorizer = *authorizer

	return &resourceGroupClient, nil
}

// CreateSQLServerClient is a helper function that will create and setup a sql server client
func CreateSQLServerClient(subscriptionID string) (*armsql.ServersClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// CreateSQLMangedInstanceClient is a helper function that will create and setup a sql managed instance client
func CreateSQLMangedInstanceClient(subscriptionID string) (*armsql.ManagedInstancesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedInstancesClient(), nil
}

// CreateSQLMangedDatabasesClient is a helper function that will create and setup a sql managed databases client
func CreateSQLMangedDatabasesClient(subscriptionID string) (*armsql.ManagedDatabasesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedDatabasesClient(), nil
}

// getArmSQLClientFactory gets an arm sql client factory
func getArmSQLClientFactory(subscriptionID string) (*armsql.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armsql.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// CreateDatabaseClient is a helper function that will create and setup a SQL DB client
func CreateDatabaseClient(subscriptionID string) (*armsql.DatabasesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabasesClient(), nil
}

// CreateMySQLServerClientE is a helper function that will setup a mysql server client.
func CreateMySQLServerClientE(subscriptionID string) (*armmysql.ServersClient, error) {
	clientFactory, err := getArmMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// getArmMySQLClientFactory gets an arm mysql client factory
func getArmMySQLClientFactory(subscriptionID string) (*armmysql.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armmysql.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// CreateDisksClientE returns a disks client.
func CreateDisksClientE(subscriptionID string) (*armcompute.DisksClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDisksClient(), nil
}

// CreateActionGroupClient creates an Action Groups client for Azure Monitor.
func CreateActionGroupClient(subscriptionID string) (*insights.ActionGroupsClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	metricAlertsClient := insights.NewActionGroupsClientWithBaseURI(baseURI, subID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	metricAlertsClient.Authorizer = *authorizer

	return &metricAlertsClient, nil
}

// CreateVMInsightsClientE gets a VM Insights client
func CreateVMInsightsClientE(subscriptionID string) (*insights.VMInsightsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	client := insights.NewVMInsightsClientWithBaseURI(baseURI, subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// CreateActivityLogAlertsClientE gets an Action Groups client in the specified Azure Subscription
func CreateActivityLogAlertsClientE(subscriptionID string) (*insights.ActivityLogAlertsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	// Get the Action Groups client
	client := insights.NewActivityLogAlertsClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// CreateDiagnosticsSettingsClientE returns a diagnostics settings client
func CreateDiagnosticsSettingsClientE(subscriptionID string) (*insights.DiagnosticSettingsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	client := insights.NewDiagnosticSettingsClientWithBaseURI(baseURI, subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// CreateNsgDefaultRulesClientE returns an NSG default (platform) rules client.
func CreateNsgDefaultRulesClientE(subscriptionID string) (*armnetwork.DefaultSecurityRulesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDefaultSecurityRulesClient(), nil
}

// CreateNsgCustomRulesClientE returns an NSG custom (user) rules client.
func CreateNsgCustomRulesClientE(subscriptionID string) (*armnetwork.SecurityRulesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSecurityRulesClient(), nil
}

// CreateNewNetworkInterfacesClientE returns a network interfaces client.
func CreateNewNetworkInterfacesClientE(subscriptionID string) (*armnetwork.InterfacesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewInterfacesClient(), nil
}

// CreateNewNetworkInterfaceIPConfigurationClientE returns a NIC IP configuration client.
func CreateNewNetworkInterfaceIPConfigurationClientE(subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewInterfaceIPConfigurationsClient(), nil
}

// CreatePublicIPAddressesClientE returns a public IP addresses client.
func CreatePublicIPAddressesClientE(subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewPublicIPAddressesClient(), nil
}

// CreateLoadBalancerClientE returns a load balancer client.
func CreateLoadBalancerClientE(subscriptionID string) (*armnetwork.LoadBalancersClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewLoadBalancersClient(), nil
}

// CreateLoadBalancerFrontendIPConfigClientE returns a load balancer frontend IP configuration client.
func CreateLoadBalancerFrontendIPConfigClientE(subscriptionID string) (*armnetwork.LoadBalancerFrontendIPConfigurationsClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewLoadBalancerFrontendIPConfigurationsClient(), nil
}

// CreateNewSubnetClientE returns a subnet client.
func CreateNewSubnetClientE(subscriptionID string) (*armnetwork.SubnetsClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSubnetsClient(), nil
}

// CreateNetworkManagementClientE returns a network management client.
func CreateNetworkManagementClientE(subscriptionID string) (*armnetwork.ManagementClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagementClient(), nil
}

// CreateNewVirtualNetworkClientE returns a virtual network client.
func CreateNewVirtualNetworkClientE(subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualNetworksClient(), nil
}

// CreateAppServiceClientE returns an App service client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	clientFactory, err := getArmAppServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWebAppsClient(), nil
}

// getArmAppServiceClientFactory gets an arm app service client factory
func getArmAppServiceClientFactory(subscriptionID string) (*armappservice.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armappservice.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// CreateContainerRegistryClientE returns an ACR client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateContainerRegistryClientE(subscriptionID string) (*containerregistry.RegistriesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return nil, err
	}

	// create client
	registryClient := containerregistry.NewRegistriesClientWithBaseURI(baseURI, subscriptionID)

	return &registryClient, nil
}

// CreateContainerInstanceClientE returns an ACI client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateContainerInstanceClientE(subscriptionID string) (*containerinstance.ContainerGroupsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return nil, err
	}

	// create client
	instanceClient := containerinstance.NewContainerGroupsClientWithBaseURI(baseURI, subscriptionID)

	return &instanceClient, nil
}

// CreateFrontDoorClientE returns an AFD client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateFrontDoorClientE(subscriptionID string) (*frontdoor.FrontDoorsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return nil, err
	}

	// create client
	client := frontdoor.NewFrontDoorsClientWithBaseURI(baseURI, subscriptionID)

	return &client, nil
}

// CreateFrontDoorFrontendEndpointClientE returns an AFD Frontend Endpoints client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateFrontDoorFrontendEndpointClientE(subscriptionID string) (*frontdoor.FrontendEndpointsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return nil, err
	}

	// create client
	client := frontdoor.NewFrontendEndpointsClientWithBaseURI(baseURI, subscriptionID)

	return &client, nil
}

// CreateSynapseWorkspaceClientE is a helper function that will setup a synapse workspace client.
func CreateSynapseWorkspaceClientE(subscriptionID string) (*armsynapse.WorkspacesClient, error) {
	clientFactory, err := getArmSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWorkspacesClient(), nil
}

// CreateSynapseSQLPoolClientE is a helper function that will setup a Synapse SQL pool client.
func CreateSynapseSQLPoolClientE(subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	clientFactory, err := getArmSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSQLPoolsClient(), nil
}

// CreateSynapseSqlPoolClientE is a helper function that will setup a Synapse SQL pool client.
//
// Deprecated: Use [CreateSynapseSQLPoolClientE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func CreateSynapseSqlPoolClientE(subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	return CreateSynapseSQLPoolClientE(subscriptionID)
}

// getArmSynapseClientFactory gets an arm synapse client factory
func getArmSynapseClientFactory(subscriptionID string) (*armsynapse.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armsynapse.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// CreateDataFactoriesClientE is a helper function that will setup a data factory client.
func CreateDataFactoriesClientE(subscriptionID string) (*armdatafactory.FactoriesClient, error) {
	clientFactory, err := getArmDataFactoryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFactoriesClient(), nil
}

// CreatePrivateDNSZonesClientE is a helper function that will setup a private DNS zone client.
func CreatePrivateDNSZonesClientE(subscriptionID string) (*privatedns.PrivateZonesClient, error) {
	// Validate Azure subscription ID
	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	// Create a private DNS zone client
	privateZonesClient := privatedns.NewPrivateZonesClientWithBaseURI(baseURI, subID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	privateZonesClient.Authorizer = *authorizer

	return &privateZonesClient, nil
}

// CreatePrivateDnsZonesClientE is a helper function that will setup a private DNS zone client.
//
// Deprecated: Use [CreatePrivateDNSZonesClientE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func CreatePrivateDnsZonesClientE(subscriptionID string) (*privatedns.PrivateZonesClient, error) {
	return CreatePrivateDNSZonesClientE(subscriptionID)
}

// CreateManagedEnvironmentsClientE creates a managed environments client for Azure Container Apps.
func CreateManagedEnvironmentsClientE(subscriptionID string) (*armappcontainers.ManagedEnvironmentsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewManagedEnvironmentsClient()

	return client, nil
}

// CreateResourceGroupClientV2E creates a v2 resource group client using the ARM SDK.
func CreateResourceGroupClientV2E(subscriptionID string) (*armresources.ResourceGroupsClient, error) {
	clientFactory, err := getArmResourcesClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewResourceGroupsClient(), nil
}

// CreateContainerAppsClientE creates a Container Apps client for Azure Container Apps.
func CreateContainerAppsClientE(subscriptionID string) (*armappcontainers.ContainerAppsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewContainerAppsClient()

	return client, nil
}

// CreateContainerAppJobsClientE creates a Container App Jobs client for Azure Container Apps.
func CreateContainerAppJobsClientE(subscriptionID string) (*armappcontainers.JobsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewJobsClient()

	return client, nil
}

// GetKeyVaultURISuffixE returns the proper KeyVault URI suffix for the configured Azure environment.
func GetKeyVaultURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()

	switch strings.ToUpper(envName) {
	case azurePublicCloud:
		return "vault.azure.net", nil
	case azureUSGovCloud:
		return "vault.usgovcloudapi.net", nil
	case azureChinaCloud:
		return "vault.azure.cn", nil
	default:
		return "", &UnknownEnvironmentError{EnvironmentName: envName}
	}
}

// GetStorageURISuffixE returns the proper storage URI suffix for the configured Azure environment.
func GetStorageURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()

	switch strings.ToUpper(envName) {
	case azurePublicCloud:
		return "core.windows.net", nil
	case azureUSGovCloud:
		return "core.usgovcloudapi.net", nil
	case azureChinaCloud:
		return "core.chinacloudapi.cn", nil
	default:
		return "", &UnknownEnvironmentError{EnvironmentName: envName}
	}
}

// getDefaultEnvironmentName returns either a configured Azure environment name, or the public default.
func getDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)
	if exists && len(envName) > 0 {
		return envName
	}

	return "AzurePublicCloud"
}

// getEnvironmentEndpointE returns the endpoint identified by the endpoint name parameter.
//
//nolint:unparam // endpointName kept as parameter for flexibility
func getEnvironmentEndpointE(endpointName string) (string, error) {
	envName := getDefaultEnvironmentName()

	env, err := autorestAzure.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}

	return getFieldValue(&env, endpointName), nil
}

// getFieldValue gets the field identified by the field parameter from the passed Environment struct
func getFieldValue(env *autorestAzure.Environment, field string) string {
	structValue := reflect.ValueOf(env)
	fieldVal := reflect.Indirect(structValue).FieldByName(field)

	return fieldVal.String()
}

// getBaseURI gets the base URI endpoint.
func getBaseURI() (string, error) {
	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return "", err
	}

	return baseURI, nil
}

// getArmResourcesClientFactory gets an arm resources client factory
func getArmResourcesClientFactory(subscriptionID string) (*armresources.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armresources.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// getArmAppContainersClientFactory gets an arm app containers client factory
func getArmAppContainersClientFactory(subscriptionID string) (*armappcontainers.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armappcontainers.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

// getArmDataFactoryClientFactory gets an arm data factory client factory
func getArmDataFactoryClientFactory(subscriptionID string) (*armdatafactory.ClientFactory, error) {
	targetSubscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	clientCloudConfig, err := getClientCloudConfig()
	if err != nil {
		return nil, err
	}

	cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
	if err != nil {
		return nil, err
	}

	return armdatafactory.NewClientFactory(targetSubscriptionID, cred, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	})
}

func getClientCloudConfig() (cloud.Configuration, error) {
	envName := getDefaultEnvironmentName()

	switch strings.ToUpper(envName) {
	case azureChinaCloud:
		return cloud.AzureChina, nil
	case azureUSGovCloud:
		return cloud.AzureGovernment, nil
	case azurePublicCloud:
		return cloud.AzurePublic, nil
	case azureStackCloud:
		adEndpoint := os.Getenv("AZURE_STACK_AD_ENDPOINT")
		rmEndpoint := os.Getenv("AZURE_STACK_RESOURCE_MANAGER_ENDPOINT")
		tokenAudience := os.Getenv("AZURE_STACK_TOKEN_AUDIENCE")

		if adEndpoint == "" || rmEndpoint == "" || tokenAudience == "" {
			return cloud.Configuration{},
				errors.New("AzureStackCloud requires AZURE_STACK_AD_ENDPOINT, " +
					"AZURE_STACK_RESOURCE_MANAGER_ENDPOINT, and AZURE_STACK_TOKEN_AUDIENCE " +
					"environment variables to be set")
		}

		return cloud.Configuration{
			ActiveDirectoryAuthorityHost: adEndpoint,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: tokenAudience,
					Endpoint: rmEndpoint,
				},
			},
		}, nil
	default:
		return cloud.Configuration{}, &UnknownEnvironmentError{EnvironmentName: envName}
	}
}
