# A manual provisioning task with neither of its optional inputs.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Round Trip"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "probe" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Probe"
  steps = [{
    step   = "manual_provisioning_task"
    inputs = jsonencode({})
  }]
}
