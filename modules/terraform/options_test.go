package terraform

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionsCloneDeepClonesEnvVars(t *testing.T) {
	t.Parallel()

	unique := random.UniqueId()
	original := Options{
		EnvVars: map[string]string{
			"unique":   unique,
			"original": unique,
		},
	}
	copied, err := original.Clone()
	require.NoError(t, err)
	copied.EnvVars["unique"] = "nullified"
	assert.Equal(t, unique, original.EnvVars["unique"])
	assert.Equal(t, unique, copied.EnvVars["original"])
}

func TestOptionsCloneDeepClonesVars(t *testing.T) {
	t.Parallel()

	unique := random.UniqueId()
	original := Options{
		Vars: map[string]interface{}{
			"unique":   unique,
			"original": unique,
		},
	}
	copied, err := original.Clone()
	require.NoError(t, err)
	copied.Vars["unique"] = "nullified"
	assert.Equal(t, unique, original.Vars["unique"])
	assert.Equal(t, unique, copied.Vars["original"])
}
