#!/usr/bin/env bash
# v2-import-rewrite.sh — rewrites terratest module imports for the v2 modularization.
#
# Enforces the rule that the /v2 SIV goes at the MODULE ROOT in the import path
# (e.g. modules/logger/v2/parser), NOT at the package leaf (modules/logger/parser/v2).
# That was the PR #1632 bug.
#
# Usage: ./scripts/v2-import-rewrite.sh [path]
#   default path: .
#
# Each rewrite rule is one line: <OLD_PREFIX>|<NEW_PREFIX>
# Apply in the order listed. Longest-prefix-first to avoid accidentally rewriting
# a sub-prefix (e.g. ".../modules/logger/parser" must be rewritten before ".../modules/logger").

set -euo pipefail

TARGET="${1:-.}"

# Order matters: more-specific paths before less-specific.
# Each line:  <OLD>|<NEW>
RULES=(
  # core collapse — sub-packages first, then bare module
  "github.com/gruntwork-io/terratest/modules/logger/parser|github.com/gruntwork-io/terratest/modules/core/v2/logger/parser"
  "github.com/gruntwork-io/terratest/modules/logger|github.com/gruntwork-io/terratest/modules/core/v2/logger"
  "github.com/gruntwork-io/terratest/modules/retry|github.com/gruntwork-io/terratest/modules/core/v2/retry"
  "github.com/gruntwork-io/terratest/modules/git|github.com/gruntwork-io/terratest/modules/core/v2/git"
  "github.com/gruntwork-io/terratest/modules/collections/v2|github.com/gruntwork-io/terratest/modules/core/v2/collections"
  "github.com/gruntwork-io/terratest/modules/random/v2|github.com/gruntwork-io/terratest/modules/core/v2/random"
  "github.com/gruntwork-io/terratest/modules/files/v2|github.com/gruntwork-io/terratest/modules/core/v2/files"
  "github.com/gruntwork-io/terratest/modules/testing/v2|github.com/gruntwork-io/terratest/modules/core/v2/testing"
  "github.com/gruntwork-io/terratest/modules/environment/v2|github.com/gruntwork-io/terratest/modules/core/v2/environment"

  # split modules — bare module name becomes module/v2
  "github.com/gruntwork-io/terratest/modules/ssh|github.com/gruntwork-io/terratest/modules/ssh/v2"
  "github.com/gruntwork-io/terratest/modules/shell|github.com/gruntwork-io/terratest/modules/shell/v2"
  "github.com/gruntwork-io/terratest/modules/http-helper|github.com/gruntwork-io/terratest/modules/http-helper/v2"
  "github.com/gruntwork-io/terratest/modules/dns-helper|github.com/gruntwork-io/terratest/modules/dns-helper/v2"
  "github.com/gruntwork-io/terratest/modules/version-checker|github.com/gruntwork-io/terratest/modules/version-checker/v2"
  "github.com/gruntwork-io/terratest/modules/docker|github.com/gruntwork-io/terratest/modules/docker/v2"
  "github.com/gruntwork-io/terratest/modules/packer|github.com/gruntwork-io/terratest/modules/packer/v2"
  "github.com/gruntwork-io/terratest/modules/database|github.com/gruntwork-io/terratest/modules/database/v2"
  "github.com/gruntwork-io/terratest/modules/slack|github.com/gruntwork-io/terratest/modules/slack/v2"
  "github.com/gruntwork-io/terratest/modules/oci|github.com/gruntwork-io/terratest/modules/oci/v2"
  "github.com/gruntwork-io/terratest/modules/opa|github.com/gruntwork-io/terratest/modules/opa/v2"
  "github.com/gruntwork-io/terratest/modules/k8s|github.com/gruntwork-io/terratest/modules/k8s/v2"
  "github.com/gruntwork-io/terratest/modules/helm|github.com/gruntwork-io/terratest/modules/helm/v2"
  "github.com/gruntwork-io/terratest/modules/aws|github.com/gruntwork-io/terratest/modules/aws/v2"
  "github.com/gruntwork-io/terratest/modules/azure|github.com/gruntwork-io/terratest/modules/azure/v2"
  "github.com/gruntwork-io/terratest/modules/gcp|github.com/gruntwork-io/terratest/modules/gcp/v2"
  "github.com/gruntwork-io/terratest/modules/terraform|github.com/gruntwork-io/terratest/modules/terraform/v2"
  "github.com/gruntwork-io/terratest/modules/terragrunt|github.com/gruntwork-io/terratest/modules/terragrunt/v2"
  "github.com/gruntwork-io/terratest/modules/test-structure|github.com/gruntwork-io/terratest/modules/test-structure/v2"
)

# Lint check: refuse to run if any rewrite target contains the bugged pattern.
for rule in "${RULES[@]}"; do
  new="${rule#*|}"
  if [[ "$new" =~ /[a-zA-Z_-]+/v2/[a-zA-Z_-]+/v2 ]]; then
    echo "::error::Rule target looks malformed (double /v2): $new" >&2
    exit 2
  fi
done

files_changed=0
total_replacements=0

while IFS= read -r -d '' file; do
  changed=0
  for rule in "${RULES[@]}"; do
    old="${rule%|*}"
    new="${rule#*|}"
    # Use \b-style boundary by matching the closing quote (imports are always quoted)
    # to avoid prefix-matching issues.
    if grep -qF "\"$old\"" "$file" 2>/dev/null; then
      sed -i '' "s|\"$old\"|\"$new\"|g" "$file"
      changed=1
      total_replacements=$((total_replacements + 1))
    fi
  done
  if [ "$changed" = 1 ]; then
    files_changed=$((files_changed + 1))
    echo "rewrote $file"
  fi
done < <(find "$TARGET" -name '*.go' -not -path '*/vendor/*' -print0)

echo
echo "summary: $files_changed files modified, $total_replacements rules applied"
