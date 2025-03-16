package terraform

import (
	"fmt"
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

func TestExtraArgsHelp(t *testing.T) {
	t.Parallel()

	testtable := []struct {
		name string
		fn   func() (string, error)
	}{
		{
			name: "apply",
			fn:   func() (string, error) { return ApplyE(t, &Options{ExtraArgs: ExtraArgs{Apply: []string{"-help"}}}) },
		},
		{
			name: "destroy",
			fn:   func() (string, error) { return DestroyE(t, &Options{ExtraArgs: ExtraArgs{Destroy: []string{"-help"}}}) },
		},
		{
			name: "get",
			fn:   func() (string, error) { return GetE(t, &Options{ExtraArgs: ExtraArgs{Get: []string{"-help"}}}) },
		},
		{
			name: "init",
			fn:   func() (string, error) { return InitE(t, &Options{ExtraArgs: ExtraArgs{Init: []string{"-help"}}}) },
		},
		{
			name: "plan",
			fn:   func() (string, error) { return PlanE(t, &Options{ExtraArgs: ExtraArgs{Plan: []string{"-help"}}}) },
		},
		{
			name: "validate",
			fn: func() (string, error) {
				return ValidateE(t, &Options{ExtraArgs: ExtraArgs{Validate: []string{"-help"}}})
			},
		},
	}

	for _, tt := range testtable {
		out, err := tt.fn()
		require.NoError(t, err)
		assert.Contains(t, out, fmt.Sprintf("Usage: terraform [global options] %s", tt.name))
	}
}
