# In Lumos, the business justification is edited and the app is pre-approved.
# The update takes only input fields, so the read is cut down to them.
id=$(state_id lumos_access_policy.this)
body=$(api GET "/v1/access-policies/$id" | jq -c '{
  name, status, is_enabled, access_condition,
  business_justification: "Edited in Lumos",
  apps: [.apps[] | {id, is_preapproved: true, entitlements: [.entitlements[]? | {id}]}]
}')
api PUT "/v1/access-policies/$id" "$body"
