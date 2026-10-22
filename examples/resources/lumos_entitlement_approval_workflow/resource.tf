# Turns on the entitlement's approval workflow override with its own approval
# template. Destroying this resource turns the override off; to stop managing
# it without turning it off, use a removed block (see below).
resource "lumos_entitlement_approval_workflow" "this" {
  entitlement_id = lumos_entitlement.deploy_admin.id
  workflow_id    = lumos_workflow.admin_approvals.id
}

# Stop managing the binding without changing Lumos.
# removed {
#   from = lumos_entitlement_approval_workflow.this
#   lifecycle {
#     destroy = false
#   }
# }
