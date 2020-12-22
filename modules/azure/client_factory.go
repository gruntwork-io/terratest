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

	"github.com/Azure/azure-sdk-for-go/profiles/preview/cosmos-db/mgmt/documentdb"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
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
)

// ClientType describes the type of client a module can create.
type ClientType int

// CreateSubscriptionsClientE returns a virtual machines client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateSubscriptionsClientE() (subscriptions.Client, error) {
	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return subscriptions.Client{}, err
	}

	// Create correct client based on type passed
	return subscriptions.NewClientWithBaseURI(baseURI), nil
}

// snippet-tag-start::client_factory_example.CreateClient

// CreateVirtualMachinesClientE returns a virtual machines client instance configured with the correct BaseURI depending on
// the Azure environment that is currently setup (or "Public", if none is setup).
func CreateVirtualMachinesClientE(subscriptionID string) (compute.VirtualMachinesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	// Create correct client based on type passed
	return compute.NewVirtualMachinesClientWithBaseURI(baseURI, subscriptionID), nil
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
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
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
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
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
	baseURI, err := getEnvironmentEndpointE(ResourceManagerEndpointName)
	if err != nil {
		return nil, err
	}

	// Create a CosmosDB client
	cosmosClient := documentdb.NewSQLResourcesClientWithBaseURI(baseURI, subscriptionID)

	return &cosmosClient, nil
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
