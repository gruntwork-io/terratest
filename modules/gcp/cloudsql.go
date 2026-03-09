package gcp

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"google.golang.org/api/sqladmin/v1"
)

// AssertCloudSQLInstanceExistsContext checks if the given Cloud SQL instance exists and fails the test if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertCloudSQLInstanceExistsContext(t testing.TestingT, ctx context.Context, projectID string, instanceName string) {
	err := AssertCloudSQLInstanceExistsContextE(t, ctx, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertCloudSQLInstanceExistsContextE checks if the given Cloud SQL instance exists and returns an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertCloudSQLInstanceExistsContextE(t testing.TestingT, ctx context.Context, projectID string, instanceName string) error {
	logger.Default.Logf(t, "Verifying Cloud SQL instance %s exists in project %s", instanceName, projectID)

	_, err := getCloudSQLInstanceE(ctx, projectID, instanceName)
	return err
}

// GetCloudSQLInstanceDatabaseVersionContext returns the database version of the given Cloud SQL instance (e.g. MYSQL_8_0, POSTGRES_14, SQLSERVER_2019_STANDARD).
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLInstanceDatabaseVersionContext(t testing.TestingT, ctx context.Context, projectID string, instanceName string) string {
	version, err := GetCloudSQLInstanceDatabaseVersionContextE(t, ctx, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}
	return version
}

// GetCloudSQLInstanceDatabaseVersionContextE returns the database version of the given Cloud SQL instance and returns an error if it could not be retrieved.
// The ctx parameter supports cancellation and timeouts.
func GetCloudSQLInstanceDatabaseVersionContextE(t testing.TestingT, ctx context.Context, projectID string, instanceName string) (string, error) {
	logger.Default.Logf(t, "Getting database version of Cloud SQL instance %s in project %s", instanceName, projectID)

	instance, err := getCloudSQLInstanceE(ctx, projectID, instanceName)
	if err != nil {
		return "", err
	}

	return instance.DatabaseVersion, nil
}

// CreateCloudSQLInstanceContext creates a new Cloud SQL instance with the given database version and waits for it to be ready.
// The ctx parameter supports cancellation and timeouts.
func CreateCloudSQLInstanceContext(t testing.TestingT, ctx context.Context, projectID string, instanceName string, databaseVersion string) {
	err := CreateCloudSQLInstanceContextE(t, ctx, projectID, instanceName, databaseVersion)
	if err != nil {
		t.Fatal(err)
	}
}

// CreateCloudSQLInstanceContextE creates a new Cloud SQL instance and returns an error if it fails.
// The ctx parameter supports cancellation and timeouts.
func CreateCloudSQLInstanceContextE(t testing.TestingT, ctx context.Context, projectID string, instanceName string, databaseVersion string) error {
	logger.Default.Logf(t, "Creating Cloud SQL instance %s in project %s", instanceName, projectID)

	service, err := newCloudSQLService(ctx)
	if err != nil {
		return err
	}

	instance := &sqladmin.DatabaseInstance{
		Name:            instanceName,
		DatabaseVersion: databaseVersion,
		Settings: &sqladmin.Settings{
			Tier: "db-f1-micro",
		},
	}

	op, err := service.Instances.Insert(projectID, instance).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create Cloud SQL instance %s in project %s: %w", instanceName, projectID, err)
	}

	return waitForCloudSQLOperation(t, ctx, service, projectID, op.Name)
}

// DeleteCloudSQLInstanceContext deletes the given Cloud SQL instance and waits for the operation to complete.
// The ctx parameter supports cancellation and timeouts.
func DeleteCloudSQLInstanceContext(t testing.TestingT, ctx context.Context, projectID string, instanceName string) {
	err := DeleteCloudSQLInstanceContextE(t, ctx, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteCloudSQLInstanceContextE deletes the given Cloud SQL instance and returns an error if it fails.
// The ctx parameter supports cancellation and timeouts.
func DeleteCloudSQLInstanceContextE(t testing.TestingT, ctx context.Context, projectID string, instanceName string) error {
	logger.Default.Logf(t, "Deleting Cloud SQL instance %s in project %s", instanceName, projectID)

	service, err := newCloudSQLService(ctx)
	if err != nil {
		return err
	}

	op, err := service.Instances.Delete(projectID, instanceName).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to delete Cloud SQL instance %s in project %s: %w", instanceName, projectID, err)
	}

	return waitForCloudSQLOperation(t, ctx, service, projectID, op.Name)
}

// getCloudSQLInstanceE is a helper that fetches the Cloud SQL instance details.
func getCloudSQLInstanceE(ctx context.Context, projectID string, instanceName string) (*sqladmin.DatabaseInstance, error) {
	service, err := newCloudSQLService(ctx)
	if err != nil {
		return nil, err
	}

	instance, err := service.Instances.Get(projectID, instanceName).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get Cloud SQL instance %s in project %s: %w", instanceName, projectID, err)
	}

	return instance, nil
}

// newCloudSQLService creates a new Cloud SQL Admin service using the global GCP auth options.
func newCloudSQLService(ctx context.Context) (*sqladmin.Service, error) {
	service, err := sqladmin.NewService(ctx, withOptions()...)
	if err != nil {
		return nil, err
	}
	return service, nil
}

// waitForCloudSQLOperation polls until the given Cloud SQL operation completes using the terratest retry module.
func waitForCloudSQLOperation(t testing.TestingT, ctx context.Context, service *sqladmin.Service, projectID string, operationName string) error {
	// Cloud SQL instances can take up to 20 minutes to provision; 120 retries x 15s = 30 minutes to be safe.
	maxRetries := 120
	sleepBetweenRetries := 15 * time.Second

	_, err := retry.DoWithRetryE(t, fmt.Sprintf("Waiting for Cloud SQL operation %s", operationName), maxRetries, sleepBetweenRetries, func() (string, error) {
		op, err := service.Operations.Get(projectID, operationName).Context(ctx).Do()
		if err != nil {
			return "", fmt.Errorf("failed to get Cloud SQL operation status: %w", err)
		}

		if op.Status != "DONE" {
			return "", fmt.Errorf("operation %s not done yet, current status: %s", operationName, op.Status)
		}

		if op.Error != nil {
			// Operation is DONE but failed — stop retrying immediately.
			return "", retry.FatalError{Underlying: fmt.Errorf("Cloud SQL operation failed: %v", op.Error.Errors)}
		}

		return "DONE", nil
	})

	return err
}
