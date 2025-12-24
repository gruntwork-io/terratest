package azure

import "github.com/gruntwork-io/terratest/modules/collections/v2"

// GetNameFromResourceID gets the Name from an Azure Resource ID.
func GetNameFromResourceID(resourceID string) string {
	id, err := GetNameFromResourceIDE(resourceID)
	if err != nil {
		return ""
	}
	return id
}

// GetNameFromResourceIDE gets the Name from an Azure Resource ID.
// This function would fail the test if there is an error.
func GetNameFromResourceIDE(resourceID string) (string, error) {
	id, err := collections.GetSliceLastValueE(resourceID, "/")
	if err != nil {
		return "", err
	}
	return id, nil
}
