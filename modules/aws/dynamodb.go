package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DynamoDBAPI is the subset of *dynamodb.Client operations used by the helpers in this file.
// Declared as an interface so tests can substitute a mock; a real *dynamodb.Client satisfies it
// automatically.
type DynamoDBAPI interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	DescribeTimeToLive(ctx context.Context, params *dynamodb.DescribeTimeToLiveInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTimeToLiveOutput, error)
	ListTagsOfResource(ctx context.Context, params *dynamodb.ListTagsOfResourceInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ListTagsOfResourceOutput, error)
}

// GetDynamoDBTableTagsContextE fetches resource tags of a specified dynamoDB table.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTagsContextE(t testing.TestingT, ctx context.Context, region string, tableName string) ([]types.Tag, error) {
	client, err := NewDynamoDBClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDynamoDBTableTagsWithClientContextE(t, ctx, client, tableName)
}

// GetDynamoDBTableTagsWithClientContextE fetches resource tags of a specified dynamoDB table using
// the provided DynamoDB client.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTagsWithClientContextE(t testing.TestingT, ctx context.Context, client DynamoDBAPI, tableName string) ([]types.Tag, error) {
	table, err := GetDynamoDBTableWithClientContextE(t, ctx, client, tableName)
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

// GetDynamoDbTableTags fetches resource tags of a specified dynamoDB table. This will fail the test if there are any errors.
//
// Deprecated: Use [GetDynamoDBTableTagsContext] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetDynamoDbTableTags(t testing.TestingT, region string, tableName string) []types.Tag {
	t.Helper()
	return GetDynamoDBTableTagsContext(t, context.Background(), region, tableName)
}

// GetDynamoDbTableTagsE fetches resource tags of a specified dynamoDB table.
//
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

	return GetDynamoDBTableTimeToLiveWithClientContextE(t, ctx, client, tableName)
}

// GetDynamoDBTableTimeToLiveWithClientContextE fetches the TTL configuration of a specified
// dynamoDB table using the provided DynamoDB client.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableTimeToLiveWithClientContextE(t testing.TestingT, ctx context.Context, client DynamoDBAPI, tableName string) (*types.TimeToLiveDescription, error) {
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

	return GetDynamoDBTableWithClientContextE(t, ctx, client, tableName)
}

// GetDynamoDBTableWithClientContextE fetches information about the specified dynamoDB table using
// the provided DynamoDB client.
// The ctx parameter supports cancellation and timeouts.
func GetDynamoDBTableWithClientContextE(t testing.TestingT, ctx context.Context, client DynamoDBAPI, tableName string) (*types.TableDescription, error) {
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
	sess, err := NewAuthenticatedSessionContext(ctx, region)
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
