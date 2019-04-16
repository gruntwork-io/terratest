package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/stretchr/testify/require"
)

// GetElb fetches information about specified ELB.
func GetElb(t *testing.T, region string, name string) *elbv2.LoadBalancer {
	elb, err := GetElbE(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
	return elb
}

// GetElbE fetches information about specified ELB.
func GetElbE(t *testing.T, region string, name string) (*elbv2.LoadBalancer, error) {
	client := NewElbClient(t, region)

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

// CreateElb creates ELB in the given region under the given name and subnets list.
func CreateElb(t *testing.T, region string, name string, subnets []*string) *elbv2.LoadBalancer {
	elb, err := CreateElbE(t, region, name, subnets)
	if err != nil {
		t.Fatal(err)
	}

	return elb
}

// CreateElbE creates ELB in the given region under the given name and subnets list.
func CreateElbE(t *testing.T, region string, name string, subnets []*string) (*elbv2.LoadBalancer, error) {
	client := NewElbClient(t, region)
	elb, err := client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(name),
		Subnets: subnets,
	})

	if err != nil {
		return nil, err
	}
	return elb.LoadBalancers[0], nil
}

// DeleteElb deletes existing ELB in the given region.
func DeleteElb(t *testing.T, region string, elb *elbv2.LoadBalancer) {
	err := DeleteElbE(t, region, elb)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteElbE deletes existing ELB in the given region.
func DeleteElbE(t *testing.T, region string, elb *elbv2.LoadBalancer) error {
	client := NewElbClient(t, region)
	_, err := client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(*elb.LoadBalancerArn),
	})

	return err
}

// NewElbClient creates en ELB client.
func NewElbClient(t *testing.T, region string) *elbv2.ELBV2 {
	client, err := NewElbClientE(t, region)
	require.NoError(t, err)

	return client
}

// NewElbClientE creates an ELB client.
func NewElbClientE(t *testing.T, region string) (*elbv2.ELBV2, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return elbv2.New(sess), nil
}
