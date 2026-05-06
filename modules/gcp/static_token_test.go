package gcp_test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GOOGLE_OAUTH_ACCESS_TOKEN short-circuits the credential lookup in the client
// constructors. When it is set, the service constructors must build a client
// without touching GOOGLE_APPLICATION_CREDENTIALS, which lets CI exercise the
// auth wiring without real GCP creds.
func TestStaticTokenShortCircuitsCredentialLookup(t *testing.T) {
	// Poison the ADC path so any fallback to default credentials would fail loudly.
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/nonexistent/credentials.json")
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "fake-access-token")

	svc, err := gcp.NewComputeServiceE(t)
	require.NoError(t, err)
	assert.NotNil(t, svc)

	cb, err := gcp.NewCloudBuildServiceE(t)
	require.NoError(t, err)
	require.NotNil(t, cb)
	require.NoError(t, cb.Close())

	osl, err := gcp.NewOSLoginServiceE(t)
	require.NoError(t, err)
	assert.NotNil(t, osl)
}

// Without the static token env var, constructors must fall back to ADC, which
// fails when the credentials file is poisoned. t.Setenv records the original
// value and registers cleanup; the subsequent Unsetenv triggers the "not set"
// branch in getStaticTokenSource for the duration of the test.
func TestMissingStaticTokenFallsBackToADC(t *testing.T) {
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", "")
	require.NoError(t, os.Unsetenv("GOOGLE_OAUTH_ACCESS_TOKEN"))

	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/nonexistent/credentials.json")

	_, err := gcp.NewCloudBuildServiceE(t)
	assert.Error(t, err)
}
