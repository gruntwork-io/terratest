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

// GetPrivateIPOfEc2InstanceContextE gets the private IP address of the given EC2 Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateIPOfEc2InstanceContextE(t testing.TestingT, ctx context.Context, instanceID string, awsRegion string) (string, error) {
	ips, err := GetPrivateIpsOfEc2InstancesContextE(t, ctx, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	ip, containsIP := ips[instanceID]

	if !containsIP {
		return "", IpForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion, Type: "private"}
	}

	return ip, nil
}

// GetPrivateIPOfEc2InstanceContext gets the private IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateIPOfEc2InstanceContext(t testing.TestingT, ctx context.Context, instanceID string, awsRegion string) string {
	t.Helper()
	ip, err := GetPrivateIPOfEc2InstanceContextE(t, ctx, instanceID, awsRegion)
	require.NoError(t, err)

	return ip
}

// GetPrivateIPOfEc2Instance gets the private IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
//
// Deprecated: Use [GetPrivateIPOfEc2InstanceContext] instead.
func GetPrivateIPOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	t.Helper()

	return GetPrivateIPOfEc2InstanceContext(t, context.Background(), instanceID, awsRegion)
}

// GetPrivateIPOfEc2InstanceE gets the private IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPrivateIPOfEc2InstanceContextE] instead.
func GetPrivateIPOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	return GetPrivateIPOfEc2InstanceContextE(t, context.Background(), instanceID, awsRegion)
}

// GetPrivateIpOfEc2Instance gets the private IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPrivateIPOfEc2Instance] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetPrivateIpOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	return GetPrivateIPOfEc2Instance(t, instanceID, awsRegion)
}

// GetPrivateIpOfEc2InstanceE gets the private IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPrivateIPOfEc2InstanceE] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetPrivateIpOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	return GetPrivateIPOfEc2InstanceE(t, instanceID, awsRegion)
}

// GetPrivateIpsOfEc2InstancesContextE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateIpsOfEc2InstancesContextE(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return getInstanceFieldMapContextE(t, ctx, instanceIDs, awsRegion, func(inst *types.Instance) string {
		return aws.ToString(inst.PrivateIpAddress)
	})
}

// GetPrivateIpsOfEc2InstancesE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
//
// Deprecated: Use [GetPrivateIpsOfEc2InstancesContextE] instead.
func GetPrivateIpsOfEc2InstancesE(t testing.TestingT, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return GetPrivateIpsOfEc2InstancesContextE(t, context.Background(), instanceIDs, awsRegion)
}

// GetPrivateIpsOfEc2InstancesContext gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateIpsOfEc2InstancesContext(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string) map[string]string {
	t.Helper()
	ips, err := GetPrivateIpsOfEc2InstancesContextE(t, ctx, instanceIDs, awsRegion)
	require.NoError(t, err)

	return ips
}

// GetPrivateIpsOfEc2Instances gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
//
// Deprecated: Use [GetPrivateIpsOfEc2InstancesContext] instead.
func GetPrivateIpsOfEc2Instances(t testing.TestingT, instanceIDs []string, awsRegion string) map[string]string {
	t.Helper()
	return GetPrivateIpsOfEc2InstancesContext(t, context.Background(), instanceIDs, awsRegion)
}

// GetPrivateHostnameOfEc2InstanceContextE gets the private IP address of the given EC2 Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateHostnameOfEc2InstanceContextE(t testing.TestingT, ctx context.Context, instanceID string, awsRegion string) (string, error) {
	hostnames, err := GetPrivateHostnamesOfEc2InstancesContextE(t, ctx, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	hostname, containsHostname := hostnames[instanceID]

	if !containsHostname {
		return "", HostnameForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion, Type: "private"}
	}

	return hostname, nil
}

