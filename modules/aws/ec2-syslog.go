package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

const syslogRetryInterval = 5 * time.Second

// GetSyslogForInstanceContextE gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
// The ctx parameter supports cancellation and timeouts.
func GetSyslogForInstanceContextE(t testing.TestingT, ctx context.Context, instanceID string, region string) (string, error) {
	description := fmt.Sprintf("Fetching syslog for Instance %s in %s", instanceID, region)
	maxRetries := 120 //nolint:mnd // max retry count for syslog availability

	logger.Default.Logf(t, "%s", description)

	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return "", err
	}

	client := ec2.NewFromConfig(*sess)

	input := ec2.GetConsoleOutputInput{
		InstanceId: aws.String(instanceID),
	}

	syslogB64, err := retry.DoWithRetryContextE(t, ctx, description, maxRetries, syslogRetryInterval, func() (string, error) {
		out, err := client.GetConsoleOutput(ctx, &input)
		if err != nil {
			return "", err
		}

		syslog := aws.ToString(out.Output)
		if syslog == "" {
			return "", fmt.Errorf("syslog is not yet available for instance %s in %s", instanceID, region)
		}

		return syslog, nil
	})
	if err != nil {
		return "", err
	}

	syslogBytes, err := base64.StdEncoding.DecodeString(syslogB64)
	if err != nil {
		return "", err
	}

	return string(syslogBytes), nil
}

// GetSyslogForInstanceContext gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSyslogForInstanceContext(t testing.TestingT, ctx context.Context, instanceID string, region string) string {
	t.Helper()

	out, err := GetSyslogForInstanceContextE(t, ctx, instanceID, region)
	require.NoError(t, err)

	return out
}

// GetSyslogForInstance (Deprecated) See the FetchContentsOfFileFromInstance method for a more powerful solution.
//
// GetSyslogForInstance gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
//
// Deprecated: Use [GetSyslogForInstanceContext] instead.
func GetSyslogForInstance(t testing.TestingT, instanceID string, awsRegion string) string {
	t.Helper()

	return GetSyslogForInstanceContext(t, context.Background(), instanceID, awsRegion)
}

// GetSyslogForInstanceE (Deprecated) See the FetchContentsOfFileFromInstanceE method for a more powerful solution.
//
// GetSyslogForInstanceE gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
//
// Deprecated: Use [GetSyslogForInstanceContextE] instead.
func GetSyslogForInstanceE(t testing.TestingT, instanceID string, region string) (string, error) {
	return GetSyslogForInstanceContextE(t, context.Background(), instanceID, region)
}

// GetSyslogForInstancesInAsgContextE gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance ID -> Syslog for that Instance.
// The ctx parameter supports cancellation and timeouts.
func GetSyslogForInstancesInAsgContextE(t testing.TestingT, ctx context.Context, asgName string, awsRegion string) (map[string]string, error) {
	logger.Default.Logf(t, "Fetching syslog for each Instance in ASG %s in %s", asgName, awsRegion)

	instanceIDs, err := GetEc2InstanceIdsByTagContextE(t, ctx, awsRegion, "aws:autoscaling:groupName", asgName)
	if err != nil {
		return nil, err
	}

	logs := map[string]string{}

	for _, id := range instanceIDs {
		syslog, err := GetSyslogForInstanceContextE(t, ctx, id, awsRegion)
		if err != nil {
			return nil, err
		}

		logs[id] = syslog
	}

	return logs, nil
}

// GetSyslogForInstancesInAsgContext gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance ID -> Syslog for that Instance.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSyslogForInstancesInAsgContext(t testing.TestingT, ctx context.Context, asgName string, awsRegion string) map[string]string {
	t.Helper()

	out, err := GetSyslogForInstancesInAsgContextE(t, ctx, asgName, awsRegion)
	require.NoError(t, err)

	return out
}

// GetSyslogForInstancesInAsg (Deprecated) See the FetchContentsOfFilesFromAsg method for a more powerful solution.
//
// GetSyslogForInstancesInAsg gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance ID -> Syslog for that Instance.
//
// Deprecated: Use [GetSyslogForInstancesInAsgContext] instead.
func GetSyslogForInstancesInAsg(t testing.TestingT, asgName string, awsRegion string) map[string]string {
	t.Helper()

	return GetSyslogForInstancesInAsgContext(t, context.Background(), asgName, awsRegion)
}

// GetSyslogForInstancesInAsgE (Deprecated) See the FetchContentsOfFilesFromAsgE method for a more powerful solution.
//
// GetSyslogForInstancesInAsgE gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance ID -> Syslog for that Instance.
//
// Deprecated: Use [GetSyslogForInstancesInAsgContextE] instead.
func GetSyslogForInstancesInAsgE(t testing.TestingT, asgName string, awsRegion string) (map[string]string, error) {
	return GetSyslogForInstancesInAsgContextE(t, context.Background(), asgName, awsRegion)
}
