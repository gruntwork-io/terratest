package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// GetCloudWatchLogEntries returns the CloudWatch log messages in the given region for the given log stream and log group.
func GetCloudWatchLogEntries(t *testing.T, awsRegion string, logStreamName string, logGroupName string, sessExists ...*session.Session) []string {
	out, err := GetCloudWatchLogEntriesE(t, awsRegion, logStreamName, logGroupName, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetCloudWatchLogEntriesE returns the CloudWatch log messages in the given region for the given log stream and log group.
func GetCloudWatchLogEntriesE(t *testing.T, awsRegion string, logStreamName string, logGroupName string, sessExists ...*session.Session) ([]string, error) {
	client, err := NewCloudWatchLogsClientE(t, awsRegion, sessExists[0])
	if err != nil {
		return nil, err
	}

	output, err := client.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})

	if err != nil {
		return nil, err
	}

	entries := []string{}
	for _, event := range output.Events {
		entries = append(entries, *event.Message)
	}

	return entries, nil
}

// NewCloudWatchLogsClient creates a new CloudWatch Logs client.
func NewCloudWatchLogsClient(t *testing.T, region string, sessExists ...*session.Session) *cloudwatchlogs.CloudWatchLogs {
	client, err := NewCloudWatchLogsClientE(t, region, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewCloudWatchLogsClientE creates a new CloudWatch Logs client.
func NewCloudWatchLogsClientE(t *testing.T, region string, sessExists ...*session.Session) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := NewAuthenticatedSession(region, sessExists[0])
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.New(sess), nil
}
