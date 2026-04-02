# Terraform GCP Pub/Sub Example

This folder contains a simple Terraform module that deploys resources in [GCP](https://cloud.google.com/) to demonstrate
how you can use Terratest to write automated tests for your GCP Terraform code. This module deploys a [Pub/Sub Topic](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.topics) and a [Pub/Sub Subscription](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions) attached to that topic.

Check out [test/gcp/terraform_gcp_pubsub_example_test.go](../../test/gcp/terraform_gcp_pubsub_example_test.go) to see how 
you can write automated tests for this module.

**WARNING**: This module and the automated tests for it deploy real resources into your GCP account which can cost you
money. The resources are typically part of the [GCP Free Tier](https://cloud.google.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all GCP charges.

## Running this module manually

1. Sign up for [GCP](https://cloud.google.com/).
1. Configure your GCP credentials using one of the [supported methods for GCP CLI
   tools](https://cloud.google.com/sdk/docs/quickstarts).
1. Install [Terraform](https://www.terraform.io/) and make sure it's in your `PATH`.
1. Ensure the desired Project ID is set: `export GOOGLE_CLOUD_PROJECT=your-project-id`.
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
1. `cd test/gcp`
1. `go test -v -tags=gcp -run TestTerraformGcpPubSubExample`
