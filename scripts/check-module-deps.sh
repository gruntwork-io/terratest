#!/usr/bin/env bash
#
# Checks the cross-module dependency graph against an explicit allowlist, and catches stale indirect requires.
#
# check-acyclic-deps.sh enforces layering: a module may not import from a strictly higher tier. That is necessary but
# not sufficient. aws and k8s are both tier 3, so k8s importing aws for a single EC2 call never violated the tier
# rule, and it went unnoticed until a user reported that depending on k8s pulled in 23 AWS service SDKs (#1875).
#
# This script enforces the complementary rule: every cross-module edge is listed below on purpose. Adding one is a
# deliberate act with a reviewer attached, not something that happens because an import was convenient.
#
# It also catches the second half of that problem. When k8s stopped requiring aws, helm kept carrying aws as an
# indirect require, because nothing tidies submodule go.mod files: the go-mod-tidy-check workflow runs at the repo
# root and diffs only the root go.mod and go.sum. helm shipped 72 aws-sdk-go-v2 entries in its go.sum for code no
# module imported. An indirect require that is not reachable through the declared direct graph is stale.
#
# No -e: accumulate every violation and report them all, rather than aborting on the first.
set -uo pipefail

# Permitted direct cross-module requires, as "importer:importee".
#
# Test-only edges still appear here, because a test dependency is a real entry in go.mod. Keeping the graph small is
# the point of the v2 split, so treat an addition as a design decision: prefer injecting behaviour over taking the
# dependency. modules/k8s/kubectl_options.go NodePublicIPLookup is the worked example.
ALLOWED_EDGES=(
  "aws:core"
  "aws:ssh"           # Ec2Keypair embeds ssh.KeyPair; SCP helpers
  "azure:core"
  "database:core"
  "dnshelper:core"
  "docker:core"
  "docker:httphelper" # test only
  "gcp:core"
  "helm:core"
  "helm:httphelper"   # test only
  "helm:k8s"
  "httphelper:core"
  "k8s:core"
  "k8s:httphelper"    # test only
  "opa:core"
  "packer:core"
  "ssh:core"
  "terraform:core"
  "terraform:httphelper" # test only
  "terraform:opa"
  "terraform:ssh"     # Options.SshAgent
  "terragrunt:core"
  "terragrunt:terraform"
  "teststructure:core"
  "teststructure:opa"
  "teststructure:terraform"
)

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
modules_dir="${repo_root}/modules"

if [ ! -d "$modules_dir" ]; then
  echo "check-module-deps: no modules/ directory; nothing to check"
  exit 0
fi

is_allowed() {
  local edge="$1"
  for allowed in "${ALLOWED_EDGES[@]}"; do
    # Strip the trailing comment before comparing.
    [ "${allowed%%[[:space:]]*}" = "$edge" ] && return 0
  done
  return 1
}

# terratest_requires <go.mod> <direct|indirect> prints one module name per line.
terratest_requires() {
  local gomod="$1" kind="$2"
  if [ "$kind" = "direct" ]; then
    grep 'gruntwork-io/terratest/modules/' "$gomod" | grep -v '^module' | grep -v '// indirect' || true
  else
    grep 'gruntwork-io/terratest/modules/' "$gomod" | grep -v '^module' | grep '// indirect' || true
  fi | sed -E 's|.*/modules/([a-z0-9]+)/v2.*|\1|' | sort -u
}

declare -A DIRECT_EDGES=()
modules=()

for dir in "$modules_dir"/*/; do
  module="$(basename "$dir")"
  [ -f "${dir}go.mod" ] || continue
  modules+=("$module")
  DIRECT_EDGES[$module]="$(terratest_requires "${dir}go.mod" direct | tr '\n' ' ')"
done

exit_code=0
seen_edges=()

# 1. Every declared direct edge must be allowlisted.
for module in "${modules[@]}"; do
  for importee in ${DIRECT_EDGES[$module]}; do
    edge="${module}:${importee}"
    seen_edges+=("$edge")

    if ! is_allowed "$edge"; then
      echo "::error file=modules/${module}/go.mod::undeclared cross-module dependency '${edge}'. If this is intended, add it to ALLOWED_EDGES in scripts/check-module-deps.sh with a short note saying why. Prefer injecting behaviour over taking the dependency."
      exit_code=1
    fi
  done
done

# 2. Every allowlisted edge must still exist, so the list documents the real graph rather than history.
for allowed in "${ALLOWED_EDGES[@]}"; do
  edge="${allowed%%[[:space:]]*}"
  found=0

  for seen in "${seen_edges[@]}"; do
    [ "$seen" = "$edge" ] && found=1 && break
  done

  if [ "$found" -eq 0 ]; then
    echo "::error file=scripts/check-module-deps.sh::allowlisted edge '${edge}' no longer exists. Remove it from ALLOWED_EDGES so the list keeps describing the real graph."
    exit_code=1
  fi
done

# 3. Every indirect terratest require must be reachable through the declared direct graph. An unreachable one is a
#    stale entry left behind when some other module dropped the dependency, and it drags that module's whole
#    transitive tree into every consumer's go.sum.
reachable_from() {
  local start="$1"
  local -A visited=()
  local queue=("$start") current

  while [ ${#queue[@]} -gt 0 ]; do
    current="${queue[0]}"
    queue=("${queue[@]:1}")
    [ -n "${visited[$current]+set}" ] && continue
    visited[$current]=1

    for next in ${DIRECT_EDGES[$current]:-}; do
      queue+=("$next")
    done
  done

  echo "${!visited[@]}"
}

for module in "${modules[@]}"; do
  indirect="$(terratest_requires "${modules_dir}/${module}/go.mod" indirect | tr '\n' ' ')"
  [ -z "$indirect" ] && continue

  reachable=" $(reachable_from "$module") "

  for importee in $indirect; do
    if [[ "$reachable" != *" ${importee} "* ]]; then
      echo "::error file=modules/${module}/go.mod::stale indirect require '${importee}': no module in ${module}'s dependency graph requires it any more. Run 'go mod tidy' for this module and commit the result."
      exit_code=1
    fi
  done
done

if [ "$exit_code" -eq 0 ]; then
  echo "module-deps check: OK (${#seen_edges[@]} cross-module edges, all allowlisted, no stale indirects)"
fi

exit "$exit_code"
