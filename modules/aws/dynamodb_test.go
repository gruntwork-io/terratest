package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
)

// mockDynamoDBClient is a test double for DynamoDBAPI.
type mockDynamoDBClient struct {
	DescribeTableOutput      *dynamodb.DescribeTableOutput
	DescribeTableErr         error
	DescribeTimeToLiveOutput *dynamodb.DescribeTimeToLiveOutput
	DescribeTimeToLiveErr    error
	ListTagsOfResourceOutput *dynamodb.ListTagsOfResourceOutput
	ListTagsOfResourceErr    error
	lastDescribeTableName    string
	lastListTagsResourceArn  string
}

func (m *mockDynamoDBClient) DescribeTable(_ context.Context, params *dynamodb.DescribeTableInput, _ ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error) {
	m.lastDescribeTableName = aws.ToString(params.TableName)
	if m.DescribeTableErr != nil {
		return nil, m.DescribeTableErr
	}
	return m.DescribeTableOutput, nil
}

func (m *mockDynamoDBClient) DescribeTimeToLive(_ context.Context, _ *dynamodb.DescribeTimeToLiveInput, _ ...func(*dynamodb.Options)) (*dynamodb.DescribeTimeToLiveOutput, error) {
	if m.DescribeTimeToLiveErr != nil {
		return nil, m.DescribeTimeToLiveErr
	}
	return m.DescribeTimeToLiveOutput, nil
}

func (m *mockDynamoDBClient) ListTagsOfResource(_ context.Context, params *dynamodb.ListTagsOfResourceInput, _ ...func(*dynamodb.Options)) (*dynamodb.ListTagsOfResourceOutput, error) {
	m.lastListTagsResourceArn = aws.ToString(params.ResourceArn)
	if m.ListTagsOfResourceErr != nil {
		return nil, m.ListTagsOfResourceErr
	}
	return m.ListTagsOfResourceOutput, nil
}

const (
	testTableArn  = "arn:aws:dynamodb:us-east-1:123456789012:table/my-table"
	testTableName = "my-table"
)

func TestGetDynamoDBTableWithClientContextE(t *testing.T) {
	t.Parallel()

	t.Run("returns table description", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{
			DescribeTableOutput: &dynamodb.DescribeTableOutput{
				Table: &types.TableDescription{
					TableArn:  aws.String(testTableArn),
					TableName: aws.String(testTableName),
				},
			},
		}
		got, err := GetDynamoDBTableWithClientContextE(t, context.Background(), client, testTableName)
		require.NoError(t, err)
		require.Equal(t, testTableArn, aws.ToString(got.TableArn))
		require.Equal(t, testTableName, client.lastDescribeTableName)
	})

	t.Run("propagates api error", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{DescribeTableErr: errors.New("ResourceNotFoundException")}
		_, err := GetDynamoDBTableWithClientContextE(t, context.Background(), client, testTableName)
		require.Error(t, err)
	})
}

func TestGetDynamoDBTableTimeToLiveWithClientContextE(t *testing.T) {
	t.Parallel()

	t.Run("returns ttl description", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{
			DescribeTimeToLiveOutput: &dynamodb.DescribeTimeToLiveOutput{
				TimeToLiveDescription: &types.TimeToLiveDescription{
					TimeToLiveStatus: types.TimeToLiveStatusEnabled,
					AttributeName:    aws.String("expiresAt"),
				},
			},
		}
		got, err := GetDynamoDBTableTimeToLiveWithClientContextE(t, context.Background(), client, testTableName)
		require.NoError(t, err)
		require.Equal(t, types.TimeToLiveStatusEnabled, got.TimeToLiveStatus)
		require.Equal(t, "expiresAt", aws.ToString(got.AttributeName))
	})

	t.Run("propagates api error", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{DescribeTimeToLiveErr: errors.New("InternalServerError")}
		_, err := GetDynamoDBTableTimeToLiveWithClientContextE(t, context.Background(), client, testTableName)
		require.Error(t, err)
	})
}

func TestGetDynamoDBTableTagsWithClientContextE(t *testing.T) {
	t.Parallel()

	describeOK := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			TableArn:  aws.String(testTableArn),
			TableName: aws.String(testTableName),
		},
	}

	t.Run("returns tags and queries by the described table arn", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{
			DescribeTableOutput: describeOK,
			ListTagsOfResourceOutput: &dynamodb.ListTagsOfResourceOutput{
				Tags: []types.Tag{
					{Key: aws.String("env"), Value: aws.String("prod")},
					{Key: aws.String("team"), Value: aws.String("platform")},
				},
			},
		}
		got, err := GetDynamoDBTableTagsWithClientContextE(t, context.Background(), client, testTableName)
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.Equal(t, testTableArn, client.lastListTagsResourceArn)
	})

	t.Run("returns empty slice when table has no tags", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{
			DescribeTableOutput:      describeOK,
			ListTagsOfResourceOutput: &dynamodb.ListTagsOfResourceOutput{},
		}
		got, err := GetDynamoDBTableTagsWithClientContextE(t, context.Background(), client, testTableName)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("propagates describe table error without calling list tags", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{DescribeTableErr: errors.New("ResourceNotFoundException")}
		_, err := GetDynamoDBTableTagsWithClientContextE(t, context.Background(), client, testTableName)
		require.Error(t, err)
		require.Empty(t, client.lastListTagsResourceArn, "ListTagsOfResource must not be called when describe fails")
	})

	t.Run("propagates list tags error", func(t *testing.T) {
		t.Parallel()

		client := &mockDynamoDBClient{
			DescribeTableOutput:   describeOK,
			ListTagsOfResourceErr: errors.New("AccessDeniedException"),
		}
		_, err := GetDynamoDBTableTagsWithClientContextE(t, context.Background(), client, testTableName)
		require.Error(t, err)
	})
}
