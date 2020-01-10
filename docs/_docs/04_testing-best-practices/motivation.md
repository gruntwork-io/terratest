---
layout: collection-browser-doc
title: Motivation
category: testing-best-practices
excerpt: >-
  Testing infrastructure as code (IaC) is hard.
tags: ["testing-best-practices"]
order: 400
nav_title: Documentation
nav_title_link: /docs/
---

Testing infrastructure as code (IaC) is hard. With general purpose programming languages (e.g., Java, Python, Ruby),
you have a "localhost" environment where you can run and test the code before you commit. You can also isolate parts
of your code from external dependencies to create fast, reliable unit tests. With IaC, neither of these advantages is
typically available, as there isn't a "localhost" equivalent for most IaC code (e.g., I can't use Terraform to deploy
an AWS VPC on my own laptop) and there's no way to isolate your code from the outside world (i.e., the whole point of
a tool like Terraform is to make calls to AWS, so if you remove AWS, there's nothing left).

That means that most of the tests are going to be integration tests that deploy into a real AWS account. This makes
the tests effective at catching real-world bugs, but it also makes them much slower and more brittle. We'll outline some best practices to minimize the downsides of this sort of testing.

1.  [Unit tests, integration tests, end-to-end tests]({{site.baseurl}}/docs/testing-best-practices/unit-integration-end-to-end-test/)
1.  [Testing environment]({{site.baseurl}}/docs/testing-best-practices/testing-environment/)
1.  [Namespacing]({{site.baseurl}}/docs/testing-best-practices/namespacing/)
1.  [Cleanup]({{site.baseurl}}/docs/testing-best-practices/cleanup/)
1.  [Timeouts and logging]({{site.baseurl}}/docs/testing-best-practices/timeouts-and-logging/)
1.  [Debugging interleaved test output]({{site.baseurl}}/docs/testing-best-practices/debugging-interleaved-test-output/)
1.  [Avoid test caching]({{site.baseurl}}/docs/testing-best-practices/avoid-test-caching/)
1.  [Error handling]({{site.baseurl}}/docs/testing-best-practices/error-handling/)
1.  [Iterating locally using Docker]({{site.baseurl}}/docs/testing-best-practices/iterating-locally-using-docker/)
1.  [Iterating locally using test stages]({{site.baseurl}}/docs/testing-best-practices/iterating-locally-using-test-stages/)
