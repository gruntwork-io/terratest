package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetEksClusterE fetches information about an EKS cluster.
func GetEksClusterE(t testing.TestingT, region string, name string) (*eks.Cluster, error) {
	client, err := NewEksClientE(t, region)
	if err != nil {
		return nil, err
	}
	input := &eks.DescribeClusterInput{
		Name: aws.String(name),
	}
	output, err := client.DescribeCluster(input)
	if err != nil {
		return nil, err
	}
	return output.Cluster, nil
}

// GetEksCluster fetches information about an EKS cluster.
func GetEksCluster(t testing.TestingT, region string, name string) *eks.Cluster {
	cluster, err := GetEksClusterE(t, region, name)
	require.NoError(t, err)
	return cluster
}

// CreateEksClusterE creates EKS cluster in the given region under the given name.
func CreateEksClusterE(t testing.TestingT, region string, name string, roleArn string, enablePrivateAccess bool, enablePublicAccess bool, subnets []*string, publicAccessCidrs []*string, securityGroupIds []*string) (*eks.Cluster, error) {
	client := NewEksClient(t, region)
	cluster, err := client.CreateCluster(&eks.CreateClusterInput{
		Name: aws.String(name),
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			EndpointPublicAccess:  &enablePublicAccess,
			EndpointPrivateAccess: &enablePrivateAccess,
			SubnetIds:             subnets,
			SecurityGroupIds:      securityGroupIds,
			PublicAccessCidrs:     publicAccessCidrs,
		},
		RoleArn: aws.String(roleArn),
	})

	if err != nil {
		return nil, err
	}
	return cluster.Cluster, nil
}

// CreateEksCluster creates EKS cluster in the given region under the given name.
func CreateEksCluster(t testing.TestingT, region string, name string, roleArn string, enablePrivateAccess bool, enablePublicAccess bool, subnets []*string, publicAccessCidrs []*string, securityGroupIds []*string) *eks.Cluster {
	cluster, err := CreateEksClusterE(t, region, name, roleArn, enablePrivateAccess, enablePublicAccess, subnets, publicAccessCidrs, securityGroupIds)
	require.NoError(t, err)
	return cluster
}

// DeleteEksClusterE deletes existing EKS cluster in the given region.
func DeleteEksClusterE(t testing.TestingT, region string, cluster *eks.Cluster) error {
	client := NewEksClient(t, region)
	_, err := client.DeleteCluster(&eks.DeleteClusterInput{
		Name: aws.String(*cluster.Name),
	})
	return err
}

// DeleteEksCluster deletes existing EKS cluster in the given region.
func DeleteEksCluster(t testing.TestingT, region string, cluster *eks.Cluster) {
	err := DeleteEksClusterE(t, region, cluster)
	require.NoError(t, err)
}

// NewEksClient creates en EKS client.
func NewEksClient(t testing.TestingT, region string) *eks.EKS {
	client, err := NewEksClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewEcsClientE creates an ECS client.
func NewEksClientE(t testing.TestingT, region string) (*eks.EKS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return eks.New(sess), nil
}
