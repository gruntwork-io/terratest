package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// You can set this environment variable to force Terratest to use a specific region rather than a random one. This is
// convenient when iterating locally.
const regionOverrideEnvVarName = "TERRATEST_REGION"

// AWS API calls typically require an AWS region. We typically require the user to set one explicitly, but in some
// cases, this doesn't make sense (e.g., for fetching the lsit of regions in an account), so for those cases, we use
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

// GetRandomStableRegion gets a randomly chosen AWS region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
func GetRandomStableRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) string {
	regionsToPickFrom := stableRegions
	if len(approvedRegions) > 0 {
		regionsToPickFrom = collections.ListIntersection(regionsToPickFrom, approvedRegions)
	}
	if len(forbiddenRegions) > 0 {
		regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	}
	return GetRandomRegion(t, regionsToPickFrom, nil)
}

// GetRandomRegion gets a randomly chosen AWS region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the AWS APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
func GetRandomRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) string {
	region, err := GetRandomRegionE(t, approvedRegions, forbiddenRegions)
	if err != nil {
		t.Fatal(err)
	}
	return region
}

// GetRandomRegionE gets a randomly chosen AWS region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the AWS APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
func GetRandomRegionE(t testing.TestingT, approvedRegions []string, forbiddenRegions []string) (string, error) {
	regionFromEnvVar := os.Getenv(regionOverrideEnvVarName)
	if regionFromEnvVar != "" {
		logger.Logf(t, "Using AWS region %s from environment variable %s", regionFromEnvVar, regionOverrideEnvVarName)
		return regionFromEnvVar, nil
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllAwsRegionsE(t)
		if err != nil {
			return "", err
		}
		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	logger.Logf(t, "Using region %s", region)
	return region, nil
}

// GetAllAwsRegions gets the list of AWS regions available in this account.
func GetAllAwsRegions(t testing.TestingT) []string {
	out, err := GetAllAwsRegionsE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetAllAwsRegionsE gets the list of AWS regions available in this account.
func GetAllAwsRegionsE(t testing.TestingT) ([]string, error) {
	logger.Log(t, "Looking up all AWS regions available in this account")

	ec2Client, err := NewEc2ClientE(t, defaultRegion)
	if err != nil {
		return nil, err
	}

	out, err := ec2Client.DescribeRegions(&ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, err
	}

	regions := []string{}
	for _, region := range out.Regions {
		regions = append(regions, aws.StringValue(region.RegionName))
	}

	return regions, nil
}

// GetAvailabilityZones gets the Availability Zones for a given AWS region. Note that for certain regions (e.g. us-east-1), different AWS
// accounts have access to different availability zones.
func GetAvailabilityZones(t testing.TestingT, region string) []string {
	out, err := GetAvailabilityZonesE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetAvailabilityZonesE gets the Availability Zones for a given AWS region. Note that for certain regions (e.g. us-east-1), different AWS
// accounts have access to different availability zones.
func GetAvailabilityZonesE(t testing.TestingT, region string) ([]string, error) {
	logger.Logf(t, "Looking up all availability zones available in this account for region %s", region)

	ec2Client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	resp, err := ec2Client.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{})
	if err != nil {
		return nil, err
	}

	var out []string
	for _, availabilityZone := range resp.AvailabilityZones {
		out = append(out, aws.StringValue(availabilityZone.ZoneName))
	}

	return out, nil
}
