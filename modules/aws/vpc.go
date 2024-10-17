package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Vpc is an Amazon Virtual Private Cloud.
type Vpc struct {
	Id                   string            // The ID of the VPC
	Name                 string            // The name of the VPC
	Subnets              []Subnet          // A list of subnets in the VPC
	Tags                 map[string]string // The tags associated with the VPC
	CidrBlock            *string           // The primary IPv4 CIDR block for the VPC.
	CidrAssociations     []*string         // Information about the IPv4 CIDR blocks associated with the VPC.
	Ipv6CidrAssociations []*string         // Information about the IPv6 CIDR blocks associated with the VPC.
}

// Subnet is a subnet in an availability zone.
type Subnet struct {
	Id               string            // The ID of the Subnet
	AvailabilityZone string            // The Availability Zone the subnet is in
	DefaultForAz     bool              // If the subnet is default for the Availability Zone
	Tags             map[string]string // The tags associated with the subnet
}

const vpcIDFilterName = "vpc-id"
const defaultForAzFilterName = "default-for-az"
const resourceTypeFilterName = "resource-id"
const resourceIdFilterName = "resource-type"
const vpcResourceTypeFilterValue = "vpc"
const subnetResourceTypeFilterValue = "subnet"
const isDefaultFilterName = "isDefault"
const isDefaultFilterValue = "true"
const defaultVPCName = "Default"

// GetDefaultVpc fetches information about the default VPC in the given region.
func GetDefaultVpc(t testing.TestingT, region string) *Vpc {
	vpc, err := GetDefaultVpcE(t, region)
	require.NoError(t, err)
	return vpc
}

// GetDefaultVpcE fetches information about the default VPC in the given region.
func GetDefaultVpcE(t testing.TestingT, region string) (*Vpc, error) {
	defaultVpcFilter := types.Filter{Name: aws.String(isDefaultFilterName), Values: []string{isDefaultFilterValue}}
	vpcs, err := GetVpcsE(t, []types.Filter{defaultVpcFilter}, region)

	numVpcs := len(vpcs)
	if numVpcs != 1 {
		return nil, fmt.Errorf("expected to find one default VPC in region %s but found %s", region, strconv.Itoa(numVpcs))
	}

	return vpcs[0], err
}

// GetVpcById fetches information about a VPC with given ID in the given region.
func GetVpcById(t testing.TestingT, vpcId string, region string) *Vpc {
	vpc, err := GetVpcByIdE(t, vpcId, region)
	require.NoError(t, err)
	return vpc
}

// GetVpcByIdE fetches information about a VPC with given ID in the given region.
func GetVpcByIdE(t testing.TestingT, vpcId string, region string) (*Vpc, error) {
	vpcIdFilter := types.Filter{Name: aws.String(vpcIDFilterName), Values: []string{vpcId}}
	vpcs, err := GetVpcsE(t, []types.Filter{vpcIdFilter}, region)

	numVpcs := len(vpcs)
	if numVpcs != 1 {
		return nil, fmt.Errorf("expected to find one VPC with ID %s in region %s but found %s", vpcId, region, strconv.Itoa(numVpcs))
	}

	return vpcs[0], err
}

