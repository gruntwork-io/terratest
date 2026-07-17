package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithy "github.com/aws/smithy-go"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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

		bucketObjects, err := s3Client.ListObjectVersions(ctx, params)
		if err != nil {
			return err
		}

		if len(bucketObjects.Versions) == 0 {
			logger.Default.Logf(t, "Bucket %s is already empty", name)

			return nil
		}

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

		deleteArray := types.Delete{Objects: objectsToDelete}
		deleteParams := &s3.DeleteObjectsInput{
			Bucket: aws.String(name),
			Delete: &deleteArray,
		}

		_, err = s3Client.DeleteObjects(ctx, deleteParams)
		if err != nil {
			return err
		}

		if *bucketObjects.IsTruncated {

			params.KeyMarker = bucketObjects.NextKeyMarker
			logger.Default.Logf(t, "Requesting next batch | %s", *(params.KeyMarker))
		} else {
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

// AssertS3BucketServerSideEncryptionContextE checks if the given S3 bucket has server-side encryption configured with
// the given algorithm, and returns an error if it does not. The ctx parameter supports cancellation and timeouts.
//
// The algorithm is matched exactly: an expectation of `aws:kms` will not match a bucket configured with
// `aws:kms:dsse`, and vice versa.
func AssertS3BucketServerSideEncryptionContextE(t testing.TestingT, ctx context.Context, region string, bucketName string, algorithm types.ServerSideEncryption) error {
	s3Client, err := NewS3ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	out, err := s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		// A bucket with no SSE configuration surfaces as ServerSideEncryptionConfigurationNotFoundError. Translate
		// that to our typed error so callers can match on the failure mode regardless of SDK version.
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ServerSideEncryptionConfigurationNotFoundError" {
			return NewBucketServerSideEncryptionNotEnabledError(bucketName, region, algorithm)
		}

		return err
	}

	if out.ServerSideEncryptionConfiguration == nil {
		return NewBucketServerSideEncryptionNotEnabledError(bucketName, region, algorithm)
	}

	for _, rule := range out.ServerSideEncryptionConfiguration.Rules {
		if rule.ApplyServerSideEncryptionByDefault == nil {
			continue
		}

		if rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm == algorithm {
			return nil
		}
	}

	return NewBucketServerSideEncryptionNotEnabledError(bucketName, region, algorithm)
}

// AssertS3BucketServerSideEncryptionContext checks if the given S3 bucket has server-side encryption configured with
// the given algorithm, and fails the test if it does not. The ctx parameter supports cancellation and timeouts.
func AssertS3BucketServerSideEncryptionContext(t testing.TestingT, ctx context.Context, region string, bucketName string, algorithm types.ServerSideEncryption) {
	t.Helper()

	err := AssertS3BucketServerSideEncryptionContextE(t, ctx, region, bucketName, algorithm)
	require.NoError(t, err)
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

// S3AccessLoggingNotEnabledErr is a custom error that occurs when acess logging hasn't been enabled on the S3 Bucket
type S3AccessLoggingNotEnabledErr struct {
	OriginBucket string
	Region       string
}

func (err S3AccessLoggingNotEnabledErr) Error() string {
	return fmt.Sprintf("Server Access Logging hasn't been enabled for S3 Bucket %s in region %s", err.OriginBucket, err.Region)
}
