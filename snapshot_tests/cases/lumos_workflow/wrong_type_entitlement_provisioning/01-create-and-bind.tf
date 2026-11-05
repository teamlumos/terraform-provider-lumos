# expect: error
# A valid entitlement-deprovisioning workflow bound to the entitlement's
# provisioning slot. Each slot accepts only its own type.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Wrong Type Entitlement Provisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_workflow" "revoke_admin" {
  app_id = lumos_app.this.id
  type   = "entitlement-deprovisioning"
  name   = "bb-${var.run_id} revoke_admin"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P7D" })
  }]
}

resource "lumos_entitlement_provisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.this.id
  workflow_id    = lumos_workflow.revoke_admin.id
}
