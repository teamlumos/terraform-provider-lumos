# expect: error
# An approval stage without its required assignees.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Bad Steps"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "chain" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} Edited badly"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1 })
    },
  ]
}
