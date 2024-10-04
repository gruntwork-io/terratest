package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetCloudWatchLogEntries returns the CloudWatch log messages in the given region for the given log stream and log group.
func GetCloudWatchLogEntries(t testing.TestingT, awsRegion string, logStreamName string, logGroupName string) []string {
	out, err := GetCloudWatchLogEntriesE(t, awsRegion, logStreamName, logGroupName)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetCloudWatchLogEntriesE returns the CloudWatch log messages in the given region for the given log stream and log group.
func GetCloudWatchLogEntriesE(t testing.TestingT, awsRegion string, logStreamName string, logGroupName string) ([]string, error) {
	client, err := NewCloudWatchLogsClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	output, err := client.GetLogEvents(context.Background(), &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})

	if err != nil {
		return nil, err
	}

	var entries []string
	for _, event := range output.Events {
		entries = append(entries, *event.Message)
	}

	return entries, nil
}

// NewCloudWatchLogsClient creates a new CloudWatch Logs client.
func NewCloudWatchLogsClient(t testing.TestingT, region string) *cloudwatchlogs.Client {
	client, err := NewCloudWatchLogsClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewCloudWatchLogsClientE creates a new CloudWatch Logs client.
func NewCloudWatchLogsClientE(t testing.TestingT, region string) (*cloudwatchlogs.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.NewFromConfig(*sess), nil
}
