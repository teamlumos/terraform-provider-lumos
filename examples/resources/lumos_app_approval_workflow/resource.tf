# Makes an approval template the approval process for requests to the app.
# Lumos always keeps a workflow in this slot: change it by setting
# workflow_id, and stop managing it with a removed block (see below).
resource "lumos_app_approval_workflow" "this" {
  app_id      = lumos_app.deploy_dashboard.id
  workflow_id = lumos_workflow.deploy_dashboard_approvals.id
}

# Stop managing the binding without changing Lumos.
# removed {
#   from = lumos_app_approval_workflow.this
#   lifecycle {
#     destroy = false
#   }
# }
