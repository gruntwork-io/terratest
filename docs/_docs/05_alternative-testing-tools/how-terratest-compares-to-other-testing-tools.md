---
layout: collection-browser-doc
title: How Terratest compares to other testing tools
category: alternative-testing-tools
excerpt: >-
  Comparing to other testing tools, Terratest doesn't check individual properties only. The question we're trying to answer is, "does the
  infrastructure actually work?"
tags: ["tools"]
order: 501
nav_title: Documentation
nav_title_link: /docs/
---

Most of the other infrastructure testing tools we've seen are focused on making it easy to check the properties of a
single server or resource. For example, the various `xxx-spec` tools offer a nice, concise language for connecting to
a server and checking if, say, `httpd` is installed and running. These tools are effectively verifying that individual
"properties" of your infrastructure meet a certain spec.

Terratest approaches the testing problem from a different angle. The question we're trying to answer is, "does the
infrastructure actually work?" Instead of checking individual server properties (e.g., is `httpd` installed and
running), we'll actually make HTTP requests to the server and check that we get the expected response; or we'll store
data in a database and make sure we can read it back out; or we'll try to deploy a new version of a Docker container
and make sure the orchestration tool can roll out the new container with no downtime.

Moreover, we use Terratest not only with individual servers, but to test entire systems. For example, the automated
tests for the [Vault module](https://github.com/hashicorp/terraform-aws-vault/tree/master/modules) do the following:

1.  Use Packer to build an AMI.
1.  Use Terraform to create self-signed TLS certificates.
1.  Use Terraform to deploy all the infrastructure: a Vault cluster (which runs the AMI from the previous step), Consul
    cluster, load balancers, security groups, S3 buckets, and so on.
1.  SSH to a Vault node to initialize the cluster.
1.  SSH to all the Vault nodes to unseal them.
1.  Use the Vault SDK to store data in Vault.
1.  Use the Vault SDK to make sure you can read the same data back out of Vault.
1.  Use Terraform to undeploy and clean up all the infrastructure.

The steps above are exactly what you would've done to test the Vault module manually. Terratest helps automate this
process. You can think of Terratest as a way to do end-to-end, acceptance or integration testing, whereas most other
tools are focused on unit or functional testing.