// GetVpcsE fetches information about VPCs from given regions limited by filters
func GetVpcsE(t testing.TestingT, filters []types.Filter, region string) ([]*Vpc, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	vpcs, err := client.DescribeVpcs(context.Background(), &ec2.DescribeVpcsInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	numVpcs := len(vpcs.Vpcs)
	retVal := make([]*Vpc, numVpcs)

	for i, vpc := range vpcs.Vpcs {
		vpcIdFilter := generateVpcIdFilter(aws.ToString(vpc.VpcId))
		subnets, err := GetSubnetsForVpcE(t, region, []types.Filter{vpcIdFilter})
		if err != nil {
			return nil, err
		}

		tags, err := GetTagsForVpcE(t, aws.ToString(vpc.VpcId), region)
		if err != nil {
			return nil, err
		}

		// cidr block associations
		var cidrBlockAssociations = func() (list []*string) {
			for _, cidr := range vpc.CidrBlockAssociationSet {
				list = append(list, cidr.CidrBlock)
			}
			return
		}()

		// ipv6 cidr block associations
		var Ipv6CidrAssociations = func() (list []*string) {
			for _, cidr := range vpc.Ipv6CidrBlockAssociationSet {
				list = append(list, cidr.Ipv6CidrBlock)
			}
			return
		}()

		retVal[i] = &Vpc{
			Id:                   aws.ToString(vpc.VpcId),
			Name:                 FindVpcName(vpc),
			Subnets:              subnets,
			Tags:                 tags,
			CidrBlock:            vpc.CidrBlock,
			CidrAssociations:     cidrBlockAssociations,
			Ipv6CidrAssociations: Ipv6CidrAssociations,
		}
	}

	return retVal, nil
}

// FindVpcName extracts the VPC name from its tags (if any). Fall back to "Default" if it's the default VPC or empty string
// otherwise.
func FindVpcName(vpc types.Vpc) string {
	for _, tag := range vpc.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}

	if *vpc.IsDefault {
		return defaultVPCName
	}

	return ""
}

// GetSubnetsForVpc gets the subnets in the specified VPC.
func GetSubnetsForVpc(t testing.TestingT, vpcID string, region string) []Subnet {
	vpcIDFilter := generateVpcIdFilter(vpcID)
	subnets, err := GetSubnetsForVpcE(t, region, []types.Filter{vpcIDFilter})
	if err != nil {
		t.Fatal(err)
	}
	return subnets
}

// GetAzDefaultSubnetsForVpc gets the default az subnets in the specified VPC.
func GetAzDefaultSubnetsForVpc(t testing.TestingT, vpcID string, region string) []Subnet {
	vpcIDFilter := generateVpcIdFilter(vpcID)
	defaultForAzFilter := types.Filter{
		Name:   aws.String(defaultForAzFilterName),
		Values: []string{"true"},
	}
	subnets, err := GetSubnetsForVpcE(t, region, []types.Filter{vpcIDFilter, defaultForAzFilter})
	if err != nil {
		t.Fatal(err)
	}
	return subnets
}

// generateVpcIdFilter is a helper method to generate vpc id filter
func generateVpcIdFilter(vpcID string) types.Filter {
	return types.Filter{Name: aws.String(vpcIDFilterName), Values: []string{vpcID}}
}

