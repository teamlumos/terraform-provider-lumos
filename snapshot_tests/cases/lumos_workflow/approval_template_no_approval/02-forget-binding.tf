# The binding is kept in Lumos and forgotten, since the app's approval slot
# can't be emptied.
# Lumos refuses to delete a workflow while it's bound, so the bound workflow is
# forgotten too, and it stays in Lumos.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Template No Approval"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

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
