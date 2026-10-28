# A template that needs no approval (no branches, no fallback), bound as the
# app's approval workflow.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Template No Approval"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "approvals" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} No approval needed"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [] })
  }]
}

resource "lumos_app_approval_workflow" "this" {
  app_id      = lumos_app.this.id
  workflow_id = lumos_workflow.approvals.id
}
