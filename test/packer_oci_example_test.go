package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/oci"
	"github.com/gruntwork-io/terratest/modules/packer"
)

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerOciExample(t *testing.T) {
	t.Parallel()

	// The Terratest CI environment does not yet have CI creds set up, so we skip these tests for now
	// https://github.com/gruntwork-io/terratest/issues/160
	if os.Getenv("CIRCLECI") != "" {
		t.Skip("The build is running on CircleCI, so skipping OCI tests.")
	}

	compartmentID := oci.GetRootCompartmentID(t)
	baseImageID := oci.GetMostRecentImageID(t, compartmentID, "Canonical Ubuntu", "18.04")
	availabilityDomain := oci.GetRandomAvailabilityDomain(t, compartmentID)
	subnetID := oci.GetRandomSubnetID(t, compartmentID, availabilityDomain)
	passPhrase := oci.GetPassPhraseFromEnvVar()

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.pkr.hcl",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"oci_compartment_ocid":    compartmentID,
			"oci_base_image_ocid":     baseImageID,
			"oci_availability_domain": availabilityDomain,
			"oci_subnet_ocid":         subnetID,
			"oci_pass_phrase":         passPhrase,
		},

		// Only build an OCI image
		Only: "oracle-oci",

		// Configure retries for intermittent errors
		RetryableErrors:    DefaultRetryablePackerErrors,
		TimeBetweenRetries: DefaultTimeBetweenPackerRetries,
		MaxRetries:         DefaultMaxPackerRetries,
	}

	// Make sure the Packer build completes successfully
	ocid := packer.BuildArtifact(t, packerOptions)

	// Delete the OCI image after we're done
	defer oci.DeleteImage(t, ocid)
}
