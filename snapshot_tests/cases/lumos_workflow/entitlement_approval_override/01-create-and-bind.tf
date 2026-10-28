# An entitlement with its own approval template, through the approval
# override binding.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Override"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "admins" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Admin approval"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
    },
  ]
}

resource "lumos_workflow" "admin_approvals" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} Admin approvals"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.admins.id })
  }]
}

resource "lumos_entitlement" "admin" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Admin"
}

resource "lumos_entitlement_approval_workflow" "this" {
  entitlement_id = lumos_entitlement.admin.id
  workflow_id    = lumos_workflow.admin_approvals.id
}
