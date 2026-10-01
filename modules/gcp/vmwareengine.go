package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/vmwareengine/v1"
)

// GetVMwareEnginePrivateCloudAttrs returns the settings Google Cloud holds for the given VMware Engine private cloud, so a test can assert on what was
// actually created rather than only that it exists.
// A private cloud is the vSphere deployment itself, so its type and its management cluster are what it is made of.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEnginePrivateCloudAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string) *vmwareengine.PrivateCloud {
	attrs, err := GetVMwareEnginePrivateCloudAttrsE(t, ctx, projectID, location, cloudID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEnginePrivateCloudAttrsE returns the settings Google Cloud holds for the given VMware Engine private cloud.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEnginePrivateCloudAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string) (*vmwareengine.PrivateCloud, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine private cloud %s in %s in project %s", cloudID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEnginePrivateCloudAttrsWithClient(ctx, service, projectID, location, cloudID)
}

// GetVMwareEnginePrivateCloudAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine private cloud using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEnginePrivateCloudAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, cloudID string) (*vmwareengine.PrivateCloud, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/privateClouds/%s", projectID, location, cloudID)

	attrs, err := service.Projects.Locations.PrivateClouds.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine private cloud %s in %s in project %s does not exist", cloudID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine private cloud %s in %s in project %s: %w", cloudID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineClusterAttrs returns the settings Google Cloud holds for the given VMware Engine cluster, so a test can assert on what was
// actually created rather than only that it exists.
// A cluster is a set of nodes inside a private cloud, so its node configuration is what is billed and what runs VMs.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string, clusterID string) *vmwareengine.Cluster {
	attrs, err := GetVMwareEngineClusterAttrsE(t, ctx, projectID, location, cloudID, clusterID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineClusterAttrsE returns the settings Google Cloud holds for the given VMware Engine cluster.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string, clusterID string) (*vmwareengine.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine cluster %s %s in %s in project %s", clusterID, cloudID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineClusterAttrsWithClient(ctx, service, projectID, location, cloudID, clusterID)
}

// GetVMwareEngineClusterAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine cluster using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineClusterAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, cloudID string, clusterID string) (*vmwareengine.Cluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/privateClouds/%s/clusters/%s", projectID, location, cloudID, clusterID)

	attrs, err := service.Projects.Locations.PrivateClouds.Clusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine cluster %s %s in %s in project %s does not exist", clusterID, cloudID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine cluster %s %s in %s in project %s: %w", clusterID, cloudID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineSubnetAttrs returns the settings Google Cloud holds for the given VMware Engine subnet, so a test can assert on what was
// actually created rather than only that it exists.
// A subnet is a range inside the private cloud's management network, so its range and type decide what can be addressed.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineSubnetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string, subnetID string) *vmwareengine.Subnet {
	attrs, err := GetVMwareEngineSubnetAttrsE(t, ctx, projectID, location, cloudID, subnetID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineSubnetAttrsE returns the settings Google Cloud holds for the given VMware Engine subnet.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineSubnetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string, subnetID string) (*vmwareengine.Subnet, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine subnet %s %s in %s in project %s", subnetID, cloudID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineSubnetAttrsWithClient(ctx, service, projectID, location, cloudID, subnetID)
}

// GetVMwareEngineSubnetAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine subnet using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineSubnetAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, cloudID string, subnetID string) (*vmwareengine.Subnet, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/privateClouds/%s/subnets/%s", projectID, location, cloudID, subnetID)

	attrs, err := service.Projects.Locations.PrivateClouds.Subnets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine subnet %s %s in %s in project %s does not exist", subnetID, cloudID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine subnet %s %s in %s in project %s: %w", subnetID, cloudID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineExternalAddressAttrs returns the settings Google Cloud holds for the given VMware Engine external address, so a test can assert on what was
// actually created rather than only that it exists.
// An external address is how one internal VM is reachable from outside, so the internal IP it maps is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineExternalAddressAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string, addressID string) *vmwareengine.ExternalAddress {
	attrs, err := GetVMwareEngineExternalAddressAttrsE(t, ctx, projectID, location, cloudID, addressID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineExternalAddressAttrsE returns the settings Google Cloud holds for the given VMware Engine external address.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineExternalAddressAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, cloudID string, addressID string) (*vmwareengine.ExternalAddress, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine external address %s %s in %s in project %s", addressID, cloudID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineExternalAddressAttrsWithClient(ctx, service, projectID, location, cloudID, addressID)
}

// GetVMwareEngineExternalAddressAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine external address using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineExternalAddressAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, cloudID string, addressID string) (*vmwareengine.ExternalAddress, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/privateClouds/%s/externalAddresses/%s", projectID, location, cloudID, addressID)

	attrs, err := service.Projects.Locations.PrivateClouds.ExternalAddresses.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine external address %s %s in %s in project %s does not exist", addressID, cloudID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine external address %s %s in %s in project %s: %w", addressID, cloudID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineNetworkAttrs returns the settings Google Cloud holds for the given VMware Engine network, so a test can assert on what was
// actually created rather than only that it exists.
// The network is what private clouds and peerings attach to, so its type and state decide what can join it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, networkID string) *vmwareengine.VmwareEngineNetwork {
	attrs, err := GetVMwareEngineNetworkAttrsE(t, ctx, projectID, location, networkID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineNetworkAttrsE returns the settings Google Cloud holds for the given VMware Engine network.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, networkID string) (*vmwareengine.VmwareEngineNetwork, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine network %s in %s in project %s", networkID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineNetworkAttrsWithClient(ctx, service, projectID, location, networkID)
}

// GetVMwareEngineNetworkAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine network using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, networkID string) (*vmwareengine.VmwareEngineNetwork, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/vmwareEngineNetworks/%s", projectID, location, networkID)

	attrs, err := service.Projects.Locations.VmwareEngineNetworks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine network %s in %s in project %s does not exist", networkID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine network %s in %s in project %s: %w", networkID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineNetworkPeeringAttrs returns the settings Google Cloud holds for the given VMware Engine network peering, so a test can assert on what was
// actually created rather than only that it exists.
// A peering joins the VMware network to another network, so which networks it names and whether routes are exported are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkPeeringAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, peeringID string) *vmwareengine.NetworkPeering {
	attrs, err := GetVMwareEngineNetworkPeeringAttrsE(t, ctx, projectID, location, peeringID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineNetworkPeeringAttrsE returns the settings Google Cloud holds for the given VMware Engine network peering.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkPeeringAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, peeringID string) (*vmwareengine.NetworkPeering, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine network peering %s in %s in project %s", peeringID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineNetworkPeeringAttrsWithClient(ctx, service, projectID, location, peeringID)
}

// GetVMwareEngineNetworkPeeringAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine network peering using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkPeeringAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, peeringID string) (*vmwareengine.NetworkPeering, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/networkPeerings/%s", projectID, location, peeringID)

	attrs, err := service.Projects.Locations.NetworkPeerings.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine network peering %s in %s in project %s does not exist", peeringID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine network peering %s in %s in project %s: %w", peeringID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineNetworkPolicyAttrs returns the settings Google Cloud holds for the given VMware Engine network policy, so a test can assert on what was
// actually created rather than only that it exists.
// A policy is what turns internet access and external IPs on for a region, so what it enables is exactly what is reachable.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *vmwareengine.NetworkPolicy {
	attrs, err := GetVMwareEngineNetworkPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineNetworkPolicyAttrsE returns the settings Google Cloud holds for the given VMware Engine network policy.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*vmwareengine.NetworkPolicy, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine network policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineNetworkPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetVMwareEngineNetworkPolicyAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine network policy using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineNetworkPolicyAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, policyID string) (*vmwareengine.NetworkPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/networkPolicies/%s", projectID, location, policyID)

	attrs, err := service.Projects.Locations.NetworkPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine network policy %s in %s in project %s does not exist", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine network policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineExternalAccessRuleAttrs returns the settings Google Cloud holds for the given VMware Engine external access rule, so a test can assert on what was
// actually created rather than only that it exists.
// A rule is the firewall entry that lets outside traffic in, so its action, priority and ports are exactly what it permits.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineExternalAccessRuleAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string, ruleID string) *vmwareengine.ExternalAccessRule {
	attrs, err := GetVMwareEngineExternalAccessRuleAttrsE(t, ctx, projectID, location, policyID, ruleID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineExternalAccessRuleAttrsE returns the settings Google Cloud holds for the given VMware Engine external access rule.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineExternalAccessRuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string, ruleID string) (*vmwareengine.ExternalAccessRule, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine external access rule %s %s in %s in project %s", ruleID, policyID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineExternalAccessRuleAttrsWithClient(ctx, service, projectID, location, policyID, ruleID)
}

// GetVMwareEngineExternalAccessRuleAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine external access rule using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineExternalAccessRuleAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, policyID string, ruleID string) (*vmwareengine.ExternalAccessRule, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/networkPolicies/%s/externalAccessRules/%s", projectID, location, policyID, ruleID)

	attrs, err := service.Projects.Locations.NetworkPolicies.ExternalAccessRules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine external access rule %s %s in %s in project %s does not exist", ruleID, policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine external access rule %s %s in %s in project %s: %w", ruleID, policyID, location, projectID, err)
	}

	return attrs, nil
}

// GetVMwareEngineDatastoreAttrs returns the settings Google Cloud holds for the given VMware Engine datastore, so a test can assert on what was
// actually created rather than only that it exists.
// A datastore is the storage a cluster's VMs sit on, so its type and capacity are what they draw from.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineDatastoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datastoreID string) *vmwareengine.Datastore {
	attrs, err := GetVMwareEngineDatastoreAttrsE(t, ctx, projectID, location, datastoreID)
	require.NoError(t, err)

	return attrs
}

// GetVMwareEngineDatastoreAttrsE returns the settings Google Cloud holds for the given VMware Engine datastore.
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineDatastoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datastoreID string) (*vmwareengine.Datastore, error) {
	logger.Default.Logf(t, "Getting settings for VMware Engine datastore %s in %s in project %s", datastoreID, location, projectID)

	service, err := NewVMwareEngineServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetVMwareEngineDatastoreAttrsWithClient(ctx, service, projectID, location, datastoreID)
}

// GetVMwareEngineDatastoreAttrsWithClient returns the settings Google Cloud holds for the given VMware Engine datastore using the supplied
// *vmwareengine.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see vmwareengine_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetVMwareEngineDatastoreAttrsWithClient(ctx context.Context, service *vmwareengine.Service, projectID string, location string, datastoreID string) (*vmwareengine.Datastore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datastores/%s", projectID, location, datastoreID)

	attrs, err := service.Projects.Locations.Datastores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware Engine datastore %s in %s in project %s does not exist", datastoreID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware Engine datastore %s in %s in project %s: %w", datastoreID, location, projectID, err)
	}

	return attrs, nil
}

// NewVMwareEngineServiceE creates a VMware Engine service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewVMwareEngineServiceE(t testing.TestingT, ctx context.Context) (*vmwareengine.Service, error) {
	return vmwareengine.NewService(ctx, append(withOptions(), option.WithScopes(vmwareengine.CloudPlatformScope))...)
}
