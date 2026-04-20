package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// s3DeleteBatchSize is the maximum number of objects to delete in a single batch.
const s3DeleteBatchSize = 1000

// FindS3BucketWithTagContextE finds the name of the S3 bucket in the given region with the given tag key=value.
// The ctx parameter supports cancellation and timeouts.
func FindS3BucketWithTagContextE(t testing.TestingT, ctx context.Context, awsRegion string, key string, value string) (string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return "", err
	}

	for _, bucket := range resp.Buckets {
		tagResponse, err := s3Client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{Bucket: bucket.Name})
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchBucket") {
				// Occasionally, the ListBuckets call will return a bucket that has been deleted by S3
				// but hasn't yet been actually removed from the backend. Listing tags on that bucket
				// will return this error. If the bucket has been deleted, it can't be the one to find,
				// so just ignore this error, and keep checking the other buckets.
				continue
			}

			if !strings.Contains(err.Error(), "AuthorizationHeaderMalformed") &&
				!strings.Contains(err.Error(), "BucketRegionError") &&
				!strings.Contains(err.Error(), "NoSuchTagSet") {
				return "", err
			}

			continue
		}

		for _, tag := range tagResponse.TagSet {
			if *tag.Key == key && *tag.Value == value {
				logger.Default.Logf(t, "Found S3 bucket %s with tag %s=%s", *bucket.Name, key, value)

				return *bucket.Name, nil
			}
		}
	}

	return "", nil
}

// FindS3BucketWithTagContext finds the name of the S3 bucket in the given region with the given tag key=value.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FindS3BucketWithTagContext(t testing.TestingT, ctx context.Context, awsRegion string, key string, value string) string {
	t.Helper()

	bucket, err := FindS3BucketWithTagContextE(t, ctx, awsRegion, key, value)
	require.NoError(t, err)

	return bucket
}

// FindS3BucketWithTag finds the name of the S3 bucket in the given region with the given tag key=value.
//
// Deprecated: Use [FindS3BucketWithTagContext] instead.
func FindS3BucketWithTag(t testing.TestingT, awsRegion string, key string, value string) string {
	t.Helper()

	return FindS3BucketWithTagContext(t, context.Background(), awsRegion, key, value)
}

// FindS3BucketWithTagE finds the name of the S3 bucket in the given region with the given tag key=value.
//
// Deprecated: Use [FindS3BucketWithTagContextE] instead.
func FindS3BucketWithTagE(t testing.TestingT, awsRegion string, key string, value string) (string, error) {
	return FindS3BucketWithTagContextE(t, context.Background(), awsRegion, key, value)
}

// GetS3BucketTagsContextE fetches the given bucket's tags and returns them as a string map of strings.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketTagsContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) (map[string]string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	out, err := s3Client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	tags := map[string]string{}
	for _, tag := range out.TagSet {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return tags, nil
}

// GetS3BucketTagsContext fetches the given bucket's tags and returns them as a string map of strings.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketTagsContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) map[string]string {
	t.Helper()

	tags, err := GetS3BucketTagsContextE(t, ctx, awsRegion, bucket)
	require.NoError(t, err)

	return tags
}

// GetS3BucketTags fetches the given bucket's tags and returns them as a string map of strings.
//
// Deprecated: Use [GetS3BucketTagsContext] instead.
func GetS3BucketTags(t testing.TestingT, awsRegion string, bucket string) map[string]string {
	t.Helper()

	return GetS3BucketTagsContext(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketTagsE fetches the given bucket's tags and returns them as a string map of strings.
//
// Deprecated: Use [GetS3BucketTagsContextE] instead.
func GetS3BucketTagsE(t testing.TestingT, awsRegion string, bucket string) (map[string]string, error) {
	return GetS3BucketTagsContextE(t, context.Background(), awsRegion, bucket)
}

// GetS3ObjectContentsContextE fetches the contents of the object in the given bucket with the given key and return it as a string.
// The ctx parameter supports cancellation and timeouts.
func GetS3ObjectContentsContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string, key string) (string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}

	contents := buf.String()

	logger.Default.Logf(t, "Read contents from s3://%s/%s", bucket, key)

	return contents, nil
}

// GetS3ObjectContentsContext fetches the contents of the object in the given bucket with the given key and return it as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3ObjectContentsContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string, key string) string {
	t.Helper()

	contents, err := GetS3ObjectContentsContextE(t, ctx, awsRegion, bucket, key)
	require.NoError(t, err)

	return contents
}

