//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/packer"
)

// Occasionally, a Packer build may fail due to intermittent issues (e.g., brief network outage or EC2 issue). We try
// to make our tests resilient to that by specifying those known common errors here and telling our builds to retry if
// they hit those errors.
var DefaultRetryablePackerErrors = map[string]string{
	"Script disconnected unexpectedly":                                                 "Occasionally, Packer seems to lose connectivity to AWS, perhaps due to a brief network outage",
	"can not open /var/lib/apt/lists/archive.ubuntu.com_ubuntu_dists_xenial_InRelease": "Occasionally, apt-get fails on ubuntu to update the cache",
}
var DefaultTimeBetweenPackerRetries = 15 * time.Second

// Regions that support running f1-micro instances
var RegionsThatSupportF1Micro = []string{"us-central1", "us-east1", "us-west1", "europe-west1"}

// Zones that support running f1-micro instances
var ZonesThatSupportF1Micro = []string{"us-central1-a", "us-east1-b", "us-west1-a", "europe-north1-a", "europe-west1-b", "europe-central2-a"}

const DefaultMaxPackerRetries = 3

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerGCPBasicExample(t *testing.T) {
	t.Parallel()

	// Get the Project Id to use
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	// Pick a random GCP zone to test in. This helps ensure your code works in all regions.
	zone := gcp.GetRandomZone(t, projectID, ZonesThatSupportF1Micro, nil, nil)

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../../examples/packer-basic-example/build-gcp.pkr.hcl",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"gcp_project_id": projectID,
			"gcp_zone":       zone,
		},

		// Only build the Google Compute Image
		Only: "googlecompute.ubuntu-bionic",

		// Configure retries for intermittent errors
		RetryableErrors:    DefaultRetryablePackerErrors,
		TimeBetweenRetries: DefaultTimeBetweenPackerRetries,
		MaxRetries:         DefaultMaxPackerRetries,
	}

	// Make sure the Packer build completes successfully
	imageName := packer.BuildArtifact(t, packerOptions)

	// Delete the Image after we're done
	image := gcp.FetchImage(t, projectID, imageName)
	defer image.DeleteImage(t)
}
