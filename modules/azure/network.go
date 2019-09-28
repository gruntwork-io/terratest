package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-10-01/network"
	"github.com/stretchr/testify/require"
)

// GetVirtualNetworkClient is a helper function to setup an Azure Virtual Network client
func GetVirtualNetworkClient(subscriptionID string) (*network.VirtualNetworksClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a VNet client
	vnetClient := network.NewVirtualNetworksClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	vnetClient.Authorizer = *authorizer

	return &vnetClient, nil
}

// GetVirtualNetwork gets the details of an Azure Virtual Network
func GetVirtualNetwork(t *testing.T, vnetName string, resGroupName string, subscriptionID string) network.VirtualNetwork {
	vnet, err := GetVirtualNetworkE(t, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vnet
}

// GetVirtualNetworkE gets the details of an Azure Virtual Network
func GetVirtualNetworkE(t *testing.T, vnetName string, resGroupName string, subscriptionID string) (network.VirtualNetwork, error) {
	vnet := network.VirtualNetwork{}

	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return vnet, err
	}

	// Create a VNet client
	vnetClient, err := GetVirtualNetworkClient(subscriptionID)
	if err != nil {
		return vnet, err
	}

	// Get the details of the Virtual Network
	vnet, err = vnetClient.Get(context.Background(), resGroupName, vnetName, "")
	if err != nil {
		return vnet, err
	}
	return vnet, nil
}

// AssertVirtualNetworkExists checks if a given Virtual Network exists in the specified Resource Group and returns an error if it does not
func AssertVirtualNetworkExists(t *testing.T, vnetName string, resGroupName string, subscriptionID string) {
	_, err := GetVirtualNetworkE(t, vnetName, resGroupName, subscriptionID)
	// If the Virtual Network does not exist, the API returns a "ResourceNotFound" error
	require.NoError(t, err)
}

// Subnet is a representation of an Azure Virtual Network subnet
type Subnet struct {
	ID            string
	Name          string
	AddressPrefix string
}

