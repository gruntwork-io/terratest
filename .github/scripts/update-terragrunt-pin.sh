#!/usr/bin/env bash
# Fetch the latest Terragrunt release and update the pin in mise.toml.

set -euo pipefail

latest=$(gh api repos/gruntwork-io/terragrunt/releases/latest --jq '.tag_name' | sed 's/^v//')
current=$(sed -nE 's/^terragrunt = "([^"]+)".*/\1/p' mise.toml)
echo "latest=$latest" >> "$GITHUB_OUTPUT"
echo "current=$current" >> "$GITHUB_OUTPUT"
sed -i -E "s/^(terragrunt = \")[^\"]+(\")/\\1${latest}\\2/" mise.toml