// GetPrivateHostnameOfEc2InstanceContext gets the private IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateHostnameOfEc2InstanceContext(t testing.TestingT, ctx context.Context, instanceID string, awsRegion string) string {
	t.Helper()
	ip, err := GetPrivateHostnameOfEc2InstanceContextE(t, ctx, instanceID, awsRegion)
	require.NoError(t, err)

	return ip
}

// GetPrivateHostnameOfEc2Instance gets the private IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPrivateHostnameOfEc2InstanceContext] instead.
func GetPrivateHostnameOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	t.Helper()
	return GetPrivateHostnameOfEc2InstanceContext(t, context.Background(), instanceID, awsRegion)
}

// GetPrivateHostnameOfEc2InstanceE gets the private IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPrivateHostnameOfEc2InstanceContextE] instead.
func GetPrivateHostnameOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	return GetPrivateHostnameOfEc2InstanceContextE(t, context.Background(), instanceID, awsRegion)
}

// GetPrivateHostnamesOfEc2InstancesContextE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateHostnamesOfEc2InstancesContextE(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return getInstanceFieldMapContextE(t, ctx, instanceIDs, awsRegion, func(inst *types.Instance) string {
		return aws.ToString(inst.PrivateDnsName)
	})
}

// GetPrivateHostnamesOfEc2InstancesContext gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPrivateHostnamesOfEc2InstancesContext(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string) map[string]string {
	t.Helper()
	ips, err := GetPrivateHostnamesOfEc2InstancesContextE(t, ctx, instanceIDs, awsRegion)
	require.NoError(t, err)

	return ips
}

// GetPrivateHostnamesOfEc2Instances gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
//
// Deprecated: Use [GetPrivateHostnamesOfEc2InstancesContext] instead.
func GetPrivateHostnamesOfEc2Instances(t testing.TestingT, instanceIDs []string, awsRegion string) map[string]string {
	t.Helper()
	return GetPrivateHostnamesOfEc2InstancesContext(t, context.Background(), instanceIDs, awsRegion)
}

// GetPrivateHostnamesOfEc2InstancesE gets the private IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
//
// Deprecated: Use [GetPrivateHostnamesOfEc2InstancesContextE] instead.
func GetPrivateHostnamesOfEc2InstancesE(t testing.TestingT, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return GetPrivateHostnamesOfEc2InstancesContextE(t, context.Background(), instanceIDs, awsRegion)
}

// GetPublicIPOfEc2InstanceContextE gets the public IP address of the given EC2 Instance in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIPOfEc2InstanceContextE(t testing.TestingT, ctx context.Context, instanceID string, awsRegion string) (string, error) {
	ips, err := GetPublicIpsOfEc2InstancesContextE(t, ctx, []string{instanceID}, awsRegion)
	if err != nil {
		return "", err
	}

	ip, containsIP := ips[instanceID]

	if !containsIP {
		return "", IpForEc2InstanceNotFound{InstanceId: instanceID, AwsRegion: awsRegion, Type: "public"}
	}

	return ip, nil
}

// GetPublicIPOfEc2InstanceContext gets the public IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIPOfEc2InstanceContext(t testing.TestingT, ctx context.Context, instanceID string, awsRegion string) string {
	t.Helper()
	ip, err := GetPublicIPOfEc2InstanceContextE(t, ctx, instanceID, awsRegion)
	require.NoError(t, err)

	return ip
}

// GetPublicIPOfEc2Instance gets the public IP address of the given EC2 Instance in the given region.
// This function will fail the test if there is an error.
//
// Deprecated: Use [GetPublicIPOfEc2InstanceContext] instead.
func GetPublicIPOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	t.Helper()

	return GetPublicIPOfEc2InstanceContext(t, context.Background(), instanceID, awsRegion)
}

// GetPublicIPOfEc2InstanceE gets the public IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPublicIPOfEc2InstanceContextE] instead.
func GetPublicIPOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	return GetPublicIPOfEc2InstanceContextE(t, context.Background(), instanceID, awsRegion)
}

