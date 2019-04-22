package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/stretchr/testify/require"
)

// GetElbV2 fetches information about specified ELB.
func GetElbV2(t *testing.T, region string, name string) *elbv2.LoadBalancer {
	elb, err := GetElbV2E(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
	return elb
}

// GetElbV2E fetches information about specified ELB.
func GetElbV2E(t *testing.T, region string, name string) (*elbv2.LoadBalancer, error) {
	client := NewElbV2Client(t, region)

	input := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{
			aws.String(name),
		},
	}
	output, err := client.DescribeLoadBalancers(input)
	if err != nil {
		return nil, err
	}

	numElb := len(output.LoadBalancers)
	if numElb != 1 {
		return nil, fmt.Errorf("Expected to find 1 ELB named '%s' in region '%v', but found '%d'",
			name, region, numElb)
	}

	return output.LoadBalancers[0], nil
}

// CreateElbV2 creates ELB in the given region under the given name and subnets list.
func CreateElbV2(t *testing.T, region string, name string, subnets []*string) *elbv2.LoadBalancer {
	elb, err := CreateElbV2E(t, region, name, subnets)
	if err != nil {
		t.Fatal(err)
	}

	return elb
}

// CreateElbV2E creates ELB in the given region under the given name and subnets list.
func CreateElbV2E(t *testing.T, region string, name string, subnets []*string) (*elbv2.LoadBalancer, error) {
	client := NewElbV2Client(t, region)
	elb, err := client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(name),
		Subnets: subnets,
	})

	if err != nil {
		return nil, err
	}
	return elb.LoadBalancers[0], nil
}

// DeleteElbV2 deletes existing ELB in the given region.
func DeleteElbV2(t *testing.T, region string, elb *elbv2.LoadBalancer) {
	err := DeleteElbV2E(t, region, elb)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteElbV2E deletes existing ELB in the given region.
func DeleteElbV2E(t *testing.T, region string, elb *elbv2.LoadBalancer) error {
	client := NewElbV2Client(t, region)
	_, err := client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(*elb.LoadBalancerArn),
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
