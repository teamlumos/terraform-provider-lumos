# The same config: the plan is empty.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Template Fallback"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "manager" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Manager approval"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = data.lumos_current_user.me.id, type = "user" }] })
    },
  ]
}

resource "lumos_workflow" "approvals" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} Everyone through the manager"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.manager.id })
  }]
}

resource "lumos_app_approval_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.approvals.id
}
