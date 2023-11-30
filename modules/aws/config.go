package aws

import (
	"log"
	"os"
)

var customEndpoints map[string]string

// SetAwsEndpointsOverrides function to override the aws endpoints globall one time
func SetAwsEndpointsOverrides(endpoints map[string]string) {

	if HasAwsCustomEndPoints() {
		log.Println("The aws custom endpoints have been already setup, check your configuration. The global endpoint override is just allow one time")
		os.Exit(1)
	}
	customEndpoints = endpoints
}

// HasAwsCustomEndPoints Function to verify if already has been setup the custom endpoints
func HasAwsCustomEndPoints() bool {
	return len(customEndpoints) > 0
}

// EndpointShouldBeOverride fucntion check if the current service evaulate has custom endpoint to be override
func EndpointShouldBeOverride(service string) (string, bool) {
	customServiceEndpoint, endpointShouldBeOverride := customEndpoints[service]

	return customServiceEndpoint, endpointShouldBeOverride
}
