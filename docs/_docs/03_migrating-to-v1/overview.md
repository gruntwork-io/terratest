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

This page lists what changed at the v0.x to v1.0.0 boundary and points at
the per-service guides where the change set is large.

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
  retry. Both variants exist for almost every helper.
- **`Context` suffix.** Added in v1. Takes an explicit `context.Context`
  as the second argument so callers can plumb timeouts, cancellation,
  and tracing through. The non-`Context` variants are now deprecated and
  internally call the `*Context*` variant with `context.Background()`.

The preferred v1 call is `FooContext` (or `FooContextE`). For example,
prefer `terraform.ApplyContext(t, ctx, opts)` over `terraform.Apply(t, opts)`.

## Migrating to the `Context` variants

The Context migration is the single largest source of deprecation
warnings you will see when upgrading. Every helper that runs an external
command, makes an SDK call, or sleeps got a `*Context` / `*ContextE`
companion. The non-`Context` variants compile and behave identically;
they just emit a `// Deprecated:` godoc warning and forward to the
`*Context*` variant with `context.Background()`.

Affected packages (rough scope):

- `modules/terraform`, `modules/helm`, `modules/dns-helper`
- `modules/k8s`, `modules/aws`, `modules/azure`, `modules/gcp`

Mechanical migration:

```go
// Before
out := terraform.Apply(t, options)
out, err := terraform.ApplyE(t, options)

// After
ctx := t.Context() // Go 1.24+; or context.Background() / context.WithTimeout(...)
out := terraform.ApplyContext(t, ctx, options)
out, err := terraform.ApplyContextE(t, ctx, options)
```

The order of arguments in `*Context*` variants is always
`(t, ctx, ...originalArgs)`. If you do not need cancellation or
timeouts yet, passing `context.Background()` is fine and gives you the
same behavior as the deprecated wrapper while moving you off the
deprecation warning.

You can do this incrementally. The deprecated wrappers will keep working
for the entire v1 line; they only disappear in v2.

## What changed by service

### Azure

The largest set of breaking changes. The whole `modules/azure` package
was moved from the archived `services/...` SDK to the actively maintained
`sdk/resourcemanager/...` SDK, plus a handful of naming cleanups landed
in the same release. See [Azure modules](./azure/) for the full migration
guide. Highlights:

- `services/...` imports become `sdk/resourcemanager/.../arm<service>`.
- Resource fields move under `.Properties`.
- Iterator-based list calls become pagers.
- 8 `Get*ClientE` getters were removed; use the `Create*ClientE`
  replacements that have been around for a while.
- 4 `CreateNew*ClientE` factories were renamed to `Create*ClientE` (the
  old names remain as deprecated aliases).
- `NsgRuleSummary.SourceAdresssPrefixes` (triple-s typo) renamed to
  `SourceAddressPrefixes`.
- New `*WithClient` functions accept a pre-built SDK client for
  injection in unit tests.

### AWS

`modules/aws/s3.go` migrated off the deprecated `s3/manager` package onto
`s3/transfermanager`. Four exported functions changed return type:

- `NewS3Uploader`
- `NewS3UploaderE`
- `NewS3UploaderWithSession`
- `NewS3UploaderWithSessionE`

These now return `*transfermanager.Client` instead of
`*manager.Uploader`. The call shape moves from
`uploader.Upload(ctx, &s3.PutObjectInput{...})` to
`client.UploadObject(ctx, &transfermanager.UploadObjectInput{...})`. The
input and output types live in
`github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager`.

All other AWS SDK v2 service deps (acm, autoscaling, cloudwatchlogs,
dynamodb, ec2, ecr, ecs, iam, kms, lambda, rds, route53, s3,
secretsmanager, sns, sqs, ssm, sts, config, credentials) were bumped to
the versions current at v1.0.0 cut. If your tests import those SDKs
directly, run `go mod tidy` after upgrading Terratest.

### GCP

`modules/gcp/pubsub.go` moved from `cloud.google.com/go/pubsub` (v1) to
`cloud.google.com/go/pubsub/v2`. The v1 client is deprecated upstream and
trips `staticcheck` SA1019.

The shape of the wrapper functions in `modules/gcp` is unchanged, but
callers that drove the underlying client directly need to switch from
`client.Topic("name")` / `client.Subscription("name")` handles to
`TopicAdminClient` / `SubscriptionAdminClient` calls that take fully
qualified resource names (`projects/<id>/topics/<name>`).

`modules/gcp/compute.go` also moved several free functions onto receiver
methods of `Instance`, `ZonalInstanceGroup`, and `RegionalInstanceGroup`
(e.g. `GetPublicIP(t, instance)` becomes `instance.GetPublicIP(t)`). The
free-function forms are kept as deprecated wrappers; switch when
convenient.

### Kubernetes

`GetKubernetesClientFromOptionsContextE` no longer falls back silently to
`rest.InClusterConfig()` when an explicit kubeconfig path or context
fails to load. It now returns the underlying `LoadAPIClientConfigE`
error.

This was a silent-failure footgun: a typo in `KubectlOptions.ConfigPath`
would cause tests to run against the test runner's in-cluster identity
(potentially a different cluster) with no error. If you relied on that
fallback, set `KubectlOptions.InClusterAuth = true` to opt in
explicitly.

## Other deprecations you can defer

A handful of smaller renames also landed with `// Deprecated:` aliases:

- `modules/azure`: `CreateNew*Client*` factories deprecated in favor of
  `Create*Client*` (the redundant `New` is dropped).
- `modules/ssh`: `SshSession` and `SshConnectionOptions` renamed to
  `SSHSession` and `SSHConnectionOptions` (Go-idiomatic acronym
  casing). The old names remain as deprecated type aliases.
- `modules/terraform`: a few legacy spellings (e.g. `CtyJsonOutput` →
  `CtyJSONOutput`).
- `modules/test-structure`: SSH-key and artifact-ID save/load helpers
  picked up consistent names (`SaveSSHKeyPair`, `LoadSSHKeyPair`,
  `SaveArtifactID`, `LoadArtifactID`).

These are pure renames: the old names forward to the new ones and stay
in the v1 line.

## Need help

Open an issue on the [Terratest
repo](https://github.com/gruntwork-io/terratest/issues) with a snippet
of the failing code and the relevant module label. If you spot a gap
in this guide, send a PR against `docs/_docs/03_migrating-to-v1/`.
