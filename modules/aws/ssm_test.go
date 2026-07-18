package aws_test

import (
	"context"
	"testing"

	terraaws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
)

func TestParameterIsFound(t *testing.T) {
	t.Parallel()

	expectedName := "test-name-" + random.UniqueID()
	awsRegion := terraaws.GetRandomRegionContext(t, context.Background(), nil, nil)
	expectedValue := "test-value-" + random.UniqueID()
	expectedDescription := "test-description-" + random.UniqueID()
	version := terraaws.PutParameterContext(t, context.Background(), awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Default.Logf(t, "Created parameter with version %d", version)
	keyValue := terraaws.GetParameterContext(t, context.Background(), awsRegion, expectedName)
	logger.Default.Logf(t, "Found key with name %s", expectedName)
	assert.Equal(t, expectedValue, keyValue)
}

func TestParameterIsDeleted(t *testing.T) {
	t.Parallel()

	expectedName := "test-name-" + random.UniqueID()
	awsRegion := terraaws.GetRandomRegionContext(t, context.Background(), nil, nil)
	expectedValue := "test-value-" + random.UniqueID()
	expectedDescription := "test-description-" + random.UniqueID()
	version := terraaws.PutParameterContext(t, context.Background(), awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Default.Logf(t, "Created parameter with version %d", version)

	terraaws.DeleteParameterContext(t, context.Background(), awsRegion, expectedName)
	logger.Default.Logf(t, "Deleted parameter %s", expectedName)

	actualValue, err := terraaws.GetParameterContextE(t, context.Background(), awsRegion, expectedName)
	assert.Empty(t, actualValue)
	assert.Error(t, err)
}
