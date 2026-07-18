package aws_test

import (
	"context"
	"testing"

	terraaws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretsManagerMethods(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	name := random.UniqueID()
	description := "This is just a secrets manager test description."
	secretOriginalValue := "This is the secret value."
	secretUpdatedValue := "This is the NEW secret value."

	secretARN := terraaws.CreateSecretStringWithDefaultKeyContext(t, context.Background(), region, description, name, secretOriginalValue)
	defer deleteSecret(t, region, secretARN)

	storedValue := terraaws.GetSecretValueContext(t, context.Background(), region, secretARN)
	assert.Equal(t, secretOriginalValue, storedValue)

	terraaws.PutSecretStringContext(t, context.Background(), region, secretARN, secretUpdatedValue)

	storedValueAfterUpdate := terraaws.GetSecretValueContext(t, context.Background(), region, secretARN)
	assert.Equal(t, secretUpdatedValue, storedValueAfterUpdate)
}

func deleteSecret(t *testing.T, region, id string) {
	t.Helper()

	terraaws.DeleteSecretContext(t, context.Background(), region, id, true)

	_, err := terraaws.GetSecretValueContextE(t, context.Background(), region, id)
	require.Error(t, err)
}
