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
	Id                   string            //nolint:staticcheck,revive // preserving existing field name
	Name                 string            // The name of the VPC
	Subnets              []Subnet          // A list of subnets in the VPC
	Tags                 map[string]string // The tags associated with the VPC
	CidrBlock            *string           // The primary IPv4 CIDR block for the VPC.
	CidrAssociations     []*string         // Information about the IPv4 CIDR blocks associated with the VPC.
	Ipv6CidrAssociations []*string         // Information about the IPv6 CIDR blocks associated with the VPC.
}

// Subnet is a subnet in an availability zone.
type Subnet struct {
	Tags             map[string]string // The tags associated with the subnet
	Id               string            //nolint:staticcheck,revive // preserving existing field name
	AvailabilityZone string            // The Availability Zone the subnet is in
	CidrBlock        string            // The CIDR block associated with the subnet
	DefaultForAz     bool              // If the subnet is default for the Availability Zone
}

const vpcIDFilterName = "vpc-id"
const defaultForAzFilterName = "default-for-az"
const resourceTypeFilterName = "resource-type"
const resourceIDFilterName = "resource-id"
const vpcResourceTypeFilterValue = "vpc"
const subnetResourceTypeFilterValue = "subnet"
const isDefaultFilterName = "isDefault"
const isDefaultFilterValue = "true"
const defaultVPCName = "Default"

// GetDefaultVpcContextE fetches information about the default VPC in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetDefaultVpcContextE(t testing.TestingT, ctx context.Context, region string) (*Vpc, error) {
	defaultVpcFilter := types.Filter{Name: aws.String(isDefaultFilterName), Values: []string{isDefaultFilterValue}}
	vpcs, err := GetVpcsContextE(t, ctx, []types.Filter{defaultVpcFilter}, region)

	numVpcs := len(vpcs)
	if numVpcs != 1 {
		return nil, fmt.Errorf("expected to find one default VPC in region %s but found %s", region, strconv.Itoa(numVpcs))
	}

	return vpcs[0], err
}

// GetDefaultVpcContext fetches information about the default VPC in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDefaultVpcContext(t testing.TestingT, ctx context.Context, region string) *Vpc {
	t.Helper()

	vpc, err := GetDefaultVpcContextE(t, ctx, region)
	require.NoError(t, err)

	return vpc
}

// GetDefaultVpc fetches information about the default VPC in the given region.
//
// Deprecated: Use [GetDefaultVpcContext] instead.
func GetDefaultVpc(t testing.TestingT, region string) *Vpc {
	t.Helper()

	return GetDefaultVpcContext(t, context.Background(), region)
}

// GetDefaultVpcE fetches information about the default VPC in the given region.
//
// Deprecated: Use [GetDefaultVpcContextE] instead.
func GetDefaultVpcE(t testing.TestingT, region string) (*Vpc, error) {
	return GetDefaultVpcContextE(t, context.Background(), region)
}

// GetVpcByIDContextE fetches information about a VPC with given ID in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetVpcByIDContextE(t testing.TestingT, ctx context.Context, vpcID string, region string) (*Vpc, error) {
	vpcIDFilter := types.Filter{Name: aws.String(vpcIDFilterName), Values: []string{vpcID}}
	vpcs, err := GetVpcsContextE(t, ctx, []types.Filter{vpcIDFilter}, region)

	numVpcs := len(vpcs)
	if numVpcs != 1 {
		return nil, fmt.Errorf("expected to find one VPC with ID %s in region %s but found %s", vpcID, region, strconv.Itoa(numVpcs))
	}

	return vpcs[0], err
}

// GetVpcByIDContext fetches information about a VPC with given ID in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVpcByIDContext(t testing.TestingT, ctx context.Context, vpcID string, region string) *Vpc {
	t.Helper()

	vpc, err := GetVpcByIDContextE(t, ctx, vpcID, region)
	require.NoError(t, err)

	return vpc
}

// GetVpcByID fetches information about a VPC with given ID in the given region.
//
// Deprecated: Use [GetVpcByIDContext] instead.
func GetVpcByID(t testing.TestingT, vpcID string, region string) *Vpc {
	t.Helper()

	return GetVpcByIDContext(t, context.Background(), vpcID, region)
}

