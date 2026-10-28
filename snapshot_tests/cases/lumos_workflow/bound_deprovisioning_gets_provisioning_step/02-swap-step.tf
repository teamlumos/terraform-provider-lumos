# expect: error
# The bound workflow's step is swapped for a provisioning step. It should be
# refused, as at create.
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
