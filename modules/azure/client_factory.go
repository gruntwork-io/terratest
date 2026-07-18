// Package azure allows users to interact with resources on the Microsoft Azure platform.
package azure

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
)

const (
	// AzureEnvironmentEnvName is the name of the Azure environment to use. Set to one of the following:
	//
	// "AzureChinaCloud":        ChinaCloud
	// "AzureGermanCloud":       GermanCloud
	// "AzurePublicCloud":       PublicCloud
	// "AzureUSGovernmentCloud": USGovernmentCloud
	// "AzureStackCloud":		 Azure stack
	AzureEnvironmentEnvName = "AZURE_ENVIRONMENT"

	// Azure environment name constants (upper-cased for case-insensitive switch matching).
	azurePublicCloud = "AZUREPUBLICCLOUD"
	azureUSGovCloud  = "AZUREUSGOVERNMENTCLOUD"
	azureChinaCloud  = "AZURECHINACLOUD"
	azureStackCloud  = "AZURESTACKCLOUD"
)

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

// CreateVirtualMachinesClientContextE returns a virtual machines client.
// The ctx parameter supports cancellation and timeouts.
func CreateVirtualMachinesClientContextE(_ context.Context, subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualMachinesClient(), nil
}

// getArmContainerServiceClientFactory creates an ARM container service client factory.
func getArmContainerServiceClientFactory(subscriptionID string) (*armcontainerservice.ClientFactory, error) {
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

	return armcontainerservice.NewClientFactory(targetSubscriptionID, cred, opts)
}

// CreateManagedClustersClientContextE returns a managed clusters client.
// The ctx parameter supports cancellation and timeouts.
func CreateManagedClustersClientContextE(_ context.Context, subscriptionID string) (*armcontainerservice.ManagedClustersClient, error) {
	clientFactory, err := getArmContainerServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedClustersClient(), nil
}

// CreateCosmosDBAccountClientContextE returns a Cosmos DB database accounts client.
// The ctx parameter supports cancellation and timeouts.
func CreateCosmosDBAccountClientContextE(_ context.Context, subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	clientFactory, err := getArmCosmosClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabaseAccountsClient(), nil
}

