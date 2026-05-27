# Terratest

[![Maintained by Gruntwork.io](https://img.shields.io/badge/maintained%20by-gruntwork.io-%235849a6.svg)](https://gruntwork.io/?ref=repo_terratest)
[![Go Report Card](https://goreportcard.com/badge/github.com/gruntwork-io/terratest)](https://goreportcard.com/report/github.com/gruntwork-io/terratest)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/gruntwork-io/terratest?tab=overview)
![go.mod version](https://img.shields.io/github/go-mod/go-version/gruntwork-io/terratest)

Terratest is a Go library that makes it easier to write automated tests for your infrastructure code. It provides a
variety of helper functions and patterns for common infrastructure testing tasks, including:

- Testing Terraform code
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

## Install

**v2 (recommended for new code):**

```bash
go get github.com/gruntwork-io/terratest/modules/<name>/v2@latest
```

Each module is installed independently — pull only what you need. See the
[v2 migration guide](https://terratest.gruntwork.io/docs/migrating-to-v2/overview/)
for the full path mapping and a codemod.

**v1 (still supported in maintenance):**

```bash
go get github.com/gruntwork-io/terratest@latest
```

Requires Go 1.26 or later. To lock to a specific release instead of `@latest`, see [Pinning a Terratest version](https://terratest.gruntwork.io/docs/getting-started/version-pinning/).

## Stability and versioning

Starting with v1.0.0, Terratest follows [semantic versioning](https://semver.org/). Breaking changes to the public API
only happen in major releases.

v2.0.0 split the library into per-domain Go submodules; runtime API is unchanged, only import paths move. v1 is frozen
at its last tag and remains served from `proxy.golang.org` for pinned consumers. Security backports land on the
`v1-maintenance` branch. See the [v2 migration guide](https://terratest.gruntwork.io/docs/migrating-to-v2/overview/).

Migrating from v0.x: see the [v1 migration guide](https://terratest.gruntwork.io/docs/migrating-to-v1/overview/).

## More info

- [Terratest Website](https://terratest.gruntwork.io)
- [Getting started with Terratest](https://terratest.gruntwork.io/docs/getting-started/quick-start/)
- [Terratest Documentation](https://terratest.gruntwork.io/docs/)
- [Contributing to Terratest](https://terratest.gruntwork.io/docs/community/contributing/)
- [Commercial Support](https://gruntwork.io/support/)

## License

This code is released under the Apache 2.0 License. Please see [LICENSE](LICENSE) and [NOTICE](NOTICE) for more details.

Copyright &copy; 2025 Gruntwork, Inc.
