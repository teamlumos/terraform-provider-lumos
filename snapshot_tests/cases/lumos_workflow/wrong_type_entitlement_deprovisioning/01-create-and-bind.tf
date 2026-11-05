# expect: error
# A valid entitlement-provisioning workflow bound to the entitlement's
# deprovisioning slot. Each slot accepts only its own type.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Wrong Type Entitlement Deprovisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_workflow" "grant_admin" {
  app_id = lumos_app.this.id
  type   = "entitlement-provisioning"
  name   = "bb-${var.run_id} grant_admin"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P1D" })
  }]
}

resource "lumos_entitlement_deprovisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.this.id
  workflow_id    = lumos_workflow.grant_admin.id
}
