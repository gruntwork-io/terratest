package aws

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// FindS3BucketWithTag finds the name of the S3 bucket in the given region with the given tag key=value.
func FindS3BucketWithTag(t testing.TestingT, awsRegion string, key string, value string) string {
	bucket, err := FindS3BucketWithTagE(t, awsRegion, key, value)
	require.NoError(t, err)

	return bucket
}

// FindS3BucketWithTagE finds the name of the S3 bucket in the given region with the given tag key=value.
func FindS3BucketWithTagE(t testing.TestingT, awsRegion string, key string, value string) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	resp, err := s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return "", err
	}

	for _, bucket := range resp.Buckets {
		tagResponse, err := s3Client.GetBucketTagging(context.Background(), &s3.GetBucketTaggingInput{Bucket: bucket.Name})

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

// GetS3BucketTags fetches the given bucket's tags and returns them as a string map of strings.
func GetS3BucketTags(t testing.TestingT, awsRegion string, bucket string) map[string]string {
	tags, err := GetS3BucketTagsE(t, awsRegion, bucket)
	require.NoError(t, err)

	return tags
}

// GetS3BucketTagsE fetches the given bucket's tags and returns them as a string map of strings.
func GetS3BucketTagsE(t testing.TestingT, awsRegion string, bucket string) (map[string]string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	out, err := s3Client.GetBucketTagging(context.Background(), &s3.GetBucketTaggingInput{
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

// GetS3ObjectContents fetches the contents of the object in the given bucket with the given key and return it as a string.
func GetS3ObjectContents(t testing.TestingT, awsRegion string, bucket string, key string) string {
	contents, err := GetS3ObjectContentsE(t, awsRegion, bucket, key)
	require.NoError(t, err)

	return contents
}

// GetS3ObjectContentsE fetches the contents of the object in the given bucket with the given key and return it as a string.
func GetS3ObjectContentsE(t testing.TestingT, awsRegion string, bucket string, key string) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
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

// CreateS3Bucket creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
func CreateS3Bucket(t testing.TestingT, region string, name string) {
	err := CreateS3BucketE(t, region, name)
	require.NoError(t, err)
}

// CreateS3BucketE creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
func CreateS3BucketE(t testing.TestingT, region string, name string) error {
	logger.Default.Logf(t, "Creating bucket %s in %s", name, region)

	s3Client, err := NewS3ClientE(t, region)
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

	_, err = s3Client.CreateBucket(context.Background(), params)
	return err
}

// PutS3BucketPolicy applies an IAM resource policy to a given S3 bucket to create its bucket policy
func PutS3BucketPolicy(t testing.TestingT, region string, bucketName string, policyJSONString string) {
	err := PutS3BucketPolicyE(t, region, bucketName, policyJSONString)
	require.NoError(t, err)
}

// PutS3BucketPolicyE applies an IAM resource policy to a given S3 bucket to create its bucket policy
func PutS3BucketPolicyE(t testing.TestingT, region string, bucketName string, policyJSONString string) error {
	logger.Default.Logf(t, "Applying bucket policy for bucket %s in %s", bucketName, region)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		return err
	}

	input := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policyJSONString),
	}

	_, err = s3Client.PutBucketPolicy(context.Background(), input)
	return err
}

// PutS3BucketVersioning creates an S3 bucket versioning configuration in the given region against the given bucket name, WITHOUT requiring MFA to remove versioning.
func PutS3BucketVersioning(t testing.TestingT, region string, bucketName string) {
	err := PutS3BucketVersioningE(t, region, bucketName)
	require.NoError(t, err)
}

// PutS3BucketVersioningE creates an S3 bucket versioning configuration in the given region against the given bucket name, WITHOUT requiring MFA to remove versioning.
func PutS3BucketVersioningE(t testing.TestingT, region string, bucketName string) error {
	logger.Default.Logf(t, "Creating bucket versioning configuration for bucket %s in %s", bucketName, region)

	s3Client, err := NewS3ClientE(t, region)
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

	_, err = s3Client.PutBucketVersioning(context.Background(), input)
	return err
}

// DeleteS3Bucket destroys the S3 bucket in the given region with the given name.
func DeleteS3Bucket(t testing.TestingT, region string, name string) {
	err := DeleteS3BucketE(t, region, name)
	require.NoError(t, err)
}

// DeleteS3BucketE destroys the S3 bucket in the given region with the given name.
func DeleteS3BucketE(t testing.TestingT, region string, name string) error {
	logger.Default.Logf(t, "Deleting bucket %s in %s", region, name)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		return err
	}

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err = s3Client.DeleteBucket(context.Background(), params)
	return err
}

