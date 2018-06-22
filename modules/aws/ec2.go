package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
)

// GetPublicIpOfEc2Instance gets the public IP address of the given EC2 Instance in the given region.
func GetPublicIpOfEc2Instance(t *testing.T, instanceID string, awsRegion string) string {
	ip, err := GetPublicIpOfEc2InstanceE(t, instanceID, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return ip
}

// GetPublicIpOfEc2InstanceE gets the public IP address of the given EC2 Instance in the given region.
func GetPublicIpOfEc2InstanceE(t *testing.T, instanceID string, awsRegion string) (string, error) {
	ips, err := GetPublicIpsOfEc2InstancesE(t, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	ip, containsIP := ips[instanceID]

	if !containsIP {
		return "", IpForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion}
	}

	return ip, nil
}

// GetPublicIpsOfEc2Instances gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPublicIpsOfEc2Instances(t *testing.T, instanceIDs []string, awsRegion string) map[string]string {
	ips, err := GetPublicIpsOfEc2InstancesE(t, instanceIDs, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return ips
}

// GetPublicIpsOfEc2InstancesE gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPublicIpsOfEc2InstancesE(t *testing.T, instanceIDs []string, awsRegion string) (map[string]string, error) {
	ec2Client := NewEc2Client(t, awsRegion)

	input := ec2.DescribeInstancesInput{InstanceIds: aws.StringSlice(instanceIDs)}
	output, err := ec2Client.DescribeInstances(&input)
	if err != nil {
		return nil, err
	}

	ips := map[string]string{}

	for _, reserveration := range output.Reservations {
		for _, instance := range reserveration.Instances {
			ips[aws.StringValue(instance.InstanceId)] = aws.StringValue(instance.PublicIpAddress)
		}
	}

	return ips, nil
}

// GetEc2InstanceIdsByTag returns all the IDs of EC2 instances in the given region with the given tag.
func GetEc2InstanceIdsByTag(t *testing.T, region string, tagName string, tagValue string) []string {
	out, err := GetEc2InstanceIdsByTagE(t, region, tagName, tagValue)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetEc2InstanceIdsByTagE returns all the IDs of EC2 instances in the given region with the given tag.
func GetEc2InstanceIdsByTagE(t *testing.T, region string, tagName string, tagValue string) ([]string, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	tagFilter := &ec2.Filter{
		Name:   aws.String(fmt.Sprintf("tag:%s", tagName)),
		Values: []*string{aws.String(tagValue)},
	}
	output, err := client.DescribeInstances(&ec2.DescribeInstancesInput{Filters: []*ec2.Filter{tagFilter}})
	if err != nil {
		return nil, err
	}

	instanceIDs := []string{}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			instanceIDs = append(instanceIDs, *instance.InstanceId)
		}
	}

	return instanceIDs, err
}

// GetTagsForEc2Instance returns all the tags for the given EC2 Instance.
func GetTagsForEc2Instance(t *testing.T, region string, instanceID string) map[string]string {
	tags, err := GetTagsForEc2InstanceE(t, region, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	return tags
}

// GetTagsForEc2InstanceE returns all the tags for the given EC2 Instance.
func GetTagsForEc2InstanceE(t *testing.T, region string, instanceID string) (map[string]string, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	input := ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-type"),
				Values: aws.StringSlice([]string{"instance"}),
			},
			{
				Name:   aws.String("resource-id"),
				Values: aws.StringSlice([]string{instanceID}),
			},
		},
	}

	out, err := client.DescribeTags(&input)
	if err != nil {
		return nil, err
	}

	tags := map[string]string{}

	for _, tag := range out.Tags {
		tags[aws.StringValue(tag.Key)] = aws.StringValue(tag.Value)
	}

	return tags, nil
}

// DeleteAmi deletes the given AMI in the given region.
func DeleteAmi(t *testing.T, region string, imageID string) {
	err := DeleteAmiE(t, region, imageID)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteAmiE deletes the given AMI in the given region.
func DeleteAmiE(t *testing.T, region string, imageID string) error {
	logger.Logf(t, "Deregistering AMI %s", imageID)

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	_, err = client.DeregisterImage(&ec2.DeregisterImageInput{ImageId: aws.String(imageID)})
	return err
}

// AddTagsToResource adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
func AddTagsToResource(t *testing.T, region string, resource string, tags map[string]string) {
	err := AddTagsToResourceE(t, region, resource, tags)
	if err != nil {
		t.Fatal(err)
	}
}

// AddTagsToResourceE adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
func AddTagsToResourceE(t *testing.T, region string, resource string, tags map[string]string) error {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	var awsTags []*ec2.Tag
	for key, value := range tags {
		awsTags = append(awsTags, &ec2.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	_, err = client.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{aws.String(resource)},
		Tags:      awsTags,
	})

	return err
}

// TerminateInstance terminates the EC2 instance with the given ID in the given region.
func TerminateInstance(t *testing.T, region string, instanceID string) {
	err := TerminateInstanceE(t, region, instanceID)
	if err != nil {
		t.Fatal(err)
	}
}

// TerminateInstanceE terminates the EC2 instance with the given ID in the given region.
func TerminateInstanceE(t *testing.T, region string, instanceID string) error {
	logger.Logf(t, "Terminating Instance %s", instanceID)

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	_, err = client.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	})

	return err
}

// NewEc2Client creates an EC2 client.
func NewEc2Client(t *testing.T, region string) *ec2.EC2 {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewEc2ClientE creates an EC2 client.
func NewEc2ClientE(t *testing.T, region string) (*ec2.EC2, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return ec2.New(sess), nil
}

// IpForEc2InstanceNotFound is an error that occurs when the IP for an EC2 instance is not found.
type IpForEc2InstanceNotFound struct {
	InstanceId string
	AwsRegion  string
}

func (err IpForEc2InstanceNotFound) Error() string {
	return fmt.Sprintf("Could not find a public IP address for EC2 Instance %s in %s", err.InstanceId, err.AwsRegion)
}
