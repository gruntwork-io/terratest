package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetEcsCluster fetches information about specified ECS cluster.
func GetEcsCluster(t testing.TestingT, region string, name string) *types.Cluster {
	cluster, err := GetEcsClusterE(t, region, name)
	require.NoError(t, err)
	return cluster
}

// GetEcsClusterE fetches information about specified ECS cluster.
func GetEcsClusterE(t testing.TestingT, region string, name string) (*types.Cluster, error) {
	return GetEcsClusterWithIncludeE(t, region, name, []types.ClusterField{})
}

// GetEcsClusterWithInclude fetches extended information about specified ECS cluster.
// The `include` parameter specifies a list of `ecs.ClusterField*` constants, such as `ecs.ClusterFieldTags`.
func GetEcsClusterWithInclude(t testing.TestingT, region string, name string, include []types.ClusterField) *types.Cluster {
	clusterInfo, err := GetEcsClusterWithIncludeE(t, region, name, include)
	require.NoError(t, err)
	return clusterInfo
}

// GetEcsClusterWithIncludeE fetches extended information about specified ECS cluster.
// The `include` parameter specifies a list of `ecs.ClusterField*` constants, such as `ecs.ClusterFieldTags`.
func GetEcsClusterWithIncludeE(t testing.TestingT, region string, name string, include []types.ClusterField) (*types.Cluster, error) {
	client, err := NewEcsClientE(t, region)
	if err != nil {
		return nil, err
	}

	input := &ecs.DescribeClustersInput{
		Clusters: []string{
			name,
		},
		Include: include,
	}
	output, err := client.DescribeClusters(context.Background(), input)
	if err != nil {
		return nil, err
	}

	numClusters := len(output.Clusters)
	if numClusters != 1 {
		return nil, fmt.Errorf("expected to find 1 ECS cluster named '%s' in region '%v', but found '%d'",
			name, region, numClusters)
	}

	return &output.Clusters[0], nil
}

// GetDefaultEcsClusterE fetches information about default ECS cluster.
func GetDefaultEcsClusterE(t testing.TestingT, region string) (*types.Cluster, error) {
	return GetEcsClusterE(t, region, "default")
}

// GetDefaultEcsCluster fetches information about default ECS cluster.
func GetDefaultEcsCluster(t testing.TestingT, region string) *types.Cluster {
	return GetEcsCluster(t, region, "default")
}

// CreateEcsCluster creates ECS cluster in the given region under the given name.
func CreateEcsCluster(t testing.TestingT, region string, name string) *types.Cluster {
	cluster, err := CreateEcsClusterE(t, region, name)
	require.NoError(t, err)
	return cluster
}

// CreateEcsClusterE creates ECS cluster in the given region under the given name.
func CreateEcsClusterE(t testing.TestingT, region string, name string) (*types.Cluster, error) {
	client := NewEcsClient(t, region)
	cluster, err := client.CreateCluster(context.Background(), &ecs.CreateClusterInput{
		ClusterName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	return cluster.Cluster, nil
}

func DeleteEcsCluster(t testing.TestingT, region string, cluster *types.Cluster) {
	err := DeleteEcsClusterE(t, region, cluster)
	require.NoError(t, err)
}

// DeleteEcsClusterE deletes existing ECS cluster in the given region.
func DeleteEcsClusterE(t testing.TestingT, region string, cluster *types.Cluster) error {
	client := NewEcsClient(t, region)
	_, err := client.DeleteCluster(context.Background(), &ecs.DeleteClusterInput{
		Cluster: aws.String(*cluster.ClusterName),
	})
	return err
}

// GetEcsService fetches information about specified ECS service.
func GetEcsService(t testing.TestingT, region string, clusterName string, serviceName string) *types.Service {
	service, err := GetEcsServiceE(t, region, clusterName, serviceName)
	require.NoError(t, err)
	return service
}

// GetEcsServiceE fetches information about specified ECS service.
func GetEcsServiceE(t testing.TestingT, region string, clusterName string, serviceName string) (*types.Service, error) {
	output, err := NewEcsClient(t, region).DescribeServices(context.Background(), &ecs.DescribeServicesInput{
		Cluster: aws.String(clusterName),
		Services: []string{
			serviceName,
		},
	})
	if err != nil {
		return nil, err
	}

	numServices := len(output.Services)
	if numServices != 1 {
		return nil, fmt.Errorf(
			"expected to find 1 ECS service named '%s' in cluster '%s' in region '%v', but found '%d'",
			serviceName, clusterName, region, numServices)
	}
	return &output.Services[0], nil
}

// GetEcsTaskDefinition fetches information about specified ECS task definition.
func GetEcsTaskDefinition(t testing.TestingT, region string, taskDefinition string) *types.TaskDefinition {
	task, err := GetEcsTaskDefinitionE(t, region, taskDefinition)
	require.NoError(t, err)
	return task
}

// GetEcsTaskDefinitionE fetches information about specified ECS task definition.
func GetEcsTaskDefinitionE(t testing.TestingT, region string, taskDefinition string) (*types.TaskDefinition, error) {
	output, err := NewEcsClient(t, region).DescribeTaskDefinition(context.Background(), &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinition),
	})
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition, nil
}

// NewEcsClient creates en ECS client.
func NewEcsClient(t testing.TestingT, region string) *ecs.Client {
	client, err := NewEcsClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewEcsClientE creates an ECS client.
func NewEcsClientE(t testing.TestingT, region string) (*ecs.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return ecs.NewFromConfig(*sess), nil
}
