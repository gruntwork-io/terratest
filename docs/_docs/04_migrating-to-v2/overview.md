---
layout: collection-browser-doc
title: Overview
category: migrating-to-v2
excerpt: >-
  Migrate from Terratest v1 to v2 — import paths, codemod, FAQ.
tags: ["migration", "v2"]
order: 400
nav_title: Documentation
nav_title_link: /docs/
---

Terratest v2 splits the library into per-domain Go submodules. The runtime
API is unchanged; only import paths move. Migration is mechanical and can
be done file by file. Existing v1 consumers stay on v1 until they choose to migrate.

Requires Go 1.26+ (older versions build via Go's automatic toolchain).

## Import path changes

**Eight tier-0 utilities collapse into `modules/core/v2`** as sub-packages:

| v1 path | v2 path |
| --- | --- |
| `modules/collections` | `modules/core/v2/collections` |
| `modules/random` | `modules/core/v2/random` |
| `modules/files` | `modules/core/v2/files` |
| `modules/testing` | `modules/core/v2/testing` |
| `modules/environment` | `modules/core/v2/environment` |
| `modules/git` | `modules/core/v2/git` |
| `modules/retry` | `modules/core/v2/retry` |
| `modules/logger` (+`/parser`) | `modules/core/v2/logger` (+`/parser`) |

**Every other module gets a `/v2` suffix:** `modules/aws/v2`, `modules/azure/v2`, `modules/gcp/v2`, `modules/k8s/v2`, `modules/helm/v2`, `modules/terraform/v2`, `modules/terragrunt/v2`, `modules/test-structure/v2`, `modules/shell/v2`, `modules/ssh/v2`, `modules/http-helper/v2`, `modules/dns-helper/v2`, `modules/version-checker/v2`, `modules/docker/v2`, `modules/packer/v2`, `modules/database/v2`, `modules/slack/v2`, `modules/oci/v2`, `modules/opa/v2`.

### `/v2` placement rule

`/v2` goes at the **module root** in the import path, not at the package leaf:

- Correct: `modules/logger/v2/parser`
- Bugged: `modules/logger/parser/v2`

This is Go's [SIV rule](https://go.dev/ref/mod#major-version-suffixes). The codemod enforces it.

## Codemod

```bash
git clone --depth 1 https://github.com/gruntwork-io/terratest /tmp/terratest
bash /tmp/terratest/scripts/v2-import-rewrite.sh /path/to/your/project
cd /path/to/your/project && go mod tidy && go build ./...
```

Idempotent; safe to re-run.

## Example

Before:
```go
import "github.com/gruntwork-io/terratest/modules/random"
import "github.com/gruntwork-io/terratest/modules/terraform"
```

After:
```go
import "github.com/gruntwork-io/terratest/modules/core/v2/random"
import "github.com/gruntwork-io/terratest/modules/terraform/v2"
```

`go.mod` requires `modules/core/v2 v2.0.0` and `modules/terraform/v2 v2.0.0` instead of one `terratest v1.0.0`.

## FAQ

- **Can I mix v1 and v2?** Yes, in the same `go.mod` and even the same package. v1 and v2 are different Go modules.
- **`go get -u`?** Stays within the current major version. Won't accidentally jump to v2.
- **API changes?** None. Same signatures, same behavior.
- **`cmd/` binaries?** `go install github.com/gruntwork-io/terratest/cmd/<name>/v2@latest`.
- **v1 support?** Frozen at the last v1 tag; security backports on the `v1-maintenance` branch. Final EOL announced separately.
