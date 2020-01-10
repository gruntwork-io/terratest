---
layout: collection-browser-doc
title: Terraform Basic Example
category: code-examples
excerpt: >-
  A simple "Hello, World" Terraform configuration.
tags: ["terraform", "example"]
image: /assets/img/logos/terraform-logo.png
order: 200
nav_title: Examples
nav_title_link: /examples/
---

This folder contains a very simple Terraform module to demonstrate how you can use Terratest to write automated tests
for your Terraform code. This module takes in an input variable called `example`, renders it using a `template_file`
data source, and outputs the result in an output variable called `example`.

Check out [test/terraform_basic_example_test.go](/test/terraform_basic_example_test.go) to see how you can write
automated tests for this simple module.

Note that this module doesn't do anything useful; it's just here to demonstrate the simplest usage pattern for
Terratest. For a slightly more complicated, real-world example of a Terraform module and the corresponding tests, see
[terraform-aws-example](/examples/terraform-aws-example).




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
1. `go test -v -run TestTerraformBasicExample`
