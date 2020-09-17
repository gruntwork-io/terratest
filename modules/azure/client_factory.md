## Azure SDK Client Factory 

This documentation provides and overview of the `client_factory.go` module, targeted use cases, and behaviors. 

## Multi-cloud environment support

In order to better support Azure customers, we needed to support both Public and sovereign cloud environments (includes USGovernment, Germany, China, and Azure Stack).  To do this, we need to use the correct base URI's for the Azure management plane when using the REST API (either directly or via SDK). For the Go SDK, this can be accomplished by using the `WithBaseURI` suffixed calls when creating service clients.

For example, when using the `VirtualMachinesClient`, a developer would normally write code like this:

```golang
import (
    "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
)

func SomeMethod() {
    subscriptionID := "your subscription ID"

    // Create a VM client and return
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)

    // Use client / etc
}
```

However, this code will not work in non-Public cloud environments, such as USGovCloud, Germany, China, or on Azure Stack.  Instead, we need to use an alternative method (provided in the Go SDK) to get our clients:

```golang
import (
    "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
)

func SomeMethod() {
    subscriptionID := "your subscription ID"
    baseURI := "management.azure.com"

    // Create a VM client and return
	vmClient := compute.NewVirtualMachinesClientWithBaseURI(baseURI, subscriptionID)

    // Use client / etc
}
```

Now we have code that can be used with any cloud environment just by changing the base URI we pass to the clients (Public is shown in above example). But where do we get these URI's?

## Environment URI Lookups

Rather than hardcode such URI's, we should look these up. Fortunately, the Autorest-GO library (used by the Go SDK) provides such functionality. The `client_factory` module makes use of the Autorest `EnvironmentFromName(envName string)` function to return the appropriate structure. 

This Environment structure has the following fields configured for each cloud environment:

```golang
ManagementPortalURL          
PublishSettingsURL           
ServiceManagementEndpoint    
ResourceManagerEndpoint      
ActiveDirectoryEndpoint      
GalleryEndpoint              
KeyVaultEndpoint             
GraphEndpoint                
ServiceBusEndpoint           
BatchManagementEndpoint      
StorageEndpointSuffix        
SQLDatabaseDNSSuffix         
TrafficManagerDNSSuffix      
KeyVaultDNSSuffix            
ServiceBusEndpointSuffix     
ServiceManagementVMDNSSuffix 
ResourceManagerVMDNSSuffix   
ContainerRegistryDNSSuffix   
CosmosDBDNSSuffix            
TokenAudience                
APIManagementHostNameSuffix  
SynapseEndpointSuffix        
ResourceIdentifiers          
```

Using these URI's, `client_factory` is able to instantiate and return properly configured SDK clients without module developers having to duplicate this code.

## Configuration and Defaults

To configure different cloud environments, we will use the same `AZURE_ENVIRONMENT` environment variable that the Go SDK uses. This can be set to one of the following values:

|Value                      |Cloud Environment  |
|---------------------------|-------------------|
|"AzureChinaCloud"          |ChinaCloud         |
|"AzureGermanCloud"         |GermanCloud        |
|"AzurePublicCloud"         |PublicCloud        |
|"AzureUSGovernmentCloud"   |USGovernmentCloud  |
|"AzureStackCloud"          |Azure stack        |


Note that when using the "AzureStackCloud" setting, you must also set the `AZURE_ENVIRONMENT_FILEPATH` variable to point to a JSON file containing your Azure Stack URI deatils.

>NOTE: The default behavior of the `client_factory` is to use the AzurePublicCloud environment. This requires no work >from the developer to configure, and ensures consistent behavior with the current SDK code.

## Usage Patterns

Modules authors will interact with the `client_factory` through the `NewClientFactory` method on the `azure` package as shown below:

```golang
	// Create a VM client
	factory := NewClientFactory()
	vmClient, err := factory.GetVirtualMachinesClientE(subscriptionID)
	if err != nil {
		return nil, err
    }
```

The `ClientFactory` interface currently exposes the following methods - these will be built out more as more modules are developed and integrated:

* `GetVirtualMachinesClientE(subscriptionID string)`: returns a configured compute client, setup for proper cloud environment use
* `GetSubscriptionClientE()`: returns a configured compute client, setup for proper cloud environment use
* `GetManagedClustersClientE(subscriptionID string)`: returns a configured compute client, setup for proper cloud environment use.