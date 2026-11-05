# expect: error
# A valid account-deprovisioning workflow bound as the app's approval
# workflow. The approval slot accepts only an approval workflow.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Wrong Type Approval"
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

resource "lumos_app_approval_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.remove_account.id
}
