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

	tfOpts := &terraform.Options{
		// website::tag::1:: Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/terraform-opa-example/pass",
	}

	opaOpts := &opa.EvalOptions{
		// website::tag::2:: Set the path to the OPA policy code that should be used.
		RulePath: "../examples/terraform-opa-example/policy/enforce_source.rego",

		// website::tag::3:: Run OPA in fail mode so that it will exit with non-zero exit code when the result query is undefined.
		FailMode: opa.FailUndefined,
	}

	// website::tag::4:: Run OPA with the configured options, querying for the allow variable. The OPAEval function automatically expects the check to pass, failing the test if opa eval exits with non-zero exit code.
	terraform.OPAEval(t, tfOpts, opaOpts, "data.enforce_source.allow")
}

// An example of how to use Terratest to run OPA policy checks on Terraform source code. This will check the module
// called `fail` against the rego policy `enforce_source` defined in the `terraform-opa-example` folder and validate
// that the module fails the OPA checks.
func TestOPAEvalTerraformModuleFailsCheck(t *testing.T) {
	t.Parallel()

	// website::tag::5:: Configure in a similar fashion to the above test, but run against the `fail` example.
	policyPath := "../examples/terraform-opa-example/policy/enforce_source.rego"
	tfOpts := &terraform.Options{TerraformDir: "../examples/terraform-opa-example/fail"}
	opaOpts := &opa.EvalOptions{
		FailMode: opa.FailUndefined,
		RulePath: policyPath,
	}

	// website::tag::6:: Here we expect the checks to fail, so we use `OPAEvalE` to check the error. Note that on the files that failed, this function will rerun `opa eval` with the query set to `data`, so you can see the values of all the variables in the policy. This is useful for debugging failures.
	require.Error(t, terraform.OPAEvalE(t, tfOpts, opaOpts, "data.enforce_source.allow"))
}

// An example of how to use Terratest to run OPA policy checks on Terraform source code using a remote OPA policy source
// file. This will check the module called `pass` against the rego policy `enforce_source` defined in the
// `terraform-opa-example` folder of the terratest repository.
func TestOPAEvalTerraformModuleRemotePolicy(t *testing.T) {
	t.Parallel()

	tfOpts := &terraform.Options{
		TerraformDir: "../examples/terraform-opa-example/pass",
	}
	opaOpts := &opa.EvalOptions{
		RulePath: "git::https://github.com/gruntwork-io/terratest.git//examples/terraform-opa-example/policy/enforce_source.rego?ref=master",
		FailMode: opa.FailUndefined,
	}
	terraform.OPAEval(t, tfOpts, opaOpts, "data.enforce_source.allow")
}
