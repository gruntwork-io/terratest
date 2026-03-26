package aws_test

import (
	"testing"

	terraaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretsManagerMethods(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegion(t, nil, nil)
	name := random.UniqueID()
	description := "This is just a secrets manager test description."
	secretOriginalValue := "This is the secret value."
	secretUpdatedValue := "This is the NEW secret value."

	secretARN := terraaws.CreateSecretStringWithDefaultKey(t, region, description, name, secretOriginalValue)
	defer deleteSecret(t, region, secretARN)

	storedValue := terraaws.GetSecretValue(t, region, secretARN)
	assert.Equal(t, secretOriginalValue, storedValue)

	terraaws.PutSecretString(t, region, secretARN, secretUpdatedValue)

	storedValueAfterUpdate := terraaws.GetSecretValue(t, region, secretARN)
	assert.Equal(t, secretUpdatedValue, storedValueAfterUpdate)
}

func deleteSecret(t *testing.T, region, id string) {
	t.Helper()

	terraaws.DeleteSecret(t, region, id, true)

	_, err := terraaws.GetSecretValueE(t, region, id)
	require.Error(t, err)
}
