# Makes an account-deprovisioning workflow the app's account workflow for
# removing accounts. Lumos always keeps a workflow in this slot: change it by
# setting workflow_id, and stop managing it with a removed block (see below).
resource "lumos_app_deprovisioning_workflow" "this" {
  app_id      = lumos_app.deploy_dashboard.id
  workflow_id = lumos_workflow.remove_account.id
}

# Stop managing the binding without changing Lumos.
# removed {
#   from = lumos_app_deprovisioning_workflow.this
#   lifecycle {
#     destroy = false
#   }
# }
