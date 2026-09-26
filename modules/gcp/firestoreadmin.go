package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/firestore/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetFirestoreIndexAttrs returns the settings Google Cloud holds for the given Firestore index, so
// a test can assert on what was actually created rather than only that it exists. Google assigns an
// index its id, so the caller passes the id it got back rather than a name it chose.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreIndexAttrs(t testing.TestingT, ctx context.Context, projectID string, databaseID string, collection string, indexID string) *firestore.GoogleFirestoreAdminV1Index {
	index, err := GetFirestoreIndexAttrsE(t, ctx, projectID, databaseID, collection, indexID)
	require.NoError(t, err)

	return index
}

// GetFirestoreIndexAttrsE returns the settings Google Cloud holds for the given Firestore index.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreIndexAttrsE(t testing.TestingT, ctx context.Context, projectID string, databaseID string, collection string, indexID string) (*firestore.GoogleFirestoreAdminV1Index, error) {
	logger.Default.Logf(t, "Getting settings for Firestore index %s on %s in database %s in project %s", indexID, collection, databaseID, projectID)

	service, err := NewFirestoreAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirestoreIndexAttrsWithClient(ctx, service, projectID, databaseID, collection, indexID)
}

// GetFirestoreIndexAttrsWithClient returns the settings Google Cloud holds for the given Firestore
// index using the supplied *firestore.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see firestoreadmin_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreIndexAttrsWithClient(ctx context.Context, service *firestore.Service, projectID string, databaseID string, collection string, indexID string) (*firestore.GoogleFirestoreAdminV1Index, error) {
	name := fmt.Sprintf("projects/%s/databases/%s/collectionGroups/%s/indexes/%s", projectID, databaseID, collection, indexID)

	index, err := service.Projects.Databases.CollectionGroups.Indexes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firestore index %s does not exist on %s in database %s in project %s", indexID, collection, databaseID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firestore index %s on %s in database %s in project %s: %w", indexID, collection, databaseID, projectID, err)
	}

	return index, nil
}

// GetFirestoreFieldAttrs returns the settings Google Cloud holds for the given Firestore field, so
// a test can assert on what was actually configured. A field exists whether or not anyone has
// configured it, so this reads the overrides rather than the field's existence.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreFieldAttrs(t testing.TestingT, ctx context.Context, projectID string, databaseID string, collection string, fieldPath string) *firestore.GoogleFirestoreAdminV1Field {
	field, err := GetFirestoreFieldAttrsE(t, ctx, projectID, databaseID, collection, fieldPath)
	require.NoError(t, err)

	return field
}

// GetFirestoreFieldAttrsE returns the settings Google Cloud holds for the given Firestore field.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreFieldAttrsE(t testing.TestingT, ctx context.Context, projectID string, databaseID string, collection string, fieldPath string) (*firestore.GoogleFirestoreAdminV1Field, error) {
	logger.Default.Logf(t, "Getting settings for Firestore field %s on %s in database %s in project %s", fieldPath, collection, databaseID, projectID)

	service, err := NewFirestoreAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirestoreFieldAttrsWithClient(ctx, service, projectID, databaseID, collection, fieldPath)
}

// GetFirestoreFieldAttrsWithClient returns the settings Google Cloud holds for the given Firestore
// field using the supplied *firestore.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see firestoreadmin_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreFieldAttrsWithClient(ctx context.Context, service *firestore.Service, projectID string, databaseID string, collection string, fieldPath string) (*firestore.GoogleFirestoreAdminV1Field, error) {
	name := fmt.Sprintf("projects/%s/databases/%s/collectionGroups/%s/fields/%s", projectID, databaseID, collection, fieldPath)

	field, err := service.Projects.Databases.CollectionGroups.Fields.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firestore field %s does not exist on %s in database %s in project %s", fieldPath, collection, databaseID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firestore field %s on %s in database %s in project %s: %w", fieldPath, collection, databaseID, projectID, err)
	}

	return field, nil
}

// GetFirestoreDocumentAttrs returns the contents Google Cloud holds for the given Firestore
// document, so a test can assert on what was actually written rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreDocumentAttrs(t testing.TestingT, ctx context.Context, projectID string, databaseID string, collection string, documentID string) *firestore.Document {
	document, err := GetFirestoreDocumentAttrsE(t, ctx, projectID, databaseID, collection, documentID)
	require.NoError(t, err)

	return document
}

// GetFirestoreDocumentAttrsE returns the contents Google Cloud holds for the given Firestore
// document.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreDocumentAttrsE(t testing.TestingT, ctx context.Context, projectID string, databaseID string, collection string, documentID string) (*firestore.Document, error) {
	logger.Default.Logf(t, "Getting the contents of Firestore document %s in %s in database %s in project %s", documentID, collection, databaseID, projectID)

	service, err := NewFirestoreAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirestoreDocumentAttrsWithClient(ctx, service, projectID, databaseID, collection, documentID)
}

// GetFirestoreDocumentAttrsWithClient returns the contents Google Cloud holds for the given
// Firestore document using the supplied *firestore.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see firestoreadmin_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreDocumentAttrsWithClient(ctx context.Context, service *firestore.Service, projectID string, databaseID string, collection string, documentID string) (*firestore.Document, error) {
	name := fmt.Sprintf("projects/%s/databases/%s/documents/%s/%s", projectID, databaseID, collection, documentID)

	document, err := service.Projects.Databases.Documents.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firestore document %s does not exist in %s in database %s in project %s", documentID, collection, databaseID, projectID)
		}

		return nil, fmt.Errorf("failed to get the contents of Firestore document %s in %s in database %s in project %s: %w", documentID, collection, databaseID, projectID, err)
	}

	return document, nil
}

// NewFirestoreAdminServiceE creates a Firestore service authenticated the same way every other
// client in this module is. It serves both the admin calls and the document calls.
// The ctx parameter supports cancellation and timeouts.
func NewFirestoreAdminServiceE(t testing.TestingT, ctx context.Context) (*firestore.Service, error) {
	return firestore.NewService(ctx, append(withOptions(), option.WithScopes(firestore.CloudPlatformScope))...)
}
