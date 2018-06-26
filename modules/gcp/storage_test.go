package gcp

import (
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

var (
	projectId = os.Getenv("GOOGLE_STORAGE_PROJECT_ID")
)

func TestCreateAndDestroyStorageBucket(t *testing.T) {
	t.Parallel()

	id := random.UniqueId()
	logger.Logf(t, "Random values selected Id = %s\n", id)

	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateStorageBucket(t, projectId, gsBucketName, nil)
	DeleteStorageBucket(t, gsBucketName)
}
