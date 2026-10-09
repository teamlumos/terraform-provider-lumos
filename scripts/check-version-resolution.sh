#!/usr/bin/env bash
# Checks which provider version `terraform init` picks once 1.x versions are
# published next to 0.x, using local filesystem mirrors in a temp directory.
# Nothing is downloaded or published.
#
#   ./scripts/check-version-resolution.sh
set -euo pipefail

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

platform=$(terraform version -json | sed -n 's/.*"platform": *"\([^"]*\)".*/\1/p')
failed=0

# Makes a mirror holding the given versions and points terraform at it.
use_mirror() {
  local mirror="$work/mirror-$RANDOM" version dir
  for version in "$@"; do
    dir="$mirror/registry.terraform.io/teamlumos/lumos/$version/$platform"
    mkdir -p "$dir"
    # init only hashes this file; it never runs it.
    printf '#!/bin/sh\n' > "$dir/terraform-provider-lumos_v$version"
    chmod +x "$dir/terraform-provider-lumos_v$version"
  done
  cat > "$mirror.tfrc" <<EOF
provider_installation {
  filesystem_mirror {
    path    = "$mirror"
    include = ["registry.terraform.io/teamlumos/lumos"]
  }
}
EOF
  export TF_CLI_CONFIG_FILE="$mirror.tfrc"
  echo "mirror: $*"
}

# Prints the version `terraform init` locks for a version constraint ("" for none).
resolve() {
  local dir="$work/case-$RANDOM" constraint_line=""
  mkdir -p "$dir"
  if [ -n "$1" ]; then constraint_line="version = \"$1\""; fi
  cat > "$dir/main.tf" <<EOF
terraform {
  required_providers {
    lumos = {
      source = "teamlumos/lumos"
      $constraint_line
    }
  }
}
EOF
  if ! terraform -chdir="$dir" init -input=false -no-color > "$dir/init.log" 2>&1; then
    cat "$dir/init.log" >&2
    return 1
  fi
  sed -n 's/^ *version *= *"\([^"]*\)".*/\1/p' "$dir/.terraform.lock.hcl"
}

expect() {
  local label=$1 constraint=$2 want=$3 got
  got=$(resolve "$constraint")
  if [ "$got" = "$want" ]; then
    echo "  ok    $label: $got"
  else
    echo "  FAIL  $label: got $got, want $want"
    failed=1
  fi
}

# While V1 is in prerelease.
use_mirror 0.11.0 1.0.0-alpha.1
expect '"~> 0.10"' '~> 0.10' 0.11.0
expect 'no constraint' '' 0.11.0
expect 'exact "1.0.0-alpha.1"' '1.0.0-alpha.1' 1.0.0-alpha.1

# After 1.0.0 is published: pinned configs stay put, unpinned ones move to 1.x.
use_mirror 0.11.0 1.0.0-alpha.1 1.0.0
expect '"~> 0.10"' '~> 0.10' 0.11.0
expect 'no constraint' '' 1.0.0

exit $failed
