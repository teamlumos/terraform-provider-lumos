# The entitlement is renamed in Lumos.
api PATCH "/v1/entitlements/$(state_id lumos_entitlement.this)" \
  "{\"label\": \"bb-$TF_VAR_run_id Renamed In Lumos\"}"
