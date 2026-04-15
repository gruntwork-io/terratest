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

// CreateVirtualMachinesClientContextE returns a virtual machines client.
// The ctx parameter supports cancellation and timeouts.
func CreateVirtualMachinesClientContextE(_ context.Context, subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualMachinesClient(), nil
}

// CreateVirtualMachinesClientE returns a virtual machines client.
//
// Deprecated: Use [CreateVirtualMachinesClientContextE] instead.
func CreateVirtualMachinesClientE(subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	return CreateVirtualMachinesClientContextE(context.Background(), subscriptionID)
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

// CreateManagedClustersClientE returns a managed clusters client.
//
// Deprecated: Use [CreateManagedClustersClientContextE] instead.
func CreateManagedClustersClientE(subscriptionID string) (*armcontainerservice.ManagedClustersClient, error) {
	return CreateManagedClustersClientContextE(context.Background(), subscriptionID)
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

// CreateCosmosDBAccountClientE returns a Cosmos DB database accounts client.
//
// Deprecated: Use [CreateCosmosDBAccountClientContextE] instead.
func CreateCosmosDBAccountClientE(subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	return CreateCosmosDBAccountClientContextE(context.Background(), subscriptionID)
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

// CreateCosmosDBSQLClientE returns a Cosmos DB SQL resources client.
//
// Deprecated: Use [CreateCosmosDBSQLClientContextE] instead.
func CreateCosmosDBSQLClientE(subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
	return CreateCosmosDBSQLClientContextE(context.Background(), subscriptionID)
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

// CreateStorageAccountClientE creates a storage account client.
//
// Deprecated: Use [CreateStorageAccountClientContextE] instead.
func CreateStorageAccountClientE(subscriptionID string) (*armstorage.AccountsClient, error) {
	return CreateStorageAccountClientContextE(context.Background(), subscriptionID)
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

// CreateStorageBlobContainerClientE creates a storage blob container client.
//
// Deprecated: Use [CreateStorageBlobContainerClientContextE] instead.
func CreateStorageBlobContainerClientE(subscriptionID string) (*armstorage.BlobContainersClient, error) {
	return CreateStorageBlobContainerClientContextE(context.Background(), subscriptionID)
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

// CreateStorageFileSharesClientE creates a storage file shares client.
//
// Deprecated: Use [CreateStorageFileSharesClientContextE] instead.
func CreateStorageFileSharesClientE(subscriptionID string) (*armstorage.FileSharesClient, error) {
	return CreateStorageFileSharesClientContextE(context.Background(), subscriptionID)
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

// CreateServiceBusNamespacesClientE returns a service bus namespaces client.
//
// Deprecated: Use [CreateServiceBusNamespacesClientContextE] instead.
func CreateServiceBusNamespacesClientE(subscriptionID string) (*armservicebus.NamespacesClient, error) {
	return CreateServiceBusNamespacesClientContextE(context.Background(), subscriptionID)
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

// CreateServiceBusTopicsClientE returns a service bus topics client.
//
// Deprecated: Use [CreateServiceBusTopicsClientContextE] instead.
func CreateServiceBusTopicsClientE(subscriptionID string) (*armservicebus.TopicsClient, error) {
	return CreateServiceBusTopicsClientContextE(context.Background(), subscriptionID)
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

// CreateServiceBusSubscriptionsClientE returns a service bus subscriptions client.
//
// Deprecated: Use [CreateServiceBusSubscriptionsClientContextE] instead.
func CreateServiceBusSubscriptionsClientE(subscriptionID string) (*armservicebus.SubscriptionsClient, error) {
	return CreateServiceBusSubscriptionsClientContextE(context.Background(), subscriptionID)
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

// CreateAvailabilitySetClientE creates a new Availability Set client.
//
// Deprecated: Use [CreateAvailabilitySetClientContextE] instead.
func CreateAvailabilitySetClientE(subscriptionID string) (*armcompute.AvailabilitySetsClient, error) {
	return CreateAvailabilitySetClientContextE(context.Background(), subscriptionID)
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

// CreateResourceGroupClientE gets a resource group client in a subscription.
//
// Deprecated: Use [CreateResourceGroupClientContextE] instead.
func CreateResourceGroupClientE(subscriptionID string) (*armresources.ResourceGroupsClient, error) {
	return CreateResourceGroupClientContextE(context.Background(), subscriptionID)
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

// CreateSQLServerClient is a helper function that will create and setup a sql server client.
//
// Deprecated: Use [CreateSQLServerClientContext] instead.
func CreateSQLServerClient(subscriptionID string) (*armsql.ServersClient, error) {
	return CreateSQLServerClientContext(context.Background(), subscriptionID)
}

// CreateSQLMangedInstanceClientContext is a helper function that will create and setup a sql managed instance client.
// The ctx parameter supports cancellation and timeouts.
func CreateSQLMangedInstanceClientContext(_ context.Context, subscriptionID string) (*armsql.ManagedInstancesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedInstancesClient(), nil
}

// CreateSQLMangedInstanceClient is a helper function that will create and setup a sql managed instance client.
//
// Deprecated: Use [CreateSQLMangedInstanceClientContext] instead.
func CreateSQLMangedInstanceClient(subscriptionID string) (*armsql.ManagedInstancesClient, error) {
	return CreateSQLMangedInstanceClientContext(context.Background(), subscriptionID)
}

// CreateSQLMangedDatabasesClientContext is a helper function that will create and setup a sql managed databases client.
// The ctx parameter supports cancellation and timeouts.
func CreateSQLMangedDatabasesClientContext(_ context.Context, subscriptionID string) (*armsql.ManagedDatabasesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedDatabasesClient(), nil
}

// CreateSQLMangedDatabasesClient is a helper function that will create and setup a sql managed databases client.
//
// Deprecated: Use [CreateSQLMangedDatabasesClientContext] instead.
func CreateSQLMangedDatabasesClient(subscriptionID string) (*armsql.ManagedDatabasesClient, error) {
	return CreateSQLMangedDatabasesClientContext(context.Background(), subscriptionID)
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

// CreateDatabaseClient is a helper function that will create and setup a SQL DB client.
//
// Deprecated: Use [CreateDatabaseClientContext] instead.
func CreateDatabaseClient(subscriptionID string) (*armsql.DatabasesClient, error) {
	return CreateDatabaseClientContext(context.Background(), subscriptionID)
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

// CreateMySQLServerClientE is a helper function that will setup a mysql server client.
//
// Deprecated: Use [CreateMySQLServerClientContextE] instead.
func CreateMySQLServerClientE(subscriptionID string) (*armmysql.ServersClient, error) {
	return CreateMySQLServerClientContextE(context.Background(), subscriptionID)
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

// CreateDisksClientE returns a disks client.
//
// Deprecated: Use [CreateDisksClientContextE] instead.
func CreateDisksClientE(subscriptionID string) (*armcompute.DisksClient, error) {
	return CreateDisksClientContextE(context.Background(), subscriptionID)
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

// CreateActionGroupClient creates an Action Groups client for Azure Monitor.
//
// Deprecated: Use [CreateActionGroupClientContext] instead.
func CreateActionGroupClient(subscriptionID string) (*armmonitor.ActionGroupsClient, error) {
	return CreateActionGroupClientContext(context.Background(), subscriptionID)
}

// CreateVMInsightsClientContextE gets a VM Insights client.
// The ctx parameter supports cancellation and timeouts.
func CreateVMInsightsClientContextE(_ context.Context, subscriptionID string) (*armmonitor.VMInsightsClient, error) {
	// Validate Azure subscription ID
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

// CreateVMInsightsClientE gets a VM Insights client.
//
// Deprecated: Use [CreateVMInsightsClientContextE] instead.
func CreateVMInsightsClientE(subscriptionID string) (*armmonitor.VMInsightsClient, error) {
	return CreateVMInsightsClientContextE(context.Background(), subscriptionID)
}

// CreateActivityLogAlertsClientContextE gets an Activity Log Alerts client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func CreateActivityLogAlertsClientContextE(_ context.Context, subscriptionID string) (*armmonitor.ActivityLogAlertsClient, error) {
	// Validate Azure subscription ID
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

// CreateActivityLogAlertsClientE gets an Activity Log Alerts client in the specified Azure Subscription.
//
// Deprecated: Use [CreateActivityLogAlertsClientContextE] instead.
func CreateActivityLogAlertsClientE(subscriptionID string) (*armmonitor.ActivityLogAlertsClient, error) {
	return CreateActivityLogAlertsClientContextE(context.Background(), subscriptionID)
}

// CreateDiagnosticsSettingsClientContextE returns a diagnostics settings client.
// The ctx parameter supports cancellation and timeouts.
func CreateDiagnosticsSettingsClientContextE(_ context.Context, subscriptionID string) (*armmonitor.DiagnosticSettingsClient, error) {
	// Validate Azure subscription ID
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

	// DiagnosticSettingsClient does not take a subscriptionID — it operates on resourceURIs.
	return armmonitor.NewDiagnosticSettingsClient(cred, opts)
}

// CreateDiagnosticsSettingsClientE returns a diagnostics settings client.
//
// Deprecated: Use [CreateDiagnosticsSettingsClientContextE] instead.
func CreateDiagnosticsSettingsClientE(subscriptionID string) (*armmonitor.DiagnosticSettingsClient, error) {
	return CreateDiagnosticsSettingsClientContextE(context.Background(), subscriptionID)
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

// CreateNsgDefaultRulesClientE returns an NSG default (platform) rules client.
//
// Deprecated: Use [CreateNsgDefaultRulesClientContextE] instead.
func CreateNsgDefaultRulesClientE(subscriptionID string) (*armnetwork.DefaultSecurityRulesClient, error) {
	return CreateNsgDefaultRulesClientContextE(context.Background(), subscriptionID)
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

// CreateNsgCustomRulesClientE returns an NSG custom (user) rules client.
//
// Deprecated: Use [CreateNsgCustomRulesClientContextE] instead.
func CreateNsgCustomRulesClientE(subscriptionID string) (*armnetwork.SecurityRulesClient, error) {
	return CreateNsgCustomRulesClientContextE(context.Background(), subscriptionID)
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

// CreateNetworkInterfacesClientE returns a network interfaces client.
//
// Deprecated: Use [CreateNetworkInterfacesClientContextE] instead.
func CreateNetworkInterfacesClientE(subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return CreateNetworkInterfacesClientContextE(context.Background(), subscriptionID)
}

// CreateNewNetworkInterfacesClientContextE is an alias for backward compatibility.
//
// Deprecated: Use [CreateNetworkInterfacesClientContextE] instead.
func CreateNewNetworkInterfacesClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return CreateNetworkInterfacesClientContextE(ctx, subscriptionID)
}

// CreateNewNetworkInterfacesClientE is an alias for backward compatibility.
//
// Deprecated: Use [CreateNetworkInterfacesClientContextE] instead.
func CreateNewNetworkInterfacesClientE(subscriptionID string) (*armnetwork.InterfacesClient, error) {
	return CreateNetworkInterfacesClientContextE(context.Background(), subscriptionID)
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

// CreateNetworkInterfaceIPConfigurationClientE returns a NIC IP configuration client.
//
// Deprecated: Use [CreateNetworkInterfaceIPConfigurationClientContextE] instead.
func CreateNetworkInterfaceIPConfigurationClientE(subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	return CreateNetworkInterfaceIPConfigurationClientContextE(context.Background(), subscriptionID)
}

// CreateNewNetworkInterfaceIPConfigurationClientContextE is an alias for backward compatibility.
//
// Deprecated: Use [CreateNetworkInterfaceIPConfigurationClientContextE] instead.
func CreateNewNetworkInterfaceIPConfigurationClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	return CreateNetworkInterfaceIPConfigurationClientContextE(ctx, subscriptionID)
}

// CreateNewNetworkInterfaceIPConfigurationClientE is an alias for backward compatibility.
//
// Deprecated: Use [CreateNetworkInterfaceIPConfigurationClientContextE] instead.
func CreateNewNetworkInterfaceIPConfigurationClientE(subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	return CreateNetworkInterfaceIPConfigurationClientContextE(context.Background(), subscriptionID)
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

// CreatePublicIPAddressesClientE returns a public IP addresses client.
//
// Deprecated: Use [CreatePublicIPAddressesClientContextE] instead.
func CreatePublicIPAddressesClientE(subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	return CreatePublicIPAddressesClientContextE(context.Background(), subscriptionID)
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

// CreateLoadBalancerClientE returns a load balancer client.
//
// Deprecated: Use [CreateLoadBalancerClientContextE] instead.
func CreateLoadBalancerClientE(subscriptionID string) (*armnetwork.LoadBalancersClient, error) {
	return CreateLoadBalancerClientContextE(context.Background(), subscriptionID)
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

// CreateLoadBalancerFrontendIPConfigClientE returns a load balancer frontend IP configuration client.
//
// Deprecated: Use [CreateLoadBalancerFrontendIPConfigClientContextE] instead.
func CreateLoadBalancerFrontendIPConfigClientE(subscriptionID string) (*armnetwork.LoadBalancerFrontendIPConfigurationsClient, error) {
	return CreateLoadBalancerFrontendIPConfigClientContextE(context.Background(), subscriptionID)
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

// CreateSubnetClientE returns a subnet client.
//
// Deprecated: Use [CreateSubnetClientContextE] instead.
func CreateSubnetClientE(subscriptionID string) (*armnetwork.SubnetsClient, error) {
	return CreateSubnetClientContextE(context.Background(), subscriptionID)
}

// CreateNewSubnetClientContextE is an alias for backward compatibility.
//
// Deprecated: Use [CreateSubnetClientContextE] instead.
func CreateNewSubnetClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.SubnetsClient, error) {
	return CreateSubnetClientContextE(ctx, subscriptionID)
}

// CreateNewSubnetClientE is an alias for backward compatibility.
//
// Deprecated: Use [CreateSubnetClientContextE] instead.
func CreateNewSubnetClientE(subscriptionID string) (*armnetwork.SubnetsClient, error) {
	return CreateSubnetClientContextE(context.Background(), subscriptionID)
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

// CreateNetworkManagementClientE returns a network management client.
//
// Deprecated: Use [CreateNetworkManagementClientContextE] instead.
func CreateNetworkManagementClientE(subscriptionID string) (*armnetwork.ManagementClient, error) {
	return CreateNetworkManagementClientContextE(context.Background(), subscriptionID)
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

// CreateVirtualNetworkClientE returns a virtual network client.
//
// Deprecated: Use [CreateVirtualNetworkClientContextE] instead.
func CreateVirtualNetworkClientE(subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	return CreateVirtualNetworkClientContextE(context.Background(), subscriptionID)
}

// CreateNewVirtualNetworkClientContextE is an alias for backward compatibility.
//
// Deprecated: Use [CreateVirtualNetworkClientContextE] instead.
func CreateNewVirtualNetworkClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	return CreateVirtualNetworkClientContextE(ctx, subscriptionID)
}

// CreateNewVirtualNetworkClientE is an alias for backward compatibility.
//
// Deprecated: Use [CreateVirtualNetworkClientContextE] instead.
func CreateNewVirtualNetworkClientE(subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	return CreateVirtualNetworkClientContextE(context.Background(), subscriptionID)
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

// CreateAppServiceClientE returns an App Service client.
//
// Deprecated: Use [CreateAppServiceClientContextE] instead.
func CreateAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	return CreateAppServiceClientContextE(context.Background(), subscriptionID)
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

// CreateContainerRegistryClientE returns an ACR registries client.
//
// Deprecated: Use [CreateContainerRegistryClientContextE] instead.
func CreateContainerRegistryClientE(subscriptionID string) (*armcontainerregistry.RegistriesClient, error) {
	return CreateContainerRegistryClientContextE(context.Background(), subscriptionID)
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

// CreateContainerInstanceClientE returns an ACI container groups client.
//
// Deprecated: Use [CreateContainerInstanceClientContextE] instead.
func CreateContainerInstanceClientE(subscriptionID string) (*armcontainerinstance.ContainerGroupsClient, error) {
	return CreateContainerInstanceClientContextE(context.Background(), subscriptionID)
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

// CreateFrontDoorClientE returns a Front Door client.
//
// Deprecated: Use [CreateFrontDoorClientContextE] instead.
func CreateFrontDoorClientE(subscriptionID string) (*armfrontdoor.FrontDoorsClient, error) {
	return CreateFrontDoorClientContextE(context.Background(), subscriptionID)
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

// CreateFrontDoorFrontendEndpointClientE returns a Front Door Frontend Endpoints client.
//
// Deprecated: Use [CreateFrontDoorFrontendEndpointClientContextE] instead.
func CreateFrontDoorFrontendEndpointClientE(subscriptionID string) (*armfrontdoor.FrontendEndpointsClient, error) {
	return CreateFrontDoorFrontendEndpointClientContextE(context.Background(), subscriptionID)
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

// CreateSynapseWorkspaceClientE is a helper function that will setup a synapse workspace client.
//
// Deprecated: Use [CreateSynapseWorkspaceClientContextE] instead.
func CreateSynapseWorkspaceClientE(subscriptionID string) (*armsynapse.WorkspacesClient, error) {
	return CreateSynapseWorkspaceClientContextE(context.Background(), subscriptionID)
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

// CreateSynapseSQLPoolClientE is a helper function that will setup a Synapse SQL pool client.
//
// Deprecated: Use [CreateSynapseSQLPoolClientContextE] instead.
func CreateSynapseSQLPoolClientE(subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	return CreateSynapseSQLPoolClientContextE(context.Background(), subscriptionID)
}

// CreateSynapseSqlPoolClientContextE is a helper function that will setup a Synapse SQL pool client.
// The ctx parameter supports cancellation and timeouts.
//
// Deprecated: Use [CreateSynapseSQLPoolClientContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func CreateSynapseSqlPoolClientContextE(ctx context.Context, subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	return CreateSynapseSQLPoolClientContextE(ctx, subscriptionID)
}

// CreateSynapseSqlPoolClientE is a helper function that will setup a Synapse SQL pool client.
//
// Deprecated: Use [CreateSynapseSQLPoolClientContextE] instead.
//
//nolint:staticcheck,revive // preserving existing function name
func CreateSynapseSqlPoolClientE(subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	return CreateSynapseSQLPoolClientContextE(context.Background(), subscriptionID)
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

// CreateDataFactoriesClientE is a helper function that will setup a data factory client.
//
// Deprecated: Use [CreateDataFactoriesClientContextE] instead.
func CreateDataFactoriesClientE(subscriptionID string) (*armdatafactory.FactoriesClient, error) {
	return CreateDataFactoriesClientContextE(context.Background(), subscriptionID)
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

// CreateManagedEnvironmentsClientE creates a managed environments client for Azure Container Apps.
//
// Deprecated: Use [CreateManagedEnvironmentsClientContextE] instead.
func CreateManagedEnvironmentsClientE(subscriptionID string) (*armappcontainers.ManagedEnvironmentsClient, error) {
	return CreateManagedEnvironmentsClientContextE(context.Background(), subscriptionID)
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

// CreateContainerAppsClientE creates a Container Apps client for Azure Container Apps.
//
// Deprecated: Use [CreateContainerAppsClientContextE] instead.
func CreateContainerAppsClientE(subscriptionID string) (*armappcontainers.ContainerAppsClient, error) {
	return CreateContainerAppsClientContextE(context.Background(), subscriptionID)
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

// CreateContainerAppJobsClientE creates a Container App Jobs client for Azure Container Apps.
//
// Deprecated: Use [CreateContainerAppJobsClientContextE] instead.
func CreateContainerAppJobsClientE(subscriptionID string) (*armappcontainers.JobsClient, error) {
	return CreateContainerAppJobsClientContextE(context.Background(), subscriptionID)
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

// GetKeyVaultURISuffixE returns the proper KeyVault URI suffix for the configured Azure environment.
//
// Deprecated: Use [GetKeyVaultURISuffixContextE] instead.
func GetKeyVaultURISuffixE() (string, error) {
	return GetKeyVaultURISuffixContextE(context.Background())
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

// GetStorageURISuffixE returns the proper storage URI suffix for the configured Azure environment.
//
// Deprecated: Use [GetStorageURISuffixContextE] instead.
func GetStorageURISuffixE() (string, error) {
	return GetStorageURISuffixContextE(context.Background())
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
