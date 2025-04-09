package aws

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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

func createTestAutoScalingGroup(t *testing.T, name string, region string, desiredCount int32) {
	azs := GetAvailabilityZones(t, region)
	ec2Client := NewEc2Client(t, region)
	imageID := GetAmazonLinuxAmi(t, region)
	template, err := ec2Client.CreateLaunchTemplate(context.Background(), &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			ImageId:      aws.String(imageID),
			InstanceType: types.InstanceType(GetRecommendedInstanceType(t, region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})),
		},
		LaunchTemplateName: aws.String(name),
	})
	require.NoError(t, err)

	asgClient := NewAsgClient(t, region)
	param := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName: &name,
		LaunchTemplate: &autoscalingTypes.LaunchTemplateSpecification{
			LaunchTemplateId: template.LaunchTemplate.LaunchTemplateId,
			Version:          aws.String("$Latest"),
		},
		AvailabilityZones: azs,
		DesiredCapacity:   aws.Int32(desiredCount),
		MinSize:           aws.Int32(1),
		MaxSize:           aws.Int32(3),
	}
	_, err = asgClient.CreateAutoScalingGroup(context.Background(), param)
	require.NoError(t, err)

	waiter := autoscaling.NewGroupExistsWaiter(asgClient)
	err = waiter.Wait(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}, 42*time.Minute)
	require.NoError(t, err)
}

func createTestEC2Instance(t *testing.T, region string, name string) types.Instance {
	ec2Client := NewEc2Client(t, region)
	imageID := GetAmazonLinuxAmi(t, region)
	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(imageID),
		InstanceType: types.InstanceType(GetRecommendedInstanceType(t, region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}
	runResult, err := ec2Client.RunInstances(context.Background(), params)
	require.NoError(t, err)

	require.NotEqual(t, len(runResult.Instances), 0)

	waiter := ec2.NewInstanceExistsWaiter(ec2Client)
	err = waiter.Wait(
		context.Background(),
		&ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("instance-id"),
					Values: []string{*runResult.Instances[0].InstanceId},
				},
			},
		},
		42*time.Minute,
	)
	require.NoError(t, err)

	// Add test tag to the created instance
	_, err = ec2Client.CreateTags(context.Background(), &ec2.CreateTagsInput{
		Resources: []string{*runResult.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
	})
	require.NoError(t, err)

	// EC2 Instance must be in a running before this function returns
	runningWaiter := ec2.NewInstanceRunningWaiter(ec2Client)
	err = runningWaiter.Wait(context.Background(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []string{*runResult.Instances[0].InstanceId},
			},
		},
	}, 42*time.Minute)
	require.NoError(t, err)

	return runResult.Instances[0]
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
	_, err := asgClient.DeleteAutoScalingGroup(context.Background(), input)
	require.NoError(t, err)

	waiter := autoscaling.NewGroupNotExistsWaiter(asgClient)
	err = waiter.Wait(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}, 40*time.Minute)
	require.NoError(t, err)

	ec2Client := NewEc2Client(t, region)
	_, err = ec2Client.DeleteLaunchTemplate(context.Background(), &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateName: aws.String(name),
	})
	require.NoError(t, err)
}

func scaleAsgToZero(t *testing.T, name string, region string) {
	asgClient := NewAsgClient(t, region)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(name),
		DesiredCapacity:      aws.Int32(0),
		MinSize:              aws.Int32(0),
		MaxSize:              aws.Int32(0),
	}
	_, err := asgClient.UpdateAutoScalingGroup(context.Background(), input)
	require.NoError(t, err)
	WaitForCapacity(t, name, region, 40, 15*time.Second)

	// There is an eventual consistency bug where even though the ASG is scaled down, AWS sometimes still views a
	// scaling activity so we add a 5-second pause here to work around it.
	time.Sleep(5 * time.Second)
}
