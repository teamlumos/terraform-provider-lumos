# Makes an entitlement-provisioning workflow the one that grants the
# entitlement. Lumos can't empty this slot: change it by setting workflow_id,
# and stop managing it with a removed block (see below).
resource "lumos_entitlement_provisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.deploy_admin.id
  workflow_id    = lumos_workflow.grant_admin.id
}

# Stop managing the binding without changing Lumos.
# removed {
#   from = lumos_entitlement_provisioning_workflow.this
#   lifecycle {
#     destroy = false
#   }
# }
