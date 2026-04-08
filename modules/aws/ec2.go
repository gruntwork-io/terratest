package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetPrivateIPOfEc2Instance gets the private IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
func GetPrivateIPOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	ip, err := GetPrivateIPOfEc2InstanceE(t, instanceID, awsRegion)
	require.NoError(t, err)

	return ip
}

// GetPrivateIPOfEc2InstanceE gets the private IP address of the given EC2 Instance in the given region.
func GetPrivateIPOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	ips, err := GetPrivateIpsOfEc2InstancesE(t, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	ip, containsIP := ips[instanceID]

	if !containsIP {
		return "", IpForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion, Type: "private"}
	}

	return ip, nil
}

// GetPrivateIpsOfEc2Instances gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPrivateIpsOfEc2Instances(t testing.TestingT, instanceIDs []string, awsRegion string) map[string]string {
	ips, err := GetPrivateIpsOfEc2InstancesE(t, instanceIDs, awsRegion)
	require.NoError(t, err)

	return ips
}

// GetPrivateIpsOfEc2InstancesE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPrivateIpsOfEc2InstancesE(t testing.TestingT, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return getInstanceFieldMapE(t, instanceIDs, awsRegion, func(inst *types.Instance) string {
		return aws.ToString(inst.PrivateIpAddress)
	})
}

// GetPrivateHostnameOfEc2Instance gets the private IP address of the given EC2 Instance in the given region.
func GetPrivateHostnameOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	ip, err := GetPrivateHostnameOfEc2InstanceE(t, instanceID, awsRegion)
	require.NoError(t, err)

	return ip
}

// GetPrivateHostnameOfEc2InstanceE gets the private IP address of the given EC2 Instance in the given region.
func GetPrivateHostnameOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	hostnames, err := GetPrivateHostnamesOfEc2InstancesE(t, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	hostname, containsHostname := hostnames[instanceID]

	if !containsHostname {
		return "", HostnameForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion, Type: "private"}
	}

	return hostname, nil
}

// GetPrivateHostnamesOfEc2Instances gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPrivateHostnamesOfEc2Instances(t testing.TestingT, instanceIDs []string, awsRegion string) map[string]string {
	ips, err := GetPrivateHostnamesOfEc2InstancesE(t, instanceIDs, awsRegion)
	require.NoError(t, err)

	return ips
}

// GetPrivateHostnamesOfEc2InstancesE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPrivateHostnamesOfEc2InstancesE(t testing.TestingT, instanceIDs []string, awsRegion string) (map[string]string, error) {
	ec2Client, err := NewEc2ClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	input := ec2.DescribeInstancesInput{InstanceIds: instanceIDs}

	hostnames := map[string]string{}

	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, &input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, reservation := range page.Reservations {
			for j := range reservation.Instances {
				instance := &reservation.Instances[j]
				hostnames[aws.ToString(instance.InstanceId)] = aws.ToString(instance.PrivateDnsName)
			}
		}
	}

	return hostnames, nil
}

// GetPublicIPOfEc2Instance gets the public IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
func GetPublicIPOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	ip, err := GetPublicIPOfEc2InstanceE(t, instanceID, awsRegion)
	require.NoError(t, err)

	return ip
}

// GetPublicIPOfEc2InstanceE gets the public IP address of the given EC2 Instance in the given region.
func GetPublicIPOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	ips, err := GetPublicIpsOfEc2InstancesE(t, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	ip, containsIP := ips[instanceID]

	if !containsIP {
		return "", IpForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion, Type: "public"}
	}

	return ip, nil
}

// GetPublicIpsOfEc2Instances gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPublicIpsOfEc2Instances(t testing.TestingT, instanceIDs []string, awsRegion string) map[string]string {
	ips, err := GetPublicIpsOfEc2InstancesE(t, instanceIDs, awsRegion)
	require.NoError(t, err)

	return ips
}

// GetPublicIpsOfEc2InstancesE gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
func GetPublicIpsOfEc2InstancesE(t testing.TestingT, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return getInstanceFieldMapE(t, instanceIDs, awsRegion, func(inst *types.Instance) string {
		return aws.ToString(inst.PublicIpAddress)
	})
}

