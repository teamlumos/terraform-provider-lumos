# expect: error
# A one-day wait inserted between the stages is refused: an approval chain can
# only contain approval stages.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Chain Two Stages"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "two_stage" {
  app_id = lumos_app.this.id
  type   = "approval-chain"
  name   = "bb-${var.run_id} User, then admins"
  steps = [
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 1, assignees = [{ uuid = data.lumos_current_user.me.id, type = "user" }] })
    },
    {
      step   = "wait"
      inputs = jsonencode({ duration = "P1D" })
    },
    {
      step   = "approval_task"
      inputs = jsonencode({ stage_number = 2, assignees = [{ uuid = "app_admins", type = "dynamic_assignee" }] })
    },
  ]
}
