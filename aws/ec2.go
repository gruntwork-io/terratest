package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/gruntwork-io/gruntwork-cli/logging"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Return all the ids of EC2 instances in the given region with the given tag
func GetEc2InstanceIdsByTag(awsRegion string, tagName string, tagValue string) ([]string, error) {
	instanceIds := []string{}
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

	tagFilter := &ec2.Filter{
		Name: aws.String(fmt.Sprintf("tag:%s", tagName)),
		Values: []*string{aws.String(tagValue)},
	}
	output, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{Filters: []*ec2.Filter{tagFilter}})
	if err != nil {
		return instanceIds, err
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			instanceIds = append(instanceIds, *instance.InstanceId)
		}
	}

	return instanceIds, err
}

//Remove an AMI
func RemoveAmi(ec2Client *ec2.EC2, imageId string) (error) {
	logger := logging.GetLogger("RemoveAmi")
	logger.Debug("Deregistering AMI %s\n", imageId)

	_, err := ec2Client.DeregisterImage(&ec2.DeregisterImageInput{ImageId: aws.String(imageId) })
	if err != nil {
		return err
	}
	return nil
}

//Terminate an instance
func TerminateInstance(ec2Client *ec2.EC2, instanceId string) (error) {
	logger := logging.GetLogger("TerminateInstance")
	logger.Debug("Terminating Instance %s", instanceId)

	_, err := ec2Client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func CreateEC2Client(awsRegion string) (*ec2.EC2, error) {
	awsConfig, err := CreateAwsConfig(awsRegion)
	if err != nil {
		return nil, err
	}

	return ec2.New(session.New(), awsConfig), nil
}
