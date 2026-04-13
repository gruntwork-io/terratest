package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// LoadBalancerExistsContext indicates whether the specified Load Balancer exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func LoadBalancerExistsContext(t testing.TestingT, ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := LoadBalancerExistsContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// LoadBalancerExistsContextE indicates whether the specified Load Balancer exists.
// The ctx parameter supports cancellation and timeouts.
func LoadBalancerExistsContextE(ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetLoadBalancerContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetLoadBalancerFrontendIPConfigNamesContext gets a list of the Frontend IP Configuration Names for the Load Balancer.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigNamesContext(t testing.TestingT, ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) []string {
	t.Helper()

	configName, err := GetLoadBalancerFrontendIPConfigNamesContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return configName
}

// GetLoadBalancerFrontendIPConfigNamesContextE gets a list of the Frontend IP Configuration Names for the Load Balancer.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigNamesContextE(ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) ([]string, error) {
	lb, err := GetLoadBalancerContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Frontend IP Configurations
	feConfigs := lb.Properties.FrontendIPConfigurations

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

// GetIPOfLoadBalancerFrontendIPConfigContext gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIPOfLoadBalancerFrontendIPConfigContext(t testing.TestingT, ctx context.Context, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType) {
	t.Helper()

	ipAddress, ipType, err := GetIPOfLoadBalancerFrontendIPConfigContextE(ctx, feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return ipAddress, ipType
}

// GetIPOfLoadBalancerFrontendIPConfigContextE gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
// The ctx parameter supports cancellation and timeouts.
func GetIPOfLoadBalancerFrontendIPConfigContextE(ctx context.Context, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType, err1 error) {
	// Get the specified Load Balancer Frontend Config
	feConfig, err := GetLoadBalancerFrontendIPConfigContextE(ctx, feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", NoIP, err
	}

	// Get the Properties of the Frontend Configuration
	feProps := feConfig.Properties

	// Check for the Public Type Frontend Config
	if feProps.PublicIPAddress != nil {
		// Get PublicIPAddress resource name from the Load Balancer Frontend Configuration
		pipName := GetNameFromResourceID(*feProps.PublicIPAddress.ID)

		// Get the Public IP of the PublicIPAddress
		ipValue, err := GetIPOfPublicIPAddressByNameContextE(ctx, pipName, resourceGroupName, subscriptionID)
		if err != nil {
			return "", NoIP, err
		}

		return ipValue, PublicIP, nil
	}

	// Return the Private IP as there are no other option available
	return *feProps.PrivateIPAddress, PrivateIP, nil
}

// GetLoadBalancerFrontendIPConfigContext gets the specified Load Balancer Frontend IP Configuration network resource.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigContext(t testing.TestingT, ctx context.Context, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) *armnetwork.FrontendIPConfiguration {
	t.Helper()

	lbFEConfig, err := GetLoadBalancerFrontendIPConfigContextE(ctx, feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return lbFEConfig
}

// GetLoadBalancerFrontendIPConfigContextE gets the specified Load Balancer Frontend IP Configuration network resource.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigContextE(ctx context.Context, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (*armnetwork.FrontendIPConfiguration, error) {
	// Validate Azure Resource Group Name
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetLoadBalancerFrontendIPConfigClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer Frontend IP Configuration
	resp, err := client.Get(ctx, resourceGroupName, loadBalancerName, feConfigName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.FrontendIPConfiguration, nil
}

// GetLoadBalancerFrontendIPConfigClientContextE gets a new Load Balancer Frontend IP Configuration client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.LoadBalancerFrontendIPConfigurationsClient, error) {
	return CreateLoadBalancerFrontendIPConfigClientContextE(ctx, subscriptionID)
}

// GetLoadBalancerFrontendIPConfigClientE gets a new Load Balancer Frontend IP Configuration client in the specified Azure Subscription.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigClientContextE] instead.
func GetLoadBalancerFrontendIPConfigClientE(subscriptionID string) (*armnetwork.LoadBalancerFrontendIPConfigurationsClient, error) {
	return GetLoadBalancerFrontendIPConfigClientContextE(context.Background(), subscriptionID)
}

// GetLoadBalancerContext gets a Load Balancer network resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerContext(t testing.TestingT, ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) *armnetwork.LoadBalancer {
	t.Helper()

	lb, err := GetLoadBalancerContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return lb
}

// GetLoadBalancerContextE gets a Load Balancer network resource in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerContextE(ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) (*armnetwork.LoadBalancer, error) {
	// Validate Azure Resource Group Name
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetLoadBalancerClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer
	resp, err := client.Get(ctx, resourceGroupName, loadBalancerName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.LoadBalancer, nil
}

// GetLoadBalancerClientContextE gets a new Load Balancer client in the specified Azure Subscription.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerClientContextE(ctx context.Context, subscriptionID string) (*armnetwork.LoadBalancersClient, error) {
	return CreateLoadBalancerClientContextE(ctx, subscriptionID)
}

// GetLoadBalancerClientE gets a new Load Balancer client in the specified Azure Subscription.
//
// Deprecated: Use [GetLoadBalancerClientContextE] instead.
func GetLoadBalancerClientE(subscriptionID string) (*armnetwork.LoadBalancersClient, error) {
	return GetLoadBalancerClientContextE(context.Background(), subscriptionID)
}
