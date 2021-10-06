---
layout: collection-browser-doc
title: Introduction
category: getting-started
toc: true
excerpt: >-
  Terratest provides a variety of helper functions and patterns for common infrastructure testing tasks. Learn more about Terratest basic usage.
tags: ["basic-usage"]
order: 100
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
- Enforcing policies with OPA
- Testing Helm Charts
- Making HTTP requests
- Running shell commands
- And much more


## Watch: “How to test infrastructure code”

Yevgeniy Brikman talks about how to write automated tests for infrastructure code, including the code written for use with tools such as Terraform, Docker, Packer, and Kubernetes. Topics covered include: unit tests, integration tests, end-to-end tests, dependency injection, test parallelism, retries and error handling, static analysis, property testing and CI / CD for infrastructure code.

This presentation was recorded at QCon San Francisco 2019: https://qconsf.com/.

<iframe width="100%" height="450" allowfullscreen src="https://www.youtube.com/embed/xhHOW0EF5u8"></iframe>

### Slides

Slides to the video can be found here: [Slides: How to test infrastructure code](https://www.slideshare.net/brikis98/how-to-test-infrastructure-code-automated-testing-for-terraform-kubernetes-docker-packer-and-more){:target="\_blank"}.


## Gruntwork

Terratest was developed at [Gruntwork](https://gruntwork.io/) to help maintain the [Infrastructure as Code
Library](https://gruntwork.io/infrastructure-as-code-library/), which contains over 300,000 lines of code written
in Terraform, Go, Python, and Bash, and is used in production by hundreds of companies.

<div class="cb-post-cta">
  <span class="title">See how to get started with Terratest</span>
  <a class="btn btn-primary" href="{{site.baseurl}}/docs/getting-started/quick-start/">Quick Start</a>
</div>
