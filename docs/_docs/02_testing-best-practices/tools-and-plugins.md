---
layout: collection-browser-doc
title: Tools and Plugins
category: testing-best-practices
excerpt: >-
  Additional tools and plugins for terratest to make integration with existing workflows easier.
tags: ["testing-best-practices", "alternatives","plugins","tooling"]
order: 214
nav_title: Documentation
nav_title_link: /docs/
---

## Tools and Plugins

This page contains a list of tools and plugins for Terratest to make integration with existing workflows easier. If you've created other tools that integrate with Terratest, a PR to add it to this list is very welcome!

- [Tools and Plugins](#tools-and-plugins)
  - [Terratest Maven Plugin](#terratest-maven-plugin)
  - [Custom AWS Endpoints](#custom-aws-endpoints)
  - [S3 Path Style Endpoints](#s3-path-style-endpoints)

### Terratest Maven Plugin

The Terratest Maven Plugin aims to bring Terratest to the JVM world. Create your Go based tests beside your Java code with Maven and run them together. You can export the results into Json or an HTML page. As the plugin is MIT licensed, it is easy and painless to integrate into any Java+Maven combination. To learn more check out the website: [Terratest Maven Plugin](https://terratest-maven-plugin.github.io) and the [GitHub repository](https://github.com/terratest-maven-plugin/terratest-maven-plugin)

### Custom AWS Endpoints

Terratest can be used with a custom AWS endpoint. One of the use cases for this is testing your automation using a mock AWS server such as [Moto in standalone server mode](https://docs.getmoto.org/en/latest/docs/getting_started.html#stand-alone-server-mode). This allows you to test your AWS based automation without an AWS account.

To configure this, set `TERRATEST_CUSTOM_AWS_ENDPOINT` to the endpoint that you want Terratest to use for accessing AWS resources. You must also configure your automation to use the same endpoint, for example [Terraform supports custom AWS endpoints as well](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/custom-service-endpoints).

### S3 Path Style Endpoints

By default, the golang client will interact with s3 in a way that uses subdomains to specify information about the bucket. If you have issues with subdomains resolving properly, set `TERRATEST_AWS_S3_USE_PATH_STYLE_ENDPOINT` to `1` to tell Terratest to use [path style access](https://docs.aws.amazon.com/AmazonS3/latest/userguide/VirtualHosting.html#path-style-access) rather than subdomains for this information.

Note that AWS is [deprecating path style endpoints](https://aws.amazon.com/blogs/aws/amazon-s3-path-deprecation-plan-the-rest-of-the-story/) someday, but it seems like enough users are pushing against it that it's not scheduled for a specific time.
