// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestCreateAndDestroyStorageBucket(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	id := random.UniqueId()
	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)
	testFilePath := fmt.Sprintf("test-file-%s.txt", random.UniqueId())
	testFileBody := "test file text"

	logger.Logf(t, "Random values selected Bucket Name = %s, Test Filepath: %s\n", gsBucketName, testFilePath)

	CreateStorageBucket(t, projectID, gsBucketName, nil)
	defer DeleteStorageBucket(t, gsBucketName)

	// Write a test file to the storage bucket
	objectURL := WriteBucketObject(t, gsBucketName, testFilePath, strings.NewReader(testFileBody), "text/plain")
	logger.Logf(t, "Got URL: %s", objectURL)

	// Then verify its contents matches the expected result
	fileReader := ReadBucketObject(t, gsBucketName, testFilePath)

	buf := new(bytes.Buffer)
	buf.ReadFrom(fileReader)
	result := buf.String()

	require.Equal(t, testFileBody, result)

	// Empty the storage bucket so we can delete it
	defer EmptyStorageBucket(t, gsBucketName)
}

func TestAssertStorageBucketExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	id := random.UniqueId()
	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)
	logger.Logf(t, "Random values selected Id = %s\n", id)

	CreateStorageBucket(t, projectID, gsBucketName, nil)
	defer DeleteStorageBucket(t, gsBucketName)

	AssertStorageBucketExists(t, gsBucketName)
}

func TestAssertStorageBucketExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	id := random.UniqueId()
	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)
	logger.Logf(t, "Random values selected Id = %s\n", id)

	// Don't create a new storage bucket so we can confirm that our function works as expected.

	err := AssertStorageBucketExistsE(t, gsBucketName)
	if err == nil {
		t.Fatalf("Function claimed that the Storage Bucket '%s' exists, but in fact it does not.", gsBucketName)
	}
}
