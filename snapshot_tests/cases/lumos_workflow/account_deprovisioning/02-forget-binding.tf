# The binding is forgotten so the case can be destroyed.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Account Deprovisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "remove_account" {
  app_id = lumos_app.this.id
  type   = "account-deprovisioning"
  name   = "bb-${var.run_id} Remove account"
  steps = [
    {
      step   = "wait"
      inputs = jsonencode({ duration = "PT12H" })
    },
    {
      step = "manual_deprovisioning_task"
      inputs = jsonencode({
        instructions = "Deactivate the account and transfer its dashboards."
        assignees    = [{ uuid = "app_admins", type = "dynamic_assignee" }]
      })
    },
  ]
}

removed {
  from = lumos_app_deprovisioning_workflow.this

  lifecycle {
    destroy = false
  }
}
