package aws

import (
	"fmt"
)

// IpForEc2InstanceNotFound is an error that occurs when the IP for an EC2 instance is not found.
//
// Deprecated: Use [IPForEc2InstanceNotFound] instead.
type IpForEc2InstanceNotFound struct { //nolint:staticcheck,revive // preserving deprecated type name
	InstanceId string //nolint:staticcheck,revive // preserving existing field name
	AwsRegion  string
	Type       string
}

// IPForEc2InstanceNotFound is an alias for [IpForEc2InstanceNotFound].
type IPForEc2InstanceNotFound = IpForEc2InstanceNotFound //nolint:staticcheck,revive // preserving deprecated type name

func (err IpForEc2InstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a %s IP address for EC2 Instance %s in %s", err.Type, err.InstanceId, err.AwsRegion)
}

// HostnameForEc2InstanceNotFound is an error that occurs when the IP for an EC2 instance is not found.
type HostnameForEc2InstanceNotFound struct {
	InstanceId string //nolint:staticcheck,revive // preserving existing field name
	AwsRegion  string
	Type       string
}

func (err HostnameForEc2InstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a %s hostname for EC2 Instance %s in %s", err.Type, err.InstanceId, err.AwsRegion)
}

// NotFoundError is returned when an expected object is not found
type NotFoundError struct {
	objectType string
	objectID   string
	region     string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("Object of type %s with id %s not found in region %s", err.objectType, err.objectID, err.region)
}

// NewNotFoundError returns a [NotFoundError] for the given object type, ID, and region.
func NewNotFoundError(objectType string, objectID string, region string) NotFoundError {
	return NotFoundError{objectType, objectID, region}
}

// AsgCapacityNotMetError is returned when the ASG capacity is not yet at the desired capacity.
type AsgCapacityNotMetError struct {
	asgName         string
	desiredCapacity int64
	currentCapacity int64
}

func (err AsgCapacityNotMetError) Error() string {
	return fmt.Sprintf(
		"ASG %s not yet at desired capacity %d (current %d)",
		err.asgName,
		err.desiredCapacity,
		err.currentCapacity,
	)
}

// NewAsgCapacityNotMetError returns an [AsgCapacityNotMetError] describing
// the given ASG's desired and current capacities.
func NewAsgCapacityNotMetError(asgName string, desiredCapacity int64, currentCapacity int64) AsgCapacityNotMetError {
	return AsgCapacityNotMetError{asgName, desiredCapacity, currentCapacity}
}

// BucketVersioningNotEnabledError is returned when an S3 bucket that should have versioning does not have it applied
type BucketVersioningNotEnabledError struct {
	s3BucketName     string
	awsRegion        string
	versioningStatus string
}

func (err BucketVersioningNotEnabledError) Error() string {
	return fmt.Sprintf(
		"Versioning status for bucket %s in the %s region is %s",
		err.s3BucketName,
		err.awsRegion,
		err.versioningStatus,
	)
}

// NewBucketVersioningNotEnabledError returns a [BucketVersioningNotEnabledError]
// for the given S3 bucket, region, and observed versioning status.
func NewBucketVersioningNotEnabledError(s3BucketName string, awsRegion string, versioningStatus string) BucketVersioningNotEnabledError {
	return BucketVersioningNotEnabledError{s3BucketName: s3BucketName, awsRegion: awsRegion, versioningStatus: versioningStatus}
}

// NoBucketPolicyError is returned when an S3 bucket that should have a policy applied does not
type NoBucketPolicyError struct {
	s3BucketName string
	awsRegion    string
	bucketPolicy string
}

func (err NoBucketPolicyError) Error() string {
	return fmt.Sprintf(
		"The policy for bucket %s in the %s region does not have a policy attached.",
		err.s3BucketName,
		err.awsRegion,
	)
}

// NewNoBucketPolicyError returns a [NoBucketPolicyError] for the given S3
// bucket, region, and bucket policy.
func NewNoBucketPolicyError(s3BucketName string, awsRegion string, bucketPolicy string) NoBucketPolicyError {
	return NoBucketPolicyError{s3BucketName: s3BucketName, awsRegion: awsRegion, bucketPolicy: bucketPolicy}
}

// NoInstanceTypeError is returned when none of the given instance type options are available in all AZs in a region
type NoInstanceTypeError struct {
	InstanceTypeOptions []string
	Azs                 []string
}

func (err NoInstanceTypeError) Error() string {
	return fmt.Sprintf(
		"None of the given instance types (%v) is available in all the AZs in this region (%v).",
		err.InstanceTypeOptions,
		err.Azs,
	)
}

// NoRdsInstanceTypeError is returned when none of the given instance types are available for the region, database engine, and database engine combination given
type NoRdsInstanceTypeError struct {
	DatabaseEngine        string
	DatabaseEngineVersion string
	InstanceTypeOptions   []string
}

func (err NoRdsInstanceTypeError) Error() string {
	return fmt.Sprintf(
		"None of the given RDS instance types (%v) is available in this region for database engine (%v) of version (%v).",
		err.InstanceTypeOptions,
		err.DatabaseEngine,
		err.DatabaseEngineVersion,
	)
}