// GetS3ObjectContents fetches the contents of the object in the given bucket with the given key and return it as a string.
//
// Deprecated: Use [GetS3ObjectContentsContext] instead.
func GetS3ObjectContents(t testing.TestingT, awsRegion string, bucket string, key string) string {
	t.Helper()

	return GetS3ObjectContentsContext(t, context.Background(), awsRegion, bucket, key)
}

// GetS3ObjectContentsE fetches the contents of the object in the given bucket with the given key and return it as a string.
//
// Deprecated: Use [GetS3ObjectContentsContextE] instead.
func GetS3ObjectContentsE(t testing.TestingT, awsRegion string, bucket string, key string) (string, error) {
	return GetS3ObjectContentsContextE(t, context.Background(), awsRegion, bucket, key)
}

// PutS3ObjectContentsContextE puts the contents of the object in the given bucket with the given key.
// The ctx parameter supports cancellation and timeouts.
func PutS3ObjectContentsContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string, key string, body io.Reader) error {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return fmt.Errorf("failed to instantiate s3 client: %w", err)
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	}

	_, err = s3Client.PutObject(ctx, params)

	return err
}

// PutS3ObjectContentsContext puts the contents of the object in the given bucket with the given key.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PutS3ObjectContentsContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string, key string, body io.Reader) {
	t.Helper()

	err := PutS3ObjectContentsContextE(t, ctx, awsRegion, bucket, key, body)
	require.NoError(t, err)
}

// PutS3ObjectContents puts the contents of the object in the given bucket with the given key.
//
// Deprecated: Use [PutS3ObjectContentsContext] instead.
func PutS3ObjectContents(t testing.TestingT, awsRegion string, bucket string, key string, body io.Reader) {
	t.Helper()

	PutS3ObjectContentsContext(t, context.Background(), awsRegion, bucket, key, body)
}

// PutS3ObjectContentsE puts the contents of the object in the given bucket with the given key.
//
// Deprecated: Use [PutS3ObjectContentsContextE] instead.
func PutS3ObjectContentsE(t testing.TestingT, awsRegion string, bucket string, key string, body io.Reader) error {
	return PutS3ObjectContentsContextE(t, context.Background(), awsRegion, bucket, key, body)
}

// CreateS3BucketContextE creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
// The ctx parameter supports cancellation and timeouts.
func CreateS3BucketContextE(t testing.TestingT, ctx context.Context, region string, name string) error {
	logger.Default.Logf(t, "Creating bucket %s in %s", name, region)

	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	params := &s3.CreateBucketInput{
		Bucket:          aws.String(name),
		ObjectOwnership: types.ObjectOwnershipObjectWriter,
	}

	if region != "us-east-1" {
		params.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		}
	}

	_, err = s3Client.CreateBucket(ctx, params)

	return err
}

// CreateS3BucketContext creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateS3BucketContext(t testing.TestingT, ctx context.Context, region string, name string) {
	t.Helper()

	err := CreateS3BucketContextE(t, ctx, region, name)
	require.NoError(t, err)
}

// CreateS3Bucket creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
//
// Deprecated: Use [CreateS3BucketContext] instead.
func CreateS3Bucket(t testing.TestingT, region string, name string) {
	t.Helper()

	CreateS3BucketContext(t, context.Background(), region, name)
}

// CreateS3BucketE creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
//
// Deprecated: Use [CreateS3BucketContextE] instead.
func CreateS3BucketE(t testing.TestingT, region string, name string) error {
	return CreateS3BucketContextE(t, context.Background(), region, name)
}

// PutS3BucketPolicyContextE applies an IAM resource policy to a given S3 bucket to create its bucket policy.
// The ctx parameter supports cancellation and timeouts.
func PutS3BucketPolicyContextE(t testing.TestingT, ctx context.Context, region string, bucketName string, policyJSONString string) error {
	logger.Default.Logf(t, "Applying bucket policy for bucket %s in %s", bucketName, region)

	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policyJSONString),
	}

	_, err = s3Client.PutBucketPolicy(ctx, input)

	return err
}

// PutS3BucketPolicyContext applies an IAM resource policy to a given S3 bucket to create its bucket policy.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PutS3BucketPolicyContext(t testing.TestingT, ctx context.Context, region string, bucketName string, policyJSONString string) {
	t.Helper()

	err := PutS3BucketPolicyContextE(t, ctx, region, bucketName, policyJSONString)
	require.NoError(t, err)
}

