package oci

import "fmt"

// NoImagesFoundError is an error that occurs when no images matching the given criteria are found in a compartment.
type NoImagesFoundError struct {
	OSName        string
	OSVersion     string
	CompartmentID string
}

func (err NoImagesFoundError) Error() string {
	return fmt.Sprintf("no %s %s images found in the %s compartment", err.OSName, err.OSVersion, err.CompartmentID)
}

// NoAvailabilityDomainsFoundError is an error that occurs when no availability domains are found in a compartment.
type NoAvailabilityDomainsFoundError struct {
	CompartmentID string
}

func (err NoAvailabilityDomainsFoundError) Error() string {
	return fmt.Sprintf("no availability domains found in the %s compartment", err.CompartmentID)
}

// NoVCNsFoundError is an error that occurs when no VCNs are found in a compartment.
type NoVCNsFoundError struct {
	CompartmentID string
}

func (err NoVCNsFoundError) Error() string {
	return fmt.Sprintf("no VCNs found in the %s compartment", err.CompartmentID)
}
