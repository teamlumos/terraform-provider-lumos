# Helpers for drift steps: NN-name.sh files in a case, which change objects in
# Lumos outside Terraform. The runner sources this file and then the step, in
# the case's workdir, with errexit on, so any failing command fails the step.

# Prints the id of a resource in state, e.g. state_id lumos_app.this
state_id() {
  local id
  id=$(terraform show -json | jq -r --arg a "$1" \
    '.values.root_module.resources[]? | select(.address == $a) | .values.id')
  [ -n "$id" ] || { echo "no $1 in state" >&2; return 1; }
  echo "$id"
}

# Calls the API and prints the response body. Fails on a status outside 2xx.
#   api GET /v1/entitlements/<id>
#   api PATCH /v1/entitlements/<id> '{"label": "x"}'
api() {
  local method=$1 path=$2 body=${3:-}
  echo "drift: $method $path ${body:+$body}" >&2
  curl -sS --fail-with-body -X "$method" \
    -H "Authorization: Bearer $LUMOS_API_TOKEN" \
    -H 'Content-Type: application/json' \
    ${body:+--data "$body"} \
    "$LUMOS_API_URL$path"
}