// PutS3BucketPolicy applies an IAM resource policy to a given S3 bucket to create its bucket policy
//
// Deprecated: Use [PutS3BucketPolicyContext] instead.
func PutS3BucketPolicy(t testing.TestingT, region string, bucketName string, policyJSONString string) {
	t.Helper()

	PutS3BucketPolicyContext(t, context.Background(), region, bucketName, policyJSONString)
}

// PutS3BucketPolicyE applies an IAM resource policy to a given S3 bucket to create its bucket policy
//
// Deprecated: Use [PutS3BucketPolicyContextE] instead.
func PutS3BucketPolicyE(t testing.TestingT, region string, bucketName string, policyJSONString string) error {
	return PutS3BucketPolicyContextE(t, context.Background(), region, bucketName, policyJSONString)
}

// PutS3BucketVersioningContextE creates an S3 bucket versioning configuration in the given region against the given bucket name, WITHOUT requiring MFA to remove versioning.
// The ctx parameter supports cancellation and timeouts.
func PutS3BucketVersioningContextE(t testing.TestingT, ctx context.Context, region string, bucketName string) error {
	logger.Default.Logf(t, "Creating bucket versioning configuration for bucket %s in %s", bucketName, region)

	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	input := &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			MFADelete: types.MFADeleteDisabled,
			Status:    types.BucketVersioningStatusEnabled,
		},
	}

	_, err = s3Client.PutBucketVersioning(ctx, input)

	return err
}

// PutS3BucketVersioningContext creates an S3 bucket versioning configuration in the given region against the given bucket name, WITHOUT requiring MFA to remove versioning.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PutS3BucketVersioningContext(t testing.TestingT, ctx context.Context, region string, bucketName string) {
	t.Helper()

	err := PutS3BucketVersioningContextE(t, ctx, region, bucketName)
	require.NoError(t, err)
}

// PutS3BucketVersioning creates an S3 bucket versioning configuration in the given region against the given bucket name, WITHOUT requiring MFA to remove versioning.
//
// Deprecated: Use [PutS3BucketVersioningContext] instead.
func PutS3BucketVersioning(t testing.TestingT, region string, bucketName string) {
	t.Helper()

	PutS3BucketVersioningContext(t, context.Background(), region, bucketName)
}

// PutS3BucketVersioningE creates an S3 bucket versioning configuration in the given region against the given bucket name, WITHOUT requiring MFA to remove versioning.
//
// Deprecated: Use [PutS3BucketVersioningContextE] instead.
func PutS3BucketVersioningE(t testing.TestingT, region string, bucketName string) error {
	return PutS3BucketVersioningContextE(t, context.Background(), region, bucketName)
}

// DeleteS3BucketContextE destroys the S3 bucket in the given region with the given name.
// The ctx parameter supports cancellation and timeouts.
func DeleteS3BucketContextE(t testing.TestingT, ctx context.Context, region string, name string) error {
	logger.Default.Logf(t, "Deleting bucket %s in %s", region, name)

	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = s3Client.DeleteBucket(ctx, params)

	return err
}

// DeleteS3BucketContext destroys the S3 bucket in the given region with the given name.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteS3BucketContext(t testing.TestingT, ctx context.Context, region string, name string) {
	t.Helper()

	err := DeleteS3BucketContextE(t, ctx, region, name)
	require.NoError(t, err)
}

// DeleteS3Bucket destroys the S3 bucket in the given region with the given name.
//
// Deprecated: Use [DeleteS3BucketContext] instead.
func DeleteS3Bucket(t testing.TestingT, region string, name string) {
	t.Helper()

	DeleteS3BucketContext(t, context.Background(), region, name)
}

// DeleteS3BucketE destroys the S3 bucket in the given region with the given name.
//
// Deprecated: Use [DeleteS3BucketContextE] instead.
func DeleteS3BucketE(t testing.TestingT, region string, name string) error {
	return DeleteS3BucketContextE(t, context.Background(), region, name)
}

