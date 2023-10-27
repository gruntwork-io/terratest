package aws

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetCapacityInfoForAsg(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueId()
	asgName := fmt.Sprintf("%s-%s", t.Name(), uniqueID)
	region := GetRandomStableRegion(t, []string{}, []string{})

	defer deleteAutoScalingGroup(t, asgName, region)
	createTestAutoScalingGroup(t, asgName, region, 2)
	WaitForCapacity(t, asgName, region, 40, 15*time.Second)

	capacityInfo := GetCapacityInfoForAsg(t, asgName, region)
	assert.Equal(t, capacityInfo.DesiredCapacity, int64(2))
	assert.Equal(t, capacityInfo.CurrentCapacity, int64(2))
	assert.Equal(t, capacityInfo.MinCapacity, int64(1))
	assert.Equal(t, capacityInfo.MaxCapacity, int64(3))
}

func TestGetInstanceIdsForAsg(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueId()
	asgName := fmt.Sprintf("%s-%s", t.Name(), uniqueID)
	region := GetRandomStableRegion(t, []string{}, []string{})

	defer deleteAutoScalingGroup(t, asgName, region)
	createTestAutoScalingGroup(t, asgName, region, 1)
	WaitForCapacity(t, asgName, region, 40, 15*time.Second)

	instanceIds := GetInstanceIdsForAsg(t, asgName, region)
	assert.Equal(t, len(instanceIds), 1)
}

// The following functions were adapted from the tests for cloud-nuke

func createTestAutoScalingGroup(t *testing.T, name string, region string, desiredCount int64) {
	instance := createTestEC2Instance(t, region, name)

	asgClient := NewAsgClient(t, region)
	param := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName: &name,
		InstanceId:           instance.InstanceId,
		DesiredCapacity:      aws.Int64(desiredCount),
		MinSize:              aws.Int64(1),
		MaxSize:              aws.Int64(3),
	}
	_, err := asgClient.CreateAutoScalingGroup(param)
	require.NoError(t, err)

	err = asgClient.WaitUntilGroupExists(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{&name},
	})
	require.NoError(t, err)
}

func createTestEC2Instance(t *testing.T, region string, name string) ec2.Instance {
	ec2Client := NewEc2Client(t, region)
	imageID := GetAmazonLinuxAmi(t, region)
	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(imageID),
		InstanceType: aws.String(GetRecommendedInstanceType(t, region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	}
	runResult, err := ec2Client.RunInstances(params)
	require.NoError(t, err)

	require.NotEqual(t, len(runResult.Instances), 0)

	err = ec2Client.WaitUntilInstanceExists(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-id"),
				Values: []*string{runResult.Instances[0].InstanceId},
			},
		},
	})
	require.NoError(t, err)

	// Add test tag to the created instance
	_, err = ec2Client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
	})
	require.NoError(t, err)

	// EC2 Instance must be in a running before this function returns
	err = ec2Client.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("instance-id"),
				Values: []*string{runResult.Instances[0].InstanceId},
			},
		},
	})
	require.NoError(t, err)

	return *runResult.Instances[0]
}

func terminateEc2InstancesByName(t *testing.T, region string, names []string) {
	for _, name := range names {
		instanceIds := GetEc2InstanceIdsByTag(t, region, "Name", name)
		for _, instanceId := range instanceIds {
			TerminateInstance(t, region, instanceId)
		}
	}
}

func deleteAutoScalingGroup(t *testing.T, name string, region string) {
	// We have to scale ASG down to 0 before we can delete it
	scaleAsgToZero(t, name, region)

	asgClient := NewAsgClient(t, region)
	input := &autoscaling.DeleteAutoScalingGroupInput{AutoScalingGroupName: aws.String(name)}
	_, err := asgClient.DeleteAutoScalingGroup(input)
	require.NoError(t, err)
	err = asgClient.WaitUntilGroupNotExists(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{aws.String(name)},
	})
	require.NoError(t, err)
}

func scaleAsgToZero(t *testing.T, name string, region string) {
	asgClient := NewAsgClient(t, region)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int64(0),
		MinSize:              aws.Int64(0),
		MaxSize:              aws.Int64(0),
	}
	_, err := asgClient.UpdateAutoScalingGroup(input)
	require.NoError(t, err)
	WaitForCapacity(t, name, region, 40, 15*time.Second)

	// There is an eventual consistency bug where even though the ASG is scaled down, AWS sometimes still views a
	// scaling activity so we add a 5 second pause here to work around it.
	time.Sleep(5 * time.Second)
}
