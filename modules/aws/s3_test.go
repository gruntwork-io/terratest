// Integration tests that validate S3-related code in AWS.
package aws

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndDestroyS3Bucket(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)

	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateS3Bucket(t, region, s3BucketName)
	DeleteS3Bucket(t, region, s3BucketName)
}

func TestAssertS3BucketExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	s3BucketName := "gruntwork-terratest-" + strings.ToLower(random.UniqueId())
	logger.Default.Logf(t, "Random values selected. Region = %s, s3BucketName = %s\n", region, s3BucketName)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)

	AssertS3BucketExists(t, region, s3BucketName)
}

func TestAssertS3BucketExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	s3BucketName := "gruntwork-terratest-" + strings.ToLower(random.UniqueId())
	logger.Default.Logf(t, "Random values selected. Region = %s, s3BucketName = %s\n", region, s3BucketName)

	// We elect not to create the S3 bucket to confirm that our function correctly reports it doesn't exist.
	// aws.CreateS3Bucket(region, s3BucketName)

	err := AssertS3BucketExistsE(t, region, s3BucketName)
	if err == nil {
		t.Fatalf("Function claimed that S3 Bucket '%s' exists, but in fact it does not.", s3BucketName)
	}
}

func TestAssertS3BucketVersioningEnabled(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	s3BucketName := "gruntwork-terratest-" + strings.ToLower(random.UniqueId())
	logger.Default.Logf(t, "Random values selected. Region = %s, s3BucketName = %s\n", region, s3BucketName)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)
	PutS3BucketVersioning(t, region, s3BucketName)

	AssertS3BucketVersioningExists(t, region, s3BucketName)
}

func TestEmptyS3Bucket(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)

	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}

	testEmptyBucket(t, s3Client, region, s3BucketName)
}

func TestEmptyS3BucketVersioned(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)

	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)

	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}

	versionInput := &s3.PutBucketVersioningInput{
		Bucket: aws.String(s3BucketName),
		VersioningConfiguration: &types.VersioningConfiguration{
			MFADelete: types.MFADeleteDisabled,
			Status:    types.BucketVersioningStatusEnabled,
		},
	}

	_, err = s3Client.PutBucketVersioning(context.Background(), versionInput)
	if err != nil {
		t.Fatal(err)
	}

	testEmptyBucket(t, s3Client, region, s3BucketName)
}

