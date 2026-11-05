# The first app's binding is forgotten so the case can be destroyed.
# Lumos refuses to delete a workflow while it's bound, so the bound workflow is
# forgotten too, and it stays in Lumos.
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

removed {
  from = lumos_app_approval_workflow.first

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_workflow.approvals

  lifecycle {
    destroy = false
  }
}