// GetEc2InstanceIdsByTag returns all the IDs of EC2 instances in the given region with the given tag.
func GetEc2InstanceIdsByTag(t testing.TestingT, region string, tagName string, tagValue string) []string {
	out, err := GetEc2InstanceIdsByTagE(t, region, tagName, tagValue)
	require.NoError(t, err)

	return out
}

// GetEc2InstanceIdsByTagE returns all the IDs of EC2 instances in the given region with the given tag.
func GetEc2InstanceIdsByTagE(t testing.TestingT, region string, tagName string, tagValue string) ([]string, error) {
	ec2Filters := map[string][]string{
		"tag:" + tagName: {tagValue},
	}

	return GetEc2InstanceIdsByFiltersE(t, region, ec2Filters)
}

// GetEc2InstanceIdsByFilters returns all the IDs of EC2 instances in the given region which match to EC2 filter list
// as per https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput.
func GetEc2InstanceIdsByFilters(t testing.TestingT, region string, ec2Filters map[string][]string) []string {
	out, err := GetEc2InstanceIdsByFiltersE(t, region, ec2Filters)
	require.NoError(t, err)

	return out
}

// GetEc2InstanceIdsByFiltersE returns all the IDs of EC2 instances in the given region which match to EC2 filter list
// as per https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput.
func GetEc2InstanceIdsByFiltersE(t testing.TestingT, region string, ec2Filters map[string][]string) ([]string, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	var ec2FilterList []types.Filter

	for name, values := range ec2Filters {
		ec2FilterList = append(ec2FilterList, types.Filter{Name: aws.String(name), Values: values})
	}

	var instanceIDs []string

	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{Filters: ec2FilterList})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, reservation := range page.Reservations {
			for j := range reservation.Instances {
				instance := &reservation.Instances[j]
				instanceIDs = append(instanceIDs, *instance.InstanceId)
			}
		}
	}

	return instanceIDs, nil
}

// GetTagsForEc2Instance returns all the tags for the given EC2 Instance.
func GetTagsForEc2Instance(t testing.TestingT, region string, instanceID string) map[string]string {
	tags, err := GetTagsForEc2InstanceE(t, region, instanceID)
	require.NoError(t, err)

	return tags
}

// GetTagsForEc2InstanceE returns all the tags for the given EC2 Instance.
func GetTagsForEc2InstanceE(t testing.TestingT, region string, instanceID string) (map[string]string, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	input := ec2.DescribeTagsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("resource-type"),
				Values: []string{"instance"},
			},
			{
				Name:   aws.String("resource-id"),
				Values: []string{instanceID},
			},
		},
	}

	tags := map[string]string{}

	paginator := ec2.NewDescribeTagsPaginator(client, &input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, tag := range page.Tags {
			tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}
	}

	return tags, nil
}

// DeleteAmi deletes the given AMI in the given region.
func DeleteAmi(t testing.TestingT, region string, imageID string) {
	require.NoError(t, DeleteAmiE(t, region, imageID))
}

// DeleteAmiE deletes the given AMI in the given region.
func DeleteAmiE(t testing.TestingT, region string, imageID string) error {
	logger.Default.Logf(t, "Deregistering AMI %s", imageID)

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	_, err = client.DeregisterImage(context.Background(), &ec2.DeregisterImageInput{ImageId: aws.String(imageID)})

	return err
}

// AddTagsToResource adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
func AddTagsToResource(t testing.TestingT, region string, resource string, tags map[string]string) {
	require.NoError(t, AddTagsToResourceE(t, region, resource, tags))
}

