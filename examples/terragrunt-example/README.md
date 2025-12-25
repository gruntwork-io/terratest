# Terragrunt Example

This folder contains a single Terragrunt module demonstrating how to test it using Terratest's `terraform` package
with `TerraformBinary: "terragrunt"`.

Check out [modules/terragrunt/terragrunt_example_test.go](/modules/terragrunt/terragrunt_example_test.go) to see how you can write automated tests
for this module.

For testing multiple Terragrunt modules with dependencies (using `--all` commands), see
[terragrunt-multi-module-example](/examples/terragrunt-multi-module-example).




## Running this module manually

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Terragrunt](https://terragrunt.gruntwork.io/) and make sure it's on your `PATH`.
1. Run `terragrunt apply`.
1. When you're done, run `terragrunt destroy`.




## Running automated tests against this module

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Terragrunt](https://terragrunt.gruntwork.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `go test -v -run TestTerragruntExample`
