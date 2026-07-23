#!/usr/bin/env bash
# release-prep-pin.sh <version>   e.g. release-prep-pin.sh v2.0.0-beta.1
#
# Pins every submodule's cross-module /v2 require from the go.work dev
# placeholder (v2.0.0-00010101000000-000000000000) to the real release version,
# in place. The v2 release workflow runs this in CI right before the lockstep
# tag push; the pinned commit is what the tags point at.
#
# This intentionally breaks the go.work workspace build (the tags do not exist
# yet) — that is expected and only lasts until the tags are pushed. Root go.mod
# and go.work are left untouched (the root module is never published in v2).
set -uo pipefail

VERSION="${1:?usage: release-prep-pin.sh <version, e.g. v2.0.0-beta.1>}"
if [[ ! "$VERSION" =~ ^v2\.[0-9]+\.[0-9]+(-[0-9A-Za-z.]+)?$ ]]; then
  echo "bad version: '$VERSION' (want v2.<minor>.<patch>[-suffix])" >&2
  exit 1
fi

PLACEHOLDER='v2.0.0-00010101000000-000000000000'
count=0
for d in modules/*/; do
  [ -f "${d}go.mod" ] || continue
  mods=$(grep -oE "github.com/gruntwork-io/terratest/modules/[a-z0-9]+/v2 ${PLACEHOLDER}" "${d}go.mod" 2>/dev/null | awk '{print $1}' || true)
  for mod in $mods; do
    go -C "$d" mod edit -require="${mod}@${VERSION}"
    echo "  ${d}go.mod: ${mod} -> ${VERSION}"
    count=$((count + 1))
  done
done
echo "pinned ${count} cross-module require(s) to ${VERSION}"
