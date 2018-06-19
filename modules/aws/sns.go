package aws

import (
	"testing"

	"github.com/Briansbum/terratest/modules/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// CreateSnsTopic creates an SNS Topic and return the ARN.
func CreateSnsTopic(t *testing.T, region string, snsTopicName string, sessExists ...*session.Session) string {
	out, err := CreateSnsTopicE(t, region, snsTopicName, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// CreateSnsTopicE creates an SNS Topic and return the ARN.
func CreateSnsTopicE(t *testing.T, region string, snsTopicName string, sessExists ...*session.Session) (string, error) {
	logger.Logf(t, "Creating SNS topic %s in %s", snsTopicName, region)

	snsClient, err := NewSnsClientE(t, region, sessExists[0])
	if err != nil {
		return "", err
	}

	createTopicInput := &sns.CreateTopicInput{
		Name: &snsTopicName,
	}

	output, err := snsClient.CreateTopic(createTopicInput)
	if err != nil {
		return "", err
	}

	return aws.StringValue(output.TopicArn), err
}

// DeleteSNSTopic deletes an SNS Topic.
func DeleteSNSTopic(t *testing.T, region string, snsTopicArn string, sessExists ...*session.Session) {
	err := DeleteSNSTopicE(t, region, snsTopicArn, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteSNSTopicE deletes an SNS Topic.
func DeleteSNSTopicE(t *testing.T, region string, snsTopicArn string, sessExists ...*session.Session) error {
	logger.Logf(t, "Deleting SNS topic %s in %s", snsTopicArn, region)

	snsClient, err := NewSnsClientE(t, region, sessExists[0])
	if err != nil {
		return err
	}

	deleteTopicInput := &sns.DeleteTopicInput{
		TopicArn: aws.String(snsTopicArn),
	}

	_, err = snsClient.DeleteTopic(deleteTopicInput)
	return err
}

// NewSnsClient creates a new SNS client.
func NewSnsClient(t *testing.T, region string, sessExists ...*session.Session) *sns.SNS {
	client, err := NewSnsClientE(t, region, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewSnsClientE creates a new SNS client.
func NewSnsClientE(t *testing.T, region string, sessExists ...*session.Session) (*sns.SNS, error) {
	sess, err := NewAuthenticatedSession(region, sessExists[0])
	if err != nil {
		return nil, err
	}

	return sns.New(sess), nil
}
