package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
)

func vnetClient(subscriptionID string) (*network.VirtualNetworksClient, error) {
	authorizer, err := DeploymentServicePrincipalAuthorizer()
	if err != nil {
		return nil, err
	}

	client := network.NewVirtualNetworksClient(subscriptionID)
	client.Authorizer = authorizer
	return &client, err
}

// VnetSubnetsListE - Return the subnets that exist wihin a given VNET
func VnetSubnetsListE(subscriptionID string, resourceGroupName string, vnetName string) ([]string, error) {

	client, err := vnetClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	vnet, err := client.Get(context.Background(), resourceGroupName, vnetName, "")
	if err != nil {
		return nil, err
	}

	subnets := make([]string, len(*vnet.VirtualNetworkPropertiesFormat.Subnets))
	for index, subnet := range *vnet.VirtualNetworkPropertiesFormat.Subnets {
		subnets[index] = *subnet.ID
	}

	return subnets, nil
}

// VnetSubnetsList - Like VnetSubnetsListE but fails in the case an error is returned
func VnetSubnetsList(t *testing.T, subscriptionID string, resourceGroupName string, vnetName string) []string {
	subnets, err := VnetSubnetsListE(subscriptionID, resourceGroupName, vnetName)
	if err != nil {
		t.Fatal(err)
	}
	return subnets
}
