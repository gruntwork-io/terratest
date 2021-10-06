package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/opa"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// An example of how to use Terratest to run OPA policy checks on Terraform source code. This will check the module
// called `pass` against the rego policy `enforce_source` defined in the `terraform-opa-example` folder.
func TestOPAEvalTerraformModulePassesCheck(t *testing.T) {
	t.Parallel()

	policyPath := "../examples/terraform-opa-example/policy/enforce_source.rego"
	tfOpts := &terraform.Options{TerraformDir: "../examples/terraform-opa-example/pass"}
	opaOpts := &opa.EvalOptions{
		FailMode: opa.FailUndefined,
		RulePath: policyPath,
	}
	terraform.OPAEval(t, tfOpts, opaOpts, "data.enforce_source.allow")
}

// An example of how to use Terratest to run OPA policy checks on Terraform source code. This will check the module
// called `fail` against the rego policy `enforce_source` defined in the `terraform-opa-example` folder and validate
// that the module fails the OPA checks.
func TestOPAEvalTerraformModuleFailsCheck(t *testing.T) {
	t.Parallel()

	policyPath := "../examples/terraform-opa-example/policy/enforce_source.rego"
	tfOpts := &terraform.Options{TerraformDir: "../examples/terraform-opa-example/fail"}
	opaOpts := &opa.EvalOptions{
		FailMode: opa.FailUndefined,
		RulePath: policyPath,
	}
	require.Error(t, terraform.OPAEvalE(t, tfOpts, opaOpts, "data.enforce_source.allow"))
}
