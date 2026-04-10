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

// GetEcsClusterWithIncludeContextE fetches extended information about specified ECS cluster.
// The `include` parameter specifies a list of `ecs.ClusterField*` constants, such as `ecs.ClusterFieldTags`.
// The ctx parameter supports cancellation and timeouts.
func GetEcsClusterWithIncludeContextE(t testing.TestingT, ctx context.Context, region string, name string, include []types.ClusterField) (*types.Cluster, error) {
	client, err := NewEcsClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	input := &ecs.DescribeClustersInput{
		Clusters: []string{
			name,
		},
		Include: include,
	}

	output, err := client.DescribeClusters(ctx, input)
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

// GetEcsClusterWithIncludeContext fetches extended information about specified ECS cluster.
// The `include` parameter specifies a list of `ecs.ClusterField*` constants, such as `ecs.ClusterFieldTags`.
// The ctx parameter supports cancellation and timeouts.
func GetEcsClusterWithIncludeContext(t testing.TestingT, ctx context.Context, region string, name string, include []types.ClusterField) *types.Cluster {
	t.Helper()
	clusterInfo, err := GetEcsClusterWithIncludeContextE(t, ctx, region, name, include)
	require.NoError(t, err)

	return clusterInfo
}

// GetEcsClusterWithInclude fetches extended information about specified ECS cluster.
// The `include` parameter specifies a list of `ecs.ClusterField*` constants, such as `ecs.ClusterFieldTags`.
//
// Deprecated: Use [GetEcsClusterWithIncludeContext] instead.
func GetEcsClusterWithInclude(t testing.TestingT, region string, name string, include []types.ClusterField) *types.Cluster {
	t.Helper()
	return GetEcsClusterWithIncludeContext(t, context.Background(), region, name, include)
}

// GetEcsClusterWithIncludeE fetches extended information about specified ECS cluster.
// The `include` parameter specifies a list of `ecs.ClusterField*` constants, such as `ecs.ClusterFieldTags`.
//
// Deprecated: Use [GetEcsClusterWithIncludeContextE] instead.
func GetEcsClusterWithIncludeE(t testing.TestingT, region string, name string, include []types.ClusterField) (*types.Cluster, error) {
	return GetEcsClusterWithIncludeContextE(t, context.Background(), region, name, include)
}

// GetEcsClusterContextE fetches information about specified ECS cluster.
// The ctx parameter supports cancellation and timeouts.
func GetEcsClusterContextE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Cluster, error) {
	return GetEcsClusterWithIncludeContextE(t, ctx, region, name, []types.ClusterField{})
}

// GetEcsClusterContext fetches information about specified ECS cluster.
// The ctx parameter supports cancellation and timeouts.
func GetEcsClusterContext(t testing.TestingT, ctx context.Context, region string, name string) *types.Cluster {
	t.Helper()
	cluster, err := GetEcsClusterContextE(t, ctx, region, name)
	require.NoError(t, err)

	return cluster
}

// GetEcsCluster fetches information about specified ECS cluster.
//
// Deprecated: Use [GetEcsClusterContext] instead.
func GetEcsCluster(t testing.TestingT, region string, name string) *types.Cluster {
	t.Helper()
	return GetEcsClusterContext(t, context.Background(), region, name)
}

// GetEcsClusterE fetches information about specified ECS cluster.
//
// Deprecated: Use [GetEcsClusterContextE] instead.
func GetEcsClusterE(t testing.TestingT, region string, name string) (*types.Cluster, error) {
	return GetEcsClusterContextE(t, context.Background(), region, name)
}

// GetDefaultEcsClusterContextE fetches information about default ECS cluster.
// The ctx parameter supports cancellation and timeouts.
func GetDefaultEcsClusterContextE(t testing.TestingT, ctx context.Context, region string) (*types.Cluster, error) {
	return GetEcsClusterContextE(t, ctx, region, "default")
}

