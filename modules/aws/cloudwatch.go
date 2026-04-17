package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// CloudWatchLogsAPI is the subset of *cloudwatchlogs.Client operations used by the helpers in this
// file. Declared as an interface so tests can substitute a mock; a real *cloudwatchlogs.Client
// satisfies it automatically.
type CloudWatchLogsAPI interface {
	GetLogEvents(ctx context.Context, params *cloudwatchlogs.GetLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error)
}

// GetCloudWatchLogEntriesContextE returns the CloudWatch log messages in the given region for the given log stream and log group.
// The ctx parameter supports cancellation and timeouts.
func GetCloudWatchLogEntriesContextE(t testing.TestingT, ctx context.Context, awsRegion string, logStreamName string, logGroupName string) ([]string, error) {
	client, err := NewCloudWatchLogsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	return GetCloudWatchLogEntriesWithClientContextE(t, ctx, client, logStreamName, logGroupName)
}

// GetCloudWatchLogEntriesWithClientContextE returns the CloudWatch log messages for the given log
// stream and log group using the provided CloudWatch Logs client.
// The ctx parameter supports cancellation and timeouts.
func GetCloudWatchLogEntriesWithClientContextE(t testing.TestingT, ctx context.Context, client CloudWatchLogsAPI, logStreamName string, logGroupName string) ([]string, error) {
	output, err := client.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
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

// GetCloudWatchLogEntriesContext returns the CloudWatch log messages in the given region for the given log stream and log group.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudWatchLogEntriesContext(t testing.TestingT, ctx context.Context, awsRegion string, logStreamName string, logGroupName string) []string {
	t.Helper()

	out, err := GetCloudWatchLogEntriesContextE(t, ctx, awsRegion, logStreamName, logGroupName)
	require.NoError(t, err)

	return out
}

// GetCloudWatchLogEntries returns the CloudWatch log messages in the given region for the given log stream and log group.
//
// Deprecated: Use [GetCloudWatchLogEntriesContext] instead.
func GetCloudWatchLogEntries(t testing.TestingT, awsRegion string, logStreamName string, logGroupName string) []string {
	t.Helper()

	return GetCloudWatchLogEntriesContext(t, context.Background(), awsRegion, logStreamName, logGroupName)
}

// GetCloudWatchLogEntriesE returns the CloudWatch log messages in the given region for the given log stream and log group.
//
// Deprecated: Use [GetCloudWatchLogEntriesContextE] instead.
func GetCloudWatchLogEntriesE(t testing.TestingT, awsRegion string, logStreamName string, logGroupName string) ([]string, error) {
	return GetCloudWatchLogEntriesContextE(t, context.Background(), awsRegion, logStreamName, logGroupName)
}

// NewCloudWatchLogsClientContextE creates a new CloudWatch Logs client.
// The ctx parameter supports cancellation and timeouts.
func NewCloudWatchLogsClientContextE(t testing.TestingT, ctx context.Context, region string) (*cloudwatchlogs.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return cloudwatchlogs.NewFromConfig(*sess), nil
}

// NewCloudWatchLogsClientContext creates a new CloudWatch Logs client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewCloudWatchLogsClientContext(t testing.TestingT, ctx context.Context, region string) *cloudwatchlogs.Client {
	t.Helper()

	client, err := NewCloudWatchLogsClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewCloudWatchLogsClient creates a new CloudWatch Logs client.
//
// Deprecated: Use [NewCloudWatchLogsClientContext] instead.
func NewCloudWatchLogsClient(t testing.TestingT, region string) *cloudwatchlogs.Client {
	t.Helper()

	return NewCloudWatchLogsClientContext(t, context.Background(), region)
}

// NewCloudWatchLogsClientE creates a new CloudWatch Logs client.
//
// Deprecated: Use [NewCloudWatchLogsClientContextE] instead.
func NewCloudWatchLogsClientE(t testing.TestingT, region string) (*cloudwatchlogs.Client, error) {
	return NewCloudWatchLogsClientContextE(t, context.Background(), region)
}