// CreateCosmosDBSQLClientContextE returns a Cosmos DB SQL resources client.
// The ctx parameter supports cancellation and timeouts.
func CreateCosmosDBSQLClientContextE(_ context.Context, subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
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

// CreateStorageAccountClientContextE creates a storage account client.
// The ctx parameter supports cancellation and timeouts.
func CreateStorageAccountClientContextE(_ context.Context, subscriptionID string) (*armstorage.AccountsClient, error) {
	clientFactory, err := getArmStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAccountsClient(), nil
}

// CreateStorageBlobContainerClientContextE creates a storage blob container client.
// The ctx parameter supports cancellation and timeouts.
func CreateStorageBlobContainerClientContextE(_ context.Context, subscriptionID string) (*armstorage.BlobContainersClient, error) {
	clientFactory, err := getArmStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBlobContainersClient(), nil
}

// CreateStorageFileSharesClientContextE creates a storage file shares client.
// The ctx parameter supports cancellation and timeouts.
func CreateStorageFileSharesClientContextE(_ context.Context, subscriptionID string) (*armstorage.FileSharesClient, error) {
	clientFactory, err := getArmStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFileSharesClient(), nil
}

// CreateServiceBusNamespacesClientContextE returns a service bus namespaces client.
// The ctx parameter supports cancellation and timeouts.
func CreateServiceBusNamespacesClientContextE(_ context.Context, subscriptionID string) (*armservicebus.NamespacesClient, error) {
	clientFactory, err := getArmServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewNamespacesClient(), nil
}

// CreateServiceBusTopicsClientContextE returns a service bus topics client.
// The ctx parameter supports cancellation and timeouts.
func CreateServiceBusTopicsClientContextE(_ context.Context, subscriptionID string) (*armservicebus.TopicsClient, error) {
	clientFactory, err := getArmServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewTopicsClient(), nil
}

// CreateServiceBusSubscriptionsClientContextE returns a service bus subscriptions client.
// The ctx parameter supports cancellation and timeouts.
func CreateServiceBusSubscriptionsClientContextE(_ context.Context, subscriptionID string) (*armservicebus.SubscriptionsClient, error) {
	clientFactory, err := getArmServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSubscriptionsClient(), nil
}

// CreateAvailabilitySetClientContextE creates a new Availability Set client.
// The ctx parameter supports cancellation and timeouts.
func CreateAvailabilitySetClientContextE(_ context.Context, subscriptionID string) (*armcompute.AvailabilitySetsClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAvailabilitySetsClient(), nil
}

// CreateResourceGroupClientContextE gets a resource group client in a subscription.
// The ctx parameter supports cancellation and timeouts.
func CreateResourceGroupClientContextE(_ context.Context, subscriptionID string) (*armresources.ResourceGroupsClient, error) {
	clientFactory, err := getArmResourcesClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewResourceGroupsClient(), nil
}

// CreateSQLServerClientContext is a helper function that will create and setup a sql server client.
// The ctx parameter supports cancellation and timeouts.
func CreateSQLServerClientContext(_ context.Context, subscriptionID string) (*armsql.ServersClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// CreateSQLManagedInstanceClientContext is a helper function that will create and setup a sql managed instance client.
// The ctx parameter supports cancellation and timeouts.
func CreateSQLManagedInstanceClientContext(_ context.Context, subscriptionID string) (*armsql.ManagedInstancesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedInstancesClient(), nil
}

// CreateSQLManagedDatabasesClientContext is a helper function that will create and setup a sql managed databases client.
// The ctx parameter supports cancellation and timeouts.
func CreateSQLManagedDatabasesClientContext(_ context.Context, subscriptionID string) (*armsql.ManagedDatabasesClient, error) {
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

// CreateDatabaseClientContext is a helper function that will create and setup a SQL DB client.
// The ctx parameter supports cancellation and timeouts.
func CreateDatabaseClientContext(_ context.Context, subscriptionID string) (*armsql.DatabasesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabasesClient(), nil
}

// CreateMySQLServerClientContextE is a helper function that will setup a mysql server client.
// The ctx parameter supports cancellation and timeouts.
func CreateMySQLServerClientContextE(_ context.Context, subscriptionID string) (*armmysql.ServersClient, error) {
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

// CreateDisksClientContextE returns a disks client.
// The ctx parameter supports cancellation and timeouts.
func CreateDisksClientContextE(_ context.Context, subscriptionID string) (*armcompute.DisksClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDisksClient(), nil
}

// CreateActionGroupClientContext creates an Action Groups client for Azure Monitor.
// The ctx parameter supports cancellation and timeouts.
func CreateActionGroupClientContext(_ context.Context, subscriptionID string) (*armmonitor.ActionGroupsClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
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

	return armmonitor.NewActionGroupsClient(subID, cred, opts)
}

// CreateVMInsightsClientContextE gets a VM Insights client.
// The ctx parameter supports cancellation and timeouts.
func CreateVMInsightsClientContextE(_ context.Context, subscriptionID string) (*armmonitor.VMInsightsClient, error) {

	_, err := getTargetAzureSubscription(subscriptionID)
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

	return armmonitor.NewVMInsightsClient(cred, opts)
}

// CreateActivityLogAlertsClientContextE gets an Activity Log Alerts client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func CreateActivityLogAlertsClientContextE(_ context.Context, subscriptionID string) (*armmonitor.ActivityLogAlertsClient, error) {

	subID, err := getTargetAzureSubscription(subscriptionID)
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

	return armmonitor.NewActivityLogAlertsClient(subID, cred, opts)
}

// CreateDiagnosticsSettingsClientContextE returns a diagnostics settings client.
// The ctx parameter supports cancellation and timeouts.
func CreateDiagnosticsSettingsClientContextE(_ context.Context, subscriptionID string) (*armmonitor.DiagnosticSettingsClient, error) {

	_, err := getTargetAzureSubscription(subscriptionID)
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

	return armmonitor.NewDiagnosticSettingsClient(cred, opts)
}

// CreateNsgDefaultRulesClientContextE returns an NSG default (platform) rules client.
// The ctx parameter supports cancellation and timeouts.
func CreateNsgDefaultRulesClientContextE(_ context.Context, subscriptionID string) (*armnetwork.DefaultSecurityRulesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDefaultSecurityRulesClient(), nil
}

// CreateNsgCustomRulesClientContextE returns an NSG custom (user) rules client.
// The ctx parameter supports cancellation and timeouts.
func CreateNsgCustomRulesClientContextE(_ context.Context, subscriptionID string) (*armnetwork.SecurityRulesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSecurityRulesClient(), nil
}

// CreateNetworkInterfacesClientContextE returns a network interfaces client.
// The ctx parameter supports cancellation and timeouts.
func CreateNetworkInterfacesClientContextE(_ context.Context, subscriptionID string) (*armnetwork.InterfacesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewInterfacesClient(), nil
}

// CreateNetworkInterfaceIPConfigurationClientContextE returns a NIC IP configuration client.
// The ctx parameter supports cancellation and timeouts.
func CreateNetworkInterfaceIPConfigurationClientContextE(_ context.Context, subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewInterfaceIPConfigurationsClient(), nil
}

// CreatePublicIPAddressesClientContextE returns a public IP addresses client.
// The ctx parameter supports cancellation and timeouts.
func CreatePublicIPAddressesClientContextE(_ context.Context, subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewPublicIPAddressesClient(), nil
}

// CreateLoadBalancerClientContextE returns a load balancer client.
// The ctx parameter supports cancellation and timeouts.
func CreateLoadBalancerClientContextE(_ context.Context, subscriptionID string) (*armnetwork.LoadBalancersClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewLoadBalancersClient(), nil
}

// CreateLoadBalancerFrontendIPConfigClientContextE returns a load balancer frontend IP configuration client.
// The ctx parameter supports cancellation and timeouts.
func CreateLoadBalancerFrontendIPConfigClientContextE(_ context.Context, subscriptionID string) (*armnetwork.LoadBalancerFrontendIPConfigurationsClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewLoadBalancerFrontendIPConfigurationsClient(), nil
}

// CreateSubnetClientContextE returns a subnet client.
// The ctx parameter supports cancellation and timeouts.
func CreateSubnetClientContextE(_ context.Context, subscriptionID string) (*armnetwork.SubnetsClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSubnetsClient(), nil
}

// CreateNetworkManagementClientContextE returns a network management client.
// The ctx parameter supports cancellation and timeouts.
func CreateNetworkManagementClientContextE(_ context.Context, subscriptionID string) (*armnetwork.ManagementClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagementClient(), nil
}

// CreateVirtualNetworkClientContextE returns a virtual network client.
// The ctx parameter supports cancellation and timeouts.
func CreateVirtualNetworkClientContextE(_ context.Context, subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	clientFactory, err := getArmNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualNetworksClient(), nil
}

// CreateAppServiceClientContextE returns an App Service client.
// The ctx parameter supports cancellation and timeouts.
func CreateAppServiceClientContextE(_ context.Context, subscriptionID string) (*armappservice.WebAppsClient, error) {
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

// getArmContainerRegistryClientFactory creates an ARM container registry client factory.
func getArmContainerRegistryClientFactory(subscriptionID string) (*armcontainerregistry.ClientFactory, error) {
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

	return armcontainerregistry.NewClientFactory(targetSubscriptionID, cred, opts)
}

// CreateContainerRegistryClientContextE returns an ACR registries client.
// The ctx parameter supports cancellation and timeouts.
func CreateContainerRegistryClientContextE(_ context.Context, subscriptionID string) (*armcontainerregistry.RegistriesClient, error) {
	clientFactory, err := getArmContainerRegistryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewRegistriesClient(), nil
}

// getArmContainerInstanceClientFactory creates an ARM container instance client factory.
func getArmContainerInstanceClientFactory(subscriptionID string) (*armcontainerinstance.ClientFactory, error) {
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

	return armcontainerinstance.NewClientFactory(targetSubscriptionID, cred, opts)
}

// CreateContainerInstanceClientContextE returns an ACI container groups client.
// The ctx parameter supports cancellation and timeouts.
func CreateContainerInstanceClientContextE(_ context.Context, subscriptionID string) (*armcontainerinstance.ContainerGroupsClient, error) {
	clientFactory, err := getArmContainerInstanceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewContainerGroupsClient(), nil
}

// CreateFrontDoorClientContextE returns a Front Door client.
// The ctx parameter supports cancellation and timeouts.
func CreateFrontDoorClientContextE(_ context.Context, subscriptionID string) (*armfrontdoor.FrontDoorsClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
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

	return armfrontdoor.NewFrontDoorsClient(subID, cred, opts)
}

// CreateFrontDoorFrontendEndpointClientContextE returns a Front Door Frontend Endpoints client.
// The ctx parameter supports cancellation and timeouts.
func CreateFrontDoorFrontendEndpointClientContextE(_ context.Context, subscriptionID string) (*armfrontdoor.FrontendEndpointsClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
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

	return armfrontdoor.NewFrontendEndpointsClient(subID, cred, opts)
}

// CreateSynapseWorkspaceClientContextE is a helper function that will setup a synapse workspace client.
// The ctx parameter supports cancellation and timeouts.
func CreateSynapseWorkspaceClientContextE(_ context.Context, subscriptionID string) (*armsynapse.WorkspacesClient, error) {
	clientFactory, err := getArmSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWorkspacesClient(), nil
}

// CreateSynapseSQLPoolClientContextE is a helper function that will setup a Synapse SQL pool client.
// The ctx parameter supports cancellation and timeouts.
func CreateSynapseSQLPoolClientContextE(_ context.Context, subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	clientFactory, err := getArmSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSQLPoolsClient(), nil
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

// CreateDataFactoriesClientContextE is a helper function that will setup a data factory client.
// The ctx parameter supports cancellation and timeouts.
func CreateDataFactoriesClientContextE(_ context.Context, subscriptionID string) (*armdatafactory.FactoriesClient, error) {
	clientFactory, err := getArmDataFactoryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFactoriesClient(), nil
}

// CreateManagedEnvironmentsClientContextE creates a managed environments client for Azure Container Apps.
// The ctx parameter supports cancellation and timeouts.
func CreateManagedEnvironmentsClientContextE(_ context.Context, subscriptionID string) (*armappcontainers.ManagedEnvironmentsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewManagedEnvironmentsClient()

	return client, nil
}

// CreateContainerAppsClientContextE creates a Container Apps client for Azure Container Apps.
// The ctx parameter supports cancellation and timeouts.
func CreateContainerAppsClientContextE(_ context.Context, subscriptionID string) (*armappcontainers.ContainerAppsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewContainerAppsClient()

	return client, nil
}

// CreateContainerAppJobsClientContextE creates a Container App Jobs client for Azure Container Apps.
// The ctx parameter supports cancellation and timeouts.
func CreateContainerAppJobsClientContextE(_ context.Context, subscriptionID string) (*armappcontainers.JobsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	client := clientFactory.NewJobsClient()

	return client, nil
}

// GetKeyVaultURISuffixContextE returns the proper KeyVault URI suffix for the configured Azure environment.
// The ctx parameter supports cancellation and timeouts.
func GetKeyVaultURISuffixContextE(_ context.Context) (string, error) {
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

// GetStorageURISuffixContextE returns the proper storage URI suffix for the configured Azure environment.
// The ctx parameter supports cancellation and timeouts.
func GetStorageURISuffixContextE(_ context.Context) (string, error) {
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
