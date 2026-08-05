#!/usr/bin/env bash
#
# Checks the cross-module dependency graph against an explicit allowlist, and catches stale indirect requires.
#
# check-acyclic-deps.sh enforces layering only. aws and k8s are both tier 3, so k8s importing aws for one EC2 call
# never violated it (#1875). This requires every edge to be listed on purpose instead.
#
# It also catches staleness: nothing tidies submodule go.mod files, so when k8s dropped aws, helm kept it as an
# indirect and shipped 72 aws-sdk-go-v2 go.sum entries for code no module imported.
#
# No -e: accumulate every violation rather than aborting on the first.
set -uo pipefail

# Permitted direct cross-module requires, as "importer:importee". Test-only edges count, since a test dependency is
# a real go.mod entry. Treat an addition as a design decision: prefer injecting behaviour over taking the dependency
# (see NodePublicIPLookup in modules/k8s/kubectl_options.go).
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

# 3. Every indirect terratest require must be reachable through the declared direct graph. An unreachable one is
#    stale, and drags a whole transitive tree into every consumer's go.sum.
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
