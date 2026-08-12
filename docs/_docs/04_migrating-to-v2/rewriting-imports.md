---
layout: collection-browser-doc
title: Rewriting imports
category: migrating-to-v2
excerpt: >-
  The mechanical part of the v2 migration, and the places a blind
  find-and-replace gets it wrong.
tags: ["migration", "v2"]
order: 401
nav_title: Documentation
nav_title_link: /docs/
---

Every v2 import path changes. This page covers doing that in bulk, then the
symbol relocations, then the cases a scripted rewrite cannot finish on its
own.

The commands below are written for BSD `sed`, which is what macOS ships.
On Linux, drop the `''` after `-i`. They use `#` as the delimiter rather
than `|`, because `|` collides with regex alternation, and they avoid `\b`,
which BSD `sed` does not support and silently ignores.

## 1. Collapse the utility packages into `core`

```bash
find . -name '*.go' -exec sed -i '' -E \
  's#gruntwork-io/terratest/modules/(random|files|logger|shell|retry|testing)#gruntwork-io/terratest/modules/core/v2/\1#g' {} +
```

## 2. Rename the three hyphenated packages

Path, package identifier, and any stale alias:

```bash
find . -name '*.go' -exec sed -i '' -E \
  -e 's#gruntwork-io/terratest/modules/http-helper#gruntwork-io/terratest/modules/httphelper/v2#g' \
  -e 's#gruntwork-io/terratest/modules/dns-helper#gruntwork-io/terratest/modules/dnshelper/v2#g' \
  -e 's#gruntwork-io/terratest/modules/test-structure#gruntwork-io/terratest/modules/teststructure/v2#g' \
  -e 's#(^|[^A-Za-z0-9_])http_helper\.#\1httphelper.#g' \
  -e 's#(^|[^A-Za-z0-9_])dns_helper\.#\1dnshelper.#g' \
  -e 's#(^|[^A-Za-z0-9_])test_structure\.#\1teststructure.#g' \
  -e 's#^(\t*)(http_helper|dns_helper|test_structure) "#\1"#' {} +
```

That last expression matters. These packages were commonly imported under
an explicit alias:

```go
test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
```

Rewriting only the path leaves the old alias bound to the new package, so
every rewritten call site fails with `undefined: teststructure`. The
expression drops the alias so the package name applies.

## 3. Add the `/v2` suffix to the rest

```bash
find . -name '*.go' -exec sed -i '' -E \
  's#gruntwork-io/terratest/modules/(aws|azure|gcp|k8s|helm|ssh|docker|packer|database|opa|terraform|terragrunt)([^/a-z]|$)#gruntwork-io/terratest/modules/\1/v2\2#g' {} +
```

## 4. Reformat

The rewrite reorders import paths alphabetically, so the blocks are no
longer sorted:

```bash
gofmt -w .
```

## Symbol relocations

Eight functions moved to the module that owns the type they operate on, so
`teststructure` no longer requires `aws`, `k8s`, `packer` and `ssh`.
Signatures and on-disk filenames are unchanged, so this is a qualifier
swap.

| v1 | v2 |
|---|---|
| `test_structure.{Save,Load}Ec2KeyPair` | `aws.{Save,Load}Ec2KeyPair` |
| `test_structure.{Save,Load}KubectlOptions` | `k8s.{Save,Load}KubectlOptions` |
| `test_structure.{Save,Load}PackerOptions` | `packer.{Save,Load}PackerOptions` |
| `test_structure.{Save,Load}SSHKeyPair` | `ssh.{Save,Load}SSHKeyPair` |

The target module is usually already imported, because the value being
saved came from it. Where it is not, add the import, and re-run `go mod
tidy` afterwards since this can pull in a module you did not previously
require.

Everything else stays in `teststructure`: `RunTestStage`,
`CopyTerraformFolderToTemp`, the Terraform option helpers,
`SaveString`/`LoadString`, `SaveInt`/`LoadInt`,
`SaveArtifactID`/`LoadArtifactID`, and the generic
`SaveTestData`/`LoadTestData`.

## Files that alias Terratest's `aws`

A file importing both the AWS SDK and Terratest's `aws` usually binds plain
`aws` to the SDK:

```go
import (
    "github.com/aws/aws-sdk-go-v2/aws"
    terraAws "github.com/gruntwork-io/terratest/modules/aws/v2"
)
```

Here the relocation above resolves to the SDK and fails to compile. Use
that file's alias: `terraAws.LoadEc2KeyPair`. This is the one place the
scripted rewrite needs a human, and the compiler will point at it.

## Verify

No v1 paths left, and everything still builds:

```bash
grep -rn 'gruntwork-io/terratest/modules/' --include='*.go' . | grep -v '/v2'
go test -run '^$' ./...
```

The `grep` should print nothing. Use `go test -run '^$'` rather than
`go build`: it compiles `_test.go` files without running anything, and for
a Terratest suite that is where all of your code lives.