// EmptyS3BucketContextE removes the contents of an S3 bucket in the given region with the given name.
// The ctx parameter supports cancellation and timeouts.
func EmptyS3BucketContextE(t testing.TestingT, ctx context.Context, region string, name string) error {
	logger.Default.Logf(t, "Emptying bucket %s in %s", name, region)

	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	params := &s3.ListObjectVersionsInput{
		Bucket: aws.String(name),
	}

	for {
		// Requesting a batch of objects from s3 bucket
		bucketObjects, err := s3Client.ListObjectVersions(ctx, params)
		if err != nil {
			return err
		}

		// Checks if the bucket is already empty
		if len(bucketObjects.Versions) == 0 {
			logger.Default.Logf(t, "Bucket %s is already empty", name)

			return nil
		}

		// creating an array of pointers of ObjectIdentifier
		objectsToDelete := make([]types.ObjectIdentifier, 0, s3DeleteBatchSize)

		for i := range bucketObjects.Versions {
			object := &bucketObjects.Versions[i]
			obj := types.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			}
			objectsToDelete = append(objectsToDelete, obj)
		}

		for i := range bucketObjects.DeleteMarkers {
			object := &bucketObjects.DeleteMarkers[i]
			obj := types.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			}
			objectsToDelete = append(objectsToDelete, obj)
		}

		// Creating JSON payload for bulk delete
		deleteArray := types.Delete{Objects: objectsToDelete}
		deleteParams := &s3.DeleteObjectsInput{
			Bucket: aws.String(name),
			Delete: &deleteArray,
		}

		// Running the Bulk delete job (limit 1000)
		_, err = s3Client.DeleteObjects(ctx, deleteParams)
		if err != nil {
			return err
		}

		if *bucketObjects.IsTruncated { // if there are more objects in the bucket, IsTruncated = true
			// params.Marker = (*deleteParams).Delete.Objects[len((*deleteParams).Delete.Objects)-1].Key
			params.KeyMarker = bucketObjects.NextKeyMarker
			logger.Default.Logf(t, "Requesting next batch | %s", *(params.KeyMarker))
		} else { // if all objects in the bucket have been cleaned up.
			break
		}
	}

	logger.Default.Logf(t, "Bucket %s is now empty", name)

	return err
}

// EmptyS3BucketContext removes the contents of an S3 bucket in the given region with the given name.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func EmptyS3BucketContext(t testing.TestingT, ctx context.Context, region string, name string) {
	t.Helper()

	err := EmptyS3BucketContextE(t, ctx, region, name)
	require.NoError(t, err)
}

// EmptyS3Bucket removes the contents of an S3 bucket in the given region with the given name.
//
// Deprecated: Use [EmptyS3BucketContext] instead.
func EmptyS3Bucket(t testing.TestingT, region string, name string) {
	t.Helper()

	EmptyS3BucketContext(t, context.Background(), region, name)
}

// EmptyS3BucketE removes the contents of an S3 bucket in the given region with the given name.
//
// Deprecated: Use [EmptyS3BucketContextE] instead.
func EmptyS3BucketE(t testing.TestingT, region string, name string) error {
	return EmptyS3BucketContextE(t, context.Background(), region, name)
}

// GetS3BucketLoggingTargetContextE fetches the given bucket's logging target bucket and returns it as the following string:
// `TargetBucket` of the `LoggingEnabled` property for an S3 bucket.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketLoggingTargetContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketLogging(ctx, &s3.GetBucketLoggingInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}

	if res.LoggingEnabled == nil {
		return "", S3AccessLoggingNotEnabledErr{bucket, awsRegion}
	}

	return aws.ToString(res.LoggingEnabled.TargetBucket), nil
}

// GetS3BucketLoggingTargetContext fetches the given bucket's logging target bucket and returns it as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketLoggingTargetContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) string {
	t.Helper()

	loggingTarget, err := GetS3BucketLoggingTargetContextE(t, ctx, awsRegion, bucket)
	require.NoError(t, err)

	return loggingTarget
}

