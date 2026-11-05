# expect: error
# A valid account-provisioning workflow bound to the deprovisioning slot.
# Each slot accepts only its own type.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Wrong Type Deprovisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "create_account" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} create_account"
  steps = [{
    step   = "manual_provisioning_task_target_user"
    inputs = jsonencode({ instructions = "Sign up with your work email." })
  }]
}

resource "lumos_app_deprovisioning_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.create_account.id
}
