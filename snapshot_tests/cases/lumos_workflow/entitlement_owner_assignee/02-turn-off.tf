# The override is turned off by removing its binding.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Owner"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_workflow" "owner_approval" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Entitlement owner approval"
  steps = [{
    step   = "approval_task"
    inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "entitlement_owner", type = "dynamic_assignee" }] })
  }]
}

resource "lumos_workflow" "owner_approvals" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} Entitlement owner approvals"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.owner_approval.id })
  }]
}