func TestAssertS3BucketPolicyExists(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)

	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)

	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)
	exampleBucketPolicy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Deny","Principal":{"AWS":["*"]},"Action":"s3:Get*","Resource":"arn:aws:s3:::%s/*","Condition":{"Bool":{"aws:SecureTransport":"false"}}}]}`, s3BucketName)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)
	PutS3BucketPolicy(t, region, s3BucketName, exampleBucketPolicy)

	AssertS3BucketPolicyExists(t, region, s3BucketName)

}

func TestGetS3BucketTags(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)
	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)

	s3Client, err := NewS3ClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s3Client.PutBucketTagging(context.Background(), &s3.PutBucketTaggingInput{
		Bucket: &s3BucketName,
		Tagging: &types.Tagging{
			TagSet: []types.Tag{
				{
					Key:   aws.String("Key1"),
					Value: aws.String("Value1"),
				},
				{
					Key:   aws.String("Key2"),
					Value: aws.String("Value2"),
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	actualTags := GetS3BucketTags(t, region, s3BucketName)
	assert.True(t, actualTags["Key1"] == "Value1")
	assert.True(t, actualTags["Key2"] == "Value2")
	assert.True(t, actualTags["NonExistentKey"] == "")
}

func testEmptyBucket(t *testing.T, s3Client *s3.Client, region string, s3BucketName string) {
	expectedFileCount := rand.Intn(1000)
	logger.Default.Logf(t, "Uploading %s files to bucket %s", strconv.Itoa(expectedFileCount), s3BucketName)

	deleted := 0

	// Upload expectedFileCount files
	for i := 1; i <= expectedFileCount; i++ {
		key := fmt.Sprintf("test-%s", strconv.Itoa(i))
		body := strings.NewReader("This is the body")

		params := &s3.PutObjectInput{
			Bucket: aws.String(s3BucketName),
			Key:    &key,
			Body:   body,
		}

		uploader := NewS3Uploader(t, region)

		_, err := uploader.Upload(context.Background(), params)
		if err != nil {
			t.Fatal(err)
		}

		// Delete the first 10 files to be able to test if all files, including delete markers are deleted
		if i < 10 {
			_, err := s3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
				Bucket: aws.String(s3BucketName),
				Key:    aws.String(key),
			})
			if err != nil {
				t.Fatal(err)
			}
			deleted++
		}

		if i != 0 && i%100 == 0 {
			logger.Default.Logf(t, "Uploaded %s files to bucket %s successfully", strconv.Itoa(i), s3BucketName)
		}
	}

	logger.Default.Logf(t, "Uploaded %s files to bucket %s successfully", strconv.Itoa(expectedFileCount), s3BucketName)

	// verify bucket contains 1 file now
	listObjectsParams := &s3.ListObjectsV2Input{
		Bucket: aws.String(s3BucketName),
	}

	logger.Default.Logf(t, "Verifying %s files were uploaded to bucket %s", strconv.Itoa(expectedFileCount), s3BucketName)
	actualCount := 0
	for {
		bucketObjects, err := s3Client.ListObjectsV2(context.Background(), listObjectsParams)
		if err != nil {
			t.Fatal(err)
		}

		pageLength := len((*bucketObjects).Contents)
		actualCount += pageLength

		if !*bucketObjects.IsTruncated {
			break
		}

		listObjectsParams.ContinuationToken = bucketObjects.NextContinuationToken
	}

	require.Equal(t, expectedFileCount-deleted, actualCount)

	// empty bucket
	logger.Default.Logf(t, "Emptying bucket %s", s3BucketName)
	EmptyS3Bucket(t, region, s3BucketName)

	// verify the bucket is empty
	bucketObjects, err := s3Client.ListObjectsV2(context.Background(), listObjectsParams)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, 0, len((*bucketObjects).Contents))
}

func TestGetS3BucketOwnershipControls(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)

	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)
	CreateS3Bucket(t, region, s3BucketName)
	t.Cleanup(func() {
		DeleteS3Bucket(t, region, s3BucketName)
	})

	t.Run("Exist", func(t *testing.T) {
		s3Client, err := NewS3ClientE(t, region)
		require.NoError(t, err)
		_, err = s3Client.PutBucketOwnershipControls(context.Background(), &s3.PutBucketOwnershipControlsInput{
			Bucket: &s3BucketName,
			OwnershipControls: &types.OwnershipControls{
				Rules: []types.OwnershipControlsRule{
					{
						ObjectOwnership: types.ObjectOwnershipBucketOwnerEnforced,
					},
				},
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			_, err := s3Client.DeleteBucketOwnershipControls(context.Background(), &s3.DeleteBucketOwnershipControlsInput{
				Bucket: &s3BucketName,
			})
			require.NoError(t, err)
		})

		controls := GetS3BucketOwnershipControls(t, region, s3BucketName)
		assert.Equal(t, 1, len(controls))
		assert.Equal(t, string(types.ObjectOwnershipBucketOwnerEnforced), controls[0])
	})

	t.Run("NotExist", func(t *testing.T) {
		_, err := GetS3BucketOwnershipControlsE(t, region, s3BucketName)
		assert.Error(t, err)
	})
}

func TestS3ObjectContents(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	id := random.UniqueId()
	logger.Default.Logf(t, "Random values selected. Region = %s, Id = %s\n", region, id)
	s3BucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateS3Bucket(t, region, s3BucketName)
	defer DeleteS3Bucket(t, region, s3BucketName)
	defer EmptyS3BucketE(t, region, s3BucketName)

	key := fmt.Sprintf("content-%s", id)
	body := make([]byte, 1024)
	rand.Read(body)

	PutS3ObjectContentsE(t, region, s3BucketName, key, bytes.NewReader(body))
	storedBody := GetS3ObjectContents(t, region, s3BucketName, key)

	assert.Equal(t, body, []byte(storedBody))
}
