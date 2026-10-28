# The bindings are forgotten so the case can be destroyed.
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

removed {
  from = lumos_app_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_app_deprovisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_entitlement_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}
