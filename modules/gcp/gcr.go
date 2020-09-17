package gcp

import (
	"fmt"

	gcrname "github.com/google/go-containerregistry/pkg/name"
	gcrgoogle "github.com/google/go-containerregistry/pkg/v1/google"
	gcrremote "github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DeleteGCRRepo deletes a GCR repository including all tagged images
func DeleteGCRRepo(t testing.TestingT, repo string) {
	err := DeleteGCRRepoE(t, repo)
	require.NoError(t, err)
}

// DeleteGCRRepoE deletes a GCR repository including all tagged images
func DeleteGCRRepoE(t testing.TestingT, repo string) error {
	// create a new auther for the API calls
	auther, err := gcrgoogle.NewEnvAuthenticator()
	if err != nil {
		return fmt.Errorf("Failed to create auther. Got error: %v", err)
	}

	gcrrepo, err := gcrname.NewRepository(repo)
	if err != nil {
		return fmt.Errorf("Failed to get repo. Got error: %v", err)
	}

	logger.Logf(t, "Retriving Image Digests %s", gcrrepo)
	tags, err := gcrgoogle.List(gcrrepo, gcrgoogle.WithAuth(auther))
	if err != nil {
		return fmt.Errorf("Failed to list tags for repo %s. Got error: %v", repo, err)
	}

	// attempt to delete the latest image tag
	latestRef := repo + ":latest"
	logger.Logf(t, "Deleting Image Ref %s", latestRef)
	if err := DeleteGCRImageRefE(t, latestRef); err != nil {
		return fmt.Errorf("Failed to delete GCR Image Reference %s. Got error: %v", latestRef, err)
	}

	// delete image references sequentially
	for k := range tags.Manifests {
		ref := repo + "@" + k
		logger.Logf(t, "Deleting Image Ref %s", ref)

		if err := DeleteGCRImageRefE(t, ref); err != nil {
			return fmt.Errorf("Failed to delete GCR Image Reference %s. Got error: %v", ref, err)
		}
	}

	return nil
}

// DeleteGCRImageRef deletes a single repo image ref/digest
func DeleteGCRImageRef(t testing.TestingT, ref string) {
	err := DeleteGCRImageRefE(t, ref)
	require.NoError(t, err)
}

// DeleteGCRImageRefE deletes a single repo image ref/digest
func DeleteGCRImageRefE(t testing.TestingT, ref string) error {
	name, err := gcrname.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("Failed to parse reference %s. Got error: %v", ref, err)
	}

	// create a new auther for the API calls
	auther, err := gcrgoogle.NewEnvAuthenticator()
	if err != nil {
		return fmt.Errorf("Failed to create auther. Got error: %v", err)
	}

	opts := gcrremote.WithAuth(auther)

	if err := gcrremote.Delete(name, opts); err != nil {
		return fmt.Errorf("Failed to delete %s. Got error: %v", name, err)
	}

	return nil
}
