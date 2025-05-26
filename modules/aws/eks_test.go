package aws

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/stretchr/testify/assert"
)

func TestEksCluster(t *testing.T) {
	t.Parallel()

	// Setting the cluster to always run in us-east-1 and only the us-east-1a, us-east-1b, and us-east-1d subnets as there
	// are seemingly randon subnet restrictions for EKS clusters.
	region := GetRandomStableRegion(t, []string{"us-east-1"}, nil)
	vpc, err := GetDefaultVpcE(t, region)
	assert.Nil(t, err)

	var subnetList []*string
	for _, subnet := range vpc.Subnets {
		if subnet.AvailabilityZone == "us-east-1a" || subnet.AvailabilityZone == "us-east-1b" || subnet.AvailabilityZone == "us-east-1d" {
			subnetList = append(subnetList, aws.String(subnet.Id))
		}
	}

	iamClient, err := NewIamClientE(t, region)
	assert.Nil(t, err)

	role, err := iamClient.GetRole(&iam.GetRoleInput{
		RoleName: aws.String("AWSServiceRoleForOrganizations"),
	})
	assert.Nil(t, err)

	clusterName := "terratest"
	c1, err := CreateEksClusterE(t, region, clusterName, *role.Role.Arn, true, true, subnetList, []*string{aws.String("0.0.0.0/0")}, []*string{})
	defer DeleteEksCluster(t, region, c1)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *c1.Name)

	maxRetries := 60
	sleepBetweenRetries := 30 * time.Second

	status := retry.DoWithRetry(t, "Ensure cluster is active", maxRetries, sleepBetweenRetries, func() (string, error) {
		cluster := GetEksCluster(t, region, clusterName)
		status := *cluster.Status

		if status != "ACTIVE" {
			return "", fmt.Errorf("Got Cluster Status %s. Retrying.\n", *cluster.Status)
		}
		return status, nil
	})

	assert.Equal(t, "ACTIVE", status)

	c2, err := GetEksClusterE(t, region, *c1.Name)

	assert.Nil(t, err)
	assert.Equal(t, "terratest", *c2.Name)
}
