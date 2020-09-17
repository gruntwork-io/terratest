package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
)

type AsgCapacityInfo struct {
	MinCapacity     int64
	MaxCapacity     int64
	CurrentCapacity int64
	DesiredCapacity int64
}

// GetCapacityInfoForAsg returns the capacity info for the queried asg as a struct, AsgCapacityInfo.
func GetCapacityInfoForAsg(t testing.TestingT, asgName string, awsRegion string) AsgCapacityInfo {
	capacityInfo, err := GetCapacityInfoForAsgE(t, asgName, awsRegion)
	require.NoError(t, err)
	return capacityInfo
}

// GetCapacityInfoForAsgE returns the capacity info for the queried asg as a struct, AsgCapacityInfo.
func GetCapacityInfoForAsgE(t testing.TestingT, asgName string, awsRegion string) (AsgCapacityInfo, error) {
	asgClient, err := NewAsgClientE(t, awsRegion)
	if err != nil {
		return AsgCapacityInfo{}, err
	}

	input := autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []*string{aws.String(asgName)}}
	output, err := asgClient.DescribeAutoScalingGroups(&input)
	if err != nil {
		return AsgCapacityInfo{}, err
	}
	groups := output.AutoScalingGroups
	if len(groups) == 0 {
		return AsgCapacityInfo{}, NewNotFoundError("ASG", asgName, awsRegion)
	}
	capacityInfo := AsgCapacityInfo{
		MinCapacity:     *groups[0].MinSize,
		MaxCapacity:     *groups[0].MaxSize,
		DesiredCapacity: *groups[0].DesiredCapacity,
		CurrentCapacity: int64(len(groups[0].Instances)),
	}
	return capacityInfo, nil
}

// GetInstanceIdsForAsg gets the IDs of EC2 Instances in the given ASG.
func GetInstanceIdsForAsg(t testing.TestingT, asgName string, awsRegion string) []string {
	ids, err := GetInstanceIdsForAsgE(t, asgName, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsForAsgE gets the IDs of EC2 Instances in the given ASG.
func GetInstanceIdsForAsgE(t testing.TestingT, asgName string, awsRegion string) ([]string, error) {
	asgClient, err := NewAsgClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	input := autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []*string{aws.String(asgName)}}
	output, err := asgClient.DescribeAutoScalingGroups(&input)
	if err != nil {
		return nil, err
	}

	instanceIDs := []string{}
	for _, asg := range output.AutoScalingGroups {
		for _, instance := range asg.Instances {
			instanceIDs = append(instanceIDs, aws.StringValue(instance.InstanceId))
		}
	}

	return instanceIDs, nil
}

// WaitForCapacity waits for the currently set desired capacity to be reached on the ASG
func WaitForCapacity(
	t testing.TestingT,
	asgName string,
	region string,
	maxRetries int,
	sleepBetweenRetries time.Duration,
) {
	err := WaitForCapacityE(t, asgName, region, maxRetries, sleepBetweenRetries)
	require.NoError(t, err)
}

// WaitForCapacityE waits for the currently set desired capacity to be reached on the ASG
func WaitForCapacityE(
	t testing.TestingT,
	asgName string,
	region string,
	maxRetries int,
	sleepBetweenRetries time.Duration,
) error {
	msg, err := retry.DoWithRetryE(
		t,
		fmt.Sprintf("Waiting for ASG %s to reach desired capacity.", asgName),
		maxRetries,
		sleepBetweenRetries,
		func() (string, error) {
			capacityInfo, err := GetCapacityInfoForAsgE(t, asgName, region)
			if err != nil {
				return "", err
			}
			if capacityInfo.CurrentCapacity != capacityInfo.DesiredCapacity {
				return "", NewAsgCapacityNotMetError(asgName, capacityInfo.DesiredCapacity, capacityInfo.CurrentCapacity)
			}
			return fmt.Sprintf("ASG %s is now at desired capacity %d", asgName, capacityInfo.DesiredCapacity), nil
		},
	)
	logger.Log(t, msg)
	return err
}

// NewAsgClient creates an Auto Scaling Group client.
func NewAsgClient(t testing.TestingT, region string) *autoscaling.AutoScaling {
	client, err := NewAsgClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewAsgClientE creates an Auto Scaling Group client.
func NewAsgClientE(t testing.TestingT, region string) (*autoscaling.AutoScaling, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return autoscaling.New(sess), nil
}
