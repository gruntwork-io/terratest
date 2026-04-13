package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// AsgCapacityInfo holds capacity information about an Auto Scaling Group.
type AsgCapacityInfo struct {
	MinCapacity     int64
	MaxCapacity     int64
	CurrentCapacity int64
	DesiredCapacity int64
}

// GetCapacityInfoForAsgContextE returns the capacity info for the queried asg as a struct, AsgCapacityInfo.
// The ctx parameter supports cancellation and timeouts.
func GetCapacityInfoForAsgContextE(t testing.TestingT, ctx context.Context, asgName string, awsRegion string) (AsgCapacityInfo, error) {
	asgClient, err := NewAsgClientContextE(t, ctx, awsRegion)
	if err != nil {
		return AsgCapacityInfo{}, err
	}

	input := autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []string{asgName}}

	output, err := asgClient.DescribeAutoScalingGroups(ctx, &input)
	if err != nil {
		return AsgCapacityInfo{}, err
	}

	groups := output.AutoScalingGroups
	if len(groups) == 0 {
		return AsgCapacityInfo{}, NewNotFoundError("ASG", asgName, awsRegion)
	}

	capacityInfo := AsgCapacityInfo{
		MinCapacity:     int64(*groups[0].MinSize),
		MaxCapacity:     int64(*groups[0].MaxSize),
		DesiredCapacity: int64(*groups[0].DesiredCapacity),
		CurrentCapacity: int64(len(groups[0].Instances)),
	}

	return capacityInfo, nil
}

// GetCapacityInfoForAsgContext returns the capacity info for the queried asg as a struct, AsgCapacityInfo.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCapacityInfoForAsgContext(t testing.TestingT, ctx context.Context, asgName string, awsRegion string) AsgCapacityInfo {
	t.Helper()

	capacityInfo, err := GetCapacityInfoForAsgContextE(t, ctx, asgName, awsRegion)
	require.NoError(t, err)

	return capacityInfo
}

// GetCapacityInfoForAsg returns the capacity info for the queried asg as a struct, AsgCapacityInfo.
//
// Deprecated: Use [GetCapacityInfoForAsgContext] instead.
func GetCapacityInfoForAsg(t testing.TestingT, asgName string, awsRegion string) AsgCapacityInfo {
	t.Helper()

	return GetCapacityInfoForAsgContext(t, context.Background(), asgName, awsRegion)
}

// GetCapacityInfoForAsgE returns the capacity info for the queried asg as a struct, AsgCapacityInfo.
//
// Deprecated: Use [GetCapacityInfoForAsgContextE] instead.
func GetCapacityInfoForAsgE(t testing.TestingT, asgName string, awsRegion string) (AsgCapacityInfo, error) {
	return GetCapacityInfoForAsgContextE(t, context.Background(), asgName, awsRegion)
}

// GetInstanceIdsForAsgContextE gets the IDs of EC2 Instances in the given ASG.
// The ctx parameter supports cancellation and timeouts.
func GetInstanceIdsForAsgContextE(t testing.TestingT, ctx context.Context, asgName string, awsRegion string) ([]string, error) {
	asgClient, err := NewAsgClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	input := autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []string{asgName}}

	output, err := asgClient.DescribeAutoScalingGroups(ctx, &input)
	if err != nil {
		return nil, err
	}

	var instanceIDs []string

	for i := range output.AutoScalingGroups {
		for _, instance := range output.AutoScalingGroups[i].Instances {
			instanceIDs = append(instanceIDs, aws.ToString(instance.InstanceId))
		}
	}

	return instanceIDs, nil
}

// GetInstanceIdsForAsgContext gets the IDs of EC2 Instances in the given ASG.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetInstanceIdsForAsgContext(t testing.TestingT, ctx context.Context, asgName string, awsRegion string) []string {
	t.Helper()

	ids, err := GetInstanceIdsForAsgContextE(t, ctx, asgName, awsRegion)
	require.NoError(t, err)

	return ids
}

