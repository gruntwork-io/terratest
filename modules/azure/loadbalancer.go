package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
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

// LoadBalancerExists indicates whether the specified Load Balancer exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [LoadBalancerExistsContext] instead.
func LoadBalancerExists(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return LoadBalancerExistsContext(t, context.Background(), loadBalancerName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// LoadBalancerExistsE indicates whether the specified Load Balancer exists.
//
// Deprecated: Use [LoadBalancerExistsContextE] instead.
func LoadBalancerExistsE(loadBalancerName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return LoadBalancerExistsContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
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

// GetLoadBalancerFrontendIPConfigNames gets a list of the Frontend IP Configuration Names for the Load Balancer.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigNamesContext] instead.
func GetLoadBalancerFrontendIPConfigNames(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetLoadBalancerFrontendIPConfigNamesContext(t, context.Background(), loadBalancerName, resourceGroupName, subscriptionID) //nolint:staticcheck
}

// GetLoadBalancerFrontendIPConfigNamesContextE gets a list of the Frontend IP Configuration Names for the Load Balancer.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigNamesContextE(ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) ([]string, error) {
	lb, err := GetLoadBalancerContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
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

// GetLoadBalancerFrontendIPConfigNamesE gets a list of the Frontend IP Configuration Names for the Load Balancer.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigNamesContextE] instead.
func GetLoadBalancerFrontendIPConfigNamesE(loadBalancerName string, resourceGroupName string, subscriptionID string) ([]string, error) {
	return GetLoadBalancerFrontendIPConfigNamesContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
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

// GetIPOfLoadBalancerFrontendIPConfig gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetIPOfLoadBalancerFrontendIPConfigContext] instead.
func GetIPOfLoadBalancerFrontendIPConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType) {
	t.Helper()

	return GetIPOfLoadBalancerFrontendIPConfigContext(t, context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID) //nolint:staticcheck
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
	feProps := *feConfig.FrontendIPConfigurationPropertiesFormat

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

// GetIPOfLoadBalancerFrontendIPConfigE gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
//
// Deprecated: Use [GetIPOfLoadBalancerFrontendIPConfigContextE] instead.
func GetIPOfLoadBalancerFrontendIPConfigE(feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType, err1 error) {
	return GetIPOfLoadBalancerFrontendIPConfigContextE(context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
}

// GetLoadBalancerFrontendIPConfigContext gets the specified Load Balancer Frontend IP Configuration network resource.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigContext(t testing.TestingT, ctx context.Context, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.FrontendIPConfiguration {
	t.Helper()

	lbFEConfig, err := GetLoadBalancerFrontendIPConfigContextE(ctx, feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return lbFEConfig
}

// GetLoadBalancerFrontendIPConfig gets the specified Load Balancer Frontend IP Configuration network resource.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigContext] instead.
func GetLoadBalancerFrontendIPConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.FrontendIPConfiguration {
	t.Helper()

	return GetLoadBalancerFrontendIPConfigContext(t, context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID) //nolint:staticcheck
}

// GetLoadBalancerFrontendIPConfigContextE gets the specified Load Balancer Frontend IP Configuration network resource.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerFrontendIPConfigContextE(ctx context.Context, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.FrontendIPConfiguration, error) {
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
	lbc, err := client.Get(ctx, resourceGroupName, loadBalancerName, feConfigName)
	if err != nil {
		return nil, err
	}

	return &lbc, nil
}

// GetLoadBalancerFrontendIPConfigE gets the specified Load Balancer Frontend IP Configuration network resource.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigContextE] instead.
func GetLoadBalancerFrontendIPConfigE(feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.FrontendIPConfiguration, error) {
	return GetLoadBalancerFrontendIPConfigContextE(context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
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

// GetLoadBalancerContext gets a Load Balancer network resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerContext(t testing.TestingT, ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.LoadBalancer {
	t.Helper()

	lb, err := GetLoadBalancerContextE(ctx, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return lb
}

// GetLoadBalancer gets a Load Balancer network resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLoadBalancerContext] instead.
func GetLoadBalancer(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.LoadBalancer {
	t.Helper()

	return GetLoadBalancerContext(t, context.Background(), loadBalancerName, resourceGroupName, subscriptionID) //nolint:staticcheck
}

// GetLoadBalancerContextE gets a Load Balancer network resource in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetLoadBalancerContextE(ctx context.Context, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.LoadBalancer, error) {
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
	lb, err := client.Get(ctx, resourceGroupName, loadBalancerName, "")
	if err != nil {
		return nil, err
	}

	return &lb, nil
}

// GetLoadBalancerE gets a Load Balancer network resource in the specified Azure Resource Group.
//
// Deprecated: Use [GetLoadBalancerContextE] instead.
func GetLoadBalancerE(loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.LoadBalancer, error) {
	return GetLoadBalancerContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
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
