---
layout: collection-browser-doc
title: Import map
category: migrating-to-v2
excerpt: >-
  Every v1 import path and what it becomes in v2.
tags: ["migration", "v2"]
order: 402
nav_title: Documentation
nav_title_link: /docs/
---

The complete v1 to v2 path mapping. See [rewriting imports]({{site.baseurl}}/docs/migrating-to-v2/rewriting-imports/)
for how to apply it in bulk.

## Collapsed into `core`

Six utility packages became subpackages of one module. The package identifier at call sites is unchanged.

| v1 | v2 |
|---|---|
| `modules/random` | `modules/core/v2/random` |
| `modules/files` | `modules/core/v2/files` |
| `modules/logger` | `modules/core/v2/logger` |
| `modules/shell` | `modules/core/v2/shell` |
| `modules/retry` | `modules/core/v2/retry` |
| `modules/testing` | `modules/core/v2/testing` |

`modules/logger/parser` becomes `modules/core/v2/logger/parser`. `internal/lib/formatting` moved to
`internal/formatting`, which was never importable.

## Renamed

Path *and* package identifier change.

| v1 | v2 |
|---|---|
| `modules/http-helper`, `http_helper.X` | `modules/httphelper/v2`, `httphelper.X` |
| `modules/dns-helper`, `dns_helper.X` | `modules/dnshelper/v2`, `dnshelper.X` |
| `modules/test-structure`, `test_structure.X` | `modules/teststructure/v2`, `teststructure.X` |

## Suffix only

Path gains `/v2`; package identifier unchanged.

| v1 | v2 |
|---|---|
| `modules/aws` | `modules/aws/v2` |
| `modules/azure` | `modules/azure/v2` |
| `modules/gcp` | `modules/gcp/v2` |
| `modules/k8s` | `modules/k8s/v2` |
| `modules/helm` | `modules/helm/v2` |
| `modules/ssh` | `modules/ssh/v2` |
| `modules/docker` | `modules/docker/v2` |
| `modules/packer` | `modules/packer/v2` |
| `modules/database` | `modules/database/v2` |
| `modules/opa` | `modules/opa/v2` |
| `modules/terraform` | `modules/terraform/v2` |
| `modules/terragrunt` | `modules/terragrunt/v2` |

## Removed

| v1 | replacement |
|---|---|
| `modules/collections` | stdlib `slices` |
| `modules/environment` | stdlib `os.Getenv` |
| `modules/git` | stdlib `os/exec` |
| `modules/slack` | none; vendor from v1 if needed |
| `modules/version-checker` | none; shell out |
| `modules/oci` | none; Oracle Cloud is not carried forward to v2 |
| `cmd/pick-instance-type` | none |
| `cmd/terratest_log_parser` | none as a binary; the library survives at `modules/core/v2/logger/parser` |

These were deprecated in v1 first and deleted at the v2 cutover. If you depend on `slack`, `version-checker` or
`oci`, v1 stays available and is the place to stay.
