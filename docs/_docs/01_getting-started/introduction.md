---
layout: collection-browser-doc
title: Introduction
category: getting-started
excerpt: >-
  Terratest provides a variety of helper functions and patterns for common infrastructure testing tasks. Learn more about Terratest basic usage.
excerpt_html: >-
  <p>Empty</p>
excerpt_md: >-
  The basic usage pattern for writing automated tests with Terratest is to:

  1.  Write tests using Go's built-in package testing: you create a file ending in
      `_test.go` and run tests with the `go test` command.
  1.  Use Terratest to execute your _real_ IaC tools (e.g., Terraform, Packer, etc.) to deploy _real_ infrastructure
      (e.g., servers) in a _real_ environment (e.g., AWS).
  1.  Validate that the infrastructure works correctly in that environment by making HTTP requests, API calls, SSH
      connections, etc.
  1.  Undeploy everything at the end of the test.

    {:.bottom-wide-block}
    ```go
    terraformOptions := &terraform.Options {
      // The path to where your Terraform code is located
      TerraformDir: "../examples/terraform-basic-example",
    }

    // At the end of the test, run `terraform destroy` to clean up any resources that were created
    defer terraform.Destroy(t, terraformOptions)

    // This will run `terraform init` and `terraform apply` and fail the test if there are any errors
    terraform.InitAndApply(t, terraformOptions)

    // Validate your code works as expected
    validateServerIsWor king(t, terraformOptions)
    ```
tags: ["basic-usage"]
order: 100
index_list:
  no_hover_enlarge_effect: true
  disable_card_link: true
  read_more_btn: true
nav_title: Documentation
nav_title_link: /docs/
---

## Introduction

Terratest is a Go library that makes it easier to write automated tests for your infrastructure code. It provides a
variety of helper functions and patterns for common infrastructure testing tasks, including:

- Testing Terraform code
- Testing Packer templates
- Testing Docker images
- Executing commands on servers over SSH
- Working with AWS APIs
- Working with Azure APIs
- Working with GCP APIs
- Working with Kubernetes APIs
- Testing Helm Charts
- Making HTTP requests
- Running shell commands
- And much more

For an introduction to Terratest, including unit tests, integration tests, end-to-end tests, dependency injection, test
parallelism, retries, error handling, and static analysis, see the talk [Automated Testing for Terraform, Docker,
Packer, Kubernetes, and More](https://www.infoq.com/presentations/automated-testing-terraform-docker-packer/) and the
blog post [Open sourcing Terratest: a swiss army knife for testing infrastructure
code](https://blog.gruntwork.io/open-sourcing-terratest-a-swiss-army-knife-for-testing-infrastructure-code-5d883336fcd5).


## Basic usage

The basic usage pattern for writing automated tests with Terratest is to:

1.  Write tests using Go's built-in [package testing](https://golang.org/pkg/testing/): you create a file ending in
    `_test.go` and run tests with the `go test` command.
1.  Use Terratest to execute your _real_ IaC tools (e.g., Terraform, Packer, etc.) to deploy _real_ infrastructure
    (e.g., servers) in a _real_ environment (e.g., AWS).
1.  Validate that the infrastructure works correctly in that environment by making HTTP requests, API calls, SSH
    connections, etc.
1.  Undeploy everything at the end of the test.

Here's a simple example of how to test some Terraform code:

```go
terraformOptions := &terraform.Options {
  // The path to where your Terraform code is located
  TerraformDir: "../examples/terraform-basic-example",
}

// At the end of the test, run `terraform destroy` to clean up any resources that were created
defer terraform.Destroy(t, terraformOptions)

// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
terraform.InitAndApply(t, terraformOptions)

// Validate your code works as expected
validateServerIsWorking(t, terraformOptions)
```

<div class="cb-post-cta">
  <span class="title">See how to start with Terratest</span>
  <a class="btn btn-primary" href="{{site.baseurl}}/docs/getting-started/quick-start/">Quick Start</a>
</div>

## Gruntwork

Terratest was developed at [Gruntwork](https://gruntwork.io/) to help maintain the [Infrastructure as Code
Library](https://gruntwork.io/infrastructure-as-code-library/), which contains over 300,000 lines of code written
in Terraform, Go, Python, and Bash, and is used in production by hundreds of companies.
