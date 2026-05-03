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

## What changed by service

### Azure

The largest set of breaking changes. The whole `modules/azure` package was
moved from the archived `services/...` SDK to the actively maintained
`sdk/resourcemanager/...` SDK, with a few naming cleanups landed in the
same release. See [Azure modules](./azure/) for the full migration guide.

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

### Kubernetes

`GetKubernetesClientFromOptionsContextE` no longer falls back silently to
`rest.InClusterConfig()` when an explicit kubeconfig path or context fails
to load. It now returns the underlying `LoadAPIClientConfigE` error.

This was a silent-failure footgun: a typo in `KubectlOptions.ConfigPath`
would cause tests to run against the test runner's in-cluster identity
(potentially a different cluster) with no error. If you relied on that
fallback, set `KubectlOptions.InClusterAuth = true` to opt in
explicitly.

## Deprecations you can defer

Throughout v1, replaced symbols carry a `// Deprecated:` godoc annotation
pointing at the new name. Examples:

- Non-`Context` variants in `modules/terraform`, `modules/helm`,
  `modules/dns-helper`, etc. (`Destroy`, `Show`, `RunTerraformCommand`, ...)
  are deprecated in favor of `*Context` / `*ContextE` variants that take
  an explicit `context.Context`.
- `CreateNew*Client*` factories in `modules/azure` are deprecated in
  favor of `Create*Client*` (the redundant `New` is dropped).

These keep working for the entire v1 line. Migrate at your convenience;
removal is a v2 concern.

## Need help

Open an issue on the [Terratest
repo](https://github.com/gruntwork-io/terratest/issues) with a snippet of
the failing code and the relevant module label. If you spot a gap in
this guide, send a PR against `docs/_docs/03_migrating-to-v1/`.
