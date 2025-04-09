package aws

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretsManagerMethods(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	name := random.UniqueId()
	description := "This is just a secrets manager test description."
	secretOriginalValue := "This is the secret value."
	secretUpdatedValue := "This is the NEW secret value."

	secretARN := CreateSecretStringWithDefaultKey(t, region, description, name, secretOriginalValue)
	defer deleteSecret(t, region, secretARN)

	storedValue := GetSecretValue(t, region, secretARN)
	assert.Equal(t, secretOriginalValue, storedValue)

	PutSecretString(t, region, secretARN, secretUpdatedValue)

	storedValueAfterUpdate := GetSecretValue(t, region, secretARN)
	assert.Equal(t, secretUpdatedValue, storedValueAfterUpdate)
}

func deleteSecret(t *testing.T, region, id string) {
	DeleteSecret(t, region, id, true)

	_, err := GetSecretValueE(t, region, id)
	require.Error(t, err)
}
