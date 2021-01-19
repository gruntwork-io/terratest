// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

func TestCreateBuild(t *testing.T) {
	t.Parallel()
	// This test performs the following steps:
	//
	// 1. Creates a tarball with a single Dockerfile
	// 2. Creates a GCS bucket
	// 3. Uploads the tarball to the GCS Bucket
	// 4. Triggers a build using the Cloud Build API
	// 5. Untags and deletes all pushed Build images
	// 6. Deletes the GCS bucket

	// Create and add some files to the archive.
	tarball := createSampleAppTarball(t)

	// Create GCS bucket
	projectID := GetGoogleProjectIDFromEnvVar(t)
	id := random.UniqueId()
	gsBucketName := "cloud-build-terratest-" + strings.ToLower(id)
	sampleAppPath := "docker-example.tar.gz"
	imagePath := fmt.Sprintf("gcr.io/%s/test-image-%s", projectID, strings.ToLower(id))

	logger.Logf(t, "Random values selected Bucket Name = %s\n", gsBucketName)

	CreateStorageBucket(t, projectID, gsBucketName, nil)
	defer DeleteStorageBucket(t, gsBucketName)

	// Write the compressed archive to the storage bucket
	objectURL := WriteBucketObject(t, gsBucketName, sampleAppPath, tarball, "application/gzip")
	logger.Logf(t, "Got URL: %s", objectURL)

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
	b := CreateBuild(t, projectID, build)

	// Delete the pushed build images
	// We could just use the `b` struct above, but we want to explicitly test
	// the `GetBuild` method.
	b2 := GetBuild(t, projectID, b.GetId())
	for _, image := range b2.GetImages() {
		DeleteGCRRepo(t, image)
	}

	// Empty the storage bucket so we can delete it
	defer EmptyStorageBucket(t, gsBucketName)
}

func createSampleAppTarball(t *testing.T) *bytes.Reader {
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
