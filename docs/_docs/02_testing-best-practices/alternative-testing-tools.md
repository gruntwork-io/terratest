---
layout: collection-browser-doc
title: Alternative testing tools
category: testing-best-practices
excerpt: >-
  Learn more about alternatives to Terratest, and how Terratest compares to other testing tools.
tags: ["testing-best-practices", "alternatives"]
order: 211
nav_title: Documentation
nav_title_link: /docs/
---

## A list of infrastructure testing tools

Below is a list of other infrastructure testing tools you may wish to use in addition to Terratest. Check out [How
Terratest compares to other testing tools]({{site.baseurl}}/docs/testing-best-practices/alternative-testing-tools/#how-terratest-compares-to-other-testing-tools) to understand the trade-offs.

1.  [kitchen-terraform](https://github.com/newcontext-oss/kitchen-terraform)
1.  [rspec-terraform](https://github.com/bsnape/rspec-terraform)
1.  [serverspec](https://serverspec.org/)
1.  [inspec](https://www.inspec.io/)
1.  [Goss](https://github.com/aelsabbahy/goss)
1.  [awspec](https://github.com/k1LoW/awspec)
1.  [Terraform's acceptance testing framework](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#acceptance-tests-testing-interactions-with-external-services)
1.  [ruby_terraform](https://github.com/infrablocks/ruby_terraform)



## Why Terratest?

Our experience with building the [Infrastructure as Code Library](https://gruntwork.io/infrastructure-as-code-library/)
is that the _only_ way to create reliable, maintainable infrastructure code is to have a thorough suite of real-world,
end-to-end acceptance tests. Without these sorts of tests, you simply cannot be confident that the infrastructure code
actually works.

This is especially important with modern DevOps, as all the tools are changing so quickly. Terratest has helped us
catch bugs not only in our own code, but also in AWS, Azure, Terraform, Packer, Kafka, Elasticsearch, CircleCI, and
so on. Moreover, by running tests nightly, we're able to catch backwards incompatible changes and
regressions in our dependencies (e.g., backwards incompatibilities in new versions of Terraform) as early as possible.



## How Terratest compares to other testing tools

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
