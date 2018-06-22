package aws

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Vpc is an Amazon Virtual Private Cloud.
type Vpc struct {
	Id      string   // The ID of the VPC
	Name    string   // The name of the VPC
	Subnets []Subnet // A list of subnets in the VPC
}

// Subnet is a subnet in an availability zone.
type Subnet struct {
	Id               string // The ID of the Subnet
	AvailabilityZone string // The Availability Zone the subnet is in
}

var vpcIDFilterName = "vpc-id"
var isDefaultFilterName = "isDefault"
var isDefaultFilterValue = "true"

// GetDefaultVpc fetches information about the default VPC in the given region.
func GetDefaultVpc(t *testing.T, region string) *Vpc {
	vpc, err := GetDefaultVpcE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return vpc
}

// GetDefaultVpcE fetches information about the default VPC in the given region.
func GetDefaultVpcE(t *testing.T, region string) (*Vpc, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	defaultVpcFilter := ec2.Filter{Name: &isDefaultFilterName, Values: []*string{&isDefaultFilterValue}}
	vpcs, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{Filters: []*ec2.Filter{&defaultVpcFilter}})
	if err != nil {
		return nil, err
	}

	numVpcs := len(vpcs.Vpcs)
	if numVpcs != 1 {
		return nil, fmt.Errorf("Expected to find one default VPC in region %s but found %s", region, strconv.Itoa(numVpcs))
	}

	defaultVpc := vpcs.Vpcs[0]

	subnets, err := GetSubnetsForVpcE(t, aws.StringValue(defaultVpc.VpcId), region)
	if err != nil {
		return nil, err
	}

	return &Vpc{Id: aws.StringValue(defaultVpc.VpcId), Name: FindVpcName(defaultVpc), Subnets: subnets}, nil
}

// FindVpcName extracts the VPC name from its tags (if any). Fall back to "Default" if it's the default VPC or empty string
// otherwise.
func FindVpcName(vpc *ec2.Vpc) string {
	for _, tag := range vpc.Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}

	if *vpc.IsDefault {
		return "Default"
	}

	return ""
}

// GetSubnetsForVpc gets the subnets in the specified VPC.
func GetSubnetsForVpc(t *testing.T, vpcID string, region string) []Subnet {
	subnets, err := GetSubnetsForVpcE(t, vpcID, region)
	if err != nil {
		t.Fatal(err)
	}
	return subnets
}

// GetSubnetsForVpcE gets the subnets in the specified VPC.
func GetSubnetsForVpcE(t *testing.T, vpcID string, region string) ([]Subnet, error) {
	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	vpcIDFilter := ec2.Filter{Name: &vpcIDFilterName, Values: []*string{&vpcID}}
	subnetOutput, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{Filters: []*ec2.Filter{&vpcIDFilter}})
	if err != nil {
		return nil, err
	}

	subnets := []Subnet{}

	for _, ec2Subnet := range subnetOutput.Subnets {
		subnet := Subnet{Id: aws.StringValue(ec2Subnet.SubnetId), AvailabilityZone: aws.StringValue(ec2Subnet.AvailabilityZone)}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
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