// AddTagsToResourceE adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
func AddTagsToResourceE(t testing.TestingT, region string, resource string, tags map[string]string) error {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	var awsTags []types.Tag

	for key, value := range tags {
		awsTags = append(awsTags, types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	_, err = client.CreateTags(context.Background(), &ec2.CreateTagsInput{
		Resources: []string{resource},
		Tags:      awsTags,
	})

	return err
}

// TerminateInstance terminates the EC2 instance with the given ID in the given region.
func TerminateInstance(t testing.TestingT, region string, instanceID string) {
	require.NoError(t, TerminateInstanceE(t, region, instanceID))
}

// TerminateInstanceE terminates the EC2 instance with the given ID in the given region.
func TerminateInstanceE(t testing.TestingT, region string, instanceID string) error {
	logger.Default.Logf(t, "Terminating Instance %s", instanceID)

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	_, err = client.TerminateInstances(context.Background(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	})

	return err
}

// GetAmiPubliclyAccessible returns whether the AMI is publicly accessible or not
func GetAmiPubliclyAccessible(t testing.TestingT, awsRegion string, amiID string) bool {
	output, err := GetAmiPubliclyAccessibleE(t, awsRegion, amiID)
	require.NoError(t, err)

	return output
}

// GetAmiPubliclyAccessibleE returns whether the AMI is publicly accessible or not
func GetAmiPubliclyAccessibleE(t testing.TestingT, awsRegion string, amiID string) (bool, error) {
	launchPermissions, err := GetLaunchPermissionsForAmiE(t, awsRegion, amiID)
	if err != nil {
		return false, err
	}

	for _, launchPermission := range launchPermissions {
		if string(launchPermission.Group) == "all" {
			return true, nil
		}
	}

	return false, nil
}

// GetAccountsWithLaunchPermissionsForAmi returns list of accounts that the AMI is shared with
func GetAccountsWithLaunchPermissionsForAmi(t testing.TestingT, awsRegion string, amiID string) []string {
	output, err := GetAccountsWithLaunchPermissionsForAmiE(t, awsRegion, amiID)
	require.NoError(t, err)

	return output
}

// GetAccountsWithLaunchPermissionsForAmiE returns list of accounts that the AMI is shared with
func GetAccountsWithLaunchPermissionsForAmiE(t testing.TestingT, awsRegion string, amiID string) ([]string, error) {
	var accountIDs []string

	launchPermissions, err := GetLaunchPermissionsForAmiE(t, awsRegion, amiID)
	if err != nil {
		return accountIDs, err
	}

	for _, launchPermission := range launchPermissions {
		if aws.ToString(launchPermission.UserId) != "" {
			accountIDs = append(accountIDs, aws.ToString(launchPermission.UserId))
		}
	}

	return accountIDs, nil
}

// GetLaunchPermissionsForAmiE returns launchPermissions as configured in AWS
func GetLaunchPermissionsForAmiE(t testing.TestingT, awsRegion string, amiID string) ([]types.LaunchPermission, error) {
	client := NewEc2Client(t, awsRegion)
	input := &ec2.DescribeImageAttributeInput{
		Attribute: types.ImageAttributeNameLaunchPermission,
		ImageId:   aws.String(amiID),
	}

	output, err := client.DescribeImageAttribute(context.Background(), input)
	if err != nil {
		return []types.LaunchPermission{}, err
	}

	return output.LaunchPermissions, nil
}

// GetRecommendedInstanceType takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs, and if you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function will fail the test if there is an error.
func GetRecommendedInstanceType(t testing.TestingT, region string, instanceTypeOptions []string) string {
	out, err := GetRecommendedInstanceTypeE(t, region, instanceTypeOptions)
	require.NoError(t, err)

	return out
}

// GetRecommendedInstanceTypeE takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
func GetRecommendedInstanceTypeE(t testing.TestingT, region string, instanceTypeOptions []string) (string, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return "", err
	}

	return GetRecommendedInstanceTypeWithClientE(t, client, instanceTypeOptions)
}

// GetRecommendedInstanceTypeWithClientE takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function expects an authenticated EC2 client from the AWS SDK Go library.
func GetRecommendedInstanceTypeWithClientE(t testing.TestingT, ec2Client *ec2.Client, instanceTypeOptions []string) (string, error) {
	availabilityZones, err := getAllAvailabilityZonesE(ec2Client)
	if err != nil {
		return "", err
	}

	instanceTypeOfferings, err := getInstanceTypeOfferingsE(ec2Client, instanceTypeOptions)
	if err != nil {
		return "", err
	}

	return PickRecommendedInstanceTypeE(availabilityZones, instanceTypeOfferings, instanceTypeOptions)
}

// PickRecommendedInstanceTypeE picks the first instance type from instanceTypeOptions that is available in all the
// given availability zones based on the given instance type offerings. Returns a NoInstanceTypeError if none of
// the options are available in all AZs.
func PickRecommendedInstanceTypeE(availabilityZones []string, instanceTypeOfferings []types.InstanceTypeOffering, instanceTypeOptions []string) (string, error) {
	// O(n^3) for the win!
	for _, instanceType := range instanceTypeOptions {
		if instanceTypeExistsInAllAzs(instanceType, availabilityZones, instanceTypeOfferings) {
			return instanceType, nil
		}
	}

	return "", NoInstanceTypeError{InstanceTypeOptions: instanceTypeOptions, Azs: availabilityZones}
}

// instanceTypeExistsInAllAzs returns true if the given instance type exists in all the given availabilityZones based
// on the availability data in instanceTypeOfferings
func instanceTypeExistsInAllAzs(instanceType string, availabilityZones []string, instanceTypeOfferings []types.InstanceTypeOffering) bool {
	if len(availabilityZones) == 0 || len(instanceTypeOfferings) == 0 {
		return false
	}

	for _, az := range availabilityZones {
		if !hasOffering(instanceTypeOfferings, az, instanceType) {
			return false
		}
	}

	return true
}

// hasOffering returns true if the given availability zone and instance type are one of the offerings in
// instanceTypeOfferings
func hasOffering(instanceTypeOfferings []types.InstanceTypeOffering, availabilityZone string, instanceType string) bool {
	for _, offering := range instanceTypeOfferings {
		if string(offering.InstanceType) == instanceType && aws.ToString(offering.Location) == availabilityZone {
			return true
		}
	}

	return false
}

// getInstanceTypeOfferingsE returns the instance types from the given list that are available in the region configured
// in the given EC2 client
func getInstanceTypeOfferingsE(client *ec2.Client, instanceTypeOptions []string) ([]types.InstanceTypeOffering, error) {
	input := ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: types.LocationTypeAvailabilityZone,
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: instanceTypeOptions,
			},
		},
	}

	out, err := client.DescribeInstanceTypeOfferings(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return out.InstanceTypeOfferings, nil
}

