package gcp_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllGCPRegionsWithClient(t *testing.T) {
	t.Parallel()

	svc := newFakeComputeService(t, respond(t, http.MethodGet, "/projects/p/regions", http.StatusOK, `{"items":[{"name":"us-central1"},{"name":"us-east1"}]}`))

	got, err := gcp.GetAllGCPRegionsWithClient(context.Background(), svc, "p")
	require.NoError(t, err)
	assert.Equal(t, []string{"us-central1", "us-east1"}, got)
}

func TestGetAllGCPZonesWithClient(t *testing.T) {
	t.Parallel()

	svc := newFakeComputeService(t, respond(t, http.MethodGet, "/projects/p/zones", http.StatusOK, `{"items":[{"name":"us-central1-a"},{"name":"us-east1-b"}]}`))

	got, err := gcp.GetAllGCPZonesWithClient(context.Background(), svc, "p")
	require.NoError(t, err)
	assert.Equal(t, []string{"us-central1-a", "us-east1-b"}, got)
}