// GetVpcByIDE fetches information about a VPC with given ID in the given region.
//
// Deprecated: Use [GetVpcByIDContextE] instead.
func GetVpcByIDE(t testing.TestingT, vpcID string, region string) (*Vpc, error) {
	return GetVpcByIDContextE(t, context.Background(), vpcID, region)
}

// GetVpcById fetches information about a VPC with given ID in the given region.
//
// Deprecated: Use [GetVpcByID] instead.
func GetVpcById(t testing.TestingT, vpcID string, region string) *Vpc { //nolint:staticcheck,revive // preserving deprecated function name
	return GetVpcByID(t, vpcID, region)
}

// GetVpcByIdE fetches information about a VPC with given ID in the given region.
//
// Deprecated: Use [GetVpcByIDE] instead.
func GetVpcByIdE(t testing.TestingT, vpcID string, region string) (*Vpc, error) { //nolint:staticcheck,revive // preserving deprecated function name
	return GetVpcByIDE(t, vpcID, region)
}

// GetVpcsContextE fetches information about VPCs from given regions limited by filters
// The ctx parameter supports cancellation and timeouts.
func GetVpcsContextE(t testing.TestingT, ctx context.Context, filters []types.Filter, region string) ([]*Vpc, error) {
	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	vpcs, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	numVpcs := len(vpcs.Vpcs)
	retVal := make([]*Vpc, numVpcs)

	for i := range vpcs.Vpcs {
		vpc := &vpcs.Vpcs[i]

		vpcIDFilter := generateVpcIDFilter(aws.ToString(vpc.VpcId))

		subnets, err := GetSubnetsForVpcContextE(t, ctx, region, []types.Filter{vpcIDFilter})
		if err != nil {
			return nil, err
		}

		tags, err := GetTagsForVpcContextE(t, ctx, aws.ToString(vpc.VpcId), region)
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
		var ipv6CidrAssociations = func() (list []*string) {
			for _, cidr := range vpc.Ipv6CidrBlockAssociationSet {
				list = append(list, cidr.Ipv6CidrBlock)
			}

			return
		}()

		retVal[i] = &Vpc{
			Id:                   aws.ToString(vpc.VpcId),
			Name:                 FindVPCName(vpc),
			Subnets:              subnets,
			Tags:                 tags,
			CidrBlock:            vpc.CidrBlock,
			CidrAssociations:     cidrBlockAssociations,
			Ipv6CidrAssociations: ipv6CidrAssociations,
		}
	}

	return retVal, nil
}

// GetVpcsE fetches information about VPCs from given regions limited by filters
//
// Deprecated: Use [GetVpcsContextE] instead.
func GetVpcsE(t testing.TestingT, filters []types.Filter, region string) ([]*Vpc, error) {
	return GetVpcsContextE(t, context.Background(), filters, region)
}

// FindVPCName extracts the VPC name from its tags (if any). Falls back to "Default" if it's the default VPC or empty
// string otherwise.
func FindVPCName(vpc *types.Vpc) string {
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

// FindVpcName extracts the VPC name from its tags (if any). Fall back to "Default" if it's the default VPC or empty string
// otherwise.
//
// Deprecated: Use [FindVPCName] instead.
func FindVpcName(vpc types.Vpc) string { //nolint:staticcheck,revive,gocritic // preserving deprecated function name
	return FindVPCName(&vpc)
}

// GetSubnetsForVpcContextE gets the subnets in the specified VPC.
// The ctx parameter supports cancellation and timeouts.
func GetSubnetsForVpcContextE(t testing.TestingT, ctx context.Context, region string, filters []types.Filter) ([]Subnet, error) {
	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	subnetOutput, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{Filters: filters})
	if err != nil {
		return nil, err
	}

	var subnets []Subnet

	for i := range subnetOutput.Subnets {
		ec2Subnet := &subnetOutput.Subnets[i]

		subnetTags := GetTagsForSubnetContext(t, ctx, *ec2Subnet.SubnetId, region)
		subnet := Subnet{Id: aws.ToString(ec2Subnet.SubnetId), AvailabilityZone: aws.ToString(ec2Subnet.AvailabilityZone), DefaultForAz: aws.ToBool(ec2Subnet.DefaultForAz), Tags: subnetTags, CidrBlock: aws.ToString(ec2Subnet.CidrBlock)}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}

// GetSubnetsForVpcContext gets the subnets in the specified VPC.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSubnetsForVpcContext(t testing.TestingT, ctx context.Context, vpcID string, region string) []Subnet {
	t.Helper()

	vpcIDFilter := generateVpcIDFilter(vpcID)

	subnets, err := GetSubnetsForVpcContextE(t, ctx, region, []types.Filter{vpcIDFilter})
	require.NoError(t, err)

	return subnets
}