// GetSubnetsForVpcE gets the subnets in the specified VPC.
func GetSubnetsForVpcE(t testing.TestingT, region string, filters []types.Filter) ([]Subnet, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	subnetOutput, err := client.DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	var subnets []Subnet

	for _, ec2Subnet := range subnetOutput.Subnets {
		subnetTags := GetTagsForSubnet(t, *ec2Subnet.SubnetId, region)
		subnet := Subnet{Id: aws.ToString(ec2Subnet.SubnetId), AvailabilityZone: aws.ToString(ec2Subnet.AvailabilityZone), DefaultForAz: aws.ToBool(ec2Subnet.DefaultForAz), Tags: subnetTags}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}

// GetTagsForVpc gets the tags for the specified VPC.
func GetTagsForVpc(t testing.TestingT, vpcID string, region string) map[string]string {
	tags, err := GetTagsForVpcE(t, vpcID, region)
	require.NoError(t, err)

	return tags
}

// GetTagsForVpcE gets the tags for the specified VPC.
func GetTagsForVpcE(t testing.TestingT, vpcID string, region string) (map[string]string, error) {
	client, err := NewEc2ClientE(t, region)
	require.NoError(t, err)

	vpcResourceTypeFilter := types.Filter{Name: aws.String(resourceIdFilterName), Values: []string{vpcResourceTypeFilterValue}}
	vpcResourceIdFilter := types.Filter{Name: aws.String(resourceTypeFilterName), Values: []string{vpcID}}
	tagsOutput, err := client.DescribeTags(context.Background(), &ec2.DescribeTagsInput{Filters: []types.Filter{vpcResourceTypeFilter, vpcResourceIdFilter}})
	require.NoError(t, err)

	tags := map[string]string{}
	for _, tag := range tagsOutput.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return tags, nil
}

// GetDefaultSubnetIDsForVpc gets the ids of the subnets that are the default subnet for the AvailabilityZone
func GetDefaultSubnetIDsForVpc(t testing.TestingT, vpc Vpc) []string {
	subnetIDs, err := GetDefaultSubnetIDsForVpcE(t, vpc)
	require.NoError(t, err)
	return subnetIDs
}

// GetDefaultSubnetIDsForVpcE gets the ids of the subnets that are the default subnet for the AvailabilityZone
func GetDefaultSubnetIDsForVpcE(t testing.TestingT, vpc Vpc) ([]string, error) {
	if vpc.Name != defaultVPCName {
		// You cannot create a default subnet in a nondefault VPC
		// https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html
		return nil, fmt.Errorf("only default VPCs have default subnets but VPC with id %s is not default VPC", vpc.Id)
	}
	var subnetIDs []string
	numSubnets := len(vpc.Subnets)
	if numSubnets == 0 {
		return nil, fmt.Errorf("expected to find at least one subnet in vpc with ID %s but found zero", vpc.Id)
	}

	for _, subnet := range vpc.Subnets {
		if subnet.DefaultForAz {
			subnetIDs = append(subnetIDs, subnet.Id)
		}
	}
	return subnetIDs, nil
}

// GetTagsForSubnet gets the tags for the specified subnet.
func GetTagsForSubnet(t testing.TestingT, subnetId string, region string) map[string]string {
	tags, err := GetTagsForSubnetE(t, subnetId, region)
	require.NoError(t, err)

	return tags
}

// GetTagsForSubnetE gets the tags for the specified subnet.
func GetTagsForSubnetE(t testing.TestingT, subnetId string, region string) (map[string]string, error) {
	client, err := NewEc2ClientE(t, region)
	require.NoError(t, err)

	subnetResourceTypeFilter := types.Filter{Name: aws.String(resourceIdFilterName), Values: []string{subnetResourceTypeFilterValue}}
	subnetResourceIdFilter := types.Filter{Name: aws.String(resourceTypeFilterName), Values: []string{subnetId}}
	tagsOutput, err := client.DescribeTags(context.Background(), &ec2.DescribeTagsInput{Filters: []types.Filter{subnetResourceTypeFilter, subnetResourceIdFilter}})
	require.NoError(t, err)

	tags := map[string]string{}
	for _, tag := range tagsOutput.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return tags, nil
}

// IsPublicSubnet returns True if the subnet identified by the given id in the provided region is public.
func IsPublicSubnet(t testing.TestingT, subnetId string, region string) bool {
	isPublic, err := IsPublicSubnetE(t, subnetId, region)
	require.NoError(t, err)
	return isPublic
}

// IsPublicSubnetE returns True if the subnet identified by the given id in the provided region is public.
func IsPublicSubnetE(t testing.TestingT, subnetId string, region string) (bool, error) {
	subnetIdFilterName := "association.subnet-id"

	subnetIdFilter := types.Filter{
		Name:   &subnetIdFilterName,
		Values: []string{subnetId},
	}

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return false, err
	}

	rts, err := client.DescribeRouteTables(context.Background(), &ec2.DescribeRouteTablesInput{Filters: []types.Filter{subnetIdFilter}})
	if err != nil {
		return false, err
	}

	if len(rts.RouteTables) == 0 {
		// Subnets not explicitly associated with any route table are implicitly associated with the main route table
		rts, err = getImplicitRouteTableForSubnetE(t, subnetId, region)
		if err != nil {
			return false, err
		}
	}

	for _, rt := range rts.RouteTables {
		for _, r := range rt.Routes {
			if strings.HasPrefix(aws.ToString(r.GatewayId), "igw-") {
				return true, nil
			}
		}
	}

	return false, nil
}

