package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
