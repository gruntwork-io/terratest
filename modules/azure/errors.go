package azure

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// SubscriptionIDNotFound is an error that occurs when the Azure Subscription ID could not be found or was not provided.
type SubscriptionIDNotFound struct{}

func (err SubscriptionIDNotFound) Error() string {
	return "could not find an Azure Subscription ID in expected environment variable " +
		AzureSubscriptionID + " and one was not provided for this test"
}

// ResourceGroupNameNotFound is an error that occurs when the target Azure Resource Group name could not be found or was not provided.
type ResourceGroupNameNotFound struct{}

func (err ResourceGroupNameNotFound) Error() string {
	return "could not find an Azure Resource Group name in expected environment variable " +
		AzureResGroupName + " and one was not provided for this test"
}

// FailedToParseError is returned when an object cannot be parsed.
type FailedToParseError struct {
	objectType string
	objectID   string
}

func (err FailedToParseError) Error() string {
	return fmt.Sprintf("failed to parse %s with ID %s", err.objectType, err.objectID)
}

// NewFailedToParseError creates a new not found error when an expected object is not found in the search space.
func NewFailedToParseError(objectType string, objectID string) FailedToParseError {
	return FailedToParseError{objectType: objectType, objectID: objectID}
}

// NotFoundError is returned when an expected object is not found in the search space.
type NotFoundError struct {
	objectType  string
	objectID    string
	searchSpace string
}

func (err NotFoundError) Error() string {
	var objIDMsg string

	if err.objectID != "Any" {
		objIDMsg = " with id " + err.objectID
	}

	return fmt.Sprintf("object of type %s%s not found in %s", err.objectType, objIDMsg, err.searchSpace)
}

// NewNotFoundError creates a new not found error when an expected object is not found in the search space.
func NewNotFoundError(objectType string, objectID string, region string) NotFoundError {
	return NotFoundError{objectType: objectType, objectID: objectID, searchSpace: region}
}

// UnknownEnvironmentError is returned when an Azure environment name is not recognized.
type UnknownEnvironmentError struct {
	EnvironmentName string
}

func (e *UnknownEnvironmentError) Error() string {
	return fmt.Sprintf("unknown Azure environment: %s. "+
		"Available values are: AzurePublicCloud (default), "+
		"AzureUSGovernmentCloud, AzureChinaCloud, or AzureStackCloud",
		e.EnvironmentName)
}

// ResourceNotFoundErrorExists checks the Service Error Code for the 'Resource Not Found' error.
func ResourceNotFoundErrorExists(err error) bool {
	if err == nil {
		return false
	}

	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) {
		return respErr.ErrorCode == "ResourceNotFound" || respErr.ErrorCode == "ResourceGroupNotFound"
	}

	return false
}
