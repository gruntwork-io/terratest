package aws

import (
	"strings"

	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gruntwork-io/terratest/modules/logger"
)

// FindS3BucketWithTag finds the name of the S3 bucket in the given region with the given tag key=value.
func FindS3BucketWithTag(t *testing.T, awsRegion string, key string, value string, sessExists ...*session.Session) string {
	bucket, err := FindS3BucketWithTagE(t, awsRegion, key, value, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return bucket
}

// FindS3BucketWithTagE finds the name of the S3 bucket in the given region with the given tag key=value.
func FindS3BucketWithTagE(t *testing.T, awsRegion string, key string, value string, sessExists ...*session.Session) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion, sessExists[0])
	if err != nil {
		return "", err
	}

	resp, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return "", err
	}

	for _, bucket := range resp.Buckets {
		tagResponse, err := s3Client.GetBucketTagging(&s3.GetBucketTaggingInput{Bucket: bucket.Name})

		if err != nil {
			if !strings.Contains(err.Error(), "AuthorizationHeaderMalformed") &&
				!strings.Contains(err.Error(), "BucketRegionError") &&
				!strings.Contains(err.Error(), "NoSuchTagSet") {
				return "", err
			}
		}

		for _, tag := range tagResponse.TagSet {
			if *tag.Key == key && *tag.Value == value {
				logger.Logf(t, "Found S3 bucket %s with tag %s=%s", *bucket.Name, key, value)
				return *bucket.Name, nil
			}
		}
	}

	return "", nil
}

// GetS3ObjectContents fetches the contents of the object in the given bucket with the given key and return it as a string.
func GetS3ObjectContents(t *testing.T, awsRegion string, bucket string, key string, sessExists ...*session.Session) string {
	contents, err := GetS3ObjectContentsE(t, awsRegion, bucket, key, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

// GetS3ObjectContentsE fetches the contents of the object in the given bucket with the given key and return it as a string.
func GetS3ObjectContentsE(t *testing.T, awsRegion string, bucket string, key string, sessExists ...*session.Session) (string, error) {
	s3Client, err := NewS3ClientE(t, awsRegion, sessExists[0])
	if err != nil {
		return "", err
	}

	res, err := s3Client.GetObject(&s3.GetObjectInput{
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
	logger.Logf(t, "Read contents from s3://%s/%s", bucket, key)

	return contents, nil
}

// CreateS3Bucket creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
func CreateS3Bucket(t *testing.T, region string, name string, sessExists ...*session.Session) {
	err := CreateS3BucketE(t, region, name, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
}

// CreateS3BucketE creates an S3 bucket in the given region with the given name. Note that S3 bucket names must be globally unique.
func CreateS3BucketE(t *testing.T, region string, name string, sessExists ...*session.Session) error {
	logger.Logf(t, "Creating bucket %s in %s", region, name)

	s3Client, err := NewS3ClientE(t, region, sessExists[0])
	if err != nil {
		return err
	}

	params := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	_, err = s3Client.CreateBucket(params)
	return err
}

// DeleteS3Bucket destroys the S3 bucket in the given region with the given name.
func DeleteS3Bucket(t *testing.T, region string, name string, sessExists ...*session.Session) {
	err := DeleteS3BucketE(t, region, name, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteS3BucketE destroys the S3 bucket in the given region with the given name.
func DeleteS3BucketE(t *testing.T, region string, name string, sessExists ...*session.Session) error {
	logger.Logf(t, "Deleting bucket %s in %s", region, name)

	s3Client, err := NewS3ClientE(t, region, sessExists[0])
	if err != nil {
		return err
	}

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err = s3Client.DeleteBucket(params)
	return err
}

// EmptyS3BucketE removes the contents of an S3 bucket in the given region with the given name.
func EmptyS3Bucket(t *testing.T, region string, name string) {
	err := EmptyS3BucketE(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
}

// EmptyS3BucketE removes the contents of an S3 bucket in the given region with the given name.
func EmptyS3BucketE(t *testing.T, region string, name string) error {
	logger.Logf(t, "Emptying bucket %s in %s", name, region)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		return err
	}

	params := &s3.ListObjectVersionsInput{
		Bucket: aws.String(name),
	}

	for {
		// Requesting a batch of objects from s3 bucket
		bucketObjects, err := s3Client.ListObjectVersions(params)
		if err != nil {
			return err
		}

		//Checks if the bucket is already empty
		if len((*bucketObjects).Versions) == 0 {
			logger.Logf(t, "Bucket %s is already empty", name)
			return nil
		}

		//creating an array of pointers of ObjectIdentifier
		objectsToDelete := make([]*s3.ObjectIdentifier, 0, 1000)
		for _, object := range (*bucketObjects).Versions {
			obj := s3.ObjectIdentifier{
				Key:       object.Key,
				VersionId: object.VersionId,
			}
			objectsToDelete = append(objectsToDelete, &obj)
		}

		//Creating JSON payload for bulk delete
		deleteArray := s3.Delete{Objects: objectsToDelete}
		deleteParams := &s3.DeleteObjectsInput{
			Bucket: aws.String(name),
			Delete: &deleteArray,
		}

		//Running the Bulk delete job (limit 1000)
		_, err = s3Client.DeleteObjects(deleteParams)
		if err != nil {
			return err
		}

		if *(*bucketObjects).IsTruncated { //if there are more objects in the bucket, IsTruncated = true
			// params.Marker = (*deleteParams).Delete.Objects[len((*deleteParams).Delete.Objects)-1].Key
			params.KeyMarker = bucketObjects.NextKeyMarker
			logger.Logf(t, "Requesting next batch | %s", *(params.KeyMarker))
		} else { //if all objects in the bucket have been cleaned up.
			break
		}
	}
	logger.Logf(t, "Bucket %s is now empty", name)
	return err
}

// AssertS3BucketExists checks if the given S3 bucket exists in the given region and fail the test if it does not.
func AssertS3BucketExists(t *testing.T, region string, name string, sessExists ...*session.Session) {
	err := AssertS3BucketExistsE(t, region, name, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
}

// AssertS3BucketExistsE checks if the given S3 bucket exists in the given region and return an error if it does not.
func AssertS3BucketExistsE(t *testing.T, region string, name string, sessExists ...*session.Session) error {
	s3Client, err := NewS3ClientE(t, region, sessExists[0])
	if err != nil {
		return err
	}

	params := &s3.HeadBucketInput{
		Bucket: aws.String(name),
	}
	_, err = s3Client.HeadBucket(params)
	return err
}

// NewS3Client creates an S3 client.
func NewS3Client(t *testing.T, region string, sessExists ...*session.Session) *s3.S3 {
	client, err := NewS3ClientE(t, region, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewS3ClientE creates an S3 client.
func NewS3ClientE(t *testing.T, region string, sessExists ...*session.Session) (*s3.S3, error) {
	sess, err := NewAuthenticatedSession(region, sessExists[0])
	if err != nil {
		return nil, err
	}

	return s3.New(sess), nil
}

// NewS3Uploader creates an S3 Uploader.
func NewS3Uploader(t *testing.T, region string) *s3manager.Uploader {
	uploader, err := NewS3UploaderE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return uploader
}

// NewS3UploaderE creates an S3 Uploader.
func NewS3UploaderE(t *testing.T, region string) (*s3manager.Uploader, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return s3manager.NewUploader(sess), nil
}
