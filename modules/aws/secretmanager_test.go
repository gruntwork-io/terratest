package aws

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestSecretIsFound(t *testing.T) {
	t.Parallel()

	expectedName := fmt.Sprintf("test-name-%s", random.UniqueId())
	awsRegion := GetRandomRegion(t, nil, nil)
	expectedValue := fmt.Sprintf("test-value-%s", random.UniqueId())
	expectedDescription := fmt.Sprintf("test-description-%s", random.UniqueId())
	version := PutSecret(t, awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Logf(t, "Created secret with version %s", version)
	keyName := GetSecret(t, awsRegion, expectedName)
	logger.Logf(t, "Found secret with name %s", expectedName)
	assert.Equal(t, expectedName, keyName)

}
