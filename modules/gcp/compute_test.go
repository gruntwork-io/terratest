// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/compute/v1"
)

const DEFAULT_MACHINE_TYPE = "f1-micro"
const DEFAULT_IMAGE_FAMILY_PROJECT_NAME = "ubuntu-os-cloud"
const DEFAULT_IMAGE_FAMILY_NAME = "family/ubuntu-1804-lts"

func TestGetPublicIpOfInstance(t *testing.T) {
	t.Parallel()

	instanceName := RandomValidGcpName()
	projectID := GetGoogleProjectIDFromEnvVar(t)
	zone := GetRandomZone(t, projectID, nil, nil, nil)

	createComputeInstance(t, projectID, zone, instanceName)
	defer deleteComputeInstance(t, projectID, zone, instanceName)

	// Now that our Instance is launched, attempt to query the public IP
	maxRetries := 10
	sleepBetweenRetries := 3 * time.Second

	ip := retry.DoWithRetry(t, "Read IP address of Compute Instance", maxRetries, sleepBetweenRetries, func() (string, error) {
		// Consider attempting to connect to the Compute Instance at this IP in the future, but for now, we just call the
		// the function to ensure we don't have errors
		instance := FetchInstance(t, projectID, instanceName)
		ip := instance.GetPublicIp(t)

		if ip == "" {
			return "", fmt.Errorf("Got blank IP. Retrying.\n")
		}
		return ip, nil
	})

	fmt.Printf("Public IP of Compute Instance %s = %s\n", instanceName, ip)
}

func TestZoneUrlToZone(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		zoneUrl      string
		expectedZone string
	}{
		{"https://www.googleapis.com/compute/v1/projects/terratest-123456/zones/asia-east1-b", "asia-east1-b"},
		{"https://www.googleapis.com/compute/v1/projects/terratest-123456/zones/us-east1-a", "us-east1-a"},
	}

	for _, tc := range testCases {
		zone := ZoneUrlToZone(tc.zoneUrl)
		assert.Equal(t, zone, tc.expectedZone, "Zone not extracted successfully from Zone URL")
	}
}

func TestGetAndSetLabels(t *testing.T) {
	t.Parallel()

	instanceName := RandomValidGcpName()
	projectID := GetGoogleProjectIDFromEnvVar(t)

	// On October 22, 2018, GCP launched the asia-east2 region, which promptly failed all our tests, so blacklist asia-east2.
	zone := GetRandomZone(t, projectID, nil, nil, []string{"asia-east2"})

	createComputeInstance(t, projectID, zone, instanceName)
	defer deleteComputeInstance(t, projectID, zone, instanceName)

	// Now that our Instance is launched, set the labels. Note that in GCP label keys and values can only contain
	// lowercase letters, numeric characters, underscores and dashes.
	instance := FetchInstance(t, projectID, instanceName)

	labelsToWrite := map[string]string{
		"context": "terratest",
	}
	instance.SetLabels(t, labelsToWrite)

	// Now attempt to read the labels we just set.
	maxRetries := 30
	sleepBetweenRetries := 3 * time.Second

	retry.DoWithRetry(t, "Read newly set labels", maxRetries, sleepBetweenRetries, func() (string, error) {
		instance := FetchInstance(t, projectID, instanceName)
		labelsFromRead := instance.GetLabels(t)
		if !reflect.DeepEqual(labelsFromRead, labelsToWrite) {
			return "", fmt.Errorf("Labels that were written did not match labels that were read. Retrying.\n")
		}

		return "", nil
	})
}

