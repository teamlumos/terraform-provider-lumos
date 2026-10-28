# Only the wait expression changes.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Unbound Workflow"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "onboarding" {
  app_id      = lumos_app.this.id
  type        = "account-provisioning"
  name        = "bb-${var.run_id} Contractor onboarding"
  description = "Not bound to any slot."
  steps = [
    {
      step   = "wait_until"
      inputs = jsonencode({ until_expression = "add_days(user.start_date, -1)" })
    },
    {
      step = "send_notifications"
      inputs = jsonencode({
        subject             = "Your Deploy Dashboard account is on its way"
        body                = "IT will set it up today."
        recipients          = [{ uuid = data.lumos_current_user.me.id, type = "user" }]
        fallback_recipients = [{ uuid = data.lumos_current_user.me.id, type = "user" }]
      })
    },
    {
      step = "manual_provisioning_task"
      inputs = jsonencode({
        instructions = "Create a contractor account."
        assignees    = [{ uuid = "app_admins", type = "dynamic_assignee" }]
      })
    },
  ]
}
