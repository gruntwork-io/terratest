package docker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestListImagesAndDeleteImage(t *testing.T) {
	t.Parallel()

	uniqueID := strings.ToLower(random.UniqueId())
	repo := "gruntwork-io/test-image"
	tag := fmt.Sprintf("v1-%s", uniqueID)
	img := fmt.Sprintf("%s:%s", repo, tag)

	options := &BuildOptions{
		Tags: []string{img},
	}
	Build(t, "../../test/fixtures/docker", options)

	assert.True(t, DoesImageExist(t, img, nil))
	DeleteImage(t, img, nil)
	assert.False(t, DoesImageExist(t, img, nil))
}
