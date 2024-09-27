# Terraform AWS Endpoint Example

This folder contains a simple Terraform module to demonstrate using custom
endpoints. It's deploying some AWS resources to `http://localhost:5000`, which
is the default port for [moto running in server
mode](http://docs.getmoto.org/en/latest/docs/server_mode.html). This allows for
testing terraform modules locally with no connection to AWS.

Check out
[test/terraform_aws_endpoint_example_test.go](/test/terraform_aws_endpoint_example_test.go)
to see how you can write automated tests for this module.

## Running this module manually

1. Run [Moto locally in server mode](http://docs.getmoto.org/en/latest/docs/server_mode.html)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Run [Moto locally in server mode](http://docs.getmoto.org/en/latest/docs/server_mode.html)
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformAwsEndpointExample`
