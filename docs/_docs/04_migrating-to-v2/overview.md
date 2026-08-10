---
layout: collection-browser-doc
title: Overview
category: migrating-to-v2
excerpt: >-
  What changes when you move from Terratest v1 to v2, and the order to do it in.
tags: ["migration", "v2"]
order: 400
nav_title: Documentation
nav_title_link: /docs/
---

Terratest v2 splits the single `github.com/gruntwork-io/terratest` module into 16 independent modules, so you depend
only on the parts you use. A test that imports `terraform` no longer pulls in the AWS SDK, client-go, and every other
provider's dependencies.

The cost is that every import path changes. This page tells you what the changes are and the order to apply them; the
per-topic pages hold the details.

## Should you migrate yet

v2 is in beta. v1 remains supported and gets fixes. Migrate now if you want the smaller dependency graph or are
starting fresh; wait for v2.0.0 if you would rather not track beta releases.

Note that v2 modules are versioned and released in lockstep, so all 16 carry the same version. Pin them together.

## The four changes

**1. Import paths gain a `/v2` suffix.** The semantic import version goes after the module root, not the end of the
path:

```go
// v1
"github.com/gruntwork-io/terratest/modules/terraform"
// v2
"github.com/gruntwork-io/terratest/modules/terraform/v2"
```

**2. Six utility packages collapse into `core`.** `random`, `files`, `formatting`, `logger`, `shell`, `retry` and
`testing` are no longer separate packages:

```go
// v1
"github.com/gruntwork-io/terratest/modules/random"
// v2
"github.com/gruntwork-io/terratest/modules/core/v2/random"
```

**3. Three packages are renamed to drop the hyphen.** This changes the package identifier at call sites, not just the
import path:

| v1 | v2 |
|---|---|
| `modules/http-helper`, `http_helper.X` | `modules/httphelper/v2`, `httphelper.X` |
| `modules/dns-helper`, `dns_helper.X` | `modules/dnshelper/v2`, `dnshelper.X` |
| `modules/test-structure`, `test_structure.X` | `modules/teststructure/v2`, `teststructure.X` |

**4. Each module needs its own `require`.** Where v1 was one line in `go.mod`, v2 needs one per module you import.

The full path mapping is in [the import map]({{site.baseurl}}/docs/migrating-to-v2/import-map/).

## Order to do it in

1. Rewrite import paths and package identifiers. This is mechanical and the compiler finds everything.
2. Add a `require` per module, then `go mod tidy`.
3. Fix the symbol relocations. Also compile errors, also mechanical. See
   [rewriting imports]({{site.baseurl}}/docs/migrating-to-v2/rewriting-imports/#symbol-relocations).
4. Review the [behaviour changes]({{site.baseurl}}/docs/migrating-to-v2/behaviour-changes/), which the compiler will
   *not* find for you. There are two, both in `k8s`.

Steps 1 to 3 are safe: if it builds, you are done. Step 4 is the one that needs reading.

## Removed packages

Six packages and two binaries are not carried forward. Most have a stdlib replacement:

| v1 | replacement |
|---|---|
| `modules/collections` | stdlib `slices` |
| `modules/environment` | stdlib `os.Getenv` |
| `modules/git` | stdlib `os/exec` |
| `modules/slack` | none; vendor from v1 |
| `modules/version-checker` | none; shell out |
| `modules/oci` | none; Oracle Cloud is not carried forward, stay on v1 |

`cmd/pick-instance-type` and `cmd/terratest_log_parser` are gone as binaries. The log parser's library survives at
`modules/core/v2/logger/parser`.

## What did not change

Function names, arguments and return types are unchanged except for the relocations noted above. The `Foo` /
`FooE` convention is unchanged. Test data written by v1 loads in v2: filenames and JSON layout are the same.

If you already migrated to v1 and adopted the `Context` variants, there is nothing further to do on that front. If
you have not, see [migrating to v1]({{site.baseurl}}/docs/migrating-to-v1/overview/) first; doing both at once is
harder than doing them in sequence.

## Need help

Open an issue at [gruntwork-io/terratest](https://github.com/gruntwork-io/terratest/issues) with the version you are
coming from and the error you hit.
