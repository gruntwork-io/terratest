package azure

import (
	"fmt"
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
	if !strings.Contains(resourceID, "/") {
		return "", fmt.Errorf("could not resolve name from resource ID %q", resourceID)
	}

	parts := strings.Split(resourceID, "/")

	return parts[len(parts)-1], nil
}
