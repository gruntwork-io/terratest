package gcp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iterator"
)

// CreateStorageBucketContext creates a Google Cloud bucket with the given BucketAttrs.
// Note that Google Storage bucket names must be globally unique.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateStorageBucketContext(t testing.TestingT, ctx context.Context, projectID string, name string, attr *storage.BucketAttrs) {
	err := CreateStorageBucketContextE(t, ctx, projectID, name, attr)
	require.NoError(t, err)
}

// CreateStorageBucketContextE creates a Google Cloud bucket with the given BucketAttrs.
// Note that Google Storage bucket names must be globally unique.
// The ctx parameter supports cancellation and timeouts.
func CreateStorageBucketContextE(t testing.TestingT, ctx context.Context, projectID string, name string, attr *storage.BucketAttrs) error {
	logger.Default.Logf(t, "Creating bucket %s", name)

	client, err := newStorageClient(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return CreateStorageBucketWithClient(ctx, client, projectID, name, attr)
}

// CreateStorageBucketWithClient creates a Google Cloud bucket with the given BucketAttrs using the
// supplied *storage.Client. Prefer this variant in unit tests where the client is backed by an
// httptest fake server (see storage_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func CreateStorageBucketWithClient(ctx context.Context, client *storage.Client, projectID string, name string, attr *storage.BucketAttrs) error {
	return client.Bucket(name).Create(ctx, projectID, attr)
}

// DeleteStorageBucketContext destroys the Google Storage bucket.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteStorageBucketContext(t testing.TestingT, ctx context.Context, name string) {
	err := DeleteStorageBucketContextE(t, ctx, name)
	require.NoError(t, err)
}

// DeleteStorageBucketContextE destroys the Google Cloud Storage bucket with the given name.
// The ctx parameter supports cancellation and timeouts.
func DeleteStorageBucketContextE(t testing.TestingT, ctx context.Context, name string) error {
	logger.Default.Logf(t, "Deleting bucket %s", name)

	client, err := newStorageClient(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return DeleteStorageBucketWithClient(ctx, client, name)
}

// DeleteStorageBucketWithClient destroys the Google Cloud Storage bucket with the given name using
// the supplied *storage.Client. Prefer this variant in unit tests where the client is backed by an
// httptest fake server (see storage_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DeleteStorageBucketWithClient(ctx context.Context, client *storage.Client, name string) error {
	return client.Bucket(name).Delete(ctx)
}

// ReadBucketObjectContext reads an object from the given Storage Bucket and returns its contents.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ReadBucketObjectContext(t testing.TestingT, ctx context.Context, bucketName string, filePath string) io.Reader {
	out, err := ReadBucketObjectContextE(t, ctx, bucketName, filePath)
	require.NoError(t, err)

	return out
}

// ReadBucketObjectContextE reads an object from the given Storage Bucket and returns its contents.
// The ctx parameter supports cancellation and timeouts.
func ReadBucketObjectContextE(t testing.TestingT, ctx context.Context, bucketName string, filePath string) (io.Reader, error) {
	logger.Default.Logf(t, "Reading object from bucket %s using path %s", bucketName, filePath)

	client, err := newStorageClient(ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = client.Close() }()

	return ReadBucketObjectWithClient(ctx, client, bucketName, filePath)
}

// ReadBucketObjectWithClient reads an object from the given Storage Bucket and returns its contents
// using the supplied *storage.Client. Prefer this variant in unit tests where the client is backed
// by an httptest fake server (see storage_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func ReadBucketObjectWithClient(ctx context.Context, client *storage.Client, bucketName string, filePath string) (io.Reader, error) {
	r, err := client.Bucket(bucketName).Object(filePath).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer func() { _ = r.Close() }()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	return &buf, nil
}

// WriteBucketObjectContext writes an object to the given Storage Bucket and returns its URL.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func WriteBucketObjectContext(t testing.TestingT, ctx context.Context, bucketName string, filePath string, body io.Reader, contentType string) string {
	out, err := WriteBucketObjectContextE(t, ctx, bucketName, filePath, body, contentType)
	require.NoError(t, err)

	return out
}

