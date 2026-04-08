/*

This file implements an Azure client factory that automatically handles setting up
credential and cloud configuration for sovereign cloud support. All clients use
the new Azure SDK (github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/).

*/

package azure

import (
	"errors"
	"fmt"
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

// getDefaultEnvironmentName returns either a configured Azure environment name, or the public default.
func getDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)
	if exists && len(envName) > 0 {
		return envName
	}

	return "AzurePublicCloud"
}

func getClientCloudConfig() (cloud.Configuration, error) {
	envName := getDefaultEnvironmentName()

	switch strings.ToUpper(envName) {
	case "AZURECHINACLOUD":
		return cloud.AzureChina, nil
	case "AZUREUSGOVERNMENTCLOUD":
		return cloud.AzureGovernment, nil
	case "AZUREPUBLICCLOUD":
		return cloud.AzurePublic, nil
	case "AZURESTACKCLOUD":
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
		return cloud.Configuration{},
			fmt.Errorf("no cloud environment matching the name: %s. "+
				"Available values are: "+
				"AzurePublicCloud (default), "+
				"AzureUSGovernmentCloud, "+
				"AzureChinaCloud or "+
				"AzureStackCloud",
				envName)
	}
}

// GetKeyVaultURISuffixE returns the proper KeyVault URI suffix for the configured Azure environment.
func GetKeyVaultURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()

	switch strings.ToUpper(envName) {
	case "AZUREPUBLICCLOUD":
		return "vault.azure.net", nil
	case "AZUREUSGOVERNMENTCLOUD":
		return "vault.usgovcloudapi.net", nil
	case "AZURECHINACLOUD":
		return "vault.azure.cn", nil
	default:
		return "", fmt.Errorf("KeyVault URI suffix not known for environment: %s", envName)
	}
}

// GetStorageURISuffixE returns the proper storage URI suffix for the configured Azure environment.
func GetStorageURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()

	switch strings.ToUpper(envName) {
	case "AZUREPUBLICCLOUD":
		return "core.windows.net", nil
	case "AZUREUSGOVERNMENTCLOUD":
		return "core.usgovcloudapi.net", nil
	case "AZURECHINACLOUD":
		return "core.chinacloudapi.cn", nil
	default:
		return "", fmt.Errorf("storage URI suffix not known for environment: %s", envName)
	}
}

// ---- Private factory functions ----

