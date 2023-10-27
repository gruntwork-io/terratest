package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	terratest_aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Occasionally, a Packer build may fail due to intermittent issues (e.g., brief network outage or EC2 issue). We try
// to make our tests resilient to that by specifying those known common errors here and telling our builds to retry if
// they hit those errors.
var DefaultRetryablePackerErrors = map[string]string{
	"Script disconnected unexpectedly":                                                 "Occasionally, Packer seems to lose connectivity to AWS, perhaps due to a brief network outage",
	"can not open /var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_xenial_InRelease": "Occasionally, apt-get fails on ubuntu to update the cache",
}
var DefaultTimeBetweenPackerRetries = 15 * time.Second

const DefaultMaxPackerRetries = 3

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerBasicExample(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := terratest_aws.GetRandomStableRegion(t, nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := terratest_aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// website::tag::1::Read Packer's template and set AWS Region variable.
	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.pkr.hcl",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"aws_region":    awsRegion,
			"ami_base_name": fmt.Sprintf("%s", random.UniqueId()),
			"instance_type": instanceType,
		},

		// Only build the AWS AMI
		Only: "amazon-ebs.ubuntu-example",

		// Configure retries for intermittent errors
		RetryableErrors:    DefaultRetryablePackerErrors,
		TimeBetweenRetries: DefaultTimeBetweenPackerRetries,
		MaxRetries:         DefaultMaxPackerRetries,
	}

	// website::tag::2::Build artifacts from Packer's template.
	// Make sure the Packer build completes successfully
	amiID := packer.BuildArtifact(t, packerOptions)

	// website::tag::4::Remove AMI after test.
	// Clean up the AMI after we're done
	defer terratest_aws.DeleteAmiAndAllSnapshots(t, awsRegion, amiID)

	// Check if AMI is shared/not shared with account
	requestingAccount := terratest_aws.CanonicalAccountId
	randomAccount := "123456789012" // Random Account
	ec2Client := terratest_aws.NewEc2Client(t, awsRegion)
	ShareAmi(t, amiID, requestingAccount, ec2Client)
	accountsWithLaunchPermissions := terratest_aws.GetAccountsWithLaunchPermissionsForAmi(t, awsRegion, amiID)
	assert.NotContains(t, accountsWithLaunchPermissions, randomAccount)
	assert.Contains(t, accountsWithLaunchPermissions, requestingAccount)

	// website::tag::3::Check AMI's properties.
	// Check if AMI is public
	MakeAmiPublic(t, amiID, ec2Client)
	amiIsPublic := terratest_aws.GetAmiPubliclyAccessible(t, awsRegion, amiID)
	assert.True(t, amiIsPublic)
}