// GetPublicIpOfEc2Instance gets the public IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPublicIPOfEc2Instance] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetPublicIpOfEc2Instance(t testing.TestingT, instanceID string, awsRegion string) string {
	return GetPublicIPOfEc2Instance(t, instanceID, awsRegion)
}

// GetPublicIpOfEc2InstanceE gets the public IP address of the given EC2 Instance in the given region.
//
// Deprecated: Use [GetPublicIPOfEc2InstanceE] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetPublicIpOfEc2InstanceE(t testing.TestingT, instanceID string, awsRegion string) (string, error) {
	return GetPublicIPOfEc2InstanceE(t, instanceID, awsRegion)
}

// GetPublicIpsOfEc2InstancesContextE gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIpsOfEc2InstancesContextE(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return getInstanceFieldMapContextE(t, ctx, instanceIDs, awsRegion, func(inst *types.Instance) string {
		return aws.ToString(inst.PublicIpAddress)
	})
}

// GetPublicIpsOfEc2InstancesE gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
//
// Deprecated: Use [GetPublicIpsOfEc2InstancesContextE] instead.
func GetPublicIpsOfEc2InstancesE(t testing.TestingT, instanceIDs []string, awsRegion string) (map[string]string, error) {
	return GetPublicIpsOfEc2InstancesContextE(t, context.Background(), instanceIDs, awsRegion)
}

// GetPublicIpsOfEc2InstancesContext gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPublicIpsOfEc2InstancesContext(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string) map[string]string {
	t.Helper()
	ips, err := GetPublicIpsOfEc2InstancesContextE(t, ctx, instanceIDs, awsRegion)
	require.NoError(t, err)

	return ips
}

// GetPublicIpsOfEc2Instances gets the public IP address of the given EC2 Instance in the given region. Returns a map of instance ID to IP address.
//
// Deprecated: Use [GetPublicIpsOfEc2InstancesContext] instead.
func GetPublicIpsOfEc2Instances(t testing.TestingT, instanceIDs []string, awsRegion string) map[string]string {
	t.Helper()
	return GetPublicIpsOfEc2InstancesContext(t, context.Background(), instanceIDs, awsRegion)
}

// GetEc2InstanceIdsByTagContextE returns all the IDs of EC2 instances in the given region with the given tag.
// The ctx parameter supports cancellation and timeouts.
func GetEc2InstanceIdsByTagContextE(t testing.TestingT, ctx context.Context, region string, tagName string, tagValue string) ([]string, error) {
	ec2Filters := map[string][]string{
		"tag:" + tagName: {tagValue},
	}

	return GetEc2InstanceIdsByFiltersContextE(t, ctx, region, ec2Filters)
}

// GetEc2InstanceIdsByTagContext returns all the IDs of EC2 instances in the given region with the given tag.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEc2InstanceIdsByTagContext(t testing.TestingT, ctx context.Context, region string, tagName string, tagValue string) []string {
	t.Helper()
	out, err := GetEc2InstanceIdsByTagContextE(t, ctx, region, tagName, tagValue)
	require.NoError(t, err)

	return out
}

// GetEc2InstanceIdsByTag returns all the IDs of EC2 instances in the given region with the given tag.
//
// Deprecated: Use [GetEc2InstanceIdsByTagContext] instead.
func GetEc2InstanceIdsByTag(t testing.TestingT, region string, tagName string, tagValue string) []string {
	t.Helper()
	return GetEc2InstanceIdsByTagContext(t, context.Background(), region, tagName, tagValue)
}

// GetEc2InstanceIdsByTagE returns all the IDs of EC2 instances in the given region with the given tag.
//
// Deprecated: Use [GetEc2InstanceIdsByTagContextE] instead.
func GetEc2InstanceIdsByTagE(t testing.TestingT, region string, tagName string, tagValue string) ([]string, error) {
	return GetEc2InstanceIdsByTagContextE(t, context.Background(), region, tagName, tagValue)
}

