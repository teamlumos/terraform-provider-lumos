# In Lumos, the app is bound to "Wait a week" instead.
api PUT "/v1/apps/$(state_id lumos_app.this)/default-provisioning-workflow" \
  "{\"workflow_id\": \"$(state_id lumos_workflow.week)\"}"
