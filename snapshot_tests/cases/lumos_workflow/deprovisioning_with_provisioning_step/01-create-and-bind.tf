# expect: error
# An account-deprovisioning workflow with a provisioning step, bound to the
# deprovisioning slot. It should be refused.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Deprovisioning With Provisioning Step"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "remove_account" {
  app_id = lumos_app.this.id
  type   = "account-deprovisioning"
  name   = "bb-${var.run_id} remove_account"
  steps = [{
    step = "manual_provisioning_task"
    inputs = jsonencode({
      instructions = "Create the account."
      assignees    = [{ uuid = "app_admins", type = "dynamic_assignee" }]
    })
  }]
}

resource "lumos_app_deprovisioning_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.remove_account.id
}