// EmptyS3Bucket removes the contents of an S3 bucket in the given region with the given name.
func EmptyS3Bucket(t testing.TestingT, region string, name string) {
	err := EmptyS3BucketE(t, region, name)
	require.NoError(t, err)
}

// EmptyS3BucketE removes the contents of an S3 bucket in the given region with the given name.
func EmptyS3BucketE(t testing.TestingT, region string, name string) error {
	logger.Default.Logf(t, "Emptying bucket %s in %s", name, region)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		return err
	}

	params := &s3.ListObjectVersionsInput{
		Bucket: aws.String(name),
	}

	for {
		// Requesting a batch of objects from s3 bucket
		bucketObjects, err := s3Client.ListObjectVersions(context.Background(), params)
		if err != nil {
			return err
		}

		// Checks if the bucket is already empty
		if len((*bucketObjects).Versions) == 0 {
			logger.Default.Logf(t, "Bucket %s is already empty", name)
			return nil
		}

		// creating an array of pointers of ObjectIdentifier
		objectsToDelete := make([]types.ObjectIdentifier, 0, 1000)
		for _, object := range (*bucketObjects).Versions {
			obj := types.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			}
			objectsToDelete = append(objectsToDelete, obj)
		}

		for _, object := range (*bucketObjects).DeleteMarkers {
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
		_, err = s3Client.DeleteObjects(context.Background(), deleteParams)
		if err != nil {
			return err
		}

		if *(*bucketObjects).IsTruncated { // if there are more objects in the bucket, IsTruncated = true
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

// GetS3BucketLoggingTarget fetches the given bucket's logging target bucket and returns it as a string
func GetS3BucketLoggingTarget(t testing.TestingT, awsRegion string, bucket string) string {
	loggingTarget, err := GetS3BucketLoggingTargetE(t, awsRegion, bucket)
	require.NoError(t, err)

	return loggingTarget
}

// GetS3BucketLoggingTargetE fetches the given bucket's logging target bucket and returns it as the following string:
// `TargetBucket` of the `LoggingEnabled` property for an S3 bucket
func GetS3BucketLoggingTargetE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketLogging(context.Background(), &s3.GetBucketLoggingInput{
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

// GetS3BucketLoggingTargetPrefix fetches the given bucket's logging object prefix and returns it as a string
func GetS3BucketLoggingTargetPrefix(t testing.TestingT, awsRegion string, bucket string) string {
	loggingObjectTargetPrefix, err := GetS3BucketLoggingTargetPrefixE(t, awsRegion, bucket)
	require.NoError(t, err)

	return loggingObjectTargetPrefix
}

// GetS3BucketLoggingTargetPrefixE fetches the given bucket's logging object prefix and returns it as the following string:
// `TargetPrefix` of the `LoggingEnabled` property for an S3 bucket
func GetS3BucketLoggingTargetPrefixE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketLogging(context.Background(), &s3.GetBucketLoggingInput{
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

// GetS3BucketVersioning fetches the given bucket's versioning configuration status and returns it as a string
func GetS3BucketVersioning(t testing.TestingT, awsRegion string, bucket string) string {
	versioningStatus, err := GetS3BucketVersioningE(t, awsRegion, bucket)
	require.NoError(t, err)

	return versioningStatus
}

// GetS3BucketVersioningE fetches the given bucket's versioning configuration status and returns it as a string
func GetS3BucketVersioningE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketVersioning(context.Background(), &s3.GetBucketVersioningInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}

	return string(res.Status), nil
}

// GetS3BucketPolicy fetches the given bucket's resource policy and returns it as a string
func GetS3BucketPolicy(t testing.TestingT, awsRegion string, bucket string) string {
	bucketPolicy, err := GetS3BucketPolicyE(t, awsRegion, bucket)
	require.NoError(t, err)

	return bucketPolicy
}

// GetS3BucketPolicyE fetches the given bucket's resource policy and returns it as a string
func GetS3BucketPolicyE(t testing.TestingT, awsRegion string, bucket string) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetBucketPolicy(context.Background(), &s3.GetBucketPolicyInput{
		Bucket: &bucket,
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(res.Policy), nil
}

// AssertS3BucketExists checks if the given S3 bucket exists in the given region and fail the test if it does not.
func AssertS3BucketExists(t testing.TestingT, region string, name string) {
	err := AssertS3BucketExistsE(t, region, name)
	require.NoError(t, err)
}

// AssertS3BucketExistsE checks if the given S3 bucket exists in the given region and return an error if it does not.
func AssertS3BucketExistsE(t testing.TestingT, region string, name string) error {
	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		return err
	}

	params := &s3.HeadBucketInput{
		Bucket: aws.String(name),
	}
	_, err = s3Client.HeadBucket(context.Background(), params)
	return err
}

// AssertS3BucketVersioningExists checks if the given S3 bucket has a versioning configuration enabled and returns an error if it does not.
func AssertS3BucketVersioningExists(t testing.TestingT, region string, bucketName string) {
	err := AssertS3BucketVersioningExistsE(t, region, bucketName)
	require.NoError(t, err)
}

// AssertS3BucketVersioningExistsE checks if the given S3 bucket has a versioning configuration enabled and returns an error if it does not.
func AssertS3BucketVersioningExistsE(t testing.TestingT, region string, bucketName string) error {
	status, err := GetS3BucketVersioningE(t, region, bucketName)
	if err != nil {
		return err
	}

	if status == "Enabled" {
		return nil
	}
	return NewBucketVersioningNotEnabledError(bucketName, region, status)
}

// AssertS3BucketPolicyExists checks if the given S3 bucket has a resource policy attached and returns an error if it does not
func AssertS3BucketPolicyExists(t testing.TestingT, region string, bucketName string) {
	err := AssertS3BucketPolicyExistsE(t, region, bucketName)
	require.NoError(t, err)
}

// AssertS3BucketPolicyExistsE checks if the given S3 bucket has a resource policy attached and returns an error if it does not
func AssertS3BucketPolicyExistsE(t testing.TestingT, region string, bucketName string) error {
	policy, err := GetS3BucketPolicyE(t, region, bucketName)
	if err != nil {
		return err
	}

	if policy == "" {
		return NewNoBucketPolicyError(bucketName, region, policy)
	}
	return nil
}

// NewS3Client creates an S3 client.
func NewS3Client(t testing.TestingT, region string) *s3.Client {
	client, err := NewS3ClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewS3ClientE creates an S3 client.
func NewS3ClientE(t testing.TestingT, region string) (*s3.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return s3.NewFromConfig(*sess), nil
}

// NewS3Uploader creates an S3 Uploader.
func NewS3Uploader(t testing.TestingT, region string) *manager.Uploader {
	uploader, err := NewS3UploaderE(t, region)
	require.NoError(t, err)
	return uploader
}

// NewS3UploaderE creates an S3 Uploader.
func NewS3UploaderE(t testing.TestingT, region string) (*manager.Uploader, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return manager.NewUploader(s3.NewFromConfig(*sess)), nil
}

// S3AccessLoggingNotEnabledErr is a custom error that occurs when acess logging hasn't been enabled on the S3 Bucket
type S3AccessLoggingNotEnabledErr struct {
	OriginBucket string
	Region       string
}

func (err S3AccessLoggingNotEnabledErr) Error() string {
	return fmt.Sprintf("Server Acess Logging hasn't been enabled for S3 Bucket %s in region %s", err.OriginBucket, err.Region)
}
