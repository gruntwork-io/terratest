#!/usr/bin/env bash
# check-no-replaces.sh — fails the build if a release commit still contains local
# `replace github.com/gruntwork-io/terratest/...` directives. Required before
# tagging the lockstep release.
#
# Invoked only when the workflow runs in "release-prep" mode (controlled by a
# CI env var or a manual workflow_dispatch). Dev branches keep their replaces.

set -euo pipefail

matches=$(grep -nH '^replace github.com/gruntwork-io/terratest' modules/*/go.mod cmd/*/go.mod 2>/dev/null || true)

if [ -n "$matches" ]; then
  echo "::error::Local terratest replace directives present in release commit:"
  echo "$matches" | sed 's/^/    /'
  echo "::error::Strip them before tagging (see docs/v2-release-runbook.md)."
  exit 1
fi

echo "no-replaces check: OK"
