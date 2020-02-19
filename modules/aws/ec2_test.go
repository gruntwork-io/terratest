package aws

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEc2InstanceIdsByTag(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ids, err := GetEc2InstanceIdsByTagE(t, region, "Name", fmt.Sprintf("nonexistent-%s", random.UniqueId()))
	require.NoError(t, err)
	assert.Equal(t, 0, len(ids))
}

func TestGetEc2InstanceIdsByFilters(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	filters := map[string][]string{
		"instance-state-name": {"running", "shutting-down"},
		"tag:Name":            {fmt.Sprintf("nonexistent-%s", random.UniqueId())},
	}

	ids, err := GetEc2InstanceIdsByFiltersE(t, region, filters)
	require.NoError(t, err)
	assert.Equal(t, 0, len(ids))
}

func TestExtractTagsFromResource(t *testing.T) {
	t.Parallel()

	env := "test"
	createdBy := "terratest"

	t.Run("verify ExtractTagsFromResource works with EC2 Instance TagSet", func(t *testing.T) {
		instanceName := "web-server"

		ec2InstanceTagSet := []ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(instanceName),
			},
			{
				Key:   aws.String("Env"),
				Value: aws.String(env),
			},
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String(createdBy),
			},
		}

		expected := map[string]string{
			"Name":      instanceName,
			"Env":       env,
			"CreatedBy": createdBy,
		}
		actual := ExtractTagsFromResource(ec2InstanceTagSet)
		assert.True(t, reflect.DeepEqual(expected, actual))
	})

	t.Run("verify ExtractTagsFromResource works with S3 Bucket TagSet", func(t *testing.T) {
		bucketName := "terratest-s3-bucket"

		s3BucketTagSet := []ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(bucketName),
			},
			{
				Key:   aws.String("Env"),
				Value: aws.String(env),
			},
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String(createdBy),
			},
		}

		expected := map[string]string{
			"Name":      bucketName,
			"Env":       env,
			"CreatedBy": createdBy,
		}
		actual := ExtractTagsFromResource(s3BucketTagSet)
		assert.True(t, reflect.DeepEqual(expected, actual))
	})
}
