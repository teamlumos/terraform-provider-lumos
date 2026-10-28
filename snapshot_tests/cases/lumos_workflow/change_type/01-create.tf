# An account-provisioning workflow that waits a day.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Change Type"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "this" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Wait a day"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P1D" })
  }]
}
