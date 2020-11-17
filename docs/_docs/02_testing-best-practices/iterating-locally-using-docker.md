---
layout: collection-browser-doc
title: Iterating locally using Docker
category: testing-best-practices
excerpt: >-
  If you're writing scripts (i.e., Bash, Python, or Go), you should be able to test them locally using Docker. Docker containers typically build 10x faster and start 100x faster than real servers.
tags: ["testing-best-practices", "docker"]
order: 209
nav_title: Documentation
nav_title_link: /docs/
---

For most infrastructure code, your only option is to deploy into a real environment such as AWS. However, if you're
writing scripts (i.e., Bash, Python, or Go), you should be able to test them locally using Docker. Docker containers
typically build 10x faster and start 100x faster than real servers, so using Docker for testing can help you iterate
much faster.

Here are some techniques we use with Docker:

- If your script is used in a Packer template, add a [Docker
  builder](https://www.packer.io/docs/builders/docker.html) to the template so you can create a Docker image from the
  same code. See the [Packer Docker Example](https://github.com/gruntwork-io/terratest/tree/master/examples/packer-docker-example) for working sample code.

- We have prebuilt Docker images for major Linux distros that have many important dependencies (e.g., curl, vim,
  tar, sudo) already installed. See the [test-docker-images folder](https://github.com/gruntwork-io/terratest/tree/master/test-docker-images) for more details.

- Create a `docker-compose.yml` to make it easier to run your Docker image with all the ports, environment variables,
  and other settings it needs. See the [Packer Docker Example](https://github.com/gruntwork-io/terratest/tree/master/examples/packer-docker-example) for working sample code.

- With scripts in Docker, you can replace _some_ real-world dependencies with mocks! One way to do this is to create
  some "mock scripts" and to bind-mount them in `docker-compose.yml` in a way that replaces the real dependency. For
  example, if your script calls the `aws` CLI, you could create a mock script called `aws` that shows up earlier in the
  `PATH`. Using mocks allows you to test 100% locally, without external dependencies such as AWS.
