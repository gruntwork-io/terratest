package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// CreateSnsTopic creates an SNS Topic and return the ARN.
func CreateSnsTopic(t testing.TestingT, region string, snsTopicName string) string {
	out, err := CreateSnsTopicE(t, region, snsTopicName)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// CreateSnsTopicE creates an SNS Topic and return the ARN.
func CreateSnsTopicE(t testing.TestingT, region string, snsTopicName string) (string, error) {
	logger.Default.Logf(t, "Creating SNS topic %s in %s", snsTopicName, region)

	snsClient, err := NewSnsClientE(t, region)
	if err != nil {
		return "", err
	}

	createTopicInput := &sns.CreateTopicInput{
		Name: &snsTopicName,
	}

	output, err := snsClient.CreateTopic(context.Background(), createTopicInput)
	if err != nil {
		return "", err
	}

	return aws.ToString(output.TopicArn), err
}

// DeleteSNSTopic deletes an SNS Topic.
func DeleteSNSTopic(t testing.TestingT, region string, snsTopicArn string) {
	err := DeleteSNSTopicE(t, region, snsTopicArn)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteSNSTopicE deletes an SNS Topic.
func DeleteSNSTopicE(t testing.TestingT, region string, snsTopicArn string) error {
	logger.Default.Logf(t, "Deleting SNS topic %s in %s", snsTopicArn, region)

	snsClient, err := NewSnsClientE(t, region)
	if err != nil {
		return err
	}

	deleteTopicInput := &sns.DeleteTopicInput{
		TopicArn: aws.String(snsTopicArn),
	}

	_, err = snsClient.DeleteTopic(context.Background(), deleteTopicInput)
	return err
}

// NewSnsClient creates a new SNS client.
func NewSnsClient(t testing.TestingT, region string) *sns.Client {
	client, err := NewSnsClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewSnsClientE creates a new SNS client.
func NewSnsClientE(t testing.TestingT, region string) (*sns.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return sns.NewFromConfig(*sess), nil
}
