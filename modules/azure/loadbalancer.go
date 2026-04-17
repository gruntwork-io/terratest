package azure

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// LoadBalancerExists indicates whether the specified Load Balancer exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [LoadBalancerExistsContext] instead.
func LoadBalancerExists(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return LoadBalancerExistsContext(t, context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
}

// LoadBalancerExistsE indicates whether the specified Load Balancer exists.
//
// Deprecated: Use [LoadBalancerExistsContextE] instead.
func LoadBalancerExistsE(loadBalancerName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return LoadBalancerExistsContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
}

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

// GetLoadBalancerFrontendIPConfigNames gets a list of the Frontend IP Configuration Names for the Load Balancer.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigNamesContext] instead.
func GetLoadBalancerFrontendIPConfigNames(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetLoadBalancerFrontendIPConfigNamesContext(t, context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
}

// GetLoadBalancerFrontendIPConfigNamesE gets a list of the Frontend IP Configuration Names for the Load Balancer.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigNamesContextE] instead.
func GetLoadBalancerFrontendIPConfigNamesE(loadBalancerName string, resourceGroupName string, subscriptionID string) ([]string, error) {
	return GetLoadBalancerFrontendIPConfigNamesContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
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

	return ExtractLoadBalancerFrontendIPConfigNames(lb), nil
}

// GetIPOfLoadBalancerFrontendIPConfig gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetIPOfLoadBalancerFrontendIPConfigContext] instead.
func GetIPOfLoadBalancerFrontendIPConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType) {
	t.Helper()

	return GetIPOfLoadBalancerFrontendIPConfigContext(t, context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
}

// GetIPOfLoadBalancerFrontendIPConfigE gets the IP and LoadBalancerIPType for the specified Load Balancer Frontend IP Configuration.
//
// Deprecated: Use [GetIPOfLoadBalancerFrontendIPConfigContextE] instead.
func GetIPOfLoadBalancerFrontendIPConfigE(feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate LoadBalancerIPType, err1 error) {
	return GetIPOfLoadBalancerFrontendIPConfigContextE(context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
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

	// Resolve the IP using a PIP client for public address lookups
	pipClient, err := GetPublicIPAddressClientContextE(ctx, subscriptionID)
	if err != nil {
		return "", NoIP, err
	}

	return GetIPOfLoadBalancerFrontendIPConfigWithClient(ctx, feConfig, pipClient, resourceGroupName)
}

// GetLoadBalancerFrontendIPConfig gets the specified Load Balancer Frontend IP Configuration network resource.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigContext] instead.
func GetLoadBalancerFrontendIPConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) *armnetwork.FrontendIPConfiguration {
	t.Helper()

	return GetLoadBalancerFrontendIPConfigContext(t, context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
}

// GetLoadBalancerFrontendIPConfigE gets the specified Load Balancer Frontend IP Configuration network resource.
//
// Deprecated: Use [GetLoadBalancerFrontendIPConfigContextE] instead.
func GetLoadBalancerFrontendIPConfigE(feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (*armnetwork.FrontendIPConfiguration, error) {
	return GetLoadBalancerFrontendIPConfigContextE(context.Background(), feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
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

	return GetLoadBalancerFrontendIPConfigWithClient(ctx, client, resourceGroupName, loadBalancerName, feConfigName)
}

// GetLoadBalancerFrontendIPConfigWithClient gets the specified Load Balancer Frontend IP Configuration
// using the provided LoadBalancerFrontendIPConfigurationsClient.
func GetLoadBalancerFrontendIPConfigWithClient(ctx context.Context, client *armnetwork.LoadBalancerFrontendIPConfigurationsClient, resourceGroupName string, loadBalancerName string, feConfigName string) (*armnetwork.FrontendIPConfiguration, error) {
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

// GetLoadBalancer gets a Load Balancer network resource in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLoadBalancerContext] instead.
func GetLoadBalancer(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) *armnetwork.LoadBalancer {
	t.Helper()

	return GetLoadBalancerContext(t, context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
}

// GetLoadBalancerE gets a Load Balancer network resource in the specified Azure Resource Group.
//
// Deprecated: Use [GetLoadBalancerContextE] instead.
func GetLoadBalancerE(loadBalancerName string, resourceGroupName string, subscriptionID string) (*armnetwork.LoadBalancer, error) {
	return GetLoadBalancerContextE(context.Background(), loadBalancerName, resourceGroupName, subscriptionID)
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

	return GetLoadBalancerWithClient(ctx, client, resourceGroupName, loadBalancerName)
}

// GetLoadBalancerWithClient gets a Load Balancer using the provided LoadBalancersClient.
func GetLoadBalancerWithClient(ctx context.Context, client *armnetwork.LoadBalancersClient, resourceGroupName string, loadBalancerName string) (*armnetwork.LoadBalancer, error) {
	resp, err := client.Get(ctx, resourceGroupName, loadBalancerName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.LoadBalancer, nil
}

// ExtractLoadBalancerFrontendIPConfigNames gets a list of the Frontend IP Configuration Names
// from a Load Balancer.
func ExtractLoadBalancerFrontendIPConfigNames(lb *armnetwork.LoadBalancer) []string {
	if lb == nil || lb.Properties == nil {
		return nil
	}

	feConfigs := lb.Properties.FrontendIPConfigurations

	if len(feConfigs) == 0 {
		return nil
	}

	configNames := make([]string, 0, len(feConfigs))

	for _, config := range feConfigs {
		if config == nil || config.Name == nil {
			continue
		}

		configNames = append(configNames, *config.Name)
	}

	return configNames
}

// GetIPOfLoadBalancerFrontendIPConfigWithClient gets the IP and LoadBalancerIPType for the
// specified Frontend IP Configuration. For public IPs it requires a PublicIPAddressesClient
// to resolve the public IP address.
func GetIPOfLoadBalancerFrontendIPConfigWithClient(ctx context.Context, feConfig *armnetwork.FrontendIPConfiguration, pipClient *armnetwork.PublicIPAddressesClient, resourceGroupName string) (string, LoadBalancerIPType, error) {
	if feConfig == nil || feConfig.Properties == nil {
		return "", NoIP, errors.New("frontend IP configuration has nil properties")
	}

	feProps := feConfig.Properties

	pip := feProps.PublicIPAddress
	if pip == nil || pip.ID == nil {
		if feProps.PrivateIPAddress == nil {
			return "", NoIP, errors.New("frontend IP configuration has no private or public IP address assigned")
		}

		return *feProps.PrivateIPAddress, PrivateIP, nil
	}

	pipName := GetNameFromResourceID(*pip.ID)

	ipValue, err := GetPublicIPAddressWithClient(ctx, pipClient, resourceGroupName, pipName)
	if err != nil {
		return "", NoIP, err
	}

	ip, err := ExtractIPOfPublicIPAddress(ipValue)
	if err != nil {
		return "", NoIP, err
	}

	return ip, PublicIP, nil
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
