# expect: error
# The bound workflow's step is swapped for a deprovisioning step. It should be
# refused, as at create.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Bound Provisioning Gets Deprovisioning Step"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "create_account" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} create_account"
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
  workflow_id = lumos_workflow.create_account.id
}
