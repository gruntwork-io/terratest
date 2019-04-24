package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/stretchr/testify/require"
)

// GetElbV2Arn fetches ARN information about specified ELB.
func GetElbV2Arn(t *testing.T, region string, name string) *string {
	arn, err := GetElbV2ArnE(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
	return arn
}

// GetElbV2ArnE fetches ARN information about specified ELB.
func GetElbV2ArnE(t *testing.T, region string, name string) (*string, error) {
	client := NewElbV2Client(t, region)

	resp, err := client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(name),
		},
	})
	if err != nil {
		return nil, err
	}

	numElb := len(resp.LoadBalancers)
	if numElb != 1 {
		return nil, fmt.Errorf("Expected to find 1 ELB named '%s' in region '%v', but found '%d'",
			name, region, numElb)
	}

	return resp.LoadBalancers[0].LoadBalancerArn, nil
}

// GetElbV2CanonicalHostedZoneID fetches Canonical Hosted Zone ID information about specified ELB.
func GetElbV2CanonicalHostedZoneID(t *testing.T, region string, name string) *string {
	arn, err := GetElbV2CanonicalHostedZoneIDE(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
	return arn
}

// GetElbV2CanonicalHostedZoneIDE fetches Canonical Hosted Zone ID information about specified ELB.
func GetElbV2CanonicalHostedZoneIDE(t *testing.T, region string, name string) (*string, error) {
	client := NewElbV2Client(t, region)

	resp, err := client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(name),
		},
	})
	if err != nil {
		return nil, err
	}

	numElb := len(resp.LoadBalancers)
	if numElb != 1 {
		return nil, fmt.Errorf("Expected to find 1 ELB named '%s' in region '%v', but found '%d'",
			name, region, numElb)
	}

	return resp.LoadBalancers[0].CanonicalHostedZoneId, nil
}

// GetElbV2DNSName fetches DNS Name information about specified ELB.
func GetElbV2DNSName(t *testing.T, region string, name string) *string {
	arn, err := GetElbV2DNSNameE(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
	return arn
}

// GetElbV2DNSNameE fetches DNS Name information about specified ELB.
func GetElbV2DNSNameE(t *testing.T, region string, name string) (*string, error) {
	client := NewElbV2Client(t, region)

	resp, err := client.DescribeLoadBalancers(&elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(name),
		},
	})
	if err != nil {
		return nil, err
	}

	numElb := len(resp.LoadBalancers)
	if numElb != 1 {
		return nil, fmt.Errorf("Expected to find 1 ELB named '%s' in region '%v', but found '%d'",
			name, region, numElb)
	}

	return resp.LoadBalancers[0].DNSName, nil
}

// CreateElbV2 creates ELB in the given region under the given name and subnets list.
func CreateElbV2(t *testing.T, region string, name string, subnets []*string) {
	err := CreateElbV2E(t, region, name, subnets)
	if err != nil {
		t.Fatal(err)
	}
}

// CreateElbV2E creates ELB in the given region under the given name and subnets list.
func CreateElbV2E(t *testing.T, region string, name string, subnets []*string) error {
	client := NewElbV2Client(t, region)
	elb, err := client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(name),
		Subnets: subnets,
	})
	if err != nil {
		return err
	}

	numElb := len(elb.LoadBalancers)
	if numElb != 1 {
		return fmt.Errorf("Expected to create 1 ELB named '%s' in region '%v', but found '%d'",
			name, region, numElb)
	}

	return nil
}

// DeleteElbV2 deletes existing ELB in the given region.
func DeleteElbV2(t *testing.T, region string, name string) {
	err := DeleteElbV2E(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteElbV2E deletes existing ELB in the given region.
func DeleteElbV2E(t *testing.T, region string, name string) error {
	client := NewElbV2Client(t, region)
	arn := GetElbV2Arn(t, region, name)
	_, err := client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: arn,
	})

	return err
}

// NewElbV2Client creates en ELB client.
func NewElbV2Client(t *testing.T, region string) *elbv2.ELBV2 {
	client, err := NewElbV2ClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewElbV2ClientE creates an ELB client.
func NewElbV2ClientE(t *testing.T, region string) (*elbv2.ELBV2, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return elbv2.New(sess), nil
}
