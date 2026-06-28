#!/usr/bin/env bash
# check-siv-placement.sh — fails CI if any Go import places the /v2 SIV at the
# package leaf instead of the module root.
#
# Correct:  github.com/gruntwork-io/terratest/modules/logger/v2/parser
# Bugged:   github.com/gruntwork-io/terratest/modules/logger/parser/v2
#
# This is the bug that stalled PR #1632. We grep for the bugged pattern under
# every known submodule's top-level package directory.

set -euo pipefail

# Submodules whose top-level package is the module root (with /v2 SIV).
# Format: short-name|sub-packages-list (space-separated). For each, check that
# imports of any sub-package put /v2 at the module-root, not at the leaf.
fail=0
# Derive the module set from disk so this gate can never silently go vacuous
# when modules are renamed/added/removed.
for dir in modules/*/; do
  name=$(basename "$dir")
  # Matches: "github.com/gruntwork-io/terratest/modules/<name>/<sub-pkg>/v2"
  # where <sub-pkg> is any path component(s) not equal to v2.
  bugged=$(grep -rEn "\"github\.com/gruntwork-io/terratest/modules/${name}/[a-zA-Z0-9_./-]+/v2\"" \
    --include='*.go' . 2>/dev/null \
    | grep -v '/v2/' \
    || true)
  if [ -n "$bugged" ]; then
    echo "::error::SIV placement bug detected under modules/${name}/ — /v2 must come right after the module name, not at the package leaf:"
    echo "$bugged" | sed 's/^/    /'
    fail=1
  fi
done

if [ "$fail" -ne 0 ]; then
  echo "::error::Run scripts/v2-import-rewrite.sh to auto-fix."
  exit 1
fi

echo "SIV-placement check: OK"
