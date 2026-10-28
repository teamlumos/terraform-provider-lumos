# The binding is forgotten so the case can be destroyed.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Template Branches"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "it_only" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} IT only"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
    },
  ]
}

resource "lumos_workflow" "manager_then_it" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Manager, then IT"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = data.lumos_current_user.me.id, type = "user" }] })
    },
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 2, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
    },
  ]
}

resource "lumos_workflow" "approvals" {
  app_id = lumos_app.this.id
  type   = "approval"
  name   = "bb-${var.run_id} Branch on team"
  steps = [{
    step = "branch"
    inputs = jsonencode({
      branches = [{
        condition = {
          and = [{
            attribute = { label = "Team", value = "identity.team" }
            in        = [{ label = "Engineering", value = "Eng" }, { label = "IT", value = "IT" }]
          }]
        }
        workflow_id = lumos_workflow.it_only.id
      }]
      fallback_workflow_id = lumos_workflow.manager_then_it.id
    })
  }]
}

removed {
  from = lumos_app_approval_workflow.this

  lifecycle {
    destroy = false
  }
}
