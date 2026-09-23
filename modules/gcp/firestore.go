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

// GetFirestoreDatabaseAttrs returns the settings Google Cloud holds for the given Firestore
// database, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreDatabaseAttrs(t testing.TestingT, ctx context.Context, projectID string, databaseID string) *firestore.GoogleFirestoreAdminV1Database {
	database, err := GetFirestoreDatabaseAttrsE(t, ctx, projectID, databaseID)
	require.NoError(t, err)

	return database
}

// GetFirestoreDatabaseAttrsE returns the settings Google Cloud holds for the given Firestore
// database.
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreDatabaseAttrsE(t testing.TestingT, ctx context.Context, projectID string, databaseID string) (*firestore.GoogleFirestoreAdminV1Database, error) {
	logger.Default.Logf(t, "Getting settings for Firestore database %s in project %s", databaseID, projectID)

	service, err := NewFirestoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirestoreDatabaseAttrsWithClient(ctx, service, projectID, databaseID)
}

// GetFirestoreDatabaseAttrsWithClient returns the settings Google Cloud holds for the given
// Firestore database using the supplied *firestore.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see firestore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirestoreDatabaseAttrsWithClient(ctx context.Context, service *firestore.Service, projectID string, databaseID string) (*firestore.GoogleFirestoreAdminV1Database, error) {
	name := fmt.Sprintf("projects/%s/databases/%s", projectID, databaseID)

	database, err := service.Projects.Databases.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firestore database %s does not exist in project %s", databaseID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firestore database %s in project %s: %w", databaseID, projectID, err)
	}

	return database, nil
}

// NewFirestoreServiceE creates a Firestore service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewFirestoreServiceE(t testing.TestingT, ctx context.Context) (*firestore.Service, error) {
	return firestore.NewService(ctx, append(withOptions(), option.WithScopes(firestore.CloudPlatformScope))...)
}
