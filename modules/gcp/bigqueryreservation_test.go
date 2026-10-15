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

func TestGetBigQueryReservationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a BigQuery reservation a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/reservations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"slotCapacity":"0","edition":"ENTERPRISE","ignoreIdleSlots":true,"autoscale":{"maxSlots":"100","currentSlots":"0"}}`))
	})

	result, err := gcp.GetBigQueryReservationAttrsWithClient(context.Background(), newFakeBigQueryReservationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(0), result.SlotCapacity)
	assert.Equal(t, "ENTERPRISE", result.Edition)
	assert.True(t, result.IgnoreIdleSlots)
	require.NotNil(t, result.Autoscale)
	assert.Equal(t, int64(100), result.Autoscale.MaxSlots)
}

func TestGetBigQueryReservationAttrsWithClientReportsAMissingBigqueryReservation(t *testing.T) {
	t.Parallel()

	// A caller who names a BigQuery reservation that is not there should read a sentence about that BigQuery reservation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetBigQueryReservationAttrsWithClient(context.Background(), newFakeBigQueryReservationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetBigQueryCapacityCommitmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a commitment the capacity commitment
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/capacityCommitments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/capacityCommitments/gw-library-test","slotCount":"100","plan":"ANNUAL","renewalPlan":"FLEX","state":"ACTIVE","edition":"ENTERPRISE"}`))
	})

	commitment, err := gcp.GetBigQueryCapacityCommitmentAttrsWithClient(context.Background(), newFakeBigQueryReservationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(100), commitment.SlotCount)
	assert.Equal(t, "ANNUAL", commitment.Plan)
	assert.Equal(t, "FLEX", commitment.RenewalPlan)
	assert.Equal(t, "ENTERPRISE", commitment.Edition)
}

func TestGetBigQueryCapacityCommitmentAttrsWithClientReportsAMissingCommitment(t *testing.T) {
	t.Parallel()

	// A caller who asks for a commitment that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetBigQueryCapacityCommitmentAttrsWithClient(context.Background(), newFakeBigQueryReservationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
