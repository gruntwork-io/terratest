---
layout: collection-browser-doc
title: Iterating locally using test stages
category: testing-best-practices
excerpt: >-
  Learn more about Terratest's `test_structure`.
tags: ["testing-best-practices", "test_structure"]
order: 210
nav_title: Documentation
nav_title_link: /docs/
---

Most automated tests written with Terratest consist of multiple "stages", such as:

1.  Build an AMI using Packer
1.  Deploy the AMI using Terraform
1.  Validate that the AMI works as expected
1.  Undeploy the AMI using Terraform

Often, while testing locally, you'll want to re-run some subset of these stages over and over again: for example, you
might want to repeatedly run the validation step while you work out the kinks. Having to run _all_ of these stages
each time you change a single line of code can be very slow.

This is where Terratest's `test_structure` package comes in handy: it allows you to explicitly break up your tests into
stages and to be able to disable any one of those stages simply by setting an environment variable. Check out the
[terraform_packer_example_test.go](https://github.com/gruntwork-io/terratest/blob/master/test/terraform_packer_example_test.go) 
for working sample code.
