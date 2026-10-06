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
	"google.golang.org/api/bigqueryreservation/v1"
	"google.golang.org/api/option"
)

// newFakeBigQueryReservationService points a real BigQuery Reservation client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigQueryReservationService(t *testing.T, handler http.Handler) *bigqueryreservation.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigqueryreservation.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigQueryBiReservationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for the reservation the terraform-google-data-analytics BI reservation module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/biReservation"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/biReservation","size":"1073741824","updateTime":"2026-09-24T12:00:00Z"}`))
	})

	reservation, err := gcp.GetBigQueryBiReservationAttrsWithClient(context.Background(), newFakeBigQueryReservationService(t, handler), "gw-library-test-project", "us-central1")
	require.NoError(t, err)

	assert.Equal(t, int64(1073741824), reservation.Size)
}
