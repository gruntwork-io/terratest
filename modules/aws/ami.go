package aws

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// These are commonly used AMI account IDs.
const (
	// CanonicalAccountID is the AWS account ID for Canonical (Ubuntu).
	CanonicalAccountID = "099720109477"
	// CentOsAccountID is the AWS account ID for CentOS.
	CentOsAccountID = "679593333241"
	// AmazonAccountID is the AWS account ID (or alias) for Amazon.
	AmazonAccountID = "amazon"

	// Deprecated: Use [CanonicalAccountID] instead.
	CanonicalAccountId = CanonicalAccountID //nolint:staticcheck,revive // preserving deprecated constant name
	// Deprecated: Use [CentOsAccountID] instead.
	CentOsAccountId = CentOsAccountID //nolint:staticcheck,revive // preserving deprecated constant name
	// Deprecated: Use [AmazonAccountID] instead.
	AmazonAccountId = AmazonAccountID //nolint:staticcheck,revive // preserving deprecated constant name
)

// DeleteAmiAndAllSnapshotsContextE will delete the given AMI along with all EBS snapshots that backed that AMI.
// The ctx parameter supports cancellation and timeouts.
func DeleteAmiAndAllSnapshotsContextE(t testing.TestingT, ctx context.Context, region string, ami string) error {
	snapshots, err := GetEbsSnapshotsForAmiContextE(t, ctx, region, ami)
	if err != nil {
		return err
	}

	err = DeleteAmiContextE(t, ctx, region, ami)
	if err != nil {
		return err
	}

	for _, snapshot := range snapshots {
		err = DeleteEbsSnapshotContextE(t, ctx, region, snapshot)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteAmiAndAllSnapshotsContext will delete the given AMI along with all EBS snapshots that backed that AMI.
// The ctx parameter supports cancellation and timeouts.
func DeleteAmiAndAllSnapshotsContext(t testing.TestingT, ctx context.Context, region string, ami string) {
	t.Helper()

	err := DeleteAmiAndAllSnapshotsContextE(t, ctx, region, ami)
	require.NoError(t, err)
}

// DeleteAmiAndAllSnapshots will delete the given AMI along with all EBS snapshots that backed that AMI.
//
// Deprecated: Use [DeleteAmiAndAllSnapshotsContext] instead.
func DeleteAmiAndAllSnapshots(t testing.TestingT, region string, ami string) {
	t.Helper()

	DeleteAmiAndAllSnapshotsContext(t, context.Background(), region, ami)
}

// DeleteAmiAndAllSnapshotsE will delete the given AMI along with all EBS snapshots that backed that AMI.
//
// Deprecated: Use [DeleteAmiAndAllSnapshotsContextE] instead.
func DeleteAmiAndAllSnapshotsE(t testing.TestingT, region string, ami string) error {
	return DeleteAmiAndAllSnapshotsContextE(t, context.Background(), region, ami)
}

// GetEbsSnapshotsForAmiContextE retrieves the EBS snapshots which back the given AMI.
// The ctx parameter supports cancellation and timeouts.
func GetEbsSnapshotsForAmiContextE(t testing.TestingT, ctx context.Context, region string, ami string) ([]string, error) {
	logger.Default.Logf(t, "Retrieving EBS snapshots backing AMI %s", ami)

	ec2Client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	images, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{
			ami,
		},
	})
	if err != nil {
		return nil, err
	}

	var snapshots []string

	for i := range images.Images {
		image := &images.Images[i]

		for _, mapping := range image.BlockDeviceMappings {
			if mapping.Ebs != nil && mapping.Ebs.SnapshotId != nil {
				snapshots = append(snapshots, aws.ToString(mapping.Ebs.SnapshotId))
			}
		}
	}

	return snapshots, nil
}

// GetEbsSnapshotsForAmiContext retrieves the EBS snapshots which back the given AMI.
// The ctx parameter supports cancellation and timeouts.
func GetEbsSnapshotsForAmiContext(t testing.TestingT, ctx context.Context, region string, ami string) []string {
	t.Helper()

	snapshots, err := GetEbsSnapshotsForAmiContextE(t, ctx, region, ami)
	require.NoError(t, err)

	return snapshots
}

// GetEbsSnapshotsForAmi retrieves the EBS snapshots which back the given AMI.
//
// Deprecated: Use [GetEbsSnapshotsForAmiContext] instead.
func GetEbsSnapshotsForAmi(t testing.TestingT, region string, ami string) []string {
	t.Helper()

	return GetEbsSnapshotsForAmiContext(t, context.Background(), region, ami)
}

