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
	secretValue := "This is the secret value."

	secretARN := CreateSecretStringWithDefaultKey(t, region, description, name, secretValue)
	defer deleteSecret(t, region, secretARN)

	storedValue := GetSecretValue(t, region, secretARN)
	assert.Equal(t, secretValue, storedValue)
}

func deleteSecret(t *testing.T, region, id string) {
	DeleteSecret(t, region, id, true)

	_, err := GetSecretValueE(t, region, id)
	require.Error(t, err)
}
