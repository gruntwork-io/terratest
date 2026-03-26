package aws_test

import (
	"context"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestEcsCluster(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)

	c1, err := aws.CreateEcsClusterE(t, region, "terratest")
	defer aws.DeleteEcsCluster(t, region, c1)

	require.NoError(t, err)
	assert.Equal(t, "terratest", *c1.ClusterName)

	c2, err := aws.GetEcsClusterE(t, region, *c1.ClusterName)

	require.NoError(t, err)
	assert.Equal(t, "terratest", *c2.ClusterName)
}

func TestEcsClusterWithInclude(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	clusterName := "terratest-" + random.UniqueID()
	tags := []types.Tag{{
		Key:   awsSDK.String("test-tag"),
		Value: awsSDK.String("hello-world"),
	}}

	client := aws.NewEcsClient(t, region)
	c1, err := client.CreateCluster(context.Background(), &ecs.CreateClusterInput{
		ClusterName: awsSDK.String(clusterName),
		Tags:        tags,
	})
	require.NoError(t, err)

	defer aws.DeleteEcsCluster(t, region, c1.Cluster)

	assert.Equal(t, clusterName, awsSDK.ToString(c1.Cluster.ClusterName))

	c2, err := aws.GetEcsClusterWithIncludeE(t, region, clusterName, []types.ClusterField{types.ClusterFieldTags})
	require.NoError(t, err)

	assert.Equal(t, clusterName, awsSDK.ToString(c2.ClusterName))
	assert.Equal(t, tags, c2.Tags)
	assert.Empty(t, c2.Statistics)

	c3, err := aws.GetEcsClusterWithIncludeE(t, region, clusterName, []types.ClusterField{types.ClusterFieldStatistics})
	require.NoError(t, err)

	assert.Equal(t, clusterName, awsSDK.ToString(c3.ClusterName))
	assert.NotEmpty(t, c3.Statistics)
	assert.Empty(t, c3.Tags)
}
