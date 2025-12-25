package gcp

import (
	"context"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2/google"
)

func TestStaticTokenClient(t *testing.T) {
	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
	require.NoError(t, err)
	token, err := creds.TokenSource.Token()
	require.NoError(t, err)
	projectID := GetGoogleProjectIDFromEnvVar(t)

	// we poison the default client instantiation with invalid file so that if it is used, it fails
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "non-existent-credentials.json")
	_, err = NewCloudBuildServiceE(t)
	require.Error(t, err)
	_, err = NewComputeServiceE(t)
	require.Error(t, err)
	_, err = newGCRAuther()
	require.Error(t, err)
	_, err = NewOSLoginServiceE(t)
	require.Error(t, err)
	_, err = newStorageClient()
	require.Error(t, err)

	// now we instantiate client with oauth2 token
	// and run several function to make sure the new client is correctly configured with access token
	t.Setenv("GOOGLE_OAUTH_ACCESS_TOKEN", token.AccessToken)
	GetAllGcpRegions(t, projectID)
	GetBuilds(t, projectID)
	GetLoginProfile(t, GetGoogleIdentityEmailEnvVar(t))
	_, err = newGCRAuther()
	require.NoError(t, err)
	bucket := "gruntwork-terratest-" + strings.ToLower(random.UniqueId())
	CreateStorageBucket(t, projectID, bucket, nil)
	defer DeleteStorageBucket(t, bucket)
}
