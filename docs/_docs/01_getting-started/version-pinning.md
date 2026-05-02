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

## Use `go get`, not subpath requires

`go get github.com/gruntwork-io/terratest@<tag>` is enough; the pin applies to every subpackage your tests import. Avoid writing subpath requires by hand (e.g. `go mod edit -require github.com/gruntwork-io/terratest/modules/terraform@<tag>`); Go searches for a `modules/terraform/<tag>` tag and fails with `unknown revision`.
