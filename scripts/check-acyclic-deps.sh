#!/usr/bin/env bash
# check-acyclic-deps.sh — fails CI if any submodule's production code imports a
# module from a strictly higher tier. Enforces the v2 layering rule:
# core → helpers → tooling → platforms → IaC, downward-only.
#
# Test files (*_test.go) are excluded; cross-module test-only imports are allowed
# (e.g. modules/core/logger/parser_test imports modules/shell/v2 — legal per the
# RFC's external _test package rule).

set -uo pipefail

# Tier assignment. Lower number = lower layer.
declare -A TIER=(
  [core]=0
  [shell]=1
  [ssh]=1
  [http-helper]=1
  [dns-helper]=1
  [docker]=2
  [packer]=2
  [database]=2
  [slack]=2
  [oci]=2
  [opa]=2
  [aws]=3
  [azure]=3
  [gcp]=3
  [k8s]=3
  [helm]=3
  [terraform]=4
  [terragrunt]=4
  [test-structure]=4
  [version-checker]=4
)

fail=0

for dir in modules/*/; do
  importer=$(basename "$dir")
  importer_tier="${TIER[$importer]:-99}"

  # Scan all .go files in the submodule recursively, excluding test files.
  while IFS= read -r gofile; do
    while IFS= read -r importee; do
      [ -z "$importee" ] && continue
      importee_tier="${TIER[$importee]:-99}"
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
