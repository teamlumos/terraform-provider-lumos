# expect: error
# A valid account-deprovisioning workflow bound to the provisioning slot. Each
# slot accepts only its own type.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Wrong Slot"
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

resource "lumos_app_provisioning_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.remove_account.id
}
