# Workflows with no steps, for the types that allow one, each bound to its
# slot.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Empty Workflows"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "account_provisioning" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Empty account-provisioning"
  steps  = []
}

resource "lumos_workflow" "account_deprovisioning" {
  app_id = lumos_app.this.id
  type   = "account-deprovisioning"
  name   = "bb-${var.run_id} Empty account-deprovisioning"
  steps  = []
}

resource "lumos_workflow" "entitlement_provisioning" {
  app_id = lumos_app.this.id
  type   = "entitlement-provisioning"
  name   = "bb-${var.run_id} Empty entitlement-provisioning"
  steps  = []
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_app_provisioning_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.account_provisioning.id
}

resource "lumos_app_deprovisioning_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.account_deprovisioning.id
}

resource "lumos_entitlement_provisioning_workflow" "this" {
  entitlement_id = lumos_entitlement.admin.id
  workflow_id    = lumos_workflow.entitlement_provisioning.id
}
