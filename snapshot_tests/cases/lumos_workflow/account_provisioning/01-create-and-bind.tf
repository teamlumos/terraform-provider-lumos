# A realistic account-provisioning workflow: tell the requester, then give IT
# a manual task. Bound as the app's account-provisioning workflow.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Account Provisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "create_account" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Create account"
  steps = [
    {
      step = "send_notifications"
      inputs = jsonencode({
        subject             = "Access to Deploy Dashboard is being set up"
        body                = "IT will create your account within a day."
        recipients          = [{ uuid = data.lumos_current_user.me.id, type = "user" }]
        fallback_recipients = [{ uuid = data.lumos_current_user.me.id, type = "user" }]
      })
    },
    {
      step = "manual_provisioning_task"
      inputs = jsonencode({
        instructions = "Create the account in the admin console, then mark this task done."
        assignees    = [{ uuid = data.lumos_current_user.me.id, type = "user" }]
      })
    },
  ]
}

resource "lumos_workflow" "create_account_target_user" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Requester creates account"
  steps = [{
    step   = "manual_provisioning_task_target_user"
    inputs = jsonencode({ instructions = "Sign up with your work email." })
  }]
}

resource "lumos_app_provisioning_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.create_account.id
}