// GetSubnetsForVpc gets the subnets in the specified VPC.
//
// Deprecated: Use [GetSubnetsForVpcContext] instead.
func GetSubnetsForVpc(t testing.TestingT, vpcID string, region string) []Subnet {
	t.Helper()

	return GetSubnetsForVpcContext(t, context.Background(), vpcID, region)
}

// GetSubnetsForVpcE gets the subnets in the specified VPC.
//
// Deprecated: Use [GetSubnetsForVpcContextE] instead.
func GetSubnetsForVpcE(t testing.TestingT, region string, filters []types.Filter) ([]Subnet, error) {
	return GetSubnetsForVpcContextE(t, context.Background(), region, filters)
}

// GetAzDefaultSubnetsForVpcContextE gets the default az subnets in the specified VPC.
// The ctx parameter supports cancellation and timeouts.
func GetAzDefaultSubnetsForVpcContextE(t testing.TestingT, ctx context.Context, vpcID string, region string) ([]Subnet, error) {
	vpcIDFilter := generateVpcIDFilter(vpcID)
	defaultForAzFilter := types.Filter{
		Name:   aws.String(defaultForAzFilterName),
		Values: []string{"true"},
	}

	return GetSubnetsForVpcContextE(t, ctx, region, []types.Filter{vpcIDFilter, defaultForAzFilter})
}

// GetAzDefaultSubnetsForVpcContext gets the default az subnets in the specified VPC.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAzDefaultSubnetsForVpcContext(t testing.TestingT, ctx context.Context, vpcID string, region string) []Subnet {
	t.Helper()

	subnets, err := GetAzDefaultSubnetsForVpcContextE(t, ctx, vpcID, region)
	require.NoError(t, err)

	return subnets
}

// GetAzDefaultSubnetsForVpc gets the default az subnets in the specified VPC.
//
// Deprecated: Use [GetAzDefaultSubnetsForVpcContext] instead.
func GetAzDefaultSubnetsForVpc(t testing.TestingT, vpcID string, region string) []Subnet {
	t.Helper()

	return GetAzDefaultSubnetsForVpcContext(t, context.Background(), vpcID, region)
}

// GetAzDefaultSubnetsForVpcE gets the default az subnets in the specified VPC.
//
// Deprecated: Use [GetAzDefaultSubnetsForVpcContextE] instead.
func GetAzDefaultSubnetsForVpcE(t testing.TestingT, vpcID string, region string) ([]Subnet, error) {
	return GetAzDefaultSubnetsForVpcContextE(t, context.Background(), vpcID, region)
}

// generateVpcIDFilter is a helper method to generate vpc id filter
func generateVpcIDFilter(vpcID string) types.Filter {
	return types.Filter{Name: aws.String(vpcIDFilterName), Values: []string{vpcID}}
}

// getTagsForResourceContextE is a helper that gets the tags for a specified EC2 resource.
// The ctx parameter supports cancellation and timeouts.
func getTagsForResourceContextE(t testing.TestingT, ctx context.Context, resourceType string, resourceID string, region string) (map[string]string, error) {
	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	resourceTypeFilter := types.Filter{Name: aws.String(resourceTypeFilterName), Values: []string{resourceType}}
	resourceIDFilter := types.Filter{Name: aws.String(resourceIDFilterName), Values: []string{resourceID}}

	tagsOutput, err := client.DescribeTags(ctx, &ec2.DescribeTagsInput{Filters: []types.Filter{resourceTypeFilter, resourceIDFilter}})
	if err != nil {
		return nil, err
	}

	tags := map[string]string{}

	for _, tag := range tagsOutput.Tags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}

	return tags, nil
}

