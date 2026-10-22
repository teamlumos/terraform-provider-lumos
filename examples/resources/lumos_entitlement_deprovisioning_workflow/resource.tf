# Makes an entitlement-deprovisioning workflow the one that revokes the
# entitlement. Lumos can't empty this slot: change it by setting workflow_id,
# and stop managing it with a removed block (see below).
resource "lumos_entitlement_deprovisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.deploy_admin.id
  workflow_id    = lumos_workflow.revoke_admin.id
}

# Stop managing the binding without changing Lumos.
# removed {
#   from = lumos_entitlement_deprovisioning_workflow.this
#   lifecycle {
#     destroy = false
#   }
# }
