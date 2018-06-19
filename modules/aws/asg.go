package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// GetInstanceIdsForAsg gets the IDs of EC2 Instances in the given ASG.
func GetInstanceIdsForAsg(t *testing.T, asgName string, awsRegion string, sessExists ...*session.Session) []string {
	ids, err := GetInstanceIdsForAsgE(t, asgName, awsRegion, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsForAsgE gets the IDs of EC2 Instances in the given ASG.
func GetInstanceIdsForAsgE(t *testing.T, asgName string, awsRegion string, sessExists ...*session.Session) ([]string, error) {
	asgClient, err := NewAsgClientE(t, awsRegion, sessExists[0])
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

// NewAsgClient creates an Auto Scaling Group client.
func NewAsgClient(t *testing.T, region string, sessExists ...*session.Session) *autoscaling.AutoScaling {
	client, err := NewAsgClientE(t, region, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewAsgClientE creates an Auto Scaling Group client.
func NewAsgClientE(t *testing.T, region string, sessExists ...*session.Session) (*autoscaling.AutoScaling, error) {
	sess, err := NewAuthenticatedSession(region, sessExists[0])
	if err != nil {
		return nil, err
	}

	return autoscaling.New(sess), nil
}
