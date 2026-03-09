package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"google.golang.org/api/sqladmin/v1"
)

// AssertCloudSQLInstanceExists checks if the given Cloud SQL instance exists and fails the test if it does not.
func AssertCloudSQLInstanceExists(t testing.TestingT, projectID string, instanceName string) {
	err := AssertCloudSQLInstanceExistsE(t, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertCloudSQLInstanceExistsE checks if the given Cloud SQL instance exists and returns an error if it does not.
func AssertCloudSQLInstanceExistsE(t testing.TestingT, projectID string, instanceName string) error {
	logger.Default.Logf(t, "Verifying Cloud SQL instance %s exists in project %s", instanceName, projectID)

	_, err := getCloudSQLInstanceE(projectID, instanceName)
	return err
}

// GetCloudSQLInstanceDatabaseVersion returns the database version of the given Cloud SQL instance (e.g. MYSQL_8_0, POSTGRES_14, SQLSERVER_2019_STANDARD).
func GetCloudSQLInstanceDatabaseVersion(t testing.TestingT, projectID string, instanceName string) string {
	version, err := GetCloudSQLInstanceDatabaseVersionE(t, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}
	return version
}

// GetCloudSQLInstanceDatabaseVersionE returns the database version of the given Cloud SQL instance and returns an error if it could not be retrieved.
func GetCloudSQLInstanceDatabaseVersionE(t testing.TestingT, projectID string, instanceName string) (string, error) {
	logger.Default.Logf(t, "Getting database version of Cloud SQL instance %s in project %s", instanceName, projectID)

	instance, err := getCloudSQLInstanceE(projectID, instanceName)
	if err != nil {
		return "", err
	}

	return instance.DatabaseVersion, nil
}

// getCloudSQLInstanceE is a helper that fetches the Cloud SQL instance details.
func getCloudSQLInstanceE(projectID string, instanceName string) (*sqladmin.DatabaseInstance, error) {
	service, err := newCloudSQLService()
	if err != nil {
		return nil, err
	}

	instance, err := service.Instances.Get(projectID, instanceName).Do()
	if err != nil {
		return nil, fmt.Errorf("cloud SQL instance %s does not exist in project %s: %v", instanceName, projectID, err)
	}

	return instance, nil
}

// newCloudSQLService creates a new Cloud SQL Admin service using the global GCP auth options.
func newCloudSQLService() (*sqladmin.Service, error) {
	ctx := context.Background()
	service, err := sqladmin.NewService(ctx, withOptions()...)
	if err != nil {
		return nil, err
	}
	return service, nil
}
