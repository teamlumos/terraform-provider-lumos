# An app and an entitlement for the invalid policies that follow.
resource "lumos_app" "this" {
  label_override       = "bb-${var.run_id} Invalid Policies"
  category_override    = "Developers"
  description_override = "Snapshot test custom app"
}

resource "lumos_entitlement" "this" {
  app_id = lumos_app.this.id
  label  = "bb-${var.run_id} Read"
}
