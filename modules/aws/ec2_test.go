package aws_test

import (
	"strings"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestGetEc2InstanceIdsByTag(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	ids, err := aws.GetEc2InstanceIdsByTagE(t, region, "Name", "nonexistent-"+random.UniqueID())
	require.NoError(t, err)
	assert.Empty(t, ids)
}

func TestGetEc2InstanceIdsByFilters(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	filters := map[string][]string{
		"instance-state-name": {"running", "shutting-down"},
		"tag:Name":            {"nonexistent-" + random.UniqueID()},
	}

	ids, err := aws.GetEc2InstanceIdsByFiltersE(t, region, filters)
	require.NoError(t, err)
	assert.Empty(t, ids)
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

		t.Run(testCase.region+"-"+strings.Join(testCase.instanceTypeOptions, "-"), func(t *testing.T) {
			t.Parallel()
			instanceType := aws.GetRecommendedInstanceType(t, testCase.region, testCase.instanceTypeOptions)
			// We could hard-code the expected result (e.g., as of July 2020, we expect eu-west-1 to return t2.micro
			// and ap-northeast-2 to return t3.micro), but the result will likely change over time, so to avoid a
			// brittle test, we simply check that we get _one_ result. Combined with the unit test below, this hopefully
			// is enough to be confident this function works correctly.
			assert.Contains(t, testCase.instanceTypeOptions, instanceType)
		})
	}
}

func TestPickRecommendedInstanceTypeHappyPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		expected              string
		availabilityZones     []string
		instanceTypeOptions   []string
		instanceTypeOfferings []types.InstanceTypeOffering
	}{
		{
			name:                  "One AZ, one instance type, available in one offering",
			availabilityZones:     []string{"us-east-1a"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro"}}),
			instanceTypeOptions:   []string{"t2.micro"},
			expected:              "t2.micro",
		},
		{
			name:                  "Three AZs, one instance type, available in all three offerings",
			availabilityZones:     []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}}),
			instanceTypeOptions:   []string{"t2.micro"},
			expected:              "t2.micro",
		},
		{
			name:                  "Three AZs, two instance types, first one available in all three offerings, the other not available at all",
			availabilityZones:     []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}}),
			instanceTypeOptions:   []string{"t2.micro", "t3.micro"},
			expected:              "t2.micro",
		},
		{
			name:                  "Three AZs, two instance types, first one available in all three offerings, the other only available in one offering in an unrequested AZ",
			availabilityZones:     []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}, "us-east-1d": {"t3.micro"}}),
			instanceTypeOptions:   []string{"t2.micro", "t3.micro"},
			expected:              "t2.micro",
		},
		{
			name:                  "Three AZs, two instance types, first one available in all three offerings, the other one available in only two offerings",
			availabilityZones:     []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro", "t3.micro"}, "us-east-1b": {"t2.micro"}, "us-east-1c": {"t2.micro"}}),
			instanceTypeOptions:   []string{"t2.micro", "t3.micro"},
			expected:              "t2.micro",
		},
		{
			name:                  "Three AZs, three instance types, first one available in two offerings, second in all three offerings, third in two offerings",
			availabilityZones:     []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro", "t3.micro", "t3.small"}, "us-east-1b": {"t3.micro"}, "us-east-1c": {"t2.micro", "t3.micro", "t3.small"}}),
			instanceTypeOptions:   []string{"t2.micro", "t3.micro", "t3.small"},
			expected:              "t3.micro",
		},
	}

	for _, testCase := range testCases {
		// The following is necessary to make sure testCase's values don't get updated due to concurrency within the
		// scope of t.Run(..) below. https://golang.org/doc/faq#closures_and_goroutines
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual, err := aws.PickRecommendedInstanceTypeE(testCase.availabilityZones, testCase.instanceTypeOfferings, testCase.instanceTypeOptions)
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestPickRecommendedInstanceTypeErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		availabilityZones     []string
		instanceTypeOfferings []types.InstanceTypeOffering
		instanceTypeOptions   []string
	}{
		{
			name:                  "All params nil",
			availabilityZones:     nil,
			instanceTypeOfferings: nil,
			instanceTypeOptions:   nil,
		},
		{
			name:                  "No AZs, one instance type, no offerings",
			availabilityZones:     nil,
			instanceTypeOfferings: nil,
			instanceTypeOptions:   []string{"t2.micro"},
		},
		{
			name:                  "One AZ, one instance type, no offerings",
			availabilityZones:     []string{"us-east-1a"},
			instanceTypeOfferings: nil,
			instanceTypeOptions:   []string{"t2.micro"},
		},
		{
			name:                  "Two AZs, one instance type, available in only one offering",
			availabilityZones:     []string{"us-east-1a", "us-east-1b"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro"}}),
			instanceTypeOptions:   []string{"t2.micro"},
		},
		{
			name:                  "Three AZs, two instance types, each available in only two of the three offerings",
			availabilityZones:     []string{"us-east-1a", "us-east-1b", "us-east-1c"},
			instanceTypeOfferings: offerings(map[string][]string{"us-east-1a": {"t2.micro"}, "us-east-1b": {"t2.micro", "t3.micro"}, "us-east-1c": {"t3.micro"}}),
			instanceTypeOptions:   []string{"t2.micro", "t3.micro"},
		},
	}

	for _, testCase := range testCases {
		// The following is necessary to make sure testCase's values don't
		// get updated due to concurrency within the scope of t.Run(..) below
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := aws.PickRecommendedInstanceTypeE(testCase.availabilityZones, testCase.instanceTypeOfferings, testCase.instanceTypeOptions)
			assert.EqualError(t, err, aws.NoInstanceTypeError{Azs: testCase.availabilityZones, InstanceTypeOptions: testCase.instanceTypeOptions}.Error())
		})
	}
}

func offerings(offerings map[string][]string) []types.InstanceTypeOffering {
	var out []types.InstanceTypeOffering

	for az, instanceTypes := range offerings {
		for _, instanceType := range instanceTypes {
			offering := types.InstanceTypeOffering{
				InstanceType: types.InstanceType(instanceType),
				Location:     awsSDK.String(az),
				LocationType: types.LocationTypeAvailabilityZone,
			}
			out = append(out, offering)
		}
	}

	return out
}