// GetEbsSnapshotsForAmiE retrieves the EBS snapshots which back the given AMI.
//
// Deprecated: Use [GetEbsSnapshotsForAmiContextE] instead.
func GetEbsSnapshotsForAmiE(t testing.TestingT, region string, ami string) ([]string, error) {
	return GetEbsSnapshotsForAmiContextE(t, context.Background(), region, ami)
}

// GetMostRecentAmiIDContextE gets the ID of the most recent AMI in the given region that has the given owner and matches
// the given filters. Each filter should correspond to the name and values of a filter supported by DescribeImagesInput:
// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeImagesInput
// The ctx parameter supports cancellation and timeouts.
func GetMostRecentAmiIDContextE(t testing.TestingT, ctx context.Context, region string, ownerID string, filters map[string][]string) (string, error) {
	ec2Client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	var ec2Filters []types.Filter

	for name, values := range filters {
		ec2Filters = append(ec2Filters, types.Filter{Name: aws.String(name), Values: values})
	}

	input := ec2.DescribeImagesInput{
		Filters:           ec2Filters,
		IncludeDeprecated: aws.Bool(true),
		Owners:            []string{ownerID},
	}

	out, err := ec2Client.DescribeImages(ctx, &input)
	if err != nil {
		return "", err
	}

	if len(out.Images) == 0 {
		return "", NoImagesFound{Filters: filters, Region: region, OwnerID: ownerID}
	}

	mostRecentImage := mostRecentAMI(out.Images)

	return aws.ToString(mostRecentImage.ImageId), nil
}

// GetMostRecentAmiIDContext gets the ID of the most recent AMI in the given region that has the given owner and matches
// the given filters. Each filter should correspond to the name and values of a filter supported by DescribeImagesInput:
// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeImagesInput
// The ctx parameter supports cancellation and timeouts.
func GetMostRecentAmiIDContext(t testing.TestingT, ctx context.Context, region string, ownerID string, filters map[string][]string) string {
	t.Helper()

	amiID, err := GetMostRecentAmiIDContextE(t, ctx, region, ownerID, filters)
	require.NoError(t, err)

	return amiID
}

// GetMostRecentAmiID gets the ID of the most recent AMI in the given region that has the given owner and matches
// the given filters. Each filter should correspond to the name and values of a filter supported by DescribeImagesInput:
// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeImagesInput
//
// Deprecated: Use [GetMostRecentAmiIDContext] instead.
func GetMostRecentAmiID(t testing.TestingT, region string, ownerID string, filters map[string][]string) string {
	t.Helper()

	return GetMostRecentAmiIDContext(t, context.Background(), region, ownerID, filters)
}

// GetMostRecentAmiIDE gets the ID of the most recent AMI in the given region that has the given owner and matches
// the given filters. Each filter should correspond to the name and values of a filter supported by DescribeImagesInput:
// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeImagesInput
//
// Deprecated: Use [GetMostRecentAmiIDContextE] instead.
func GetMostRecentAmiIDE(t testing.TestingT, region string, ownerID string, filters map[string][]string) (string, error) {
	return GetMostRecentAmiIDContextE(t, context.Background(), region, ownerID, filters)
}

// GetMostRecentAmiId gets the ID of the most recent AMI in the given region that has the given owner and matches
// the given filters.
//
// Deprecated: Use [GetMostRecentAmiID] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetMostRecentAmiId(t testing.TestingT, region string, ownerId string, filters map[string][]string) string {
	return GetMostRecentAmiID(t, region, ownerId, filters)
}

// GetMostRecentAmiIdE gets the ID of the most recent AMI in the given region that has the given owner and matches
// the given filters.
//
// Deprecated: Use [GetMostRecentAmiIDE] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetMostRecentAmiIdE(t testing.TestingT, region string, ownerId string, filters map[string][]string) (string, error) {
	return GetMostRecentAmiIDE(t, region, ownerId, filters)
}

// Image sorting code borrowed from: https://github.com/hashicorp/packer/blob/7f4112ba229309cfc0ebaa10ded2abdfaf1b22c8/builder/amazon/common/step_source_ami_info.go
type imageSort []types.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	iTime, _ := time.Parse(time.RFC3339, *a[i].CreationDate)
	jTime, _ := time.Parse(time.RFC3339, *a[j].CreationDate)

	return iTime.Unix() < jTime.Unix()
}

