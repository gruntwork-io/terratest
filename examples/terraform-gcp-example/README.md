# Terraform GCP Example

This folder contains a simple Terraform module that deploys resources in [GCP](https://cloud.google.com/) to demonstrate
how you can use Terratest to write automated tests for your GCP Terraform code. This module deploys a [Compute
Instance](https://cloud.google.com/compute/) and gives that Instance a `Name` with the value specified in the
`instance_name` variable. It also creates a Cloud Storage Bucket using the `bucket_name` and `bucket_location` variables.

Check out [test/terraform_gcp_example_test.go](/test/gcp/terraform_gcp_example_test.go) to see how you can write
automated tests for this module.

Note that the Compute Instance in this module doesn't actually do anything; it just runs a Vanilla Ubuntu 16.04 Image for
demonstration purposes. For slightly more complicated, real-world examples of Terraform modules, see
[terraform-http-example](/examples/terraform-http-example) and [terraform-ssh-example](/examples/terraform-ssh-example).

**WARNING**: This module and the automated tests for it deploy real resources into your GCP account which can cost you
money. The resources are all part of the [GCP Free Tier](https://cloud.google.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all GCP charges.

## Running this module manually

1. Sign up for [GCP](https://cloud.google.com/).
1. Configure your GCP credentials using one of the [supported methods for GCP CLI
   tools](https://cloud.google.com/sdk/docs/quickstarts).
1. Install [Terraform](https://www.terraform.io/) and make sure it's in your `PATH`.
1. Ensure the desired Project ID is set: `export GOOGLE_CLOUD_PROJECT=terratest-ABCXYZ`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. When you're done, run `terraform destroy`.

## Running automated tests against this module

1. Sign up for [GCP](https://cloud.google.com/free/).
1. Configure your GCP credentials using the [GCP CLI
   tools](https://cloud.google.com/sdk/docs/quickstarts).
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. Set `GOOGLE_CLOUD_PROJECT` environment variable to your project name.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformGcpExample`