// GetEc2InstanceIdsByFiltersContextE returns all the IDs of EC2 instances in the given region which match to EC2 filter list
// as per https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput.
// The ctx parameter supports cancellation and timeouts.
func GetEc2InstanceIdsByFiltersContextE(t testing.TestingT, ctx context.Context, region string, ec2Filters map[string][]string) ([]string, error) {
	client, err := NewEc2ClientContextE(t, ctx, region)
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
		page, err := paginator.NextPage(ctx)
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

// GetEc2InstanceIdsByFiltersContext returns all the IDs of EC2 instances in the given region which match to EC2 filter list
// as per https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEc2InstanceIdsByFiltersContext(t testing.TestingT, ctx context.Context, region string, ec2Filters map[string][]string) []string {
	t.Helper()
	out, err := GetEc2InstanceIdsByFiltersContextE(t, ctx, region, ec2Filters)
	require.NoError(t, err)

	return out
}

// GetEc2InstanceIdsByFilters returns all the IDs of EC2 instances in the given region which match to EC2 filter list
// as per https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput.
//
// Deprecated: Use [GetEc2InstanceIdsByFiltersContext] instead.
func GetEc2InstanceIdsByFilters(t testing.TestingT, region string, ec2Filters map[string][]string) []string {
	t.Helper()
	return GetEc2InstanceIdsByFiltersContext(t, context.Background(), region, ec2Filters)
}

// GetEc2InstanceIdsByFiltersE returns all the IDs of EC2 instances in the given region which match to EC2 filter list
// as per https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput.
//
// Deprecated: Use [GetEc2InstanceIdsByFiltersContextE] instead.
func GetEc2InstanceIdsByFiltersE(t testing.TestingT, region string, ec2Filters map[string][]string) ([]string, error) {
	return GetEc2InstanceIdsByFiltersContextE(t, context.Background(), region, ec2Filters)
}

// GetTagsForEc2InstanceContextE returns all the tags for the given EC2 Instance.
// The ctx parameter supports cancellation and timeouts.
func GetTagsForEc2InstanceContextE(t testing.TestingT, ctx context.Context, region string, instanceID string) (map[string]string, error) {
	client, err := NewEc2ClientContextE(t, ctx, region)
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
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, tag := range page.Tags {
			tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}
	}

	return tags, nil
}

// GetTagsForEc2InstanceContext returns all the tags for the given EC2 Instance.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTagsForEc2InstanceContext(t testing.TestingT, ctx context.Context, region string, instanceID string) map[string]string {
	t.Helper()
	tags, err := GetTagsForEc2InstanceContextE(t, ctx, region, instanceID)
	require.NoError(t, err)

	return tags
}

// GetTagsForEc2Instance returns all the tags for the given EC2 Instance.
//
// Deprecated: Use [GetTagsForEc2InstanceContext] instead.
func GetTagsForEc2Instance(t testing.TestingT, region string, instanceID string) map[string]string {
	t.Helper()
	return GetTagsForEc2InstanceContext(t, context.Background(), region, instanceID)
}

// GetTagsForEc2InstanceE returns all the tags for the given EC2 Instance.
//
// Deprecated: Use [GetTagsForEc2InstanceContextE] instead.
func GetTagsForEc2InstanceE(t testing.TestingT, region string, instanceID string) (map[string]string, error) {
	return GetTagsForEc2InstanceContextE(t, context.Background(), region, instanceID)
}

// DeleteAmiContextE deletes the given AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func DeleteAmiContextE(t testing.TestingT, ctx context.Context, region string, imageID string) error {
	logger.Default.Logf(t, "Deregistering AMI %s", imageID)

	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	_, err = client.DeregisterImage(ctx, &ec2.DeregisterImageInput{ImageId: aws.String(imageID)})

	return err
}

