package aws_test

import (
	"context"
	"testing"
	"time"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	autoscalingTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetCapacityInfoForAsg(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	asgName := t.Name() + "-" + uniqueID
	region := aws.GetRandomStableRegion(t, []string{}, []string{})

	defer deleteAutoScalingGroup(t, asgName, region)

	createTestAutoScalingGroup(t, asgName, region, 2)
	aws.WaitForCapacity(t, asgName, region, 40, 15*time.Second)

	capacityInfo := aws.GetCapacityInfoForAsg(t, asgName, region)
	assert.Equal(t, int64(2), capacityInfo.DesiredCapacity)
	assert.Equal(t, int64(2), capacityInfo.CurrentCapacity)
	assert.Equal(t, int64(1), capacityInfo.MinCapacity)
	assert.Equal(t, int64(3), capacityInfo.MaxCapacity)
}

func TestGetInstanceIdsForAsg(t *testing.T) {
	t.Parallel()

	uniqueID := random.UniqueID()
	asgName := t.Name() + "-" + uniqueID
	region := aws.GetRandomStableRegion(t, []string{}, []string{})

	defer deleteAutoScalingGroup(t, asgName, region)

	createTestAutoScalingGroup(t, asgName, region, 1)
	aws.WaitForCapacity(t, asgName, region, 40, 15*time.Second)

	instanceIDs := aws.GetInstanceIdsForAsg(t, asgName, region)
	assert.Len(t, instanceIDs, 1)
}

func createTestAutoScalingGroup(t *testing.T, name string, region string, desiredCount int32) {
	t.Helper()

	azs := aws.GetAvailabilityZones(t, region)
	ec2Client := aws.NewEc2Client(t, region)
	imageID := aws.GetAmazonLinuxAmi(t, region)
	template, err := ec2Client.CreateLaunchTemplate(context.Background(), &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			ImageId:      awsSDK.String(imageID),
			InstanceType: types.InstanceType(aws.GetRecommendedInstanceType(t, region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})),
		},
		LaunchTemplateName: awsSDK.String(name),
	})
	require.NoError(t, err)

	asgClient := aws.NewAsgClient(t, region)
	param := &autoscaling.CreateAutoScalingGroupInput{
		AutoScalingGroupName: &name,
		LaunchTemplate: &autoscalingTypes.LaunchTemplateSpecification{
			LaunchTemplateId: template.LaunchTemplate.LaunchTemplateId,
			Version:          awsSDK.String("$Latest"),
		},
		AvailabilityZones: azs,
		DesiredCapacity:   awsSDK.Int32(desiredCount),
		MinSize:           awsSDK.Int32(1),
		MaxSize:           awsSDK.Int32(3),
	}
	_, err = asgClient.CreateAutoScalingGroup(context.Background(), param)
	require.NoError(t, err)

	waiter := autoscaling.NewGroupExistsWaiter(asgClient)
	err = waiter.Wait(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}, 42*time.Minute)
	require.NoError(t, err)
}

func deleteAutoScalingGroup(t *testing.T, name string, region string) {
	t.Helper()

	// We have to scale ASG down to 0 before we can delete it
	scaleAsgToZero(t, name, region)

	asgClient := aws.NewAsgClient(t, region)
	input := &autoscaling.DeleteAutoScalingGroupInput{AutoScalingGroupName: awsSDK.String(name)}
	_, err := asgClient.DeleteAutoScalingGroup(context.Background(), input)
	require.NoError(t, err)

	waiter := autoscaling.NewGroupNotExistsWaiter(asgClient)
	err = waiter.Wait(context.Background(), &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	}, 40*time.Minute)
	require.NoError(t, err)

	ec2Client := aws.NewEc2Client(t, region)
	_, err = ec2Client.DeleteLaunchTemplate(context.Background(), &ec2.DeleteLaunchTemplateInput{
		LaunchTemplateName: awsSDK.String(name),
	})
	require.NoError(t, err)
}

func scaleAsgToZero(t *testing.T, name string, region string) {
	t.Helper()

	asgClient := aws.NewAsgClient(t, region)
	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: awsSDK.String(name),
		DesiredCapacity:      awsSDK.Int32(0),
		MinSize:              awsSDK.Int32(0),
		MaxSize:              awsSDK.Int32(0),
	}
	_, err := asgClient.UpdateAutoScalingGroup(context.Background(), input)
	require.NoError(t, err)
	aws.WaitForCapacity(t, name, region, 40, 15*time.Second)

	// There is an eventual consistency bug where even though the ASG is scaled down, AWS sometimes still views a
	// scaling activity so we add a 5-second pause here to work around it.
	time.Sleep(5 * time.Second)
}