// GetS3BucketLoggingTarget fetches the given bucket's logging target bucket and returns it as a string
//
// Deprecated: Use [GetS3BucketLoggingTargetContext] instead.
func GetS3BucketLoggingTarget(t testing.TestingT, awsRegion string, bucket string) string {
	t.Helper()

	return GetS3BucketLoggingTargetContext(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketLoggingTargetE fetches the given bucket's logging target bucket and returns it as the following string:
// `TargetBucket` of the `LoggingEnabled` property for an S3 bucket
//
// Deprecated: Use [GetS3BucketLoggingTargetContextE] instead.
func GetS3BucketLoggingTargetE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	return GetS3BucketLoggingTargetContextE(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketLoggingTargetPrefixContextE fetches the given bucket's logging object prefix and returns it as the following string:
// `TargetPrefix` of the `LoggingEnabled` property for an S3 bucket.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketLoggingTargetPrefixContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketLogging(ctx, &s3.GetBucketLoggingInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}

	if res.LoggingEnabled == nil {
		return "", S3AccessLoggingNotEnabledErr{bucket, awsRegion}
	}

	return aws.ToString(res.LoggingEnabled.TargetPrefix), nil
}

// GetS3BucketLoggingTargetPrefixContext fetches the given bucket's logging object prefix and returns it as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketLoggingTargetPrefixContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) string {
	t.Helper()

	loggingObjectTargetPrefix, err := GetS3BucketLoggingTargetPrefixContextE(t, ctx, awsRegion, bucket)
	require.NoError(t, err)

	return loggingObjectTargetPrefix
}

// GetS3BucketLoggingTargetPrefix fetches the given bucket's logging object prefix and returns it as a string
//
// Deprecated: Use [GetS3BucketLoggingTargetPrefixContext] instead.
func GetS3BucketLoggingTargetPrefix(t testing.TestingT, awsRegion string, bucket string) string {
	t.Helper()

	return GetS3BucketLoggingTargetPrefixContext(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketLoggingTargetPrefixE fetches the given bucket's logging object prefix and returns it as the following string:
// `TargetPrefix` of the `LoggingEnabled` property for an S3 bucket
//
// Deprecated: Use [GetS3BucketLoggingTargetPrefixContextE] instead.
func GetS3BucketLoggingTargetPrefixE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	return GetS3BucketLoggingTargetPrefixContextE(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketVersioningContextE fetches the given bucket's versioning configuration status and returns it as a string.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketVersioningContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}

	return string(res.Status), nil
}

// GetS3BucketVersioningContext fetches the given bucket's versioning configuration status and returns it as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketVersioningContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) string {
	t.Helper()

	versioningStatus, err := GetS3BucketVersioningContextE(t, ctx, awsRegion, bucket)
	require.NoError(t, err)

	return versioningStatus
}

// GetS3BucketVersioning fetches the given bucket's versioning configuration status and returns it as a string
//
// Deprecated: Use [GetS3BucketVersioningContext] instead.
func GetS3BucketVersioning(t testing.TestingT, awsRegion string, bucket string) string {
	t.Helper()

	return GetS3BucketVersioningContext(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketVersioningE fetches the given bucket's versioning configuration status and returns it as a string
//
// Deprecated: Use [GetS3BucketVersioningContextE] instead.
func GetS3BucketVersioningE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	return GetS3BucketVersioningContextE(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketPolicyContextE fetches the given bucket's resource policy and returns it as a string.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketPolicyContextE(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(res.Policy), nil
}

// GetS3BucketPolicyContext fetches the given bucket's resource policy and returns it as a string.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketPolicyContext(t testing.TestingT, ctx context.Context, awsRegion string, bucket string) string {
	t.Helper()

	bucketPolicy, err := GetS3BucketPolicyContextE(t, ctx, awsRegion, bucket)
	require.NoError(t, err)

	return bucketPolicy
}

// GetS3BucketPolicy fetches the given bucket's resource policy and returns it as a string
//
// Deprecated: Use [GetS3BucketPolicyContext] instead.
func GetS3BucketPolicy(t testing.TestingT, awsRegion string, bucket string) string {
	t.Helper()

	return GetS3BucketPolicyContext(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketPolicyE fetches the given bucket's resource policy and returns it as a string
//
// Deprecated: Use [GetS3BucketPolicyContextE] instead.
func GetS3BucketPolicyE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	return GetS3BucketPolicyContextE(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketOwnershipControlsContextE fetches the given bucket's ownership controls and returns them as a slice of strings.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketOwnershipControlsContextE(t testing.TestingT, ctx context.Context, awsRegion, bucket string) ([]string, error) {
	s3Client, err := NewS3ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	out, err := s3Client.GetBucketOwnershipControls(ctx, &s3.GetBucketOwnershipControlsInput{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	rules := make([]string, 0, len(out.OwnershipControls.Rules))
	for _, rule := range out.OwnershipControls.Rules {
		rules = append(rules, string(rule.ObjectOwnership))
	}

	return rules, nil
}

// GetS3BucketOwnershipControlsContext fetches the given bucket's ownership controls and returns them as a slice of strings.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetS3BucketOwnershipControlsContext(t testing.TestingT, ctx context.Context, awsRegion, bucket string) []string {
	t.Helper()

	rules, err := GetS3BucketOwnershipControlsContextE(t, ctx, awsRegion, bucket)
	require.NoError(t, err)

	return rules
}

// GetS3BucketOwnershipControls fetches the given bucket's ownership controls and returns them as a slice of strings.
//
// Deprecated: Use [GetS3BucketOwnershipControlsContext] instead.
func GetS3BucketOwnershipControls(t testing.TestingT, awsRegion, bucket string) []string {
	t.Helper()

	return GetS3BucketOwnershipControlsContext(t, context.Background(), awsRegion, bucket)
}

// GetS3BucketOwnershipControlsE fetches the given bucket's ownership controls and returns them as a slice of strings.
//
// Deprecated: Use [GetS3BucketOwnershipControlsContextE] instead.
func GetS3BucketOwnershipControlsE(t testing.TestingT, awsRegion, bucket string) ([]string, error) {
	return GetS3BucketOwnershipControlsContextE(t, context.Background(), awsRegion, bucket)
}

// AssertS3BucketExistsContextE checks if the given S3 bucket exists in the given region and return an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertS3BucketExistsContextE(t testing.TestingT, ctx context.Context, region string, name string) error {
	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	params := &s3.HeadBucketInput{
		Bucket: aws.String(name),
	}

	_, err = s3Client.HeadBucket(ctx, params)

	return err
}

// AssertS3BucketExistsContext checks if the given S3 bucket exists in the given region and fail the test if it does not.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func AssertS3BucketExistsContext(t testing.TestingT, ctx context.Context, region string, name string) {
	t.Helper()

	err := AssertS3BucketExistsContextE(t, ctx, region, name)
	require.NoError(t, err)
}

// AssertS3BucketExists checks if the given S3 bucket exists in the given region and fail the test if it does not.
//
// Deprecated: Use [AssertS3BucketExistsContext] instead.
func AssertS3BucketExists(t testing.TestingT, region string, name string) {
	t.Helper()

	AssertS3BucketExistsContext(t, context.Background(), region, name)
}

// AssertS3BucketExistsE checks if the given S3 bucket exists in the given region and return an error if it does not.
//
// Deprecated: Use [AssertS3BucketExistsContextE] instead.
func AssertS3BucketExistsE(t testing.TestingT, region string, name string) error {
	return AssertS3BucketExistsContextE(t, context.Background(), region, name)
}

// AssertS3BucketVersioningExistsContextE checks if the given S3 bucket has a versioning configuration enabled and returns an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertS3BucketVersioningExistsContextE(t testing.TestingT, ctx context.Context, region string, bucketName string) error {
	status, err := GetS3BucketVersioningContextE(t, ctx, region, bucketName)
	if err != nil {
		return err
	}

	if status == "Enabled" {
		return nil
	}

	return NewBucketVersioningNotEnabledError(bucketName, region, status)
}

// AssertS3BucketVersioningExistsContext checks if the given S3 bucket has a versioning configuration enabled and fails the test if it does not.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func AssertS3BucketVersioningExistsContext(t testing.TestingT, ctx context.Context, region string, bucketName string) {
	t.Helper()

	err := AssertS3BucketVersioningExistsContextE(t, ctx, region, bucketName)
	require.NoError(t, err)
}

// AssertS3BucketVersioningExists checks if the given S3 bucket has a versioning configuration enabled and returns an error if it does not.
//
// Deprecated: Use [AssertS3BucketVersioningExistsContext] instead.
func AssertS3BucketVersioningExists(t testing.TestingT, region string, bucketName string) {
	t.Helper()

	AssertS3BucketVersioningExistsContext(t, context.Background(), region, bucketName)
}

// AssertS3BucketVersioningExistsE checks if the given S3 bucket has a versioning configuration enabled and returns an error if it does not.
//
// Deprecated: Use [AssertS3BucketVersioningExistsContextE] instead.
func AssertS3BucketVersioningExistsE(t testing.TestingT, region string, bucketName string) error {
	return AssertS3BucketVersioningExistsContextE(t, context.Background(), region, bucketName)
}

// AssertS3BucketPolicyExistsContextE checks if the given S3 bucket has a resource policy attached and returns an error if it does not.
// The ctx parameter supports cancellation and timeouts.
func AssertS3BucketPolicyExistsContextE(t testing.TestingT, ctx context.Context, region string, bucketName string) error {
	policy, err := GetS3BucketPolicyContextE(t, ctx, region, bucketName)
	if err != nil {
		return err
	}

	if policy == "" {
		return NewNoBucketPolicyError(bucketName, region, policy)
	}

	return nil
}

// AssertS3BucketPolicyExistsContext checks if the given S3 bucket has a resource policy attached and fails the test if it does not.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func AssertS3BucketPolicyExistsContext(t testing.TestingT, ctx context.Context, region string, bucketName string) {
	t.Helper()

	err := AssertS3BucketPolicyExistsContextE(t, ctx, region, bucketName)
	require.NoError(t, err)
}

// AssertS3BucketPolicyExists checks if the given S3 bucket has a resource policy attached and returns an error if it does not
//
// Deprecated: Use [AssertS3BucketPolicyExistsContext] instead.
func AssertS3BucketPolicyExists(t testing.TestingT, region string, bucketName string) {
	t.Helper()

	AssertS3BucketPolicyExistsContext(t, context.Background(), region, bucketName)
}

// AssertS3BucketPolicyExistsE checks if the given S3 bucket has a resource policy attached and returns an error if it does not
//
// Deprecated: Use [AssertS3BucketPolicyExistsContextE] instead.
func AssertS3BucketPolicyExistsE(t testing.TestingT, region string, bucketName string) error {
	return AssertS3BucketPolicyExistsContextE(t, context.Background(), region, bucketName)
}

// NewS3ClientContextE creates an S3 client.
// The ctx parameter supports cancellation and timeouts.
func NewS3ClientContextE(t testing.TestingT, ctx context.Context, region string) (*s3.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(*sess), nil
}

// NewS3ClientContext creates an S3 client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewS3ClientContext(t testing.TestingT, ctx context.Context, region string) *s3.Client {
	t.Helper()

	client, err := NewS3ClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewS3Client creates an S3 client.
//
// Deprecated: Use [NewS3ClientContext] instead.
func NewS3Client(t testing.TestingT, region string) *s3.Client {
	t.Helper()

	return NewS3ClientContext(t, context.Background(), region)
}

// NewS3ClientE creates an S3 client.
//
// Deprecated: Use [NewS3ClientContextE] instead.
func NewS3ClientE(t testing.TestingT, region string) (*s3.Client, error) {
	return NewS3ClientContextE(t, context.Background(), region)
}

// NewS3UploaderContextE creates an S3 transfer manager client for uploading objects.
// The ctx parameter supports cancellation and timeouts.
func NewS3UploaderContextE(t testing.TestingT, ctx context.Context, region string) (*transfermanager.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return transfermanager.New(s3.NewFromConfig(*sess)), nil
}

// NewS3UploaderContext creates an S3 transfer manager client for uploading objects.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewS3UploaderContext(t testing.TestingT, ctx context.Context, region string) *transfermanager.Client {
	t.Helper()

	uploader, err := NewS3UploaderContextE(t, ctx, region)
	require.NoError(t, err)

	return uploader
}

// NewS3Uploader creates an S3 transfer manager client for uploading objects.
//
// Deprecated: Use [NewS3UploaderContext] instead.
func NewS3Uploader(t testing.TestingT, region string) *transfermanager.Client {
	t.Helper()

	return NewS3UploaderContext(t, context.Background(), region)
}

// NewS3UploaderE creates an S3 transfer manager client for uploading objects.
//
// Deprecated: Use [NewS3UploaderContextE] instead.
func NewS3UploaderE(t testing.TestingT, region string) (*transfermanager.Client, error) {
	return NewS3UploaderContextE(t, context.Background(), region)
}

// S3AccessLoggingNotEnabledErr is a custom error that occurs when acess logging hasn't been enabled on the S3 Bucket
type S3AccessLoggingNotEnabledErr struct {
	OriginBucket string
	Region       string
}

func (err S3AccessLoggingNotEnabledErr) Error() string {
	return fmt.Sprintf("Server Access Logging hasn't been enabled for S3 Bucket %s in region %s", err.OriginBucket, err.Region)
}
