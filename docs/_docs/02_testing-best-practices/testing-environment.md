---
layout: collection-browser-doc
title: Testing environment
category: testing-best-practices
excerpt: >-
  Learn more about testing environments.
tags: ["testing-best-practices"]
order: 202
nav_title: Documentation
nav_title_link: /docs/
---

Since most automated tests written with Terratest can make potentially destructive changes in your environment, we
strongly recommend running tests in an environment that is totally separate from production. For example, if you are
testing infrastructure code for AWS, you should run your tests in a completely separate AWS account.

This means that you will have to write your infrastructure code in such a way that you can plug in ([dependency
injection](https://en.wikipedia.org/wiki/Dependency_injection)) environment-specific details, such as account IDs,
domain names, IP addresses, etc. Adding support for this will typically make your code cleaner and more flexible.
