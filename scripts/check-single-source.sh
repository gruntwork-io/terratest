#!/usr/bin/env bash
# check-single-source.sh — fails CI if two different go.mod files declare the
# exact same module path. That's the only state Go's loader treats as truly
# ambiguous — nested module paths (e.g. root `github.com/foo/bar` and submodule
# `github.com/foo/bar/baz/v2`) are fine, because each go.mod creates a hard
# boundary that the loader respects.

set -uo pipefail

declare -A SEEN=()
fail=0

while IFS= read -r gomod; do
  mod_path=$(grep '^module' "$gomod" | head -1 | awk '{print $2}' || true)
  [ -z "$mod_path" ] && continue

  if [ -n "${SEEN[$mod_path]:-}" ]; then
    echo "::error::Duplicate module path '$mod_path' declared in:"
    echo "    ${SEEN[$mod_path]}"
    echo "    $gomod"
    fail=1
  else
    SEEN[$mod_path]="$gomod"
  fi
done < <(find . -name go.mod -not -path '*/vendor/*' 2>/dev/null)

if [ "$fail" -ne 0 ]; then
  echo "::error::Single-source-of-truth check failed."
  exit 1
fi

echo "single-source check: OK (${#SEEN[@]} unique module paths across the repo)"
