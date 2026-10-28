# An unbound workflow.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Workflow Deleted In Lumos"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "this" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Wait a day"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P1D" })
  }]
}
