---
layout: collection-browser-doc
title: Rewriting imports
category: migrating-to-v2
excerpt: >-
  The mechanical part of the v2 migration, and the two places a blind find-and-replace gets it wrong.
tags: ["migration", "v2"]
order: 401
nav_title: Documentation
nav_title_link: /docs/
---

Every v2 import path changes. This page covers doing that in bulk, and the two cases where a naive replace produces
code that does not compile or, worse, compiles and is wrong.

## The rewrite

Three transformations, in this order.

**Collapse the tier-0 utilities into `core`.** Do this first, because the next rule would otherwise rewrite these
paths to something that does not exist:

```bash
find . -name '*.go' -exec sed -i '' -E \
  's|gruntwork-io/terratest/modules/(random\|files\|formatting\|logger\|shell\|retry\|testing)|gruntwork-io/terratest/modules/core/v2/\1|g' {} +
```

**Rename the three hyphenated packages,** both the path and the identifier:

```bash
find . -name '*.go' -exec sed -i '' \
  -e 's|terratest/modules/http-helper|terratest/modules/httphelper/v2|g' \
  -e 's|terratest/modules/dns-helper|terratest/modules/dnshelper/v2|g' \
  -e 's|terratest/modules/test-structure|terratest/modules/teststructure/v2|g' \
  -e 's|\bhttp_helper\.|httphelper.|g' \
  -e 's|\bdns_helper\.|dnshelper.|g' \
  -e 's|\btest_structure\.|teststructure.|g' {} +
```

**Add `/v2` to everything else:**

```bash
find . -name '*.go' -exec sed -i '' -E \
  's|gruntwork-io/terratest/modules/(aws\|azure\|gcp\|k8s\|helm\|ssh\|docker\|packer\|database\|opa\|terraform\|terragrunt)([^/]\|$)|gruntwork-io/terratest/modules/\1/v2\2|g' {} +
```

On Linux, drop the `''` after `-i`.

Then let the compiler find the rest:

```bash
go mod tidy && go build ./...
```

## Where a blind replace goes wrong

**Files that alias Terratest's `aws`.** A file importing both the AWS SDK and Terratest's `aws` package usually binds
plain `aws` to the SDK and aliases Terratest's:

```go
import (
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    terraAws "github.com/gruntwork-io/terratest/modules/aws/v2"
)
```

The path rewrite is fine here, but the [symbol relocations](#symbol-relocations) below are not: rewriting
`test_structure.LoadEc2KeyPair` to `aws.LoadEc2KeyPair` in this file resolves to the SDK and fails to compile. Use
that file's alias. This is a compile error, so you will not miss it, but it is the one place the scripted rewrite
needs a human.

**`testing` is a common word.** The first script rewrites
`gruntwork-io/terratest/modules/testing`, which is specific enough to be safe. Do not broaden it to match bare
`testing`, or you will rewrite the stdlib import in every test file.

## Symbol relocations

Eight functions moved to the module that owns the type they operate on, so `teststructure` no longer requires `aws`,
`k8s`, `packer` and `ssh`. Signatures and on-disk filenames are unchanged.

| v1 | v2 |
|---|---|
| `test_structure.{Save,Load}Ec2KeyPair` | `aws.{Save,Load}Ec2KeyPair` |
| `test_structure.{Save,Load}KubectlOptions` | `k8s.{Save,Load}KubectlOptions` |
| `test_structure.{Save,Load}PackerOptions` | `packer.{Save,Load}PackerOptions` |
| `test_structure.{Save,Load}SSHKeyPair` | `ssh.{Save,Load}SSHKeyPair` |

The target module is usually already imported, because the value being saved came from it. Where it is not, add the
import.

Everything else stays in `teststructure`: `RunTestStage`, `CopyTerraformFolderToTemp`, the Terraform option helpers,
`SaveString`/`LoadString`, `SaveInt`/`LoadInt`, `SaveArtifactID`/`LoadArtifactID`, and the generic
`SaveTestData`/`LoadTestData`.

## go.mod

One `require` per module you import, all at the same version:

```
require (
    github.com/gruntwork-io/terratest/modules/terraform/v2 v2.0.0-beta.2
    github.com/gruntwork-io/terratest/modules/aws/v2 v2.0.0-beta.2
    github.com/gruntwork-io/terratest/modules/core/v2 v2.0.0-beta.2
)
```

The modules are released in lockstep and their cross-module requires are pinned to the release version, so mixing
versions is not supported.
