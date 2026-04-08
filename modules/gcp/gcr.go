package gcp

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	gcrname "github.com/google/go-containerregistry/pkg/name"
	gcrgoogle "github.com/google/go-containerregistry/pkg/v1/google"
	gcrremote "github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DeleteGCRRepo deletes a GCR repository including all tagged images.
// This will fail the test if there is an error.
//
// Deprecated: Use [DeleteGCRRepoContext] instead.
func DeleteGCRRepo(t testing.TestingT, repo string) {
	DeleteGCRRepoContext(t, context.Background(), repo)
}

// DeleteGCRRepoContext deletes a GCR repository including all tagged images.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteGCRRepoContext(t testing.TestingT, ctx context.Context, repo string) {
	err := DeleteGCRRepoContextE(t, ctx, repo)
	require.NoError(t, err)
}

// DeleteGCRRepoE deletes a GCR repository including all tagged images.
//
// Deprecated: Use [DeleteGCRRepoContextE] instead.
func DeleteGCRRepoE(t testing.TestingT, repo string) error {
	return DeleteGCRRepoContextE(t, context.Background(), repo)
}

// DeleteGCRRepoContextE deletes a GCR repository including all tagged images.
// The ctx parameter supports cancellation and timeouts.
func DeleteGCRRepoContextE(t testing.TestingT, ctx context.Context, repo string) error {
	// create a new authenticator for the API calls
	authenticator, err := newGCRAuthenticator() //nolint:contextcheck // newGCRAuthenticator is a pure credential helper
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	gcrrepo, err := gcrname.NewRepository(repo)
	if err != nil {
		return fmt.Errorf("failed to get repo: %w", err)
	}

	logger.Default.Logf(t, "Retrieving Image Digests %s", gcrrepo)

	tags, err := gcrgoogle.List(gcrrepo, gcrgoogle.WithAuth(authenticator), gcrgoogle.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to list tags for repo %s: %w", repo, err)
	}

	// attempt to delete the latest image tag
	latestRef := repo + ":latest"
	logger.Default.Logf(t, "Deleting Image Ref %s", latestRef)

	if err := DeleteGCRImageRefContextE(t, ctx, latestRef); err != nil {
		return fmt.Errorf("failed to delete GCR image reference %s: %w", latestRef, err)
	}

	// delete image references sequentially
	for k := range tags.Manifests {
		ref := repo + "@" + k
		logger.Default.Logf(t, "Deleting Image Ref %s", ref)

		if err := DeleteGCRImageRefContextE(t, ctx, ref); err != nil {
			return fmt.Errorf("failed to delete GCR image reference %s: %w", ref, err)
		}
	}

	return nil
}

// DeleteGCRImageRef deletes a single repo image ref/digest.
// This will fail the test if there is an error.
//
// Deprecated: Use [DeleteGCRImageRefContext] instead.
func DeleteGCRImageRef(t testing.TestingT, ref string) {
	DeleteGCRImageRefContext(t, context.Background(), ref)
}

// DeleteGCRImageRefContext deletes a single repo image ref/digest.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteGCRImageRefContext(t testing.TestingT, ctx context.Context, ref string) {
	err := DeleteGCRImageRefContextE(t, ctx, ref)
	require.NoError(t, err)
}

// DeleteGCRImageRefE deletes a single repo image ref/digest.
//
// Deprecated: Use [DeleteGCRImageRefContextE] instead.
func DeleteGCRImageRefE(t testing.TestingT, ref string) error {
	return DeleteGCRImageRefContextE(t, context.Background(), ref)
}

// DeleteGCRImageRefContextE deletes a single repo image ref/digest.
// The ctx parameter supports cancellation and timeouts.
func DeleteGCRImageRefContextE(t testing.TestingT, ctx context.Context, ref string) error {
	name, err := gcrname.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", ref, err)
	}

	// create a new authenticator for the API calls
	authenticator, err := newGCRAuthenticator() //nolint:contextcheck // newGCRAuthenticator is a pure credential helper
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	if err := gcrremote.Delete(name, gcrremote.WithAuth(authenticator), gcrremote.WithContext(ctx)); err != nil {
		return fmt.Errorf("failed to delete %s: %w", name, err)
	}

	return nil
}

func newGCRAuthenticator() (authn.Authenticator, error) {
	if ts, ok := getStaticTokenSource(); ok {
		return gcrgoogle.NewTokenSourceAuthenticator(ts), nil
	}

	return gcrgoogle.NewEnvAuthenticator(context.Background())
}