// mostRecentAMI returns the most recent AMI out of a slice of images.
func mostRecentAMI(images []types.Image) types.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))

	return sortedImages[len(sortedImages)-1]
}

// GetUbuntu1404AmiContextE gets the ID of the most recent Ubuntu 14.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu1404AmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*ubuntu-trusty-14.04-amd64-server-*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, CanonicalAccountID, filters)
}

// GetUbuntu1404AmiContext gets the ID of the most recent Ubuntu 14.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu1404AmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetUbuntu1404AmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetUbuntu1404Ami gets the ID of the most recent Ubuntu 14.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu1404AmiContext] instead.
func GetUbuntu1404Ami(t testing.TestingT, region string) string {
	t.Helper()

	return GetUbuntu1404AmiContext(t, context.Background(), region)
}

// GetUbuntu1404AmiE gets the ID of the most recent Ubuntu 14.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu1404AmiContextE] instead.
func GetUbuntu1404AmiE(t testing.TestingT, region string) (string, error) {
	return GetUbuntu1404AmiContextE(t, context.Background(), region)
}

// GetUbuntu1604AmiContextE gets the ID of the most recent Ubuntu 16.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu1604AmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*ubuntu-xenial-16.04-amd64-server-*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, CanonicalAccountID, filters)
}

// GetUbuntu1604AmiContext gets the ID of the most recent Ubuntu 16.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu1604AmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetUbuntu1604AmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetUbuntu1604Ami gets the ID of the most recent Ubuntu 16.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu1604AmiContext] instead.
func GetUbuntu1604Ami(t testing.TestingT, region string) string {
	t.Helper()

	return GetUbuntu1604AmiContext(t, context.Background(), region)
}

// GetUbuntu1604AmiE gets the ID of the most recent Ubuntu 16.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu1604AmiContextE] instead.
func GetUbuntu1604AmiE(t testing.TestingT, region string) (string, error) {
	return GetUbuntu1604AmiContextE(t, context.Background(), region)
}

// GetUbuntu2004AmiContextE gets the ID of the most recent Ubuntu 20.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu2004AmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*ubuntu-focal-20.04-amd64-server-*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, CanonicalAccountID, filters)
}

// GetUbuntu2004AmiContext gets the ID of the most recent Ubuntu 20.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu2004AmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetUbuntu2004AmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetUbuntu2004Ami gets the ID of the most recent Ubuntu 20.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu2004AmiContext] instead.
func GetUbuntu2004Ami(t testing.TestingT, region string) string {
	t.Helper()

	return GetUbuntu2004AmiContext(t, context.Background(), region)
}

// GetUbuntu2004AmiE gets the ID of the most recent Ubuntu 20.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu2004AmiContextE] instead.
func GetUbuntu2004AmiE(t testing.TestingT, region string) (string, error) {
	return GetUbuntu2004AmiContextE(t, context.Background(), region)
}

// GetUbuntu2204AmiContextE gets the ID of the most recent Ubuntu 22.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu2204AmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*ubuntu-jammy-22.04-amd64-server-*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, CanonicalAccountID, filters)
}

// GetUbuntu2204AmiContext gets the ID of the most recent Ubuntu 22.04 HVM x86_64 EBS GP2 AMI in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetUbuntu2204AmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetUbuntu2204AmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetUbuntu2204Ami gets the ID of the most recent Ubuntu 22.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu2204AmiContext] instead.
func GetUbuntu2204Ami(t testing.TestingT, region string) string {
	t.Helper()

	return GetUbuntu2204AmiContext(t, context.Background(), region)
}

// GetUbuntu2204AmiE gets the ID of the most recent Ubuntu 22.04 HVM x86_64 EBS GP2 AMI in the given region.
//
// Deprecated: Use [GetUbuntu2204AmiContextE] instead.
func GetUbuntu2204AmiE(t testing.TestingT, region string) (string, error) {
	return GetUbuntu2204AmiContextE(t, context.Background(), region)
}

// GetCentos7AmiContextE returns a CentOS 7 public AMI from the given region.
// WARNING: you may have to accept the terms & conditions of this AMI in AWS MarketPlace for your AWS Account before
// you can successfully launch the AMI.
// The ctx parameter supports cancellation and timeouts.
func GetCentos7AmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*CentOS Linux 7 x86_64 HVM EBS*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, CentOsAccountID, filters)
}

