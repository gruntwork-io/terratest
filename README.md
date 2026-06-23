# Terratest

[![Maintained by Gruntwork.io](https://img.shields.io/badge/maintained%20by-gruntwork.io-%235849a6.svg)](https://gruntwork.io/?ref=repo_terratest)
[![Go Report Card](https://goreportcard.com/badge/github.com/gruntwork-io/terratest)](https://goreportcard.com/report/github.com/gruntwork-io/terratest)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/gruntwork-io/terratest?tab=overview)
![go.mod version](https://img.shields.io/github/go-mod/go-version/gruntwork-io/terratest)

Terratest is a Go library that makes it easier to write automated tests for your infrastructure code. It provides a
variety of helper functions and patterns for common infrastructure testing tasks, including:

- Testing OpenTofu and Terraform code
- Testing Packer templates
- Testing Docker images
- Executing commands on servers over SSH
- Working with AWS APIs
- Working with Azure APIs
- Working with GCP APIs
- Working with Kubernetes APIs
- Testing Helm Charts
- Making HTTP requests
- Running shell commands
- And much more

## What Terratest is for

Terratest is a Go library for writing automated tests of infrastructure code. It covers five workflows that, together,
let you test infrastructure end to end:

- **Deploy** OpenTofu, Terragrunt, Packer, or Docker from Go and capture their output.
- **Inspect** what got deployed by calling cloud provider APIs (AWS, Azure, GCP, Kubernetes).
- **Interact** with it over the network: SSH, HTTP, DNS, and database checks that cloud SDKs alone can't do.
- **Validate** behavior and policy: OPA against OpenTofu plans, test-stage orchestration, retry-with-backoff for
  eventual consistency.
- **Tear down** with `tofu destroy` and cleanup helpers.

Terratest is deliberately scoped. It is not a unit-testing framework (Go's standard `testing` covers that), a mocking
library, a general-purpose utility collection, or a CI/notification tool. Helpers that fall outside the five workflows
above, including ones the standard library already covers, are being deprecated and removed in v2. See
[What Terratest is for](https://terratest.gruntwork.io/docs/getting-started/what-terratest-is-for/) for the full picture.

## Install

```bash
go get github.com/gruntwork-io/terratest@latest
```

Requires Go 1.26 or later. To lock to a specific release instead of `@latest`, see [Pinning a Terratest version](https://terratest.gruntwork.io/docs/getting-started/version-pinning/).

## Stability and versioning

Starting with v1.0.0, Terratest follows [semantic versioning](https://semver.org/). Breaking changes to the public API
only happen in major releases (e.g. v2.0.0).

Symbols renamed or replaced in v1 are kept with `// Deprecated:` annotations pointing at the new name; removals happen
in v2. Migrating from v0.x: see the [v1 migration guide](https://terratest.gruntwork.io/docs/migrating-to-v1/overview/).

## More info

- [Terratest Website](https://terratest.gruntwork.io)
- [Getting started with Terratest](https://terratest.gruntwork.io/docs/getting-started/quick-start/)
- [Terratest Documentation](https://terratest.gruntwork.io/docs/)
- [Contributing to Terratest](https://terratest.gruntwork.io/docs/community/contributing/)
- [Commercial Support](https://gruntwork.io/support/)

## License

This code is released under the Apache 2.0 License. Please see [LICENSE](LICENSE) and [NOTICE](NOTICE) for more details.

Copyright &copy; 2025 Gruntwork, Inc.
