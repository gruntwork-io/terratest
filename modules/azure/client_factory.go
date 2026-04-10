/*

This file implements an Azure client factory that automatically handles setting up
credential and cloud configuration for sovereign cloud support. All clients use
the new Azure SDK (github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/).

*/

package azure

import (
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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/privatedns/armprivatedns"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
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
	clientCloudConfig, err := GetClientCloudConfig()
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
	clientCloudConfig, err := GetClientCloudConfig()
	if err != nil {
		return nil, err
	}

	return &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: clientCloudConfig,
		},
	}, nil
}

// GetDefaultEnvironmentName returns either a configured Azure environment name, or the public default.
func GetDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)
	if exists && len(envName) > 0 {
		return envName
	}

	return "AzurePublicCloud"
}

// GetClientCloudConfig returns the cloud.Configuration for the currently configured Azure environment.
func GetClientCloudConfig() (cloud.Configuration, error) {
	envName := GetDefaultEnvironmentName()

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

// GetKeyVaultURISuffixE returns the proper KeyVault URI suffix for the configured Azure environment.
func GetKeyVaultURISuffixE() (string, error) {
	envName := GetDefaultEnvironmentName()

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
	envName := GetDefaultEnvironmentName()

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

// ---- Private factory functions ----

func getARMSubscriptionsClientFactory() (*armsubscriptions.ClientFactory, error) {
	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armsubscriptions.NewClientFactory(cred, opts)
}

func getARMComputeClientFactory(subscriptionID string) (*armcompute.ClientFactory, error) {
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

func getARMNetworkClientFactory(subscriptionID string) (*armnetwork.ClientFactory, error) {
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

func getARMStorageClientFactory(subscriptionID string) (*armstorage.ClientFactory, error) {
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

func getARMContainerServiceClientFactory(subscriptionID string) (*armcontainerservice.ClientFactory, error) {
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

func getARMContainerRegistryClientFactory(subscriptionID string) (*armcontainerregistry.ClientFactory, error) {
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

func getARMContainerInstanceClientFactory(subscriptionID string) (*armcontainerinstance.ClientFactory, error) {
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

func getARMServiceBusClientFactory(subscriptionID string) (*armservicebus.ClientFactory, error) {
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

func getARMOperationalInsightsClientFactory(subscriptionID string) (*armoperationalinsights.ClientFactory, error) {
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

	return armoperationalinsights.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMRecoveryServicesClientFactory(subscriptionID string) (*armrecoveryservices.ClientFactory, error) {
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

	return armrecoveryservices.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMRecoveryServicesBackupClientFactory(subscriptionID string) (*armrecoveryservicesbackup.ClientFactory, error) {
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

	return armrecoveryservicesbackup.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMFrontDoorClientFactory(subscriptionID string) (*armfrontdoor.ClientFactory, error) {
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

	return armfrontdoor.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMCosmosClientFactory(subscriptionID string) (*armcosmos.ClientFactory, error) {
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

func getARMMonitorClientFactory(subscriptionID string) (*armmonitor.ClientFactory, error) {
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

	return armmonitor.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMPrivateDNSClientFactory(subscriptionID string) (*armprivatedns.ClientFactory, error) {
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

	return armprivatedns.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMKeyVaultClientFactory(subscriptionID string) (*armkeyvault.ClientFactory, error) {
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

	return armkeyvault.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMPostgreSQLClientFactory(subscriptionID string) (*armpostgresql.ClientFactory, error) {
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

	return armpostgresql.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMSQLClientFactory(subscriptionID string) (*armsql.ClientFactory, error) {
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

	return armsql.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMMySQLClientFactory(subscriptionID string) (*armmysql.ClientFactory, error) {
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

	return armmysql.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMAppServiceClientFactory(subscriptionID string) (*armappservice.ClientFactory, error) {
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

	return armappservice.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMSynapseClientFactory(subscriptionID string) (*armsynapse.ClientFactory, error) {
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

	return armsynapse.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMResourcesClientFactory(subscriptionID string) (*armresources.ClientFactory, error) {
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

	return armresources.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMAppContainersClientFactory(subscriptionID string) (*armappcontainers.ClientFactory, error) {
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

	return armappcontainers.NewClientFactory(targetSubscriptionID, cred, opts)
}

func getARMDataFactoryClientFactory(subscriptionID string) (*armdatafactory.ClientFactory, error) {
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

	return armdatafactory.NewClientFactory(targetSubscriptionID, cred, opts)
}

// ---- Public client creator functions ----

// CreateSubscriptionsClientE returns a subscriptions client.
func CreateSubscriptionsClientE() (*armsubscriptions.Client, error) {
	clientFactory, err := getARMSubscriptionsClientFactory()
	if err != nil {
		return nil, err
	}

	return clientFactory.NewClient(), nil
}

// CreateVirtualMachinesClientE returns a virtual machines client.
func CreateVirtualMachinesClientE(subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	clientFactory, err := getARMComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualMachinesClient(), nil
}

// CreateDisksClientE returns a disks client.
func CreateDisksClientE(subscriptionID string) (*armcompute.DisksClient, error) {
	clientFactory, err := getARMComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDisksClient(), nil
}

// CreateAvailabilitySetClientE creates a new Availability Set client.
func CreateAvailabilitySetClientE(subscriptionID string) (*armcompute.AvailabilitySetsClient, error) {
	clientFactory, err := getARMComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAvailabilitySetsClient(), nil
}

// CreateNsgDefaultRulesClientE returns an NSG default (platform) rules client.
func CreateNsgDefaultRulesClientE(subscriptionID string) (*armnetwork.DefaultSecurityRulesClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDefaultSecurityRulesClient(), nil
}

// CreateNsgCustomRulesClientE returns an NSG custom (user) rules client.
func CreateNsgCustomRulesClientE(subscriptionID string) (*armnetwork.SecurityRulesClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSecurityRulesClient(), nil
}

// CreateNewNetworkInterfacesClientE returns a network interfaces client.
func CreateNewNetworkInterfacesClientE(subscriptionID string) (*armnetwork.InterfacesClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewInterfacesClient(), nil
}

// CreateNewNetworkInterfaceIPConfigurationClientE returns a NIC IP configuration client.
func CreateNewNetworkInterfaceIPConfigurationClientE(subscriptionID string) (*armnetwork.InterfaceIPConfigurationsClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewInterfaceIPConfigurationsClient(), nil
}

// CreatePublicIPAddressesClientE returns a public IP addresses client.
func CreatePublicIPAddressesClientE(subscriptionID string) (*armnetwork.PublicIPAddressesClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewPublicIPAddressesClient(), nil
}

// CreateLoadBalancerClientE returns a load balancer client.
func CreateLoadBalancerClientE(subscriptionID string) (*armnetwork.LoadBalancersClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewLoadBalancersClient(), nil
}

// CreateLoadBalancerFrontendIPConfigClientE returns a load balancer frontend IP configuration client.
func CreateLoadBalancerFrontendIPConfigClientE(subscriptionID string) (*armnetwork.LoadBalancerFrontendIPConfigurationsClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewLoadBalancerFrontendIPConfigurationsClient(), nil
}

// CreateNewSubnetClientE returns a subnet client.
func CreateNewSubnetClientE(subscriptionID string) (*armnetwork.SubnetsClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSubnetsClient(), nil
}

// CreateNetworkManagementClientE returns a network management client (used for DNS name availability checks etc.).
func CreateNetworkManagementClientE(subscriptionID string) (*armnetwork.ManagementClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagementClient(), nil
}

// CreateNewVirtualNetworkClientE returns a virtual network client.
func CreateNewVirtualNetworkClientE(subscriptionID string) (*armnetwork.VirtualNetworksClient, error) {
	clientFactory, err := getARMNetworkClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualNetworksClient(), nil
}

// CreateStorageAccountClientE creates a storage account client.
func CreateStorageAccountClientE(subscriptionID string) (*armstorage.AccountsClient, error) {
	clientFactory, err := getARMStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAccountsClient(), nil
}

// CreateStorageBlobContainerClientE creates a storage blob container client.
func CreateStorageBlobContainerClientE(subscriptionID string) (*armstorage.BlobContainersClient, error) {
	clientFactory, err := getARMStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBlobContainersClient(), nil
}

// CreateStorageFileSharesClientE creates a storage file shares client.
func CreateStorageFileSharesClientE(subscriptionID string) (*armstorage.FileSharesClient, error) {
	clientFactory, err := getARMStorageClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFileSharesClient(), nil
}

// CreateManagedClustersClientE returns a managed clusters (AKS) client.
func CreateManagedClustersClientE(subscriptionID string) (*armcontainerservice.ManagedClustersClient, error) {
	clientFactory, err := getARMContainerServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedClustersClient(), nil
}

// CreateContainerRegistryClientE returns a container registry client.
func CreateContainerRegistryClientE(subscriptionID string) (*armcontainerregistry.RegistriesClient, error) {
	clientFactory, err := getARMContainerRegistryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewRegistriesClient(), nil
}

// CreateContainerInstanceClientE returns a container instance client.
func CreateContainerInstanceClientE(subscriptionID string) (*armcontainerinstance.ContainerGroupsClient, error) {
	clientFactory, err := getARMContainerInstanceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewContainerGroupsClient(), nil
}

// CreateServiceBusNamespacesClientE returns a service bus namespaces client.
func CreateServiceBusNamespacesClientE(subscriptionID string) (*armservicebus.NamespacesClient, error) {
	clientFactory, err := getARMServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewNamespacesClient(), nil
}

// CreateServiceBusTopicsClientE returns a service bus topics client.
func CreateServiceBusTopicsClientE(subscriptionID string) (*armservicebus.TopicsClient, error) {
	clientFactory, err := getARMServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewTopicsClient(), nil
}

// CreateServiceBusSubscriptionsClientE returns a service bus subscriptions client.
func CreateServiceBusSubscriptionsClientE(subscriptionID string) (*armservicebus.SubscriptionsClient, error) {
	clientFactory, err := getARMServiceBusClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSubscriptionsClient(), nil
}

// CreateLogAnalyticsWorkspacesClientE returns a log analytics workspaces client.
func CreateLogAnalyticsWorkspacesClientE(subscriptionID string) (*armoperationalinsights.WorkspacesClient, error) {
	clientFactory, err := getARMOperationalInsightsClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWorkspacesClient(), nil
}

// CreateRecoveryServicesVaultsClientE returns a recovery services vaults client.
func CreateRecoveryServicesVaultsClientE(subscriptionID string) (*armrecoveryservices.VaultsClient, error) {
	clientFactory, err := getARMRecoveryServicesClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVaultsClient(), nil
}

// CreateBackupPoliciesClientE returns a backup policies client.
func CreateBackupPoliciesClientE(subscriptionID string) (*armrecoveryservicesbackup.BackupPoliciesClient, error) {
	clientFactory, err := getARMRecoveryServicesBackupClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBackupPoliciesClient(), nil
}

// CreateBackupProtectedItemsClientE returns a backup protected items client.
func CreateBackupProtectedItemsClientE(subscriptionID string) (*armrecoveryservicesbackup.BackupProtectedItemsClient, error) {
	clientFactory, err := getARMRecoveryServicesBackupClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBackupProtectedItemsClient(), nil
}

// CreateFrontDoorClientE returns a front door client.
func CreateFrontDoorClientE(subscriptionID string) (*armfrontdoor.FrontDoorsClient, error) {
	clientFactory, err := getARMFrontDoorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFrontDoorsClient(), nil
}

// CreateFrontDoorFrontendEndpointClientE returns a front door frontend endpoint client.
func CreateFrontDoorFrontendEndpointClientE(subscriptionID string) (*armfrontdoor.FrontendEndpointsClient, error) {
	clientFactory, err := getARMFrontDoorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFrontendEndpointsClient(), nil
}

// CreateCosmosDBAccountClientE returns a Cosmos DB database accounts client.
func CreateCosmosDBAccountClientE(subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	clientFactory, err := getARMCosmosClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabaseAccountsClient(), nil
}

// CreateCosmosDBSQLClientE returns a Cosmos DB SQL resources client.
func CreateCosmosDBSQLClientE(subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
	clientFactory, err := getARMCosmosClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSQLResourcesClient(), nil
}

// CreateDiagnosticsSettingsClientE returns a diagnostics settings client.
func CreateDiagnosticsSettingsClientE(subscriptionID string) (*armmonitor.DiagnosticSettingsClient, error) {
	clientFactory, err := getARMMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDiagnosticSettingsClient(), nil
}

// CreateVMInsightsClientE returns a VM insights client.
func CreateVMInsightsClientE(subscriptionID string) (*armmonitor.VMInsightsClient, error) {
	clientFactory, err := getARMMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVMInsightsClient(), nil
}

// CreateActivityLogAlertsClientE returns an activity log alerts client.
func CreateActivityLogAlertsClientE(subscriptionID string) (*armmonitor.ActivityLogAlertsClient, error) {
	clientFactory, err := getARMMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewActivityLogAlertsClient(), nil
}

// CreateActionGroupClientE returns an action group client.
func CreateActionGroupClientE(subscriptionID string) (*armmonitor.ActionGroupsClient, error) {
	clientFactory, err := getARMMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewActionGroupsClient(), nil
}

// CreatePrivateDNSZonesClientE returns a private DNS zones client.
func CreatePrivateDNSZonesClientE(subscriptionID string) (*armprivatedns.PrivateZonesClient, error) {
	clientFactory, err := getARMPrivateDNSClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewPrivateZonesClient(), nil
}

// CreateResourceGroupClientE gets a resource group client.
func CreateResourceGroupClientE(subscriptionID string) (*armresources.ResourceGroupsClient, error) {
	clientFactory, err := getARMResourcesClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewResourceGroupsClient(), nil
}

// CreateSQLServerClient creates a SQL server client.
func CreateSQLServerClient(subscriptionID string) (*armsql.ServersClient, error) {
	clientFactory, err := getARMSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// CreateSQLMangedInstanceClient creates a SQL managed instance client.
func CreateSQLMangedInstanceClient(subscriptionID string) (*armsql.ManagedInstancesClient, error) {
	clientFactory, err := getARMSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedInstancesClient(), nil
}

// CreateSQLMangedDatabasesClient creates a SQL managed databases client.
func CreateSQLMangedDatabasesClient(subscriptionID string) (*armsql.ManagedDatabasesClient, error) {
	clientFactory, err := getARMSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedDatabasesClient(), nil
}

// CreateDatabaseClient creates a SQL databases client.
func CreateDatabaseClient(subscriptionID string) (*armsql.DatabasesClient, error) {
	clientFactory, err := getARMSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabasesClient(), nil
}

// CreateMySQLServerClientE creates a MySQL server client.
func CreateMySQLServerClientE(subscriptionID string) (*armmysql.ServersClient, error) {
	clientFactory, err := getARMMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// CreateAppServiceClientE returns an App Service web apps client.
func CreateAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	clientFactory, err := getARMAppServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWebAppsClient(), nil
}

// CreateSynapseWorkspaceClientE returns a Synapse workspaces client.
func CreateSynapseWorkspaceClientE(subscriptionID string) (*armsynapse.WorkspacesClient, error) {
	clientFactory, err := getARMSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWorkspacesClient(), nil
}

// CreateSynapseSQLPoolClientE returns a Synapse SQL pools client.
func CreateSynapseSQLPoolClientE(subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	clientFactory, err := getARMSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSQLPoolsClient(), nil
}

// CreateDataFactoriesClientE creates a data factory client.
func CreateDataFactoriesClientE(subscriptionID string) (*armdatafactory.FactoriesClient, error) {
	clientFactory, err := getARMDataFactoryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFactoriesClient(), nil
}

// CreateManagedEnvironmentsClientE creates a managed environments client for Azure Container Apps.
func CreateManagedEnvironmentsClientE(subscriptionID string) (*armappcontainers.ManagedEnvironmentsClient, error) {
	clientFactory, err := getARMAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedEnvironmentsClient(), nil
}

// CreateContainerAppsClientE creates a Container Apps client.
func CreateContainerAppsClientE(subscriptionID string) (*armappcontainers.ContainerAppsClient, error) {
	clientFactory, err := getARMAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewContainerAppsClient(), nil
}

// CreateContainerAppJobsClientE creates a Container App Jobs client.
func CreateContainerAppJobsClientE(subscriptionID string) (*armappcontainers.JobsClient, error) {
	clientFactory, err := getARMAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewJobsClient(), nil
}
