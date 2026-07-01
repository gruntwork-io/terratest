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

# Pre-split guard: the /v2 SIV only exists once the submodules are created, so
# there is nothing meaningful to check on the flat pre-split tree.
if ! ls modules/*/go.mod >/dev/null 2>&1; then
  echo "SIV-placement check: skipped (no /v2 submodules present yet)"
  exit 0
fi

fail=0
# Derive the module set from disk so this gate can never silently go vacuous
# when modules are renamed/added/removed.
for dir in modules/*/; do
  name=$(basename "$dir")
  # Matches: "github.com/gruntwork-io/terratest/modules/<name>/<sub-pkg>/v2"
  # where <sub-pkg> is any path component(s) not equal to v2. `-o` prints only
  # the matched import token so the `/v2/` filter below inspects the import
  # itself, not the whole source line (a comment or a second import on the same
  # line must not mask a real bug).
  bugged=$(grep -rEno "\"github\.com/gruntwork-io/terratest/modules/${name}/[a-zA-Z0-9_./-]+/v2\"" \
    --include='*.go' . 2>/dev/null \
    | grep -v '/v2/' \
    || true)
  if [ -n "$bugged" ]; then
    echo "::error::SIV placement bug under modules/${name}/. The /v2 must come right after the module name, not at the package leaf:"
    sed 's/^/    /' <<< "$bugged"
    fail=1
  fi
done

if [ "$fail" -ne 0 ]; then
  echo "::error::Move the /v2 SIV to the module root in the import path."
  exit 1
fi

echo "SIV-placement check: OK"
