# The binding is forgotten so the case can be destroyed.
# Lumos refuses to delete a workflow while it's bound, so the bound workflow is
# forgotten too, and it stays in Lumos.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Account Deprovisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

removed {
  from = lumos_app_deprovisioning_workflow.this

  lifecycle {
    destroy = false
  }
}

removed {
  from = lumos_workflow.remove_account

  lifecycle {
    destroy = false
  }
}
