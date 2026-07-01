#!/usr/bin/env bash
# check-release-mode.sh — validates the v2 modules build in release mode
# (GOWORK=off) WITHOUT needing published tags.
#
# For each module it: (1) adds throwaway internal `replace`s for every sibling,
# (2) seeds the module with the ROOT go.mod's exact dependency versions so naive
# tidy cannot float a dep to an incompatible newer release, then (3) tidies.
# Finally it builds an external consumer that imports every module. All changes
# are reverted on exit; nothing is committed.
# No -e: main() accumulates a fail flag across all modules and reports them all.
set -uo pipefail

MODULE_BASE="github.com/gruntwork-io/terratest/modules"
# Every submodule that currently has its own go.mod. Derived from disk so the
# gate works during an incremental split (core only) and after the full split.
ORDER=$(for d in modules/*/; do [ -f "${d}go.mod" ] && basename "$d"; done | tr '\n' ' ')

# Scratch directory (a fresh `mktemp -d`) holding per-module tidy stderr and the
# throwaway consumer module. Set by main(); removed wholesale by cleanup().
WORKDIR=""

cleanup() {
  # Revert each path independently. `git checkout -- a b` aborts entirely if any
  # pathspec matches no tracked file, so a missing go.work.sum would otherwise
  # leave the rewritten modules/*/go.mod files dirty.
  git checkout -- modules/ >/dev/null 2>&1 || true
  git checkout -- go.work.sum >/dev/null 2>&1 || true
  git status --porcelain 2>/dev/null | awk '/^\?\?.*modules\/.*\/go\.sum$/{print $2}' | xargs -r rm -f
  [ -n "$WORKDIR" ] && rm -rf "$WORKDIR"
}

# Emit the consumer's import block from the single ORDER list: core exposes leaf
# packages, every other module is imported at its /v2 root.
consumer_imports() {
  local s pkg
  for s in $ORDER; do
    if [ "$s" = core ]; then
      # Keep in sync with core/v2's public leaf packages.
      for pkg in random files collections formatting logger shell retry testing; do
        echo "  _ \"$MODULE_BASE/core/v2/$pkg\""
      done
    else
      echo "  _ \"$MODULE_BASE/$s/v2\""
    fi
  done
}

# Refuse to run against a dirty tree: the cleanup trap reverts modules/ wholesale
# via `git checkout`, which would discard a developer's uncommitted module work.
ensure_clean_worktree() {
  local dirty
  dirty=$(git status --porcelain -- modules go.work.sum)
  if [ -n "$dirty" ]; then
    echo "::error::release-mode check needs a clean modules/ and go.work.sum worktree (cleanup reverts them). Commit or stash first:"
    sed 's/^/    /' <<< "$dirty"
    return 1
  fi
}

main() {
  local root reqargs fail=0 m s
  root=$(git rev-parse --show-toplevel)
  cd "$root" || { echo "::error::cannot cd to repo root '$root'"; return 1; }

  export GOFLAGS="${GOFLAGS:--tags=aws,azure,azure_ci_excluded,azureslim,compute,gcp,helm,kubeall,kubernetes,network}"

  # Pre-split guard: until the /v2 submodules exist, there is nothing to validate.
  # This lets the gate land and run green before the modularization commit, and do
  # the real release-mode build afterward.
  if ! ls modules/*/go.mod >/dev/null 2>&1; then
    echo "release-mode check: skipped (no /v2 submodules present yet)"
    return 0
  fi

  ensure_clean_worktree || return 1

  trap cleanup EXIT
  WORKDIR=$(mktemp -d)

  # Exact dependency versions the code is tested against, from the root module.
  reqargs=$(go mod edit -json go.mod | jq -r '.Require[]? | "-require=\(.Path)@\(.Version)"' | tr '\n' ' ')

  for m in $ORDER; do
    ( cd "modules/$m"
      for s in $ORDER; do [ "$s" = "$m" ] || go mod edit -replace="$MODULE_BASE/$s/v2=../$s"; done
      # shellcheck disable=SC2086
      go mod edit $reqargs
      GOWORK=off go mod tidy ) 2>"$WORKDIR/tidy_$m.err" \
      || { echo "::error::release-mode tidy failed: modules/$m"; tail -4 "$WORKDIR/tidy_$m.err"; fail=1; }
  done
  [ "$fail" -ne 0 ] && { echo "release-mode check: FAILED (per-module pin)"; return 1; }

  # Build an external consumer that imports every module, in release mode.
  { echo "module releasecheckconsumer"; echo "go 1.26"; } > "$WORKDIR/go.mod"
  for s in $ORDER; do
    go mod edit -modfile="$WORKDIR/go.mod" -require="$MODULE_BASE/$s/v2@v2.0.0" -replace="$MODULE_BASE/$s/v2=$root/modules/$s"
  done
  { echo "package main"; echo "import ("; consumer_imports; echo ")"; echo "func main() {}"; } > "$WORKDIR/main.go"
  ( cd "$WORKDIR" && GOWORK=off go mod tidy && GOWORK=off go build ./... ) 2>"$WORKDIR/consumer.err" \
    || { echo "::error::external consumer failed to build in release mode"; tail -10 "$WORKDIR/consumer.err"; echo "release-mode check: FAILED (consumer)"; return 1; }

  echo "release-mode check: OK ($(echo "$ORDER" | wc -w | tr -d ' ') module(s) pin at root versions + external consumer builds with GOWORK=off)"
}

main "$@"
