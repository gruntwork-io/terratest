package aws

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// (Deprecated) See the FetchContentsOfFileFromInstance method for a more powerful solution.
//
// GetSyslogForInstance gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
func GetSyslogForInstance(t testing.TestingT, instanceID string, awsRegion string) string {
	out, err := GetSyslogForInstanceE(t, instanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// (Deprecated) See the FetchContentsOfFileFromInstanceE method for a more powerful solution.
//
// GetSyslogForInstanceE gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
func GetSyslogForInstanceE(t testing.TestingT, instanceID string, region string) (string, error) {
	description := fmt.Sprintf("Fetching syslog for Instance %s in %s", instanceID, region)
	maxRetries := 120
	timeBetweenRetries := 5 * time.Second

	logger.Log(t, description)

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return "", err
	}

	input := ec2.GetConsoleOutputInput{
		InstanceId: aws.String(instanceID),
	}

	syslogB64, err := retry.DoWithRetryE(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		out, err := client.GetConsoleOutput(&input)
		if err != nil {
			return "", err
		}

		syslog := aws.StringValue(out.Output)
		if syslog == "" {
			return "", fmt.Errorf("Syslog is not yet available for instance %s in %s", instanceID, region)
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

// (Deprecated) See the FetchContentsOfFilesFromAsg method for a more powerful solution.
//
// GetSyslogForInstancesInAsg gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance Id -> Syslog for that Instance.
func GetSyslogForInstancesInAsg(t testing.TestingT, asgName string, awsRegion string) map[string]string {
	out, err := GetSyslogForInstancesInAsgE(t, asgName, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// (Deprecated) See the FetchContentsOfFilesFromAsgE method for a more powerful solution.
//
// GetSyslogForInstancesInAsgE gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance Id -> Syslog for that Instance.
func GetSyslogForInstancesInAsgE(t testing.TestingT, asgName string, awsRegion string) (map[string]string, error) {
	logger.Logf(t, "Fetching syslog for each Instance in ASG %s in %s", asgName, awsRegion)

	instanceIDs, err := GetEc2InstanceIdsByTagE(t, awsRegion, "aws:autoscaling:groupName", asgName)
	if err != nil {
		return nil, err
	}

	logs := map[string]string{}
	for _, id := range instanceIDs {
		syslog, err := GetSyslogForInstanceE(t, id, awsRegion)
		if err != nil {
			return nil, err
		}
		logs[id] = syslog
	}

	return logs, nil
}
