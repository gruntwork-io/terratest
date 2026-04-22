// Package oci allows you to interact with Oracle Cloud Infrastructure (OCI) resources.
package oci

import (
	"context"
	"sort"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/stretchr/testify/require"
)

// DeleteImageContextE deletes a custom image with given OCID.
// The ctx parameter supports cancellation and timeouts.
func DeleteImageContextE(t testing.TestingT, ctx context.Context, ocid string) error {
	logger.Default.Logf(t, "Deleting image with OCID %s", ocid)

	configProvider := common.DefaultConfigProvider()

	client, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		return err
	}

	request := core.DeleteImageRequest{ImageId: &ocid}
	_, err = client.DeleteImage(ctx, request)

	return err
}

// DeleteImageContext deletes a custom image with given OCID.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteImageContext(t testing.TestingT, ctx context.Context, ocid string) {
	t.Helper()

	err := DeleteImageContextE(t, ctx, ocid)
	require.NoError(t, err)
}

// DeleteImage deletes a custom image with given OCID.
//
// Deprecated: Use [DeleteImageContext] instead.
func DeleteImage(t testing.TestingT, ocid string) {
	t.Helper()

	DeleteImageContext(t, context.Background(), ocid)
}

// DeleteImageE deletes a custom image with given OCID.
//
// Deprecated: Use [DeleteImageContextE] instead.
func DeleteImageE(t testing.TestingT, ocid string) error {
	return DeleteImageContextE(t, context.Background(), ocid)
}

// GetMostRecentImageIDContextE gets the OCID of the most recent image in the given compartment that has the given OS name and version.
// The ctx parameter supports cancellation and timeouts.
func GetMostRecentImageIDContextE(t testing.TestingT, ctx context.Context, compartmentID string, osName string, osVersion string) (string, error) {
	configProvider := common.DefaultConfigProvider()

	client, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		return "", err
	}

	request := core.ListImagesRequest{
		CompartmentId:          &compartmentID,
		OperatingSystem:        &osName,
		OperatingSystemVersion: &osVersion,
	}

	var allItems []core.Image

	for {
		response, err := client.ListImages(ctx, request)
		if err != nil {
			return "", err
		}

		allItems = append(allItems, response.Items...)

		// Stop when no next page, when the server returns an empty token, or
		// when it returns the same token we just requested (defensive: prevents
		// an infinite loop on a misbehaving server).
		if response.OpcNextPage == nil || *response.OpcNextPage == "" {
			break
		}

		if request.Page != nil && *request.Page == *response.OpcNextPage {
			break
		}

		request.Page = response.OpcNextPage
	}

	if len(allItems) == 0 {
		return "", NoImagesFoundError{OSName: osName, OSVersion: osVersion, CompartmentID: compartmentID}
	}

	mostRecentImage := mostRecentImage(allItems)

	return *mostRecentImage.Id, nil
}

// GetMostRecentImageIDContext gets the OCID of the most recent image in the given compartment that has the given OS name and version.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMostRecentImageIDContext(t testing.TestingT, ctx context.Context, compartmentID string, osName string, osVersion string) string {
	t.Helper()

	ocid, err := GetMostRecentImageIDContextE(t, ctx, compartmentID, osName, osVersion)
	require.NoError(t, err)

	return ocid
}

// GetMostRecentImageID gets the OCID of the most recent image in the given compartment that has the given OS name and version.
//
// Deprecated: Use [GetMostRecentImageIDContext] instead.
func GetMostRecentImageID(t testing.TestingT, compartmentID string, osName string, osVersion string) string {
	t.Helper()

	return GetMostRecentImageIDContext(t, context.Background(), compartmentID, osName, osVersion)
}

// GetMostRecentImageIDE gets the OCID of the most recent image in the given compartment that has the given OS name and version.
//
// Deprecated: Use [GetMostRecentImageIDContextE] instead.
func GetMostRecentImageIDE(t testing.TestingT, compartmentID string, osName string, osVersion string) (string, error) {
	return GetMostRecentImageIDContextE(t, context.Background(), compartmentID, osName, osVersion)
}

// Image sorting code borrowed from: https://github.com/hashicorp/packer/blob/7f4112ba229309cfc0ebaa10ded2abdfaf1b22c8/builder/amazon/common/step_source_ami_info.go
type imageSort []core.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	switch {
	case a[i].TimeCreated == nil:
		return true
	case a[j].TimeCreated == nil:
		return false
	default:
		return a[i].TimeCreated.Unix() < a[j].TimeCreated.Unix()
	}
}

// mostRecentImage returns the most recent image out of a slice of images.
func mostRecentImage(images []core.Image) core.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))

	return sortedImages[len(sortedImages)-1]
}