// DeleteAmiContext deletes the given AMI in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteAmiContext(t testing.TestingT, ctx context.Context, region string, imageID string) {
	t.Helper()
	require.NoError(t, DeleteAmiContextE(t, ctx, region, imageID))
}

// DeleteAmi deletes the given AMI in the given region.
//
// Deprecated: Use [DeleteAmiContext] instead.
func DeleteAmi(t testing.TestingT, region string, imageID string) {
	t.Helper()
	DeleteAmiContext(t, context.Background(), region, imageID)
}

// DeleteAmiE deletes the given AMI in the given region.
//
// Deprecated: Use [DeleteAmiContextE] instead.
func DeleteAmiE(t testing.TestingT, region string, imageID string) error {
	return DeleteAmiContextE(t, context.Background(), region, imageID)
}

// AddTagsToResourceContextE adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
// The ctx parameter supports cancellation and timeouts.
func AddTagsToResourceContextE(t testing.TestingT, ctx context.Context, region string, resource string, tags map[string]string) error {
	client, err := NewEc2ClientContextE(t, ctx, region)
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

	_, err = client.CreateTags(ctx, &ec2.CreateTagsInput{
		Resources: []string{resource},
		Tags:      awsTags,
	})

	return err
}

// AddTagsToResourceContext adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func AddTagsToResourceContext(t testing.TestingT, ctx context.Context, region string, resource string, tags map[string]string) {
	t.Helper()
	require.NoError(t, AddTagsToResourceContextE(t, ctx, region, resource, tags))
}

// AddTagsToResource adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
//
// Deprecated: Use [AddTagsToResourceContext] instead.
func AddTagsToResource(t testing.TestingT, region string, resource string, tags map[string]string) {
	t.Helper()
	AddTagsToResourceContext(t, context.Background(), region, resource, tags)
}

// AddTagsToResourceE adds the tags to the given taggable AWS resource such as EC2, AMI or VPC.
//
// Deprecated: Use [AddTagsToResourceContextE] instead.
func AddTagsToResourceE(t testing.TestingT, region string, resource string, tags map[string]string) error {
	return AddTagsToResourceContextE(t, context.Background(), region, resource, tags)
}

// TerminateInstanceContextE terminates the EC2 instance with the given ID in the given region.
// The ctx parameter supports cancellation and timeouts.
func TerminateInstanceContextE(t testing.TestingT, ctx context.Context, region string, instanceID string) error {
	logger.Default.Logf(t, "Terminating Instance %s", instanceID)

	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	_, err = client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{
			instanceID,
		},
	})

	return err
}

// TerminateInstanceContext terminates the EC2 instance with the given ID in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func TerminateInstanceContext(t testing.TestingT, ctx context.Context, region string, instanceID string) {
	t.Helper()
	require.NoError(t, TerminateInstanceContextE(t, ctx, region, instanceID))
}

// TerminateInstance terminates the EC2 instance with the given ID in the given region.
//
// Deprecated: Use [TerminateInstanceContext] instead.
func TerminateInstance(t testing.TestingT, region string, instanceID string) {
	t.Helper()
	TerminateInstanceContext(t, context.Background(), region, instanceID)
}

// TerminateInstanceE terminates the EC2 instance with the given ID in the given region.
//
// Deprecated: Use [TerminateInstanceContextE] instead.
func TerminateInstanceE(t testing.TestingT, region string, instanceID string) error {
	return TerminateInstanceContextE(t, context.Background(), region, instanceID)
}

