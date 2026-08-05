# Terratest v2 Import Map

Status: FROZEN for import paths. The one open decision (renames) is resolved: the three hyphenated packages are renamed to idiomatic Go names at the `/v2` boundary.

Rewriting import paths is not sufficient on its own: some symbols also moved between modules during the beta. See [Symbols relocated during the v2 beta](#symbols-relocated-during-the-v2-beta).

Built from the actual v1 layout at tag `v1.0.1-test` (27 `modules/` packages, 2 `cmd/` binaries, 1 `internal/lib` tree).

Module base path: `github.com/gruntwork-io/terratest`

## Transformation rule

For any surviving import path, the rewrite is a prefix replacement that also applies to every subpackage:

- `modules/<name>/...` -> `modules/<name>/v2/...` (the `/v2` SIV goes after the module root; directory layout unchanged except for the three renames below)
- The six tier-0 utilities collapse under one module: `modules/<util>/...` -> `modules/core/v2/<util>/...`
- Three packages are also renamed to drop the hyphen: `http-helper` -> `httphelper`, `dns-helper` -> `dnshelper`, `test-structure` -> `teststructure`. The package identifier loses its underscore too (`http_helper` -> `httphelper`), so call sites change, not just the import path. The codemod handles both.

So e.g. `modules/logger/parser` -> `modules/core/v2/logger/parser`, `modules/aws/foo` -> `modules/aws/v2/foo`, and `modules/http-helper` -> `modules/httphelper/v2`.

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

| v1 import path | v2 import path |
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
| `modules/http-helper` | `modules/httphelper/v2` |
| `modules/dns-helper` | `modules/dnshelper/v2` |
| `modules/test-structure` | `modules/teststructure/v2` |

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

## Symbols relocated during the v2 beta

These functions moved to the module that owns the type they operate on, so that `teststructure` no longer requires
`aws`, `k8s`, `packer` and `ssh` (#1877). Name, arguments, return type and on-disk filename are unchanged. Each is a
compile error, so `go build ./...` finds every call site.

| v2.0.0-beta.1 | v2.0.0-beta.2 onwards |
|---|---|
| `teststructure.SaveEc2KeyPair` | `aws.SaveEc2KeyPair` |
| `teststructure.LoadEc2KeyPair` | `aws.LoadEc2KeyPair` |
| `teststructure.SaveKubectlOptions` | `k8s.SaveKubectlOptions` |
| `teststructure.LoadKubectlOptions` | `k8s.LoadKubectlOptions` |
| `teststructure.SavePackerOptions` | `packer.SavePackerOptions` |
| `teststructure.LoadPackerOptions` | `packer.LoadPackerOptions` |
| `teststructure.SaveSSHKeyPair` | `ssh.SaveSSHKeyPair` |
| `teststructure.LoadSSHKeyPair` | `ssh.LoadSSHKeyPair` |

Nothing else moved. `RunTestStage`, `CopyTerraformFolderToTemp`, the Terraform helpers, `SaveString`/`LoadString`,
`SaveInt`/`LoadInt`, `SaveArtifactID`/`LoadArtifactID` and the generic `SaveTestData`/`LoadTestData` all stay in
`teststructure`. The generic primitives now live in `modules/core/v2/teststate` and are re-exported, so those calls
are unaffected.

Migration is a qualifier swap, with two cautions:

- The target module is usually already imported, since the value came from it. Where it is not, add the import.
- A file importing both the AWS SDK and Terratest's `aws` may bind plain `aws` to the SDK and alias Terratest's. Use
  that file's alias.

## Behaviour changes during the v2 beta

No signature changed, but two behaviours did.

**Node addresses prefer `ExternalIP`** (#1878). `k8s.FindNodeHostnameContextE`, and `GetServiceEndpoint` for a
NodePort service, now return the `ExternalIP` on the Node object when present, falling back to the internal hostname
as before. On EKS that is the instance's public IP, so the common case needs no EC2 call and no
`ec2:DescribeInstances` permission. Clusters that advertise an `ExternalIP` now resolve to it rather than the
internal hostname.

`FindNodeHostnameContext[E]` keep their original signatures. A new `FindNodeHostnameWithOptionsContext[E]` pair takes
`*KubectlOptions` and consults `NodePublicIPLookup`, for clusters that do not advertise an `ExternalIP`:
`options.NodePublicIPLookup = aws.GetPublicIpsOfEc2InstancesContextE`.

**`KubectlOptions` carrying a `RestConfig` no longer serializes** (#1879). `MarshalJSON` returns
`k8s.ErrRestConfigNotSerializable` rather than dropping the config and leaving reloaded options with no cluster
identity. This failed before too, with an opaque `json: unsupported type: transport.WrapperFunc`. For staged tests
use `NewKubectlOptions` or `NewKubectlOptionsWithInClusterAuth`; both round trip intact.

## Accounting

27 `modules/` packages = 6 collapsed into core + 15 standalone submodules + 6 removed. Plus 2 removed `cmd/` binaries and 1 internal flatten. Submodule count: 16.

## Open decisions

None. The map is frozen.

## Resolved

- **`oci`** (Oracle Cloud Infrastructure): not carried forward to v2. Niche provider; removed alongside the other dropped packages. Oracle Cloud users stay on frozen v1.
- **Renames.** Decided to rename the three hyphenated packages at the v2 boundary: `http-helper` -> `httphelper`, `dns-helper` -> `dnshelper`, `test-structure` -> `teststructure`. Consumers already rewrite every import for the `/v2` bump, so folding the rename into that same edit adds no separate migration, and it drops the non-idiomatic underscore package names (`http_helper`, currently suppressed with `//nolint:staticcheck`). The rename rides along with the modularization import rewrite, and the codemod covers both the path and the package identifier.
