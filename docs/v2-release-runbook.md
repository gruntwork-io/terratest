# Terratest v2 Release Runbook

How to cut a coordinated release of the v2 submodules. Read this before tagging.

The tag push and proxy verification are automated by the `v2 Release` workflow
(`.github/workflows/v2-release.yml`, run via workflow_dispatch). This runbook
explains the procedure that workflow follows and how to do it by hand if needed.

## Layout

v2 is split into per-domain modules under `modules/<name>/`, each declaring
`module github.com/gruntwork-io/terratest/modules/<name>/v2`. Local development
uses the root `go.work`, which resolves every submodule to its local tree, so the
submodule `go.mod` files do not need internal `require` lines or `replace`
directives during normal development.

## Why pinning is a release-time step (do not commit it early)

The submodules' `go.mod` files are deliberately left without cross-module
`require` lines on `main`. Pinning a sibling `require` to the to-be-published
version (e.g. `core/v2 v2.0.0-beta.1`) BREAKS the workspace build until that tag
actually exists: `go.work` does not shadow an unpublished required version, so
`go build` and `go work sync` try to fetch the missing revision and fail. The pin
must therefore happen on a short-lived release-prep branch, immediately before the
tags are pushed, never on the modularization PR.

CI validates the release-mode build continuously without committing the pin: the
`GOWORK=off` check generates the pinned state with throwaway `replace` directives,
builds a consumer, and discards it (see `scripts/`).

## Pre-flight (on a release-prep branch)

1. Choose the version, e.g. `v2.0.0-beta.1`.
2. Pin each module, in dependency order (core first, then helpers, tooling,
   platforms, k8s/helm, IaC). For each module:
   - Add a temporary `replace` for EVERY sibling it transitively needs, not just
     its direct imports. Tidy follows transitive edges, so a partial replace set
     fails with `unknown revision` on a deeper sibling.
   - Seed the module with the root `go.mod`'s exact dependency versions before
     tidying (copy its `require` lines in). Otherwise tidy floats deps to the
     latest compatible release and you get silent version drift, e.g. the Azure
     SDK advancing to a release that drops a symbol the code uses
     (`undefined: armmonitor.DiagnosticSettingsClient`). The code is tested
     against the root's pinned versions; the submodules must inherit them.
   - `GOWORK=off go mod tidy` to populate external `require`s and `go.sum`.

`scripts/check-release-mode.sh` performs this exact pin (transitive replaces +
root-version seeding + an all-module external consumer build under `GOWORK=off`)
in a throwaway and reverts it, and runs in CI on every PR so the release-mode
build is validated continuously without committing the pin or needing tags.
3. Set every internal `require` to the exact version being tagged, then DROP all
   internal `replace` directives. Do not run `go work sync` against the unpinned
   tree.
4. CI guard: run `bash scripts/check-no-replaces.sh`; it must pass before tagging.
   It catches both single-line and block-form `replace ( ... )` directives, which
   a plain `grep '^replace'` would miss.
5. Move `test/` to its own module here too if not already done, and pin it the
   same way (it is test-only, so committed `replace`s are acceptable for it).

## Tag push order

All tags point at the same release commit. Each tag name puts the `/v2` SIV in the
tag itself: `modules/<name>/v2/<version>` (NOT `modules/<name>/<version>`, which
the proxy cannot associate with the `/v2` module path). Push in dependency order:

1. `modules/core/v2/v2.0.0-beta.1`
2. helpers: `ssh`, `httphelper`, `dnshelper`
3. tooling: `docker`, `packer`, `database`, `opa`
4. platforms: `aws`, `azure`, `gcp`, then `k8s`, `helm`
5. IaC: `terraform`, `terragrunt`, `teststructure`

After each tier, probe the proxy before continuing:
`curl -o /dev/null -w '%{http_code}' https://proxy.golang.org/github.com/gruntwork-io/terratest/modules/<name>/v2/@v/<version>.info`
should return 200.

## Verify

Build a throwaway external consumer that imports every published module with
`GOWORK=off`. This is what `scripts/check-release-mode.sh` generates in a temp
directory, except now it resolves against the real published tags rather than
local `replace`s. A clean consumer should resolve, build, and test green with
zero local references.

## If a tag is wrong

Proxy tags are immutable. Recover by cutting the next patch (`v2.0.0-beta.2`),
never by editing in place. The pre-flight checks exist to keep this rare.

## Beta to GA

After the beta soaks (suggested two weeks minimum), repeat the same procedure at
`v2.0.0` with no suffix: same release commit shape, same tag sequence, same proxy
verification, then announce.