// GetAmiPubliclyAccessibleContextE returns whether the AMI is publicly accessible or not
// The ctx parameter supports cancellation and timeouts.
func GetAmiPubliclyAccessibleContextE(t testing.TestingT, ctx context.Context, awsRegion string, amiID string) (bool, error) {
	launchPermissions, err := GetLaunchPermissionsForAmiContextE(t, ctx, awsRegion, amiID)
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

// GetAmiPubliclyAccessibleContext returns whether the AMI is publicly accessible or not
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAmiPubliclyAccessibleContext(t testing.TestingT, ctx context.Context, awsRegion string, amiID string) bool {
	t.Helper()
	output, err := GetAmiPubliclyAccessibleContextE(t, ctx, awsRegion, amiID)
	require.NoError(t, err)

	return output
}

// GetAmiPubliclyAccessible returns whether the AMI is publicly accessible or not
//
// Deprecated: Use [GetAmiPubliclyAccessibleContext] instead.
func GetAmiPubliclyAccessible(t testing.TestingT, awsRegion string, amiID string) bool {
	t.Helper()
	return GetAmiPubliclyAccessibleContext(t, context.Background(), awsRegion, amiID)
}

// GetAmiPubliclyAccessibleE returns whether the AMI is publicly accessible or not
//
// Deprecated: Use [GetAmiPubliclyAccessibleContextE] instead.
func GetAmiPubliclyAccessibleE(t testing.TestingT, awsRegion string, amiID string) (bool, error) {
	return GetAmiPubliclyAccessibleContextE(t, context.Background(), awsRegion, amiID)
}

// GetAccountsWithLaunchPermissionsForAmiContextE returns list of accounts that the AMI is shared with
// The ctx parameter supports cancellation and timeouts.
func GetAccountsWithLaunchPermissionsForAmiContextE(t testing.TestingT, ctx context.Context, awsRegion string, amiID string) ([]string, error) {
	var accountIDs []string

	launchPermissions, err := GetLaunchPermissionsForAmiContextE(t, ctx, awsRegion, amiID)
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

// GetAccountsWithLaunchPermissionsForAmiContext returns list of accounts that the AMI is shared with
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAccountsWithLaunchPermissionsForAmiContext(t testing.TestingT, ctx context.Context, awsRegion string, amiID string) []string {
	t.Helper()
	output, err := GetAccountsWithLaunchPermissionsForAmiContextE(t, ctx, awsRegion, amiID)
	require.NoError(t, err)

	return output
}

// GetAccountsWithLaunchPermissionsForAmi returns list of accounts that the AMI is shared with
//
// Deprecated: Use [GetAccountsWithLaunchPermissionsForAmiContext] instead.
func GetAccountsWithLaunchPermissionsForAmi(t testing.TestingT, awsRegion string, amiID string) []string {
	t.Helper()
	return GetAccountsWithLaunchPermissionsForAmiContext(t, context.Background(), awsRegion, amiID)
}

// GetAccountsWithLaunchPermissionsForAmiE returns list of accounts that the AMI is shared with
//
// Deprecated: Use [GetAccountsWithLaunchPermissionsForAmiContextE] instead.
func GetAccountsWithLaunchPermissionsForAmiE(t testing.TestingT, awsRegion string, amiID string) ([]string, error) {
	return GetAccountsWithLaunchPermissionsForAmiContextE(t, context.Background(), awsRegion, amiID)
}

// GetLaunchPermissionsForAmiContextE returns launchPermissions as configured in AWS
// The ctx parameter supports cancellation and timeouts.
func GetLaunchPermissionsForAmiContextE(t testing.TestingT, ctx context.Context, awsRegion string, amiID string) ([]types.LaunchPermission, error) {
	client, err := NewEc2ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return []types.LaunchPermission{}, err
	}

	input := &ec2.DescribeImageAttributeInput{
		Attribute: types.ImageAttributeNameLaunchPermission,
		ImageId:   aws.String(amiID),
	}

	output, err := client.DescribeImageAttribute(ctx, input)
	if err != nil {
		return []types.LaunchPermission{}, err
	}

	return output.LaunchPermissions, nil
}

// GetLaunchPermissionsForAmiContext returns launchPermissions as configured in AWS.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLaunchPermissionsForAmiContext(t testing.TestingT, ctx context.Context, awsRegion string, amiID string) []types.LaunchPermission {
	t.Helper()
	output, err := GetLaunchPermissionsForAmiContextE(t, ctx, awsRegion, amiID)
	require.NoError(t, err)

	return output
}

