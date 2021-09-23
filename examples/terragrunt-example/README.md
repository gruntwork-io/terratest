# Terragrunt Example

This folder contains the simplest possible Terragrunt module—one that just passes inputs to terraform-to demonstrate how
you can use Terratest to write automated tests for your Terragrunt code.

Check out [test/terragrunt_example_test.go](/test/terragrunt_example_test.go) to see how you can
write automated tests for this simple module.

Note that this module doesn't do anything useful; it's just here to demonstrate the simplest usage pattern for
Terratest.




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
