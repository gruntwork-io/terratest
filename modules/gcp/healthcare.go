package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/healthcare/v1"
	"google.golang.org/api/option"
)

// GetHealthcareDatasetAttrs returns the settings Google Cloud holds for the given Cloud Healthcare
// dataset, so a test can assert on what was actually created rather than only that it exists. A
// dataset is regional, so the location it was created in has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDatasetAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string) *healthcare.Dataset {
	dataset, err := GetHealthcareDatasetAttrsE(t, ctx, projectID, location, datasetID)
	require.NoError(t, err)

	return dataset
}

// GetHealthcareDatasetAttrsE returns the settings Google Cloud holds for the given Cloud Healthcare
// dataset.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDatasetAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string) (*healthcare.Dataset, error) {
	logger.Default.Logf(t, "Getting settings for healthcare dataset %s in project %s", datasetID, projectID)

	service, err := NewHealthcareServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcareDatasetAttrsWithClient(ctx, service, projectID, location, datasetID)
}

// GetHealthcareDatasetAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Healthcare dataset using the supplied *healthcare.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see healthcare_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDatasetAttrsWithClient(ctx context.Context, service *healthcare.Service, projectID string, location string, datasetID string) (*healthcare.Dataset, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s", projectID, location, datasetID)

	dataset, err := service.Projects.Locations.Datasets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("healthcare dataset %s does not exist in project %s", datasetID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for healthcare dataset %s in project %s: %w", datasetID, projectID, err)
	}

	return dataset, nil
}

// NewHealthcareServiceE creates a Cloud Healthcare service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewHealthcareServiceE(t testing.TestingT, ctx context.Context) (*healthcare.Service, error) {
	return healthcare.NewService(ctx, append(withOptions(), option.WithScopes(healthcare.CloudPlatformScope))...)
}

// GetHealthcareDicomStoreAttrs returns the settings Google Cloud holds for the given DICOM store, so a test can assert on
// what was actually created rather than only that it exists. A store lives in a healthcare dataset, so it is named by
// the dataset's location and id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDicomStoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) *healthcare.DicomStore {
	store, err := GetHealthcareDicomStoreAttrsE(t, ctx, projectID, location, datasetID, storeID)
	require.NoError(t, err)

	return store
}

// GetHealthcareDicomStoreAttrsE returns the settings Google Cloud holds for the given DICOM store.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDicomStoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) (*healthcare.DicomStore, error) {
	logger.Default.Logf(t, "Getting settings for DICOM store %s in dataset %s in %s in project %s", storeID, datasetID, location, projectID)

	service, err := NewHealthcareServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcareDicomStoreAttrsWithClient(ctx, service, projectID, location, datasetID, storeID)
}

// GetHealthcareDicomStoreAttrsWithClient returns the settings Google Cloud holds for the given DICOM store using the supplied
// *healthcare.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see healthcare_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareDicomStoreAttrsWithClient(ctx context.Context, service *healthcare.Service, projectID string, location string, datasetID string, storeID string) (*healthcare.DicomStore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/dicomStores/%s", projectID, location, datasetID, storeID)

	store, err := service.Projects.Locations.Datasets.DicomStores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the DICOM store %s does not exist in dataset %s in %s in project %s", storeID, datasetID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for DICOM store %s in dataset %s in %s in project %s: %w", storeID, datasetID, location, projectID, err)
	}

	return store, nil
}

// GetHealthcareFhirStoreAttrs returns the settings Google Cloud holds for the given FHIR store, so a test can assert on
// what was actually created rather than only that it exists. A store lives in a healthcare dataset, so it is named by
// the dataset's location and id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareFhirStoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) *healthcare.FhirStore {
	store, err := GetHealthcareFhirStoreAttrsE(t, ctx, projectID, location, datasetID, storeID)
	require.NoError(t, err)

	return store
}

// GetHealthcareFhirStoreAttrsE returns the settings Google Cloud holds for the given FHIR store.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareFhirStoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) (*healthcare.FhirStore, error) {
	logger.Default.Logf(t, "Getting settings for FHIR store %s in dataset %s in %s in project %s", storeID, datasetID, location, projectID)

	service, err := NewHealthcareServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcareFhirStoreAttrsWithClient(ctx, service, projectID, location, datasetID, storeID)
}

