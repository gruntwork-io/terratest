---
layout: collection-browser-doc
title: Terraform HTTP Example
category: code-examples
excerpt: >-
  A more complicated Terraform configuration that deploys a simple web server that responds to HTTP requests in AWS.
tags: ["terraform", "example"]
image: /assets/img/logos/terraform-logo.png
order: 201
nav_title: Examples
nav_title_link: /examples/
---

This folder contains a simple Terraform module that deploys resources in [AWS](https://aws.amazon.com/) to demonstrate
how you can use Terratest to write automated tests for your AWS Terraform code. This module deploys an [EC2
Instance](https://aws.amazon.com/ec2/) in the AWS region specified in the `aws_region` variable. The EC2 Instance runs
a simple web server that listens for HTTP requests on the port specified by the `instance_port` variable and returns
the text specified by the `instance_text` variable.

Check out [test/terraform_http_example_test.go](/test/terraform_http_example_test.go) to see how you can write
automated tests for this module.

Note that the example in this module is still fairly simplified, as the "web server" we run just servers up a static
`index.html`, and not in a particularly production-ready manner! For a more complicated, real-world, end-to-end
example of a Terraform module and web server, see [terraform-packer-example](/examples/terraform-packer-example).

**WARNING**: This module and the automated tests for it deploy real resources into your AWS account which can cost you
money. The resources are all part of the [AWS Free Tier](https://aws.amazon.com/free/), so if you haven't used that up,
it should be free, but you are completely responsible for all AWS charges.





## Running this module manually

1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. If you're using the `~/.aws/config` file for profiles then export `AWS_SDK_LOAD_CONFIG` as "True".
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Run `terraform init`.
1. Run `terraform apply`.
1. The `instance_url` output variable shows you the URL of the web server. Try opening it in your browser!
1. When you're done, run `terraform destroy`.




## Running automated tests against this module

1. Sign up for [AWS](https://aws.amazon.com/).
1. Configure your AWS credentials using one of the [supported methods for AWS CLI
   tools](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html), such as setting the
   `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables. If you're using the `~/.aws/config` file for profiles then export `AWS_SDK_LOAD_CONFIG` as "True".
1. Install [Terraform](https://www.terraform.io/) and make sure it's on your `PATH`.
1. Install [Golang](https://golang.org/) and make sure this code is checked out into your `GOPATH`.
1. `cd test`
1. `dep ensure`
1. `go test -v -run TestTerraformHttpExample`
