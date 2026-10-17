package gcp

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/servicenetworking/v1"
)

// serviceNetworkingParent is the service a private connection is made to. Every connection a
// Terraform module creates is to this one, which is what serves Cloud SQL, Memorystore and the
// other managed products that live in Google's own network.
const serviceNetworkingParent = "services/servicenetworking.googleapis.com"

// GetServiceNetworkingConnectionAttrs returns the settings Google Cloud holds for the private
// connection between the given network and Google's services, so a test can assert on what was
// actually created rather than only that it exists. A connection has no name of its own and cannot
// be fetched one at a time, so this lists the network's connections and returns the only one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingConnectionAttrs(t testing.TestingT, ctx context.Context, projectID string, networkName string) *servicenetworking.Connection {
	connection, err := GetServiceNetworkingConnectionAttrsE(t, ctx, projectID, networkName)
	require.NoError(t, err)

	return connection
}

// GetServiceNetworkingConnectionAttrsE returns the settings Google Cloud holds for the private
// connection between the given network and Google's services.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingConnectionAttrsE(t testing.TestingT, ctx context.Context, projectID string, networkName string) (*servicenetworking.Connection, error) {
	logger.Default.Logf(t, "Getting settings for the private connection on network %s in project %s", networkName, projectID)

	service, err := NewServiceNetworkingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceNetworkingConnectionAttrsWithClient(ctx, service, projectID, networkName)
}

// GetServiceNetworkingConnectionAttrsWithClient returns the settings Google Cloud holds for the
// private connection between the given network and Google's services using the supplied
// *servicenetworking.APIService. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see servicenetworking_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingConnectionAttrsWithClient(ctx context.Context, service *servicenetworking.APIService, projectID string, networkName string) (*servicenetworking.Connection, error) {
	network := fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)

	response, err := service.Services.Connections.List(serviceNetworkingParent).Network(network).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list the private connections on network %s in project %s: %w", networkName, projectID, err)
	}

	// The list is filtered to one network, so anything but a single connection is the test's answer
	// rather than a choice this read should make.
	switch len(response.Connections) {
	case 0:
		return nil, fmt.Errorf("no private connection exists on network %s in project %s", networkName, projectID)
	case 1:
		return response.Connections[0], nil
	default:
		names := make([]string, 0, len(response.Connections))
		for _, connection := range response.Connections {
			names = append(names, connection.Peering)
		}

		return nil, fmt.Errorf("network %s in project %s has %d private connections, %s, and not the one this read expects", networkName, projectID, len(response.Connections), strings.Join(names, ", "))
	}
}

// GetServiceNetworkingPeeredDNSDomainAttrs returns the settings Google Cloud holds for the given
// peered DNS domain, so a test can assert on what was actually created rather than only that it
// exists.
// The API offers no get for one domain, only a list over the network, so this reads the list and
// picks the domain out of it.
// Service Networking names the network under the service rather than under the project, and it
// refuses a parent that carries the project id: pass the project number.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingPeeredDNSDomainAttrs(t testing.TestingT, ctx context.Context, projectNumber string, network string, domainID string) *servicenetworking.PeeredDnsDomain {
	domain, err := GetServiceNetworkingPeeredDNSDomainAttrsE(t, ctx, projectNumber, network, domainID)
	require.NoError(t, err)

	return domain
}

// GetServiceNetworkingPeeredDNSDomainAttrsE returns the settings Google Cloud holds for the given
// peered DNS domain.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingPeeredDNSDomainAttrsE(t testing.TestingT, ctx context.Context, projectNumber string, network string, domainID string) (*servicenetworking.PeeredDnsDomain, error) {
	logger.Default.Logf(t, "Getting settings for peered DNS domain %s on network %s in project %s", domainID, network, projectNumber)

	service, err := NewServiceNetworkingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceNetworkingPeeredDNSDomainAttrsWithClient(ctx, service, projectNumber, network, domainID)
}

// GetServiceNetworkingPeeredDNSDomainAttrsWithClient returns the settings Google Cloud holds for the
// given peered DNS domain using the supplied *servicenetworking.APIService. Prefer this variant in
// unit tests where the service is backed by an httptest fake server (see servicenetworking_test.go
// for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingPeeredDNSDomainAttrsWithClient(ctx context.Context, service *servicenetworking.APIService, projectNumber string, network string, domainID string) (*servicenetworking.PeeredDnsDomain, error) {
	parent := fmt.Sprintf("%s/projects/%s/global/networks/%s", serviceNetworkingParent, projectNumber, network)

	response, err := service.Services.Projects.Global.Networks.PeeredDnsDomains.List(parent).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the network %s in project %s has no service networking connection", network, projectNumber)
		}

		return nil, fmt.Errorf("failed to list the peered DNS domains on network %s in project %s: %w", network, projectNumber, err)
	}

	for _, domain := range response.PeeredDnsDomains {
		// A listed domain's name comes back in full, so the id the caller passed is its last segment.
		if domain.Name == domainID || path.Base(domain.Name) == domainID {
			return domain, nil
		}
	}

	return nil, fmt.Errorf("the peered DNS domain %s does not exist on network %s in project %s", domainID, network, projectNumber)
}

// GetServiceNetworkingVPCServiceControlsAttrs returns the consumer configuration Google Cloud holds
// for the given producer network peering, so a test can assert whether VPC Service Controls are
// enforced on it rather than only that the apply finished.
// The setting is part of the connection's consumer config rather than a resource of its own, so the
// whole config is returned and the caller asserts on the field it cares about.
// Service Networking names the network under the service rather than under the project, and it
// refuses a parent that carries the project id: pass the project number.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingVPCServiceControlsAttrs(t testing.TestingT, ctx context.Context, projectNumber string, network string) *servicenetworking.ConsumerConfig {
	config, err := GetServiceNetworkingVPCServiceControlsAttrsE(t, ctx, projectNumber, network)
	require.NoError(t, err)

	return config
}

// GetServiceNetworkingVPCServiceControlsAttrsE returns the consumer configuration Google Cloud holds
// for the given producer network peering.
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingVPCServiceControlsAttrsE(t testing.TestingT, ctx context.Context, projectNumber string, network string) (*servicenetworking.ConsumerConfig, error) {
	logger.Default.Logf(t, "Getting the service networking consumer config for network %s in project %s", network, projectNumber)

	service, err := NewServiceNetworkingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceNetworkingVPCServiceControlsAttrsWithClient(ctx, service, projectNumber, network)
}

// GetServiceNetworkingVPCServiceControlsAttrsWithClient returns the consumer configuration Google
// Cloud holds for the given producer network peering using the supplied
// *servicenetworking.APIService. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see servicenetworking_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceNetworkingVPCServiceControlsAttrsWithClient(ctx context.Context, service *servicenetworking.APIService, projectNumber string, network string) (*servicenetworking.ConsumerConfig, error) {
	name := fmt.Sprintf("%s/projects/%s/global/networks/%s", serviceNetworkingParent, projectNumber, network)

	config, err := service.Services.Projects.Global.Networks.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the network %s in project %s has no service networking connection", network, projectNumber)
		}

		return nil, fmt.Errorf("failed to get the service networking consumer config for network %s in project %s: %w", network, projectNumber, err)
	}

	return config, nil
}

// NewServiceNetworkingServiceE creates a Service Networking service authenticated the same way
// every other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewServiceNetworkingServiceE(t testing.TestingT, ctx context.Context) (*servicenetworking.APIService, error) {
	return servicenetworking.NewService(ctx, append(withOptions(), option.WithScopes(servicenetworking.CloudPlatformScope))...)
}
