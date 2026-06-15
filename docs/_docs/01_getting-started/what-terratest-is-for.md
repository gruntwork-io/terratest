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

Terratest is a Go library for writing automated tests of infrastructure code. You write standard `go test` files;
Terratest gives you the helpers that make infrastructure-as-code (IaC) tests practical. Knowing what the library is for
helps you decide which helpers to reach for, and explains why some packages are being trimmed over time.

## The five workflows

Terratest covers five workflows that, together, let you test infrastructure end to end:

- **Deploy.** Run `terraform`, `terragrunt`, `packer`, or `docker` from Go and capture their output.
- **Inspect.** Call cloud provider APIs (AWS, Azure, GCP, Kubernetes) to verify resources were created as expected.
- **Interact.** SSH into instances, hit HTTP endpoints, resolve DNS records, query databases: the post-deploy
  connectivity checks that cloud SDKs alone can't do.
- **Validate.** Policy validation against Terraform plans (OPA), test-stage orchestration (`test-structure`), and
  retry-with-backoff for eventual consistency.
- **Tear down.** `terraform destroy` and cleanup helpers.

### Underpinning all five: generic test primitives

Every IaC test reaches for a handful of generic test primitives: unique IDs for resource names, file and fixture
helpers, shell execution, the `TestingT` interface, retry-with-backoff, and a small logging wrapper around
`*testing.T`'s `Logf`. None of these are IaC-specific, but every Terratest test depends on them, so they ship as part of
the library (`random`, `files`, `shell`, `testing`, `retry`, `logger`). They are convenience utilities for people
already writing IaC tests, not a reason to adopt Terratest on their own.

## What Terratest is not

Terratest is deliberately scoped. It is **not**:

- A unit-testing framework. Go's standard `testing` package covers that.
- A mocking library.
- A general-purpose utility collection. Helpers that the standard library already covers (for example, slice and map
  helpers, environment-variable lookups, or thin wrappers around `git`) don't belong here.
- A CI or notification tool.

Everything in Terratest should serve one of the five workflows above or directly support them. Helpers that don't
(standard-library wrappers, standalone CLI tools, and helpers unrelated to IaC testing) are candidates for removal.

## Scope and deprecations

As Terratest moves toward v2, packages that fall outside the scope above are being deprecated and will be removed.
Deprecated packages carry a `// Deprecated:` note in their GoDoc that points at the recommended replacement (usually the
standard library). They keep working for the rest of v1, so you have a full release cycle to migrate. See
[Pinning a Terratest version]({{ site.baseurl }}/docs/getting-started/version-pinning/) if you need to stay on a
specific release while you do.
