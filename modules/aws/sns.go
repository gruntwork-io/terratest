package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// CreateSnsTopicContextE creates an SNS Topic and return the ARN.
// The ctx parameter supports cancellation and timeouts.
func CreateSnsTopicContextE(t testing.TestingT, ctx context.Context, region string, snsTopicName string) (string, error) {
	logger.Default.Logf(t, "Creating SNS topic %s in %s", snsTopicName, region)

	snsClient, err := NewSnsClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	createTopicInput := &sns.CreateTopicInput{
		Name: &snsTopicName,
	}

	output, err := snsClient.CreateTopic(ctx, createTopicInput)
	if err != nil {
		return "", err
	}

	return aws.ToString(output.TopicArn), nil
}

// CreateSnsTopicContext creates an SNS Topic and return the ARN.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateSnsTopicContext(t testing.TestingT, ctx context.Context, region string, snsTopicName string) string {
	t.Helper()

	out, err := CreateSnsTopicContextE(t, ctx, region, snsTopicName)
	require.NoError(t, err)

	return out
}

// CreateSnsTopic creates an SNS Topic and return the ARN.
//
// Deprecated: Use [CreateSnsTopicContext] instead.
func CreateSnsTopic(t testing.TestingT, region string, snsTopicName string) string {
	t.Helper()

	return CreateSnsTopicContext(t, context.Background(), region, snsTopicName)
}

// CreateSnsTopicE creates an SNS Topic and return the ARN.
//
// Deprecated: Use [CreateSnsTopicContextE] instead.
func CreateSnsTopicE(t testing.TestingT, region string, snsTopicName string) (string, error) {
	return CreateSnsTopicContextE(t, context.Background(), region, snsTopicName)
}

// DeleteSNSTopicContextE deletes an SNS Topic.
// The ctx parameter supports cancellation and timeouts.
func DeleteSNSTopicContextE(t testing.TestingT, ctx context.Context, region string, snsTopicArn string) error {
	logger.Default.Logf(t, "Deleting SNS topic %s in %s", snsTopicArn, region)

	snsClient, err := NewSnsClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	deleteTopicInput := &sns.DeleteTopicInput{
		TopicArn: aws.String(snsTopicArn),
	}

	_, err = snsClient.DeleteTopic(ctx, deleteTopicInput)

	return err
}

// DeleteSNSTopicContext deletes an SNS Topic.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteSNSTopicContext(t testing.TestingT, ctx context.Context, region string, snsTopicArn string) {
	t.Helper()

	err := DeleteSNSTopicContextE(t, ctx, region, snsTopicArn)
	require.NoError(t, err)
}

// DeleteSNSTopic deletes an SNS Topic.
//
// Deprecated: Use [DeleteSNSTopicContext] instead.
func DeleteSNSTopic(t testing.TestingT, region string, snsTopicArn string) {
	t.Helper()

	DeleteSNSTopicContext(t, context.Background(), region, snsTopicArn)
}

// DeleteSNSTopicE deletes an SNS Topic.
//
// Deprecated: Use [DeleteSNSTopicContextE] instead.
func DeleteSNSTopicE(t testing.TestingT, region string, snsTopicArn string) error {
	return DeleteSNSTopicContextE(t, context.Background(), region, snsTopicArn)
}

// NewSnsClientContextE creates a new SNS client.
// The ctx parameter supports cancellation and timeouts.
func NewSnsClientContextE(t testing.TestingT, ctx context.Context, region string) (*sns.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return sns.NewFromConfig(*sess), nil
}

// NewSnsClientContext creates a new SNS client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewSnsClientContext(t testing.TestingT, ctx context.Context, region string) *sns.Client {
	t.Helper()

	client, err := NewSnsClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewSnsClient creates a new SNS client.
//
// Deprecated: Use [NewSnsClientContext] instead.
func NewSnsClient(t testing.TestingT, region string) *sns.Client {
	t.Helper()

	return NewSnsClientContext(t, context.Background(), region)
}

// NewSnsClientE creates a new SNS client.
//
// Deprecated: Use [NewSnsClientContextE] instead.
func NewSnsClientE(t testing.TestingT, region string) (*sns.Client, error) {
	return NewSnsClientContextE(t, context.Background(), region)
}