func getImplicitRouteTableForSubnetE(t testing.TestingT, subnetId string, region string) (*ec2.DescribeRouteTablesOutput, error) {
	mainRouteFilterName := "association.main"
	mainRouteFilterValue := "true"
	subnetFilterName := "subnet-id"

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	subnetFilter := types.Filter{
		Name:   &subnetFilterName,
		Values: []string{subnetId},
	}
	subnetOutput, err := client.DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{Filters: []types.Filter{subnetFilter}})
	if err != nil {
		return nil, err
	}
	numSubnets := len(subnetOutput.Subnets)
	if numSubnets != 1 {
		return nil, fmt.Errorf("expected to find one subnet with id %s but found %s", subnetId, strconv.Itoa(numSubnets))
	}

	mainRouteFilter := types.Filter{
		Name:   &mainRouteFilterName,
		Values: []string{mainRouteFilterValue},
	}
	vpcFilter := types.Filter{
		Name:   aws.String(vpcIDFilterName),
		Values: []string{*subnetOutput.Subnets[0].VpcId},
	}
	return client.DescribeRouteTables(context.Background(), &ec2.DescribeRouteTablesInput{Filters: []types.Filter{mainRouteFilter, vpcFilter}})
}

// GetRandomPrivateCidrBlock gets a random CIDR block from the range of acceptable private IP addresses per RFC 1918
// (https://tools.ietf.org/html/rfc1918#section-3)
// The routingPrefix refers to the "/28" in 1.2.3.4/28.
// Note that, as written, this function will return a subset of all valid ranges. Since we will probably use this function
// mostly for generating random CIDR ranges for VPCs and Subnets, having comprehensive set coverage is not essential.
func GetRandomPrivateCidrBlock(routingPrefix int) string {

	var o1, o2, o3, o4 int

	switch routingPrefix {
	case 32:
		o1 = random.RandomInt([]int{10, 172, 192})

		switch o1 {
		case 10:
			o2 = random.Random(0, 255)
			o3 = random.Random(0, 255)
			o4 = random.Random(0, 255)
		case 172:
			o2 = random.Random(16, 31)
			o3 = random.Random(0, 255)
			o4 = random.Random(0, 255)
		case 192:
			o2 = 168
			o3 = random.Random(0, 255)
			o4 = random.Random(0, 255)
		}

	case 31, 30, 29, 28, 27, 26, 25:
		fallthrough
	case 24:
		o1 = random.RandomInt([]int{10, 172, 192})

		switch o1 {
		case 10:
			o2 = random.Random(0, 255)
			o3 = random.Random(0, 255)
			o4 = 0
		case 172:
			o2 = 16
			o3 = 0
			o4 = 0
		case 192:
			o2 = 168
			o3 = 0
			o4 = 0
		}
	case 23, 22, 21, 20, 19:
		fallthrough
	case 18:
		o1 = random.RandomInt([]int{10, 172, 192})

		switch o1 {
		case 10:
			o2 = 0
			o3 = 0
			o4 = 0
		case 172:
			o2 = 16
			o3 = 0
			o4 = 0
		case 192:
			o2 = 168
			o3 = 0
			o4 = 0
		}
	}
	return fmt.Sprintf("%d.%d.%d.%d/%d", o1, o2, o3, o4, routingPrefix)
}

// GetFirstTwoOctets gets the first two octets from a CIDR block.
func GetFirstTwoOctets(cidrBlock string) string {
	ipAddr := strings.Split(cidrBlock, "/")[0]
	octets := strings.Split(ipAddr, ".")
	return octets[0] + "." + octets[1]
}
