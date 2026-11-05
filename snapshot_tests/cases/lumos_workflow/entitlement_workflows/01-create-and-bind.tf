# Entitlement workflows on a custom app, where no integration steps exist:
# provisioning waits a day, then for the user's third day; deprovisioning
# waits a week. Both bound to the entitlement. (An approval_task isn't offered
# in an entitlement-provisioning workflow, so Lumos refuses one there.)
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Workflows"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_workflow" "grant_admin" {
  app_id = lumos_app.this.id
  type   = "entitlement-provisioning"
  name   = "bb-${var.run_id} Grant admin"
  steps = [
    {
      step   = "wait"
      inputs = jsonencode({ duration = "P1D" })
    },
    {
      step   = "wait_until"
      inputs = jsonencode({ until_expression = "add_days(user.start_date, 3)" })
    },
  ]
}

resource "lumos_workflow" "revoke_admin" {
  app_id = lumos_app.this.id
  type   = "entitlement-deprovisioning"
  name   = "bb-${var.run_id} Revoke admin"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P7D" })
  }]
}

resource "lumos_entitlement_provisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.admin.id
  workflow_id    = lumos_workflow.grant_admin.id
}

resource "lumos_entitlement_deprovisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.admin.id
  workflow_id    = lumos_workflow.revoke_admin.id
}
