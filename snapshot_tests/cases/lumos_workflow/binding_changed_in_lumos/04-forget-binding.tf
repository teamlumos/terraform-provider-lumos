# The binding is forgotten so the case can be destroyed.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Binding Changed In Lumos"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_workflow" "day" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Wait a day"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P1D" })
  }]
}

resource "lumos_workflow" "week" {
  app_id = lumos_app.this.id
  type   = "account-provisioning"
  name   = "bb-${var.run_id} Wait a week"
  steps = [{
    step   = "wait"
    inputs = jsonencode({ duration = "P7D" })
  }]
}

removed {
  from = lumos_app_provisioning_workflow.this

  lifecycle {
    destroy = false
  }
}
