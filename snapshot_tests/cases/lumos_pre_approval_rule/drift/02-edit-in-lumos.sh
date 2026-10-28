# In Lumos, the rule's justification is edited.
api PATCH "/v1/apps/$(state_id lumos_app.this)/pre-approval-rules/$(state_id lumos_pre_approval_rule.this)" \
  '{"justification": "Edited in Lumos"}'
