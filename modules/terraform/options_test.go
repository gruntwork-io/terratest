package terraform

import (
	"fmt"
	"regexp"
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
		{
			name: "validate-inputs",
			fn: func() (string, error) {
				return ValidateInputsE(t, &Options{
					ExtraArgs: ExtraArgs{ValidateInputs: []string{"-help"}}, TerraformBinary: "terragrunt"})
			},
		},
	}

	for _, tt := range testtable {
		out, err := tt.fn()
		require.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile(fmt.Sprintf(`(Usage|USAGE):\s+\S+\s+(\[global options\]\s+)?%s`, tt.name)), out)
	}
}

func TestExtraArgsWorkspace(t *testing.T) {
	name := t.Name()

	t.Run("New", func(t *testing.T) {
		// set to default
		WorkspaceSelectOrNew(t, &Options{}, "default")

		// after adding -help, the function did not create the workspace
		out, err := WorkspaceSelectOrNewE(t, &Options{ExtraArgs: ExtraArgs{
			WorkspaceNew: []string{"-help"},
		}}, random.UniqueId())
		require.NoError(t, err)
		require.Equal(t, "default", out)
	})

	out, err := WorkspaceSelectOrNewE(t, &Options{}, name)
	require.NoError(t, err)
	require.Equal(t, name, out)
	t.Run("Select", func(t *testing.T) {
		// set to default
		WorkspaceSelectOrNew(t, &Options{}, "default")

		// after adding -help to select, the function did not select the workspace
		out, err := WorkspaceSelectOrNewE(t, &Options{ExtraArgs: ExtraArgs{
			WorkspaceSelect: []string{"-help"},
		}}, name)
		require.NoError(t, err)
		require.Equal(t, "default", out)
	})

	t.Run("Delete", func(t *testing.T) {
		// after adding -help to select, the function did not delete the workspace
		_, err := WorkspaceDeleteE(t, &Options{ExtraArgs: ExtraArgs{
			WorkspaceDelete: []string{"-help"},
		}}, name)
		require.NoError(t, err)

		// the workspace should still exist
		out, err := RunTerraformCommandE(t, &Options{}, "workspace", "list")
		require.NoError(t, err)
		assert.Contains(t, out, name)
	})
}

func TestOptionsCloneDeepClonesMixedVars(t *testing.T) {
	t.Parallel()

	unique := random.UniqueId()
	original := Options{
		MixedVars: []Var{VarFile(unique), VarInline("unique", unique)},
	}
	copied, err := original.Clone()
	require.NoError(t, err)
	copied.MixedVars[1] = VarInline("unique", "nullified")
	assert.Equal(t, VarFile(unique), copied.MixedVars[0])
	assert.Equal(t, VarInline("unique", unique), original.MixedVars[1])
}
