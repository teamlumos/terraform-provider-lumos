# expect: error
# An entitlement-deprovisioning workflow with no steps is refused.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Empty Entitlement Deprovisioning"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

data "lumos_current_user" "me" {}

resource "lumos_workflow" "entitlement_deprovisioning" {
  app_id = lumos_app.this.id
  type   = "entitlement-deprovisioning"
  name   = "bb-${var.run_id} Empty entitlement-deprovisioning"
  steps  = []
}