// GetHealthcareFhirStoreAttrsWithClient returns the settings Google Cloud holds for the given FHIR store using the supplied
// *healthcare.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see healthcare_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareFhirStoreAttrsWithClient(ctx context.Context, service *healthcare.Service, projectID string, location string, datasetID string, storeID string) (*healthcare.FhirStore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/fhirStores/%s", projectID, location, datasetID, storeID)

	store, err := service.Projects.Locations.Datasets.FhirStores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the FHIR store %s does not exist in dataset %s in %s in project %s", storeID, datasetID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for FHIR store %s in dataset %s in %s in project %s: %w", storeID, datasetID, location, projectID, err)
	}

	return store, nil
}

// GetHealthcareHl7V2StoreAttrs returns the settings Google Cloud holds for the given HL7v2 store, so a test can assert on
// what was actually created rather than only that it exists. A store lives in a healthcare dataset, so it is named by
// the dataset's location and id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareHl7V2StoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) *healthcare.Hl7V2Store {
	store, err := GetHealthcareHl7V2StoreAttrsE(t, ctx, projectID, location, datasetID, storeID)
	require.NoError(t, err)

	return store
}

// GetHealthcareHl7V2StoreAttrsE returns the settings Google Cloud holds for the given HL7v2 store.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareHl7V2StoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) (*healthcare.Hl7V2Store, error) {
	logger.Default.Logf(t, "Getting settings for HL7v2 store %s in dataset %s in %s in project %s", storeID, datasetID, location, projectID)

	service, err := NewHealthcareServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcareHl7V2StoreAttrsWithClient(ctx, service, projectID, location, datasetID, storeID)
}

// GetHealthcareHl7V2StoreAttrsWithClient returns the settings Google Cloud holds for the given HL7v2 store using the supplied
// *healthcare.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see healthcare_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareHl7V2StoreAttrsWithClient(ctx context.Context, service *healthcare.Service, projectID string, location string, datasetID string, storeID string) (*healthcare.Hl7V2Store, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/hl7V2Stores/%s", projectID, location, datasetID, storeID)

	store, err := service.Projects.Locations.Datasets.Hl7V2Stores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the HL7v2 store %s does not exist in dataset %s in %s in project %s", storeID, datasetID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for HL7v2 store %s in dataset %s in %s in project %s: %w", storeID, datasetID, location, projectID, err)
	}

	return store, nil
}

// GetHealthcareConsentStoreAttrs returns the settings Google Cloud holds for the given consent store, so a test can assert on
// what was actually created rather than only that it exists. A store lives in a healthcare dataset, so it is named by
// the dataset's location and id as well as its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareConsentStoreAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) *healthcare.ConsentStore {
	store, err := GetHealthcareConsentStoreAttrsE(t, ctx, projectID, location, datasetID, storeID)
	require.NoError(t, err)

	return store
}

// GetHealthcareConsentStoreAttrsE returns the settings Google Cloud holds for the given consent store.
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareConsentStoreAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, datasetID string, storeID string) (*healthcare.ConsentStore, error) {
	logger.Default.Logf(t, "Getting settings for consent store %s in dataset %s in %s in project %s", storeID, datasetID, location, projectID)

	service, err := NewHealthcareServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetHealthcareConsentStoreAttrsWithClient(ctx, service, projectID, location, datasetID, storeID)
}

// GetHealthcareConsentStoreAttrsWithClient returns the settings Google Cloud holds for the given consent store using the supplied
// *healthcare.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see healthcare_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetHealthcareConsentStoreAttrsWithClient(ctx context.Context, service *healthcare.Service, projectID string, location string, datasetID string, storeID string) (*healthcare.ConsentStore, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/datasets/%s/consentStores/%s", projectID, location, datasetID, storeID)

	store, err := service.Projects.Locations.Datasets.ConsentStores.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the consent store %s does not exist in dataset %s in %s in project %s", storeID, datasetID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for consent store %s in dataset %s in %s in project %s: %w", storeID, datasetID, location, projectID, err)
	}

	return store, nil
}