// Set custom metadata on a Compute Instance, and then verify it was set as expected
func TestGetAndSetMetadata(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)
	instanceName := RandomValidGcpName()

	// On October 22, 2018, GCP launched the asia-east2 region, which promptly failed all our tests, so blacklist asia-east2.
	zone := GetRandomZone(t, projectID, nil, nil, []string{"asia-east2"})

	// Create a new Compute Instance
	createComputeInstance(t, projectID, zone, instanceName)
	defer deleteComputeInstance(t, projectID, zone, instanceName)

	// Set the metadata
	instance := FetchInstance(t, projectID, instanceName)

	metadataToWrite := map[string]string{
		"foo": "bar",
	}
	instance.SetMetadata(t, metadataToWrite)

	// Now attempt to read the metadata we just set
	maxRetries := 30
	sleepBetweenRetries := 3 * time.Second

	retry.DoWithRetry(t, "Read newly set metadata", maxRetries, sleepBetweenRetries, func() (string, error) {
		instance := FetchInstance(t, projectID, instanceName)
		metadataFromRead := instance.GetMetadata(t)
		for _, metadataItem := range metadataFromRead {
			for key, val := range metadataToWrite {
				if metadataItem.Key == key && *metadataItem.Value == val {
					return "", nil
				}
			}
		}

		fmt.Printf("Metadata to write: %+v\nMetadata from read: %+v\n", metadataToWrite, metadataFromRead)

		return "", fmt.Errorf("Metadata that was written was not found in metadata that was read. Retrying.\n")
	})
}

// Helper function to launch a Compute Instance. This function is useful for quickly iterating on automated tests. But
// if you're writing a test that resembles real-world code that Terratest users may write, you should create a Compute
// Instance using a Terraform apply, similar to the tests in /test.
func createComputeInstance(t *testing.T, projectID string, zone string, name string) {
	t.Logf("Launching new Compute Instance %s\n", name)

	// This RegEx was pulled straight from the GCP API error messages that complained when it's not honored
	validNameExp := `^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$`
	regEx := regexp.MustCompile(validNameExp)

	if !regEx.MatchString(name) {
		t.Fatalf("Invalid Compute Instance name: %s. Must match RegEx %s\n", name, validNameExp)
	}

	machineType := DEFAULT_MACHINE_TYPE
	sourceImageFamilyProjectName := DEFAULT_IMAGE_FAMILY_PROJECT_NAME
	sourceImageFamilyName := DEFAULT_IMAGE_FAMILY_NAME

	// Per GCP docs (https://cloud.google.com/compute/docs/reference/rest/v1/instances/setMachineType), the MachineType
	// is actually specified as a partial URL
	machineTypeURL := fmt.Sprintf("zones/%s/machineTypes/%s", zone, machineType)
	sourceImageURL := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/images/%s", sourceImageFamilyProjectName, sourceImageFamilyName)

	// Based on the properties listed as required at https://cloud.google.com/compute/docs/reference/rest/v1/instances/insert
	// plus a somewhat painful cycle of add-next-property-try-fix-error-message-repeat.
	instanceConfig := &compute.Instance{
		Name:        name,
		MachineType: machineTypeURL,
		NetworkInterfaces: []*compute.NetworkInterface{
			&compute.NetworkInterface{
				AccessConfigs: []*compute.AccessConfig{
					&compute.AccessConfig{},
				},
			},
		},
		Disks: []*compute.AttachedDisk{
			&compute.AttachedDisk{
				AutoDelete: true,
				Boot:       true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: sourceImageURL,
				},
			},
		},
	}

	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Compute Instance
	ctx := context.Background()
	_, err = service.Instances.Insert(projectID, zone, instanceConfig).Context(ctx).Do()
	if err != nil {
		t.Fatalf("Error launching new Compute Instance: %s", err)
	}
}

// Helper function that destroys the given Compute Instance and all of its attached disks.
func deleteComputeInstance(t *testing.T, projectID string, zone string, name string) {
	t.Logf("Deleting Compute Instance %s\n", name)

	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the Compute Instance
	ctx := context.Background()
	_, err = service.Instances.Delete(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		t.Fatalf("Error deleting Compute Instance: %s", err)
	}
}

// TODO: Add additional automated tests to cover remaining functions in compute.go
