package aws

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetSyslogForInstance (Deprecated) See the FetchContentsOfFileFromInstance method for a more powerful solution.
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

// GetSyslogForInstanceE (Deprecated) See the FetchContentsOfFileFromInstanceE method for a more powerful solution.
//
// GetSyslogForInstanceE gets the syslog for the Instance with the given ID in the given region. This should be available ~1 minute after an
// Instance boots and is very useful for debugging boot-time issues, such as an error in User Data.
func GetSyslogForInstanceE(t testing.TestingT, instanceID string, region string) (string, error) {
	return "", fmt.Errorf("(Deprecated) use FetchContentsOfFileFromInstanceE method instead")
}

// GetSyslogForInstancesInAsg (Deprecated) See the FetchContentsOfFilesFromAsg method for a more powerful solution.
//
// GetSyslogForInstancesInAsg gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance ID -> Syslog for that Instance.
func GetSyslogForInstancesInAsg(t testing.TestingT, asgName string, awsRegion string) map[string]string {
	out, err := GetSyslogForInstancesInAsgE(t, asgName, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetSyslogForInstancesInAsgE (Deprecated) See the FetchContentsOfFilesFromAsgE method for a more powerful solution.
//
// GetSyslogForInstancesInAsgE gets the syslog for each of the Instances in the given ASG in the given region. These logs should be available ~1
// minute after the Instance boots and are very useful for debugging boot-time issues, such as an error in User Data.
// Returns a map of Instance ID -> Syslog for that Instance.
func GetSyslogForInstancesInAsgE(t testing.TestingT, asgName string, awsRegion string) (map[string]string, error) {
	return nil, fmt.Errorf("(Deprecated) use FetchContentsOfFilesFromAsgE method instead")
}
