---
layout: collection-browser-doc
title: Pinning a Terratest version
category: getting-started
excerpt: Lock your tests to a known-good Terratest release and upgrade safely.
tags: ["versioning", "go-modules", "pinning"]
order: 106
nav_title: Documentation
nav_title_link: /docs/
---

## Why pin a version?

Terratest is a Go library, so the version your tests build against is whatever your test module's `go.mod` resolves to. Without an explicit pin, `go get` will pick up the newest published release the next time the module graph is recomputed, which means a regression in a future release can fail builds you did not change. Pinning makes the version a deliberate choice and keeps CI reproducible.

Starting with v1.0.0, Terratest follows [semantic versioning](https://semver.org/) and breaking changes only happen in major releases. Pinning is still useful when you want to control _when_ you take a minor or patch bump.

## Pin to a specific release

From the directory that contains your test module's `go.mod`, run:

```bash
go get github.com/gruntwork-io/terratest@v0.56.0
go mod tidy
```

Replace `v0.56.0` with the release tag you want. See the [Releases page](https://github.com/gruntwork-io/terratest/releases) for the full list. Always commit both `go.mod` and `go.sum` so the pin and its checksums travel with your code. Anyone who runs `go test` against the same `go.mod`/`go.sum` will resolve to the exact same Terratest version.

To downgrade or revert to an earlier release, run the same command with the older tag.

## Use the root module path

Use the root module path `github.com/gruntwork-io/terratest`, not submodule paths like `github.com/gruntwork-io/terratest/modules/terraform`. Terratest publishes a single Go module at the repository root; submodule-style paths such as `modules/terraform/v0.51.0` are not valid Go module versions and will fail with `unknown revision`:

```text
reading github.com/gruntwork-io/terratest/modules/terraform/go.mod
at revision modules/terraform/v0.51.0: unknown revision
```

If you import individual subpackages in your tests (for example `github.com/gruntwork-io/terratest/modules/terraform`), that is fine. They all resolve through the single root module pin.

## Upgrading

When you are ready to take a newer release, bump the pin explicitly and run the test suite:

```bash
go get github.com/gruntwork-io/terratest@v0.56.0
go mod tidy
go test ./...
```

If a release introduces a regression, revert the pin to the last known-good version and open an issue at [github.com/gruntwork-io/terratest/issues](https://github.com/gruntwork-io/terratest/issues). For history of breaking changes between v0.x and v1.0.0, see [`MIGRATION.md`](https://github.com/gruntwork-io/terratest/blob/main/MIGRATION.md).