// GetDefaultEcsClusterContext fetches information about default ECS cluster.
// The ctx parameter supports cancellation and timeouts.
func GetDefaultEcsClusterContext(t testing.TestingT, ctx context.Context, region string) *types.Cluster {
	t.Helper()
	return GetEcsClusterContext(t, ctx, region, "default")
}

// GetDefaultEcsClusterE fetches information about default ECS cluster.
//
// Deprecated: Use [GetDefaultEcsClusterContextE] instead.
func GetDefaultEcsClusterE(t testing.TestingT, region string) (*types.Cluster, error) {
	return GetDefaultEcsClusterContextE(t, context.Background(), region)
}

// GetDefaultEcsCluster fetches information about default ECS cluster.
//
// Deprecated: Use [GetDefaultEcsClusterContext] instead.
func GetDefaultEcsCluster(t testing.TestingT, region string) *types.Cluster {
	t.Helper()
	return GetDefaultEcsClusterContext(t, context.Background(), region)
}

// CreateEcsClusterContextE creates ECS cluster in the given region under the given name.
// The ctx parameter supports cancellation and timeouts.
func CreateEcsClusterContextE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Cluster, error) {
	client, err := NewEcsClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	cluster, err := client.CreateCluster(ctx, &ecs.CreateClusterInput{
		ClusterName: aws.String(name),
	})
	if err != nil {
		return nil, err
	}

	return cluster.Cluster, nil
}

// CreateEcsClusterContext creates ECS cluster in the given region under the given name.
// The ctx parameter supports cancellation and timeouts.
func CreateEcsClusterContext(t testing.TestingT, ctx context.Context, region string, name string) *types.Cluster {
	t.Helper()
	cluster, err := CreateEcsClusterContextE(t, ctx, region, name)
	require.NoError(t, err)

	return cluster
}

// CreateEcsCluster creates ECS cluster in the given region under the given name.
//
// Deprecated: Use [CreateEcsClusterContext] instead.
func CreateEcsCluster(t testing.TestingT, region string, name string) *types.Cluster {
	t.Helper()
	return CreateEcsClusterContext(t, context.Background(), region, name)
}

// CreateEcsClusterE creates ECS cluster in the given region under the given name.
//
// Deprecated: Use [CreateEcsClusterContextE] instead.
func CreateEcsClusterE(t testing.TestingT, region string, name string) (*types.Cluster, error) {
	return CreateEcsClusterContextE(t, context.Background(), region, name)
}

// DeleteEcsClusterContextE deletes existing ECS cluster in the given region.
// The ctx parameter supports cancellation and timeouts.
func DeleteEcsClusterContextE(t testing.TestingT, ctx context.Context, region string, cluster *types.Cluster) error {
	client, err := NewEcsClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	_, err = client.DeleteCluster(ctx, &ecs.DeleteClusterInput{
		Cluster: aws.String(*cluster.ClusterName),
	})

	return err
}

// DeleteEcsClusterContext deletes existing ECS cluster in the given region.
// The ctx parameter supports cancellation and timeouts.
func DeleteEcsClusterContext(t testing.TestingT, ctx context.Context, region string, cluster *types.Cluster) {
	t.Helper()
	err := DeleteEcsClusterContextE(t, ctx, region, cluster)
	require.NoError(t, err)
}

// DeleteEcsCluster deletes existing ECS cluster in the given region.
//
// Deprecated: Use [DeleteEcsClusterContext] instead.
func DeleteEcsCluster(t testing.TestingT, region string, cluster *types.Cluster) {
	t.Helper()
	DeleteEcsClusterContext(t, context.Background(), region, cluster)
}

// DeleteEcsClusterE deletes existing ECS cluster in the given region.
//
// Deprecated: Use [DeleteEcsClusterContextE] instead.
func DeleteEcsClusterE(t testing.TestingT, region string, cluster *types.Cluster) error {
	return DeleteEcsClusterContextE(t, context.Background(), region, cluster)
}

