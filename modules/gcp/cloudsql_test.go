//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertCloudSQLInstanceExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	instanceName := fmt.Sprintf("terratest-cloudsql-%s", random.UniqueId())
	logger.Logf(t, "Creating Cloud SQL instance %s to verify existence check works", instanceName)

	CreateCloudSQLInstanceContext(t, context.Background(), projectID, instanceName, "MYSQL_8_0")
	defer DeleteCloudSQLInstanceContext(t, context.Background(), projectID, instanceName)

	AssertCloudSQLInstanceExistsContext(t, context.Background(), projectID, instanceName)

	version := GetCloudSQLInstanceDatabaseVersionContext(t, context.Background(), projectID, instanceName)
	assert.Equal(t, "MYSQL_8_0", version)
}

func TestAssertCloudSQLInstanceExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	instanceName := fmt.Sprintf("terratest-cloudsql-%s", random.UniqueId())
	logger.Logf(t, "Checking that non-existent Cloud SQL instance %s returns an error", instanceName)

	err := AssertCloudSQLInstanceExistsContextE(t, context.Background(), projectID, instanceName)
	require.Error(t, err, "Expected an error for non-existent Cloud SQL instance, but got none")
}

func TestGetCloudSQLInstanceDatabaseVersionNoFalsePositive(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	instanceName := fmt.Sprintf("terratest-cloudsql-%s", random.UniqueId())
	logger.Logf(t, "Checking that GetCloudSQLInstanceDatabaseVersionContextE returns an error for non-existent instance %s", instanceName)

	_, err := GetCloudSQLInstanceDatabaseVersionContextE(t, context.Background(), projectID, instanceName)
	require.Error(t, err, "Expected an error for non-existent Cloud SQL instance, but got none")
}
