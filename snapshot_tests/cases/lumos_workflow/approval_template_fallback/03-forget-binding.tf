# The binding is forgotten so the case can be destroyed.
# Lumos refuses to delete a workflow while it's bound, so the bound workflow is
# forgotten too, and it stays in Lumos.
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

removed {
  from = lumos_app_approval_workflow.this

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
