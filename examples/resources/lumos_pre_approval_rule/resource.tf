resource "lumos_pre_approval_rule" "my_preapprovalrule" {
  app_id        = "...my_app_id..."
  justification = "...my_justification..."
  preapproval_webhooks = [
    {
      id = "...my_id..."
    }
  ]
  preapproved_groups = [
    {
      app_id                  = "...my_app_id..."
      id                      = "...my_id..."
      integration_specific_id = "...my_integration_specific_id..."
    }
  ]
  preapproved_permissions = [
    {
      id = "...my_id..."
    }
  ]
  time_based_access = [
    "4 hours"
  ]
}