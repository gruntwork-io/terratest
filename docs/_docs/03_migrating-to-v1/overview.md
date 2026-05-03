---
layout: collection-browser-doc
title: Overview
category: migrating-to-v1
excerpt: >-
  Summary of breaking changes in Terratest v1.0.0 and where to look first.
tags: ["migration", "v1"]
order: 300
nav_title: Documentation
nav_title_link: /docs/
---

Terratest v1.0.0 is the first stable release. Once you are on v1, breaking
changes to the public API only happen in major releases (e.g. v2.0.0), per
[semver](https://semver.org/). Renamed or replaced symbols stay around as
deprecated aliases inside v1; full removal is deferred to v2.

This page is the orientation map for v0.x to v1.0.0. It tells you what
shape the changes take and where to look; the per-service guides hold the
mechanical details.

## Prerequisites

- **Go 1.26 or newer.** v1.0.0 raised the minimum.
- After upgrading, run `go mod tidy` from the directory that holds your
  test module's `go.mod`. The AWS, Azure, and GCP SDK pins all moved.

## Function naming conventions

Most public Terratest helpers come in up to four variants. Knowing which
to call is the most common source of confusion when reading v1 godoc:

| Suffix | Takes `context.Context` | On error |
| --- | --- | --- |
| `Foo` | no | calls `t.Fatal` (fails the test) |
| `FooE` | no | returns `error` to the caller |
| `FooContext` | yes | calls `t.Fatal` |
| `FooContextE` | yes | returns `error` |

Two independent suffixes:

- **`E` suffix.** Long-standing Terratest convention. Use the bare name
  (`Apply`) when you want any failure to fail the test, and the `E`
  variant (`ApplyE`) when you want the error back to assert on it or
  retry.
- **`Context` suffix.** Added in v1. Takes an explicit `context.Context`
  as the second argument so callers can plumb timeouts, cancellation,
  and tracing through. The non-`Context` variants are now deprecated
  in favor of their `*Context*` counterparts.

The preferred v1 call is `FooContext` (or `FooContextE`). For example,
prefer `terraform.ApplyContext(t, ctx, opts)` over
`terraform.Apply(t, opts)`.

A small number of helpers do not yet expose all four variants (some
packages added `Context` and dropped the bare `Foo` form, others have
not grown a `Context` variant at all). Trust godoc and the deprecation
warnings over the table when they disagree.

## Migrating to the `Context` variants

The Context migration is the single largest source of deprecation
warnings you will see when upgrading. It touches nearly every helper
package: `terraform`, `helm`, `dns-helper`, `http-helper`, `packer`,
`docker`, `ssh`, `oci`, `k8s`, `aws`, `azure`, `gcp`. The non-`Context`
variants compile and behave the same as before; they just emit a
`// Deprecated:` godoc warning.

Mechanical migration:

```go
// Before
out := terraform.Apply(t, options)
out, err := terraform.ApplyE(t, options)

// After
ctx := context.Background() // or context.WithTimeout(...) for cancellation
out := terraform.ApplyContext(t, ctx, options)
out, err := terraform.ApplyContextE(t, ctx, options)
```

The `*Context*` variants always take `(t, ctx, ...originalArgs)`. If
you do not need cancellation or timeouts, `context.Background()` gives
you the same behavior as the deprecated wrapper. When `t` is
`*testing.T` (Go 1.24+), `t.Context()` is an even better default
because it ties the context lifetime to the test.

You can do this incrementally. The deprecated wrappers will keep
working for the entire v1 line; they only disappear in v2.

## What changed by service

### Azure

The largest set of breaking changes by far. The whole `modules/azure`
package was moved from the archived `services/...` SDK to the actively
maintained `sdk/resourcemanager/...` SDK, plus a handful of naming
cleanups. Updating imports and the resulting compile errors is the bulk
of the work; per-service tables and a search-and-replace cheatsheet are
in [Azure modules](./azure/).

### AWS

`modules/aws/s3.go` migrated off the deprecated
`s3/manager` package onto `s3/transfermanager`. Four exported functions
that returned `*manager.Uploader` now return `*transfermanager.Client`:
`NewS3Uploader`, `NewS3UploaderE`, `NewS3UploaderContext`, and
`NewS3UploaderContextE`. The call shape moves from
`uploader.Upload(ctx, &s3.PutObjectInput{...})` to
`client.UploadObject(ctx, &transfermanager.UploadObjectInput{...})`,
with the input/output types under
`github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager`.

Every direct `github.com/aws/aws-*` dependency was bumped to the
versions current at the v1.0.0 cut. If your tests import the AWS SDKs
directly, expect to run `go mod tidy` and resolve a small number of
type renames at the SDK level.

### GCP

`modules/gcp/pubsub.go` moved from `cloud.google.com/go/pubsub` (v1) to
`cloud.google.com/go/pubsub/v2`. The wrapper functions in `modules/gcp`
are unchanged in shape, but callers that drove the underlying client
directly need to switch from `client.Topic("name")` /
`client.Subscription("name")` handles to `TopicAdminClient` /
`SubscriptionAdminClient` calls that take fully qualified resource
names (`projects/<id>/topics/<name>`).

A new family of `*WithClient` helpers was added across `modules/gcp`
(compute, oslogin, region, pubsub, storage, cloudbuild, gcr) so tests
can inject a pre-built SDK client. This parallels the Azure
`*WithClient` change and is purely additive.

### Kubernetes

`GetKubernetesClientFromOptionsContextE` no longer falls back silently
to `rest.InClusterConfig()` when an explicit kubeconfig path or
context fails to load. It now returns the load error.

This was a silent-failure footgun: a typo in
`KubectlOptions.ConfigPath` would cause tests to run against the test
runner's in-cluster identity (potentially a different cluster) with no
error returned. If you relied on the fallback, set
`KubectlOptions.InClusterAuth = true` to opt in explicitly.

## Other deprecations you can defer

A large set of legacy spellings picked up Go-idiomatic replacements in
v1, all preserved as deprecated aliases. Most follow common acronym
casing (`Id` → `ID`, `Ip` → `IP`, `Json` → `JSON`, `Url` → `URL`,
`Ssh` → `SSH`, `Gcp` → `GCP`); a few drop redundant prefixes (Azure's
`CreateNew*Client*` becomes `Create*Client*`); and a few rename for
clarity (e.g. `SaveAmiId` / `LoadAmiId` became
`SaveArtifactID` / `LoadArtifactID` to reflect that the helpers are
not AMI-specific).

You do not need to track these one by one. Run `go vet` or
`staticcheck` against your test module after upgrading; the deprecated
aliases all carry `// Deprecated:` annotations and the linter will
list them with the replacement to use. Aliases stay for the v1 line
and are removed in v2.

## Need help

Open an issue on the [Terratest
repo](https://github.com/gruntwork-io/terratest/issues) with a snippet
of the failing code and the relevant module label. If you spot a gap
in this guide, send a PR against `docs/_docs/03_migrating-to-v1/`.
