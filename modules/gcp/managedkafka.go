package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/managedkafka/v1"
	"google.golang.org/api/option"
)

// GetManagedKafkaClusterAttrs returns the settings Google Cloud holds for the given Managed Kafka cluster, so a test can assert on what was
// actually created rather than only that it exists.
// A cluster is the broker set itself, so its capacity and the subnet it is reachable on decide what can produce to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *managedkafka.Cluster {
	attrs, err := GetManagedKafkaClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return attrs
}

// GetManagedKafkaClusterAttrsE returns the settings Google Cloud holds for the given Managed Kafka cluster.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*managedkafka.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for Managed Kafka cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewManagedKafkaServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetManagedKafkaClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetManagedKafkaClusterAttrsWithClient returns the settings Google Cloud holds for the given Managed Kafka cluster using the supplied
// *managedkafka.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see managedkafka_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaClusterAttrsWithClient(ctx context.Context, service *managedkafka.Service, projectID string, location string, clusterID string) (*managedkafka.Cluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, location, clusterID)

	attrs, err := service.Projects.Locations.Clusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Managed Kafka cluster %s in %s in project %s does not exist", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Managed Kafka cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetManagedKafkaTopicAttrs returns the settings Google Cloud holds for the given Managed Kafka topic, so a test can assert on what was
// actually created rather than only that it exists.
// A topic's partition count and replication factor cannot be lowered later, so what they were created as matters.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaTopicAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, topicID string) *managedkafka.Topic {
	attrs, err := GetManagedKafkaTopicAttrsE(t, ctx, projectID, location, clusterID, topicID)
	require.NoError(t, err)

	return attrs
}

// GetManagedKafkaTopicAttrsE returns the settings Google Cloud holds for the given Managed Kafka topic.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaTopicAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, topicID string) (*managedkafka.Topic, error) {
	logger.Default.Logf(t, "Getting settings for Managed Kafka topic %s %s in %s in project %s", topicID, clusterID, location, projectID)

	service, err := NewManagedKafkaServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetManagedKafkaTopicAttrsWithClient(ctx, service, projectID, location, clusterID, topicID)
}

// GetManagedKafkaTopicAttrsWithClient returns the settings Google Cloud holds for the given Managed Kafka topic using the supplied
// *managedkafka.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see managedkafka_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaTopicAttrsWithClient(ctx context.Context, service *managedkafka.Service, projectID string, location string, clusterID string, topicID string) (*managedkafka.Topic, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s/topics/%s", projectID, location, clusterID, topicID)

	attrs, err := service.Projects.Locations.Clusters.Topics.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Managed Kafka topic %s %s in %s in project %s does not exist", topicID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Managed Kafka topic %s %s in %s in project %s: %w", topicID, clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetManagedKafkaACLAttrs returns the settings Google Cloud holds for the given Managed Kafka ACL, so a test can assert on what was
// actually created rather than only that it exists.
// An ACL is what decides who may read or write a topic, so the entries it carries are exactly the access it grants.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaACLAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, aclID string) *managedkafka.Acl {
	attrs, err := GetManagedKafkaACLAttrsE(t, ctx, projectID, location, clusterID, aclID)
	require.NoError(t, err)

	return attrs
}

// GetManagedKafkaACLAttrsE returns the settings Google Cloud holds for the given Managed Kafka ACL.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaACLAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, aclID string) (*managedkafka.Acl, error) {
	logger.Default.Logf(t, "Getting settings for Managed Kafka ACL %s %s in %s in project %s", aclID, clusterID, location, projectID)

	service, err := NewManagedKafkaServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetManagedKafkaACLAttrsWithClient(ctx, service, projectID, location, clusterID, aclID)
}

// GetManagedKafkaACLAttrsWithClient returns the settings Google Cloud holds for the given Managed Kafka ACL using the supplied
// *managedkafka.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see managedkafka_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaACLAttrsWithClient(ctx context.Context, service *managedkafka.Service, projectID string, location string, clusterID string, aclID string) (*managedkafka.Acl, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s/acls/%s", projectID, location, clusterID, aclID)

	attrs, err := service.Projects.Locations.Clusters.Acls.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Managed Kafka ACL %s %s in %s in project %s does not exist", aclID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Managed Kafka ACL %s %s in %s in project %s: %w", aclID, clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetManagedKafkaConnectClusterAttrs returns the settings Google Cloud holds for the given Managed Kafka Connect cluster, so a test can assert on what was
// actually created rather than only that it exists.
// A Connect cluster is where connectors run, so its capacity and the Kafka cluster it attaches to are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaConnectClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *managedkafka.ConnectCluster {
	attrs, err := GetManagedKafkaConnectClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return attrs
}

// GetManagedKafkaConnectClusterAttrsE returns the settings Google Cloud holds for the given Managed Kafka Connect cluster.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaConnectClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*managedkafka.ConnectCluster, error) {
	logger.Default.Logf(t, "Getting settings for Managed Kafka Connect cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewManagedKafkaServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetManagedKafkaConnectClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetManagedKafkaConnectClusterAttrsWithClient returns the settings Google Cloud holds for the given Managed Kafka Connect cluster using the supplied
// *managedkafka.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see managedkafka_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaConnectClusterAttrsWithClient(ctx context.Context, service *managedkafka.Service, projectID string, location string, clusterID string) (*managedkafka.ConnectCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connectClusters/%s", projectID, location, clusterID)

	attrs, err := service.Projects.Locations.ConnectClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Managed Kafka Connect cluster %s in %s in project %s does not exist", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Managed Kafka Connect cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return attrs, nil
}

// GetManagedKafkaConnectorAttrs returns the settings Google Cloud holds for the given Managed Kafka connector, so a test can assert on what was
// actually created rather than only that it exists.
// A connector is what moves data in or out, so its configuration and its state decide whether anything is flowing.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaConnectorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, connectorID string) *managedkafka.Connector {
	attrs, err := GetManagedKafkaConnectorAttrsE(t, ctx, projectID, location, clusterID, connectorID)
	require.NoError(t, err)

	return attrs
}

// GetManagedKafkaConnectorAttrsE returns the settings Google Cloud holds for the given Managed Kafka connector.
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaConnectorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, connectorID string) (*managedkafka.Connector, error) {
	logger.Default.Logf(t, "Getting settings for Managed Kafka connector %s %s in %s in project %s", connectorID, clusterID, location, projectID)

	service, err := NewManagedKafkaServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetManagedKafkaConnectorAttrsWithClient(ctx, service, projectID, location, clusterID, connectorID)
}

// GetManagedKafkaConnectorAttrsWithClient returns the settings Google Cloud holds for the given Managed Kafka connector using the supplied
// *managedkafka.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see managedkafka_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetManagedKafkaConnectorAttrsWithClient(ctx context.Context, service *managedkafka.Service, projectID string, location string, clusterID string, connectorID string) (*managedkafka.Connector, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connectClusters/%s/connectors/%s", projectID, location, clusterID, connectorID)

	attrs, err := service.Projects.Locations.ConnectClusters.Connectors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Managed Kafka connector %s %s in %s in project %s does not exist", connectorID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Managed Kafka connector %s %s in %s in project %s: %w", connectorID, clusterID, location, projectID, err)
	}

	return attrs, nil
}

// NewManagedKafkaServiceE creates a Managed Service for Apache Kafka service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewManagedKafkaServiceE(t testing.TestingT, ctx context.Context) (*managedkafka.Service, error) {
	return managedkafka.NewService(ctx, append(withOptions(), option.WithScopes(managedkafka.CloudPlatformScope))...)
}
