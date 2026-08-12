---
layout: collection-browser-doc
title: v2 overview
category: migrating-to-v2
excerpt: >-
  What changes when you move from Terratest v1 to v2, and the order to do
  it in.
tags: ["migration", "v2"]
order: 400
nav_title: Documentation
nav_title_link: /docs/
---

Terratest v2 splits the single `github.com/gruntwork-io/terratest` module
into 16 independent modules, so you depend only on the parts you use. A
test that imports `terraform` no longer pulls in the AWS SDK, client-go,
and every other provider's dependencies.

The cost is that every import path changes. This page tells you what the
changes are and the order to apply them; the per-topic pages hold the
details.

## Should you migrate yet

v2 is in beta. v1 is in maintenance and receives security fixes only,
until 12 months after v2.0.0 reaches general availability. Migrate now if
you want
the smaller dependency graph or are starting fresh; wait for v2.0.0 if you
would rather not track beta releases.

You can migrate incrementally. v1 and v2 import paths differ, so both can
coexist in one module while you convert package by package. The change
touches only `.go` files, `go.mod` and `go.sum`, so `git checkout` undoes
it.

The minimum Go version is unchanged.

## The five changes

**1. Everything v1 deprecated is gone.** v1 kept deprecated aliases
alongside their replacements; v2 deletes them. This is the largest edit in
the migration, and it is not only the `Context` variants:

- non-`Context` wrappers: `terraform.Apply` is now only `ApplyContext`
- initialism renames: `random.UniqueId` is now `UniqueID`,
  `aws.GetAccountIdE` is now `GetAccountIDContextE`
- reshaped helpers: `packer.BuildAmi` is now `BuildArtifactContextE`

The reliable way to find all of it is to run staticcheck against your
existing v1 code and clear every SA1019 (deprecated symbol) warning before
you touch imports. Once v1 is warning-free, the rest of this guide applies.

```go
// v1
out := terraform.Apply(t, options)
// v2
out := terraform.ApplyContext(t, t.Context(), options)
```

The `Context` variants always take `(t, ctx, ...originalArgs)`.
`t.Context()` is the best default; `context.Background()` also works. The
[v1 guide]({{site.baseurl}}/docs/migrating-to-v1/overview/) covers this
migration in detail, and doing it on v1 first, where both forms still
compile, is easier than doing it at the same time as the import rewrite.

**2. Import paths gain a `/v2` suffix.** The `/v2` goes after the module
root, not at the end of the path:

```go
// v1
"github.com/gruntwork-io/terratest/modules/terraform"
// v2
"github.com/gruntwork-io/terratest/modules/terraform/v2"
```

**3. Six utility packages collapse into `core`.** `random`, `files`,
`logger`, `shell`, `retry` and `testing` are no longer separate packages:

```go
// v1
"github.com/gruntwork-io/terratest/modules/random"
// v2
"github.com/gruntwork-io/terratest/modules/core/v2/random"
```

**4. Three packages are renamed to drop the hyphen.** This changes the
package identifier at call sites, not just the import path:

`http-helper` becomes `httphelper`, `dns-helper` becomes `dnshelper`, and
`test-structure` becomes `teststructure`. The full table is in the [import
map]({{site.baseurl}}/docs/migrating-to-v2/import-map/#renamed).

**5. Each module needs its own `require`.** Where v1 was one line in
`go.mod`, v2 needs one per module you import. Add them with `go get`
rather than by hand:

```bash
go get github.com/gruntwork-io/terratest/modules/terraform/v2@v2.0.0-beta.2
go get github.com/gruntwork-io/terratest/modules/aws/v2@v2.0.0-beta.2
```

Each module is tagged `modules/<name>/vX.Y.Z`, so the tag for the command
above is `modules/terraform/v2.0.0-beta.2`. The current version is on the
[releases page](https://github.com/gruntwork-io/terratest/releases). All 16
modules are released together and their cross-module requires are pinned to
the release version, so keep them on the same version.

The full path mapping is in [the import
map]({{site.baseurl}}/docs/migrating-to-v2/import-map/).

## Order to do it in

1. Move to the `Context` variants, ideally while still on v1 so both forms
   compile.
2. Rewrite import paths and package identifiers. This is mechanical and the
   compiler finds everything.
3. Fix the symbol relocations. Also compile errors, also mechanical. They
   are listed under [rewriting
   imports]({{site.baseurl}}/docs/migrating-to-v2/rewriting-imports/).
4. Add a `require` per module, then `go mod tidy`.
5. Review the [behavior
   changes]({{site.baseurl}}/docs/migrating-to-v2/behavior-changes/), which
   the compiler will *not* find for you. There are two, both in `k8s`.

Steps 1 to 4 are compiler-detectable, so the build tells you when they are
complete. Step 5 is not: the code compiles either way, so it needs reading
before you call the migration done.

## Removed packages

Six packages and two binaries are not carried forward. Three have a
standard library replacement:

| v1 | Replacement |
|---|---|
| `modules/collections` | stdlib `slices` |
| `modules/environment` | stdlib `os.Getenv` |
| `modules/git` | stdlib `os/exec` |
| `modules/slack` | none; vendor from v1 |
| `modules/version-checker` | none; shell out |
| `modules/oci` | none; Oracle Cloud is not carried forward, stay on v1 |

`cmd/pick-instance-type` and `cmd/terratest_log_parser` are gone as
binaries. The log parser's library survives at
`modules/core/v2/logger/parser`.

## What did not change

The `Foo` / `FooE` convention is unchanged: `FooContext` fails the test,
`FooContextE` returns an error. Beyond dropping the non-`Context` wrappers
and the relocations above, function arguments and return types are the
same. Test data written by v1 loads in v2, since filenames and JSON layout
are unchanged.

## Need help

Open an issue on the [Terratest
repo](https://github.com/gruntwork-io/terratest/issues) with the version
you are coming from and the error you hit. If you spot a gap in this
guide, send a PR against `docs/_docs/04_migrating-to-v2/`.
