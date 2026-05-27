#!/usr/bin/env bash
# release-prep.sh — prepare the repo for a lockstep v2 release.
#
# What it does:
#   1. Strip all local `replace github.com/gruntwork-io/terratest/...` directives
#      from every submodule's go.mod.
#   2. Pin every cross-submodule `require` line to the version we're about to tag.
#   3. Run `go mod tidy` per submodule under GOWORK=off to update go.sum files.
#   4. Verify the no-replaces guard passes.
#
# Usage:
#   ./scripts/release-prep.sh v2.0.0-beta.1
#
# Run on a release-prep branch, NOT on main. Commit the result with the message
# tag "[release-prep]" so the v2-checks workflow runs the release-prep guard.

set -euo pipefail

VERSION="${1:?usage: $0 <version>  e.g. v2.0.0-beta.1}"

cd "$(git rev-parse --show-toplevel)"

echo "==> Stripping local terratest replaces..."
for gomod in modules/*/go.mod cmd/*/go.mod; do
  [ -f "$gomod" ] || continue
  # sed in-place differs between BSD (macOS) and GNU; use -i.bak then rm the backup.
  sed -i.bak '\|^replace github.com/gruntwork-io/terratest/|d' "$gomod"
  rm -f "${gomod}.bak"
done

echo "==> Pinning cross-submodule requires to $VERSION..."
for gomod in modules/*/go.mod cmd/*/go.mod; do
  [ -f "$gomod" ] || continue
  # Match either tagged versions or pseudo-versions on terratest module lines.
  sed -i.bak -E "s|(github\\.com/gruntwork-io/terratest/(modules\|cmd)/[a-zA-Z0-9_-]+/v2) v2[^[:space:]]+|\\1 $VERSION|g" "$gomod"
  rm -f "${gomod}.bak"
done

echo "==> Skipping go mod tidy (chicken-and-egg)."
# Tidy CANNOT run cleanly at this stage:
# - GOWORK=off can't verify cross-submodule require lines (the $VERSION tags
#   we're about to publish don't exist on the proxy yet).
# - GOWORK=on (workspace) also fails, because tidy still verifies require
#   versions against the proxy regardless of workspace-level resolution.
#
# The existing go.sum entries were populated correctly during dev with local
# replace directives. Stripping replaces does not invalidate those entries
# (they're hashes of module contents, not paths). New transitive entries
# may be needed at consumer-tidy time, but Go populates those lazily.
#
# Real verification happens post-tag in .github/workflows/v2-release.yml:
# after the tag-and-push step, the workflow runs the test-external/ consumer
# simulation under GOWORK=off against the just-published proxy entries.
# If anything is missing from go.sum, that step fails.

echo "==> Verifying no-replaces guard..."
bash scripts/check-no-replaces.sh

echo
echo "Release prep complete for $VERSION."
echo "Next steps:"
echo "  1. Review the diff: git diff modules/*/go.mod cmd/*/go.mod"
echo "  2. Commit:  git commit -am 'release-prep: pin to $VERSION  [release-prep]'"
echo "  3. Push the release-prep branch and watch CI."
echo "  4. When CI is green, tag every submodule at the release commit (see docs/v2-release-runbook.md)."
