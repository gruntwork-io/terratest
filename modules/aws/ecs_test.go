package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestEcsCluster(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	c1, err := CreateEcsClusterE(t, region, "terratest")
	defer DeleteEcsCluster(t, region, c1)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *c1.ClusterName)

	c2, err := GetEcsClusterE(t, region, *c1.ClusterName)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *c2.ClusterName)
}

func TestEcsClusterWithInclude(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	clusterName := "terratest-" + random.UniqueId()
	tags := []types.Tag{{
		Key:   aws.String("test-tag"),
		Value: aws.String("hello-world"),
	}}

	client := NewEcsClient(t, region)
	c1, err := client.CreateCluster(context.Background(), &ecs.CreateClusterInput{
		ClusterName: aws.String(clusterName),
		Tags:        tags,
	})
	assert.NoError(t, err)

	defer DeleteEcsCluster(t, region, c1.Cluster)

	assert.Equal(t, clusterName, aws.ToString(c1.Cluster.ClusterName))

	c2, err := GetEcsClusterWithIncludeE(t, region, clusterName, []types.ClusterField{types.ClusterFieldTags})
	assert.NoError(t, err)

	assert.Equal(t, clusterName, aws.ToString(c2.ClusterName))
	assert.Equal(t, tags, c2.Tags)
	assert.Empty(t, c2.Statistics)

	c3, err := GetEcsClusterWithIncludeE(t, region, clusterName, []types.ClusterField{types.ClusterFieldStatistics})
	assert.NoError(t, err)

	assert.Equal(t, clusterName, aws.ToString(c3.ClusterName))
	assert.NotEmpty(t, c3.Statistics)
	assert.Empty(t, c3.Tags)
}
