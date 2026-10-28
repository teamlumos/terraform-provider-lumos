# A key set to null counts as left out, so this should plan and read back like
# a task with no assignees.
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
    inputs = jsonencode({ instructions = "Create the account.", assignees = null })
  }]
}
