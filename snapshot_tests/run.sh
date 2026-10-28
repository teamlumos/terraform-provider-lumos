#!/usr/bin/env bash
# Runs the snapshot test cases against a real Lumos API: applies each case's
# steps in order and compares every plan with the case's snapshots. See
# README.md.
#
#   ./snapshot_tests/run.sh [--update] [case ...]
#
# A case is named <resource>/<case>, and shell globs work: lumos_app/* or '*'.
set -uo pipefail

here=$(cd "$(dirname "$0")" && pwd)
update=0
patterns=()
for arg in "$@"; do
  case "$arg" in
    --update) update=1 ;;
    -h | --help) sed -n '2,8p' "$0"; exit 0 ;;
    -*) echo "unknown option: $arg" >&2; exit 2 ;;
    *) patterns+=("$arg") ;;
  esac
done
[ ${#patterns[@]} -eq 0 ] && patterns=('*')

for name in LUMOS_API_TOKEN LUMOS_API_URL; do
  if [ -z "${!name:-}" ]; then
    echo "$name must be set; see snapshot_tests/README.md" >&2
    exit 2
  fi
done
for tool in terraform perl curl jq; do
  command -v "$tool" >/dev/null || { echo "$tool not found" >&2; exit 2; }
done

export LUMOS_ACCESS_TOKEN=$LUMOS_API_TOKEN
export TF_IN_AUTOMATION=1 TF_INPUT=0
export TF_VAR_server_url=$LUMOS_API_URL
if [ -n "${GROUP_ID:-}" ]; then export TF_VAR_group_id=$GROUP_ID; fi
if [ -n "${LUMOS_TEST_TBA_LABEL:-}" ]; then export TF_VAR_tba_label=$LUMOS_TEST_TBA_LABEL; fi

passed=0 failed=0 skipped=0
# The provider's set attributes, which normalize.pl puts in a stable order.
# Read from the provider's schema with the first case's workdir.
sets=$(mktemp)
recorded_replans=()
work='' run_id='' case_name='' snapshots=''

# With a dev_overrides entry for the provider, terraform runs the local build
# without init; init would try to download the pinned version instead.
uses_dev_override() {
  local rc=${TF_CLI_CONFIG_FILE:-$HOME/.terraformrc}
  [ -f "$rc" ] && grep -q 'dev_overrides' "$rc" && grep -q 'teamlumos/lumos' "$rc"
}

# Runs a terraform command in the case's workdir. Its output goes to last.log
# and is appended to terraform.log.
tf() {
  local cmd=$1 rc
  shift
  (cd "$work" && terraform "$cmd" -no-color "$@") >"$work/last.log" 2>&1
  rc=$?
  cat "$work/last.log" >>"$work/terraform.log"
  return "$rc"
}

# Prints the names of the provider's set attributes, nested ones included,
# from the schema of the provider in the case's workdir.
set_attributes() {
  (cd "$work" && terraform providers schema -json 2>>terraform.log) |
    jq -r '[.provider_schemas[].resource_schemas[].block | .. | objects
      | select(has("attributes")) | .attributes | to_entries[]
      | select(((.value.type | type) == "array" and .value.type[0] == "set")
        or .value.nested_type.nesting_mode == "set")
      | .key] | unique[]'
}

report() {
  echo "  $1"
  echo "  last lines of $work/terraform.log:"
  tail -n 20 "$work/terraform.log" | sed 's/^/    /'
}

# Prints the variables named on "# requires:" lines in the case that aren't set.
missing_vars() {
  local name
  for name in $(sed -n 's/^# requires://p' "$1"/*.tf | tr ' ' '\n' | sort -u); do
    [ -n "${!name:-}" ] || printf '%s ' "$name"
  done
}

# Compares a normalized plan with its snapshot, or records it with --update.
compare() {
  local got=$1 want=$snapshots/$2.txt label=$case_name/snapshots/$2.txt
  if [ "$update" = 1 ]; then
    if ! cmp -s "$got" "$want"; then
      mkdir -p "$snapshots"
      cp "$got" "$want"
      echo "  updated $label"
    fi
    return 0
  fi
  if [ ! -f "$want" ]; then
    echo "  missing $label (record it with --update)"
    return 1
  fi
  diff -u --label "$label" --label "actual" "$want" "$got" | sed 's/^/  /'
  return "${PIPESTATUS[0]}"
}

# Plans (with any extra plan arguments), saves the plan as tfplan and compares
# it with the snapshot named by the first argument. Returns 1 when the snapshot
# doesn't match and 2 when the plan itself fails.
check_plan() {
  local snapshot=$1
  shift
  rm -f "$work/tfplan"
  if ! tf plan -out=tfplan "$@"; then
    report "$snapshot: plan failed"
    return 2
  fi
  show_plan "$snapshot"
  compare "$work/$snapshot.txt" "$snapshot"
}

# Writes the saved plan, normalized, to <snapshot>.txt in the workdir.
show_plan() {
  (cd "$work" && terraform show -no-color tfplan 2>>terraform.log) |
    perl "$here/normalize.pl" "$work/uuids" "$run_id" "$sets" >"$work/$1.txt"
}

# Plans again after an apply. An empty plan passes. A plan with changes fails,
# unless the case has a snapshot for it: a diff the provider is known to leave.
check_replan() {
  local snapshot=$1 label=$case_name/snapshots/$1.txt rc
  rm -f "$work/tfplan"
  tf plan -detailed-exitcode -out=tfplan
  rc=$?
  case "$rc" in
    0)
      [ -f "$snapshots/$snapshot.txt" ] || return 0
      if [ "$update" = 1 ]; then
        rm "$snapshots/$snapshot.txt"
        echo "  removed $label"
        return 0
      fi
      echo "  $snapshot: no changes, but $label expects some"
      return 1
      ;;
    2)
      show_plan "$snapshot"
      if [ "$update" = 0 ] && [ ! -f "$snapshots/$snapshot.txt" ]; then
        echo "  $snapshot: the plan after apply isn't empty:"
        sed 's/^/    /' "$work/$snapshot.txt"
        return 1
      fi
      compare "$work/$snapshot.txt" "$snapshot"
      if [ "$update" = 1 ]; then
        echo "  WARNING: recorded $label, a plan that isn't empty after apply."
        echo "  Keep it only for a diff the provider is known to leave; otherwise delete it."
        recorded_replans+=("$label")
      fi
      ;;
    *)
      report "$snapshot: plan failed"
      return 2
      ;;
  esac
}

# Fails when a boolean output isn't true after an apply. Cases check what they
# read back with outputs like `value = data.x.y.id == resource.z.id`, whose
# values are often known only after apply, so no plan snapshot shows them.
check_outputs() {
  local false_outputs
  false_outputs=$( (cd "$work" && terraform output -json 2>>terraform.log) |
    jq -r 'to_entries[] | select(.value.type == "bool" and .value.value != true) | .key')
  if [ -n "$false_outputs" ]; then
    echo "  $1: outputs that aren't true: $(echo $false_outputs)"
    return 1
  fi
}

# Runs a step marked "# expect: error", which must fail at plan or apply. A
# plan that succeeds is compared with <step>.plan.txt, and the error with
# <step>.error.txt. Returns 1 when a snapshot doesn't match and 2 when the step
# succeeds.
check_error_step() {
  local base=$1 status=0
  rm -f "$work/tfplan"
  if tf plan -out=tfplan; then
    show_plan "$base.plan"
    compare "$work/$base.plan.txt" "$base.plan" || status=1
    if tf apply tfplan; then
      echo "  $base: expected an error, but the apply succeeded"
      return 2
    fi
  elif [ -f "$snapshots/$base.plan.txt" ]; then
    if [ "$update" = 1 ]; then
      rm "$snapshots/$base.plan.txt"
      echo "  removed $case_name/snapshots/$base.plan.txt"
    else
      echo "  $base: the plan failed, but $case_name/snapshots/$base.plan.txt expects it to succeed"
      status=1
    fi
  fi
  # The error is everything from Terraform's first "Error:" line on.
  sed -n '/^Error: /,$p' "$work/last.log" |
    perl "$here/normalize.pl" "$work/uuids" "$run_id" "$sets" >"$work/$base.error.txt"
  if [ ! -s "$work/$base.error.txt" ]; then
    report "$base: failed without an error message"
    return 1
  fi
  compare "$work/$base.error.txt" "$base.error" || status=1
  return "$status"
}

destroy() {
  if ! tf destroy -auto-approve; then
    report "destroy failed; objects may be left in Lumos (names contain $run_id)"
    return 1
  fi
}

interrupted() {
  if [ -n "$work" ]; then
    echo "interrupted; destroying $case_name"
    destroy
  fi
  rm -f "$sets"
  exit 130
}
trap interrupted INT TERM

run_case() {
  local dir=$1 step base status=0 missing rc applied=''
  missing=$(missing_vars "$dir")
  if [ -n "$missing" ]; then
    echo "SKIP $case_name (needs $missing)"
    skipped=$((skipped + 1))
    return
  fi

  echo "RUN  $case_name"
  snapshots=$dir/snapshots
  run_id=$(LC_ALL=C tr -dc 'a-z0-9' </dev/urandom | head -c 6)
  export TF_VAR_run_id=$run_id
  work=$(mktemp -d)
  cp "$here"/shared/*.tf "$work"/

  if ! uses_dev_override && ! tf init; then
    report "init failed"
    if grep -q 'no available releases match' "$work/terraform.log"; then
      echo "  The pinned version isn't on the registry. To run a local build, see"
      echo "  \"Which provider build runs\" in snapshot_tests/README.md."
    fi
    status=1
  else
    [ -s "$sets" ] || set_attributes >"$sets"
    for step in "$dir"/[0-9]*; do
      if [[ "$step" == *.sh ]]; then
        base=$(basename "$step" .sh)
        (cd "$work" && bash -euo pipefail -c '. "$1"; . "$2"' drift "$here/drift.sh" "$step") \
          >"$work/last.log" 2>&1
        rc=$?
        cat "$work/last.log" >>"$work/terraform.log"
        if [ "$rc" != 0 ]; then
          report "$base: drift step failed"
          status=1
          break
        fi
        continue
      fi
      [[ "$step" == *.tf ]] || continue
      base=$(basename "$step" .tf)
      cp "$step" "$work/main.tf"
      if grep -q '^# expect: error' "$step"; then
        check_error_step "$base"
        rc=$?
        [ "$rc" = 0 ] || status=1
        [ "$rc" = 2 ] && break
        continue
      fi
      check_plan "$base.plan"
      rc=$?
      [ "$rc" = 0 ] || status=1
      [ "$rc" = 2 ] && break
      if ! tf apply tfplan; then
        report "$base: apply failed"
        status=1
        break
      fi
      applied=$step
      check_replan "$base.replan" || status=1
      check_outputs "$base" || status=1
    done
    # Destroy from the last config that applied, not from a step that was
    # expected to fail.
    if [ -n "$applied" ]; then
      cp "$applied" "$work/main.tf"
    else
      : >"$work/main.tf"
    fi
    check_plan destroy.plan -destroy || status=1
    destroy || status=1
  fi

  if [ "$status" = 0 ]; then
    echo "PASS $case_name"
    passed=$((passed + 1))
    rm -rf "$work"
  else
    echo "FAIL $case_name (workdir kept: $work)"
    failed=$((failed + 1))
  fi
  work=''
}

for dir in "$here"/cases/*/*/; do
  dir=${dir%/}
  case_name=${dir#"$here"/cases/}
  for pattern in "${patterns[@]}"; do
    # shellcheck disable=SC2254 # patterns are globs on purpose
    case "$case_name" in
      $pattern | $pattern/*) run_case "$dir"; break ;;
    esac
  done
done

rm -f "$sets"
echo
echo "$passed passed, $failed failed, $skipped skipped"
if [ ${#recorded_replans[@]} -gt 0 ]; then
  echo
  echo "WARNING: recorded ${#recorded_replans[@]} plans that aren't empty after apply:"
  printf '  %s\n' "${recorded_replans[@]}"
fi
[ "$failed" = 0 ]
