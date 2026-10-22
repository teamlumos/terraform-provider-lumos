# Makes an account-provisioning workflow the app's account workflow for
# creating accounts. Lumos always keeps a workflow in this slot: change it by
# setting workflow_id, and stop managing it with a removed block (see below).
resource "lumos_app_provisioning_workflow" "this" {
  app_id      = lumos_app.deploy_dashboard.id
  workflow_id = lumos_workflow.create_account.id
}

# Stop managing the binding without changing Lumos.
# removed {
#   from = lumos_app_provisioning_workflow.this
#   lifecycle {
#     destroy = false
#   }
# }
