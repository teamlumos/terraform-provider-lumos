# The bindings are forgotten so the case can be destroyed.
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
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
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

removed {
  from = lumos_entitlement_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_entitlement_deprovisioning_workflow.this

  lifecycle {
    destroy = false
  }
}
