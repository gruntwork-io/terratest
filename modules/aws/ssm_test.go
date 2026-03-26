package aws_test

import (
	"testing"

	terraaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestParameterIsFound(t *testing.T) {
	t.Parallel()

	expectedName := "test-name-" + random.UniqueID()
	awsRegion := terraaws.GetRandomRegion(t, nil, nil)
	expectedValue := "test-value-" + random.UniqueID()
	expectedDescription := "test-description-" + random.UniqueID()
	version := terraaws.PutParameter(t, awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Default.Logf(t, "Created parameter with version %d", version)
	keyValue := terraaws.GetParameter(t, awsRegion, expectedName)
	logger.Default.Logf(t, "Found key with name %s", expectedName)
	assert.Equal(t, expectedValue, keyValue)
}

func TestParameterIsDeleted(t *testing.T) {
	t.Parallel()

	expectedName := "test-name-" + random.UniqueID()
	awsRegion := terraaws.GetRandomRegion(t, nil, nil)
	expectedValue := "test-value-" + random.UniqueID()
	expectedDescription := "test-description-" + random.UniqueID()
	version := terraaws.PutParameter(t, awsRegion, expectedName, expectedDescription, expectedValue)
	logger.Default.Logf(t, "Created parameter with version %d", version)

	terraaws.DeleteParameter(t, awsRegion, expectedName)
	logger.Default.Logf(t, "Deleted parameter %s", expectedName)

	actualValue, err := terraaws.GetParameterE(t, awsRegion, expectedName)
	assert.Empty(t, actualValue)
	assert.Error(t, err)
}
