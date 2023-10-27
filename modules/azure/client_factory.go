/*

This file implements an Azure client factory that automatically handles setting up Base URI
values for sovereign cloud support. Note the list of clients below is not initially exhaustive;
rather, additional clients will be added as-needed.

*/

package azure

// snippet-tag-start::client_factory_example.imports

import (
	"os"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/frontdoor/mgmt/frontdoor"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/mysql/mgmt/mysql"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/sql/mgmt/sql"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/cosmos-db/mgmt/documentdb"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	kvmng "github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	sqlmi "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2020-12-01/synapse"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
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
)

// ClientType describes the type of client a module can create.
type ClientType int

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

// snippet-tag-start::client_factory_example.CreateClient

// CreateVirtualMachinesClientE returns a virtual machines client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateVirtualMachinesClientE(subscriptionID string) (*compute.VirtualMachinesClient, error) {
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

	// Create correct client based on type passed
	vmClient := compute.NewVirtualMachinesClientWithBaseURI(baseURI, subscriptionID)

	return &vmClient, nil
}

// snippet-tag-end::client_factory_example.CreateClient

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

// CreateCosmosDBAccountClientE is a helper function that will setup a CosmosDB account client with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateCosmosDBAccountClientE(subscriptionID string) (*documentdb.DatabaseAccountsClient, error) {
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

	// Create a CosmosDB client
	cosmosClient := documentdb.NewDatabaseAccountsClientWithBaseURI(baseURI, subscriptionID)

	return &cosmosClient, nil
}

// CreateCosmosDBSQLClientE is a helper function that will setup a CosmosDB SQL client with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateCosmosDBSQLClientE(subscriptionID string) (*documentdb.SQLResourcesClient, error) {
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

	// Create a CosmosDB client
	cosmosClient := documentdb.NewSQLResourcesClientWithBaseURI(baseURI, subscriptionID)

	return &cosmosClient, nil
}

// CreateKeyVaultManagementClientE is a helper function that will setup a key vault management client with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateKeyVaultManagementClientE(subscriptionID string) (*kvmng.VaultsClient, error) {
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

	//create keyvault management clinet
	vaultClient := kvmng.NewVaultsClientWithBaseURI(baseURI, subscriptionID)

	return &vaultClient, nil
}

// CreateStorageAccountClientE creates a storage account client.
func CreateStorageAccountClientE(subscriptionID string) (*storage.AccountsClient, error) {
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

	storageAccountClient := storage.NewAccountsClientWithBaseURI(baseURI, subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	storageAccountClient.Authorizer = *authorizer
	return &storageAccountClient, nil
}

// CreateStorageBlobContainerClientE creates a storage container client.
func CreateStorageBlobContainerClientE(subscriptionID string) (*storage.BlobContainersClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	blobContainerClient := storage.NewBlobContainersClientWithBaseURI(baseURI, subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}
	blobContainerClient.Authorizer = *authorizer
	return &blobContainerClient, nil
}

// CreateStorageFileSharesClientE creates a storage file share client.
func CreateStorageFileSharesClientE(subscriptionID string) (*storage.FileSharesClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Lookup environment URI
	baseURI, err := getBaseURI()
	if err != nil {
		return nil, err
	}

	fileShareClient := storage.NewFileSharesClientWithBaseURI(baseURI, subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}
	fileShareClient.Authorizer = *authorizer
	return &fileShareClient, nil
}

// CreateAvailabilitySetClientE creates a new Availability Set client in the specified Azure Subscription
func CreateAvailabilitySetClientE(subscriptionID string) (*compute.AvailabilitySetsClient, error) {
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

	// Get the Availability Set client
	client := compute.NewAvailabilitySetsClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
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
func CreateSQLServerClient(subscriptionID string) (*sql.ServersClient, error) {
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

	// Create a sql server client
	sqlClient := sql.NewServersClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	sqlClient.Authorizer = *authorizer

	return &sqlClient, nil
}

// CreateSQLMangedInstanceClient is a helper function that will create and setup a sql server client
func CreateSQLMangedInstanceClient(subscriptionID string) (*sqlmi.ManagedInstancesClient, error) {
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

	// Create a sql server client
	sqlmiClient := sqlmi.NewManagedInstancesClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	sqlmiClient.Authorizer = *authorizer

	return &sqlmiClient, nil
}

// CreateSQLMangedDatabasesClient is a helper function that will create and setup a sql server client
func CreateSQLMangedDatabasesClient(subscriptionID string) (*sqlmi.ManagedDatabasesClient, error) {
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

	// Create a sql server client
	sqlmidbClient := sqlmi.NewManagedDatabasesClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	sqlmidbClient.Authorizer = *authorizer

	return &sqlmidbClient, nil
}

// CreateDatabaseClient is a helper function that will create and setup a SQL DB client
func CreateDatabaseClient(subscriptionID string) (*sql.DatabasesClient, error) {
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

	// Create a sql DB client
	sqlDBClient := sql.NewDatabasesClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	sqlDBClient.Authorizer = *authorizer

	return &sqlDBClient, nil
}

// CreateMySQLServerClientE is a helper function that will setup a mysql server client.
func CreateMySQLServerClientE(subscriptionID string) (*mysql.ServersClient, error) {
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

	// Create a mysql server client
	mysqlClient := mysql.NewServersClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	mysqlClient.Authorizer = *authorizer

	return &mysqlClient, nil
}

// CreateDisksClientE returns a new Disks client in the specified Azure Subscription
func CreateDisksClientE(subscriptionID string) (*compute.DisksClient, error) {
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

	// Get the Disks client
	client := compute.NewDisksClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

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

// CreateNsgDefaultRulesClientE returns an NSG default (platform) rules client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateNsgDefaultRulesClientE(subscriptionID string) (*network.DefaultSecurityRulesClient, error) {
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

	// Create new client
	nsgClient := network.NewDefaultSecurityRulesClientWithBaseURI(baseURI, subscriptionID)
	return &nsgClient, nil
}

// CreateNsgCustomRulesClientE returns an NSG custom (user) rules client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateNsgCustomRulesClientE(subscriptionID string) (*network.SecurityRulesClient, error) {
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

	// Create new client
	nsgClient := network.NewSecurityRulesClientWithBaseURI(baseURI, subscriptionID)
	return &nsgClient, nil
}

// CreateNewNetworkInterfacesClientE returns an NIC client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateNewNetworkInterfacesClientE(subscriptionID string) (*network.InterfacesClient, error) {
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
	nicClient := network.NewInterfacesClientWithBaseURI(baseURI, subscriptionID)
	return &nicClient, nil
}

