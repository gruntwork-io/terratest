package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetDynamoDBTableTagsContextE fetches resource tags of a specified dynamoDB table.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTagsContextE(t testing.TestingT, ctx context.Context, region string, tableName string) ([]types.Tag, error) {
	table, err := GetDynamoDBTableContextE(t, ctx, region, tableName)
	if err != nil {
		return nil, err
	}

	client, err := NewDynamoDBClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	out, err := client.ListTagsOfResource(ctx, &dynamodb.ListTagsOfResourceInput{
		ResourceArn: table.TableArn,
	})
	if err != nil {
		return nil, err
	}

	return out.Tags, nil
}

// GetDynamoDBTableTagsContext fetches resource tags of a specified dynamoDB table. This will fail the test if there are any errors.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTagsContext(t testing.TestingT, ctx context.Context, region string, tableName string) []types.Tag {
	t.Helper()
	tags, err := GetDynamoDBTableTagsContextE(t, ctx, region, tableName)
	require.NoError(t, err)

	return tags
}

// GetDynamoDBTableTags fetches resource tags of a specified dynamoDB table. This will fail the test if there are any errors.
//
// Deprecated: Use [GetDynamoDBTableTagsContext] instead.
func GetDynamoDBTableTags(t testing.TestingT, region string, tableName string) []types.Tag {
	t.Helper()
	return GetDynamoDBTableTagsContext(t, context.Background(), region, tableName)
}

// GetDynamoDBTableTagsE fetches resource tags of a specified dynamoDB table.
//
// Deprecated: Use [GetDynamoDBTableTagsContextE] instead.
func GetDynamoDBTableTagsE(t testing.TestingT, region string, tableName string) ([]types.Tag, error) {
	return GetDynamoDBTableTagsContextE(t, context.Background(), region, tableName)
}

// Deprecated: Use [GetDynamoDBTableTagsContext] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetDynamoDbTableTags(t testing.TestingT, region string, tableName string) []types.Tag {
	t.Helper()
	return GetDynamoDBTableTagsContext(t, context.Background(), region, tableName)
}

// Deprecated: Use [GetDynamoDBTableTagsContextE] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetDynamoDbTableTagsE(t testing.TestingT, region string, tableName string) ([]types.Tag, error) {
	return GetDynamoDBTableTagsContextE(t, context.Background(), region, tableName)
}

// GetDynamoDBTableTimeToLiveContextE fetches information about the TTL configuration of a specified dynamoDB table.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTimeToLiveContextE(t testing.TestingT, ctx context.Context, region string, tableName string) (*types.TimeToLiveDescription, error) {
	client, err := NewDynamoDBClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	out, err := client.DescribeTimeToLive(ctx, &dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	return out.TimeToLiveDescription, nil
}

// GetDynamoDBTableTimeToLiveContext fetches information about the TTL configuration of a specified dynamoDB table. This will fail the test if there are any errors.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTimeToLiveContext(t testing.TestingT, ctx context.Context, region string, tableName string) *types.TimeToLiveDescription {
	t.Helper()
	ttl, err := GetDynamoDBTableTimeToLiveContextE(t, ctx, region, tableName)
	require.NoError(t, err)

	return ttl
}

// GetDynamoDBTableTimeToLive fetches information about the TTL configuration of a specified dynamoDB table. This will fail the test if there are any errors.
//
// Deprecated: Use [GetDynamoDBTableTimeToLiveContext] instead.
func GetDynamoDBTableTimeToLive(t testing.TestingT, region string, tableName string) *types.TimeToLiveDescription {
	t.Helper()
	return GetDynamoDBTableTimeToLiveContext(t, context.Background(), region, tableName)
}

// GetDynamoDBTableTimeToLiveE fetches information about the TTL configuration of a specified dynamoDB table.
//
// Deprecated: Use [GetDynamoDBTableTimeToLiveContextE] instead.
func GetDynamoDBTableTimeToLiveE(t testing.TestingT, region string, tableName string) (*types.TimeToLiveDescription, error) {
	return GetDynamoDBTableTimeToLiveContextE(t, context.Background(), region, tableName)
}

// GetDynamoDBTableContextE fetches information about the specified dynamoDB table.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableContextE(t testing.TestingT, ctx context.Context, region string, tableName string) (*types.TableDescription, error) {
	client, err := NewDynamoDBClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	out, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	return out.Table, nil
}

// GetDynamoDBTableContext fetches information about the specified dynamoDB table. This will fail the test if there are any errors.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableContext(t testing.TestingT, ctx context.Context, region string, tableName string) *types.TableDescription {
	t.Helper()
	table, err := GetDynamoDBTableContextE(t, ctx, region, tableName)
	require.NoError(t, err)

	return table
}

// GetDynamoDBTable fetches information about the specified dynamoDB table. This will fail the test if there are any errors.
//
// Deprecated: Use [GetDynamoDBTableContext] instead.
func GetDynamoDBTable(t testing.TestingT, region string, tableName string) *types.TableDescription {
	t.Helper()
	return GetDynamoDBTableContext(t, context.Background(), region, tableName)
}

// GetDynamoDBTableE fetches information about the specified dynamoDB table.
//
// Deprecated: Use [GetDynamoDBTableContextE] instead.
func GetDynamoDBTableE(t testing.TestingT, region string, tableName string) (*types.TableDescription, error) {
	return GetDynamoDBTableContextE(t, context.Background(), region, tableName)
}

// NewDynamoDBClientContextE creates a DynamoDB client.
// The ctx parameter supports cancellation and timeouts.
func NewDynamoDBClientContextE(t testing.TestingT, ctx context.Context, region string) (*dynamodb.Client, error) {
	sess, err := NewAuthConfigContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(*sess), nil
}

// NewDynamoDBClientContext creates a DynamoDB client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewDynamoDBClientContext(t testing.TestingT, ctx context.Context, region string) *dynamodb.Client {
	t.Helper()
	client, err := NewDynamoDBClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewDynamoDBClient creates a DynamoDB client.
//
// Deprecated: Use [NewDynamoDBClientContext] instead.
func NewDynamoDBClient(t testing.TestingT, region string) *dynamodb.Client {
	t.Helper()
	return NewDynamoDBClientContext(t, context.Background(), region)
}

// NewDynamoDBClientE creates a DynamoDB client.
//
// Deprecated: Use [NewDynamoDBClientContextE] instead.
func NewDynamoDBClientE(t testing.TestingT, region string) (*dynamodb.Client, error) {
	return NewDynamoDBClientContextE(t, context.Background(), region)
}
