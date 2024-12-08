// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"strings"
	"testing"

	"cloud.google.com/go/bigquery"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestAssertDatasetExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	id := random.UniqueId()
	datasetName := "gruntwork-terratest-" + strings.ToLower(id)
	logger.Logf(t, "Random values selected Id = %s\n", id)

	CreateDataset(t, projectID, datasetName, bigquery.DatasetMetadata{})
	defer DeleteDataset(t, projectID, datasetName)

	AssertStorageBucketExists(t, datasetName)
}

func TestAssertDatasetExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	id := random.UniqueId()
	datasetName := "gruntwork-terratest-" + strings.ToLower(id)
	logger.Logf(t, "Random values selected Id = %s\n", id)

	// Don't create a new dataset so we can confirm that our function works as expected.

	err := AssertDatasetExistsE(t, projectID, datasetName)
	if err == nil {
		t.Fatalf("Function claimed that the Dataset '%s' exists, but in fact it does not.", datasetName)
	}
}