// CreateNewNetworkInterfaceIPConfigurationClientE returns an NIC IP configuration client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateNewNetworkInterfaceIPConfigurationClientE(subscriptionID string) (*network.InterfaceIPConfigurationsClient, error) {
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
	ipConfigClient := network.NewInterfaceIPConfigurationsClientWithBaseURI(baseURI, subscriptionID)
	return &ipConfigClient, nil
}

// CreatePublicIPAddressesClientE returns a public IP address client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreatePublicIPAddressesClientE(subscriptionID string) (*network.PublicIPAddressesClient, error) {
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

	// Create client
	client := network.NewPublicIPAddressesClientWithBaseURI(baseURI, subscriptionID)
	return &client, nil
}

// CreateLoadBalancerClientE returns a load balancer client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateLoadBalancerClientE(subscriptionID string) (*network.LoadBalancersClient, error) {
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

	//create LB client
	client := network.NewLoadBalancersClientWithBaseURI(baseURI, subscriptionID)
	return &client, nil
}

// CreateNewSubnetClientE returns a Subnet client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateNewSubnetClientE(subscriptionID string) (*network.SubnetsClient, error) {
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
	subnetClient := network.NewSubnetsClientWithBaseURI(baseURI, subscriptionID)
	return &subnetClient, nil
}

// CreateNewVirtualNetworkClientE returns a Virtual Network client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateNewVirtualNetworkClientE(subscriptionID string) (*network.VirtualNetworksClient, error) {
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
	vnetClient := network.NewVirtualNetworksClientWithBaseURI(baseURI, subscriptionID)
	return &vnetClient, nil
}

// CreateAppServiceClientE returns an App service client instance configured with the
// correct BaseURI depending on the Azure environment that is currently setup (or "Public", if none is setup).
func CreateAppServiceClientE(subscriptionID string) (*web.AppsClient, error) {

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
	appsClient := web.NewAppsClientWithBaseURI(baseURI, subscriptionID)
	return &appsClient, nil
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

// CreateSynapseWorkspaceClientE is a helper function that will setup a synapse client.
func CreateSynapseWorkspaceClientE(subscriptionID string) (*synapse.WorkspacesClient, error) {
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

	// Create a synapse client
	synapseWorkspaceClient := synapse.NewWorkspacesClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	synapseWorkspaceClient.Authorizer = *authorizer

	return &synapseWorkspaceClient, nil
}

// CreateSynapseSqlPoolClientE is a helper function that will setup a synapse client.
func CreateSynapseSqlPoolClientE(subscriptionID string) (*synapse.SQLPoolsClient, error) {
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

	// Create a synapse client
	synapseSqlPoolClient := synapse.NewSQLPoolsClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	synapseSqlPoolClient.Authorizer = *authorizer

	return &synapseSqlPoolClient, nil
}

// CreateDataFactoriesClientE is a helper function that will setup a synapse client.
func CreateDataFactoriesClientE(subscriptionID string) (*datafactory.FactoriesClient, error) {
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

	// Create a synapse client
	dataFactoryClient := datafactory.NewFactoriesClientWithBaseURI(baseURI, subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	dataFactoryClient.Authorizer = *authorizer

	return &dataFactoryClient, nil
}

// GetKeyVaultURISuffixE returns the proper KeyVault URI suffix for the configured Azure environment.
// This function would fail the test if there is an error.
func GetKeyVaultURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorestAzure.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.KeyVaultDNSSuffix, nil
}

// getDefaultEnvironmentName returns either a configured Azure environment name, or the public default
func getDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)

	if exists && len(envName) > 0 {
		return envName
	}

	return autorestAzure.PublicCloud.Name
}

// getEnvironmentEndpointE returns the endpoint identified by the endpoint name parameter.
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
