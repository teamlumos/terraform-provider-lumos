# In Lumos, a setting the config sets is changed, and redirect_url, which the
# config leaves out, is set.
api PATCH "/v1/apps/$(state_id lumos_app.this)/appstore-settings" \
  '{"request_instructions": "Edited in Lumos", "allow_multiple_permission_selection": true, "redirect_url": "https://example.com/edited-in-lumos"}'
