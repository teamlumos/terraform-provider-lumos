# Same config as step 01: the plan creates the workflow again, instead of
# failing to read it.
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
