#!/usr/bin/env bash
# check-acyclic-deps.sh — fails CI if any submodule's production code imports a
# module from a strictly higher tier. Enforces the v2 layering rule:
# core → helpers → tooling → platforms → IaC, downward-only.
#
# Test files (*_test.go) are excluded; cross-module test-only imports are allowed
# (e.g. modules/core/logger/parser_test imports modules/shell/v2 — legal per the
# RFC's external _test package rule).

# No -e: we accumulate every tier violation and report them all, rather than
# aborting on the first one.
set -uo pipefail

# Tier assignment. Lower number = lower layer.
declare -A TIER=(
  [core]=0
  [ssh]=1
  [httphelper]=1
  [dnshelper]=1
  [docker]=2
  [packer]=2
  [database]=2
  [opa]=2
  [aws]=3
  [azure]=3
  [gcp]=3
  [k8s]=3
  [helm]=3
  [terraform]=4
  [terragrunt]=4
  [teststructure]=4
)

# Pre-split guard: the tier rule only holds once packages are relocated into
# their /v2 submodules. Until then the flat tree has no tiers to enforce.
if ! ls modules/*/go.mod >/dev/null 2>&1; then
  echo "acyclic-deps check: skipped (no /v2 submodules present yet)"
  exit 0
fi

fail=0

for dir in modules/*/; do
  importer=$(basename "$dir")
  # Only real submodules carry a go.mod and a tier; skip anything else.
  [ -f "$dir/go.mod" ] || continue
  if [ -z "${TIER[$importer]+set}" ]; then
    echo "::error file=${dir}::module '$importer' has no tier; add it to the TIER map in check-acyclic-deps.sh"
    fail=1
    continue
  fi
  importer_tier="${TIER[$importer]}"

  # Scan all .go files in the submodule recursively, excluding test files.
  while IFS= read -r gofile; do
    while IFS= read -r importee; do
      [ -z "$importee" ] && continue
      if [ -z "${TIER[$importee]+set}" ]; then
        echo "::error file=${gofile}::imports unknown module '$importee'; add it to the TIER map in check-acyclic-deps.sh"
        fail=1
        continue
      fi
      importee_tier="${TIER[$importee]}"
      if [ "$importee_tier" -gt "$importer_tier" ]; then
        echo "::error file=${gofile}::tier violation — $importer (tier $importer_tier) imports $importee (tier $importee_tier)"
        fail=1
      fi
    done < <(grep -oE '"github\.com/gruntwork-io/terratest/modules/[a-z][a-z0-9-]*' "$gofile" 2>/dev/null \
      | awk -F'/' '{print $NF}' \
      | sort -u)
  done < <(find "$dir" -name '*.go' -not -name '*_test.go' 2>/dev/null)
done

if [ "$fail" -ne 0 ]; then
  echo "::error::Tier-violation imports detected. Imports must flow downward only."
  exit 1
fi

echo "acyclic-deps check: OK"
