# An approval chain with one approval stage assigned to a group. It runs
# through an approval template that references it, once that template is
# bound to the app or one of its entitlements.
resource "lumos_workflow" "it_approval" {
  app_id      = lumos_app.deploy_dashboard.id
  type        = "approval-chain"
  name        = "IT approval"
  description = "One approval stage, assigned to the IT group."
  steps = [{
    step = "approval_task"
    inputs = jsonencode({
      stage_number = 1
      assignees    = [{ uuid = "<group id>", type = "group" }]
    })
  }]
}