// GetLaunchPermissionsForAmi returns launchPermissions as configured in AWS.
//
// Deprecated: Use [GetLaunchPermissionsForAmiContext] instead.
func GetLaunchPermissionsForAmi(t testing.TestingT, awsRegion string, amiID string) []types.LaunchPermission {
	t.Helper()
	return GetLaunchPermissionsForAmiContext(t, context.Background(), awsRegion, amiID)
}

// GetLaunchPermissionsForAmiE returns launchPermissions as configured in AWS
//
// Deprecated: Use [GetLaunchPermissionsForAmiContextE] instead.
func GetLaunchPermissionsForAmiE(t testing.TestingT, awsRegion string, amiID string) ([]types.LaunchPermission, error) {
	return GetLaunchPermissionsForAmiContextE(t, context.Background(), awsRegion, amiID)
}

// GetRecommendedInstanceTypeContextE takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedInstanceTypeContextE(t testing.TestingT, ctx context.Context, region string, instanceTypeOptions []string) (string, error) {
	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	return GetRecommendedInstanceTypeWithClientContextE(t, ctx, client, instanceTypeOptions)
}

// GetRecommendedInstanceTypeContext takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs, and if you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedInstanceTypeContext(t testing.TestingT, ctx context.Context, region string, instanceTypeOptions []string) string {
	t.Helper()
	out, err := GetRecommendedInstanceTypeContextE(t, ctx, region, instanceTypeOptions)
	require.NoError(t, err)

	return out
}

// GetRecommendedInstanceType takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs, and if you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function will fail the test if there is an error.
//
// Deprecated: Use [GetRecommendedInstanceTypeContext] instead.
func GetRecommendedInstanceType(t testing.TestingT, region string, instanceTypeOptions []string) string {
	t.Helper()
	return GetRecommendedInstanceTypeContext(t, context.Background(), region, instanceTypeOptions)
}

// GetRecommendedInstanceTypeE takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
//
// Deprecated: Use [GetRecommendedInstanceTypeContextE] instead.
func GetRecommendedInstanceTypeE(t testing.TestingT, region string, instanceTypeOptions []string) (string, error) {
	return GetRecommendedInstanceTypeContextE(t, context.Background(), region, instanceTypeOptions)
}

// GetRecommendedInstanceTypeWithClientContextE takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function expects an authenticated EC2 client from the AWS SDK Go library.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedInstanceTypeWithClientContextE(t testing.TestingT, ctx context.Context, ec2Client *ec2.Client, instanceTypeOptions []string) (string, error) {
	availabilityZones, err := getAllAvailabilityZonesContextE(ctx, ec2Client)
	if err != nil {
		return "", err
	}

	instanceTypeOfferings, err := getInstanceTypeOfferingsContextE(ctx, ec2Client, instanceTypeOptions)
	if err != nil {
		return "", err
	}

	return PickRecommendedInstanceTypeE(availabilityZones, instanceTypeOfferings, instanceTypeOptions)
}

// GetRecommendedInstanceTypeWithClientContext takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function expects an authenticated EC2 client from the AWS SDK Go library.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecommendedInstanceTypeWithClientContext(t testing.TestingT, ctx context.Context, ec2Client *ec2.Client, instanceTypeOptions []string) string {
	t.Helper()
	out, err := GetRecommendedInstanceTypeWithClientContextE(t, ctx, ec2Client, instanceTypeOptions)
	require.NoError(t, err)

	return out
}

