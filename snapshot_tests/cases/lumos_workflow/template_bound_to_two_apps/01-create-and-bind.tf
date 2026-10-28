# expect: error
# One approval template bound as the approval workflow of two apps. A template
# belongs to one app, so the second binding is refused.
resource "lumos_app" "first" {
  label_override       = "bb-${var.run_id} Template Two Apps, first"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_app" "second" {
  label_override       = "bb-${var.run_id} Template Two Apps, second"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "chain" {
  app_id = lumos_app.first.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Admin approval"
  steps = [{
    step   = "approval_task"
    inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
  }]
}

resource "lumos_workflow" "approvals" {
  app_id = lumos_app.first.id
  type   = "approval"
  name   = "bb-${var.run_id} Shared approvals"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.chain.id })
  }]
}

resource "lumos_app_approval_workflow" "first" {
  app_id      = lumos_app.first.id
  workflow_id = lumos_workflow.approvals.id
}

resource "lumos_app_approval_workflow" "second" {
  app_id      = lumos_app.second.id
  workflow_id = lumos_workflow.approvals.id

  depends_on = [lumos_app_approval_workflow.first]
}
