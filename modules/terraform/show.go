package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// ShowContext calls terraform show in json mode with the given options and returns stdout from the command. If
// PlanFilePath is set on the options, this will show the plan file. Otherwise, this will show the current state of the
// terraform module at options.TerraformDir. The context argument can be used for cancellation or timeout control.
// This will fail the test if there is an error in the command.
func ShowContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ShowContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ShowContextE calls terraform show in json mode with the given options and returns stdout from the command. If
// PlanFilePath is set on the options, this will show the plan file. Otherwise, this will show the current state of the
// terraform module at options.TerraformDir. The context argument can be used for cancellation or timeout control.
func ShowContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {

	args := []string{"show", "-no-color", "-json"}

	if options.PlanFilePath != "" {
		args = append(args, options.PlanFilePath)
	}

	return RunTerraformCommandAndGetStdoutContextE(t, ctx, options, prepend(options.ExtraArgs.Show, args...)...)
}

// ShowWithStructContext calls terraform show in json mode with the given options, parses the json result into a
// PlanStruct, and returns it. If PlanFilePath is set on the options, this will show the plan file. Otherwise, this
// will show the current state of the terraform module at options.TerraformDir. The context argument can be used for
// cancellation or timeout control. This will fail the test if there is an error in the command.
func ShowWithStructContext(t testing.TestingT, ctx context.Context, options *Options) *PlanStruct {
	out, err := ShowWithStructContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ShowWithStructContextE calls terraform show in json mode with the given options, parses the json result into a
// PlanStruct, and returns it. If PlanFilePath is set on the options, this will show the plan file. Otherwise, this
// will show the current state of the terraform module at options.TerraformDir. The context argument can be used for
// cancellation or timeout control.
func ShowWithStructContextE(t testing.TestingT, ctx context.Context, options *Options) (*PlanStruct, error) {
	json, err := ShowContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	planStruct, err := ParsePlanJSON(json)
	if err != nil {
		return nil, err
	}

	return planStruct, nil
}
