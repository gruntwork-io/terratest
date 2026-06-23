package azure

import (
	"strings"
)

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
	i := strings.LastIndex(resourceID, "/")
	if i == -1 {
		return "", NewResourceIDNameNotFoundError(resourceID)
	}

	return resourceID[i+1:], nil
}
