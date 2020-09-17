package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetEcsCluster fetches information about specified ECS cluster.
func GetEcsCluster(t testing.TestingT, region string, name string) *ecs.Cluster {
	cluster, err := GetEcsClusterE(t, region, name)
	require.NoError(t, err)
	return cluster
}

// GetEcsClusterE fetches information about specified ECS cluster.
func GetEcsClusterE(t testing.TestingT, region string, name string) (*ecs.Cluster, error) {
	client, err := NewEcsClientE(t, region)
	if err != nil {
		return nil, err
	}
	input := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(name),
		},
	}
	output, err := client.DescribeClusters(input)
	if err != nil {
		return nil, err
	}

	numClusters := len(output.Clusters)
	if numClusters != 1 {
		return nil, fmt.Errorf("Expected to find 1 ECS cluster named '%s' in region '%v', but found '%d'",
			name, region, numClusters)
	}

	return output.Clusters[0], nil
}

// GetDefaultEcsClusterE fetches information about default ECS cluster.
func GetDefaultEcsClusterE(t testing.TestingT, region string) (*ecs.Cluster, error) {
	return GetEcsClusterE(t, region, "default")
}

// GetDefaultEcsCluster fetches information about default ECS cluster.
func GetDefaultEcsCluster(t testing.TestingT, region string) *ecs.Cluster {
	return GetEcsCluster(t, region, "default")
}

// CreateEcsCluster creates ECS cluster in the given region under the given name.
func CreateEcsCluster(t testing.TestingT, region string, name string) *ecs.Cluster {
	cluster, err := CreateEcsClusterE(t, region, name)
	require.NoError(t, err)
	return cluster
}

// CreateEcsClusterE creates ECS cluster in the given region under the given name.
func CreateEcsClusterE(t testing.TestingT, region string, name string) (*ecs.Cluster, error) {
	client := NewEcsClient(t, region)
	cluster, err := client.CreateCluster(&ecs.CreateClusterInput{
		ClusterName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}
	return cluster.Cluster, nil
}

func DeleteEcsCluster(t testing.TestingT, region string, cluster *ecs.Cluster) {
	err := DeleteEcsClusterE(t, region, cluster)
	require.NoError(t, err)
}

// DeleteEcsClusterE deletes existing ECS cluster in the given region.
func DeleteEcsClusterE(t testing.TestingT, region string, cluster *ecs.Cluster) error {
	client := NewEcsClient(t, region)
	_, err := client.DeleteCluster(&ecs.DeleteClusterInput{
		Cluster: aws.String(*cluster.ClusterName),
	})
	return err
}

// GetEcsService fetches information about specified ECS service.
func GetEcsService(t testing.TestingT, region string, clusterName string, serviceName string) *ecs.Service {
	service, err := GetEcsServiceE(t, region, clusterName, serviceName)
	require.NoError(t, err)
	return service
}

// GetEcsServiceE fetches information about specified ECS service.
func GetEcsServiceE(t testing.TestingT, region string, clusterName string, serviceName string) (*ecs.Service, error) {
	output, err := NewEcsClient(t, region).DescribeServices(&ecs.DescribeServicesInput{
		Cluster: aws.String(clusterName),
		Services: []*string{
			aws.String(serviceName),
		},
	})
	if err != nil {
		return nil, err
	}

	numServices := len(output.Services)
	if numServices != 1 {
		return nil, fmt.Errorf(
			"Expected to find 1 ECS service named '%s' in cluster '%s' in region '%v', but found '%d'",
			serviceName, clusterName, region, numServices)
	}
	return output.Services[0], nil
}

// GetEcsTaskDefinition fetches information about specified ECS task definition.
func GetEcsTaskDefinition(t testing.TestingT, region string, taskDefinition string) *ecs.TaskDefinition {
	task, err := GetEcsTaskDefinitionE(t, region, taskDefinition)
	require.NoError(t, err)
	return task
}

// GetEcsTaskDefinitionE fetches information about specified ECS task definition.
func GetEcsTaskDefinitionE(t testing.TestingT, region string, taskDefinition string) (*ecs.TaskDefinition, error) {
	output, err := NewEcsClient(t, region).DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinition),
	})
	if err != nil {
		return nil, err
	}
	return output.TaskDefinition, nil
}

// NewEcsClient creates en ECS client.
func NewEcsClient(t testing.TestingT, region string) *ecs.ECS {
	client, err := NewEcsClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewEcsClientE creates an ECS client.
func NewEcsClientE(t testing.TestingT, region string) (*ecs.ECS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return ecs.New(sess), nil
}