// GetCentos7AmiContext returns a CentOS 7 public AMI from the given region.
// WARNING: you may have to accept the terms & conditions of this AMI in AWS MarketPlace for your AWS Account before
// you can successfully launch the AMI.
// The ctx parameter supports cancellation and timeouts.
func GetCentos7AmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetCentos7AmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetCentos7Ami returns a CentOS 7 public AMI from the given region.
// WARNING: you may have to accept the terms & conditions of this AMI in AWS MarketPlace for your AWS Account before
// you can successfully launch the AMI.
//
// Deprecated: Use [GetCentos7AmiContext] instead.
func GetCentos7Ami(t testing.TestingT, region string) string {
	t.Helper()

	return GetCentos7AmiContext(t, context.Background(), region)
}

// GetCentos7AmiE returns a CentOS 7 public AMI from the given region.
// WARNING: you may have to accept the terms & conditions of this AMI in AWS MarketPlace for your AWS Account before
// you can successfully launch the AMI.
//
// Deprecated: Use [GetCentos7AmiContextE] instead.
func GetCentos7AmiE(t testing.TestingT, region string) (string, error) {
	return GetCentos7AmiContextE(t, context.Background(), region)
}

// GetAmazonLinuxAmiContextE returns an Amazon Linux AMI HVM, SSD Volume Type public AMI for the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAmazonLinuxAmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*amzn2-ami-hvm-*-x86_64*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, AmazonAccountID, filters)
}

// GetAmazonLinuxAmiContext returns an Amazon Linux AMI HVM, SSD Volume Type public AMI for the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAmazonLinuxAmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetAmazonLinuxAmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetAmazonLinuxAmi returns an Amazon Linux AMI HVM, SSD Volume Type public AMI for the given region.
//
// Deprecated: Use [GetAmazonLinuxAmiContext] instead.
func GetAmazonLinuxAmi(t testing.TestingT, region string) string {
	t.Helper()

	return GetAmazonLinuxAmiContext(t, context.Background(), region)
}

// GetAmazonLinuxAmiE returns an Amazon Linux AMI HVM, SSD Volume Type public AMI for the given region.
//
// Deprecated: Use [GetAmazonLinuxAmiContextE] instead.
func GetAmazonLinuxAmiE(t testing.TestingT, region string) (string, error) {
	return GetAmazonLinuxAmiContextE(t, context.Background(), region)
}

// GetEcsOptimizedAmazonLinuxAmiContextE returns an Amazon ECS-Optimized Amazon Linux AMI for the given region. This AMI is useful for running an ECS cluster.
// The ctx parameter supports cancellation and timeouts.
func GetEcsOptimizedAmazonLinuxAmiContextE(t testing.TestingT, ctx context.Context, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*amzn-ami*amazon-ecs-optimized*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIDContextE(t, ctx, region, AmazonAccountID, filters)
}

// GetEcsOptimizedAmazonLinuxAmiContext returns an Amazon ECS-Optimized Amazon Linux AMI for the given region. This AMI is useful for running an ECS cluster.
// The ctx parameter supports cancellation and timeouts.
func GetEcsOptimizedAmazonLinuxAmiContext(t testing.TestingT, ctx context.Context, region string) string {
	t.Helper()

	amiID, err := GetEcsOptimizedAmazonLinuxAmiContextE(t, ctx, region)
	require.NoError(t, err)

	return amiID
}

// GetEcsOptimizedAmazonLinuxAmi returns an Amazon ECS-Optimized Amazon Linux AMI for the given region. This AMI is useful for running an ECS cluster.
//
// Deprecated: Use [GetEcsOptimizedAmazonLinuxAmiContext] instead.
func GetEcsOptimizedAmazonLinuxAmi(t testing.TestingT, region string) string {
	t.Helper()

	return GetEcsOptimizedAmazonLinuxAmiContext(t, context.Background(), region)
}

// GetEcsOptimizedAmazonLinuxAmiE returns an Amazon ECS-Optimized Amazon Linux AMI for the given region. This AMI is useful for running an ECS cluster.
//
// Deprecated: Use [GetEcsOptimizedAmazonLinuxAmiContextE] instead.
func GetEcsOptimizedAmazonLinuxAmiE(t testing.TestingT, region string) (string, error) {
	return GetEcsOptimizedAmazonLinuxAmiContextE(t, context.Background(), region)
}

// NoImagesFound is an error that occurs if no images were found.
type NoImagesFound struct {
	Filters map[string][]string
	Region  string
	OwnerID string //nolint:staticcheck,revive // preserving existing field name
}

func (err NoImagesFound) Error() string {
	return fmt.Sprintf("No AMIs found in %s for owner ID %s and filters: %v", err.Region, err.OwnerID, err.Filters)
}