// WriteBucketObjectContextE writes an object to the given Storage Bucket and returns its URL.
// The ctx parameter supports cancellation and timeouts.
func WriteBucketObjectContextE(t testing.TestingT, ctx context.Context, bucketName string, filePath string, body io.Reader, contentType string) (string, error) {

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	logger.Default.Logf(t, "Writing object to bucket %s using path %s and content type %s", bucketName, filePath, contentType)

	client, err := newStorageClient(ctx)
	if err != nil {
		return "", err
	}

	defer func() { _ = client.Close() }()

	return WriteBucketObjectWithClient(ctx, client, bucketName, filePath, body, contentType)
}

// WriteBucketObjectWithClient writes an object to the given Storage Bucket and returns its URL
// using the supplied *storage.Client. Prefer this variant in unit tests where the client is backed
// by an httptest fake server (see storage_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func WriteBucketObjectWithClient(ctx context.Context, client *storage.Client, bucketName string, filePath string, body io.Reader, contentType string) (string, error) {

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w := client.Bucket(bucketName).Object(filePath).NewWriter(ctx)
	w.ContentType = contentType

	if _, err := io.Copy(w, body); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"

	return fmt.Sprintf(publicURL, bucketName, filePath), nil
}

// EmptyStorageBucketContext removes the contents of a storage bucket with the given name.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func EmptyStorageBucketContext(t testing.TestingT, ctx context.Context, name string) {
	err := EmptyStorageBucketContextE(t, ctx, name)
	require.NoError(t, err)
}

// EmptyStorageBucketContextE removes the contents of a storage bucket with the given name.
// The ctx parameter supports cancellation and timeouts.
func EmptyStorageBucketContextE(t testing.TestingT, ctx context.Context, name string) error {
	logger.Default.Logf(t, "Emptying storage bucket %s", name)

	client, err := newStorageClient(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	bucket := client.Bucket(name)

	it := bucket.Objects(ctx, nil)

	for {
		objectAttrs, err := it.Next()

		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return err
		}

		logger.Default.Logf(t, "Deleting storage bucket object %s", objectAttrs.Name)

		if err := bucket.Object(objectAttrs.Name).Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}

// EmptyStorageBucketWithClient removes the contents of a storage bucket with the given name using
// the supplied *storage.Client. Prefer this variant in unit tests where the client is backed by an
// httptest fake server (see storage_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func EmptyStorageBucketWithClient(ctx context.Context, client *storage.Client, name string) error {
	bucket := client.Bucket(name)

	it := bucket.Objects(ctx, nil)

	for {
		objectAttrs, err := it.Next()

		if errors.Is(err, iterator.Done) {
			break
		}

		if err != nil {
			return err
		}

		if err := bucket.Object(objectAttrs.Name).Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}

// AssertStorageBucketExistsContext checks if the given storage bucket exists and fails the test if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertStorageBucketExistsContext(t testing.TestingT, ctx context.Context, name string) {
	err := AssertStorageBucketExistsContextE(t, ctx, name)
	require.NoError(t, err)
}

// AssertStorageBucketExistsContextE checks if the given storage bucket exists and returns an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertStorageBucketExistsContextE(t testing.TestingT, ctx context.Context, name string) error {
	logger.Default.Logf(t, "Finding bucket %s", name)

	client, err := newStorageClient(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = client.Close() }()

	return AssertStorageBucketExistsWithClient(ctx, client, name)
}

// AssertStorageBucketExistsWithClient checks if the given storage bucket exists and returns an error
// if it does not, using the supplied *storage.Client. Prefer this variant in unit tests where the
// client is backed by an httptest fake server (see storage_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func AssertStorageBucketExistsWithClient(ctx context.Context, client *storage.Client, name string) error {
	bucket := client.Bucket(name)

	if _, err := bucket.Attrs(ctx); err != nil {

		return err
	}

	it := bucket.Objects(ctx, nil)

	if _, err := it.Next(); errors.Is(err, storage.ErrBucketNotExist) {
		return err
	}

	return nil
}

func newStorageClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, withOptions()...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