// GetEcsServiceContextE fetches information about specified ECS service.
// The ctx parameter supports cancellation and timeouts.
func GetEcsServiceContextE(t testing.TestingT, ctx context.Context, region string, clusterName string, serviceName string) (*types.Service, error) {
	client, err := NewEcsClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	output, err := client.DescribeServices(ctx, &ecs.DescribeServicesInput{
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

// GetEcsServiceContext fetches information about specified ECS service.
// The ctx parameter supports cancellation and timeouts.
func GetEcsServiceContext(t testing.TestingT, ctx context.Context, region string, clusterName string, serviceName string) *types.Service {
	t.Helper()
	service, err := GetEcsServiceContextE(t, ctx, region, clusterName, serviceName)
	require.NoError(t, err)

	return service
}

// GetEcsService fetches information about specified ECS service.
//
// Deprecated: Use [GetEcsServiceContext] instead.
func GetEcsService(t testing.TestingT, region string, clusterName string, serviceName string) *types.Service {
	t.Helper()
	return GetEcsServiceContext(t, context.Background(), region, clusterName, serviceName)
}

// GetEcsServiceE fetches information about specified ECS service.
//
// Deprecated: Use [GetEcsServiceContextE] instead.
func GetEcsServiceE(t testing.TestingT, region string, clusterName string, serviceName string) (*types.Service, error) {
	return GetEcsServiceContextE(t, context.Background(), region, clusterName, serviceName)
}

// GetEcsTaskDefinitionContextE fetches information about specified ECS task definition.
// The ctx parameter supports cancellation and timeouts.
func GetEcsTaskDefinitionContextE(t testing.TestingT, ctx context.Context, region string, taskDefinition string) (*types.TaskDefinition, error) {
	client, err := NewEcsClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	output, err := client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinition),
	})
	if err != nil {
		return nil, err
	}

	return output.TaskDefinition, nil
}

// GetEcsTaskDefinitionContext fetches information about specified ECS task definition.
// The ctx parameter supports cancellation and timeouts.
func GetEcsTaskDefinitionContext(t testing.TestingT, ctx context.Context, region string, taskDefinition string) *types.TaskDefinition {
	t.Helper()
	task, err := GetEcsTaskDefinitionContextE(t, ctx, region, taskDefinition)
	require.NoError(t, err)

	return task
}

// GetEcsTaskDefinition fetches information about specified ECS task definition.
//
// Deprecated: Use [GetEcsTaskDefinitionContext] instead.
func GetEcsTaskDefinition(t testing.TestingT, region string, taskDefinition string) *types.TaskDefinition {
	t.Helper()
	return GetEcsTaskDefinitionContext(t, context.Background(), region, taskDefinition)
}

// GetEcsTaskDefinitionE fetches information about specified ECS task definition.
//
// Deprecated: Use [GetEcsTaskDefinitionContextE] instead.
func GetEcsTaskDefinitionE(t testing.TestingT, region string, taskDefinition string) (*types.TaskDefinition, error) {
	return GetEcsTaskDefinitionContextE(t, context.Background(), region, taskDefinition)
}

// NewEcsClientContextE creates an ECS client.
// The ctx parameter supports cancellation and timeouts.
func NewEcsClientContextE(t testing.TestingT, ctx context.Context, region string) (*ecs.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return ecs.NewFromConfig(*sess), nil
}

// NewEcsClientContext creates an ECS client.
// The ctx parameter supports cancellation and timeouts.
func NewEcsClientContext(t testing.TestingT, ctx context.Context, region string) *ecs.Client {
	t.Helper()
	client, err := NewEcsClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewEcsClient creates en ECS client.
//
// Deprecated: Use [NewEcsClientContext] instead.
func NewEcsClient(t testing.TestingT, region string) *ecs.Client {
	t.Helper()
	return NewEcsClientContext(t, context.Background(), region)
}

// NewEcsClientE creates an ECS client.
//
// Deprecated: Use [NewEcsClientContextE] instead.
func NewEcsClientE(t testing.TestingT, region string) (*ecs.Client, error) {
	return NewEcsClientContextE(t, context.Background(), region)
}
