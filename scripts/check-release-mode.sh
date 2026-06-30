#!/usr/bin/env bash
# check-release-mode.sh — validates the v2 modules build in release mode
# (GOWORK=off) WITHOUT needing published tags.
#
# For each module it: (1) adds throwaway internal `replace`s for every sibling,
# (2) seeds the module with the ROOT go.mod's exact dependency versions so naive
# tidy cannot float a dep to an incompatible newer release, then (3) tidies.
# Finally it builds an external consumer that imports every module. All changes
# are reverted on exit; nothing is committed.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"
export GOFLAGS="${GOFLAGS:--tags=aws,azure,azure_ci_excluded,azureslim,compute,gcp,helm,kubeall,kubernetes,network}"
B="github.com/gruntwork-io/terratest/modules"
ORDER="core ssh httphelper dnshelper docker packer database opa aws azure gcp k8s helm terraform terragrunt teststructure"

# Pre-split guard: until the /v2 submodules exist, there is nothing to validate.
# This lets the gate land and run green before the modularization commit, and do
# the real release-mode build afterward.
if ! ls modules/*/go.mod >/dev/null 2>&1; then
  echo "release-mode check: skipped (no /v2 submodules present yet)"
  exit 0
fi

cleanup() {
  # Revert each path independently. `git checkout -- a b` aborts entirely if any
  # pathspec matches no tracked file, so a missing go.work.sum would otherwise
  # leave the rewritten modules/*/go.mod files dirty.
  git checkout -- modules/ >/dev/null 2>&1 || true
  git checkout -- go.work.sum >/dev/null 2>&1 || true
  git status --porcelain 2>/dev/null | awk '/^\?\?.*modules\/.*\/go\.sum$/{print $2}' | xargs -r rm -f
  [ -n "${TMP:-}" ] && rm -rf "$TMP"
}
trap cleanup EXIT

# Exact dependency versions the code is tested against, from the root module.
REQARGS=$(go mod edit -json go.mod | jq -r '.Require[]? | "-require=\(.Path)@\(.Version)"' | tr '\n' ' ')

fail=0
for m in $ORDER; do
  ( cd "modules/$m"
    for s in $ORDER; do [ "$s" = "$m" ] || go mod edit -replace="$B/$s/v2=../$s"; done
    # shellcheck disable=SC2086
    go mod edit $REQARGS
    GOWORK=off go mod tidy ) 2>/tmp/cr_$m.err || { echo "::error::release-mode tidy failed: modules/$m"; tail -4 /tmp/cr_$m.err; fail=1; }
done
[ "$fail" -ne 0 ] && { echo "release-mode check: FAILED (per-module pin)"; exit 1; }

TMP=$(mktemp -d)
{ echo "module releasecheckconsumer"; echo "go 1.26"; } > "$TMP/go.mod"
for s in $ORDER; do
  go mod edit -modfile="$TMP/go.mod" -require="$B/$s/v2@v2.0.0" -replace="$B/$s/v2=$(pwd)/modules/$s"
done
{
  echo "package main"; echo "import ("
  echo "  _ \"$B/core/v2/random\""; echo "  _ \"$B/core/v2/files\""; echo "  _ \"$B/core/v2/collections\""; echo "  _ \"$B/core/v2/formatting\""
  for s in ssh httphelper dnshelper docker packer database opa aws azure gcp k8s helm terraform terragrunt teststructure; do echo "  _ \"$B/$s/v2\""; done
  echo ")"; echo "func main() {}"
} > "$TMP/main.go"
( cd "$TMP" && GOWORK=off go mod tidy && GOWORK=off go build ./... ) 2>/tmp/cr_consumer.err \
  || { echo "::error::external consumer failed to build in release mode"; tail -10 /tmp/cr_consumer.err; echo "release-mode check: FAILED (consumer)"; exit 1; }

echo "release-mode check: OK (all 16 modules pin at root versions + external consumer builds with GOWORK=off)"
