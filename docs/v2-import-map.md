# Terratest v2 Import Map

Status: DRAFT, pending one decision (see "Open decisions"). Not yet frozen.

Built from the actual v1 layout at tag `v1.0.1-test` (27 `modules/` packages, 2 `cmd/` binaries, 1 `internal/lib` tree).

Module base path: `github.com/gruntwork-io/terratest`

## Transformation rule

For any surviving import path, the rewrite is a prefix replacement that also applies to every subpackage:

- `modules/<name>/...` -> `modules/<name>/v2/...` (the `/v2` SIV goes after the module root, directory layout unchanged)
- The six tier-0 utilities collapse under one module: `modules/<util>/...` -> `modules/core/v2/<util>/...`

So e.g. `modules/logger/parser` -> `modules/core/v2/logger/parser`, and `modules/aws/foo` -> `modules/aws/v2/foo`.

## core collapse (6 v1 packages -> one `modules/core/v2`)

| v1 import path | v2 import path |
|---|---|
| `modules/logger` | `modules/core/v2/logger` |
| `modules/testing` | `modules/core/v2/testing` |
| `modules/retry` | `modules/core/v2/retry` |
| `modules/random` | `modules/core/v2/random` |
| `modules/files` | `modules/core/v2/files` |
| `modules/shell` | `modules/core/v2/shell` |

## Standalone `/v2` submodules

| v1 import path | v2 import path | note |
|---|---|---|
| `modules/aws` | `modules/aws/v2` | |
| `modules/azure` | `modules/azure/v2` | |
| `modules/gcp` | `modules/gcp/v2` | |
| `modules/k8s` | `modules/k8s/v2` | |
| `modules/helm` | `modules/helm/v2` | |
| `modules/ssh` | `modules/ssh/v2` | |
| `modules/docker` | `modules/docker/v2` | |
| `modules/packer` | `modules/packer/v2` | |
| `modules/database` | `modules/database/v2` | |
| `modules/opa` | `modules/opa/v2` | |
| `modules/terraform` | `modules/terraform/v2` | |
| `modules/terragrunt` | `modules/terragrunt/v2` | |
| `modules/http-helper` | `modules/httphelper/v2` | rename (provisional) |
| `modules/dns-helper` | `modules/dnshelper/v2` | rename (provisional) |
| `modules/test-structure` | `modules/teststructure/v2` | rename (provisional) |

## Removed in v2.0.0 (deprecated in v1 first, deleted at cutover)

| v1 import path | replacement |
|---|---|
| `modules/collections` | stdlib `slices` |
| `modules/environment` | stdlib `os.Getenv` |
| `modules/git` | stdlib `os/exec` |
| `modules/slack` | none, vendor from frozen v1 if needed |
| `modules/version-checker` | none, shell out |
| `modules/oci` | none, Oracle Cloud support not carried forward; remains in frozen v1, vendor if needed |
| `cmd/pick-instance-type` | none, standalone binary, out of scope |
| `cmd/terratest_log_parser` | none, standalone binary (its `logger/parser` lib survives under `modules/core/v2/logger/parser`) |

## Internal flatten (non-importable, not consumer-facing)

| v1 | v2 |
|---|---|
| `internal/lib/formatting` | `internal/formatting` |

## Accounting

27 `modules/` packages = 6 collapsed into core + 15 standalone submodules + 6 removed. Plus 2 removed `cmd/` binaries and 1 internal flatten. Submodule count: 16.

## Open decisions (must close before freezing)

1. **Renames.** `http-helper`/`dns-helper`/`test-structure` rewrites are marked provisional. The team is still weighing whether to do renames at all in v2. This must be settled before the map is frozen, because the map encodes the final names.

## Resolved

- **`oci`** (Oracle Cloud Infrastructure): not carried forward to v2. Niche provider; removed alongside the other dropped packages. Oracle Cloud users stay on frozen v1.
