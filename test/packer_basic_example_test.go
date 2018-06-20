package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/packer"
)

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerBasicExample(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomRegion(t, nil, nil)

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.json",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"aws_region": awsRegion,
		},
	}

	// Make sure the Packer build completes successfully
	amiID := packer.BuildAmi(t, packerOptions)

	// Clean up the AMI after we're done
	defer aws.DeleteAmiAndAllSnapshots(t, awsRegion, amiID)
}
