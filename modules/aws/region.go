package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// You can set this environment variable to force Terratest to use a specific region rather than a random one. This is
// convenient when iterating locally.
const regionOverrideEnvVarName = "TERRATEST_REGION"

// AWS API calls typically require an AWS region. We typically require the user to set one explicitly, but in some
// cases, this doesn't make sense (e.g., for fetching the list of regions in an account), so for those cases, we use
// this region as a default.
const defaultRegion = "us-east-1"

// Reference for launch dates: https://aws.amazon.com/about-aws/global-infrastructure/
var stableRegions = []string{
	"us-east-1",      // Launched 2006
	"us-east-2",      // Launched 2016
	"us-west-1",      // Launched 2009
	"us-west-2",      // Launched 2011
	"ca-central-1",   // Launched 2016
	"sa-east-1",      // Launched 2011
	"eu-west-1",      // Launched 2007
	"eu-west-2",      // Launched 2016
	"eu-west-3",      // Launched 2017
	"eu-central-1",   // Launched 2014
	"ap-southeast-1", // Launched 2010
	"ap-southeast-2", // Launched 2012
	"ap-northeast-1", // Launched 2011
	"ap-northeast-2", // Launched 2016
	"ap-south-1",     // Launched 2016
	"eu-north-1",     // Launched 2018
}

// GetRandomStableRegionContextE gets a randomly chosen AWS region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
// The ctx parameter supports cancellation and timeouts.
func GetRandomStableRegionContextE(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string) (string, error) {
	regionsToPickFrom := stableRegions

	if len(approvedRegions) > 0 {
		regionsToPickFrom = collections.ListIntersection(regionsToPickFrom, approvedRegions)
	}

	if len(forbiddenRegions) > 0 {
		regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	}

	return GetRandomRegionContextE(t, ctx, regionsToPickFrom, nil)
}

// GetRandomStableRegionContext gets a randomly chosen AWS region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
// The ctx parameter supports cancellation and timeouts.
func GetRandomStableRegionContext(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string) string {
	t.Helper()

	region, err := GetRandomStableRegionContextE(t, ctx, approvedRegions, forbiddenRegions)
	require.NoError(t, err)

	return region
}

// GetRandomStableRegion gets a randomly chosen AWS region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
//
// Deprecated: Use [GetRandomStableRegionContext] instead.
func GetRandomStableRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) string {
	t.Helper()

	return GetRandomStableRegionContext(t, context.Background(), approvedRegions, forbiddenRegions)
}

// GetRandomStableRegionE gets a randomly chosen AWS region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
//
// Deprecated: Use [GetRandomStableRegionContextE] instead.
func GetRandomStableRegionE(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) (string, error) {
	return GetRandomStableRegionContextE(t, context.Background(), approvedRegions, forbiddenRegions)
}

// GetRandomRegionContextE gets a randomly chosen AWS region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the AWS APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionContextE(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string) (string, error) {
	regionFromEnvVar := os.Getenv(regionOverrideEnvVarName)
	if regionFromEnvVar != "" {
		logger.Default.Logf(t, "Using AWS region %s from environment variable %s", regionFromEnvVar, regionOverrideEnvVarName)

		return regionFromEnvVar, nil
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllAwsRegionsContextE(t, ctx)
		if err != nil {
			return "", err
		}

		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	logger.Default.Logf(t, "Using region %s", region)

	return region, nil
}

// GetRandomRegionContext gets a randomly chosen AWS region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the AWS APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionContext(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string) string {
	t.Helper()

	region, err := GetRandomRegionContextE(t, ctx, approvedRegions, forbiddenRegions)
	require.NoError(t, err)

	return region
}

// GetRandomRegion gets a randomly chosen AWS region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the AWS APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
//
// Deprecated: Use [GetRandomRegionContext] instead.
func GetRandomRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) string {
	t.Helper()

	return GetRandomRegionContext(t, context.Background(), approvedRegions, forbiddenRegions)
}

// GetRandomRegionE gets a randomly chosen AWS region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the AWS APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
//
// Deprecated: Use [GetRandomRegionContextE] instead.
func GetRandomRegionE(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) (string, error) {
	return GetRandomRegionContextE(t, context.Background(), approvedRegions, forbiddenRegions)
}

// GetAllAwsRegionsContextE gets the list of AWS regions available in this account.
// The ctx parameter supports cancellation and timeouts.
func GetAllAwsRegionsContextE(t testing.TestingT, ctx context.Context) ([]string, error) {
	logger.Default.Logf(t, "Looking up all AWS regions available in this account")

	ec2Client, err := NewEc2ClientContextE(t, ctx, defaultRegion)
	if err != nil {
		return nil, err
	}

	out, err := ec2Client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	var regions []string

	for _, region := range out.Regions {
		regions = append(regions, aws.ToString(region.RegionName))
	}

	return regions, nil
}

// GetAllAwsRegionsContext gets the list of AWS regions available in this account.
// The ctx parameter supports cancellation and timeouts.
func GetAllAwsRegionsContext(t testing.TestingT, ctx context.Context) []string {
	t.Helper()

	out, err := GetAllAwsRegionsContextE(t, ctx)
	require.NoError(t, err)

	return out
}

// GetAllAwsRegions gets the list of AWS regions available in this account.
//
// Deprecated: Use [GetAllAwsRegionsContext] instead.
func GetAllAwsRegions(t testing.TestingT) []string {
	t.Helper()

	return GetAllAwsRegionsContext(t, context.Background())
}