// GetInstanceIdsForAsg gets the IDs of EC2 Instances in the given ASG.
//
// Deprecated: Use [GetInstanceIdsForAsgContext] instead.
func GetInstanceIdsForAsg(t testing.TestingT, asgName string, awsRegion string) []string {
	t.Helper()

	return GetInstanceIdsForAsgContext(t, context.Background(), asgName, awsRegion)
}

// GetInstanceIdsForAsgE gets the IDs of EC2 Instances in the given ASG.
//
// Deprecated: Use [GetInstanceIdsForAsgContextE] instead.
func GetInstanceIdsForAsgE(t testing.TestingT, asgName string, awsRegion string) ([]string, error) {
	return GetInstanceIdsForAsgContextE(t, context.Background(), asgName, awsRegion)
}

// WaitForCapacityContextE waits for the currently set desired capacity to be reached on the ASG.
// The ctx parameter supports cancellation and timeouts.
func WaitForCapacityContextE(
	t testing.TestingT,
	ctx context.Context,
	asgName string,
	region string,
	maxRetries int,
	sleepBetweenRetries time.Duration,
) error {
	msg, err := retry.DoWithRetryContextE(
		t,
		ctx,
		fmt.Sprintf("Waiting for ASG %s to reach desired capacity.", asgName),
		maxRetries,
		sleepBetweenRetries,
		func() (string, error) {
			capacityInfo, err := GetCapacityInfoForAsgContextE(t, ctx, asgName, region)
			if err != nil {
				return "", err
			}

			if capacityInfo.CurrentCapacity != capacityInfo.DesiredCapacity {
				return "", NewAsgCapacityNotMetError(asgName, capacityInfo.DesiredCapacity, capacityInfo.CurrentCapacity)
			}

			return fmt.Sprintf("ASG %s is now at desired capacity %d", asgName, capacityInfo.DesiredCapacity), nil
		},
	)
	logger.Default.Logf(t, "%s", msg)

	return err
}

// WaitForCapacityContext waits for the currently set desired capacity to be reached on the ASG.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func WaitForCapacityContext(
	t testing.TestingT,
	ctx context.Context,
	asgName string,
	region string,
	maxRetries int,
	sleepBetweenRetries time.Duration,
) {
	t.Helper()

	err := WaitForCapacityContextE(t, ctx, asgName, region, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
}

// WaitForCapacity waits for the currently set desired capacity to be reached on the ASG
//
// Deprecated: Use [WaitForCapacityContext] instead.
func WaitForCapacity(
	t testing.TestingT,
	asgName string,
	region string,
	maxRetries int,
	sleepBetweenRetries time.Duration,
) {
	t.Helper()

	WaitForCapacityContext(t, context.Background(), asgName, region, maxRetries, sleepBetweenRetries)
}

// WaitForCapacityE waits for the currently set desired capacity to be reached on the ASG
//
// Deprecated: Use [WaitForCapacityContextE] instead.
func WaitForCapacityE(
	t testing.TestingT,
	asgName string,
	region string,
	maxRetries int,
	sleepBetweenRetries time.Duration,
) error {
	return WaitForCapacityContextE(t, context.Background(), asgName, region, maxRetries, sleepBetweenRetries)
}

// NewAsgClientContextE creates an Auto Scaling Group client.
// The ctx parameter supports cancellation and timeouts.
func NewAsgClientContextE(t testing.TestingT, ctx context.Context, region string) (*autoscaling.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return autoscaling.NewFromConfig(*sess), nil
}

// NewAsgClientContext creates an Auto Scaling Group client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewAsgClientContext(t testing.TestingT, ctx context.Context, region string) *autoscaling.Client {
	t.Helper()

	client, err := NewAsgClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewAsgClient creates an Auto Scaling Group client.
//
// Deprecated: Use [NewAsgClientContext] instead.
func NewAsgClient(t testing.TestingT, region string) *autoscaling.Client {
	t.Helper()

	return NewAsgClientContext(t, context.Background(), region)
}

// NewAsgClientE creates an Auto Scaling Group client.
//
// Deprecated: Use [NewAsgClientContextE] instead.
func NewAsgClientE(t testing.TestingT, region string) (*autoscaling.Client, error) {
	return NewAsgClientContextE(t, context.Background(), region)
}
