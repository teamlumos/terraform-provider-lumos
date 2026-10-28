# The binding is forgotten so the case can be destroyed.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Bound Deprovisioning Gets Provisioning Step"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "remove_account" {
  app_id = lumos_app.this.id
  type   = "account-deprovisioning"
  name   = "bb-${var.run_id} remove_account"
  steps = [{
    step = "manual_deprovisioning_task"
    inputs = jsonencode({
      instructions = "Deactivate the account."
      assignees    = [{ uuid = "app_admins", type = "dynamic_assignee" }]
    })
  }]
}

removed {
  from = lumos_app_deprovisioning_workflow.this

  lifecycle {
    destroy = false
  }
}