// GetSubnetsForVirtualNetwork gets all subnets in a given Azure Virtual Network
func GetSubnetsForVirtualNetwork(t *testing.T, vnetName string, resGroupName string, subscriptionID string) []Subnet {
	subnets, err := GetSubnetsForVirtualNetworkE(t, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return subnets
}

// GetSubnetsForVirtualNetworkE gets all subnets in a given Azure Virtual Network
func GetSubnetsForVirtualNetworkE(t *testing.T, vnetName string, resGroupName string, subscriptionID string) ([]Subnet, error) {
	subnets := []Subnet{}

	vnet, err := GetVirtualNetworkE(t, vnetName, resGroupName, subscriptionID)
	if err != nil {
		return subnets, err
	}

	for _, subnet := range *(vnet.Subnets) {
		subnets = append(subnets, Subnet{ID: *(subnet.ID), Name: *(subnet.Name), AddressPrefix: *(subnet.AddressPrefix)})
	}

	return subnets, nil
}

// GetTagsForVirtualNetwork gets the tags of the given Virtual Network as a map
func GetTagsForVirtualNetwork(t *testing.T, vnetName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetTagsForVirtualNetworkE(t, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetTagsForVirtualNetworkE gets the tags of the given Virtual Network as a map
func GetTagsForVirtualNetworkE(t *testing.T, vnetName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	vnet, err := GetVirtualNetworkE(t, vnetName, resGroupName, subscriptionID)
	if err != nil {
		return tags, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vnet.Tags {
		tags[k] = *v
	}

	return tags, nil
}

// GetSubnetClient is a helper function to setup an Azure subnet client
func GetSubnetClient(subscriptionID string) (*network.SubnetsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a subnet client
	subnetClient := network.NewSubnetsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	subnetClient.Authorizer = *authorizer

	return &subnetClient, nil
}

// NetworkSecurityGroup is a representation of an Azure Network Security Group
type NetworkSecurityGroup struct {
	ID                 string
	Name               string
	SecurityRulesNames []string
}

// GetNetworkSecurityGroupForSubnet gets the Network Security Group associated with a given Virtual Network subnet
func GetNetworkSecurityGroupForSubnet(t *testing.T, subnetName string, vnetName string, resGroupName string, subscriptionID string) NetworkSecurityGroup {
	nsg, err := GetNetworkSecurityGroupForSubnetE(t, subnetName, vnetName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nsg
}

// GetNetworkSecurityGroupForSubnetE gets the Network Security Group associated with a given Virtual Network subnet
func GetNetworkSecurityGroupForSubnetE(t *testing.T, subnetName string, vnetName string, resGroupName string, subscriptionID string) (NetworkSecurityGroup, error) {
	nsg := NetworkSecurityGroup{}

	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nsg, err
	}

	// Create a subnet client
	subnetClient, err := GetSubnetClient(subscriptionID)
	if err != nil {
		return nsg, err
	}

	// Get the details of the subnet and its associated Network Security Group
	subnet, err := subnetClient.Get(context.Background(), resGroupName, vnetName, subnetName, "networkSecurityGroup")
	if err != nil {
		return nsg, err
	}
	if subnet.NetworkSecurityGroup != nil {
		rulesNames := []string{}
		for _, r := range *(subnet.NetworkSecurityGroup.SecurityRules) {
			rulesNames = append(rulesNames, *(r.Name))
		}

		nsg = NetworkSecurityGroup{
			ID:                 *(subnet.NetworkSecurityGroup.ID),
			Name:               *(subnet.NetworkSecurityGroup.Name),
			SecurityRulesNames: rulesNames,
		}
	}

	return nsg, nil
}

// GetPublicIPClient is a helper function to setup an Azure Public IP client
func GetPublicIPClient(subscriptionID string) (*network.PublicIPAddressesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a subnet client
	ipClient := network.NewPublicIPAddressesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	ipClient.Authorizer = *authorizer

	return &ipClient, nil
}

// PublicIP is a representation of an Azure Public IP Address resource
type PublicIP struct {
	ID               string
	Name             string
	IPAddress        string
	IPAddressVersion string
	AllocationMethod string
	FullDNSName      string
}

// GetPublicIP gets an Azure Public IP Address resource
func GetPublicIP(t *testing.T, resGroupName string, publicIPName string, subscriptionID string) PublicIP {
	ip, err := GetPublicIPE(t, resGroupName, publicIPName, subscriptionID)
	require.NoError(t, err)

	return ip
}

// GetPublicIPE gets an Azure Public IP Address resource
func GetPublicIPE(t *testing.T, resGroupName string, publicIPName string, subscriptionID string) (PublicIP, error) {
	ip := PublicIP{}

	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return ip, err
	}

	// Create a public IP client
	ipClient, err := GetPublicIPClient(subscriptionID)
	if err != nil {
		return ip, err
	}

	azurePublicIP, err := ipClient.Get(context.Background(), resGroupName, publicIPName, "ipConfiguration")
	if err != nil {
		return ip, err
	}

	ip = PublicIP{
		ID:               *(azurePublicIP.ID),
		Name:             *(azurePublicIP.Name),
		IPAddress:        *(azurePublicIP.IPAddress),
		IPAddressVersion: string(azurePublicIP.PublicIPAddressVersion),
		AllocationMethod: string(azurePublicIP.PublicIPAllocationMethod),
		FullDNSName:      *(azurePublicIP.DNSSettings.Fqdn),
	}

	return ip, nil
}

// CheckPublicDNSNameAvailability checks whether a domain name in the cloudapp.azure.com zone is available for use
func CheckPublicDNSNameAvailability(t *testing.T, location string, domainNameLabel string, subscriptionID string) (available bool) {
	available, err := CheckPublicDNSNameAvailabilityE(t, location, domainNameLabel, subscriptionID)
	require.NoError(t, err)

	return
}

// CheckPublicDNSNameAvailabilityE checks whether a domain name in the cloudapp.azure.com zone is available for use
func CheckPublicDNSNameAvailabilityE(t *testing.T, location string, domainNameLabel string, subscriptionID string) (available bool, err error) {
	ipClient, err := GetPublicIPClient(subscriptionID)
	if err != nil {
		return
	}

	res, err := (*ipClient).CheckDNSNameAvailability(context.Background(), location, domainNameLabel)
	if err != nil {
		return
	}
	available = *(res.Available)
	return
}
