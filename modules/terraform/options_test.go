package terraform_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionsCloneDeepClonesEnvVars(t *testing.T) {
	t.Parallel()

	unique := random.UniqueID()
	original := terraform.Options{
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

	unique := random.UniqueID()
	original := terraform.Options{
		Vars: map[string]any{
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
		fn   func() (string, error)
		name string
	}{
		{
			name: "apply",
			fn: func() (string, error) {
				return terraform.ApplyE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{Apply: []string{"-help"}}})
			},
		},
		{
			name: "destroy",
			fn: func() (string, error) {
				return terraform.DestroyE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{Destroy: []string{"-help"}}})
			},
		},
		{
			name: "get",
			fn: func() (string, error) {
				return terraform.GetE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{Get: []string{"-help"}}})
			},
		},
		{
			name: "init",
			fn: func() (string, error) {
				return terraform.InitE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{Init: []string{"-help"}}})
			},
		},
		{
			name: "plan",
			fn: func() (string, error) {
				return terraform.PlanE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{Plan: []string{"-help"}}})
			},
		},
		{
			name: "validate",
			fn: func() (string, error) {
				return terraform.ValidateE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{Validate: []string{"-help"}}})
			},
		},
	}

	for _, tt := range testtable {
		out, err := tt.fn()
		require.NoError(t, err)
		assert.Regexp(t, `(Usage|USAGE):\s+\S+\s+(\[global options\]\s+)?`+tt.name, out)
	}
}

func TestExtraArgsWorkspace(t *testing.T) {
	t.Parallel()

	name := t.Name()

	t.Run("New", func(t *testing.T) {
		t.Parallel()

		// set to default
		terraform.WorkspaceSelectOrNew(t, &terraform.Options{}, "default")

		// after adding -help, the function did not create the workspace
		out, err := terraform.WorkspaceSelectOrNewE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{
			WorkspaceNew: []string{"-help"},
		}}, random.UniqueID())
		require.NoError(t, err)
		require.Equal(t, "default", out)
	})

	out, err := terraform.WorkspaceSelectOrNewE(t, &terraform.Options{}, name)
	require.NoError(t, err)
	require.Equal(t, name, out)

	t.Run("Select", func(t *testing.T) {
		t.Parallel()

		// set to default
		terraform.WorkspaceSelectOrNew(t, &terraform.Options{}, "default")

		// after adding -help to select, the function did not select the workspace
		out, err := terraform.WorkspaceSelectOrNewE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{
			WorkspaceSelect: []string{"-help"},
		}}, name)
		require.NoError(t, err)
		require.Equal(t, "default", out)
	})

	t.Run("Delete", func(t *testing.T) {
		t.Parallel()

		// after adding -help to select, the function did not delete the workspace
		_, err := terraform.WorkspaceDeleteE(t, &terraform.Options{ExtraArgs: terraform.ExtraArgs{
			WorkspaceDelete: []string{"-help"},
		}}, name)
		require.NoError(t, err)

		// the workspace should still exist
		out, err := terraform.RunTerraformCommandE(t, &terraform.Options{}, "workspace", "list")
		require.NoError(t, err)
		assert.Contains(t, out, name)
	})
}

func TestOptionsCloneDeepClonesMixedVars(t *testing.T) {
	t.Parallel()

	unique := random.UniqueID()
	original := terraform.Options{
		MixedVars: []terraform.Var{terraform.VarFile(unique), terraform.VarInline("unique", unique)},
	}
	copied, err := original.Clone()
	require.NoError(t, err)

	copied.MixedVars[1] = terraform.VarInline("unique", "nullified")
	assert.Equal(t, terraform.VarFile(unique), copied.MixedVars[0])
	assert.Equal(t, terraform.VarInline("unique", unique), original.MixedVars[1])
}
