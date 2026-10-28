# expect: error
# One chain referenced by two templates. A chain can be used by only one
# workflow, so the second template is refused.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Chain Two Templates"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "chain" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Shared chain"
  steps = [{
    step   = "approval_task"
    inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
  }]
}

resource "lumos_workflow" "first" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} First template"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.chain.id })
  }]
}

resource "lumos_workflow" "second" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} Second template"
  steps = [{
    step   = "branch"
    inputs = jsonencode({ branches = [], fallback_workflow_id = lumos_workflow.chain.id })
  }]

  depends_on = [lumos_workflow.first]
}
