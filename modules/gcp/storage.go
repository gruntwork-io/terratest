package gcp

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/gruntwork-io/terratest/modules/logger"
)

// CreateStorageBucket creates a Google Cloud bucket with the given BucketAttrs. Note that Google Storage bucket names must be globally unique.
func CreateStorageBucket(t *testing.T, projectId string, name string, attr *storage.BucketAttrs) {
	err := CreateStorageBucketE(t, projectId, name, attr)
	if err != nil {
		t.Fatal(err)
	}
}

// CreateStorageBucketE creates a Google Cloud bucket with the given BucketAttrs. Note that Google Storage bucket names must be globally unique.
func CreateStorageBucketE(t *testing.T, projectId string, name string, attr *storage.BucketAttrs) error {
	logger.Logf(t, "Creating bucket %s", name)

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := name

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates the new bucket.
	if err := bucket.Create(ctx, projectId, attr); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	return err
}

// DeleteStorageBucket destroys the Google Storage bucket.
func DeleteStorageBucket(t *testing.T, name string) {

	err := DeleteStorageBucketE(t, name)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteS3BucketE destroys the S3 bucket in the given region with the given name.
func DeleteStorageBucketE(t *testing.T, name string) error {
	logger.Logf(t, "Deleting bucket %s", name)

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	if err := client.Bucket(name).Delete(ctx); err != nil {
		return err
	}

	return err
}
