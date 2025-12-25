package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/opa/v2"
	"github.com/gruntwork-io/terratest/modules/terraform/v2"
)

// TestOPAEvalTerraformModuleWithExtraArgs demonstrates how to pass extra command line arguments to OPA,
// such as --v0-compatible for backwards compatibility with OPA v0.x.
func TestOPAEvalTerraformModuleWithExtraArgs(t *testing.T) {
	t.Parallel()

	tfOpts := &terraform.Options{
		TerraformDir: "../examples/terraform-opa-example/pass",
	}

	opaOpts := &opa.EvalOptions{
		RulePath: "../examples/terraform-opa-example/policy/enforce_source_v0.rego",
		FailMode: opa.FailUndefined,
		// Pass extra command line arguments to OPA eval subcommand
		ExtraArgs: []string{"--v0-compatible"},
	}

	// This will run: opa eval --v0-compatible --fail -i <jsonfile> -d <rulepath> data.enforce_source.allow
	terraform.OPAEval(t, tfOpts, opaOpts, "data.enforce_source.allow")
}
