# expect: error
# A chain whose only stage is numbered 2. Lumos checks stage numbers when the
# chain is created, so the chain itself is refused.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Template Stage Gap"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "gap" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Starts at stage 2"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 2, assignees = [{ uuid = data.lumos_current_user.me.id, type = "user" }] })
    },
  ]
}

resource "lumos_workflow" "approvals" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} Uses the gap chain"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.gap.id })
  }]
}

resource "lumos_app_approval_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.approvals.id
}