// GetTagsForVpcContextE gets the tags for the specified VPC.
// The ctx parameter supports cancellation and timeouts.
func GetTagsForVpcContextE(t testing.TestingT, ctx context.Context, vpcID string, region string) (map[string]string, error) {
	return getTagsForResourceContextE(t, ctx, vpcResourceTypeFilterValue, vpcID, region)
}

// GetTagsForVpcContext gets the tags for the specified VPC.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTagsForVpcContext(t testing.TestingT, ctx context.Context, vpcID string, region string) map[string]string {
	t.Helper()

	tags, err := GetTagsForVpcContextE(t, ctx, vpcID, region)
	require.NoError(t, err)

	return tags
}

// GetTagsForVpc gets the tags for the specified VPC.
//
// Deprecated: Use [GetTagsForVpcContext] instead.
func GetTagsForVpc(t testing.TestingT, vpcID string, region string) map[string]string {
	t.Helper()

	return GetTagsForVpcContext(t, context.Background(), vpcID, region)
}

// GetTagsForVpcE gets the tags for the specified VPC.
//
// Deprecated: Use [GetTagsForVpcContextE] instead.
func GetTagsForVpcE(t testing.TestingT, vpcID string, region string) (map[string]string, error) {
	return GetTagsForVpcContextE(t, context.Background(), vpcID, region)
}

// GetDefaultSubnetIDsForVpcPContextE gets the ids of the subnets that are the default subnet for the AvailabilityZone.
// The P suffix differentiates this function (which accepts *Vpc pointer) from the deprecated
// GetDefaultSubnetIDsForVpcE which accepts Vpc by value.
// The ctx parameter is accepted for API consistency with other Context functions.
func GetDefaultSubnetIDsForVpcPContextE(t testing.TestingT, ctx context.Context, vpc *Vpc) ([]string, error) {
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

// GetDefaultSubnetIDsForVpcPContext gets the ids of the subnets that are the default subnet for the AvailabilityZone.
// This function will fail the test if there is an error.
// The P suffix differentiates this function (which accepts *Vpc pointer) from the deprecated
// GetDefaultSubnetIDsForVpc which accepts Vpc by value.
// The ctx parameter is accepted for API consistency with other Context functions.
func GetDefaultSubnetIDsForVpcPContext(t testing.TestingT, ctx context.Context, vpc *Vpc) []string {
	t.Helper()

	subnetIDs, err := GetDefaultSubnetIDsForVpcPContextE(t, ctx, vpc)
	require.NoError(t, err)

	return subnetIDs
}

// GetDefaultSubnetIDsForVpcP gets the ids of the subnets that are the default subnet for the AvailabilityZone.
//
// Deprecated: Use [GetDefaultSubnetIDsForVpcPContext] instead.
func GetDefaultSubnetIDsForVpcP(t testing.TestingT, vpc *Vpc) []string {
	t.Helper()

	return GetDefaultSubnetIDsForVpcPContext(t, context.Background(), vpc)
}

// GetDefaultSubnetIDsForVpcPE gets the ids of the subnets that are the default subnet for the AvailabilityZone.
//
// Deprecated: Use [GetDefaultSubnetIDsForVpcPContextE] instead.
func GetDefaultSubnetIDsForVpcPE(t testing.TestingT, vpc *Vpc) ([]string, error) {
	return GetDefaultSubnetIDsForVpcPContextE(t, context.Background(), vpc)
}

// GetDefaultSubnetIDsForVpc gets the ids of the subnets that are the default subnet for the AvailabilityZone.
//
// Deprecated: Use [GetDefaultSubnetIDsForVpcPContext] instead.
func GetDefaultSubnetIDsForVpc(t testing.TestingT, vpc Vpc) []string { //nolint:gocritic // preserving deprecated function signature
	t.Helper()

	return GetDefaultSubnetIDsForVpcPContext(t, context.Background(), &vpc)
}

// GetDefaultSubnetIDsForVpcE gets the ids of the subnets that are the default subnet for the AvailabilityZone.
//
// Deprecated: Use [GetDefaultSubnetIDsForVpcPContextE] instead.
func GetDefaultSubnetIDsForVpcE(t testing.TestingT, vpc Vpc) ([]string, error) { //nolint:gocritic // preserving deprecated function signature
	return GetDefaultSubnetIDsForVpcPContextE(t, context.Background(), &vpc)
}

// GetTagsForSubnetContextE gets the tags for the specified subnet.
// The ctx parameter supports cancellation and timeouts.
func GetTagsForSubnetContextE(t testing.TestingT, ctx context.Context, subnetID string, region string) (map[string]string, error) {
	return getTagsForResourceContextE(t, ctx, subnetResourceTypeFilterValue, subnetID, region)
}

// GetTagsForSubnetContext gets the tags for the specified subnet.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTagsForSubnetContext(t testing.TestingT, ctx context.Context, subnetID string, region string) map[string]string {
	t.Helper()

	tags, err := GetTagsForSubnetContextE(t, ctx, subnetID, region)
	require.NoError(t, err)

	return tags
}

