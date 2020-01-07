---
layout: collection-browser-doc
title: Why Terratest?
category: alternative-testing-tools
excerpt: >-
  Are you wondering why you should use Terratest?s
tags: ["tools"]
order: 500
nav_title: Documentation
nav_title_link: /docs/
---

Our experience with building the [Infrastructure as Code Library](https://gruntwork.io/infrastructure-as-code-library/)
is that the _only_ way to create reliable, maintainable infrastructure code is to have a thorough suite of real-world,
end-to-end acceptance tests. Without these sorts of tests, you simply cannot be confident that the infrastructure code
actually works.

This is especially important with modern DevOps, as all the tools are changing so quickly. Terratest has helped us
catch bugs not only in our own code, but also in AWS, Azure, Terraform, Packer, Kafka, Elasticsearch, CircleCI, and
so on. Moreover, by running tests nightly, we're able to catch backwards incompatible changes and
regressions in our dependencies (e.g., backwards incompatibilities in new versions of Terraform) as early as possible.