// getAllAvailabilityZonesE returns all the available AZs in the region configured in the given EC2 client
func getAllAvailabilityZonesE(client *ec2.Client) ([]string, error) {
	input := ec2.DescribeAvailabilityZonesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
		},
	}

	out, err := client.DescribeAvailabilityZones(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	var azs []string

	for i := range out.AvailabilityZones {
		az := &out.AvailabilityZones[i]
		azs = append(azs, aws.ToString(az.ZoneName))
	}

	return azs, nil
}

// getInstanceFieldMapE is a shared helper that paginates through DescribeInstances and builds a map
// of instance ID to a string field extracted by the given function.
func getInstanceFieldMapE(t testing.TestingT, instanceIDs []string, awsRegion string, extractField func(*types.Instance) string) (map[string]string, error) {
	ec2Client := NewEc2Client(t, awsRegion)
	input := ec2.DescribeInstancesInput{InstanceIds: instanceIDs}

	result := map[string]string{}

	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, &input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		for _, reservation := range page.Reservations {
			for j := range reservation.Instances {
				instance := &reservation.Instances[j]
				result[aws.ToString(instance.InstanceId)] = extractField(instance)
			}
		}
	}

	return result, nil
}

// NewEc2Client creates an EC2 client.
func NewEc2Client(t testing.TestingT, region string) *ec2.Client {
	client, err := NewEc2ClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewEc2ClientE creates an EC2 client.
func NewEc2ClientE(t testing.TestingT, region string) (*ec2.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(*sess), nil
}
