package gcp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	gcrname "github.com/google/go-containerregistry/pkg/name"
	gcrgoogle "github.com/google/go-containerregistry/pkg/v1/google"
	gcrremote "github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GCRClient bundles the credentials and transport the go-containerregistry library needs.
// It exists so that GCR helpers can follow terratest's (ctx, client, ...args) WithClient convention
// instead of exposing the two pieces separately. Transport is optional — nil uses the library default.
type GCRClient struct {
	Authenticator authn.Authenticator
	Transport     http.RoundTripper
}

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
	authenticator, err := newGCRAuthenticator() //nolint:contextcheck // newGCRAuthenticator is a pure credential helper
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	logger.Default.Logf(t, "Retrieving Image Digests %s", repo)

	return DeleteGCRRepoWithClient(ctx, &GCRClient{Authenticator: authenticator}, repo)
}

// DeleteGCRRepoWithClient deletes a GCR repository including all tagged images using the supplied
// client. Prefer this variant in unit tests where the client's Transport points at an httptest
// fake server (see gcr_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DeleteGCRRepoWithClient(ctx context.Context, client *GCRClient, repo string) error {
	gcrrepo, err := gcrname.NewRepository(repo)
	if err != nil {
		return fmt.Errorf("failed to get repo: %w", err)
	}

	listOpts := []gcrgoogle.Option{gcrgoogle.WithAuth(client.Authenticator), gcrgoogle.WithContext(ctx)}
	if client.Transport != nil {
		listOpts = append(listOpts, gcrgoogle.WithTransport(client.Transport))
	}

	tags, err := gcrgoogle.List(gcrrepo, listOpts...)
	if err != nil {
		return fmt.Errorf("failed to list tags for repo %s: %w", repo, err)
	}

	latestRef := repo + ":latest"
	if err := DeleteGCRImageRefWithClient(ctx, client, latestRef); err != nil {
		return fmt.Errorf("failed to delete GCR image reference %s: %w", latestRef, err)
	}

	for k := range tags.Manifests {
		ref := repo + "@" + k
		if err := DeleteGCRImageRefWithClient(ctx, client, ref); err != nil {
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
	authenticator, err := newGCRAuthenticator() //nolint:contextcheck // newGCRAuthenticator is a pure credential helper
	if err != nil {
		return fmt.Errorf("failed to create authenticator: %w", err)
	}

	logger.Default.Logf(t, "Deleting Image Ref %s", ref)

	return DeleteGCRImageRefWithClient(ctx, &GCRClient{Authenticator: authenticator}, ref)
}

// DeleteGCRImageRefWithClient deletes a single repo image ref/digest using the supplied client.
// Prefer this variant in unit tests where the client's Transport points at an httptest fake
// server (see gcr_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func DeleteGCRImageRefWithClient(ctx context.Context, client *GCRClient, ref string) error {
	name, err := gcrname.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("failed to parse reference %s: %w", ref, err)
	}

	opts := []gcrremote.Option{gcrremote.WithAuth(client.Authenticator), gcrremote.WithContext(ctx)}
	if client.Transport != nil {
		opts = append(opts, gcrremote.WithTransport(client.Transport))
	}

	if err := gcrremote.Delete(name, opts...); err != nil {
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
