resource "lumos_workflow" "my_workflow" {
  app_id      = "...my_app_id..."
  description = "...my_description..."
  name        = "...my_name..."
  steps = [
    {
      inputs = {
        # ...
      }
      step = "...my_step..."
    }
  ]
  type = "approval-chain"
}