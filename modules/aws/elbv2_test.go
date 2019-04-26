package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// getSubnetIdsPerAZ gets the subnets ids in a slice given the specified VPC.
func getSubnetIdsPerAZ(t *testing.T, vpc *Vpc) []string {
	var subnetIds []string

	lookUp := make(map[string]bool)
	for _, subnet := range vpc.Subnets {
		_, ok := lookUp[subnet.AvailabilityZone]
		if ok {
			continue
		}
		lookUp[subnet.AvailabilityZone] = true
		subnetIds = append(subnetIds, subnet.Id)
	}

	return subnetIds
}

func TestElbV2(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	vpc := GetDefaultVpc(t, region)
	subnets := getSubnetIdsPerAZ(t, vpc) // To create ELB you must specify subnets from at least two Availability Zones.
	elbName := "terratest"

	lb, err := CreateElbV2E(t, region, elbName, subnets)
	defer DeleteElbV2(t, region, lb)

	assert.Nil(t, err)
}
