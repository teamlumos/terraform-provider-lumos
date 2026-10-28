# The binding is kept in Lumos and forgotten, since the app's approval slot
# can't be emptied.
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

removed {
  from = lumos_app_approval_workflow.this

  lifecycle {
    destroy = false
  }
}
