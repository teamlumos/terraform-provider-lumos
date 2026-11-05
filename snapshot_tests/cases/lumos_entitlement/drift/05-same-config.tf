# Same config as step 01: the plan creates the entitlement again, and
# changes nothing else.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Drift"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}
