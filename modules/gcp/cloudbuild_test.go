//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
	"testing"

	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestCreateBuild(t *testing.T) {
	t.Parallel()
	// This test performs the following steps:
	//
	// 1. Creates a tarball with a single Dockerfile
	// 2. Creates a GCS bucket
	// 3. Uploads the tarball to the GCS Bucket
	// 4. Triggers a build using the Cloud Build API
	// 5. Attempts to untag and delete all pushed Build images (best-effort cleanup)
	// 6. Deletes the GCS bucket

	// Create and add some files to the archive.
	tarball := createSampleAppTarball(t)

	// Create GCS bucket
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)
	id := random.UniqueID()
	gsBucketName := "cloud-build-terratest-" + strings.ToLower(id)
	sampleAppPath := "docker-example.tar.gz"
	imagePath := fmt.Sprintf("gcr.io/%s/test-image-%s", projectID, strings.ToLower(id))

	logger.Default.Logf(t, "Random values selected Bucket Name = %s\n", gsBucketName)

	gcp.CreateStorageBucket(t, projectID, gsBucketName, nil)
	defer gcp.DeleteStorageBucket(t, gsBucketName)

	// Write the compressed archive to the storage bucket
	objectURL := gcp.WriteBucketObject(t, gsBucketName, sampleAppPath, tarball, "application/gzip")
	logger.Default.Logf(t, "Got URL: %s", objectURL)

	// Create a new build
	build := &cloudbuildpb.Build{
		Source: &cloudbuildpb.Source{
			Source: &cloudbuildpb.Source_StorageSource{
				StorageSource: &cloudbuildpb.StorageSource{
					Bucket: gsBucketName,
					Object: sampleAppPath,
				},
			},
		},
		Steps: []*cloudbuildpb.BuildStep{{
			Name: "gcr.io/cloud-builders/docker",
			Args: []string{"build", "-t", imagePath, "."},
		}},
		Images: []string{imagePath},
	}

	// CreateBuild blocks until the build is complete
	b := gcp.CreateBuild(t, projectID, build)

	// Attempt to delete the pushed build images (best-effort cleanup).
	// Note: GCR (gcr.io) has been deprecated in favor of Artifact Registry.
	// The cleanup may fail due to permission changes, but this doesn't affect
	// the validity of the Cloud Build test itself.
	// We could just use the `b` struct above, but we want to explicitly test
	// the `GetBuild` method.
	b2 := gcp.GetBuild(t, projectID, b.GetId())
	for _, image := range b2.GetImages() {
		if err := gcp.DeleteGCRRepoE(t, image); err != nil {
			logger.Default.Logf(t, "Warning: Failed to delete image %s (this may be expected due to GCR deprecation): %v", image, err)
		}
	}

	// Empty the storage bucket so we can delete it
	defer gcp.EmptyStorageBucket(t, gsBucketName)
}

func createSampleAppTarball(t *testing.T) *bytes.Reader {
	t.Helper()

	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)

	file := `FROM busybox:latest
MAINTAINER Rob Morgan (rob@gruntwork.io)
	`

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(file)),
	}

	err := tw.WriteHeader(hdr)
	require.NoError(t, err)

	_, werr := tw.Write([]byte(file))
	require.NoError(t, werr)

	cerr := tw.Close()
	require.NoError(t, cerr)

	// gzip the tar archive
	var zbuf bytes.Buffer

	gzw := gzip.NewWriter(&zbuf)
	_, gwerr := gzw.Write(buf.Bytes())
	require.NoError(t, gwerr)

	gcerr := gzw.Close()
	require.NoError(t, gcerr)

	// return the compressed buffer
	return bytes.NewReader(zbuf.Bytes())
}
