# Terraform Lambda Example

This folder contains a Terraform module to demonstrate how you can use Terratest to deploy a lambda function
for your Terraform code. This module takes in an input variable called `function_name`, and uses the function name as
an identifier for the lambda and associated resources (e.g. IAM role).

Check out [test/terraform_aws_lambda_example_test.go](/test/terraform_aws_lambda_example_test.go) to see how you can write
automated tests for this simple module.

The function that this module creates is a simple one whose input can cause it to error or echo messages it receives.

## Running this module manually

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformAwsLambdaExample`
