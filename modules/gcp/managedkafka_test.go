package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/managedkafka/v1"
	"google.golang.org/api/option"
)

// newFakeManagedKafkaService points a real Managed Service for Apache Kafka client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeManagedKafkaService(t *testing.T, handler http.Handler) *managedkafka.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := managedkafka.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetManagedKafkaClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Managed Kafka cluster the terraform-google-messaging Managed Kafka cluster module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test","state":"ACTIVE","capacityConfig":{"vcpuCount":"3","memoryBytes":"3221225472"},"labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetManagedKafkaClusterAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetManagedKafkaClusterAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Managed Kafka cluster that is not there should read a sentence about that Managed Kafka cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetManagedKafkaClusterAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetManagedKafkaTopicAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Managed Kafka topic the terraform-google-messaging Managed Kafka topic module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-parent/topics/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-parent/topics/gw-library-test","partitionCount":3,"replicationFactor":3}`))
	})

	attrs, err := gcp.GetManagedKafkaTopicAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), attrs.PartitionCount)
	assert.Equal(t, int64(3), attrs.ReplicationFactor)
}

func TestGetManagedKafkaTopicAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Managed Kafka topic that is not there should read a sentence about that Managed Kafka topic, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetManagedKafkaTopicAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetManagedKafkaACLAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Managed Kafka ACL the terraform-google-messaging Managed Kafka ACL module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-parent/acls/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-parent/acls/gw-library-test","aclEntries":[{"principal":"User:terratest","operation":"READ","permissionType":"ALLOW","host":"*"}],"etag":"BwXhqw=="}`))
	})

	attrs, err := gcp.GetManagedKafkaACLAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "User:terratest", attrs.AclEntries[0].Principal)
	assert.Equal(t, "READ", attrs.AclEntries[0].Operation)
}

func TestGetManagedKafkaACLAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Managed Kafka ACL that is not there should read a sentence about that Managed Kafka ACL, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetManagedKafkaACLAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetManagedKafkaConnectClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Managed Kafka Connect cluster the terraform-google-messaging Managed Kafka Connect cluster module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connectClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connectClusters/gw-library-test","kafkaCluster":"projects/gw-library-test-project/locations/us-central1/clusters/gw-library-parent","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetManagedKafkaConnectClusterAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/clusters/gw-library-parent", attrs.KafkaCluster)
}

func TestGetManagedKafkaConnectClusterAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Managed Kafka Connect cluster that is not there should read a sentence about that Managed Kafka Connect cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetManagedKafkaConnectClusterAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetManagedKafkaConnectorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Managed Kafka connector the terraform-google-messaging Managed Kafka connector module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connectClusters/gw-library-parent/connectors/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connectClusters/gw-library-parent/connectors/gw-library-test","state":"RUNNING","configs":{"connector.class":"io.confluent.connect.gcp.pubsub.PubSubSinkConnector"}}`))
	})

	attrs, err := gcp.GetManagedKafkaConnectorAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "RUNNING", attrs.State)
}

func TestGetManagedKafkaConnectorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Managed Kafka connector that is not there should read a sentence about that Managed Kafka connector, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetManagedKafkaConnectorAttrsWithClient(context.Background(), newFakeManagedKafkaService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
