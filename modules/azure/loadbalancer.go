package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// LoadBalancerExists indicates whether the specified Load Balancer exists.
// This function would fail the test if there is an error.
func LoadBalancerExists(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := LoadBalancerExistsE(loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// LoadBalancerExistsE indicates whether the specified Load Balancer exists.
func LoadBalancerExistsE(loadBalancerName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetLoadBalancerE(loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetLoadBalancerFrontendIPConfigNames gets a list of the Frontend IP Configuration Names for the Load Balancer.
// This function would fail the test if there is an error.
func GetLoadBalancerFrontendIPConfigNames(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) []string {
	configName, err := GetLoadBalancerFrontendIPConfigNamesE(loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return configName
}

// GetLoadBalancerFrontendIPConfigNamesE ConfigNamesE gets a list of the Frontend IP Configuration Names for the Load Balancer.
func GetLoadBalancerFrontendIPConfigNamesE(loadBalancerName string, resourceGroupName string, subscriptionID string) ([]string, error) {
	lb, err := GetLoadBalancerE(loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Frontend IP Configurations
	lbProps := lb.LoadBalancerPropertiesFormat
	feConfigs := *lbProps.FrontendIPConfigurations
	if len(feConfigs) == 0 {
		// No Frontend IP Configuration present
		return nil, nil
	}

	// Get the names of the Frontend IP Configurations present
	configNames := make([]string, len(feConfigs))
	for i, config := range feConfigs {
		configNames[i] = *config.Name
	}

	return configNames, nil
}

// GetIPOfLoadBalancerFrontendIPConfig gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
// This function would fail the test if there is an error.
func GetIPOfLoadBalancerFrontendIPConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType) {
	ipAddress, ipType, err := GetIPOfLoadBalancerFrontendIPConfigE(feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return ipAddress, ipType
}

// GetIPOfLoadBalancerFrontendIPConfigE gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
func GetIPOfLoadBalancerFrontendIPConfigE(feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType, err1 error) {
	// Get the specified Load Balancer Frontend Config
	feConfig, err := GetLoadBalancerFrontendIPConfigE(feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", NoIP, err
	}

	// Get the Properties of the Frontend Configuration
	feProps := *feConfig.FrontendIPConfigurationPropertiesFormat

	// Check for the Public Type Frontend Config
	if feProps.PublicIPAddress != nil {
		// Get PublicIPAddress resource name from the Load Balancer Frontend Configuration
		pipName := GetNameFromResourceID(*feProps.PublicIPAddress.ID)

		// Get the Public IP of the PublicIPAddress
		ipValue, err := GetIPOfPublicIPAddressByNameE(pipName, resourceGroupName, subscriptionID)
		if err != nil {
			return "", NoIP, err
		}

		return ipValue, PublicIP, nil
	}

	// Return the Private IP as there are no other option available
	return *feProps.PrivateIPAddress, PrivateIP, nil

}

// GetLoadBalancerFrontendIPConfig gets the specified Load Balancer Frontend IP Configuration network resource.
// This function would fail the test if there is an error.
func GetLoadBalancerFrontendIPConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.FrontendIPConfiguration {
	lbFEConfig, err := GetLoadBalancerFrontendIPConfigE(feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return lbFEConfig
}

// GetLoadBalancerFrontendIPConfigE gets the specified Load Balancer Frontend IP Configuration network resource.
func GetLoadBalancerFrontendIPConfigE(feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.FrontendIPConfiguration, error) {
	// Validate Azure Resource Group Name
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetLoadBalancerFrontendIPConfigClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer Frontend IP Configuration
	lbc, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, feConfigName)
	if err != nil {
		return nil, err
	}

	return &lbc, nil
}

// GetLoadBalancerFrontendIPConfigClientE gets a new Load Balancer Frontend IP Configuration client in the specified Azure Subscription.
func GetLoadBalancerFrontendIPConfigClientE(subscriptionID string) (*network.LoadBalancerFrontendIPConfigurationsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer Frontend Configuration client
	client := network.NewLoadBalancerFrontendIPConfigurationsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}

// GetLoadBalancer gets a Load Balancer network resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
func GetLoadBalancer(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.LoadBalancer {
	lb, err := GetLoadBalancerE(loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return lb
}

// GetLoadBalancerE gets a Load Balancer network resource in the specified Azure Resource Group.
func GetLoadBalancerE(loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.LoadBalancer, error) {
	// Validate Azure Resource Group Name
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetLoadBalancerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer
	lb, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, "")
	if err != nil {
		return nil, err
	}

	return &lb, nil
}

// GetLoadBalancerClientE gets a new Load Balancer client in the specified Azure Subscription.
func GetLoadBalancerClientE(subscriptionID string) (*network.LoadBalancersClient, error) {
	// Get the Load Balancer client
	client, err := CreateLoadBalancerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return client, nil
}
