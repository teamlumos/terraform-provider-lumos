# The binding is forgotten instead, so the case can be destroyed.
# Lumos refuses to delete a workflow while it's bound, so the bound workflow is
# forgotten too, and it stays in Lumos.
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
        instructions = "Create the account with the standard role, then mark this task done."
        assignees    = [{ uuid = data.lumos_current_user.me.id, type = "user" }]
      })
    },
  ]
}

removed {
  from = lumos_app_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_workflow.create_account_target_user

  lifecycle {
    destroy = false
  }
}