// An example of how to test the Packer template in examples/packer-basic-example using Terratest
// with the VarFiles option. This test generates a temporary *.json file containing the value
// for the `aws_region` variable.
func TestPackerBasicExampleWithVarFile(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := terratest_aws.GetRandomStableRegion(t, nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := terratest_aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	// Create temporary packer variable file to store aws region
	varFile, err := os.CreateTemp("", "*.json")
	require.NoError(t, err, "Did not expect temp file creation to cause error")

	// Be sure to clean up temp file
	defer os.Remove(varFile.Name())

	// Write the vars we need to a temporary json file
	varFileContent := []byte(fmt.Sprintf(`{"aws_region": "%s", "instance_type": "%s"}`, awsRegion, instanceType))
	_, err = varFile.Write(varFileContent)
	require.NoError(t, err, "Did not expect writing to temp file %s to cause error", varFile.Name())

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.pkr.hcl",

		// Variable file to to pass to our Packer build using -var-file option
		VarFiles: []string{
			varFile.Name(),
		},

		// Environment settings to avoid plugin conflicts
		Env: map[string]string{
			"PACKER_PLUGIN_PATH": "../examples/packer-basic-example/.packer.d/plugins",
		},

		// Only build the AWS AMI
		Only: "amazon-ebs.ubuntu-example",

		// Configure retries for intermittent errors
		RetryableErrors:    DefaultRetryablePackerErrors,
		TimeBetweenRetries: DefaultTimeBetweenPackerRetries,
		MaxRetries:         DefaultMaxPackerRetries,
	}

	// Make sure the Packer build completes successfully
	amiID := packer.BuildArtifact(t, packerOptions)

	// Clean up the AMI after we're done
	defer terratest_aws.DeleteAmiAndAllSnapshots(t, awsRegion, amiID)

	// Check if AMI is shared/not shared with account
	requestingAccount := terratest_aws.CanonicalAccountId
	randomAccount := "123456789012" // Random Account
	ec2Client := terratest_aws.NewEc2Client(t, awsRegion)
	ShareAmi(t, amiID, requestingAccount, ec2Client)
	accountsWithLaunchPermissions := terratest_aws.GetAccountsWithLaunchPermissionsForAmi(t, awsRegion, amiID)
	assert.NotContains(t, accountsWithLaunchPermissions, randomAccount)
	assert.Contains(t, accountsWithLaunchPermissions, requestingAccount)

	// Check if AMI is public
	MakeAmiPublic(t, amiID, ec2Client)
	amiIsPublic := terratest_aws.GetAmiPubliclyAccessible(t, awsRegion, amiID)
	assert.True(t, amiIsPublic)
}

func TestPackerMultipleConcurrentAmis(t *testing.T) {
	t.Parallel()

	// Build a map of 3 randomId <-> packer.Options, in 3 random AWS Regions
	// then build all of these AMIs in parallel and make sure that there are
	// no errors.
	var identifierToOptions = map[string]*packer.Options{}
	for i := 0; i < 3; i++ {
		// Pick a random AWS region to test in. This helps ensure your code works in all regions.
		awsRegion := terratest_aws.GetRandomStableRegion(t, nil, nil)

		// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
		instanceType := terratest_aws.GetRecommendedInstanceType(t, awsRegion, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

		packerOptions := &packer.Options{
			// The path to where the Packer template is located
			Template: "../examples/packer-basic-example/build.pkr.hcl",

			// Variables to pass to our Packer build using -var options
			Vars: map[string]string{
				"aws_region":    awsRegion,
				"ami_base_name": fmt.Sprintf("%s", random.UniqueId()),
				"instance_type": instanceType,
			},

			// Only build the AWS AMI
			Only: "amazon-ebs.ubuntu-example",

			// Configure retries for intermittent errors
			RetryableErrors:    DefaultRetryablePackerErrors,
			TimeBetweenRetries: DefaultTimeBetweenPackerRetries,
			MaxRetries:         DefaultMaxPackerRetries,
		}

		identifierToOptions[random.UniqueId()] = packerOptions
	}

	resultMap := packer.BuildArtifacts(t, identifierToOptions)

	// Clean up the AMIs after we're done
	for key, amiId := range resultMap {
		awsRegion := identifierToOptions[key].Vars["aws_region"]
		terratest_aws.DeleteAmiAndAllSnapshots(t, awsRegion, amiId)
	}
}

func ShareAmi(t *testing.T, amiID string, accountID string, ec2Client *ec2.EC2) {
	input := &ec2.ModifyImageAttributeInput{
		ImageId: aws.String(amiID),
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: []*ec2.LaunchPermission{
				{
					UserId: aws.String(accountID),
				},
			},
		},
	}
	_, err := ec2Client.ModifyImageAttribute(input)
	if err != nil {
		t.Fatal(err)
	}
}

func MakeAmiPublic(t *testing.T, amiID string, ec2Client *ec2.EC2) {
	input := &ec2.ModifyImageAttributeInput{
		ImageId: aws.String(amiID),
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: []*ec2.LaunchPermission{
				{
					Group: aws.String("all"),
				},
			},
		},
	}
	_, err := ec2Client.ModifyImageAttribute(input)
	if err != nil {
		t.Fatal(err)
	}
}
