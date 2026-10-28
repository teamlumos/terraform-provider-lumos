# Two entitlements with the same label on one app, in one apply.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Entitlement Duplicates"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "first" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Same Label"
}

resource "lumos_entitlement" "second" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Same Label"
}
