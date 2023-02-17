# Packer Basic Example

This folder contains a very simple Packer template to demonstrate how you can use Terratest to write automated tests
for your Packer templates. The template just creates an up-to-date Ubuntu AMI by running `apt-get update` and
`apt-get upgrade`.

Check out [test/packer_basic_example_test.go](/test/packer_basic_example_test.go) to see how you can write
automated tests for this simple template.

Note that this template doesn't do anything useful; it's just here to demonstrate the simplest usage pattern for
Terratest. For slightly more complicated, real-world examples of Packer templates and the corresponding tests, see
[packer-docker-example](/examples/packer-docker-example) and
[terraform-packer-example](/examples/terraform-packer-example).



## Installation steps
1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.

## Building the Packer template manually (Packer >= 1.7.0)
1. Run `packer init build.pkr.hcl`.  # Use build-gcp.pkr.hcl if using GCP
1. Run `packer build build.pkr.hcl`. # Use build-gcp.pkr.hcl if using GCP

## Building the Packer template manually (Packer < 1.7.0)
1. Run `packer build build.json`.



## Running automated tests against this Packer template

1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestPackerBasicExample`




## Running automated tests against this Packer template for the GCP builder

1. Sign up for [GCP](https://cloud.google.com/).
1. Configure your GCP credentials using one of the
   [Authentication](https://www.packer.io/docs/builders/googlecompute.html#authentication) methods.
1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestPackerGCPBasicExample`




## Running automated tests against this Packer template for the OCI builder

1. Sign up for [OCI](https://cloud.oracle.com/cloud-infrastructure).
1. Configure your OCI credentials via [CLI Configuration
   Information](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm).
1. Create [VCN](https://docs.cloud.oracle.com/iaas/Content/GSG/Tasks/creatingnetwork.htm) and subnet
   resources in your tenancy (a.k.a. a root compartment).
1. (Optional) Create `TF_VAR_pass_phrase` environment property with the pass phrase for decrypting of the OCI [API signing
      key](https://docs.cloud.oracle.com/iaas/Content/API/Concepts/apisigningkey.htm) (can be omitted
      if the key is not protected).
1. Install [Packer](https://www.packer.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestPackerOciExample`
