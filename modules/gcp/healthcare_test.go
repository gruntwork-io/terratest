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
	"google.golang.org/api/healthcare/v1"
	"google.golang.org/api/option"
)

// newFakeHealthcareService points a real *healthcare.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeHealthcareService(t *testing.T, handler http.Handler) *healthcare.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := healthcare.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetHealthcareDatasetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-healthcare dataset module sets, because the point
	// of reading settings back is asserting a module configured the dataset it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/us-central1/datasets/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test",
			"timeZone":"America/New_York"
		}`))
	})

	dataset, err := gcp.GetHealthcareDatasetAttrsWithClient(context.Background(), newFakeHealthcareService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test", dataset.Name)
	assert.Equal(t, "America/New_York", dataset.TimeZone)
}

func TestGetHealthcareDatasetAttrsWithClientMissingDataset(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the dataset and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetHealthcareDatasetAttrsWithClient(context.Background(), newFakeHealthcareService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetHealthcareDicomStoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a store the terraform-google-healthcare DICOM store module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/dicomStores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/dicomStores/gw-library-test","labels":{"purpose":"terratest"}}`))
	})

	store, err := gcp.GetHealthcareDicomStoreAttrsWithClient(context.Background(), newFakeHealthcareService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest", store.Labels["purpose"])
}

func TestGetHealthcareFhirStoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a store the terraform-google-healthcare FHIR store module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/fhirStores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/fhirStores/gw-library-test","version":"R4","enableUpdateCreate":true,"disableReferentialIntegrity":true,"labels":{"purpose":"terratest"}}`))
	})

	store, err := gcp.GetHealthcareFhirStoreAttrsWithClient(context.Background(), newFakeHealthcareService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "R4", store.Version)
	assert.True(t, store.EnableUpdateCreate)
	assert.True(t, store.DisableReferentialIntegrity)
}

func TestGetHealthcareHl7V2StoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a store the terraform-google-healthcare HL7v2 store module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/hl7V2Stores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/hl7V2Stores/gw-library-test","rejectDuplicateMessage":true,"parserConfig":{"version":"V3","allowNullHeader":true},"labels":{"purpose":"terratest"}}`))
	})

	store, err := gcp.GetHealthcareHl7V2StoreAttrsWithClient(context.Background(), newFakeHealthcareService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.True(t, store.RejectDuplicateMessage)
	require.NotNil(t, store.ParserConfig)
	assert.Equal(t, "V3", store.ParserConfig.Version)
}

func TestGetHealthcareConsentStoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a store the terraform-google-healthcare consent store module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/consentStores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/datasets/gw-library-test/consentStores/gw-library-test","defaultConsentTtl":"90000s","enableConsentCreateOnUpdate":true,"labels":{"purpose":"terratest"}}`))
	})

	store, err := gcp.GetHealthcareConsentStoreAttrsWithClient(context.Background(), newFakeHealthcareService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "90000s", store.DefaultConsentTtl)
	assert.True(t, store.EnableConsentCreateOnUpdate)
}
