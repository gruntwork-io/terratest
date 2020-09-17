package gcp

import (
	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/testing"
)

var credsEnvVars = []string{
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var projectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GOOGLE_CLOUD_PROJECT",
	"GOOGLE_CLOUD_PROJECT_ID",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var regionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

var googleIdentityEmailEnvVars = []string{
	"GOOGLE_IDENTITY_EMAIL",
}

// GetGoogleCredentialsFromEnvVar returns the Credentials for use with testing.
func GetGoogleCredentialsFromEnvVar(t testing.TestingT) string {
	return environment.GetFirstNonEmptyEnvVarOrEmptyString(t, credsEnvVars)
}

// GetGoogleProjectIDFromEnvVar returns the Project Id for use with testing.
func GetGoogleProjectIDFromEnvVar(t testing.TestingT) string {
	return environment.GetFirstNonEmptyEnvVarOrFatal(t, projectEnvVars)
}

// GetGoogleRegionFromEnvVar returns the Region for use with testing.
func GetGoogleRegionFromEnvVar(t testing.TestingT) string {
	return environment.GetFirstNonEmptyEnvVarOrFatal(t, regionEnvVars)
}

// GetGoogleIdentityEmailEnvVar returns a Google identity (user) for use with testing.
func GetGoogleIdentityEmailEnvVar(t testing.TestingT) string {
	return environment.GetFirstNonEmptyEnvVarOrFatal(t, googleIdentityEmailEnvVars)
}