func getArmSubscriptionsClientFactory() (*armsubscriptions.ClientFactory, error) {
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

func getArmOperationalInsightsClientFactory(subscriptionID string) (*armoperationalinsights.ClientFactory, error) {
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

func getArmRecoveryServicesClientFactory(subscriptionID string) (*armrecoveryservices.ClientFactory, error) {
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

func getArmRecoveryServicesBackupClientFactory(subscriptionID string) (*armrecoveryservicesbackup.ClientFactory, error) {
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

func getArmFrontDoorClientFactory(subscriptionID string) (*armfrontdoor.ClientFactory, error) {
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

func getArmMonitorClientFactory(subscriptionID string) (*armmonitor.ClientFactory, error) {
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

func getArmPrivateDNSClientFactory(subscriptionID string) (*armprivatedns.ClientFactory, error) {
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

func getArmKeyVaultClientFactory(subscriptionID string) (*armkeyvault.ClientFactory, error) {
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

func getArmPostgreSQLClientFactory(subscriptionID string) (*armpostgresql.ClientFactory, error) {
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

func getArmSQLClientFactory(subscriptionID string) (*armsql.ClientFactory, error) {
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

func getArmMySQLClientFactory(subscriptionID string) (*armmysql.ClientFactory, error) {
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

func getArmAppServiceClientFactory(subscriptionID string) (*armappservice.ClientFactory, error) {
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

func getArmSynapseClientFactory(subscriptionID string) (*armsynapse.ClientFactory, error) {
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

func getArmResourcesClientFactory(subscriptionID string) (*armresources.ClientFactory, error) {
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

func getArmAppContainersClientFactory(subscriptionID string) (*armappcontainers.ClientFactory, error) {
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

func getArmDataFactoryClientFactory(subscriptionID string) (*armdatafactory.ClientFactory, error) {
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
	clientFactory, err := getArmSubscriptionsClientFactory()
	if err != nil {
		return nil, err
	}

	return clientFactory.NewClient(), nil
}

// CreateVirtualMachinesClientE returns a virtual machines client.
func CreateVirtualMachinesClientE(subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVirtualMachinesClient(), nil
}

// CreateDisksClientE returns a disks client.
func CreateDisksClientE(subscriptionID string) (*armcompute.DisksClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDisksClient(), nil
}

// CreateAvailabilitySetClientE creates a new Availability Set client.
func CreateAvailabilitySetClientE(subscriptionID string) (*armcompute.AvailabilitySetsClient, error) {
	clientFactory, err := getArmComputeClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewAvailabilitySetsClient(), nil
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

// CreateNetworkManagementClientE returns a network management client (used for DNS name availability checks etc.).
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

// CreateManagedClustersClientE returns a managed clusters (AKS) client.
func CreateManagedClustersClientE(subscriptionID string) (*armcontainerservice.ManagedClustersClient, error) {
	clientFactory, err := getArmContainerServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedClustersClient(), nil
}

// CreateContainerRegistryClientE returns a container registry client.
func CreateContainerRegistryClientE(subscriptionID string) (*armcontainerregistry.RegistriesClient, error) {
	clientFactory, err := getArmContainerRegistryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewRegistriesClient(), nil
}

// CreateContainerInstanceClientE returns a container instance client.
func CreateContainerInstanceClientE(subscriptionID string) (*armcontainerinstance.ContainerGroupsClient, error) {
	clientFactory, err := getArmContainerInstanceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewContainerGroupsClient(), nil
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

// CreateLogAnalyticsWorkspacesClientE returns a log analytics workspaces client.
func CreateLogAnalyticsWorkspacesClientE(subscriptionID string) (*armoperationalinsights.WorkspacesClient, error) {
	clientFactory, err := getArmOperationalInsightsClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWorkspacesClient(), nil
}

// CreateRecoveryServicesVaultsClientE returns a recovery services vaults client.
func CreateRecoveryServicesVaultsClientE(subscriptionID string) (*armrecoveryservices.VaultsClient, error) {
	clientFactory, err := getArmRecoveryServicesClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVaultsClient(), nil
}

// CreateBackupPoliciesClientE returns a backup policies client.
func CreateBackupPoliciesClientE(subscriptionID string) (*armrecoveryservicesbackup.BackupPoliciesClient, error) {
	clientFactory, err := getArmRecoveryServicesBackupClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBackupPoliciesClient(), nil
}

// CreateBackupProtectedItemsClientE returns a backup protected items client.
func CreateBackupProtectedItemsClientE(subscriptionID string) (*armrecoveryservicesbackup.BackupProtectedItemsClient, error) {
	clientFactory, err := getArmRecoveryServicesBackupClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewBackupProtectedItemsClient(), nil
}

// CreateFrontDoorClientE returns a front door client.
func CreateFrontDoorClientE(subscriptionID string) (*armfrontdoor.FrontDoorsClient, error) {
	clientFactory, err := getArmFrontDoorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFrontDoorsClient(), nil
}

// CreateFrontDoorFrontendEndpointClientE returns a front door frontend endpoint client.
func CreateFrontDoorFrontendEndpointClientE(subscriptionID string) (*armfrontdoor.FrontendEndpointsClient, error) {
	clientFactory, err := getArmFrontDoorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFrontendEndpointsClient(), nil
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

// CreateDiagnosticsSettingsClientE returns a diagnostics settings client.
func CreateDiagnosticsSettingsClientE(subscriptionID string) (*armmonitor.DiagnosticSettingsClient, error) {
	clientFactory, err := getArmMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDiagnosticSettingsClient(), nil
}

// CreateVMInsightsClientE returns a VM insights client.
func CreateVMInsightsClientE(subscriptionID string) (*armmonitor.VMInsightsClient, error) {
	clientFactory, err := getArmMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewVMInsightsClient(), nil
}

// CreateActivityLogAlertsClientE returns an activity log alerts client.
func CreateActivityLogAlertsClientE(subscriptionID string) (*armmonitor.ActivityLogAlertsClient, error) {
	clientFactory, err := getArmMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewActivityLogAlertsClient(), nil
}

// CreateActionGroupClientE returns an action group client.
func CreateActionGroupClientE(subscriptionID string) (*armmonitor.ActionGroupsClient, error) {
	clientFactory, err := getArmMonitorClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewActionGroupsClient(), nil
}

// CreatePrivateDNSZonesClientE returns a private DNS zones client.
func CreatePrivateDNSZonesClientE(subscriptionID string) (*armprivatedns.PrivateZonesClient, error) {
	clientFactory, err := getArmPrivateDNSClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewPrivateZonesClient(), nil
}

// CreateResourceGroupClientE gets a resource group client.
func CreateResourceGroupClientE(subscriptionID string) (*armresources.ResourceGroupsClient, error) {
	clientFactory, err := getArmResourcesClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewResourceGroupsClient(), nil
}

// CreateSQLServerClient creates a SQL server client.
func CreateSQLServerClient(subscriptionID string) (*armsql.ServersClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// CreateSQLMangedInstanceClient creates a SQL managed instance client.
func CreateSQLMangedInstanceClient(subscriptionID string) (*armsql.ManagedInstancesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedInstancesClient(), nil
}

// CreateSQLMangedDatabasesClient creates a SQL managed databases client.
func CreateSQLMangedDatabasesClient(subscriptionID string) (*armsql.ManagedDatabasesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedDatabasesClient(), nil
}

// CreateDatabaseClient creates a SQL databases client.
func CreateDatabaseClient(subscriptionID string) (*armsql.DatabasesClient, error) {
	clientFactory, err := getArmSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabasesClient(), nil
}

// CreateMySQLServerClientE creates a MySQL server client.
func CreateMySQLServerClientE(subscriptionID string) (*armmysql.ServersClient, error) {
	clientFactory, err := getArmMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// CreateAppServiceClientE returns an App Service web apps client.
func CreateAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	clientFactory, err := getArmAppServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWebAppsClient(), nil
}

// CreateSynapseWorkspaceClientE returns a Synapse workspaces client.
func CreateSynapseWorkspaceClientE(subscriptionID string) (*armsynapse.WorkspacesClient, error) {
	clientFactory, err := getArmSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWorkspacesClient(), nil
}

// CreateSynapseSQLPoolClientE returns a Synapse SQL pools client.
func CreateSynapseSQLPoolClientE(subscriptionID string) (*armsynapse.SQLPoolsClient, error) {
	clientFactory, err := getArmSynapseClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewSQLPoolsClient(), nil
}

// CreateDataFactoriesClientE creates a data factory client.
func CreateDataFactoriesClientE(subscriptionID string) (*armdatafactory.FactoriesClient, error) {
	clientFactory, err := getArmDataFactoryClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewFactoriesClient(), nil
}

// CreateManagedEnvironmentsClientE creates a managed environments client for Azure Container Apps.
func CreateManagedEnvironmentsClientE(subscriptionID string) (*armappcontainers.ManagedEnvironmentsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewManagedEnvironmentsClient(), nil
}

// CreateContainerAppsClientE creates a Container Apps client.
func CreateContainerAppsClientE(subscriptionID string) (*armappcontainers.ContainerAppsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewContainerAppsClient(), nil
}

// CreateContainerAppJobsClientE creates a Container App Jobs client.
func CreateContainerAppJobsClientE(subscriptionID string) (*armappcontainers.JobsClient, error) {
	clientFactory, err := getArmAppContainersClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewJobsClient(), nil
}