// GetTagsForSubnet gets the tags for the specified subnet.
//
// Deprecated: Use [GetTagsForSubnetContext] instead.
func GetTagsForSubnet(t testing.TestingT, subnetID string, region string) map[string]string {
	t.Helper()

	return GetTagsForSubnetContext(t, context.Background(), subnetID, region)
}

// GetTagsForSubnetE gets the tags for the specified subnet.
//
// Deprecated: Use [GetTagsForSubnetContextE] instead.
func GetTagsForSubnetE(t testing.TestingT, subnetID string, region string) (map[string]string, error) {
	return GetTagsForSubnetContextE(t, context.Background(), subnetID, region)
}

// IsPublicSubnetContextE returns True if the subnet identified by the given id in the provided region is public.
// The ctx parameter supports cancellation and timeouts.
func IsPublicSubnetContextE(t testing.TestingT, ctx context.Context, subnetID string, region string) (bool, error) {
	subnetIDFilterName := "association.subnet-id"

	subnetIDFilter := types.Filter{
		Name:   &subnetIDFilterName,
		Values: []string{subnetID},
	}

	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return false, err
	}

	rts, err := client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{Filters: []types.Filter{subnetIDFilter}})
	if err != nil {
		return false, err
	}

	if len(rts.RouteTables) == 0 {
		// Subnets not explicitly associated with any route table are implicitly associated with the main route table
		rts, err = getImplicitRouteTableForSubnetContextE(t, ctx, subnetID, region)
		if err != nil {
			return false, err
		}
	}

	for i := range rts.RouteTables {
		rt := &rts.RouteTables[i]

		for j := range rt.Routes {
			r := &rt.Routes[j]

			if strings.HasPrefix(aws.ToString(r.GatewayId), "igw-") {
				return true, nil
			}
		}
	}

	return false, nil
}

// IsPublicSubnetContext returns True if the subnet identified by the given id in the provided region is public.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func IsPublicSubnetContext(t testing.TestingT, ctx context.Context, subnetID string, region string) bool {
	t.Helper()

	isPublic, err := IsPublicSubnetContextE(t, ctx, subnetID, region)
	require.NoError(t, err)

	return isPublic
}

// IsPublicSubnet returns True if the subnet identified by the given id in the provided region is public.
//
// Deprecated: Use [IsPublicSubnetContext] instead.
func IsPublicSubnet(t testing.TestingT, subnetID string, region string) bool {
	t.Helper()

	return IsPublicSubnetContext(t, context.Background(), subnetID, region)
}

// IsPublicSubnetE returns True if the subnet identified by the given id in the provided region is public.
//
// Deprecated: Use [IsPublicSubnetContextE] instead.
func IsPublicSubnetE(t testing.TestingT, subnetID string, region string) (bool, error) {
	return IsPublicSubnetContextE(t, context.Background(), subnetID, region)
}

