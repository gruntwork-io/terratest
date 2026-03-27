package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetDynamoDBTableTags fetches resource tags of a specified dynamoDB table. This will fail the test if there are any errors.
func GetDynamoDBTableTags(t testing.TestingT, region string, tableName string) []types.Tag {
	tags, err := GetDynamoDBTableTagsE(t, region, tableName)
	require.NoError(t, err)

	return tags
}

// GetDynamoDBTableTagsE fetches resource tags of a specified dynamoDB table.
func GetDynamoDBTableTagsE(t testing.TestingT, region string, tableName string) ([]types.Tag, error) {
	table, err := GetDynamoDBTableE(t, region, tableName)
	if err != nil {
		return nil, err
	}

	client, err := NewDynamoDBClientE(t, region)
	if err != nil {
		return nil, err
	}

	out, err := client.ListTagsOfResource(context.Background(), &dynamodb.ListTagsOfResourceInput{
		ResourceArn: table.TableArn,
	})
	if err != nil {
		return nil, err
	}

	return out.Tags, err
}

// Deprecated: Use GetDynamoDBTableTags instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetDynamoDbTableTags(t testing.TestingT, region string, tableName string) []types.Tag {
	return GetDynamoDBTableTags(t, region, tableName)
}

// Deprecated: Use GetDynamoDBTableTagsE instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetDynamoDbTableTagsE(t testing.TestingT, region string, tableName string) ([]types.Tag, error) {
	return GetDynamoDBTableTagsE(t, region, tableName)
}

// GetDynamoDBTableTimeToLive fetches information about the TTL configuration of a specified dynamoDB table. This will fail the test if there are any errors.
func GetDynamoDBTableTimeToLive(t testing.TestingT, region string, tableName string) *types.TimeToLiveDescription {
	ttl, err := GetDynamoDBTableTimeToLiveE(t, region, tableName)
	require.NoError(t, err)

	return ttl
}

// GetDynamoDBTableTimeToLiveE fetches information about the TTL configuration of a specified dynamoDB table.
func GetDynamoDBTableTimeToLiveE(t testing.TestingT, region string, tableName string) (*types.TimeToLiveDescription, error) {
	client, err := NewDynamoDBClientE(t, region)
	if err != nil {
		return nil, err
	}

	out, err := client.DescribeTimeToLive(context.Background(), &dynamodb.DescribeTimeToLiveInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	return out.TimeToLiveDescription, err
}

// GetDynamoDBTable fetches information about the specified dynamoDB table. This will fail the test if there are any errors.
func GetDynamoDBTable(t testing.TestingT, region string, tableName string) *types.TableDescription {
	table, err := GetDynamoDBTableE(t, region, tableName)
	require.NoError(t, err)

	return table
}

// GetDynamoDBTableE fetches information about the specified dynamoDB table.
func GetDynamoDBTableE(t testing.TestingT, region string, tableName string) (*types.TableDescription, error) {
	client, err := NewDynamoDBClientE(t, region)
	if err != nil {
		return nil, err
	}

	out, err := client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}

	return out.Table, err
}

// NewDynamoDBClient creates a DynamoDB client.
func NewDynamoDBClient(t testing.TestingT, region string) *dynamodb.Client {
	client, err := NewDynamoDBClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewDynamoDBClientE creates a DynamoDB client.
func NewDynamoDBClientE(t testing.TestingT, region string) (*dynamodb.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(*sess), nil
}