// GetRecommendedInstanceTypeWithClient takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function expects an authenticated EC2 client from the AWS SDK Go library.
// This function will fail the test if there is an error.
//
// Deprecated: Use [GetRecommendedInstanceTypeWithClientContext] instead.
func GetRecommendedInstanceTypeWithClient(t testing.TestingT, ec2Client *ec2.Client, instanceTypeOptions []string) string {
	t.Helper()
	return GetRecommendedInstanceTypeWithClientContext(t, context.Background(), ec2Client, instanceTypeOptions)
}

// GetRecommendedInstanceTypeWithClientE takes in a list of EC2 instance types (e.g., "t2.micro", "t3.micro") and returns the
// first instance type in the list that is available in all Availability Zones (AZs) in the given region. If there's no
// instance available in all AZs, this function exits with an error. This is useful because certain instance types,
// such as t2.micro, are not available in some of the newer AZs, while t3.micro is not available in some of the older
// AZs. If you have code that needs to run on a "small" instance across all AZs in many different regions, you can
// use this function to automatically figure out which instance type you should use.
// This function expects an authenticated EC2 client from the AWS SDK Go library.
//
// Deprecated: Use [GetRecommendedInstanceTypeWithClientContextE] instead.
func GetRecommendedInstanceTypeWithClientE(t testing.TestingT, ec2Client *ec2.Client, instanceTypeOptions []string) (string, error) {
	return GetRecommendedInstanceTypeWithClientContextE(t, context.Background(), ec2Client, instanceTypeOptions)
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

// getInstanceTypeOfferingsContextE returns the instance types from the given list that are available in the region configured
// in the given EC2 client
func getInstanceTypeOfferingsContextE(ctx context.Context, client *ec2.Client, instanceTypeOptions []string) ([]types.InstanceTypeOffering, error) {
	input := ec2.DescribeInstanceTypeOfferingsInput{
		LocationType: types.LocationTypeAvailabilityZone,
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: instanceTypeOptions,
			},
		},
	}

	out, err := client.DescribeInstanceTypeOfferings(ctx, &input)
	if err != nil {
		return nil, err
	}

	return out.InstanceTypeOfferings, nil
}

// getAllAvailabilityZonesContextE returns all the available AZs in the region configured in the given EC2 client
func getAllAvailabilityZonesContextE(ctx context.Context, client *ec2.Client) ([]string, error) {
	input := ec2.DescribeAvailabilityZonesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("state"),
				Values: []string{"available"},
			},
		},
	}

	out, err := client.DescribeAvailabilityZones(ctx, &input)
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

// getInstanceFieldMapContextE is a shared helper that paginates through DescribeInstances and builds a map
// of instance ID to a string field extracted by the given function.
func getInstanceFieldMapContextE(t testing.TestingT, ctx context.Context, instanceIDs []string, awsRegion string, extractField func(*types.Instance) string) (map[string]string, error) {
	ec2Client, err := NewEc2ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	input := ec2.DescribeInstancesInput{InstanceIds: instanceIDs}

	result := map[string]string{}

	paginator := ec2.NewDescribeInstancesPaginator(ec2Client, &input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
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

// NewEc2ClientContextE creates an EC2 client.
// The ctx parameter supports cancellation and timeouts.
func NewEc2ClientContextE(t testing.TestingT, ctx context.Context, region string) (*ec2.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return ec2.NewFromConfig(*sess), nil
}

// NewEc2Client creates an EC2 client.
//
// Deprecated: Use [NewEc2ClientContext] instead.
func NewEc2Client(t testing.TestingT, region string) *ec2.Client {
	t.Helper()

	return NewEc2ClientContext(t, context.Background(), region)
}

// NewEc2ClientE creates an EC2 client.
//
// Deprecated: Use [NewEc2ClientContextE] instead.
func NewEc2ClientE(t testing.TestingT, region string) (*ec2.Client, error) {
	return NewEc2ClientContextE(t, context.Background(), region)
}

// NewEc2ClientContext creates an EC2 client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewEc2ClientContext(t testing.TestingT, ctx context.Context, region string) *ec2.Client {
	t.Helper()
	client, err := NewEc2ClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}