// getImplicitRouteTableForSubnetContextE gets the implicit route table for a subnet.
func getImplicitRouteTableForSubnetContextE(t testing.TestingT, ctx context.Context, subnetID string, region string) (*ec2.DescribeRouteTablesOutput, error) {
	mainRouteFilterName := "association.main"
	mainRouteFilterValue := "true"
	subnetFilterName := "subnet-id"

	client, err := NewEc2ClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	subnetFilter := types.Filter{
		Name:   &subnetFilterName,
		Values: []string{subnetID},
	}

	subnetOutput, err := client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{Filters: []types.Filter{subnetFilter}})
	if err != nil {
		return nil, err
	}

	numSubnets := len(subnetOutput.Subnets)
	if numSubnets != 1 {
		return nil, fmt.Errorf("expected to find one subnet with id %s but found %s", subnetID, strconv.Itoa(numSubnets))
	}

	mainRouteFilter := types.Filter{
		Name:   &mainRouteFilterName,
		Values: []string{mainRouteFilterValue},
	}
	vpcFilter := types.Filter{
		Name:   aws.String(vpcIDFilterName),
		Values: []string{*subnetOutput.Subnets[0].VpcId},
	}

	return client.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{Filters: []types.Filter{mainRouteFilter, vpcFilter}})
}

// GetRandomPrivateCidrBlock gets a random CIDR block from the range of acceptable private IP addresses per RFC 1918
// (https://tools.ietf.org/html/rfc1918#section-3)
// The routingPrefix refers to the "/28" in 1.2.3.4/28.
// Note that, as written, this function will return a subset of all valid ranges. Since we will probably use this function
// mostly for generating random CIDR ranges for VPCs and Subnets, having comprehensive set coverage is not essential.
func GetRandomPrivateCidrBlock(routingPrefix int) string {
	var o1, o2, o3, o4 int

	switch routingPrefix {
	case 32: //nolint:mnd // RFC 1918 private address range
		o1 = random.RandomInt([]int{10, 172, 192}) //nolint:mnd // RFC 1918 private address range

		switch o1 {
		case 10: //nolint:mnd // RFC 1918 private address range
			o2 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
			o3 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
			o4 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
		case 172: //nolint:mnd // RFC 1918 private address range
			o2 = random.Random(16, 31) //nolint:mnd // RFC 1918 private address range
			o3 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
			o4 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
		case 192: //nolint:mnd // RFC 1918 private address range
			o2 = 168                   //nolint:mnd // RFC 1918 private address range
			o3 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
			o4 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
		}

	case 31, 30, 29, 28, 27, 26, 25: //nolint:mnd // RFC 1918 private address range
		fallthrough
	case 24: //nolint:mnd // RFC 1918 private address range
		o1 = random.RandomInt([]int{10, 172, 192}) //nolint:mnd // RFC 1918 private address range

		switch o1 {
		case 10: //nolint:mnd // RFC 1918 private address range
			o2 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
			o3 = random.Random(0, 255) //nolint:mnd // RFC 1918 private address range
			o4 = 0
		case 172: //nolint:mnd // RFC 1918 private address range
			o2 = 16 //nolint:mnd // RFC 1918 private address range
			o3 = 0
			o4 = 0
		case 192: //nolint:mnd // RFC 1918 private address range
			o2 = 168 //nolint:mnd // RFC 1918 private address range
			o3 = 0
			o4 = 0
		}
	case 23, 22, 21, 20, 19: //nolint:mnd // RFC 1918 private address range
		fallthrough
	case 18: //nolint:mnd // RFC 1918 private address range
		o1 = random.RandomInt([]int{10, 172, 192}) //nolint:mnd // RFC 1918 private address range

		switch o1 {
		case 10: //nolint:mnd // RFC 1918 private address range
			o2 = 0
			o3 = 0
			o4 = 0
		case 172: //nolint:mnd // RFC 1918 private address range
			o2 = 16 //nolint:mnd // RFC 1918 private address range
			o3 = 0
			o4 = 0
		case 192: //nolint:mnd // RFC 1918 private address range
			o2 = 168 //nolint:mnd // RFC 1918 private address range
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
