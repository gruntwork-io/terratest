package aws

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestParameterIsFound(t *testing.T) {
	t.Parallel()

	expectedName := fmt.Sprintf("test-name-%s", random.UniqueId())
	awsRegion := GetRandomRegion(t, nil, nil)
	expectedValue := fmt.Sprintf("test-value-%s", random.UniqueId())
	expectedDescription := fmt.Sprintf("test-description-%s", random.UniqueId())
	version := PutParameter(t, awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Default.Logf(t, "Created parameter with version %d", version)
	keyValue := GetParameter(t, awsRegion, expectedName)
	logger.Default.Logf(t, "Found key with name %s", expectedName)
	assert.Equal(t, expectedValue, keyValue)
}

func TestParameterIsDeleted(t *testing.T) {
	expectedName := fmt.Sprintf("test-name-%s", random.UniqueId())
	awsRegion := GetRandomRegion(t, nil, nil)
	expectedValue := fmt.Sprintf("test-value-%s", random.UniqueId())
	expectedDescription := fmt.Sprintf("test-description-%s", random.UniqueId())
	version := PutParameter(t, awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Default.Logf(t, "Created parameter with version %d", version)

	DeleteParameter(t, awsRegion, expectedName)
	logger.Default.Logf(t, "Deleted paramter %s", expectedName)

	actualValue, err := GetParameterE(t, awsRegion, expectedName)
	assert.Equal(t, actualValue, "")
	assert.Error(t, err)
}
