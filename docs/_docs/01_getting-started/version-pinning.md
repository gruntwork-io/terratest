---
layout: collection-browser-doc
title: Pinning a Terratest version
category: getting-started
excerpt: Lock your tests to a specific Terratest release.
tags: ["versioning", "go-modules", "pinning"]
order: 106
nav_title: Documentation
nav_title_link: /docs/
---

Pin Terratest if you need reproducible test builds and deliberate control over when you adopt a new release.

## Pin to a release

From the directory that contains your test module's `go.mod`:

```bash
go get github.com/gruntwork-io/terratest@v0.56.0
go mod tidy
```

Replace `v0.56.0` with the tag you want from the [Releases page](https://github.com/gruntwork-io/terratest/releases). The same command upgrades, downgrades, or reverts the pin. Always commit `go.mod` and `go.sum` so the pin travels with your code; avoid `@latest` if you need reproducibility.

## Use the root module path

Pin `github.com/gruntwork-io/terratest`, not submodule paths like `github.com/gruntwork-io/terratest/modules/terraform`. Terratest publishes a single Go module at the repository root, so submodule-style versions (e.g. `modules/terraform/v0.51.0`) fail with `unknown revision`. Subpackage imports in your test code are fine; they all resolve through the root module pin.
