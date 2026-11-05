# Someone tries to delete the bound workflow in Lumos. Lumos refuses, since the
# workflow is the app's default provisioning workflow, so nothing drifts.
id=$(state_id lumos_workflow.this)
if out=$(api DELETE "/v1/workflows/$id" 2>&1); then
  echo "Lumos deleted the bound workflow; it should have refused" >&2
  exit 1
fi
echo "$out" | grep -q -E 'workflow\.(default|bound)_delete_not_allowed' || {
  echo "unexpected response: $out" >&2
  exit 1
}
