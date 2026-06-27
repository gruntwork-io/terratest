package environment

import (
	"os"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// GetFirstNonEmptyEnvVarOrFatal returns the first non-empty environment variable from envVarNames, or fails the test if all of them are empty.
//
// Deprecated: scheduled for removal in Terratest v2. Range over os.Getenv at the call
// site, e.g.:
//
//	var value string
//	for _, name := range envVarNames {
//		if value = os.Getenv(name); value != "" {
//			break
//		}
//	}
//	require.NotEmptyf(t, value, "all of %v are empty", envVarNames)
func GetFirstNonEmptyEnvVarOrFatal(t testing.TestingT, envVarNames []string) string {
	value := GetFirstNonEmptyEnvVarOrEmptyString(t, envVarNames)
	if value == "" {
		t.Fatalf("All of the following env vars %v are empty. At least one must be non-empty.", envVarNames)
	}

	return value
}

// GetFirstNonEmptyEnvVarOrEmptyString returns the first non-empty environment variable from envVarNames, or returns the
// empty string
//
// Deprecated: scheduled for removal in Terratest v2. Range over os.Getenv at the call
// site, e.g.:
//
//	for _, name := range envVarNames {
//		if value := os.Getenv(name); value != "" {
//			return value
//		}
//	}
func GetFirstNonEmptyEnvVarOrEmptyString(t testing.TestingT, envVarNames []string) string {
	for _, name := range envVarNames {
		if value := os.Getenv(name); value != "" {
			return value
		}
	}

	return ""
}

// RequireEnvVar fails the test if the specified environment variable is not defined or is blank.
//
// Deprecated: scheduled for removal in Terratest v2. Use require with os.Getenv at the
// call site, e.g.:
//
//	require.NotEmptyf(t, os.Getenv(name), "environment variable %s must be set", name)
func RequireEnvVar(t testing.TestingT, envVarName string) {
	require.NotEmptyf(t, os.Getenv(envVarName), "Environment variable %s must be set for this test.", envVarName)
}
