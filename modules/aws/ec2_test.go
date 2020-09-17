package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEc2InstanceIdsByTag(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ids, err := GetEc2InstanceIdsByTagE(t, region, "Name", fmt.Sprintf("nonexistent-%s", random.UniqueId()))
	require.NoError(t, err)
	assert.Equal(t, 0, len(ids))
}

func TestGetEc2InstanceIdsByFilters(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	filters := map[string][]string{
		"instance-state-name": {"running", "shutting-down"},
		"tag:Name":            {fmt.Sprintf("nonexistent-%s", random.UniqueId())},
	}

	ids, err := GetEc2InstanceIdsByFiltersE(t, region, filters)
	require.NoError(t, err)
	assert.Equal(t, 0, len(ids))
}

func TestGetRecommendedInstanceType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		region              string
		instanceTypeOptions []string
	}{
		{"eu-west-1", []string{"t2.micro", "t3.micro"}},
		{"ap-northeast-2", []string{"t2.micro", "t3.micro"}},
		{"us-east-1", []string{"t2.large", "t3.large"}},
	}

	for _, testCase := range testCases {
		// The following is necessary to make sure testCase's values don't get updated due to concurrency within the
		// scope of t.Run(..) below. https://golang.org/doc/faq#closures_and_goroutines
		testCase := testCase

		t.Run(fmt.Sprintf("%s-%s", testCase.region, strings.Join(testCase.instanceTypeOptions, "-")), func(t *testing.T) {
			t.Parallel()
			instanceType := GetRecommendedInstanceType(t, testCase.region, testCase.instanceTypeOptions)
			// We could hard-code the expected result (e.g., as of July, 2020, we expect eu-west-1 to return t2.micro
			// and ap-northeast-2 to return t3.micro), but the result will likely change over time, so to avoid a
			// brittle test, we simply check that we get _one_ result. Combined with the unit test below, this hopefully
			// is enough to be confident this function works correctly.
			assert.Contains(t, testCase.instanceTypeOptions, instanceType)
		})
	}
}

func TestPickRecommendedInstanceTypeHappyPath(t *testing.T) {
	testCases := []struct {
		name                  string
		availabilityZones     []string
		instanceTypeOfferings []*ec2.InstanceTypeOffering
		instanceTypeOptions   []string
		expected              string
	}{
		{
			"One AZ, one instance type, available in one offering",
			[]string{"us-east-1a"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro"}}),
			[]string{"t2.micro"},
			"t2.micro",
		},
		{
			"Three AZs, one instance type, available in all three offerings",
			[]string{"us-east-1a", "us-east-1b", "us-east-1c"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}}),
			[]string{"t2.micro"},
			"t2.micro",
		},
		{
			"Three AZs, two instance types, first one available in all three offerings, the other not available at all",
			[]string{"us-east-1a", "us-east-1b", "us-east-1c"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}}),
			[]string{"t2.micro", "t3.micro"},
			"t2.micro",
		},
		{
			"Three AZs, two instance types, first one available in all three offerings, the other only available in one offering in an unrequested AZ",
			[]string{"us-east-1a", "us-east-1b", "us-east-1c"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}, "us-east-1d": {"t3.micro"}}),
			[]string{"t2.micro", "t3.micro"},
			"t2.micro",
		},
		{
			"Three AZs, two instance types, first one available in all three offerings, the other one available in only two offerings",
			[]string{"us-east-1a", "us-east-1b", "us-east-1c"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro", "t3.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}}),
			[]string{"t2.micro", "t3.micro"},
			"t2.micro",
		},
		{
			"Three AZs, three instance types, first one available in two offerings, second in all three offerings, third in two offerings",
			[]string{"us-east-1a", "us-east-1b", "us-east-1c"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro", "t3.micro", "t3.small"}, "us-east-1b": {"t3.micro"}, "us-east-1c": {"t2.micro", "t3.micro", "t3.small"}}),
			[]string{"t2.micro", "t3.micro", "t3.small"},
			"t3.micro",
		},
	}

	for _, testCase := range testCases {
		// The following is necessary to make sure testCase's values don't get updated due to concurrency within the
		// scope of t.Run(..) below. https://golang.org/doc/faq#closures_and_goroutines
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual, err := pickRecommendedInstanceTypeE(testCase.availabilityZones, testCase.instanceTypeOfferings, testCase.instanceTypeOptions)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestPickRecommendedInstanceTypeErrors(t *testing.T) {
	testCases := []struct {
		name                  string
		availabilityZones     []string
		instanceTypeOfferings []*ec2.InstanceTypeOffering
		instanceTypeOptions   []string
	}{
		{
			"All params nil",
			nil,
			nil,
			nil,
		},
		{
			"No AZs, one instance type, no offerings",
			nil,
			nil,
			[]string{"t2.micro"},
		},
		{
			"One AZ, one instance type, no offerings",
			[]string{"us-east-1a"},
			nil,
			[]string{"t2.micro"},
		},
		{
			"Two AZs, one instance type, available in only one offering",
			[]string{"us-east-1a", "us-east-1b"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro"}}),
			[]string{"t2.micro"},
		},
		{
			"Three AZs, two instance types, each available in only two of the three offerings",
			[]string{"us-east-1a", "us-east-1b", "us-east-1c"},
			offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro", "t3.micro"}, "us-east-1c": {"t3.micro"}}),
			[]string{"t2.micro", "t3.micro"},
		},
	}

	for _, testCase := range testCases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := pickRecommendedInstanceTypeE(testCase.availabilityZones, testCase.instanceTypeOfferings, testCase.instanceTypeOptions)
			assert.EqualError(t, err, NoInstanceTypeError{Azs: testCase.availabilityZones, InstanceTypeOptions: testCase.instanceTypeOptions}.Error())
		})
	}
}

func offerings(offerings map[string][]string) []*ec2.InstanceTypeOffering {
	var out []*ec2.InstanceTypeOffering

	for az, instanceTypes := range offerings {
		for _, instanceType := range instanceTypes {
			offering := &ec2.InstanceTypeOffering{
				InstanceType: aws.String(instanceType),
				Location:     aws.String(az),
				LocationType: aws.String(ec2.LocationTypeAvailabilityZone),
			}
			out = append(out, offering)
		}
	}

	return out
}