// GetAllAwsRegionsE gets the list of AWS regions available in this account.
//
// Deprecated: Use [GetAllAwsRegionsContextE] instead.
func GetAllAwsRegionsE(t testing.TestingT) ([]string, error) {
	return GetAllAwsRegionsContextE(t, context.Background())
}

// GetAvailabilityZonesContextE gets the Availability Zones for a given AWS region. Note that for certain regions (e.g. us-east-1), different AWS
// accounts have access to different availability zones.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilityZonesContextE(t testing.TestingT, ctx context.Context, region string) ([]string, error) {
	logger.Default.Logf(t, "Looking up all availability zones available in this account for region %s", region)

	ec2Client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	resp, err := ec2Client.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, err
	}

	var out []string

	for i := range resp.AvailabilityZones {
		out = append(out, aws.ToString(resp.AvailabilityZones[i].ZoneName))
	}

	return out, nil
}

// GetAvailabilityZonesContext gets the Availability Zones for a given AWS region. Note that for certain regions (e.g. us-east-1), different AWS
// accounts have access to different availability zones.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilityZonesContext(t testing.TestingT, ctx context.Context, region string) []string {
	t.Helper()

	out, err := GetAvailabilityZonesContextE(t, ctx, region)
	require.NoError(t, err)

	return out
}

// GetAvailabilityZones gets the Availability Zones for a given AWS region. Note that for certain regions (e.g. us-east-1), different AWS
// accounts have access to different availability zones.
//
// Deprecated: Use [GetAvailabilityZonesContext] instead.
func GetAvailabilityZones(t testing.TestingT, region string) []string {
	t.Helper()

	return GetAvailabilityZonesContext(t, context.Background(), region)
}

// GetAvailabilityZonesE gets the Availability Zones for a given AWS region. Note that for certain regions (e.g. us-east-1), different AWS
// accounts have access to different availability zones.
//
// Deprecated: Use [GetAvailabilityZonesContextE] instead.
func GetAvailabilityZonesE(t testing.TestingT, region string) ([]string, error) {
	return GetAvailabilityZonesContextE(t, context.Background(), region)
}

// GetRegionsForServiceContextE gets all AWS regions in which a service is available and returns errors.
// See https://docs.aws.amazon.com/systems-manager/latest/userguide/parameter-store-public-parameters-global-infrastructure.html
// The ctx parameter supports cancellation and timeouts.
func GetRegionsForServiceContextE(t testing.TestingT, ctx context.Context, serviceName string) ([]string, error) {
	// These values are available in any region, defaulting to us-east-1 since it's the oldest
	ssmClient, err := NewSsmClientContextE(t, ctx, "us-east-1")
	if err != nil {
		return nil, err
	}

	paramPath := "/aws/service/global-infrastructure/services/%s/regions"

	resp, err := ssmClient.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
		Path: aws.String(fmt.Sprintf(paramPath, serviceName)),
	})
	if err != nil {
		return nil, err
	}

	var availableRegions []string

	for _, p := range resp.Parameters {
		availableRegions = append(availableRegions, *p.Value)
	}

	return availableRegions, nil
}

// GetRegionsForServiceContext gets all AWS regions in which a service is available.
// The ctx parameter supports cancellation and timeouts.
func GetRegionsForServiceContext(t testing.TestingT, ctx context.Context, serviceName string) []string {
	t.Helper()

	out, err := GetRegionsForServiceContextE(t, ctx, serviceName)
	require.NoError(t, err)

	return out
}

// GetRegionsForService gets all AWS regions in which a service is available.
//
// Deprecated: Use [GetRegionsForServiceContext] instead.
func GetRegionsForService(t testing.TestingT, serviceName string) []string {
	t.Helper()

	return GetRegionsForServiceContext(t, context.Background(), serviceName)
}

// GetRegionsForServiceE gets all AWS regions in which a service is available and returns errors.
// See https://docs.aws.amazon.com/systems-manager/latest/userguide/parameter-store-public-parameters-global-infrastructure.html
//
// Deprecated: Use [GetRegionsForServiceContextE] instead.
func GetRegionsForServiceE(t testing.TestingT, serviceName string) ([]string, error) {
	return GetRegionsForServiceContextE(t, context.Background(), serviceName)
}

// GetRandomRegionForServiceContextE retrieves a list of AWS regions in which a service is available
// Then returns one region randomly from the list.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionForServiceContextE(t testing.TestingT, ctx context.Context, serviceName string) (string, error) {
	availableRegions, err := GetRegionsForServiceContextE(t, ctx, serviceName)
	if err != nil {
		return "", err
	}

	return GetRandomRegionContextE(t, ctx, availableRegions, nil)
}

// GetRandomRegionForServiceContext retrieves a list of AWS regions in which a service is available
// Then returns one region randomly from the list.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionForServiceContext(t testing.TestingT, ctx context.Context, serviceName string) string {
	t.Helper()

	region, err := GetRandomRegionForServiceContextE(t, ctx, serviceName)
	require.NoError(t, err)

	return region
}

// GetRandomRegionForService retrieves a list of AWS regions in which a service is available
// Then returns one region randomly from the list
//
// Deprecated: Use [GetRandomRegionForServiceContext] instead.
func GetRandomRegionForService(t testing.TestingT, serviceName string) string {
	t.Helper()

	return GetRandomRegionForServiceContext(t, context.Background(), serviceName)
}

// GetRandomRegionForServiceE retrieves a list of AWS regions in which a service is available
// Then returns one region randomly from the list and returns errors.
//
// Deprecated: Use [GetRandomRegionForServiceContextE] instead.
func GetRandomRegionForServiceE(t testing.TestingT, serviceName string) (string, error) {
	return GetRandomRegionForServiceContextE(t, context.Background(), serviceName)
}
