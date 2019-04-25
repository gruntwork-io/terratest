package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/stretchr/testify/require"
)

// LoadBalancer is an Amazon load balancer.
type LoadBalancer struct {
	Name                  string // The name of the load balancer.
	ARN                   string // The Amazon Resource Name (ARN) of the load balancer.
	CanonicalHostedZoneID string // The ID of the Amazon Route 53 hosted zone associated with the load balancer.
	DNSName               string // The public DNS name of the load balancer.
}

// CreateElbV2 creates ELB in the given region under the given name and subnets list.
func CreateElbV2(t *testing.T, region string, name string, subnets []*string) *LoadBalancer {
	lb, err := CreateElbV2E(t, region, name, subnets)
	if err != nil {
		t.Fatal(err)
	}
	return lb
}

// CreateElbV2E creates ELB in the given region under the given name and subnets list.
func CreateElbV2E(t *testing.T, region string, name string, subnets []*string) (*LoadBalancer, error) {
	client := NewElbV2Client(t, region)
	resp, err := client.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:    aws.String(name),
		Subnets: subnets,
	})
	if err != nil {
		return nil, err
	}

	numLb := len(resp.LoadBalancers)
	if numLb != 1 {
		return nil, fmt.Errorf("Expected to create 1 ELB named '%s' in region '%v', but found '%d'",
			name, region, numLb)
	}
	lb := resp.LoadBalancers[0]

	return &LoadBalancer{
		Name: aws.StringValue(lb.LoadBalancerName),
		ARN:  aws.StringValue(lb.LoadBalancerArn),
		CanonicalHostedZoneID: aws.StringValue(lb.CanonicalHostedZoneId),
		DNSName:               aws.StringValue(lb.DNSName),
	}, nil
}

// DeleteElbV2 deletes existing ELB in the given region.
func DeleteElbV2(t *testing.T, region string, lb *LoadBalancer) {
	err := DeleteElbV2E(t, region, lb)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteElbV2E deletes existing ELB in the given region.
func DeleteElbV2E(t *testing.T, region string, lb *LoadBalancer) error {
	client := NewElbV2Client(t, region)
	_, err := client.DeleteLoadBalancer(&elbv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(lb.ARN),
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
