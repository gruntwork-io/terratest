package docker_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
)

func TestListImagesAndDeleteImage(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueID())
	repo := "gruntwork-io/test-image"
	tag := "v1-" + uniqueID
	img := fmt.Sprintf("%s:%s", repo, tag)

	options := &docker.BuildOptions{
		Tags: []string{img},
	}

	ctx := t.Context()
	docker.BuildContext(t, ctx, "../../test/fixtures/docker", options)

	assert.True(t, docker.DoesImageExistContext(t, ctx, img, nil))
	docker.DeleteImageContext(t, ctx, img, nil)
	assert.False(t, docker.DoesImageExistContext(t, ctx, img, nil))
}
