package gcp

import (
	"context"

	"cloud.google.com/go/bigquery"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// CreateDataset creates a BigQuery Dataset with the given DatasetMetadata.
func CreateDataset(t testing.TestingT, projectID, datasetID string, dm bigquery.DatasetMetadata) {
	err := CreateDatasetE(t, projectID, datasetID, dm)
	if err != nil {
		t.Fatal(err)
	}
}

// CreateDatasetE creates a BigQuery Dataset with the given DatasetMetadata.
func CreateDatasetE(t testing.TestingT, projectID, datasetID string, dm bigquery.DatasetMetadata) error {
	logger.Logf(t, "Creating dataset %s", datasetID)

	ctx := context.Background()

	// Creates a client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	// Creates a Dataset handle
	dataset := client.Dataset(datasetID)

	// Creates the new dataset.
	return dataset.Create(ctx, &dm)
}

// DeleteDataset destroys the dataset with the given name.
func DeleteDataset(t testing.TestingT, projectID, datasetID string) {
	err := DeleteDatasetE(t, projectID, datasetID)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteDatasetE destroys the dataset with the given name.
func DeleteDatasetE(t testing.TestingT, projectID, datasetID string) error {
	logger.Logf(t, "Deleting dataset %s", datasetID)

	ctx := context.Background()

	// Creates a client
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	// Deletes the dataset
	return client.Dataset(datasetID).Delete(ctx)
}

// AssertDatasetExists checks if the given dataset exists and fails the test if it does not.
func AssertDatasetExists(t testing.TestingT, projectID, datasetID string) {
	err := AssertDatasetExistsE(t, projectID, datasetID)

	if err != nil {
		t.Fatal(err)
	}
}

// AssertDatasetExistsE checks if the given dataset exists and returns an error if it does not.
func AssertDatasetExistsE(t testing.TestingT, projectID, datasetID string) error {

	logger.Logf(t, "Finding dataset %s", datasetID)

	ctx := context.Background()

	// Creates a client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}

	// Seeks metadata for Dataset
	if _, err := client.Dataset(datasetID).Metadata(ctx); err != nil {
		return err
	}
	return nil

}
