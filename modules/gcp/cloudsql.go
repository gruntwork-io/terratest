package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sqladmin/v1"
)

// GetCloudSQLInstanceAttrs returns the settings Google Cloud holds for the given Cloud SQL
// instance, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string) *sqladmin.DatabaseInstance {
	instance, err := GetCloudSQLInstanceAttrsE(t, ctx, projectID, instanceID)
	require.NoError(t, err)

	return instance
}

// GetCloudSQLInstanceAttrsE returns the settings Google Cloud holds for the given Cloud SQL
// instance.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string) (*sqladmin.DatabaseInstance, error) {
	logger.Default.Logf(t, "Getting settings for Cloud SQL instance %s in project %s", instanceID, projectID)

	service, err := NewCloudSQLServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudSQLInstanceAttrsWithClient(ctx, service, projectID, instanceID)
}

// GetCloudSQLInstanceAttrsWithClient returns the settings Google Cloud holds for the given Cloud SQL
// instance using the supplied *sqladmin.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see cloudsql_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLInstanceAttrsWithClient(ctx context.Context, service *sqladmin.Service, projectID string, instanceID string) (*sqladmin.DatabaseInstance, error) {
	instance, err := service.Instances.Get(projectID, instanceID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud SQL instance %s does not exist in project %s", instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud SQL instance %s in project %s: %w", instanceID, projectID, err)
	}

	return instance, nil
}

// GetCloudSQLDatabaseAttrs returns the settings Google Cloud holds for the given database on a Cloud
// SQL instance, so a test can assert on what was actually created rather than only that it exists. A
// database is named by its instance as well as its own name.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLDatabaseAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, databaseName string) *sqladmin.Database {
	database, err := GetCloudSQLDatabaseAttrsE(t, ctx, projectID, instanceID, databaseName)
	require.NoError(t, err)

	return database
}

// GetCloudSQLDatabaseAttrsE returns the settings Google Cloud holds for the given database on a
// Cloud SQL instance.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLDatabaseAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, databaseName string) (*sqladmin.Database, error) {
	logger.Default.Logf(t, "Getting settings for database %s on Cloud SQL instance %s in project %s", databaseName, instanceID, projectID)

	service, err := NewCloudSQLServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudSQLDatabaseAttrsWithClient(ctx, service, projectID, instanceID, databaseName)
}

// GetCloudSQLDatabaseAttrsWithClient returns the settings Google Cloud holds for the given database
// on a Cloud SQL instance using the supplied *sqladmin.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see cloudsql_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLDatabaseAttrsWithClient(ctx context.Context, service *sqladmin.Service, projectID string, instanceID string, databaseName string) (*sqladmin.Database, error) {
	database, err := service.Databases.Get(projectID, instanceID, databaseName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the database %s does not exist on Cloud SQL instance %s in project %s", databaseName, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for database %s on Cloud SQL instance %s in project %s: %w", databaseName, instanceID, projectID, err)
	}

	return database, nil
}

// GetCloudSQLUserAttrs returns the settings Google Cloud holds for the given user on a Cloud SQL
// instance, so a test can assert on what was actually created rather than only that it exists. On
// MySQL a user is named by its host as well as its name, and passing the wrong host, or none, reads
// as absent; on PostgreSQL and SQL Server there is no host and the argument is empty. A user's
// password is never returned, so nothing here can leak one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLUserAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, userName string, host string) *sqladmin.User {
	user, err := GetCloudSQLUserAttrsE(t, ctx, projectID, instanceID, userName, host)
	require.NoError(t, err)

	return user
}

// GetCloudSQLUserAttrsE returns the settings Google Cloud holds for the given user on a Cloud SQL
// instance.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLUserAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, userName string, host string) (*sqladmin.User, error) {
	logger.Default.Logf(t, "Getting settings for user %s at host %s on Cloud SQL instance %s in project %s", userName, host, instanceID, projectID)

	service, err := NewCloudSQLServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudSQLUserAttrsWithClient(ctx, service, projectID, instanceID, userName, host)
}

// GetCloudSQLUserAttrsWithClient returns the settings Google Cloud holds for the given user on a
// Cloud SQL instance using the supplied *sqladmin.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see cloudsql_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLUserAttrsWithClient(ctx context.Context, service *sqladmin.Service, projectID string, instanceID string, userName string, host string) (*sqladmin.User, error) {
	user, err := service.Users.Get(projectID, instanceID, userName).Host(host).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the user %s at host %s does not exist on Cloud SQL instance %s in project %s", userName, host, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for user %s at host %s on Cloud SQL instance %s in project %s: %w", userName, host, instanceID, projectID, err)
	}

	return user, nil
}

// NewCloudSQLServiceE creates a Cloud SQL Admin service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudSQLServiceE(t testing.TestingT, ctx context.Context) (*sqladmin.Service, error) {
	return sqladmin.NewService(ctx, append(withOptions(), option.WithScopes(sqladmin.CloudPlatformScope))...)
}
