# Terraform OPA Example

This folder contains an [OPA](https://www.openpolicyagent.org/) policy that validates that all module blocks use a
source that comes from the `gruntwork-io` GitHub org (the [enforce_source.rego](./policy/enforce_source.rego) file).
To test this policy, we provided two Terraform modules, [pass](./pass) and [fail](./fail), which will demonstrate how
OPA looks when run against a module that passes the checks, and one that fails the checks.

Check out [test/terraform_opa_example_test.go](/test/terraform_opa_example_test.go) to see how you can write automated
tests for this module.


## Running this module manually

1. Install [OPA](https://www.openpolicyagent.org/) and make sure it's on your `PATH`.
1. Install [hcl2json](https://github.com/tmccombs/hcl2json) and make sure it's on your `PATH`. We need this to convert
   the terraform source code to json as OPA currently doesn't support parsing HCL.
1. Convert each terraform source code in the `pass` or `fail` folder to json by feeding it to `hcl2json`:

       hcl2json pass/main.tf > pass/main.json

1. Run each converted terraform json file against the OPA policy:

       opa eval --fail \
         -i pass/main.json \
         -d policy/enforce_source.rego \
         'data.enforce_source.allow'


## Running automated tests against this module

1. Install [OPA](https://www.openpolicyagent.org/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/).
1. `cd test`
1. `go test -v -run TestOPAEvalTerraformModule`
