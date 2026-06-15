---
layout: collection-browser-doc
title: What Terratest is for
category: getting-started
toc: true
excerpt: >-
  Terratest is a Go library for writing automated tests of infrastructure code. Learn the five workflows it covers, and what it is deliberately not.
tags: ["basic-usage"]
order: 101
nav_title: Documentation
nav_title_link: /docs/
---

Terratest is a Go library for writing automated tests of infrastructure code. You write ordinary `_test.go` files, and
Terratest gives you the helpers that make infrastructure-as-code (IaC) testing practical. Once you know what the library
is for, it is easier to decide which helpers to reach for, and to understand why some packages are being trimmed over
time.

## The five workflows

Terratest covers five workflows. Together, they let you test infrastructure from end to end.

- **Deploy.** Run `tofu`, `terragrunt`, `packer`, or `docker` from Go and capture their output.
- **Inspect.** Call cloud provider APIs (AWS, Azure, GCP, and Kubernetes) to verify that resources were created as
  expected.
- **Interact.** Run the post-deploy connectivity checks that cloud SDKs alone can't do: SSH into instances, hit HTTP
  endpoints, resolve DNS records, and query databases.
- **Validate.** Check OpenTofu plans against policy with OPA, orchestrate test stages with `test-structure`, and retry
  with backoff to handle eventual consistency.
- **Tear down.** Run `tofu destroy` and clean up with related helpers.

## The primitives underneath

Every IaC test leans on a handful of generic building blocks: unique IDs for resource names, file and fixture helpers,
shell execution, the `TestingT` interface, retry with backoff, and a small logging wrapper around the `Logf` method on
`*testing.T`. None of these are specific to infrastructure, but every Terratest test depends on them, so they ship with
the library in the `random`, `files`, `shell`, `testing`, `retry`, and `logger` packages.

## What Terratest is not

Terratest is deliberately narrow in scope. It is not:

- A unit-testing framework. Go's standard `testing` package already covers that.
- A mocking library.
- A general-purpose utility collection. Anything the standard library already handles (slice and map helpers,
  environment-variable lookups, thin wrappers around `git`, and the like) does not belong here.
- A CI or notification tool.

The rule of thumb: everything in Terratest should serve one of the five workflows above, or directly support them.
Helpers that don't (standard-library wrappers, standalone CLI tools, and anything unrelated to IaC testing) are
candidates for removal.

## Scope and deprecations

As Terratest moves toward v2, packages that fall outside this scope are being deprecated and will eventually be removed.
A deprecated package carries a `// Deprecated:` note in its GoDoc that points to the recommended replacement, which is
usually the standard library. These packages keep working for the rest of v1, so you get a full release cycle to
migrate. If you need to stay on a specific release while you do, see
[Pinning a Terratest version]({{ site.baseurl }}/docs/getting-started/version-pinning/).
