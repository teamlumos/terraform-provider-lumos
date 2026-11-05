# The entitlement is hidden from the App Store in Lumos.
api PATCH "/v1/entitlements/$(state_id lumos_entitlement.this)/appstore-settings" \
  '{"visible_in_appstore": false}'
